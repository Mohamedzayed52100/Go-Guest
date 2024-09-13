package database

import (
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"log"
	"os"
	"time"

	"github.com/goplaceapp/goplace-guest/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	GetSharedDB func() *gorm.DB
	GetTenantDB func(tenant string) *gorm.DB
}

var ignoreRecordNotFoundError = true

func init() {
	if os.Getenv("ENVIRONMENT") != meta.ProdEnvironment {
		ignoreRecordNotFoundError = false
	}
}

var gormLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags),
	logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  utils.GetLogLevel(),
		IgnoreRecordNotFoundError: ignoreRecordNotFoundError,
		ParameterizedQueries:      true,
		Colorful:                  true,
	},
)
