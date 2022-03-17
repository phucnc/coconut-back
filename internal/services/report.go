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
	"nft-backend/internal/repositories"
	//"strconv"
	"database/sql"
	"strings"
)

func (s *ApiServer) handleReportPath() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.CreateReport(w, r)
		case http.MethodPut:
			s.Update(w, r)
		case http.MethodGet:
			s.Pages(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (s *ApiServer) CreateReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	account_id := r.Form.Get("account_id")
	collectible_id := r.Form.Get("collectible_id")
	report_type_id, err := strconv.Atoi(r.Form.Get("report_type_id"))
	content := strings.TrimSpace(r.Form.Get("content"))

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
			"account_id", account_id,
			"collectible_id", collectible_id,
			"report_type_id", report_type_id,
			"content", content,
		)
	}()

	if account_id == "" || report_type_id == 0 {
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

	report := &entities.CollectibleReport{
		CollectibleID: collectible_check.Id,
		AccountID:     account_check.Id,
		ReportTypeID:  report_type_id,
		Content:       content,
	}

	err = s.ReportRepo.Insert(r.Context(), s.postgres.Pool, report)
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

	var account_owner_id sql.NullInt64
	account_owner_id.Int64 = account_check.Id
	account_owner_id.Valid = true
	var collectible_id_64 sql.NullInt64
	collectible_id_64.Int64 = collectible_check.Id
	collectible_id_64.Valid = true

	notice_receive := &entities.Notice{
		AccountID:     account_owner.Id,
		FromAccountID: account_owner_id,
		CollectibleID: collectible_id_64,
		Content:       entities.Notice_report_receive,
	}

	var account_check_id sql.NullInt64
	account_check_id.Int64 = account_check.Id
	account_check_id.Valid = true

	notice_send := &entities.Notice{
		AccountID:     account_check.Id,
		FromAccountID: account_check_id,
		CollectibleID: collectible_id_64,
		Content:       entities.Notice_report_send,
	}

	err = s.NoticeRepo.Insert(r.Context(), s.postgres.Pool, notice_receive)
	err = s.NoticeRepo.Insert(r.Context(), s.postgres.Pool, notice_send)

	w.WriteHeader(http.StatusOK)
	return

	// fmt.Println("");
}

type ReportPagingResp struct {
	Reports []*entities.Report `json:"reports"`
}

func (s *ApiServer) Pages(w http.ResponseWriter, r *http.Request) {
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

	collectible_id := r.Form.Get("collectible_id")
	status, err := strconv.ParseInt(r.Form.Get("status"), 10, 64)
	var collectible_Id = int64(0)
	var account_Id = int64(0)
	account_id := r.Form.Get("account_id")

	fromTime, _ := strconv.Atoi(r.Form.Get("from"))
	toTime, _ := strconv.Atoi(r.Form.Get("to"))

	if collectible_id != "" {
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

		collectible_Id = collectible_check.Id
	}

	if account_id != "" {
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

		account_Id = account_check.Id
	}

	limitParam := r.Form.Get("limit")

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

	ofsetParam := r.Form.Get("offset")
	offset, err := strconv.ParseInt(ofsetParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.zapLogger.Sugar().Error(err)
		return
	}
	if offset <= 0 {
		offset = 0
	}

	paging := &repositories.ReportPaging{
		Collectible: collectible_Id,
		Account:     account_Id,
		Status:      status,
		Limit:       limit,
		Offset:      offset,
		From:        fromTime,
		To:          toTime,
	}

	reports, err := s.ReportRepo.Paging(r.Context(), s.postgres.Pool, paging)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.zapLogger.Sugar().Error(err)
		return
	}

	for index, report := range reports {
		//fmt.Println(index)

		account, err := s.AccountRepo.Get(r.Context(), s.postgres.Pool, report.Account_Id)
		if err != nil {

			account = nil

		}
		report.Account = account

		//like

		collectible, err := s.CollectibleRepo.GetById(r.Context(), s.postgres.Pool, report.Collectible_Id)
		if err != nil {

			collectible = nil

		}
		report.Collectible = collectible

		rtype, err := s.ReportTypeRepo.Get(r.Context(), s.postgres.Pool, report.Report_Type_Id)
		if err != nil {

			rtype = nil
			report.Description = ""

		} else {
			report.Description = rtype.Description
		}

		reports[index] = report

	}

	resp := &ReportPagingResp{
		Reports: reports,
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

func (s *ApiServer) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	status, err := strconv.ParseInt(r.Form.Get("status"), 10, 64)
	id, err := strconv.ParseInt(r.Form.Get("id"), 10, 64)

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"id", id,
			"status", status,
		)
	}()

	if id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		err := (&errResponse{Error: errInvalidParameter.Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err = s.ReportRepo.Update(r.Context(), s.postgres.Pool, id, status)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

	// fmt.Println("");
}
