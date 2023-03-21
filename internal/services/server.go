package services

import (
	"context"
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"nft-backend/internal/database"
	"nft-backend/internal/entities"
	"nft-backend/internal/repositories"
	"nft-backend/postgres"
	"os"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type ApiServer struct {
	zapLogger  *zap.Logger
	httpServer *http.Server
	postgres   *postgres.Postgres

	uploader  Uploader
	TokenRepo interface {
		GetByName(ctx context.Context, db database.QueryExecer, name string) (*entities.Token, error)
	}
	CollectibleRepo interface {
		Insert(ctx context.Context, db database.QueryExecer, collectible *entities.Collectible) error
		Get(ctx context.Context, db database.QueryExecer, id uuid.UUID) (*entities.Collectible, error)
		GetById(ctx context.Context, db database.QueryExecer, id int64) (*entities.Collectible, error)
		GetByName(ctx context.Context, db database.QueryExecer, name string, limit int) ([]*entities.Collectible, error)
		Paging(ctx context.Context, db database.QueryExecer, paging *repositories.CollectiblePaging) ([]*entities.Collectible, error)
		GetTotalInDuration(ctx context.Context, db database.QueryExecer, duration time.Duration) (decimal.Decimal, error)
		Check(ctx context.Context, db database.QueryExecer, id uuid.UUID) (bool, error)
		UpdateView(ctx context.Context, db database.QueryExecer, collectible_id int64) error
		UpdateFakView(ctx context.Context, db database.QueryExecer, collectible_id int64, count uint64) error
		UpdateFakLike(ctx context.Context, db database.QueryExecer, collectible_id int64, count uint64) error
		UpdateResell(ctx context.Context, db database.QueryExecer, collectible_id int64, price decimal.Decimal, quote_token int32) error
		Delete(ctx context.Context, db database.QueryExecer, collectible_id int64) error
		Block(ctx context.Context, db database.QueryExecer, collectible_id int64) error
		UpdateStatus(ctx context.Context, db database.QueryExecer, collectible_id int64, status int) error
	}

	CollectibleLikeRepo interface {
		Insert(ctx context.Context, db database.QueryExecer, collectible_id int64, account_id int64) error
		Delete(ctx context.Context, db database.QueryExecer, collectible_id int64, account_id int64) error
		Check(ctx context.Context, db database.QueryExecer, collectible_id int64, account_id int64) (bool, error)
		Count(ctx context.Context, db database.QueryExecer, collectible_id int64) int
	}
	CategoryRepo interface {
		GetAll(ctx context.Context, db database.QueryExecer) (entities.Categories, error)
		GetCategoryByNames(ctx context.Context, db database.QueryExecer, names []string) ([]*entities.Category, error)
		InsertCategory(ctx context.Context, db database.QueryExecer, name string) error
		UpdateCategory(ctx context.Context, db database.QueryExecer, id int64, name string) error
	}
	CollectibleCategoryRepo interface {
		Insert(ctx context.Context, db database.QueryExecer, collectible *entities.Collectible, category *entities.Category) (*entities.CollectibleCategory, error)
	}
	ExchangeEventRepo interface {
		GetTotalSoldInDuration(ctx context.Context, db database.QueryExecer, duration time.Duration, quoteToken int16) (decimal.Decimal, error)
	}

	AccountRepo interface {
		Insert(ctx context.Context, db database.QueryExecer, account *entities.Account) error
		Update(ctx context.Context, db database.QueryExecer, address string, username string, info string) error
		UpdateAvatar(ctx context.Context, db database.QueryExecer, address string, avatar string) error
		UpdateCover(ctx context.Context, db database.QueryExecer, address string, avatar string) error
		CheckUserName(ctx context.Context, db database.QueryExecer, address string, username string) (bool, error)
		Get(ctx context.Context, db database.QueryExecer, id int64) (*entities.Account, error)
		GetByAddress(ctx context.Context, db database.QueryExecer, address string) (*entities.Account, error)
		GetByUsername(ctx context.Context, db database.QueryExecer, username string) (*entities.Account, error)
		SearchPaging(ctx context.Context, db database.QueryExecer, paging *repositories.AccountPaging) ([]*entities.Account, error)
		Paging(ctx context.Context, db database.QueryExecer, paging *repositories.AccountPaging) ([]*entities.Account, error)
		Block(ctx context.Context, db database.QueryExecer, address string) error
	}

	CommentRepo interface {
		Insert(ctx context.Context, db database.QueryExecer, comment *entities.Comment) error
		Check(ctx context.Context, db database.QueryExecer, id int64, account_id int64) (bool, error)
		Delete(ctx context.Context, db database.QueryExecer, id int64) error
		GetByCollectible(ctx context.Context, db database.QueryExecer, collectible_id int64, limit int, offset int) ([]*entities.Comment, error)
	}

	ReportRepo interface {
		Insert(ctx context.Context, db database.QueryExecer, report *entities.CollectibleReport) error
		Update(ctx context.Context, db database.QueryExecer, id int64, status int64) error
		Paging(ctx context.Context, db database.QueryExecer, paging *repositories.ReportPaging) ([]*entities.Report, error)
	}

	ReportTypeRepo interface {
		GetAllReportTypes(ctx context.Context, db database.QueryExecer) (entities.ReportTypes, error)
		Get(ctx context.Context, db database.QueryExecer, id int64) (*entities.ReportType, error)
		Insert(ctx context.Context, db database.QueryExecer, description string) error
		Update(ctx context.Context, db database.QueryExecer, id int64, description string) error
	}

	TrendRepo interface {
		Get(ctx context.Context, db database.QueryExecer) (entities.Trends, error)
		Insert(ctx context.Context, db database.QueryExecer, trend *entities.Trend) error
		Delete(ctx context.Context, db database.QueryExecer, id int64) error
	}

	EventRepo interface {
		Get(ctx context.Context, db database.QueryExecer, id int64) (*entities.Event, error)
		Insert(ctx context.Context, db database.QueryExecer, trend *entities.Event) error
		Update(ctx context.Context, db database.QueryExecer, id int64, event *entities.Event) error
		Delete(ctx context.Context, db database.QueryExecer, id int64) error
		Paging(ctx context.Context, db database.QueryExecer, paging *repositories.EventPaging) ([]*entities.Event, error)
	}

	NoticeRepo interface {
		Get(ctx context.Context, db database.QueryExecer, id int64) (*entities.Notice, error)
		Insert(ctx context.Context, db database.QueryExecer, trend *entities.Notice) error
		Read(ctx context.Context, db database.QueryExecer, id int) error
		Delete(ctx context.Context, db database.QueryExecer, id int64) error
		Paging(ctx context.Context, db database.QueryExecer, paging *repositories.NoticePaging) ([]*entities.Notice, error)
	}

	BannerRepo interface {
		Get(ctx context.Context, db database.QueryExecer, id int64) (*entities.Banner, error)
		Insert(ctx context.Context, db database.QueryExecer, trend *entities.Banner) error
		Update(ctx context.Context, db database.QueryExecer, id int64, event *entities.Banner) error
		Delete(ctx context.Context, db database.QueryExecer, id int64) error
		Paging(ctx context.Context, db database.QueryExecer, paging *repositories.BannerPaging) ([]*entities.Banner, error)
	}

	KYCRepo interface {
		Insert(ctx context.Context, db database.QueryExecer, kyc *entities.KYC) error
	}
}

func NewServer(ctx context.Context, addr string, zapLogger *zap.Logger, pg *postgres.Postgres, uploader Uploader) *ApiServer {
	httpServer := &http.Server{
		Addr: addr,
	}
	apiServer := &ApiServer{
		httpServer:              httpServer,
		uploader:                uploader,
		postgres:                pg,
		zapLogger:               zapLogger,
		TokenRepo:               &repositories.TokenRepository{},
		CollectibleRepo:         &repositories.CollectibleRepository{},
		CategoryRepo:            &repositories.CategoryRepo{},
		CollectibleCategoryRepo: &repositories.CollectibleCategoryRepo{},
		ExchangeEventRepo:       &repositories.ExchangeEventRepo{},
		AccountRepo:             &repositories.AccountRepository{},
		CollectibleLikeRepo:     &repositories.CollectibleLikeRepo{},
		CommentRepo:             &repositories.CommentRepo{},
		ReportRepo:              &repositories.ReportRepo{},
		ReportTypeRepo:          &repositories.ReportTypeRepo{},
		TrendRepo:               &repositories.TrendRepo{},
		EventRepo:               &repositories.EventRepo{},
		NoticeRepo:              &repositories.NoticeRepo{},
		KYCRepo:                 &repositories.KYCRepository{},
		BannerRepo:              &repositories.BannerRepository{},
	}

	r := mux.NewRouter()
	r.HandleFunc("/", enableCORS(apiServer.notFoundHandler))
	r.HandleFunc("/v1/upload", enableCORS(apiServer.uploadHandler))
	r.HandleFunc("/test-files/upload", enableCORS(apiServer.uploadTemplate))
	r.HandleFunc("/nft", enableCORS(apiServer.handleCollectiblePath()))
	r.HandleFunc("/nft/like", enableCORS(apiServer.Like()))
	r.HandleFunc("/nft/fak-like", enableCORS(apiServer.FakLike()))
	r.HandleFunc("/nft/block", enableCORS(apiServer.BlockCollective()))
	r.HandleFunc("/nft/total-mint", enableCORS(apiServer.TotalMint()))
	r.HandleFunc("/nft/collectible-search", enableCORS(apiServer.GetByName()))
	r.HandleFunc("/nft/collectible-paging", enableCORS(apiServer.CollectiblePaging))
	r.HandleFunc("/nft/{id}", enableCORS(apiServer.GetCollectibleMux))
	r.HandleFunc("/exchange/total-sold", enableCORS(apiServer.TotalSold24h()))
	r.HandleFunc("/account", enableCORS(apiServer.handleAccountPath()))
	r.HandleFunc("/account/cover", enableCORS(apiServer.Cover()))
	r.HandleFunc("/account/avatar", enableCORS(apiServer.Avatar()))
	r.HandleFunc("/comment", enableCORS(apiServer.handleCommentPath()))
	r.HandleFunc("/view", enableCORS(apiServer.View()))
	r.HandleFunc("/fak-view", enableCORS(apiServer.FakView()))
	r.HandleFunc("/status", enableCORS(apiServer.Status()))
	r.HandleFunc("/report", enableCORS(apiServer.handleReportPath()))
	r.HandleFunc("/account/search-paging", enableCORS(apiServer.AccountSearchPaging))
	r.HandleFunc("/account/paging", enableCORS(apiServer.AccountPaging))
	r.HandleFunc("/account/block", enableCORS(apiServer.Block()))
	r.HandleFunc("/comment/paging", enableCORS(apiServer.GetCommentByCollectible))
	r.HandleFunc("/history", enableCORS(apiServer.GetHistory))
	r.HandleFunc("/report-type", enableCORS(apiServer.handleReportTypePath()))
	r.HandleFunc("/category", enableCORS(apiServer.handleCategoryPath()))
	r.HandleFunc("/trend", enableCORS(apiServer.handleTrendPath()))
	r.HandleFunc("/event", enableCORS(apiServer.handleEventPath()))
	r.HandleFunc("/event/paging", enableCORS(apiServer.EventPages))
	r.HandleFunc("/notice", enableCORS(apiServer.handleNoticePath()))
	r.HandleFunc("/notice/paging", enableCORS(apiServer.NoticePages))
	r.HandleFunc("/kyc", enableCORS(apiServer.handleKYCPath()))
	r.HandleFunc("/banner", enableCORS(apiServer.handleBannerPath()))
	r.HandleFunc("/banner/paging", enableCORS(apiServer.BannerPages))

	httpServer.Handler = r

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	return apiServer
}

func (s *ApiServer) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func enableCORS(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == http.MethodOptions {
			return
		}
		nextHandler(w, r)
	}
}

//go:embed "upload.html"
var uploadTemplate embed.FS

func (s *ApiServer) uploadTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		executeAndReturnInternalIfErr(w, r, func() error {
			fmt.Println(os.Getwd())
			t, err := template.ParseFS(uploadTemplate, "*.html")
			if err != nil {
				s.zapLogger.Sugar().Error(err)
				return err
			}

			err = t.Execute(w, nil)
			if err != nil {
				s.zapLogger.Sugar().Error(err)
				return err
			}
			return nil
		})
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func executeAndReturnInternalIfErr(w http.ResponseWriter, r *http.Request, tx func() error) {
	err := tx()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *ApiServer) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
