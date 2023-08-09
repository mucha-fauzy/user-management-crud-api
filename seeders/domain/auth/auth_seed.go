package main

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/infras"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `db:"id"`
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

	uuid := "ce560f71-1495-4d6c-b6bf-f9f1b3b2f212"
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
	user := User{
		ID:        uuid,
		Username:  username,
		Password:  string(hashedPassword),
		Role:      role,
		CreatedAt: CreatedAt,
		CreatedBy: username,
		UpdatedAt: updatedAt,
		UpdatedBy: username,
	}

	// Insert the user data into the database
	query := "INSERT INTO ums_users (id, username, password, role, created_at, created_by, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = mysqlConn.Write.Exec(
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
		log.Error().Err(err).Msg("Failed to insert data")
		return
	}
}
