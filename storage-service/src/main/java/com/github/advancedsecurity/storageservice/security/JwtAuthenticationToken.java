package com.github.advancedsecurity.storageservice.security;

import java.util.Collection;
import java.util.ArrayList;

import org.springframework.security.authentication.AbstractAuthenticationToken;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.core.GrantedAuthority;

import com.github.advancedsecurity.storageservice.models.Profile;

import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.Jws;
import io.jsonwebtoken.Claims;

public class JwtAuthenticationToken extends AbstractAuthenticationToken {
    private Claims claims;

    public JwtAuthenticationToken(Jws<Claims> jws) {
        super(new ArrayList<GrantedAuthority>());
        this.claims = jws.getBody();
    }

    @Override
    public String getName() {
        return claims.get("profile", Profile.class).login;
    }

    @Override
    public Object getPrincipal() {
        return claims.get("profile", Profile.class);
    }

    @Override
    public Object getCredentials() {
        return claims.get("profile", Profile.class);
    }

    @Override
    public Object getDetails() {
        return claims.get("profile", Profile.class);
    }

    @Override
    public boolean isAuthenticated() {
        return true;
    }

    @Override
    public Collection<GrantedAuthority> getAuthorities() {
        return new ArrayList<GrantedAuthority>();
    }

}