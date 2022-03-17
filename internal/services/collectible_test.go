package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"nft-backend/internal/database"
	"nft-backend/internal/entities"
	"nft-backend/internal/repositories"
	"nft-backend/internal/storage"
	"nft-backend/postgres"
	"os"
	"testing"
	"time"
)

func createMultipartFormData(t *require.Assertions, params map[string][]string, fieldName, fileName string) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer
	file := mustOpen(fileName)
	fw, err = w.CreateFormFile(fieldName, file.Name())
	t.NoError(err)

	_, err = io.Copy(fw, file)
	t.NoError(err)

	for key, values := range params {
		for _, value := range values {
			t.NoError(w.WriteField(key, value))
		}
	}

	_ = w.Close()
	return b, w
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Println("PWD: ", pwd)
		panic(err)
	}
	return r
}

func TestCreateCollectible(t *testing.T) {
	test := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config := database.NewPostgreConfigFromEnv()
	pg, err := postgres.NewPostgres(ctx, config.ToConnStr())
	test.NoError(err)

	uploader, err := storage.NewUploader(ctx, storage.LoadConfigFromEnv())
	test.NoError(err)

	s := NewServer(ctx, os.Getenv("BACKEND_PORT"), zap.NewNop(), pg, uploader)
	defer func() {
		_ = s.Shutdown(ctx)
	}()

	time.Sleep(50 * time.Millisecond)

	////
	params := map[string][]string{
		"unlock_once_purchased": {"true"},
		"instant_sale_price":    {"12"},
		"royalty_percent":       {"30"},
		"title":                 {"test-title"},
		"description":           {"test-description"},
		"categories":            {"art", "game"},
	}
	b, w := createMultipartFormData(test, params, "upload_file", "../../test-files/test.txt")

	req, err := http.NewRequest("POST", "http://localhost:12000/nft", &b)
	test.NoError(err)
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Simulate client calls
	resp, err := (&http.Client{}).Do(req)
	test.NoError(err)
	test.Equal(http.StatusOK, resp.StatusCode)
	defer func() {
		_ = resp.Body.Close()
	}()

	// Check resp can decode
	createCollectibleResp := &CreateCollectibleResponse{}
	test.NoError(json.NewDecoder(resp.Body).Decode(createCollectibleResp))

	id, err := uuid.FromString(createCollectibleResp.ID)
	test.NoError(err)

	insertedCollectible, err := (&repositories.CollectibleRepository{}).Get(ctx, pg.Pool, id)
	test.NoError(err)
	//fmt.Println(insertedCollectible)

	test.Equal(params[formCreateCollectible_Title][0], insertedCollectible.Title)
	test.Equal(params[formCreateCollectible_Description][0], insertedCollectible.Description)
	//test.Equal(params[formCreateCollectible_RoyaltyPercent][0], insertedCollectible.Royalties.Mul(decimal.NewFromInt(100)).String())
	//test.Equal(params[formCreateCollectible_UnlockOncePurchased][0], strconv.FormatBool(insertedCollectible.UnlockOncePurchased))
	test.NotEmpty(insertedCollectible.UploadFile)
	test.Equal(len(params[formCreateCollectible_Categories]), len(insertedCollectible.Categories))
	categoryNames := make(map[string]bool)
	for _, categoryName := range params[formCreateCollectible_Categories] {
		categoryNames[categoryName] = true
	}
	for _, category := range insertedCollectible.Categories {
		_, ok := categoryNames[category.Name]
		test.True(ok)
	}

	test.Equal(createCollectibleResp.ID, insertedCollectible.GUID.String())
}

func TestGetCollectible(t *testing.T) {
	test := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	zapLogger, _ := zap.NewDevelopment()

	config := database.NewPostgreConfigFromEnv()
	pg, err := postgres.NewPostgres(ctx, config.ToConnStr())
	test.NoError(err)

	uploader, err := storage.NewUploader(ctx, storage.LoadConfigFromEnv())
	test.NoError(err)

	s := NewServer(ctx, os.Getenv("BACKEND_PORT"), zapLogger, pg, uploader)
	defer func() {
		_ = s.Shutdown(ctx)
	}()

	time.Sleep(100 * time.Millisecond)

	insertCollectible := &entities.Collectible{
		UnlockOncePurchased: true,
		Title:               "test-name",
		Description:         "test-description",
		Royalties:           decimal.NewFromFloat(0.4),
		InstantSalePrice:    decimal.NewFromInt(16),
		QuoteToken: &entities.Token{
			ID: 1,
		},
	}
	test.NoError((&repositories.CollectibleRepository{}).Insert(ctx, pg.Pool, insertCollectible))
	//insertCollectible.Royalties = decimal.NewFromBigInt(insertCollectible.Royalties.BigInt(), 0)

	categories, err := (&repositories.CategoryRepo{}).GetCategoryByNames(ctx, pg.Pool, []string{"art"})
	for _, category := range categories {
		_, err = (&repositories.CollectibleCategoryRepo{}).Insert(ctx, pg.Pool, insertCollectible, category)
		test.NoError(err)
	}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:12000/nft", nil)
	test.NoError(err)
	params := url.Values{}
	params.Set("id", insertCollectible.GUID.String())
	req.URL.RawQuery = params.Encode()

	resp, err := (&http.Client{}).Do(req)
	test.NoError(err)
	defer func() {
		_ = resp.Body.Close()
	}()

	test.Equal(http.StatusOK, resp.StatusCode)

	getCollectibleResp := &GetCollectibleResp{}
	test.NoError(json.NewDecoder(resp.Body).Decode(getCollectibleResp))

	//fmt.Println(fmt.Sprintf("%+v", getCollectibleResp))
	//fmt.Println(fmt.Sprintf("%+v", insertCollectible))

	test.Equal(insertCollectible.GUID, getCollectibleResp.Id)
	test.Equal(insertCollectible.Title, getCollectibleResp.Title)
	test.Equal(insertCollectible.Description, getCollectibleResp.Description)
	test.Equal(insertCollectible.UploadFile, getCollectibleResp.UploadFile)
	test.Equal(insertCollectible.UnlockOncePurchased, getCollectibleResp.UnlockOncePurchased)
	test.Equal(insertCollectible.InstantSalePrice.String(), getCollectibleResp.InstantSalePrice.String())
	test.Equal(insertCollectible.Royalties.String(), getCollectibleResp.RoyaltyPercent.Div(decimal.NewFromInt(100)).String())
	test.Equal(len(categories), len(getCollectibleResp.Categories))

	categoryNames := make(map[string]bool)
	for _, category := range categories {
		categoryNames[category.Name] = true
	}

	for _, category := range getCollectibleResp.Categories {
		_, ok := categoryNames[category]
		test.True(ok)
	}
}
