package http_server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"io"
	"net/http"
	"passwordvault/internal/storage/file_store"
	"strings"
)

type DownloadFileReqBody struct {
	Filename string `json:"filename"`
}

func (s *HttpServer) WithLoggingHTTP(handlerFunc runtime.HandlerFunc) runtime.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request, m map[string]string) {
		s.logger.Sugar().Infof("URL: %s; Content-Disposition: %s", req.URL.Path, req.Header.Get("Content-Disposition"))
		handlerFunc(rw, req, m)
	}
}

func (s *HttpServer) WithCheckCredentials(handlerFunc runtime.HandlerFunc) runtime.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request, m map[string]string) {
		scheme, token, found := strings.Cut(req.Header.Get("Authorization"), " ")
		if strings.ToLower(scheme) != "bearer" || !found {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId, err := s.db.UserCheckLoggedIn(token)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		r := req.WithContext(context.WithValue(req.Context(), "LoggedUserId", userId))

		handlerFunc(rw, r, m)
	}
}

func (s *HttpServer) DownloadFile(w http.ResponseWriter, r *http.Request, dummy map[string]string) {

	fileStoreId := strings.Replace(r.Context().Value("LoggedUserId").(string), "-", "", -1)

	pBody := DownloadFileReqBody{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading body: " + err.Error()))
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, &pBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error parsing body: " + err.Error()))
		return
	}
	filename := pBody.Filename

	key, err := s.db.GetFileStoreKey(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error getting file storage key: %v", err)))
		return
	}

	ms, err := file_store.NewMinioStorage(r.Context(), s.config.MinioEndPoint, "securestorageservice", fileStoreId, key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error creating minio storage: %v", err)))
	}
	obj, err := ms.Download(r.Context(), fmt.Sprintf(`%s\%s`, fileStoreId, filename))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("content-type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%s`, filename))

	_, err = io.Copy(w, obj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
}

func (s *HttpServer) UploadFile(w http.ResponseWriter, r *http.Request, dummy map[string]string) {

	fileStoreId := strings.Replace(r.Context().Value("LoggedUserId").(string), "-", "", -1)

	_, filename, found := strings.Cut(r.Header.Get("Content-Disposition"), "attachment; filename=")
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing attachment filename"))
		return
	}

	key, err := s.db.GetFileStoreKey(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error getting file storage key: %v", err)))
		return
	}

	ms, err := file_store.NewMinioStorage(r.Context(), s.config.MinioEndPoint, "securestorageservice", fileStoreId, key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	err = ms.Upload(r.Context(), r.Body, fmt.Sprintf(`%s\%s`, fileStoreId, filename))
	if err != nil {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

}
