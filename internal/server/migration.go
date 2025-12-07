package server

import (
	"context"
	"go-nunu/internal/model"
	"go-nunu/pkg/log"
	"os"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MigrateServer struct {
	db  *gorm.DB
	log *log.Logger
}

func NewMigrateServer(db *gorm.DB, log *log.Logger) *MigrateServer {
	return &MigrateServer{
		db:  db,
		log: log,
	}
}
func (m *MigrateServer) Start(ctx context.Context) error {
	m.db.Migrator().DropTable(
		&model.User{},
		&model.Role{},
		&model.Permission{},
	)
	if err := m.db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
	); err != nil {
		m.log.Error("err: ", zap.Error(err))
		return err
	}
	if err := m.initAdminUser(); err != nil {
		m.log.Error("initAdminUser error", zap.Error(err))
		return err
	}
	m.log.Info("AutoMigrate success")
	os.Exit(0)
	return nil
}
func (m *MigrateServer) Stop(ctx context.Context) error {
	m.log.Info("AutoMigrate stop")
	return nil
}
func (m *MigrateServer) initAdminUser() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		m.log.Error("bcrypt.GenerateFromPassword error", zap.Error(err))
		return err
	}
	return m.db.Create(&model.User{
		BaseModel: model.BaseModel{ID: 1},
		UserId:    model.AdminUserID,
		Password:  string(hashedPassword),
		Email:     "admin@gmail.com",
		Name:      "Administrator",
	}).Error
}
