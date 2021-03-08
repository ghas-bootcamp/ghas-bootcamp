package com.github.advancedsecurity.storageservice.controllers;

import java.util.UUID;
import java.util.Arrays;
import java.io.InputStream;
import java.io.ObjectInputStream;
import java.io.OutputStream;
import java.io.ObjectOutputStream;
import java.io.IOException;
import java.io.ByteArrayInputStream;

import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.core.io.Resource;
import org.springframework.core.io.WritableResource;
import org.springframework.core.io.ResourceLoader;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.http.HttpStatus;
import org.springframework.web.multipart.MultipartFile;

import com.github.advancedsecurity.storageservice.security.JwtAuthenticationToken;
import com.github.advancedsecurity.storageservice.models.Blob;
import com.github.advancedsecurity.storageservice.models.Profile;

@RestController
@CrossOrigin
public class BlobController {
    @Autowired
    private ResourceLoader resourceLoader;

    @Value("${cloud.aws.s3.bucket}")
    private String bucket;

    @Value("${blob.allowed-content-types}")
    private String[] allowedContentTypes;

    private Blob deserializeBlob(InputStream inputStream) throws IOException, ClassNotFoundException{
        ObjectInputStream in = new ObjectInputStream(inputStream);
        Blob b = (Blob)in.readObject();
        in.close();
        return b;
    }

    @GetMapping("/blob/{id}")
    public Blob get(@PathVariable UUID id, JwtAuthenticationToken token) {
        Profile profile = (Profile)token.getPrincipal();
        try {
            Blob blob = null;
            Resource resource = this.resourceLoader.getResource(String.format("s3://%s/%s/%s", bucket, profile.name, id.toString()));
            InputStream inputStream = resource.getInputStream();
            blob = deserializeBlob(inputStream);
            inputStream.close();
            return blob;
        } catch (IOException i) {
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "Failed to retrieve blob!");
        } catch (ClassNotFoundException c) {
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "Failed to retrieve blob!");
        }
    }

    @PostMapping("/blob")
    public UUID post(@RequestParam MultipartFile file, @RequestParam(defaultValue="false") boolean isBlob, JwtAuthenticationToken token) {
        Profile profile = (Profile)token.getPrincipal();
        try {
            UUID id = UUID.randomUUID();
            Blob blob = null;
            if (isBlob) {
                blob = deserializeBlob(new ByteArrayInputStream(file.getBytes()));
            } else {
                blob = new Blob(file.getContentType(), file.getBytes());
            }

            if(!Arrays.asList(this.allowedContentTypes).contains(blob.getMimeType())) {
                throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Disallowed content type");
            }
            
            Resource resource = this.resourceLoader.getResource(String.format("s3://%s/%s/%s", bucket, profile.name, id.toString()));
            WritableResource writableResource = (WritableResource) resource;
            OutputStream outputStream = writableResource.getOutputStream();
            ObjectOutputStream out = new ObjectOutputStream(outputStream);
            out.writeObject(blob);
            out.close();
            outputStream.close();
            return id;
        } catch (IOException i) {
           throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "Failed to store blob!");
        } catch (ClassNotFoundException c) {
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "Failed to retrieve blob!");
        }
    }
}