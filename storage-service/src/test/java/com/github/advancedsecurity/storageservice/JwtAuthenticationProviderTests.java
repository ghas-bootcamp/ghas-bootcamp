package com.github.advancedsecurity.storageservice;

import org.springframework.security.oauth2.server.resource.BearerTokenAuthenticationToken;

import org.springframework.boot.test.context.SpringBootTest;
import com.github.advancedsecurity.storageservice.security.JwtAuthenticationProvider;

import org.junit.jupiter.api.Test;

import static org.hamcrest.Matchers.is;
import static org.hamcrest.MatcherAssert.assertThat;

@SpringBootTest
class JwtAuthenticationProviderTests {

    @Test
    void shouldSupportBearerTokenAuthenticationToken() throws Exception {
        JwtAuthenticationProvider provider = new JwtAuthenticationProvider();
        assertThat(provider.supports(BearerTokenAuthenticationToken.class), is(true));
    }
}
