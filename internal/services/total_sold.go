package services

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
	"time"
)

type TotalSoldResp struct {
	TotalSold decimal.Decimal `json:"total_sold"`
	//Duration  string          `json:"duration"`
}

func (s *ApiServer) TotalSold24h() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.Method {
		case http.MethodGet:
			status, resp, err := s.totalSold24h(writer, request)
			writer.WriteHeader(status)
			_, err = writer.Write(resp)
			if err != nil {
				s.zapLogger.Sugar().Error(err)
			}
			//writer.WriteHeader(code)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (s *ApiServer) totalSold24h(w http.ResponseWriter, r *http.Request) (int, []byte, error) {
	if err := r.ParseForm(); err != nil {
		s.zapLogger.Sugar().Error(err)
		return http.StatusBadRequest, nil, err
	}

	quoteToken, err := strconv.ParseInt(r.Form.Get("quote_token"), 10, 64)
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		return http.StatusBadRequest, nil, err
	}

	durationStr := r.Form.Get("duration")

	var duration time.Duration
	switch durationStr {
	case "24h", "24H":
		duration = 24 * time.Hour
	case "168h", "168H":
		duration = 168 * time.Hour
	case "all":
		duration = 100000 * time.Hour
	default:
		s.zapLogger.Sugar().Error(err)
		return http.StatusBadRequest, nil, err
	}

	totalSold, err := s.ExchangeEventRepo.GetTotalSoldInDuration(r.Context(), s.postgres.Pool, duration, int16(quoteToken))
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		return http.StatusInternalServerError, nil, err
	}

	resp := &TotalSoldResp{
		TotalSold: totalSold,
		//Duration:  durationStr,
	}

	respData, err := json.Marshal(resp)
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, respData, nil
}
