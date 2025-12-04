package casbinPkg

import (
	"log"
	"os"
	"path/filepath"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func NewEnforcer(db *gorm.DB) *casbin.CachedEnforcer {
	a, _ := gormadapter.NewAdapterByDB(db)
	cwd, _ := os.Getwd()
	m, err := model.NewModelFromFile(filepath.Join(cwd, "config", "model.conf"))
	if err != nil {
		log.Fatalf("failed to create Casbin enforcer. Check paths: Model=%v. Error: %v", m, err)
	}
	e, _ := casbin.NewCachedEnforcer(m, a)
	e.SetExpireTime(60 * 60)
	e.LoadPolicy()
	return e
}
