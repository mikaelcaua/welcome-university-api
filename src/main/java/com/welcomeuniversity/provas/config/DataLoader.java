package com.welcomeuniversity.provas.config;

import org.springframework.boot.CommandLineRunner;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import com.welcomeuniversity.provas.model.Course;
import com.welcomeuniversity.provas.model.Exam;
import com.welcomeuniversity.provas.model.ExamType;
import com.welcomeuniversity.provas.model.State;
import com.welcomeuniversity.provas.model.Subject;
import com.welcomeuniversity.provas.model.University;
import com.welcomeuniversity.provas.repository.CourseRepository;
import com.welcomeuniversity.provas.repository.ExamRepository;
import com.welcomeuniversity.provas.repository.StateRepository;
import com.welcomeuniversity.provas.repository.SubjectRepository;
import com.welcomeuniversity.provas.repository.UniversityRepository;

@Configuration
public class DataLoader {

    @Bean
    CommandLineRunner init(StateRepository stateRepo,
                           UniversityRepository uniRepo,
                           CourseRepository courseRepo,
                           SubjectRepository subjectRepo,
                           ExamRepository examRepo) {
        return args -> {


            State ma = stateRepo.save(new State("MA", "Maranhão"));
            State sp = stateRepo.save(new State("SP", "São Paulo"));

            University ufma = uniRepo.save(new University("UFMA", ma));
            University usp = uniRepo.save(new University("USP", sp));

            Course ccUfma = courseRepo.save(new Course("Ciência da Computação", ufma));
            Course meusp = courseRepo.save(new Course("Engenharia", usp));

            Subject prog = subjectRepo.save(new Subject("Programação", ccUfma));
            Subject algebra = subjectRepo.save(new Subject("Álgebra Linear", ccUfma));
            Subject calculo = subjectRepo.save(new Subject("Cálculo I", meusp));

            createSemesterExams(examRepo, prog, 2025, 1);
            createSemesterExams(examRepo, prog, 2025, 2);
            createSemesterExams(examRepo, algebra, 2025, 1);
            createSemesterExams(examRepo, algebra, 2025, 2);
            createSemesterExams(examRepo, calculo, 2025, 1);
        };
    }

    private void createSemesterExams(ExamRepository examRepo, Subject subject, int year, int semester){
        examRepo.save(new Exam("Prova 1", year, semester, ExamType.PROVA1, "https://example.com/pdfs/"+subject.getName()+"-"+year+"."+semester+"-p1.pdf", subject));
        examRepo.save(new Exam("Prova 2", year, semester, ExamType.PROVA2, "https://example.com/pdfs/"+subject.getName()+"-"+year+"."+semester+"-p2.pdf", subject));
        examRepo.save(new Exam("Prova 3", year, semester, ExamType.PROVA3, "https://example.com/pdfs/"+subject.getName()+"-"+year+"."+semester+"-p3.pdf", subject));
        examRepo.save(new Exam("Recuperação", year, semester, ExamType.RECUPERACAO, "https://example.com/pdfs/"+subject.getName()+"-"+year+"."+semester+"-rec.pdf", subject));
        examRepo.save(new Exam("Final", year, semester, ExamType.FINAL, "https://example.com/pdfs/"+subject.getName()+"-"+year+"."+semester+"-final.pdf", subject));
    }
}
