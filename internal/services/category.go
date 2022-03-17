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


func (s *ApiServer) handleCategoryPath() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			 s.AddCategory(w,r)
		case http.MethodPut:
			s.UpdateCategory(w,r)
		case http.MethodGet:
			 s.GetCategory(w,r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

type CategoriesResp struct {
	Categories []*entities.Category `json:"category"`
}


func (s *ApiServer) AddCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    
  	if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

	name  := strings.TrimSpace(r.Form.Get("name"))
	

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"name", name,
		)
	}()

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


	 err := s.CategoryRepo.InsertCategory(r.Context(), s.postgres.Pool,  name)
	  if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}

			return 

	// fmt.Println("");
}


func (s *ApiServer) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    
  	if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

	name  := strings.TrimSpace(r.Form.Get("name"))
	id, err := strconv.ParseInt(r.Form.Get("id"), 10,64) 

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"id", id,
			"name", name,
		)
	}()

	 if name == "" || id ==0  {
			w.WriteHeader(http.StatusBadRequest)
			err := (&errResponse{Error: errInvalidParameter.Error()}).writeResponse(w)
			if err != nil {
				s.zapLogger.Sugar().Warn(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

	 err = s.CategoryRepo.UpdateCategory(r.Context(), s.postgres.Pool, id, name)
	  if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}

		w.WriteHeader(http.StatusOK)
			return 

	// fmt.Println("");
}

func (s *ApiServer) GetCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    
	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
		)
	}()

	categories, err := s.CategoryRepo.GetAll(r.Context(), s.postgres.Pool)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.zapLogger.Sugar().Error(err)
		return
	}
	
	resp := &CategoriesResp{
		Categories: categories,
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