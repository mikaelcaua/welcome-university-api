package examusecase

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/subject"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/storage"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/erros"
)

const maxPendingExamsForRegularUser = 5
const unidentifiedPeriodLabel = "PERIODO_NAO_IDENTIFICADO"

var supportedUploadExtensions = map[string]bool{"pdf": true, "png": true, "jpg": true, "jpeg": true, "webp": true, "gif": true}

type UploadPayloadOptimizer interface {
	Optimize(payload storagecontract.UploadPayload) storagecontract.UploadPayload
}

type UploadInput struct {
	ExamYear           int
	Semester           int
	Type               exam.ExamType
	SubjectID          int64
	PeriodUnidentified bool
	OriginalFilename   string
	ContentType        string
	FileBytes          []byte
}

type ReviewInput struct {
	Status     exam.ExamStatus
	ReviewNote *string
}

func parsePeriod(period string) (int, int, error) {
	parts := strings.Split(period, ".")
	if len(parts) != 2 {
		return 0, 0, domainerrors.Validation("Period deve seguir o formato AAAA.S.")
	}
	examYear, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, domainerrors.Validation("Period deve conter apenas numeros.")
	}
	semester, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, domainerrors.Validation("Period deve conter apenas numeros.")
	}
	if semester < 1 || semester > 2 {
		return 0, 0, domainerrors.Validation("Semestre invalido.")
	}
	return examYear, semester, nil
}

func validateUploadInput(input UploadInput) error {
	if input.ExamYear < 2000 {
		return domainerrors.Validation("Ano da prova invalido.")
	}
	if input.Semester < 1 || input.Semester > 2 {
		return domainerrors.Validation("Semestre invalido.")
	}
	if !input.Type.IsValid() {
		return domainerrors.Validation("Tipo de prova invalido.")
	}
	if len(input.FileBytes) == 0 {
		return domainerrors.Validation("Arquivo obrigatorio.")
	}
	lowerFilename := strings.ToLower(input.OriginalFilename)
	lowerContentType := strings.ToLower(input.ContentType)
	extension := strings.TrimPrefix(strings.ToLower(strings.TrimPrefix(fileExtension(lowerFilename), ".")), ".")
	validFile := lowerContentType == "application/pdf" || strings.HasPrefix(lowerContentType, "image/") || supportedUploadExtensions[extension]
	if !validFile {
		return domainerrors.Validation("Somente arquivos PDF ou imagens sao aceitos.")
	}
	return nil
}

func buildExamName(subject subject.Subject, input UploadInput) string {
	if input.PeriodUnidentified {
		return fmt.Sprintf("%s - %s - %s", strings.TrimSpace(subject.Name), unidentifiedPeriodLabel, formatExamType(input.Type))
	}
	return fmt.Sprintf("%s - %d.%d - %s", strings.TrimSpace(subject.Name), input.ExamYear, input.Semester, formatExamType(input.Type))
}

func formatExamType(examType exam.ExamType) string { return strings.ToUpper(strings.ReplaceAll(string(examType), "_", " ")) }
func fileExtension(filename string) string {
	lastDotIndex := strings.LastIndex(filename, ".")
	if lastDotIndex < 0 || lastDotIndex == len(filename)-1 {
		return ""
	}
	return filename[lastDotIndex+1:]
}
func normalizeOptionalText(value *string) *string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil
	}
	normalized := strings.TrimSpace(*value)
	return &normalized
}
func sha256Hex(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}
