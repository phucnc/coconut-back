package services

import (
	"bytes"
	//"context"
	//"encoding/json"
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
	//"nft-backend/internal/repositories"
	//"strconv"
	//"math/rand"
	"mime/multipart"
	"strings"
	//"time"
)

func (s *ApiServer) handleKYCPath() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.CreateKYC(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (s *ApiServer) CreateKYC(w http.ResponseWriter, r *http.Request) {
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

	fullname := strings.TrimSpace(r.Form.Get("fullname"))
	email := strings.TrimSpace(r.Form.Get("email"))
	birthday := strings.TrimSpace(r.Form.Get("birthday"))
	city := strings.TrimSpace(r.Form.Get("city"))
	country := strings.TrimSpace(r.Form.Get("country"))

	account_address := r.Form.Get("account_id")

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"account_address", account_address,
			"fullname", fullname,
			"email", email,
		)
	}()

	account_check, err := s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, account_address)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if account_check == nil {

		w.WriteHeader(http.StatusBadRequest)
		err := (&errResponse{Error: errAccountNotFound.Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		return
	}

	account_id := account_check.Id

	var front_id = ""
	var back_id = ""
	var selfie_note = ""
	file, _, err := r.FormFile("front_id")
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

		uploadFileName := fmt.Sprintf("kyc_back_%s", account_check.Address)

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
		front_id, err = s.uploader.Upload(r.Context(), uploadFileName, bytes.NewReader(uploadBuff))
		fmt.Println("File URL " + front_id)
		if err != nil {
			s.zapLogger.Sugar().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	file, _, err = r.FormFile("selfie_note")
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

		uploadFileName := fmt.Sprintf("kyc_selfie_%s", account_check.Address)

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
		selfie_note, err = s.uploader.Upload(r.Context(), uploadFileName, bytes.NewReader(uploadBuff))
		fmt.Println("File URL " + selfie_note)
		if err != nil {
			s.zapLogger.Sugar().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	file, _, err = r.FormFile("back_id")
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

		uploadFileName := fmt.Sprintf("kyc_front_%s", account_check.Address)

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
		back_id, err = s.uploader.Upload(r.Context(), uploadFileName, bytes.NewReader(uploadBuff))
		fmt.Println("File URL " + back_id)
		if err != nil {
			s.zapLogger.Sugar().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	kyc := &entities.KYC{
		Account_id: account_id,
		Fullname:   fullname,
		Email:      email,
		Birthday:   birthday,
		City:       city,
		Country:    country,
		FrontId:    front_id,
		BackId:     back_id,
		Selfienote: selfie_note,
	}

	err = s.KYCRepo.Insert(r.Context(), s.postgres.Pool, kyc)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}

func (s *ApiServer) ResetKYC(w http.ResponseWriter, r *http.Request) {
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

	notice, err := s.NoticeRepo.Get(r.Context(), s.postgres.Pool, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if notice == nil {
		w.WriteHeader(http.StatusBadRequest)
		err = (&errResponse{Error: errors.New("Not found").Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		return
	}

	err = s.NoticeRepo.Delete(r.Context(), s.postgres.Pool, id)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}
