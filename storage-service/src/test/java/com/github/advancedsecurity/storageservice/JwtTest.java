package com.github.advancedsecurity.storageservice;

import java.util.Date;
import org.junit.jupiter.api.Test;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.SignatureAlgorithm;

import static org.hamcrest.Matchers.is;
import static org.hamcrest.MatcherAssert.assertThat;

public class JwtTest {
    @Test
	public void testJWT() {
		String token = generateJwtToken();
		assertThat(token != null, is(true));
		System.out.println(token);
	}

	@SuppressWarnings("deprecation")
	private String generateJwtToken() {
		String token = Jwts.builder().setSubject("githubMona")
				.setExpiration(new Date(2929, 11, 25))
				.setIssuer("ghasuser@githubtest.com")
				.claim("groups", new String[] { "user", "admin" })
				// 48199327 repeated n times to satisfy HS256 base64 encoded
                .signWith(SignatureAlgorithm.HS256, "NDgxOTkzMjc0ODE5OTMyNzQ4MTk5MzI3NDgxOTkzMjc0ODE5OTMyNzQ4MTk5MzI3NDgxOTkzMjc0ODE5OTMyNzQ4MTk5MzI3NDgxOTkzMjc0ODE5OTMyNzQ4MTk5MzI3NDgxOTkzMjc0ODE5OTMyNzQ4MTk5MzI3NDgxOTkzMjc0ODE5OTMyNzQ4MTk5MzI3NDgxOTkzMjc0ODE5OTMyNzQ4MTk5MzI3NDgxOTkzMjc0ODE5OTMyNzQ4MTk5MzI3NDgxOTkzMjc=")
                .compact();
		return token;
	}
}
