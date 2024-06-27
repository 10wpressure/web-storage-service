package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {

	r := http.NewServeMux()

	r.HandleFunc("GET /health", s.HealthHandler) // TODO: оставить доступ только для админов

	r.HandleFunc("POST /api/auth", s.AuthHandler)

	r.Handle("POST /api/upload-asset/{name}", s.AuthMiddleware(http.HandlerFunc(s.UploadAssetHandler)))
	r.Handle("PUT /api/update-asset/{name}", s.AuthMiddleware(http.HandlerFunc(s.UpdateAssetHandler)))
	r.Handle("GET /api/asset/{name}", s.AuthMiddleware(http.HandlerFunc(s.DownloadAssetHandler)))
	r.Handle("PUT /api/delete-asset/{name}", s.AuthMiddleware(http.HandlerFunc(s.SoftDeleteAssetHandler)))
	r.Handle("DELETE /api/delete-asset/{name}", s.AuthMiddleware(http.HandlerFunc(s.HardDeleteAssetHandler)))
	r.Handle("GET /api/assets", s.AuthMiddleware(http.HandlerFunc(s.ListAssetsHandler)))

	return r
}
