package com.welcomeuniversity.provas.dto.exam;

import org.springframework.web.multipart.MultipartFile;

import com.welcomeuniversity.provas.model.ExamType;

import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotNull;

public class ExamUploadRequest {

    @NotNull
    @Min(2000)
    private Integer examYear;

    @NotNull
    @Min(1)
    @Max(2)
    private Integer semester;

    @NotNull
    private ExamType type;

    @NotNull
    private Long subjectId;

    @NotNull
    private MultipartFile file;

    private Boolean periodUnidentified = false;

    public Integer getExamYear() {
        return examYear;
    }

    public void setExamYear(Integer examYear) {
        this.examYear = examYear;
    }

    public Integer getSemester() {
        return semester;
    }

    public void setSemester(Integer semester) {
        this.semester = semester;
    }

    public ExamType getType() {
        return type;
    }

    public void setType(ExamType type) {
        this.type = type;
    }

    public Long getSubjectId() {
        return subjectId;
    }

    public void setSubjectId(Long subjectId) {
        this.subjectId = subjectId;
    }

    public MultipartFile getFile() {
        return file;
    }

    public void setFile(MultipartFile file) {
        this.file = file;
    }

    public Boolean getPeriodUnidentified() {
        return periodUnidentified;
    }

    public void setPeriodUnidentified(Boolean periodUnidentified) {
        this.periodUnidentified = periodUnidentified;
    }
}
