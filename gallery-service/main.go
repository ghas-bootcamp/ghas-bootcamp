package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"errors"
	"database/sql"
	"strconv"
	"os"

	"github.com/gorilla/mux"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"

)

type Configuration struct {
	Host 		string	`json:"host"`
	Port 		int		`json:"port"`

	Database	string	`json:"database"`
	Secret 		string	`json:"secret"`
	AllowedOrigins []string `json:"allowed-origins"`
}

var configuration *Configuration

func LoadConfiguration(path string) (*Configuration, error) {
	configuration := &Configuration{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(configuration)
	if err != nil {
		return nil, err
	}

	return configuration, nil
}

var db *sql.DB

func GetDb() *sql.DB {
	if db != nil {
		return db
	}

	db, err := sql.Open("sqlite3", configuration.Database)
	if err != nil {
		panic(err)
	}

	return db
}

func InitializeDb() {
	db := GetDb()

	tables := [...]string {
		`CREATE TABLE IF NOT EXISTS 'gallery'
		(id INTEGER PRIMARY KEY AUTOINCREMENT, login TEXT UNIQUE, title TEXT NOT NULL,
			description TEXT, created_at DATETIME, updated_at DATETIME);`,
		`CREATE TRIGGER IF NOT EXISTS 'gallery_after_insert' 
			AFTER INSERT ON 'gallery' 
			BEGIN 
				UPDATE 'gallery' 
				SET
					created_at = DATETIME('NOW'),
					updated_at = DATETIME('NOW');
			END`,
		`CREATE TRIGGER IF NOT EXISTS 'gallery_after_update' 
			AFTER UPDATE ON 'gallery' 
			BEGIN 
				UPDATE 'gallery' 
				SET
					updated_at = DATETIME('NOW');
			END`,
		`CREATE TABLE IF NOT EXISTS 'art_piece' 
			(id INTEGER PRIMARY KEY AUTOINCREMENT, gallery_id INTEGER, uri TEXT, title TEXT,
				description TEXT, stars INTEGER CHECK (stars BETWEEN 0 AND 3) DEFAULT 0, created_at DATETIME, updated_at DATETIME,
				FOREIGN KEY(gallery_id) REFERENCES gallery(id) ON DELETE CASCADE);`,
		`CREATE TRIGGER IF NOT EXISTS 'art_piece_after_insert' 
		AFTER INSERT ON 'art_piece' 
		BEGIN 
			UPDATE 'art_piece' 
			SET
				created_at = DATETIME('NOW'),
				updated_at = DATETIME('NOW');
		END`,
		`CREATE TRIGGER IF NOT EXISTS 'art_piece_after_update' 
			AFTER UPDATE ON 'art_piece' 
			BEGIN 
				UPDATE 'art_piece' 
				SET
					updated_at = DATETIME('NOW');
			END`,
	}

	for _, q := range tables {
		_, err := db.Exec(q)
		if err != nil {
			log.Printf("Failed query: %s", q)
			panic(err)
		}
	}
}

func Contains(needle string, haystack *[]string) bool {
	for _, candidate := range *haystack {
		if needle == candidate {
			return true
		}
	}

	return false
}

func SetCORSPolicy(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Headers", "authorization")
	if Contains(origin, &configuration.AllowedOrigins) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")	
}

type Gallery struct {
	ID 				int64 		`json:"id"`
	Title 			string		`json:"title"`
	Description		string		`json:"description"`
}

func (g *Gallery) Get(profile *OctoProfile) error {
	db := GetDb()

	stmt, err := db.Prepare(`SELECT id, title, description FROM gallery WHERE login = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rs, err := stmt.Query(profile.Login)
	if err != nil {
		return err
	}
	defer rs.Close()
	
	if rs.Next() {
		rs.Scan(&g.ID, &g.Title, &g.Description)

		return nil
	} else {
		g.ID = -1
		g.Title = "The Empty gallery"
		g.Description = "This gallery is in desparate need of some art!" 
		
		return g.Create(profile)
	}
}

func (g *Gallery) Create(profile *OctoProfile) error {
	db := GetDb()

	stmt, err := db.Prepare("INSERT INTO gallery (login, title, description) VALUES(?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	r , err := stmt.Exec(profile.Login, g.Title, g.Description)
	if err != nil {
		return err
	}

	if i, err := r.RowsAffected(); err != nil || i != 1 {
		return errors.New("Unable to create gallery")
	}

	newId, err := r.LastInsertId()
	if err != nil {
		return err
	}
	g.ID = newId

	return nil
}

func (g Gallery) Update(profile *OctoProfile) error {
	db := GetDb()

	stmt, err := db.Prepare(fmt.Sprintf("UPDATE gallery SET title = '%s', description = '%s' WHERE id = %d and login = '%s'", g.Title, g.Description, g.ID, profile.Login))
	if err != nil {
		return err
	}
	defer stmt.Close()

	r , err := stmt.Exec()
	if err != nil {
		return err
	}

	if i, err := r.RowsAffected(); err != nil || i != 1 {
		return errors.New("Unable to update gallery")
	}

	return nil
}

func (g Gallery) GetArtPiece(id int64) (*ArtPiece, error) {
	db := GetDb()

	stmt, err := db.Prepare(`SELECT id, title, description, stars, uri FROM art_piece WHERE gallery_id = ? AND id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rs, err := stmt.Query(g.ID, id)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	
	if rs.Next() {
		art_piece := &ArtPiece{}
		rs.Scan(&art_piece.ID, &art_piece.Title, &art_piece.Description, &art_piece.Stars, &art_piece.Uri)

		return art_piece, nil
	}

	return nil, errors.New("Art piece seems to be missing!")
}

func (g Gallery) GetAllArtPieces() ([]ArtPiece, error)  {
	db := GetDb()

	stmt, err := db.Prepare(`SELECT id, title, description, stars, uri FROM art_piece WHERE gallery_id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rs, err := stmt.Query(g.ID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	art_pieces := []ArtPiece{}
	
	for rs.Next() {
		art_piece := &ArtPiece{}
		rs.Scan(&art_piece.ID, &art_piece.Title, &art_piece.Description, &art_piece.Stars, &art_piece.Uri)

		art_pieces = append(art_pieces, *art_piece)
	}

	return art_pieces, nil
}

type ArtPiece struct {
	ID				int64		`json:"id"`
	Title			string		`json:"title"`
	Description		string		`json:"description"`
	Stars			int			`json:"stars"`
	Uri				string		`json:"uri"`
}

func (p *ArtPiece) Create(gallery Gallery) error {
	db := GetDb()

	stmt, err := db.Prepare("INSERT INTO art_piece (title, description, uri, gallery_id) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	r , err := stmt.Exec(p.Title, p.Description, p.Uri, gallery.ID)
	if err != nil {
		return err
	}

	if i, err := r.RowsAffected(); err != nil || i != 1 {
		return errors.New("Unable to create art piece")
	}

	newId, err := r.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = newId

	return nil
}

func (p ArtPiece) Update(gallery Gallery) error {
	db := GetDb()

	stmt, err := db.Prepare("UPDATE art_piece SET title = ?, description = ?, stars = ?, uri = ? WHERE id = ? and gallery_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	r , err := stmt.Exec(p.Title, p.Description, p.Stars, p.Uri, p.ID, gallery.ID)
	if err != nil {
		return err
	}

	if i, err := r.RowsAffected(); err != nil || i != 1 {
		return errors.New("Unable to update art piece")
	}

	return nil
}

func (p ArtPiece) Delete(gallery Gallery) error {
	db := GetDb()

	stmt, err := db.Prepare("DELETE FROM art_piece WHERE id = ? and gallery_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	r , err := stmt.Exec(p.ID, gallery.ID)
	if err != nil {
		return err
	}

	if i, err := r.RowsAffected(); err != nil || i != 1 {
		return errors.New("Unable to delete art piece")
	}

	return nil
}

func GetProfile(r *http.Request) *OctoProfile {
	profile := new(OctoProfile)

	profile.Login = r.Header.Get(GitHubLoginHeader.String())
	profile.Name = r.Header.Get(GitHubNameHeader.String())
	profile.Email = r.Header.Get(GitHubEmailHeader.String())

	return profile
}

func GalleryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			GetGalleryHandler(w, r)
		case http.MethodPut:
			UpdateGalleryHandler(w, r)
		case http.MethodOptions:
			return
		default:
			http.Error(w, "Method not permitted", http.StatusBadRequest)
	}
}

func UpdateGalleryHandler(w http.ResponseWriter, r *http.Request) {
	profile := GetProfile(r)

	var gallery Gallery

	err := json.NewDecoder(r.Body).Decode(&gallery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = gallery.Update(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetGalleryHandler(w http.ResponseWriter, r *http.Request) {
	profile := GetProfile(r)
	var gallery Gallery

	err := gallery.Get(profile)

	if err == nil {
		json.NewEncoder(w).Encode(gallery)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GalleryArtsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			GetGalleryAllArtHandler(w, r)
		case http.MethodPost:
			AddGalleryArtHandler(w, r)
		case http.MethodOptions:
			return
		default:
			http.Error(w, "Method not permitted", http.StatusBadRequest)
	}
}

func GetGalleryAllArtHandler(w http.ResponseWriter, r *http.Request) {
	profile := GetProfile(r)
	var gallery Gallery

	err := gallery.Get(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	art_pieces, err := gallery.GetAllArtPieces()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(art_pieces)
}

func AddGalleryArtHandler(w http.ResponseWriter, r *http.Request) {
	profile := GetProfile(r)
	var gallery Gallery
	var art_piece ArtPiece

	err := gallery.Get(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&art_piece)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = art_piece.Create(gallery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GalleryArtHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			GetGalleryArtHandler(w, r)
		case http.MethodPut:
			UpdateGalleryArtHandler(w, r)
		case http.MethodDelete:
			DeleteGalleryArtHandler(w, r)
		case http.MethodOptions:
			return
		default:
			http.Error(w, "Method not permitted", http.StatusBadRequest)
	}
}

func GetGalleryArtHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	art_piece_id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile := GetProfile(r)
	var gallery Gallery

	err = gallery.Get(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	art_pieces, err := gallery.GetArtPiece(art_piece_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(art_pieces)
}

func UpdateGalleryArtHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	art_piece_id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile := GetProfile(r)
	var gallery Gallery

	err = gallery.Get(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	old_art_piece, err := gallery.GetArtPiece(art_piece_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var art_piece ArtPiece
	err = json.NewDecoder(r.Body).Decode(&art_piece)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	art_piece.ID = old_art_piece.ID

	err = art_piece.Update(gallery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteGalleryArtHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	art_piece_id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile := GetProfile(r)
	var gallery Gallery

	err = gallery.Get(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	art_piece, err := gallery.GetArtPiece(art_piece_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = art_piece.Delete(gallery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HomeLinkHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

type OctoProfile struct {
	Login			string		`json:"login"`
	Name			string		`json:"name"`
	Email			string		`json:"email"`
}

type OctoClaims struct {
	jwt.StandardClaims

	Profile			OctoProfile	`json:"profile"`
}

func (claims OctoClaims) Valid() error {

	log.Printf("Validating standard claims")
	if err := claims.StandardClaims.Valid(); err != nil {
		log.Printf("Failed standard claim validation with error %s", err)
		return err
	}

	vErr := new(jwt.ValidationError)

	log.Printf("Validating private claims")
	if claims.Issuer != "OctoGallery" {
		log.Printf("Invalid issuer: %s", claims.Issuer)
		vErr.Inner = errors.New("Invalid issuer!")
		vErr.Errors |= jwt.ValidationErrorIssuer
	}


	if vErr.Errors == 0 {
		log.Printf("Validated all claims without errors.")
		return nil
	}

	return vErr
}
type ProfileHeader int

const (
	GitHubLoginHeader ProfileHeader = iota
	GitHubNameHeader
	GitHubEmailHeader
)

func (p ProfileHeader) String() string {
	return [...]string{"X-GitHub-Login", "X-GitHub-Name", "X-GitHub-Email"}[p]
}

func authnMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetCORSPolicy(w, r)

		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		r.Header.Del(GitHubLoginHeader.String())
		r.Header.Del(GitHubNameHeader.String())
		r.Header.Del(GitHubEmailHeader.String())

		authz := r.Header.Get("Authorization")

		if strings.HasPrefix(authz, "Bearer") {
			tokenString := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer"))
			log.Printf("AuthN: Received bearer token %s", tokenString)
			token, err := jwt.ParseWithClaims(tokenString, &OctoClaims{}, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				//if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				//	return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				//}
			
				// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
				return []byte(configuration.Secret), nil
			})

			if err != nil {
				log.Printf("AuthN: Invalid token %s", err)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			} 
			
			
			if claims, ok := token.Claims.(*OctoClaims); ok && token.Valid {
				log.Printf("AuthN: Received valid token %s", authz)

				log.Printf("AuthN: Adding %s %s", GitHubLoginHeader, claims.Profile.Login)
				r.Header.Add(GitHubLoginHeader.String(), claims.Profile.Login)
				log.Printf("AuthN: Adding %s %s", GitHubNameHeader, claims.Profile.Name)
				r.Header.Add(GitHubNameHeader.String(), claims.Profile.Name)
				log.Printf("AuthN: Adding %s %s", GitHubEmailHeader, claims.Profile.Email)
				r.Header.Add(GitHubEmailHeader.String(), claims.Profile.Email)
				next.ServeHTTP(w, r)
				return
			} else {
				log.Printf("AuthN: Received token with invalid claims")
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}
		log.Printf("AuthN: Received invalid authorization value")
		http.Error(w, "Forbidden", http.StatusForbidden)

	})
}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage %s CONFIG_FILE\n", os.Args[0])
		return
	}
	
	c, err := LoadConfiguration(os.Args[1])
	if err != nil {
		panic(err)
	}
	configuration = c

	InitializeDb()

	router := mux.NewRouter().StrictSlash(true)
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(authnMiddleware)
	router.HandleFunc("/", HomeLinkHandler)
	router.HandleFunc("/gallery", GalleryHandler).Methods(http.MethodGet, http.MethodPut, http.MethodOptions)
	router.HandleFunc("/gallery/art", GalleryArtsHandler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc("/gallery/art/{id}", GalleryArtHandler).Methods(http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions)

	addr := fmt.Sprintf("%s:%d", configuration.Host, configuration.Port)
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}