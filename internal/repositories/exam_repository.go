package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mikaelcaua/welcome-university-api/internal/models"
)

type ApprovedExamFilter struct {
	SubjectID *int64
	ExamYear  *int
	Semester  *int
}

type PostgresExamRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresExamRepository(pool *pgxpool.Pool) *PostgresExamRepository {
	return &PostgresExamRepository{pool: pool}
}

func (repository *PostgresExamRepository) ListApproved(ctx context.Context, filter ApprovedExamFilter) ([]models.Exam, error) {
	query := baseExamQuery() + ` WHERE e.status = 'APPROVED'`
	args := []any{}
	if filter.SubjectID != nil {
		args = append(args, *filter.SubjectID)
		query += ` AND e.subject_id = $1`
	}
	if filter.ExamYear != nil && filter.Semester != nil {
		args = append(args, *filter.ExamYear, *filter.Semester)
		query += ` AND e.exam_year = $2 AND e.semester = $3`
	}
	query += ` ORDER BY e.created_at DESC`
	return repository.list(ctx, query, args...)
}

func (repository *PostgresExamRepository) ListPendingBySubject(ctx context.Context, subjectID int64) ([]models.Exam, error) {
	return repository.list(ctx, baseExamQuery()+` WHERE e.status = 'PENDING' AND e.subject_id = $1 ORDER BY e.id ASC`, subjectID)
}

func (repository *PostgresExamRepository) ListPendingByUploader(ctx context.Context, userID int64) ([]models.Exam, error) {
	return repository.list(ctx, baseExamQuery()+` WHERE e.status = 'PENDING' AND e.uploaded_by_id = $1 ORDER BY e.id ASC`, userID)
}

func (repository *PostgresExamRepository) CountPendingByUploader(ctx context.Context, userID int64) (int64, error) {
	var total int64
	err := repository.pool.QueryRow(ctx, `SELECT count(*) FROM exam WHERE uploaded_by_id = $1 AND status = 'PENDING'`, userID).Scan(&total)
	return total, err
}

func (repository *PostgresExamRepository) FileHashExists(ctx context.Context, fileHash string) (bool, error) {
	return queryExists(repository.pool, ctx, `SELECT EXISTS(SELECT 1 FROM exam WHERE file_hash = $1)`, fileHash)
}

func (repository *PostgresExamRepository) Create(ctx context.Context, exam models.Exam) (models.Exam, error) {
	var subjectID any
	var uploadedByID any
	var reviewedByID any
	if exam.Subject != nil {
		subjectID = exam.Subject.ID
	}
	if exam.UploadedBy != nil {
		uploadedByID = exam.UploadedBy.ID
	}
	if exam.ReviewedBy != nil {
		reviewedByID = exam.ReviewedBy.ID
	}
	err := repository.pool.QueryRow(ctx, `
		INSERT INTO exam (
			name, exam_year, semester, period_label, pdf_url, storage_key, file_hash, type, status,
			subject_id, uploaded_by_id, reviewed_by_id, review_note, reviewed_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		RETURNING id, created_at, updated_at
	`,
		exam.Name, exam.ExamYear, exam.Semester, exam.PeriodLabel, exam.PDFURL, exam.StorageKey,
		exam.FileHash, exam.Type, exam.Status, subjectID, uploadedByID, reviewedByID, exam.ReviewNote, exam.ReviewedAt,
	).Scan(&exam.ID, &exam.CreatedAt, &exam.UpdatedAt)
	return exam, err
}

func (repository *PostgresExamRepository) FindByID(ctx context.Context, examID int64) (models.Exam, bool, error) {
	exams, err := repository.list(ctx, baseExamQuery()+` WHERE e.id = $1`, examID)
	if err != nil {
		return models.Exam{}, false, err
	}
	if len(exams) == 0 {
		return models.Exam{}, false, nil
	}
	return exams[0], true, nil
}

