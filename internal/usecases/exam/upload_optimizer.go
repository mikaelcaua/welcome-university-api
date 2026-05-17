package examusecase

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/storage"
)

type UploadOptimizer struct{}

func NewUploadOptimizer() *UploadOptimizer { return &UploadOptimizer{} }

func (optimizer *UploadOptimizer) Optimize(payload storagecontract.UploadPayload) storagecontract.UploadPayload {
	lowerFilename := strings.ToLower(payload.OriginalFilename)
	lowerContentType := strings.ToLower(payload.ContentType)
	if lowerContentType == "image/jpeg" || lowerContentType == "image/jpg" || strings.HasSuffix(lowerFilename, ".jpg") || strings.HasSuffix(lowerFilename, ".jpeg") {
		return optimizer.optimizeJPEG(payload)
	}
	if lowerContentType == "image/png" || strings.HasSuffix(lowerFilename, ".png") {
		return optimizer.optimizePNG(payload)
	}
	return payload
}

func (optimizer *UploadOptimizer) optimizeJPEG(payload storagecontract.UploadPayload) storagecontract.UploadPayload {
	sourceImage, _, err := image.Decode(bytes.NewReader(payload.Bytes))
	if err != nil {
		return payload
	}
	var output bytes.Buffer
	if err := jpeg.Encode(&output, sourceImage, &jpeg.Options{Quality: 78}); err != nil || output.Len() >= len(payload.Bytes) {
		return payload
	}
	payload.Bytes = output.Bytes()
	payload.ContentType = "image/jpeg"
	return payload
}

func (optimizer *UploadOptimizer) optimizePNG(payload storagecontract.UploadPayload) storagecontract.UploadPayload {
	sourceImage, _, err := image.Decode(bytes.NewReader(payload.Bytes))
	if err != nil {
		return payload
	}
	var output bytes.Buffer
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	if err := encoder.Encode(&output, sourceImage); err != nil || output.Len() >= len(payload.Bytes) {
		return payload
	}
	payload.Bytes = output.Bytes()
	payload.ContentType = "image/png"
	return payload
}
