package auth

import (
	"database/sql"
	"strings"
	"time"

	"github.com/evermos/boilerplate-go/infras"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	Register(user *User) error
	GetUserByUsername(username string) (*Access, error)
	IsExist(username string) (bool, error)
}

type AuthRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideAuthRepositoryMySQL(db *infras.MySQLConn) *AuthRepositoryMySQL {
	return &AuthRepositoryMySQL{
		DB: db,
	}
}

func (r *AuthRepositoryMySQL) GetUserByUsername(username string) (*Access, error) {
	query := "SELECT id, username, password, role FROM ums_users WHERE username = ? LIMIT 1"

	var access Access
	err := r.DB.Read.Get(&access, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Msg("No user found")
			return nil, err
		}
		log.Error().Err(err).Msg("Failed to get user by username")
		return nil, err
	}
	return &access, nil
}

func (r *AuthRepositoryMySQL) IsExist(username string) (bool, error) {
	query := "SELECT EXISTS(SELECT username FROM ums_users WHERE username = ? LIMIT 1)"

	var exists bool
	err := r.DB.Read.Get(&exists, query, username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check user existence")
		return false, err
	}

	return exists, nil
}

func (r *AuthRepositoryMySQL) Register(user *User) error {
	query :=
		`
	INSERT INTO ums_users (id, username, password, role, created_at, created_by, updated_at, updated_by) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	user.ID = uuid.New().String()
	// Hash the user's password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to encrypt password")
		return err
	}
	user.Password = string(hashedPassword)
	user.Role = strings.ToLower(user.Role)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err = r.DB.Write.Exec(
		query,
		user.ID,
		user.Username,
		user.Password,
		user.Role,
		user.CreatedAt,
		user.CreatedBy,
		user.UpdatedAt,
		user.UpdatedBy,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert user into db")
		return err
	}

	return nil
}
