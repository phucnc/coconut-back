package services

import (

	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	//"github.com/gofrs/uuid"
	//"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	//"github.com/shopspring/decimal"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"nft-backend/internal/database"
	"nft-backend/internal/entities"
	"nft-backend/internal/repositories"
	//"strconv"
	"strings"
	"database/sql"
	"time"
)

type CreateAccountResponse struct {
	ID int64 `json:"id"`
}

type AccountResp struct {
	Id                 int64           `json:"id"`
	Username               sql.NullString              `json:"username"`
	Address         string              `json:"address"`
	Avatar          sql.NullString              `json:"avartar"`
	Cover      sql.NullString     `json:"cover"`
	Info             sql.NullString `json:"info"`
	Facebook sql.NullString                `json:"facebook"`
	Tiktok    sql.NullString     `json:"tiktok"`
	Twitter    sql.NullString     `json:"twitter"`
	Instagram    sql.NullString     `json:"instagram"`
	Status             int `json:"status"`

}


func NewAccountResp(account *entities.Account) *AccountResp {
	resp := &AccountResp{
		Id:                  account.Id,
		Username:               account.Username,
		Address:         account.Address,
		Avatar:          account.Avatar,
		Cover:      account.Cover,
		Info: account.Info,
		Facebook:    account.Facebook,
		Twitter:          account.Twitter,
		Tiktok:          account.Tiktok,
		Instagram:             account.Instagram,
		Status:          account.Status,
	}
	return resp
}



func (s *ApiServer) handleAccountPath() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.CreateOrLogin(w, r)
		case http.MethodPut:
			s.Updateinfo(w, r)
		case http.MethodDelete:
			fmt.Println("DELETE")
			w.WriteHeader(http.StatusOK)
		case http.MethodGet:
			s.GetInfo(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (s *ApiServer) Cover() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			s.UpdateCover(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

}

func (s *ApiServer) Avatar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			s.UpdateAvatar(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}


func (s *ApiServer) Block() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			s.DoBlock(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

var (
	errAddressCannotBeEmpty            = errors.New("Address cannot be empty")
	errUpdateInfoEmpty            = errors.New("UpdateInfo empty")
	errIdInvalid			= errors.New("Id invalid")
	errNotFound	= errors.New("Not found")
	errUserNameExist = errors.New("Username exist")
)

const (
	formAccount_Address               = "address"
	formAccount_Username        = "username"
	formAccount_Info        = "info"
	formAccount_avatar          = "avatar"
	formAccount_cover          = "cover"
	formAccount_Id         = "id"
	formAccount_UploadFile = "upload_file"
)


func newAccountFromForm(form url.Values) (*entities.Account, error) {
	address := form.Get(formAccount_Address) 
	if address == "" {
		return nil, errAddressCannotBeEmpty
	}

	account := &entities.Account{
		Address:               form.Get(formAccount_Address),
	}

	fmt.Printf("%+v\n",account)

	return account, nil
}



func (s *ApiServer) CreateOrLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    
  	if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
	address := r.Form.Get(formAccount_Address)

   defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			formAccount_Address, address,
		)
	}()

	if address == "" {
			w.WriteHeader(http.StatusBadRequest)
			err := (&errResponse{Error: errAddressCannotBeEmpty.Error()}).writeResponse(w)
			if err != nil {
				s.zapLogger.Sugar().Warn(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

	account_check, err :=  s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, address)
	var resp []byte
 
	if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}
	

	if(account_check != nil) {
		
       resp, err = json.Marshal(NewAccountResp(account_check))
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

		s.zapLogger.Sugar().Infow(
			"Login successful",
			"addr", r.RemoteAddr,
			"req", fmt.Sprintf("%+v", account_check),
		)
	} else {

	account, err := newAccountFromForm(r.Form)
	var httpCode int
	err = database.ExecInTx(r.Context(), s.postgres.Pool, func(ctx context.Context, tx pgx.Tx) error {
		err = s.AccountRepo.Insert(ctx, tx, account)
		if err != nil {
			httpCode = http.StatusInternalServerError
			return errors.Wrap(err, "s.AccountRepo.Insert")
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

    account_check, err =  s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, address)
    resp, err = json.Marshal(NewAccountResp(account_check))

    if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}


    if(account_check != nil) {
		

       resp, err = json.Marshal(NewAccountResp(account_check))
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

		s.zapLogger.Sugar().Infow(
			"Create account successful",
			"addr", r.RemoteAddr,
			"req", fmt.Sprintf("%+v", account_check),
		)
	} 

}
}

func (s *ApiServer) GetInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    
  	if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
	id := r.Form.Get(formAccount_Id)



   defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			formAccount_Id, id,
		)
	}()


	account_check, err :=  s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, id)
 
	if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}
	
    if(account_check != nil) {
		
       resp, err := json.Marshal(NewAccountResp(account_check))
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

		s.zapLogger.Sugar().Infow(
			"Get Account info successful",
			"addr", r.RemoteAddr,
		)
	  }  else {
	  	w.WriteHeader(http.StatusBadRequest)
			err := (&errResponse{Error: errNotFound.Error()}).writeResponse(w)
			if err != nil {
				s.zapLogger.Sugar().Warn(err)
				w.WriteHeader(http.StatusNotFound )
				return
			}
			return
	  }

	
}


