package databases

import (
	"github.com/go-sql-driver/mysql"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"treafik-api/pkg/logger"
)

type MysqlDBConfig struct {
	Dsn          string `yaml:"dns" json:"dns" mapstructure:"dns" toml:"dns"`
	MaxOpenCount int    `yaml:"max_open_count" json:"max_open_count" mapstructure:"max_open_count"  validate:"required" toml:"max_open_count"`
	MaxIdleCount int    `yaml:"max_idle_count" json:"max_idle_count" mapstructure:"max_idle_count"     validate:"required" toml:"max_idle_count"`
	Tracing      bool   `yaml:"tracing" json:"tracing" mapstructure:"tracing" toml:"tracing"`
}

func InitGorm(cfg *MysqlDBConfig) (*gorm.DB, error) {
	opt, err := mysql.ParseDSN(cfg.Dsn)
	if err != nil {
		return nil, err
	}
	if !opt.ParseTime {
		logger.Warnw("InitGormV2: parseTime is disabled")
	}
	if opt.Loc.String() != "UTC" {
		logger.Infow("using non UTC timezone for parseTime", "timezone_used", opt.Loc.String())
	} else {
		logger.Infow("using UTC timezone for parseTime")
	}

	db, err := gorm.Open(gormmysql.Open(cfg.Dsn), &gorm.Config{
		// Logger: gormlogger.NewDefault(logger.GetLogger()),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	if cfg.Tracing {
		plugin := otelgorm.NewPlugin(
			otelgorm.WithDBName(opt.DBName),
			otelgorm.WithAttributes(semconv.DBSystemMySQL, attribute.String("db.addr", opt.Addr)),
		)
		db.Use(plugin)
	}

	// otelConfig conn pool https://gorm.io/docs/connecting_to_the_database.html#Connection-Pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	if cfg.MaxIdleCount > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleCount)
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	if cfg.MaxOpenCount > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenCount)
	}

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	// sqlDB.SetConnMaxLifetime(time.Hour)

	// the maximum amount of time a connection may be idle
	// SetConnMaxIdleTime; added in Go 1.15
	// sqlDB.SetConnMaxIdleTime(time.Second * 3600)

	return db, nil
}
