package services

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"nft-backend/internal/database"
	"nft-backend/internal/entities"
	"nft-backend/internal/repositories"
	"strconv"
	"strings"
	"time"
)

type CreateCollectibleResponse struct {
	ID string `json:"id"`
}

type TotalResp struct {
	Total decimal.Decimal `json:"total"`
	//Duration  string          `json:"duration"`
}

func (s *ApiServer) handleCollectiblePath() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.CreateCollectible(w, r)
		case http.MethodPut:
			s.ResellUpdate(w, r)
		case http.MethodGet:
			s.GetCollectible(w, r)
		case http.MethodDelete:
			s.Delete(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

var minRoyalties = decimal.NewFromInt(0)
var maxRoyalties = decimal.NewFromInt(1)

func (s *ApiServer) Like() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.DoLike(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (s *ApiServer) View() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.UpdateView(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (s *ApiServer) Status() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.UpdateStatus(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (s *ApiServer) History() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.GetHistory(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (s *ApiServer) BlockCollective() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			s.DoBlockCollective(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func newCollectibleFromForm(form url.Values) (*entities.Collectible, error) {
	unlockOncePurchased, err := strconv.ParseBool(form.Get(formCreateCollectible_UnlockOncePurchased))
	if err != nil {
		unlockOncePurchased = false
		//return nil, errUnlockOncePurchasedIsNotValid
	}
	royalties, err := decimal.NewFromString(form.Get(formCreateCollectible_RoyaltyPercent))
	if err != nil {
		royalties = minRoyalties
		//return nil, errRoyaltiesIsNotValid
	}
	royalties = royalties.Div(decimal.NewFromInt(100))

	if royalties.LessThan(minRoyalties) || royalties.GreaterThan(maxRoyalties) {
		return nil, errRoyaltiesIsNotValid
	}

	title := form.Get(formCreateCollectible_Title)
	if title == "" {
		return nil, errTitleCannotBeEmpty
	}

	instantSalePrice, err := decimal.NewFromString(form.Get(formCreateCollectible_InstantSalePrice))
	if err != nil {
		return nil, errInstantSalePriceIsNotValid
	}

	creator := form.Get(formCreateCollectible_Creator)

	if creator == "" {
		return nil, errAccountNotFound
	}

	quoteTokenName := form.Get("quote_token")
	if quoteTokenName == "" {
		quoteTokenName = "CONT"
		//return nil, errTitleCannotBeEmpty
	}

	collectible := &entities.Collectible{
		Title:               form.Get(formCreateCollectible_Title),
		Description:         form.Get(formCreateCollectible_Description),
		Royalties:           royalties,
		UnlockOncePurchased: unlockOncePurchased,
		InstantSalePrice:    instantSalePrice,
		UploadFile:          "",
		Properties:          nil,
		QuoteToken: &entities.Token{
			Name: quoteTokenName,
		},
		Categories: make([]*entities.Category, 0),
		Creator:    creator,
	}
	return collectible, nil
}

const (
	_byte    = 1
	kilobyte = 1024 * _byte
	megabyte = 1024 * kilobyte
)

var (
	errIdCannotBeEmpty               = errors.New("id cannot be empty")
	errIdIsNotValid                  = errors.New("id is not valid")
	errTitleCannotBeEmpty            = errors.New("title cannot be empty")
	errCategoriesCannotBeEmpty       = errors.New("categories cannot be empty")
	errCategoryIsNotValid            = errors.New("category is not valid")
	errUnlockOncePurchasedIsNotValid = errors.New("unlockOncePurchased is not valid")
	errRoyaltiesIsNotValid           = errors.New("royalties is not valid, must be >= 0 and <= 100")
	errInstantSalePriceIsNotValid    = errors.New("instant sale price is not valid, must be >= 0")
	errUploadFileDoesNotExist        = errors.New("upload file does not exist")
	errUploadFileCannotBeEmpty       = errors.New("upload file cannot be empty")
	errInvalidParameter              = errors.New("invalid Paramater")
	errAccountNotFound               = errors.New("Account Not found")
	errCollectibleNotFound           = errors.New("Collectible Not found")
)

const (
	formCreateCollectible_UnlockOncePurchased = "unlock_once_purchased"
	formCreateCollectible_InstantSalePrice    = "instant_sale_price"
	formCreateCollectible_RoyaltyPercent      = "royalty_percent"
	formCreateCollectible_Title               = "title"
	formCreateCollectible_Description         = "description"
	formCreateCollectible_Categories          = "categories"
	formCreateCollectible_UploadFile          = "upload_file"
	formCreateCollectible_Creator             = "creator"
)

func (s *ApiServer) CreateCollectible(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			formCreateCollectible_UnlockOncePurchased, r.Form.Get(formCreateCollectible_UnlockOncePurchased),
			formCreateCollectible_InstantSalePrice, r.Form.Get(formCreateCollectible_InstantSalePrice),
			formCreateCollectible_RoyaltyPercent, r.Form.Get(formCreateCollectible_RoyaltyPercent),
			formCreateCollectible_Title, r.Form.Get(formCreateCollectible_Title),
			formCreateCollectible_Description, r.Form.Get(formCreateCollectible_Description),
			formCreateCollectible_Categories, r.Form[formCreateCollectible_Categories],
			formCreateCollectible_Creator, r.Form[formCreateCollectible_Creator],
		)
	}()

	switch err := r.ParseMultipartForm(50 * megabyte); err {
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

	collectible, err := newCollectibleFromForm(r.Form)
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusBadRequest)
		if err := (&errResponse{Error: err.Error()}).writeResponse(w); err != nil {
			s.zapLogger.Sugar().Error(err)
		}
		return
	}

	account_check, err := s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, collectible.Creator)
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

	file, _, err := r.FormFile(formCreateCollectible_UploadFile)
	switch err {
	case nil:
		defer func() {
			_ = file.Close()
		}()
		break
	case http.ErrMissingFile:
		w.WriteHeader(http.StatusBadRequest)
		if err := (&errResponse{Error: errUploadFileDoesNotExist.Error()}).writeResponse(w); err != nil {
			s.zapLogger.Sugar().Error(err)
		}
		return
	default:
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	/*s.zapLogger.Sugar().Infow(
		"handle upload file",
		"name", handler.Filename,
		"size", handler.Size,
		"header", handler.Header,
	)*/

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

	id, err := uuid.NewV4()
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	uploadFileName := id.String()

	if len(uploadBuff) >= 512 {
		contentType := http.DetectContentType(uploadBuff[:512])
		var fileExt string
		switch contentType {
		case "image/jpeg", "image/jpg":
			fileExt = ".jpg"
		case "image/png":
			fileExt = ".png"
		case "image/gif":
			fileExt = ".gif"
		case "application/pdf":
			fileExt = ".pdf"
		default:
			splitStr := strings.Split(contentType, "/")
			if len(splitStr) == 2 {
				fileExt = fmt.Sprintf(".%s", splitStr[1])
			}
		}
		uploadFileName += fileExt
	}

	fileUrl, err := s.uploader.Upload(r.Context(), uploadFileName, bytes.NewReader(uploadBuff))
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	collectible.UploadFile = fileUrl

	categoryNames, ok := r.Form[formCreateCollectible_Categories]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		if err := (&errResponse{Error: errCategoriesCannotBeEmpty.Error()}).writeResponse(w); err != nil {
			s.zapLogger.Sugar().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	var httpCode int
	err = database.ExecInTx(r.Context(), s.postgres.Pool, func(ctx context.Context, tx pgx.Tx) error {
		quoteToken, err := s.TokenRepo.GetByName(ctx, tx, collectible.QuoteToken.Name)
		if err != nil {
			httpCode = http.StatusInternalServerError
			return errors.Wrap(err, "s.TokenRepo.GetByName")
		}
		if quoteToken == nil {
			httpCode = http.StatusBadRequest
			return errors.New("invalid quote token name")
		}
		collectible.QuoteToken = quoteToken

		err = s.CollectibleRepo.Insert(ctx, tx, collectible)
		if err != nil {
			httpCode = http.StatusInternalServerError
			return errors.Wrap(err, "s.CollectibleRepo.Insert")
		}

		categories, err := s.CategoryRepo.GetCategoryByNames(ctx, tx, categoryNames)
		if err != nil {
			httpCode = http.StatusInternalServerError
			return errors.Wrap(err, "s.CategoryRepo.GetCategoryByNames")
		}
		if len(categories) == 0 || len(categories) != len(categoryNames) {
			httpCode = http.StatusBadRequest
			return errCategoryIsNotValid
		}

		for _, category := range categories {
			_, err := s.CollectibleCategoryRepo.Insert(ctx, tx, collectible, category)
			if err != nil {
				httpCode = http.StatusInternalServerError
				return errors.Wrap(err, "s.CollectibleCategoryRepo.Insert")
			}
			collectible.Categories = append(collectible.Categories, category)
		}
		return nil
	})

	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(httpCode)
		if httpCode == http.StatusBadRequest {
			if err := (&errResponse{Error: err.Error()}).writeResponse(w); err != nil {
				s.zapLogger.Sugar().Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		return
	}

	resp := &CreateCollectibleResponse{ID: collectible.GUID.String()}
	respData, err := json.Marshal(resp)
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(respData); err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.zapLogger.Sugar().Infow(
		"serve create nft successfully",
		"addr", r.RemoteAddr,
		"req", fmt.Sprintf("%+v", collectible),
	)
}

func (s *ApiServer) ResellUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	collectible_id := r.Form.Get("collectible_id")
	account := r.Form.Get("account")
	quoteTokenName := r.Form.Get("quote_token")
	if quoteTokenName == "" {
		quoteTokenName = "CONT"
		//return nil, errTitleCannotBeEmpty
	}

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"collectible_id", r.Form.Get("collectible_id"),
			"account", r.Form.Get("account"),
			"quote_token", r.Form.Get("quote_token"),
			"price", r.Form.Get("instant_sale_price"),
		)
	}()

	instantSalePrice, err := decimal.NewFromString(r.Form.Get("instant_sale_price"))
	if err != nil {
		s.zapLogger.Sugar().Warn(errInstantSalePriceIsNotValid)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uid, err := uuid.FromString(collectible_id)
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

	quoteToken, err := s.TokenRepo.GetByName(r.Context(), s.postgres.Pool, quoteTokenName)
	if err != nil {
		s.zapLogger.Sugar().Warn("s.TokenRepo.GetByName")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if quoteToken == nil {
		s.zapLogger.Sugar().Warn("invalid quote token name")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	collectible_check, err := s.CollectibleRepo.Get(r.Context(), s.postgres.Pool, uid)
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

	if !strings.EqualFold(account, *collectible_check.TokenOwner) {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = s.CollectibleRepo.UpdateResell(r.Context(), s.postgres.Pool, collectible_check.Id,
		instantSalePrice, quoteToken.ID)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

type GetCollectibleResp struct {
	Id                  uuid.UUID                     `json:"id"`
	Title               string                        `json:"title"`
	Description         string                        `json:"description"`
	UploadFile          string                        `json:"upload_file"`
	RoyaltyPercent      decimal.Decimal               `json:"royalty_percent"`
	UnlockOncePurchased bool                          `json:"unlock_once_purchased"`
	InstantSalePrice    decimal.Decimal               `json:"instant_sale_price"`
	Categories          []string                      `json:"categories"`
	Properties          map[string]string             `json:"properties"`
	TokenId             decimal.NullDecimal           `json:"token_id"`
	TokenOwner          *string                       `json:"token_owner"`
	Token               *string                       `json:"token"`
	Status              int                           `json:"status"`
	View                int                           `json:"view"`
	QuoteToken          string                        `json:"quote_token"`
	Creator_acc         *entities.Account             `json:"creator_acc"`
	Owner               *entities.Account             `json:"owner"`
	Like                *entities.CollectibleLikeResp `json:"like"`
}

func NewGetCollectibleResp(collectible *entities.Collectible) *GetCollectibleResp {
	categoryNames := make([]string, 0, len(collectible.Categories))
	for _, categoryName := range collectible.Categories {
		categoryNames = append(categoryNames, categoryName.Name)
	}

	//resp := collectible
	resp := &GetCollectibleResp{
		Id:                  collectible.GUID,
		Title:               collectible.Title,
		Description:         collectible.Description,
		UploadFile:          collectible.UploadFile,
		RoyaltyPercent:      collectible.Royalties.Mul(decimal.NewFromInt(100)),
		UnlockOncePurchased: collectible.UnlockOncePurchased,
		InstantSalePrice:    collectible.InstantSalePrice,
		Categories:          categoryNames,
		Properties:          collectible.Properties,
		TokenId:             collectible.TokenId,
		TokenOwner:          collectible.TokenOwner,
		View:                collectible.View,
		Token:               collectible.Token,
		Status:              collectible.Status,
		QuoteToken:          collectible.QuoteToken.Name,
		Creator_acc:         collectible.Creator_acc,
		Owner:               collectible.Owner,
		Like:                collectible.Like,
	}
	return resp
}

func (s *ApiServer) getCollectible(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			err := (&errResponse{Error: errIdCannotBeEmpty.Error()}).writeResponse(w)
			if err != nil {
				s.zapLogger.Sugar().Warn(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		uid, err := uuid.FromString(id)
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

		address := r.Form.Get("address")

		account := &entities.Account{}
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

		collectible, err := s.CollectibleRepo.Get(r.Context(), s.postgres.Pool, uid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var resp []byte
		if collectible == nil {
			resp, err = json.Marshal(make(map[string]string))
			if err != nil {
				s.zapLogger.Sugar().Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {

			//fmt.Println(collectible.View)

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

			resp, err = json.Marshal(NewGetCollectibleResp(collectible))
			if err != nil {
				s.zapLogger.Sugar().Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if _, err := w.Write(resp); err != nil {
			s.zapLogger.Sugar().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.zapLogger.Sugar().Infow(
			"serve get nft successfully",
			"addr", r.RemoteAddr,
			//"req", fmt.Sprintf("%+v", collectible),
		)
	}
}

func (s *ApiServer) GetCollectible(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := r.Form.Get("id")
	s.getCollectible(id)(w, r)
}

func (s *ApiServer) GetCollectibleMux(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	id := vars["id"]
	s.getCollectible(id)(w, r)
}

type CollectiblePagingResp struct {
	Collectibles []*entities.Collectible `json:"collectibles"`
	PrevCursor   string                  `json:"prev_cursor"`
	NextCursor   string                  `json:"next_cursor"`
}

func (s *ApiServer) CollectiblePaging(w http.ResponseWriter, r *http.Request) {
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

	filter := r.Form.Get("filter")
	if filter == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var cursor uuid.UUID
	var err error
	switch cursorParam := r.Form.Get("cursor"); cursorParam {
	case "", "null":
		cursor = uuid.Nil
	default:
		cursor, err = uuid.FromString(cursorParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	sort := r.Form.Get("sort")
	if sort == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	limitParam := r.Form.Get("limit")
	if limitParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	options := r.Form.Get("options")

	address := r.Form.Get("address")

	status, _ := strconv.Atoi(r.Form.Get("status"))

	account := &entities.Account{}

	if address != "" {
		account, err = s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, address)
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

	categories := r.Form["categories"]

	paging := &repositories.CollectiblePaging{
		Cursor:     cursor,
		Filter:     filter,
		Sort:       sort,
		Limit:      int(limit),
		Title:      r.Form.Get("title"),
		Categories: categories,
		Options:    options,
		Address:    address,
		Status:     status,
	}
	collectibles, err := s.CollectibleRepo.Paging(r.Context(), s.postgres.Pool, paging)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.zapLogger.Sugar().Error(err)
		return
	}

	account_id, err := strconv.ParseInt("-1", 10, 64)
	if err != nil {
		fmt.Printf("%d of type %T", account_id, account_id)
	}

	if account != nil {
		account_id = account.Id
	}

	for index, collectible := range collectibles {
		//fmt.Println(index)

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

		collectibles[index] = collectible

	}

	resp := &CollectiblePagingResp{
		Collectibles: collectibles,
	}
	if cursor != uuid.Nil {
		resp.PrevCursor = cursor.String()
	}
	if len(collectibles) > 0 && len(collectibles) == int(limit) {
		resp.NextCursor = collectibles[len(collectibles)-1].GUID.String()
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

func (s *ApiServer) GetByName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			status, resp := s.getByName(w, r)
			w.WriteHeader(status)
			if resp != nil {
				respData, err := json.Marshal(resp)
				if err != nil {
					s.zapLogger.Sugar().Error(err)
					return
				}
				_, err = w.Write(respData)
				if err != nil {
					s.zapLogger.Sugar().Error(err)
					return
				}
			}
		default:

		}
	}
}

func (s *ApiServer) getByName(w http.ResponseWriter, r *http.Request) (int, interface{}) {
	if err := r.ParseForm(); err != nil {
		s.zapLogger.Sugar().Error(err)
		return http.StatusBadRequest, nil
	}

	limit, err := strconv.ParseInt(r.Form.Get("limit"), 10, 64)
	if err != nil {
		return http.StatusBadRequest, nil
	}
	if limit <= 0 {
		return http.StatusBadRequest, nil
	}
	if limit > 500 {
		limit = 500
	}

	collectibles, err := s.CollectibleRepo.GetByName(r.Context(), s.postgres.Pool, r.Form.Get("name"), int(limit))
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, collectibles
}

func (s *ApiServer) TotalMint() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.Method {
		case http.MethodGet:
			status, resp, err := s.totalMint(writer, request)
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

func (s *ApiServer) totalMint(w http.ResponseWriter, r *http.Request) (int, []byte, error) {
	if err := r.ParseForm(); err != nil {
		s.zapLogger.Sugar().Error(err)
		return http.StatusBadRequest, nil, err
	}

	durationStr := r.Form.Get("duration")

	var duration time.Duration
	switch durationStr {
	case "24h", "24H":
		duration = 24 * time.Hour
	default:
		duration = 0
	}

	totalMint, err := s.CollectibleRepo.GetTotalInDuration(r.Context(), s.postgres.Pool, duration)

	if err != nil {
		s.zapLogger.Sugar().Error(err)
		return http.StatusInternalServerError, nil, err
	}

	resp := &TotalResp{
		Total: totalMint,
	}

	respData, err := json.Marshal(resp)
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, respData, nil
}

func (s *ApiServer) DoLike(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	account_id := r.Form.Get("account_id")
	collectible_id := r.Form.Get("collectible_id")
	action, err := strconv.ParseInt(r.Form.Get("action"), 10, 64)

	defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"account_id", account_id,
			"collectible_id", collectible_id,
			"action", action,
		)
	}()

	if account_id == "" || collectible_id == "" {
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

	collectblelike_check, err := s.CollectibleLikeRepo.Check(r.Context(), s.postgres.Pool, collectible_check.Id, account_check.Id)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//fmt.Println("total_like", total_like)

	if action == 1 {

		account_owner, err := s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, *collectible_check.TokenOwner)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var account_id sql.NullInt64
		account_id.Int64 = account_check.Id
		account_id.Valid = true
		var collectible_id sql.NullInt64
		collectible_id.Int64 = collectible_check.Id
		collectible_id.Valid = true

		notice := &entities.Notice{
			AccountID:     account_owner.Id,
			FromAccountID: account_id,
			CollectibleID: collectible_id,
			Content:       entities.Notice_like_nft,
		}

		//fmt.Println("%v", notice)

		err = s.NoticeRepo.Insert(r.Context(), s.postgres.Pool, notice)

		if collectblelike_check == true {
			w.WriteHeader(http.StatusOK)
			return

		} else {
			err := s.CollectibleLikeRepo.Insert(r.Context(), s.postgres.Pool, collectible_check.Id, account_check.Id)
			if err != nil {
				s.zapLogger.Sugar().Warn(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			total_like := s.CollectibleLikeRepo.Count(r.Context(), s.postgres.Pool, collectible_check.Id)
			(&repositories.CollectibleRepository{}).UpdateTotalLike(r.Context(), s.postgres.Pool, collectible_check.Id, total_like)

			w.WriteHeader(http.StatusOK)
			return

		}

	} else if action == 0 { //dislike
		if collectblelike_check == true {
			err := s.CollectibleLikeRepo.Delete(r.Context(), s.postgres.Pool, collectible_check.Id, account_check.Id)
			if err != nil {
				s.zapLogger.Sugar().Warn(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			total_like := s.CollectibleLikeRepo.Count(r.Context(), s.postgres.Pool, collectible_check.Id)
			(&repositories.CollectibleRepository{}).UpdateTotalLike(r.Context(), s.postgres.Pool, collectible_check.Id, total_like)

			w.WriteHeader(http.StatusOK)
			return

		} else {
			w.WriteHeader(http.StatusOK)
			return

		}
	}

}

func (s *ApiServer) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := r.Form.Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := (&errResponse{Error: errIdCannotBeEmpty.Error()}).writeResponse(w)
		if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	uid, err := uuid.FromString(id)
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
			"collectible_id", uid,
		)
	}()

	collectible, err := s.CollectibleRepo.Get(r.Context(), s.postgres.Pool, uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.CollectibleRepo.Delete(r.Context(), s.postgres.Pool, collectible.Id)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}

func (s *ApiServer) UpdateView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	collectible_id := r.Form.Get("collectible_id")
	/*defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			"collectible_id", collectible_id,
		)
	}()*/
	uid, err := uuid.FromString(collectible_id)

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

	collectible_check, err := s.CollectibleRepo.Get(r.Context(), s.postgres.Pool, uid)
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

	err = s.CollectibleRepo.UpdateView(r.Context(), s.postgres.Pool, collectible_check.Id)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}

func (s *ApiServer) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	collectible_id := r.Form.Get("collectible_id")
	status, err := strconv.Atoi(r.Form.Get("status"))
	uid, err := uuid.FromString(collectible_id)
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
			"collectible_id", uid,
			"status", status,
		)
	}()

	collectible_check, err := s.CollectibleRepo.Get(r.Context(), s.postgres.Pool, uid)
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

	err = s.CollectibleRepo.UpdateStatus(r.Context(), s.postgres.Pool, collectible_check.Id, status)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}

type HistoryResp struct {
	TokenEvents []*entities.TokenEvent `json:"history"`
}

func (s *ApiServer) GetHistory(w http.ResponseWriter, r *http.Request) {
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

	events, err := (&repositories.ExchangeEventRepo{}).GetHistory(r.Context(), s.postgres.Pool, collectible_check.TokenId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.zapLogger.Sugar().Error(err)
		return
	}

	for index, event := range events {
		//fmt.Println("event:", event)
		//fmt.Println("event.Account:", event.Account)
		from_address := ""
		if event.Account != nil {
			from_address = *event.Account
			//fmt.Println("event.Account:", *event.Account)
		}

		//fmt.Println("from_address:", from_address)
		account, err := s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, from_address)
		if err != nil {

			account = nil
		}
		event.From = account
		events[index] = event

	}

	resp := &HistoryResp{
		TokenEvents: events,
	}

	respData, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.zapLogger.Sugar().Error(err)
		return
	}

	//fmt.Println("respData:", respData)

	_, err = w.Write(respData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.zapLogger.Sugar().Error(err)
		return
	}

	//fmt.Println("logs:", events)

}

func (s *ApiServer) DoBlockCollective(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	collectible_id := r.Form.Get("collectible_id")
	uid, err := uuid.FromString(collectible_id)
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
			"collectible_id", uid,
		)
	}()

	collectible_check, err := s.CollectibleRepo.Get(r.Context(), s.postgres.Pool, uid)
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

	err = s.CollectibleRepo.Block(r.Context(), s.postgres.Pool, collectible_check.Id)
	if err != nil {
		s.zapLogger.Sugar().Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