func (s *ApiServer) Updateinfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    
  	if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
	username := r.Form.Get(formAccount_Username)
	info := r.Form.Get(formAccount_Info)
	id := r.Form.Get(formAccount_Id)


   defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			formAccount_Id, id,
			formAccount_Username, username,
			formAccount_Info, info,
		)
	}()

	if username == "" && info == "" {
			w.WriteHeader(http.StatusBadRequest)
			err := (&errResponse{Error: errUpdateInfoEmpty.Error()}).writeResponse(w)
			if err != nil {
				s.zapLogger.Sugar().Warn(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

	if(username != "") {

		check, err :=  s.AccountRepo.CheckUserName(r.Context(), s.postgres.Pool, id, username)
		if err != nil {
				s.zapLogger.Sugar().Warn(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		if(check == false) {

			w.WriteHeader(http.StatusBadRequest)
			err := (&errResponse{Error: errUserNameExist.Error()}).writeResponse(w)
			if err != nil {
				s.zapLogger.Sugar().Warn(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return

		}
	}


	account_check, err :=  s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, id)
	
	if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}
	

	if(account_check == nil) {

		w.WriteHeader(http.StatusBadRequest)
			err := (&errResponse{Error: errNotFound.Error()}).writeResponse(w)
			if err != nil {
				s.zapLogger.Sugar().Warn(err)
				w.WriteHeader(http.StatusNotFound )
				return
			}
			return
		

	} else {


		err :=  s.AccountRepo.Update(r.Context(), s.postgres.Pool, id, username, info)

        if err != nil {
				s.zapLogger.Sugar().Warn(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			
		account_updated, err :=  s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, id)


    if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}

    if(account_updated != nil) {
		
       resp, err := json.Marshal(NewAccountResp(account_updated))
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

		s.zapLogger.Sugar().Infow(
			"Update Account info successful",
			"addr", r.RemoteAddr,
		)
	  } 

	}
}


func (s *ApiServer) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
		defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			formAccount_Id, r.Form.Get(formAccount_Id),
		)
	}()

	switch err := r.ParseMultipartForm(2 * megabyte); err {
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
    

	id := r.Form.Get(formAccount_Id);


   defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			formAccount_Id, id,
		)
	}()


	account_check, err :=  s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, id)
	if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}
	
	if(account_check == nil) {

		w.WriteHeader(http.StatusBadRequest)
			err := (&errResponse{Error: errNotFound.Error()}).writeResponse(w)
			if err != nil {
				s.zapLogger.Sugar().Warn(err)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			return
		

	} else {
     	file, _, err := r.FormFile(formAccount_UploadFile)
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

	/*file_id, err := uuid.NewV4()
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}*/
	//uploadFileName := "avatar_"+account_check.Address;

	t := time.Now()
	uploadFileName :=  fmt.Sprintf("avatar%d%d%d%d%d_%s",t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(),account_check.Address)


	if len(uploadBuff) >= 512 {
		contentType := http.DetectContentType(uploadBuff[:512])
		var fileExt string
		switch contentType {
		case "image/jpeg", "image/jpg":
			fileExt = ".jpg"
		case "image/png":
			fileExt = ".png"
		default:
			splitStr := strings.Split(contentType, "/")
			if len(splitStr) == 2 {
				fileExt = fmt.Sprintf(".%s", splitStr[1])
			}
		}
		uploadFileName += fileExt
	}
	fmt.Println("File name "+ uploadFileName);
	fileUrl, err := s.uploader.Upload(r.Context(), uploadFileName, bytes.NewReader(uploadBuff))
	fmt.Println("File URL "+ fileUrl);
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	 err = s.AccountRepo.UpdateAvatar(r.Context(), s.postgres.Pool, id, fileUrl);

	if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		} else {

	account_updated, err :=  s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, id)

    if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}

    if(account_updated != nil) {
		
       resp, err := json.Marshal(NewAccountResp(account_updated))
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

		s.zapLogger.Sugar().Infow(
			"Update Avatar info successful",
			"addr", r.RemoteAddr,
		)
	  }

	}

	}
}


