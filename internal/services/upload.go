package services

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
)

const maxUploadSize = 10 << 20

var (
	ErrFileSizeExceedsLimit = errors.New("file size exceeds the allowable limit")
)

type Uploader interface {
	Upload(ctx context.Context, fileName string, fileSrc io.Reader) (string, error)
}

type UploadResponse struct {
	Url string `json:"url"`
}

func (s *ApiServer) uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch err := r.ParseMultipartForm(maxUploadSize); err {
	case nil:
		break
	case multipart.ErrMessageTooLarge:
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		if err := (&errResponse{Error: ErrFileSizeExceedsLimit.Error()}).writeResponse(w); err != nil {
			s.zapLogger.Sugar().Error(err)
		}
		return
	default:
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.zapLogger.Sugar().Infow(
		"handle upload file",
		"name", handler.Filename,
		"size", handler.Size,
		"header", handler.Header,
	)

	defer func() {
		_ = file.Close()
	}()

	uploadBuff, err := io.ReadAll(file)
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		return
	}
	uploadFileName := id.String()

	if len(uploadBuff) >= 512 {
		contentType := http.DetectContentType(uploadBuff[:512])
		var fileExt string
		switch contentType {
		case "image/jpeg", "image/jpg":
			fileExt = ".jpg"
		case "image/png":
			fileExt = ".png"
		case "image/gif":
			fileExt = ".gif"
		case "application/pdf":
			fileExt = ".pdf"
		}
		uploadFileName += fileExt
	}

	url, err := s.uploader.Upload(r.Context(), uploadFileName, bytes.NewReader(uploadBuff))
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := &UploadResponse{Url: url}
	responseData, err := json.Marshal(response)
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(responseData); err != nil {
		s.zapLogger.Sugar().Error(err)
	}
}
