package com.github.advancedsecurity.storageservice.security;

import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.authentication.AbstractAuthenticationToken;
import org.springframework.security.authentication.AuthenticationProvider;
import org.springframework.security.authentication.AuthenticationServiceException;
import org.springframework.security.oauth2.server.resource.BearerTokenAuthenticationToken;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.Jws;
import io.jsonwebtoken.Claims;
import io.jsonwebtoken.lang.Maps;
import io.jsonwebtoken.JwtException;
import io.jsonwebtoken.jackson.io.JacksonDeserializer;

import com.github.advancedsecurity.storageservice.models.Profile;

@Component
public class JwtAuthenticationProvider implements AuthenticationProvider {
    private final Logger logger = LoggerFactory.getLogger(this.getClass());

    @Value("${jwt.secret:secret}")
    private String secret;

    @Value("${jwt.issuer}")
    private String issuer;

    @Override
    public boolean supports(Class<?> authentication) {
        return BearerTokenAuthenticationToken.class.isAssignableFrom(authentication);
    }

    @Override
	public Authentication authenticate(Authentication authentication) throws AuthenticationException {
        logger.debug("Verifying key with secret '" + secret + "'");
        BearerTokenAuthenticationToken bearer = (BearerTokenAuthenticationToken) authentication;
		try {
            SecretKeySpec secretKey = new SecretKeySpec(this.secret.getBytes(StandardCharsets.UTF_8), "HmacSHA256");
		    
            Jws<Claims> jws = Jwts.parserBuilder()
            .deserializeJsonWith(new JacksonDeserializer(Maps.of("profile", Profile.class).build()))
            .requireIssuer(this.issuer)
            .setSigningKey(secretKey)
            .build()
            .parseClaimsJws(bearer.getToken());

            return new JwtAuthenticationToken(jws);
        } catch (JwtException ex) {
            logger.error("Failed to authenticate JWT token with error", ex);
            throw new JwtAuthenticationException("Invalid JWT token", ex);
        }
    }
}