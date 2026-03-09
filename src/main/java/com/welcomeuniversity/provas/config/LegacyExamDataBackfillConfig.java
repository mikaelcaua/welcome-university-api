package com.welcomeuniversity.provas.config;

import java.time.Instant;

import org.springframework.boot.CommandLineRunner;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import com.welcomeuniversity.provas.model.Exam;
import com.welcomeuniversity.provas.model.ExamStatus;
import com.welcomeuniversity.provas.repository.ExamRepository;
import com.welcomeuniversity.provas.service.S3Service;

@Configuration
public class LegacyExamDataBackfillConfig {

    @Bean
    CommandLineRunner backfillLegacyExamMetadata(ExamRepository examRepository, S3Service s3Service) {
        return args -> {
            for (Exam exam : examRepository.findAll()) {
                boolean dirty = false;
                Instant now = Instant.now();

                if (exam.getStatus() == null) {
                    exam.setStatus(ExamStatus.APPROVED);
                    dirty = true;
                }

                if (exam.getCreatedAt() == null) {
                    exam.setCreatedAt(now);
                    dirty = true;
                }

                if (exam.getUpdatedAt() == null) {
                    exam.setUpdatedAt(exam.getCreatedAt() != null ? exam.getCreatedAt() : now);
                    dirty = true;
                }

                if (exam.getStorageKey() != null && !exam.getStorageKey().isBlank()) {
                    String updatedUrl = s3Service.buildPublicUrlFromStorageKey(exam.getStorageKey());
                    if (updatedUrl != null && !updatedUrl.equals(exam.getPdfUrl())) {
                        exam.setPdfUrl(updatedUrl);
                        dirty = true;
                    }
                }

                if (dirty) {
                    examRepository.save(exam);
                }
            }
        };
    }
}
