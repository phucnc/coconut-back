package services

import (
	"bytes"
	//"context"
	"encoding/json"
	"fmt"
	"strconv"
	//"github.com/gofrs/uuid"
	//"github.com/gorilla/mux"
	//"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	//"github.com/shopspring/decimal"
	"io"
	"net/http"
	//"net/url"
	//"nft-backend/internal/database"
	"nft-backend/internal/entities"
	"nft-backend/internal/repositories"
	//"strconv"
	"math/rand"
	"mime/multipart"
	"strings"
	"time"
)

func (s *ApiServer) handleEventPath() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.CreateEvent(w, r)
		case http.MethodPut:
			s.UpdateEvent(w, r)
		case http.MethodDelete:
			s.DeleteEvent(w, r)
		case http.MethodGet:
			s.GetEvent(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (s *ApiServer) CreateEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"title", r.Form.Get("title"),
			"admin_created", r.Form.Get("admin_created"),
		)
	}()

	switch err := r.ParseMultipartForm(2 * megabyte); err {
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	title := r.Form.Get("title")
	content := r.Form.Get("content")

	if title == "" || content == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := (&errResponse{Error: errInvalidParameter.Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	var fileUrl = ""
	file, _, err := r.FormFile("upload_file")
	switch err {
	case nil:
		defer func() {
			_ = file.Close()
		}()
		break
	default:
		break
	}

	if err == nil {
		uploadBuff, err := io.ReadAll(file)
		if err != nil {
			s.zapLogger.Sugar().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(uploadBuff) <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			if err := (&errResponse{Error: errUploadFileCannotBeEmpty.Error()}).writeResponse(w); err != nil {
				s.zapLogger.Sugar().Error(err)
			}
		}

		t := time.Now()
		uploadFileName := fmt.Sprintf("event%d%d%d%d%d_%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), rand.Intn(100))

		if len(uploadBuff) >= 512 {
			contentType := http.DetectContentType(uploadBuff[:512])
			var fileExt string
			switch contentType {
			case "image/jpeg", "image/jpg":
				fileExt = ".jpg"
			case "image/png":
				fileExt = ".png"
			default:
				splitStr := strings.Split(contentType, "/")
				if len(splitStr) == 2 {
					fileExt = fmt.Sprintf(".%s", splitStr[1])
				}
			}
			uploadFileName += fileExt
		}
		fmt.Println("File name " + uploadFileName)
		fileUrl, err = s.uploader.Upload(r.Context(), uploadFileName, bytes.NewReader(uploadBuff))
		fmt.Println("File URL " + fileUrl)
		if err != nil {
			s.zapLogger.Sugar().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	event := &entities.Event{
		Title:   title,
		Content: content,
		Banner:  fileUrl,
	}

	err = s.EventRepo.Insert(r.Context(), s.postgres.Pool, event)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

	// fmt.Println("");
}

func (s *ApiServer) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseForm()
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, _ := strconv.ParseInt(r.Form.Get("id"), 10, 64)

	if id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		err := (&errResponse{Error: errIdCannotBeEmpty.Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	event, err := s.EventRepo.Get(r.Context(), s.postgres.Pool, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if event == nil {
		w.WriteHeader(http.StatusBadRequest)
		err = (&errResponse{Error: errors.New("Not found").Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		return
	}

	err = s.EventRepo.Delete(r.Context(), s.postgres.Pool, id)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}

func (s *ApiServer) GetEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseForm()
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, _ := strconv.ParseInt(r.Form.Get("id"), 10, 64)

	if id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		err := (&errResponse{Error: errIdCannotBeEmpty.Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	event, err := s.EventRepo.Get(r.Context(), s.postgres.Pool, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if event == nil {
		w.WriteHeader(http.StatusBadRequest)
		err = (&errResponse{Error: errors.New("Not found").Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		return
	}

	resp, err := json.Marshal(event)
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(resp); err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

type EventPagingResp struct {
	Events []*entities.Event `json:"reports"`
}

func (s *ApiServer) EventPages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.zapLogger.Sugar().Error(err)
		return
	}

	status, err := strconv.Atoi(r.Form.Get("status"))
	limit, _ := strconv.Atoi(r.Form.Get("limit"))

	if limit > 100 {
		limit = 100
	}

	if limit < 10 {
		limit = 10
	}

	offset, err := strconv.Atoi(r.Form.Get("offset"))
	if offset <= 0 {
		offset = 0
	}

	paging := &repositories.EventPaging{
		Status: status,
		Limit:  limit,
		Offset: offset,
	}

	events, err := s.EventRepo.Paging(r.Context(), s.postgres.Pool, paging)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.zapLogger.Sugar().Error(err)
		return
	}

	resp := &EventPagingResp{
		Events: events,
	}

	respData, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.zapLogger.Sugar().Error(err)
		return
	}

	_, err = w.Write(respData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.zapLogger.Sugar().Error(err)
		return
	}
}

func (s *ApiServer) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch err := r.ParseMultipartForm(2 * megabyte); err {
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"id", r.Form.Get("id"),
			"title", r.Form.Get("title"),
		)
	}()

	id, _ := strconv.ParseInt(r.Form.Get("id"), 10, 64)

	title := r.Form.Get("title")
	content := r.Form.Get("content")

	if title == "" || content == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := (&errResponse{Error: errInvalidParameter.Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	event, err := s.EventRepo.Get(r.Context(), s.postgres.Pool, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if event == nil {
		w.WriteHeader(http.StatusBadRequest)
		err = (&errResponse{Error: errors.New("Not found").Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		return
	}

	var fileUrl = event.Banner
	file, _, err := r.FormFile("upload_file")
	switch err {
	case nil:
		defer func() {
			_ = file.Close()
		}()
		break
	default:
		break
	}

	if err == nil {
		uploadBuff, err := io.ReadAll(file)
		if err != nil {
			s.zapLogger.Sugar().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(uploadBuff) <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			if err := (&errResponse{Error: errUploadFileCannotBeEmpty.Error()}).writeResponse(w); err != nil {
				s.zapLogger.Sugar().Error(err)
			}
		}

		t := time.Now()
		uploadFileName := fmt.Sprintf("event%d%d%d%d%d_%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), rand.Intn(100))

		if len(uploadBuff) >= 512 {
			contentType := http.DetectContentType(uploadBuff[:512])
			var fileExt string
			switch contentType {
			case "image/jpeg", "image/jpg":
				fileExt = ".jpg"
			case "image/png":
				fileExt = ".png"
			default:
				splitStr := strings.Split(contentType, "/")
				if len(splitStr) == 2 {
					fileExt = fmt.Sprintf(".%s", splitStr[1])
				}
			}
			uploadFileName += fileExt
		}
		fmt.Println("File name " + uploadFileName)
		fileUrl, err = s.uploader.Upload(r.Context(), uploadFileName, bytes.NewReader(uploadBuff))
		fmt.Println("File URL " + fileUrl)
		if err != nil {
			s.zapLogger.Sugar().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	eventUpdate := &entities.Event{
		Title:   title,
		Content: content,
		Banner:  fileUrl,
	}

	err = s.EventRepo.Update(r.Context(), s.postgres.Pool, id, eventUpdate)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

	// fmt.Println("");
}
