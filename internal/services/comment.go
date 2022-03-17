package services

import (

	//"bytes"
	//"context"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"strconv"
	//"github.com/gorilla/mux"
	//"github.com/jackc/pgx/v4"
	//"github.com/pkg/errors"
	//"github.com/shopspring/decimal"
	//"io"
	"net/http"
	//"net/url"
	//"nft-backend/internal/database"
	"nft-backend/internal/entities"
	//"nft-backend/internal/repositories"
	//"strconv"
	"database/sql"
	"strings"
)

type CreateCommentResponse struct {
	ID int64 `json:"id"`
}

func (s *ApiServer) handleCommentPath() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.Create(w, r)
		case http.MethodGet:
			s.GetCommentByCollectible(w, r)
		case http.MethodDelete:
			s.DeleteComment(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (s *ApiServer) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	account_id := r.Form.Get("account_id")
	collectible_id := r.Form.Get("collectible_id")
	content := strings.TrimSpace(r.Form.Get("content"))

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"account_id", account_id,
			"collectible_id", collectible_id,
			"content", content,
		)
	}()

	collectible_uid, err := uuid.FromString(collectible_id)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		if err := (&errResponse{Error: errIdIsNotValid.Error()}).writeResponse(w); err != nil {
			s.zapLogger.Sugar().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	if account_id == "" || content == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := (&errResponse{Error: errInvalidParameter.Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
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

	collectible_check, err := s.CollectibleRepo.Get(r.Context(), s.postgres.Pool, collectible_uid)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if collectible_check == nil {
		w.WriteHeader(http.StatusBadRequest)
		err := (&errResponse{Error: errCollectibleNotFound.Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		return
	}

	comment := &entities.Comment{
		CollectibleID: collectible_check.Id,
		AccountID:     account_check.Id,
		Content:       content,
	}

	err = s.CommentRepo.Insert(r.Context(), s.postgres.Pool, comment)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	account_owner, err := s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, *collectible_check.TokenOwner)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var account_id_64 sql.NullInt64
	account_id_64.Int64 = account_check.Id
	account_id_64.Valid = true
	var collectible_id_64 sql.NullInt64
	collectible_id_64.Int64 = collectible_check.Id
	collectible_id_64.Valid = true

	notice := &entities.Notice{
		AccountID:     account_owner.Id,
		FromAccountID: account_id_64,
		CollectibleID: collectible_id_64,
		Content:       entities.Notice_someone_comment,
	}

	err = s.NoticeRepo.Insert(r.Context(), s.postgres.Pool, notice)

	w.WriteHeader(http.StatusOK)
	return

	// fmt.Println("");
}

func (s *ApiServer) DeleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	account_id := r.Form.Get("account_id")
	id, err := strconv.ParseInt(r.Form.Get("id"), 10, 64)

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"account_id", account_id,
			"id", id,
		)
	}()

	if account_id == "" || id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		err := (&errResponse{Error: errInvalidParameter.Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
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

	comment_check, err := s.CommentRepo.Check(r.Context(), s.postgres.Pool, id, account_check.Id)

	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if comment_check == false {
		w.WriteHeader(http.StatusBadRequest)
		err := (&errResponse{Error: errNotFound.Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		return
	} else {

		err = s.CommentRepo.Delete(r.Context(), s.postgres.Pool, id)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return

	}

	// fmt.Println("");
}

type CommentPagingResp struct {
	Comments []*entities.Comment `json:"comments"`
}

func (s *ApiServer) GetCommentByCollectible(w http.ResponseWriter, r *http.Request) {
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

	limitParam := r.Form.Get("limit")
	if limitParam == "" {
		limitParam = "10"
	}

	offsetParam := r.Form.Get("offset")
	if offsetParam == "" {
		offsetParam = "0"
	}

	limit, err := strconv.ParseInt(limitParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.zapLogger.Sugar().Error(err)
		return
	}
	if limit < 1 {
		w.WriteHeader(http.StatusBadRequest)
		s.zapLogger.Sugar().Error(err)
		return
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.ParseInt(offsetParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.zapLogger.Sugar().Error(err)
		return
	}
	if offset < 1 {
		offset = 0
	}

	collectible_id := r.Form.Get("collectible_id")

	collectible_uid, err := uuid.FromString(collectible_id)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		if err := (&errResponse{Error: errIdIsNotValid.Error()}).writeResponse(w); err != nil {
			s.zapLogger.Sugar().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"collectible_id", collectible_id,
			"limit", limit,
			"offset", offset,
		)
	}()

	collectible_check, err := s.CollectibleRepo.Get(r.Context(), s.postgres.Pool, collectible_uid)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if collectible_check == nil {
		w.WriteHeader(http.StatusBadRequest)
		err := (&errResponse{Error: errCollectibleNotFound.Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		return
	}

	comments, err := s.CommentRepo.GetByCollectible(r.Context(), s.postgres.Pool, collectible_check.Id, int(limit), int(offset))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.zapLogger.Sugar().Error(err)
		return
	}

	for index, comment := range comments {
		//fmt.Println(index)

		account, err := s.AccountRepo.Get(r.Context(), s.postgres.Pool, comment.AccountID)
		if err != nil {

			account = nil

		}
		comment.Account = account
		//like
		comments[index] = comment

	}

	resp := &CommentPagingResp{
		Comments: comments,
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
