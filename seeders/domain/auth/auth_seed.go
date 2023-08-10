package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/infras"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type Profile struct {
	ProfileID string    `db:"profile_id"`
	CreatedAt time.Time `db:"created_at"`
	CreatedBy string    `db:"created_by"`
	UpdatedAt time.Time `db:"updated_at"`
	UpdatedBy string    `db:"updated_by"`
}

type Status struct {
	StatusID  string    `db:"status_id"`
	CreatedAt time.Time `db:"created_at"`
	CreatedBy string    `db:"created_by"`
	UpdatedAt time.Time `db:"updated_at"`
	UpdatedBy string    `db:"updated_by"`
}

type User struct {
	ID        string    `db:"id"`
	ProfileID string    `db:"profile_id"`
	StatusID  string    `db:"status_id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	CreatedBy string    `db:"created_by"`
	UpdatedAt time.Time `db:"updated_at"`
	UpdatedBy string    `db:"updated_by"`
}

func main() {
	config := configs.Get()

	mysqlConn := infras.ProvideMySQLConn(config)

	id := uuid.New().String()
	profileID := uuid.New().String()
	statusID := uuid.New().String()
	username := "admin_fauzy"
	password := "passwordkuat"
	// Generate the bcrypt hash for the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to encrypt password")
		return
	}
	role := "admin"
	CreatedAt := time.Now()
	updatedAt := time.Now()

	profile := Profile{
		ProfileID: profileID,
		CreatedAt: CreatedAt,
		CreatedBy: username,
		UpdatedAt: updatedAt,
		UpdatedBy: username,
	}

	status := Status{
		StatusID:  statusID,
		CreatedAt: CreatedAt,
		CreatedBy: username,
		UpdatedAt: updatedAt,
		UpdatedBy: username,
	}

	user := User{
		ID:        id,
		ProfileID: profileID,
		StatusID:  statusID,
		Username:  username,
		Password:  string(hashedPassword),
		Role:      role,
		CreatedAt: CreatedAt,
		CreatedBy: username,
		UpdatedAt: updatedAt,
		UpdatedBy: username,
	}

	// Insert profile_id into ums_profiles
	profileQuery := "INSERT INTO ums_profiles (id, created_at, created_by, updated_at, updated_by) VALUES (?,?,?,?,?)"
	_, err = mysqlConn.Write.Exec(
		profileQuery,
		profile.ProfileID,
		profile.CreatedAt,
		profile.CreatedBy,
		profile.UpdatedAt,
		profile.UpdatedBy,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert profile data")
		return
	}

	// Insert status_id into ums_status
	statusQuery := "INSERT INTO ums_status (id, created_at, created_by, updated_at, updated_by) VALUES (?,?,?,?,?)"
	_, err = mysqlConn.Write.Exec(
		statusQuery,
		status.StatusID,
		status.CreatedAt,
		status.CreatedBy,
		status.UpdatedAt,
		status.UpdatedBy,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert status data")
		return
	}

	// Insert the user data into ums_users
	userQuery := "INSERT INTO ums_users (id, profile_id, status_id, username, password, role, created_at, created_by, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = mysqlConn.Write.Exec(
		userQuery,
		user.ID,
		user.ProfileID,
		user.StatusID,
		user.Username,
		user.Password,
		user.Role,
		user.CreatedAt,
		user.CreatedBy,
		user.UpdatedAt,
		user.UpdatedBy,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert user data")
		return
	}
}
