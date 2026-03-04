package com.welcomeuniversity.provas;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.ConfigurationPropertiesScan;

@SpringBootApplication
@ConfigurationPropertiesScan
public class ProvasApiApplication {
    public static void main(String[] args) {
        SpringApplication.run(ProvasApiApplication.class, args);
    }
}
