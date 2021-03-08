package com.github.advancedsecurity.storageservice.security;

import java.io.IOException;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.servlet.FilterChain;
import javax.servlet.ServletException;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.web.authentication.AbstractAuthenticationProcessingFilter;
import org.springframework.security.oauth2.server.resource.web.BearerTokenResolver;
import org.springframework.security.oauth2.server.resource.web.DefaultBearerTokenResolver;
import org.springframework.security.oauth2.server.resource.BearerTokenAuthenticationToken;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.beans.factory.annotation.Autowired;

public class BearerAuthenticationFilter extends AbstractAuthenticationProcessingFilter {

    private final Logger logger = LoggerFactory.getLogger(this.getClass());
    private BearerTokenResolver bearerTokenResolver = new DefaultBearerTokenResolver();

    public BearerAuthenticationFilter(AuthenticationManager authenticationManager, String defaultFilterProcessesUrl) {
        super(defaultFilterProcessesUrl);

        setAuthenticationManager(authenticationManager);
    }

    @Override
    public Authentication attemptAuthentication(HttpServletRequest request, HttpServletResponse response) 
        throws AuthenticationException, IOException, ServletException 
    {
        logger.debug("Bearer authentication attempt.");
        String token;
		try {
            logger.debug("Resolving bearer token.");
			token = this.bearerTokenResolver.resolve(request);
		}
		catch (OAuth2AuthenticationException invalid) {
			throw new JwtAuthenticationException("Invalid Bearer token!");
		}

        if (token == null) {
            throw new JwtAuthenticationException("Invalid Bearer token!"); 
        }

        logger.debug("Constructing Bearer authentication token.");
        BearerTokenAuthenticationToken authenticationRequest = new BearerTokenAuthenticationToken(token);

        logger.debug("Authenticating with Bearer authentication token.");
        return getAuthenticationManager().authenticate(authenticationRequest);
    }

    @Override
    protected void successfulAuthentication(HttpServletRequest request, HttpServletResponse response, FilterChain chain, Authentication authResult)
            throws IOException, ServletException {
        super.successfulAuthentication(request, response, chain, authResult);

        // As this authentication is in HTTP header, after success we need to continue the request normally
        // and return the response as if the resource was not secured at all
        chain.doFilter(request, response);
    }
}