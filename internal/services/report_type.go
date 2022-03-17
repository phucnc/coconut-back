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
	//"github.com/pkg/errors"
	//"github.com/shopspring/decimal"
	//"io"
	"net/http"
	//"net/url"
	//"nft-backend/internal/database"
	"nft-backend/internal/entities"
	//"nft-backend/internal/repositories"
	//"strconv"
	"strings"
)


func (s *ApiServer) handleReportTypePath() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			 s.AddReportType(w,r)
		case http.MethodPut:
			s.UpdateReportType(w,r)
		case http.MethodGet:
			 s.GetReportType(w,r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

type ReportTypesResp struct {
	ReportTypes []*entities.ReportType `json:"report_types"`
}


func (s *ApiServer) AddReportType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    
  	if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

	description  := strings.TrimSpace(r.Form.Get("description"))
	

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"description", description,
		)
	}()

	if description == "" {
			w.WriteHeader(http.StatusBadRequest)
			err := (&errResponse{Error: errInvalidParameter.Error()}).writeResponse(w)
			if err != nil {
				s.zapLogger.Sugar().Warn(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}


	 err := s.ReportTypeRepo.Insert(r.Context(), s.postgres.Pool,  description)
	  if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}
		 w.WriteHeader(http.StatusOK)
		return 

	// fmt.Println("");
}


func (s *ApiServer) UpdateReportType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    
  	if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

	description  := strings.TrimSpace(r.Form.Get("description"))
	id, err := strconv.ParseInt(r.Form.Get("id"), 10,64) 

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"id", id,
			"description", description,
		)
	}()

	 if description == "" || id ==0  {
			w.WriteHeader(http.StatusBadRequest)
			err := (&errResponse{Error: errInvalidParameter.Error()}).writeResponse(w)
			if err != nil {
				s.zapLogger.Sugar().Warn(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

	 err = s.ReportTypeRepo.Update(r.Context(), s.postgres.Pool, id, description)
	  if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}

		w.WriteHeader(http.StatusOK)
			return 

	// fmt.Println("");
}

func (s *ApiServer) GetReportType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    
	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
		)
	}()

	reporttypes, err := s.ReportTypeRepo.GetAllReportTypes(r.Context(), s.postgres.Pool)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.zapLogger.Sugar().Error(err)
		return
	}
	
	resp := &ReportTypesResp{
		ReportTypes: reporttypes,
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