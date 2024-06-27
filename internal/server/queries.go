package server

import (
	"context"
	"log"
	"net/http"
	"time"
	"web-storage-service/internal/dto"
	"web-storage-service/internal/models"

	"github.com/jackc/pgx/v5"
)

func (s *Server) ListAssetsQuery(ctx context.Context, dto dto.ListAssets) (pgx.Rows, error) {
	query := `SELECT name, data FROM assets WHERE uid=$1 AND deleted=FALSE LIMIT $2 OFFSET $3`
	return s.db.Query(ctx, query, dto.UserID, dto.Limit, dto.Offset)
}

func (s *Server) GetAssetByNameQuery(ctx context.Context, dto dto.GetAssetByName) ([]byte, error) {
	var data []byte
	query := `SELECT data FROM assets WHERE uid=$1 AND name=$2 AND deleted=FALSE`
	err := s.db.QueryRow(ctx, query, dto.UserID, dto.Name).Scan(&data)
	return data, err
}

func (s *Server) UploadAssetQuery(ctx context.Context, name string, userID int, data []byte) error {
	query := `
        INSERT INTO assets (name, uid, data, created_at) 
        VALUES ($1, $2, $3, NOW())
        ON CONFLICT (name, uid) 
        DO NOTHING;
    `
	_, err := s.db.Exec(ctx, query, name, userID, data)
	return err
}

func (s *Server) UpdateAssetQuery(ctx context.Context, name string, userID int, data []byte) error {
	query := `
        INSERT INTO assets (name, uid, data, created_at) 
        VALUES ($1, $2, $3, NOW())
        ON CONFLICT (name, uid) 
        DO UPDATE SET data = EXCLUDED.data, created_at = EXCLUDED.created_at;
    `
	_, err := s.db.Exec(ctx, query, name, userID, data)
	return err
}

func (s *Server) SoftDeleteAssetQuery(ctx context.Context, dto dto.DeleteAsset) error {
	query := `UPDATE assets SET deleted = TRUE WHERE name = $1 AND uid = $2`
	_, err := s.db.Exec(ctx, query, dto.Name, dto.UserID)
	return err
}

func (s *Server) HardDeleteAssetQuery(ctx context.Context, dto dto.DeleteAsset) error {
	query := `DELETE FROM assets WHERE name = $1 AND uid = $2`
	_, err := s.db.Exec(ctx, query, dto.Name, dto.UserID)
	return err
}

func (s *Server) GetUserByLogin(ctx context.Context, credentials dto.Credentials) (models.User, error) {
	var user models.User
	err := s.db.QueryRow(ctx, "SELECT id, login, password_hash FROM users WHERE login=$1", credentials.Login).Scan(&user.ID, &user.Login, &user.PasswordHash)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *Server) CreateNewSessionQueryWithPostgresTrigger(ctx context.Context, dto dto.CreateNewSession) error {
	query := `INSERT INTO sessions (id, uid, expires_at, ip_address, active) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.Exec(ctx, query, dto.Token, dto.UserID, dto.ExpiresAt, dto.IPAddress, false)
	return err
}

func (s *Server) CreateNewSessionWithTransaction(ctx context.Context, dto dto.CreateNewSession) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Printf("tx rollback error: %v", rollbackErr)
			}
		}
	}()

	// Устанавливаем все предыдущие сессии как неактивные
	_, err = tx.Exec(ctx, "UPDATE sessions SET active = FALSE WHERE uid = $1", dto.UserID)
	if err != nil {
		return err
	}

	// Вставляем новую сессию как активную
	_, err = tx.Exec(ctx, "INSERT INTO sessions (id, uid, expires_at, ip_address, active) VALUES ($1, $2, $3, $4, TRUE)", dto.Token, dto.UserID, dto.ExpiresAt, dto.IPAddress)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Server) DeactivateExpiredSessionsAndReturnUserID(token string) (userID int, err error) {
	var expiresAt time.Time
	var active bool

	err = s.db.QueryRow(context.Background(), "SELECT uid, expires_at, active FROM sessions WHERE id=$1", token).Scan(&userID, &expiresAt, &active)
	if err != nil {
		return 0, err
	}

	// Если сессия активна и еще не истекла
	if active && time.Now().Before(expiresAt) {
		return userID, nil
	}

	// Если сессия по времени истекла, но статус active = TRUE - делаем сессию неактивной
	if active && time.Now().After(expiresAt) {
		_, err = s.db.Exec(context.Background(), "UPDATE sessions SET active = FALSE WHERE id = $1", token)
		if err != nil {
			log.Printf("Error deactivating session: %v", err)
			return 0, err
		}
		return 0, http.ErrNoCookie
	}

	return 0, http.ErrNoCookie
}
