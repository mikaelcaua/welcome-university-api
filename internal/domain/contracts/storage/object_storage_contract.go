package storagecontract

import "context"

type UploadPayload struct {
	OriginalFilename string
	ContentType      string
	Bytes            []byte
}

type StoredObject struct {
	Key string
	URL string
}

type ObjectStorageContract interface {
	UploadExam(ctx context.Context, payload UploadPayload, subjectID int64) (StoredObject, error)
	DeleteObjectByKey(ctx context.Context, storageKey *string) error
}
