package com.welcomeuniversity.provas.model;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.ManyToOne;

@Entity
public class Exam {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private String name;

    @Column(name = "exam_year")
    private int examYear;

    private int semester;

    private String pdfUrl;

    @Enumerated(EnumType.STRING)
    private ExamType type;

    @ManyToOne
    private Subject subject;

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

    public String getPdfUrl() {
        return pdfUrl;
    }
    public void setPdfUrl(String pdfUrl) {
        this.pdfUrl = pdfUrl;
    }

    public ExamType getType() {
        return type;
    }
    public void setType(ExamType type) {
        this.type = type;
    }

    public Subject getSubject() {
        return subject;
    }
    public void setSubject(Subject subject) {
        this.subject = subject;
    }
}