func (repository *PostgresExamRepository) UpdateReview(ctx context.Context, exam models.Exam) (models.Exam, error) {
	var reviewerID any
	if exam.ReviewedBy != nil {
		reviewerID = exam.ReviewedBy.ID
	}
	err := repository.pool.QueryRow(ctx, `
		UPDATE exam
		SET status = $2,
			reviewed_by_id = $3,
			reviewed_at = $4,
			review_note = $5,
			pdf_url = $6,
			storage_key = $7,
			updated_at = now()
		WHERE id = $1
		RETURNING updated_at
	`, exam.ID, exam.Status, reviewerID, exam.ReviewedAt, exam.ReviewNote, exam.PDFURL, exam.StorageKey).Scan(&exam.UpdatedAt)
	return exam, err
}

func (repository *PostgresExamRepository) list(ctx context.Context, query string, args ...any) ([]models.Exam, error) {
	rows, err := repository.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exams := []models.Exam{}
	for rows.Next() {
		exam, err := scanExam(rows)
		if err != nil {
			return nil, err
		}
		exams = append(exams, exam)
	}
	if errors.Is(rows.Err(), pgx.ErrNoRows) {
		return exams, nil
	}
	return exams, rows.Err()
}

func baseExamQuery() string {
	return `
		SELECT
			e.id, e.name, e.exam_year, e.semester, e.period_label, e.type, e.pdf_url, e.storage_key,
			e.file_hash, e.status, e.review_note, e.created_at, e.reviewed_at, e.updated_at,
			s.id, s.name,
			up.id, up.name, up.email, up.role, up.created_at,
			rv.id, rv.name, rv.email, rv.role, rv.created_at
		FROM exam e
		LEFT JOIN subject s ON s.id = e.subject_id
		LEFT JOIN app_user up ON up.id = e.uploaded_by_id
		LEFT JOIN app_user rv ON rv.id = e.reviewed_by_id
	`
}

func scanExam(rows pgx.Rows) (models.Exam, error) {
	var exam models.Exam
	var subjectID sql.NullInt64
	var subjectName sql.NullString
	var uploaderID sql.NullInt64
	var uploaderName sql.NullString
	var uploaderEmail sql.NullString
	var uploaderRole sql.NullString
	var uploaderCreatedAt sql.NullTime
	var reviewerID sql.NullInt64
	var reviewerName sql.NullString
	var reviewerEmail sql.NullString
	var reviewerRole sql.NullString
	var reviewerCreatedAt sql.NullTime
	var reviewedAt sql.NullTime
	var periodLabel sql.NullString
	var storageKey sql.NullString
	var fileHash sql.NullString
	var reviewNote sql.NullString
	var createdAt time.Time
	var updatedAt time.Time

	err := rows.Scan(
		&exam.ID, &exam.Name, &exam.ExamYear, &exam.Semester, &periodLabel, &exam.Type, &exam.PDFURL,
		&storageKey, &fileHash, &exam.Status, &reviewNote, &createdAt, &reviewedAt, &updatedAt,
		&subjectID, &subjectName,
		&uploaderID, &uploaderName, &uploaderEmail, &uploaderRole, &uploaderCreatedAt,
		&reviewerID, &reviewerName, &reviewerEmail, &reviewerRole, &reviewerCreatedAt,
	)
	if err != nil {
		return models.Exam{}, err
	}

	exam.CreatedAt = createdAt
	exam.UpdatedAt = updatedAt
	exam.PeriodLabel = nullStringPointer(periodLabel)
	exam.StorageKey = nullStringPointer(storageKey)
	exam.FileHash = nullStringPointer(fileHash)
	exam.ReviewNote = nullStringPointer(reviewNote)
	if reviewedAt.Valid {
		exam.ReviewedAt = &reviewedAt.Time
	}
	if subjectID.Valid {
		exam.Subject = &models.Subject{ID: subjectID.Int64, Name: subjectName.String}
	}
	if uploaderID.Valid {
		exam.UploadedBy = &models.User{
			ID:        uploaderID.Int64,
			Name:      uploaderName.String,
			Email:     uploaderEmail.String,
			Role:      models.Role(uploaderRole.String),
			CreatedAt: uploaderCreatedAt.Time,
		}
	}
	if reviewerID.Valid {
		exam.ReviewedBy = &models.User{
			ID:        reviewerID.Int64,
			Name:      reviewerName.String,
			Email:     reviewerEmail.String,
			Role:      models.Role(reviewerRole.String),
			CreatedAt: reviewerCreatedAt.Time,
		}
	}
	return exam, nil
}

func nullStringPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}
