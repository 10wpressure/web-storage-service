package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"web-storage-service/internal/dto"
	"web-storage-service/pkg"

	_ "github.com/jackc/pgx/v5/pgxpool"
)

func (s *Server) AuthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var credentials = dto.Credentials{}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		BadRequestError(w)
		return
	}

	user, err := s.GetUserByLogin(ctx, credentials)
	if err != nil || !pkg.ValidatePassword(credentials.Password, user.PasswordHash) {
		UnauthorizedError(w, "invalid login/password")
		return
	}

	token := pkg.GenerateToken(user.Login)
	err = s.CreateNewSessionWithTransaction(ctx, dto.CreateNewSession{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IPAddress: strings.Split(r.RemoteAddr, ":")[0],
	})
	if err != nil {
		log.Printf("error creating new session: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"token": token})
	if err != nil {
		InternalServerError(w)
		return
	}
}

func (s *Server) ListAssetsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(UserIDKey).(int)

	// Получение параметров пагинации
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 10
	}

	offset := (page - 1) * size

	rows, err := s.ListAssetsQuery(ctx, dto.ListAssets{
		UserID: userID,
		Offset: offset,
		Limit:  size,
	})
	if err != nil {
		InternalServerError(w)
		return
	}
	defer rows.Close()

	assets := make(map[string]string)
	for rows.Next() {
		var name string
		var data []byte
		if err = rows.Scan(&name, &data); err != nil {
			InternalServerError(w)
			return
		}
		assets[name] = pkg.TrimData(data)
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"page":    page,
		"size":    size,
		"assets":  assets,
		"hasMore": len(assets) == size,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		InternalServerError(w)
		return
	}
}

func (s *Server) DownloadAssetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(UserIDKey).(int)
	name := r.PathValue("name")
	if name == "" {
		BadRequestError(w)
		return
	}

	data, err := s.GetAssetByNameQuery(ctx, dto.GetAssetByName{UserID: userID, Name: name})
	if err != nil {
		NotFoundError(w, "asset")
		return
	}

	assets := make(map[string]string)
	assets[name] = pkg.TrimData(data)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(assets)
	if err != nil {
		InternalServerError(w)
		return
	}
}

func (s *Server) UploadAssetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(UserIDKey).(int)
	name := r.PathValue("name")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = s.UploadAssetQuery(ctx, name, userID, data)
	if err != nil {
		InternalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	if err != nil {
		InternalServerError(w)
		return
	}
}

func (s *Server) UpdateAssetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(UserIDKey).(int)
	name := r.PathValue("name")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = s.UpdateAssetQuery(ctx, name, userID, data)
	if err != nil {
		InternalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	if err != nil {
		InternalServerError(w)
		return
	}
}

func (s *Server) SoftDeleteAssetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(UserIDKey).(int)
	name := r.PathValue("name")

	err := s.SoftDeleteAssetQuery(ctx, dto.DeleteAsset{
		Name:   name,
		UserID: userID,
	})
	if err != nil {
		InternalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{strconv.Itoa(userID): name, "status": "soft deleted"})
	if err != nil {
		InternalServerError(w)
		return
	}
}

func (s *Server) HardDeleteAssetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(UserIDKey).(int)
	name := r.PathValue("name")

	err := s.HardDeleteAssetQuery(ctx, dto.DeleteAsset{
		Name:   name,
		UserID: userID,
	})
	if err != nil {
		InternalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{strconv.Itoa(userID): name, "status": "hard deleted"})
	if err != nil {
		InternalServerError(w)
		return
	}
}
func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
