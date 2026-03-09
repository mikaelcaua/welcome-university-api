package com.welcomeuniversity.provas.security;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.config.annotation.authentication.configuration.AuthenticationConfiguration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;

@Configuration
public class SecurityConfig {

    @Bean
    SecurityFilterChain securityFilterChain(HttpSecurity http, JwtAuthenticationFilter jwtAuthenticationFilter)
        throws Exception {
        http
            .csrf(AbstractHttpConfigurer::disable)
            .formLogin(AbstractHttpConfigurer::disable)
            .httpBasic(AbstractHttpConfigurer::disable)
            .sessionManagement(session -> session.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
            .authorizeHttpRequests(auth -> auth
                .requestMatchers(
                    "/v3/api-docs/**",
                    "/swagger-ui/**",
                    "/swagger-ui.html",
                    "/actuator/health/**",
                    "/actuator/info"
                ).permitAll()
                .requestMatchers(HttpMethod.POST, "/auth/register", "/auth/login", "/auth/refresh").permitAll()
                .requestMatchers(HttpMethod.GET, "/exams/pending").hasAnyRole("APPROVER", "ADMIN", "DEV")
                .requestMatchers(HttpMethod.GET, "/states/**", "/universities/**", "/courses/**", "/subjects/**", "/exams")
                    .permitAll()
                .requestMatchers(HttpMethod.POST, "/states").hasAnyRole("ADMIN", "DEV")
                .requestMatchers(HttpMethod.POST, "/states/*/universities").hasAnyRole("ADMIN", "DEV")
                .requestMatchers(HttpMethod.POST, "/universities/*/courses").hasAnyRole("ADMIN", "DEV")
                .requestMatchers(HttpMethod.POST, "/courses/*/subjects").hasAnyRole("ADMIN", "DEV")
                .requestMatchers(HttpMethod.POST, "/exams").hasAnyRole("USER", "APPROVER", "ADMIN", "DEV")
                .requestMatchers(HttpMethod.PATCH, "/exams/*/status").hasAnyRole("APPROVER", "ADMIN", "DEV")
                .requestMatchers(HttpMethod.GET, "/users/me").hasAnyRole("USER", "APPROVER", "ADMIN", "DEV")
                .requestMatchers(HttpMethod.GET, "/users").hasAnyRole("ADMIN", "DEV")
                .requestMatchers(HttpMethod.PATCH, "/users/*/role").hasAnyRole("ADMIN", "DEV")
                .anyRequest().authenticated()
            )
            .addFilterBefore(jwtAuthenticationFilter, UsernamePasswordAuthenticationFilter.class);

        return http.build();
    }

    @Bean
    PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }

    @Bean
    AuthenticationManager authenticationManager(AuthenticationConfiguration authenticationConfiguration)
        throws Exception {
        return authenticationConfiguration.getAuthenticationManager();
    }
}