func (s *ApiServer) UpdateCover(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
		defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			formAccount_Id, r.Form.Get(formAccount_Id),
		)
	}()

	switch err := r.ParseMultipartForm(2 * megabyte); err {
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
    
	id := r.Form.Get(formAccount_Id)

   defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			formAccount_Id, id,
		)
	}()


	account_check, err :=  s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, id)
	if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}
	
	if(account_check == nil) {

		w.WriteHeader(http.StatusBadRequest)
			err := (&errResponse{Error: errNotFound.Error()}).writeResponse(w)
			if err != nil {
				s.zapLogger.Sugar().Warn(err)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			return
		

	} else {
     	file, _, err := r.FormFile(formAccount_UploadFile)
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

	/*file_id, err := uuid.NewV4()
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}*/
	//uploadFileName := "cover_"+account_check.Address;
	t := time.Now()
	uploadFileName :=  fmt.Sprintf("cover%d%d%d%d%d_%s",t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(),account_check.Address)

	if len(uploadBuff) >= 512 {
		contentType := http.DetectContentType(uploadBuff[:512])
		var fileExt string
		switch contentType {
		case "image/jpeg", "image/jpg":
			fileExt = ".jpg"
		case "image/png":
			fileExt = ".png"
		default:
			splitStr := strings.Split(contentType, "/")
			if len(splitStr) == 2 {
				fileExt = fmt.Sprintf(".%s", splitStr[1])
			}
		}
		uploadFileName += fileExt
	}
	fmt.Println("File name "+ uploadFileName);
	fileUrl, err := s.uploader.Upload(r.Context(), uploadFileName, bytes.NewReader(uploadBuff))
	//fmt.Println("File URL "+ fileUrl);
	if err != nil {
		s.zapLogger.Sugar().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	 err = s.AccountRepo.UpdateCover(r.Context(), s.postgres.Pool, id, fileUrl);

	if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		} else {

	account_updated, err :=  s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, id)

    if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}

    if(account_updated != nil) {
		
       resp, err := json.Marshal(NewAccountResp(account_updated))
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

		s.zapLogger.Sugar().Infow(
			"Update Avatar info successful",
			"addr", r.RemoteAddr,
		)
	  }

	}

	}
}

type AccountPagingResp struct {
	Accounts []*entities.Account `json:"accounts"`
}

func (s *ApiServer) AccountSearchPaging(w http.ResponseWriter, r *http.Request) {
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


	paging := &repositories.AccountPaging{
		Limit:      int(limit),
		Offset:     int(offset),
		Keys:       r.Form.Get("keys"),
	}

	accounts, err := s.AccountRepo.SearchPaging(r.Context(), s.postgres.Pool, paging)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.zapLogger.Sugar().Error(err)
		return
	}
	
	resp := &AccountPagingResp{
		Accounts: accounts,
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


func (s *ApiServer) AccountPaging(w http.ResponseWriter, r *http.Request) {
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

	status, _ := strconv.Atoi(r.Form.Get("status"))


	paging := &repositories.AccountPaging{
		Limit:      int(limit),
		Offset:     int(offset),
		Status: status,
	}

	accounts, err := s.AccountRepo.Paging(r.Context(), s.postgres.Pool, paging)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.zapLogger.Sugar().Error(err)
		return
	}
	
	resp := &AccountPagingResp{
		Accounts: accounts,
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



func (s *ApiServer) DoBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
		defer func() {
		s.zapLogger.Sugar().Infow(
			r.URL.Path,
			"method", r.Method,
			formAccount_Id, r.Form.Get(formAccount_Id),
		)
	}()

	if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

	account_id := r.Form.Get(formAccount_Id)
	account_check, err :=  s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, account_id)
	if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}
	
	if(account_check == nil) {

		w.WriteHeader(http.StatusBadRequest)
			err := (&errResponse{Error: errNotFound.Error()}).writeResponse(w)
			if err != nil {
				s.zapLogger.Sugar().Warn(err)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			return
		

	} else {

	 err = s.AccountRepo.Block(r.Context(), s.postgres.Pool, account_id);

	if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		} else {

	account_updated, err :=  s.AccountRepo.GetByAddress(r.Context(), s.postgres.Pool, account_id)

    if err != nil {
			s.zapLogger.Sugar().Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}

    if(account_updated != nil) {
		
       resp, err := json.Marshal(NewAccountResp(account_updated))
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

		s.zapLogger.Sugar().Infow(
			"Block  successful",
			"addr", r.RemoteAddr,
		)
	  }

	}

	}
}
