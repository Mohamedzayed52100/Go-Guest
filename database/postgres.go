package database

import (
	"context"
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/config"
	"github.com/goplaceapp/goplace-user/pkg/tenantservice/domain"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresService struct {
	Db                  *gorm.DB
	TenantDbConnections map[string]*gorm.DB
	SvcCfg              *config.Config
}

var SharedPostgresService *PostgresService

func NewService() *PostgresService {
	if SharedPostgresService == nil {
		// parse environment variables
		cfg := &config.Config{}
		if err := env.Parse(cfg); err != nil {
			panic(fmt.Errorf("failed to parse environment variables, %w", err))
		}

		// establish database connection
		db, err := NewDatabaseInstance(cfg.DbPostgresHost, cfg.DbPostgresUser, cfg.DbPostgresPassword, cfg.DbPostgresDb, cfg.DbSSLMode)
		if err != nil {
			panic(fmt.Errorf("failed to connect to the database instance, %w", err))
		}

		// run shared database migrations
		if err := SharedMigrationsUp(db); err != nil {
			panic(fmt.Errorf("failed to run shared database migrations, %w", err))
		}

		// run tenant database migrations
		tenants, err := GetAllTenants(db)

		if err != nil {
			panic(fmt.Errorf("failed to get tenants, %w", err))
		}

		// create tenant database connections
		tenantDbConnections := make(map[string]*gorm.DB)
		for _, tenant := range tenants {
			var tenantDb *gorm.DB

			tenantDb, err = NewTenantDatabaseInstance(cfg.DbPostgresHost, cfg.DbPostgresUser, cfg.DbPostgresPassword, tenant.DbName, cfg.DbSSLMode)
			if err != nil {
				panic(fmt.Errorf("failed to connect to the tenant database instance, %w", err))
			}

			if err := TenantMigrationsUp(tenantDb); err != nil {
				panic(fmt.Errorf("failed to run tenant database migrations, %w", err))
			}

			tenantDbConnections[tenant.DbName] = tenantDb
		}

		SharedPostgresService = &PostgresService{
			Db:                  db,
			TenantDbConnections: tenantDbConnections,
			SvcCfg:              cfg,
		}

		return SharedPostgresService
	}

	return SharedPostgresService
}

func NewDatabaseInstance(host, user, password, dbName, sslmode string) (*gorm.DB, error) {
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		host, user, password, dbName, sslmode)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger:         gormLogger,
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewTenantDatabaseInstance(host, user, password, dbName, sslmode string) (*gorm.DB, error) {
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres sslmode=%s", host, user, password, sslmode)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger:         gormLogger,
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}

	var count int64
	result := db.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", dbName).Scan(&count)
	if result.Error != nil {
		return nil, result.Error
	}

	connectionString = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", host, user, password, dbName, sslmode)

	if count > 0 {
		db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
			Logger:         gormLogger,
			TranslateError: true,
		})
	} else {
		result := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if result.Error != nil {
			return nil, result.Error
		}

		db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
			Logger:         gormLogger,
			TranslateError: true,
		})
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (s *PostgresService) GetSharedDB() *gorm.DB {
	return s.Db
}

func (s *PostgresService) GetTenantDB(ctx context.Context) *gorm.DB {
	tenantDBNameValue := ctx.Value(meta.TenantDBNameContextKey.String())
	if tenantDBNameValue == nil {
		logger.Default().Error("tenant DB name not found in context")
	}
	tenantDBName, ok := tenantDBNameValue.(string)
	if !ok {
		logger.Default().Error("tenant DB name not found in context")
	}

	// Check if the connection already exists in the pool
	if db, ok := s.TenantDbConnections[tenantDBName]; ok {
		return db
	}

	// If not, create a new connection and add it to the pool
	db, err := NewTenantDatabaseInstance(s.SvcCfg.DbPostgresHost, s.SvcCfg.DbPostgresUser, s.SvcCfg.DbPostgresPassword, tenantDBName, s.SvcCfg.DbSSLMode)
	if err != nil {
		logger.Default().Error("failed to connect to the database instance (GetTenantDbConnection), %w", zap.Error(err))
	}

	return db
}

func GetAllTenants(db *gorm.DB) ([]*domain.Tenant, error) {
	var tenants []*domain.Tenant
	if err := db.Find(&tenants).Error; err != nil {
		return nil, err
	}

	return tenants, nil
}
