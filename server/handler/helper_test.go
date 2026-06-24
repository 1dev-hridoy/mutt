package handler

import (
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/dishan1223/mutt/internal/config"
	"github.com/dishan1223/mutt/internal/service"
	"github.com/dishan1223/mutt/models"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var mr *miniredis.Miniredis

func TestMain(m *testing.M) {
	os.Setenv("PORT", "3000")

	var err error
	mr, err = miniredis.Run()
	if err != nil {
		panic("failed to start miniredis: " + err.Error())
	}

	config.RDB = redis.NewClient(&redis.Options{Addr: mr.Addr()})

	code := m.Run()
	mr.Close()
	os.Exit(code)
}

func setupTestDB(t *testing.T) {
	t.Helper()
	var err error
	config.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	config.DB.AutoMigrate(&models.User{}, &models.Project{}, &models.ErrorGroup{}, &models.Error{})
	service.MustInitJWT("test-secret-key-for-testing")
}

func teardownTestDB(t *testing.T) {
	t.Helper()
	sqlDB, _ := config.DB.DB()
	sqlDB.Close()
	mr.FlushAll()
}

func seedUser(t *testing.T, username, email, password, phone string) models.User {
	t.Helper()
	hashed, _ := service.HashPassword(password)
	user := models.User{
		Username: username,
		Email:    email,
		Password: hashed,
		Phone:    phone,
	}
	config.DB.Create(&user)
	return user
}
