package com.github.advancedsecurity.storageservice.models;

import java.io.File;
import java.util.Base64;
import org.apache.commons.io.FileUtils;

public class Blob {
    private String mimeType;
    private String base64EncodedData;

    public Blob(String mimeType, byte[] data) {
        this.mimeType = mimeType;
        this.base64EncodedData = Base64.getEncoder().encodeToString(data);
    }

    public Blob(File file) throws java.net.MalformedURLException, java.io.IOException {
        this.mimeType = file.toURL().openConnection().getContentType();
        this.base64EncodedData = Base64.getEncoder().encodeToString(FileUtils.readFileToByteArray(file));
    }

    public String getMimeType() {
        return this.mimeType;
    }

    public String getData() {
        return this.base64EncodedData;
    }
}