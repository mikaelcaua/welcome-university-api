package com.welcomeuniversity.provas.config;

import java.util.Arrays;
import java.util.List;

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

            List<State> states = Arrays.asList(
                new State("AC", "Acre"),
                new State("AL", "Alagoas"),
                new State("AP", "Amapá"),
                new State("AM", "Amazonas"),
                new State("BA", "Bahia"),
                new State("CE", "Ceará"),
                new State("DF", "Distrito Federal"),
                new State("ES", "Espírito Santo"),
                new State("GO", "Goiás"),
                new State("MA", "Maranhão"),
                new State("MT", "Mato Grosso"),
                new State("MS", "Mato Grosso do Sul"),
                new State("MG", "Minas Gerais"),
                new State("PA", "Pará"),
                new State("PB", "Paraíba"),
                new State("PR", "Paraná"),
                new State("PE", "Pernambuco"),
                new State("PI", "Piauí"),
                new State("RJ", "Rio de Janeiro"),
                new State("RN", "Rio Grande do Norte"),
                new State("RS", "Rio Grande do Sul"),
                new State("RO", "Rondônia"),
                new State("RR", "Roraima"),
                new State("SC", "Santa Catarina"),
                new State("SP", "São Paulo"),
                new State("SE", "Sergipe"),
                new State("TO", "Tocantins")
            );

            stateRepo.saveAll(states);

            State ma = states.stream().filter(s -> s.getCode().equals("MA")).findFirst().get();
            State sp = states.stream().filter(s -> s.getCode().equals("SP")).findFirst().get();

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