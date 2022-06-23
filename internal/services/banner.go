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

func (s *ApiServer) handleBannerPath() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.CreateBanner(w, r)
		case http.MethodPut:
			s.UpdateBanner(w, r)
		case http.MethodDelete:
			s.DeleteBanner(w, r)
		case http.MethodGet:
			s.GetBanner(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (s *ApiServer) CreateBanner(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"name", r.Form.Get("name"),
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

	name := r.Form.Get("name")
	order, _ := strconv.Atoi(r.Form.Get("order"))
	status, _ := strconv.Atoi(r.Form.Get("status"))
	link := r.Form.Get("link")

	if name == "" {
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
		uploadFileName := fmt.Sprintf("banner%d%d%d%d%d_%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), rand.Intn(100))

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

	banner := &entities.Banner{
		Name:    name,
		Status:  status,
		Link:    link,
		Order:   order,
		Picture: fileUrl,
	}

	err = s.BannerRepo.Insert(r.Context(), s.postgres.Pool, banner)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

	// fmt.Println("");
}

func (s *ApiServer) DeleteBanner(w http.ResponseWriter, r *http.Request) {
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

	banner, err := s.BannerRepo.Get(r.Context(), s.postgres.Pool, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if banner == nil {
		w.WriteHeader(http.StatusBadRequest)
		err = (&errResponse{Error: errors.New("Not found").Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		return
	}

	err = s.BannerRepo.Delete(r.Context(), s.postgres.Pool, id)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}

func (s *ApiServer) GetBanner(w http.ResponseWriter, r *http.Request) {
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

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"id", id,
		)
	}()

	banner, err := s.BannerRepo.Get(r.Context(), s.postgres.Pool, id)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("%+v\n", banner)

	if banner == nil {
		w.WriteHeader(http.StatusBadRequest)
		err = (&errResponse{Error: errors.New("Not found").Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		return
	}

	fmt.Printf("%+v\n", banner)

	resp, err := json.Marshal(banner)
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

type BannerPagingResp struct {
	Banners []*entities.Banner `json:"banners"`
}

func (s *ApiServer) BannerPages(w http.ResponseWriter, r *http.Request) {
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

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"status", status,
		)
	}()

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

	paging := &repositories.BannerPaging{
		Status: status,
		Limit:  limit,
		Offset: offset,
	}

	banners, err := s.BannerRepo.Paging(r.Context(), s.postgres.Pool, paging)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.zapLogger.Sugar().Error(err)
		return
	}

	resp := &BannerPagingResp{
		Banners: banners,
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

func (s *ApiServer) UpdateBanner(w http.ResponseWriter, r *http.Request) {
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
			"name", r.Form.Get("name"),
		)
	}()

	id, _ := strconv.ParseInt(r.Form.Get("id"), 10, 64)

	name := r.Form.Get("name")
	order, _ := strconv.Atoi(r.Form.Get("order"))
	status, _ := strconv.Atoi(r.Form.Get("status"))
	link := r.Form.Get("link")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := (&errResponse{Error: errInvalidParameter.Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	banner, err := s.BannerRepo.Get(r.Context(), s.postgres.Pool, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if banner == nil {
		w.WriteHeader(http.StatusBadRequest)
		err = (&errResponse{Error: errors.New("Not found").Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		return
	}

	var fileUrl = banner.Picture
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
		uploadFileName := fmt.Sprintf("banner%d%d%d%d%d_%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), rand.Intn(100))

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

	bannerUpdate := &entities.Banner{
		Name:    name,
		Link:    link,
		Order:   order,
		Status:  status,
		Picture: fileUrl,
	}

	//fmt.Println(bannerUpdate)

	err = s.BannerRepo.Update(r.Context(), s.postgres.Pool, id, bannerUpdate)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

	// fmt.Println("");
}
