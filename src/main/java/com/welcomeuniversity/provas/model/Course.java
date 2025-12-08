package com.welcomeuniversity.provas.model;

import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.ManyToOne;

@Entity
public class Course {
    @Id @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    private String name;

    @ManyToOne
    private University university;

    public Course() {}
    public Course(String name, University university){ this.name = name; this.university = university; }

    public Long getId(){ return id; }
    public String getName(){ return name; }
    public void setName(String name){ this.name = name; }
    public University getUniversity(){ return university; }
    public void setUniversity(University university){ this.university = university; }
}
