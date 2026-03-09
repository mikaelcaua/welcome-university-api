package com.welcomeuniversity.provas.model;

import java.time.Instant;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.PrePersist;
import jakarta.persistence.PreUpdate;

@Entity
public class Exam {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private String name;

    @Column(name = "exam_year")
    private int examYear;

    private int semester;

    private String periodLabel;

    @Column(nullable = false)
    private String pdfUrl;

    private String storageKey;

    @Column(unique = true, length = 64)
    private String fileHash;

    @Enumerated(EnumType.STRING)
    private ExamType type;

    @Enumerated(EnumType.STRING)
    private ExamStatus status;

    @ManyToOne
    private Subject subject;

    @ManyToOne
    private AppUser uploadedBy;

    @ManyToOne
    private AppUser reviewedBy;

    private String reviewNote;

    @Column(updatable = false)
    private Instant createdAt;

    private Instant reviewedAt;

    private Instant updatedAt;

    public Exam() {}

    public Exam(
        String name,
        int examYear,
        int semester,
        ExamType type,
        String pdfUrl,
        Subject subject
    ) {
        this.name = name;
        this.examYear = examYear;
        this.semester = semester;
        this.type = type;
        this.pdfUrl = pdfUrl;
        this.subject = subject;
        this.status = ExamStatus.APPROVED;
    }

    @PrePersist
    void onCreate() {
        Instant now = Instant.now();
        createdAt = now;
        updatedAt = now;
    }

    @PreUpdate
    void onUpdate() {
        updatedAt = Instant.now();
    }

    public Long getId() {
        return id;
    }

    public String getName() {
        return name;
    }
    public void setName(String name) {
        this.name = name;
    }

    public int getExamYear() {
        return examYear;
    }
    public void setExamYear(int examYear) {
        this.examYear = examYear;
    }

    public int getSemester() {
        return semester;
    }
    public void setSemester(int semester) {
        this.semester = semester;
    }

    public String getPeriodLabel() {
        return periodLabel;
    }
    public void setPeriodLabel(String periodLabel) {
        this.periodLabel = periodLabel;
    }

    public String getPdfUrl() {
        return pdfUrl;
    }
    public void setPdfUrl(String pdfUrl) {
        this.pdfUrl = pdfUrl;
    }

    public String getStorageKey() {
        return storageKey;
    }
    public void setStorageKey(String storageKey) {
        this.storageKey = storageKey;
    }

    public String getFileHash() {
        return fileHash;
    }
    public void setFileHash(String fileHash) {
        this.fileHash = fileHash;
    }

    public ExamType getType() {
        return type;
    }
    public void setType(ExamType type) {
        this.type = type;
    }

    public ExamStatus getStatus() {
        return status;
    }
    public void setStatus(ExamStatus status) {
        this.status = status;
    }

    public Subject getSubject() {
        return subject;
    }
    public void setSubject(Subject subject) {
        this.subject = subject;
    }

    public AppUser getUploadedBy() {
        return uploadedBy;
    }
    public void setUploadedBy(AppUser uploadedBy) {
        this.uploadedBy = uploadedBy;
    }

    public AppUser getReviewedBy() {
        return reviewedBy;
    }
    public void setReviewedBy(AppUser reviewedBy) {
        this.reviewedBy = reviewedBy;
    }

    public String getReviewNote() {
        return reviewNote;
    }
    public void setReviewNote(String reviewNote) {
        this.reviewNote = reviewNote;
    }

    public Instant getCreatedAt() {
        return createdAt;
    }
    public void setCreatedAt(Instant createdAt) {
        this.createdAt = createdAt;
    }

    public Instant getReviewedAt() {
        return reviewedAt;
    }
    public void setReviewedAt(Instant reviewedAt) {
        this.reviewedAt = reviewedAt;
    }

    public Instant getUpdatedAt() {
        return updatedAt;
    }
    public void setUpdatedAt(Instant updatedAt) {
        this.updatedAt = updatedAt;
    }
}
