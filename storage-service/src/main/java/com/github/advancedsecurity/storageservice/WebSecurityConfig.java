package com.github.advancedsecurity.storageservice;

import java.util.Arrays;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configuration.WebSecurityConfigurerAdapter;
import org.springframework.security.config.annotation.web.configurers.CsrfConfigurer;
import org.springframework.web.cors.CorsConfiguration;
import org.springframework.web.cors.CorsConfigurationSource;
import org.springframework.web.cors.UrlBasedCorsConfigurationSource;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;
import org.springframework.security.config.annotation.authentication.builders.AuthenticationManagerBuilder;

import com.github.advancedsecurity.storageservice.security.BearerAuthenticationFilter;
import com.github.advancedsecurity.storageservice.security.JwtAuthenticationEntryPoint;
import com.github.advancedsecurity.storageservice.security.JwtAuthenticationProvider;
import com.github.advancedsecurity.storageservice.security.JwtAccessDeniedHandler;
import com.github.advancedsecurity.storageservice.security.JwtAuthenticationSuccessHandler;

@Configuration
@EnableWebSecurity
public class WebSecurityConfig extends WebSecurityConfigurerAdapter {
	private String antPattern = "/**";

	@Autowired
	private JwtAuthenticationProvider jwtAuthenticatinProvider;

	@Autowired
	private JwtAccessDeniedHandler jwtAccessDeniedHandler;

	@Autowired
	private JwtAuthenticationEntryPoint jwtAuthenticationEntryPoint;

	@Autowired
	private JwtAuthenticationSuccessHandler jwtAuthenticationSuccessHandler;

	@Override
	protected void configure(HttpSecurity http) throws Exception {
		BearerAuthenticationFilter filter = new BearerAuthenticationFilter(authenticationManager(), this.antPattern);
		filter.setAuthenticationSuccessHandler(jwtAuthenticationSuccessHandler);
		http.cors().and()
			 .csrf().disable()
			 .authorizeRequests().antMatchers(this.antPattern).authenticated().and()
			 .addFilterBefore(filter, UsernamePasswordAuthenticationFilter.class)
			 .sessionManagement().sessionCreationPolicy(SessionCreationPolicy.STATELESS);
	}

	@Override
  	public void configure(AuthenticationManagerBuilder auth) throws Exception {
		  auth.authenticationProvider(jwtAuthenticatinProvider);
  	}
	
	@Bean 
	public CorsConfigurationSource corsConfigurationSource() {
        UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource();
        source.registerCorsConfiguration(this.antPattern, new CorsConfiguration().applyPermitDefaultValues());
        return source;
	}
}