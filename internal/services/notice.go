package services

import (
	//"bytes"
	//"context"
	"encoding/json"
	"fmt"
	"strconv"
	//"github.com/gofrs/uuid"
	//"github.com/gorilla/mux"
	//"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	//"github.com/shopspring/decimal"
	//"io"
	"net/http"
	//"net/url"
	//"nft-backend/internal/database"
	"nft-backend/internal/entities"
	"nft-backend/internal/repositories"
	//"strconv"
	//"math/rand"
	//"mime/multipart"
	"strings"
	//"time"
)

func (s *ApiServer) handleNoticePath() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.CreateNotice(w, r)
		case http.MethodPut:
			s.Read(w, r)
		case http.MethodDelete:
			s.DeleteNotice(w, r)
		case http.MethodGet:
			s.GetNotice(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (s *ApiServer) CreateNotice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	content := strings.TrimSpace(r.Form.Get("content"))

	account_id, err := strconv.ParseInt(r.Form.Get("account_id"), 10, 64)
	if err != nil {
		panic(err)
	}

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"account_id", account_id,
			"content", content,
		)
	}()

	if content == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := (&errResponse{Error: errInvalidParameter.Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	account_check, err := s.AccountRepo.Get(r.Context(), s.postgres.Pool, account_id)
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

	notice := &entities.Notice{
		AccountID: account_check.Id,
		Content:   content,
	}

	err = s.NoticeRepo.Insert(r.Context(), s.postgres.Pool, notice)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

	// fmt.Println("");
}

func (s *ApiServer) DeleteNotice(w http.ResponseWriter, r *http.Request) {
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

func (s *ApiServer) GetNotice(w http.ResponseWriter, r *http.Request) {
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

	resp, err := json.Marshal(notice)
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

type NoticePagingResp struct {
	Notices []*entities.Notice `json:"notices"`
}

func (s *ApiServer) NoticePages(w http.ResponseWriter, r *http.Request) {
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
	account_id := r.Form.Get("account_id")

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"account_id", account_id,
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

	account_check, err := s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, account_id)
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

	//account_id, err := strconv.ParseInt("12", 10, 64)

	paging := &repositories.NoticePaging{
		Status:  status,
		Account: account_check.Id,
		Limit:   limit,
		Offset:  offset,
	}

	notices, err := s.NoticeRepo.Paging(r.Context(), s.postgres.Pool, paging)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.zapLogger.Sugar().Error(err)
		return
	}

	for index, notice := range notices {

		from_account, err := s.AccountRepo.Get(r.Context(), s.postgres.Pool, notice.FromAccountID.Int64)
		if err != nil {

			from_account = nil

		}

		collectibe, err := s.CollectibleRepo.GetById(r.Context(), s.postgres.Pool, notice.CollectibleID.Int64)
		if err != nil {
			collectibe = nil
		}

		notice.FromAccount = from_account
		notice.Collectible = collectibe

		notices[index] = notice
	}

	resp := &NoticePagingResp{
		Notices: notices,
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

func (s *ApiServer) Read(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	notice_id, err := strconv.Atoi(r.Form.Get("id"))
	if err == nil {
		fmt.Printf("%d of type %T", notice_id, notice_id)
	}

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"id", notice_id,
		)
	}()

	err = s.NoticeRepo.Read(r.Context(), s.postgres.Pool, notice_id)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}
