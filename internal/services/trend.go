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
)

func (s *ApiServer) handleTrendPath() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.AddTrend(w, r)
		case http.MethodDelete:
			s.DeleteTrend(w, r)
		case http.MethodGet:
			s.GetTrend(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

type TrendResp struct {
	Trends []*entities.Trend `json:"trend"`
}

func (s *ApiServer) AddTrend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
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

	in_order, _ := strconv.ParseInt(r.Form.Get("order"), 10, 64)
	advertisement, _ := strconv.ParseBool((r.Form.Get("advertisement")))

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"collectible_id", collectible_id,
			"in_order", in_order,
			"adverisement", advertisement,
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

	trend := &entities.Trend{
		CollectibleID: collectible_check.Id,
		In_Order:      in_order,
		Advertisement: advertisement,
	}

	err = s.TrendRepo.Insert(r.Context(), s.postgres.Pool, trend)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}

func (s *ApiServer) GetTrend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
		)
	}()

	address := r.Form.Get("address")
	account := &entities.Account{}

	trends, err := s.TrendRepo.Get(r.Context(), s.postgres.Pool)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.zapLogger.Sugar().Error(err)
		return
	}

	if address != "" {
		account, err = s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, address)
	}

	account_id, err := strconv.ParseInt("-1", 10, 64)
	if err != nil {
		fmt.Printf("%d of type %T", account_id, account_id)
	}

	if account != nil {
		account_id = account.Id
	}

	for index, trend := range trends {
		//fmt.Println(index)

		//like

		collectible, err := s.CollectibleRepo.GetById(r.Context(), s.postgres.Pool, trend.CollectibleID)
		if err != nil {

			collectible = nil

		}

		if collectible != nil {

			creator, err := s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, collectible.Creator)
			if err != nil {

				creator = nil

			}

			owner, err := s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, *collectible.TokenOwner)
			if err != nil {

				owner = nil

			}

			collectible.Creator_acc = creator
			collectible.Owner = owner
			//like

			liked, err := s.CollectibleLikeRepo.Check(r.Context(), s.postgres.Pool, collectible.Id, account_id)
			if err != nil {
				s.zapLogger.Sugar().Error(err)
				liked = false
			}
			total := s.CollectibleLikeRepo.Count(r.Context(), s.postgres.Pool, collectible.Id)

			like := &entities.CollectibleLikeResp{
				Liked: liked,
				Total: total,
			}
			collectible.Like = like
		}

		trend.Collectible = collectible

		trends[index] = trend
	}

	resp := &TrendResp{
		Trends: trends,
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

func (s *ApiServer) DeleteTrend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	/*if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}*/
	r.ParseForm()

	id, err := strconv.ParseInt(r.Form.Get("id"), 10, 64)

	fmt.Println(id)
	fmt.Println(r.Form.Get("id2"))

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"id", id,
		)
	}()

	err = s.TrendRepo.Delete(r.Context(), s.postgres.Pool, id)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

	// fmt.Println("");
}
