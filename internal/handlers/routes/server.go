package routes

import (
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/idempotency"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/saveblush/reraw-api/internal/core/breaker"
	"github.com/saveblush/reraw-api/internal/core/config"
	"github.com/saveblush/reraw-api/internal/core/connection/cache"
	"github.com/saveblush/reraw-api/internal/core/connection/sql"
	"github.com/saveblush/reraw-api/internal/core/utils/logger"
	"github.com/saveblush/reraw-api/internal/handlers/middlewares"
	"github.com/saveblush/reraw-api/internal/models"
)

const (
	// MaximumSize10MB body limit 1 mb.
	MaximumSize10MB = 10 * 1024 * 1024
	// MaximumSize1MB body limit 1 mb.
	MaximumSize1MB = 1 * 1024 * 1024
	// Timeout timeout 120 seconds
	Timeout120s = 120 * time.Second
	// Timeout timeout 10 seconds
	Timeout10s = 10 * time.Second
)

type Server struct {
	*fiber.App
}

// NewServer new server
func NewServer() (*Server, error) {
	// New source
	err := newSource()
	if err != nil {
		return nil, err
	}

	// New fiber app
	app := fiber.New(fiber.Config{
		AppName:           config.CF.App.ProjectName,
		ServerHeader:      config.CF.App.ProjectName,
		BodyLimit:         MaximumSize10MB,
		IdleTimeout:       Timeout120s,
		ReadTimeout:       Timeout10s,
		WriteTimeout:      Timeout10s,
		ReduceMemoryUsage: true,
		CaseSensitive:     true,
		JSONEncoder:       sonic.Marshal,
		JSONDecoder:       sonic.Unmarshal,
	})

	// Middlewares
	app.Use(
		compress.New(compress.Config{
			Level: compress.LevelBestCompression,
		}),
		cors.New(),
		requestid.New(),
		idempotency.New(),
		pprof.New(),
		recover.New(),
	)

	// Limiter
	if config.CF.HTTPServer.RateLimit.Enable {
		app.Use(limiter.New(limiter.Config{
			Max:        config.CF.HTTPServer.RateLimit.Max,
			Expiration: config.CF.HTTPServer.RateLimit.Expiration,
		}))
	}

	// Middlewares custom
	app.Use(
		middlewares.Logger(),
		middlewares.WrapError(),
	)

	// New router
	newRouter(app)

	return &Server{app}, nil
}

// Close close server
func (s *Server) Close() error {
	// shutdown server
	err := s.Shutdown()
	if err != nil {
		return err
	}

	// Cleanup tasks
	logger.Log.Info("Running cleanup tasks...")

	// Close db
	if config.CF.Database.RelaySQL.Enable {
		go sql.CloseConnection(sql.Database)
	}
	logger.Log.Info("Database connection closed")

	return nil
}

// newSource new source
func newSource() error {
	// Init Circuit Breaker
	breaker.Init()

	// New connection database
	err := newDatabase()
	if err != nil {
		return err
	}

	// New cache
	err = newCache()
	if err != nil {
		return err
	}

	return nil
}

// newDatabase new connection database
func newDatabase() error {
	if config.CF.Database.RelaySQL.Enable {
		configuration := &sql.Configuration{
			Host:         config.CF.Database.RelaySQL.Host,
			Port:         config.CF.Database.RelaySQL.Port,
			Username:     config.CF.Database.RelaySQL.Username,
			Password:     config.CF.Database.RelaySQL.Password,
			DatabaseName: config.CF.Database.RelaySQL.DatabaseName,
			DriverName:   config.CF.Database.RelaySQL.DriverName,
			Charset:      config.CF.Database.RelaySQL.Charset,
			MaxIdleConns: config.CF.Database.RelaySQL.MaxIdleConns,
			MaxOpenConns: config.CF.Database.RelaySQL.MaxOpenConns,
			MaxLifetime:  config.CF.Database.RelaySQL.MaxLifetime,
		}
		session, err := sql.InitConnection(configuration)
		if err != nil {
			return err
		}
		sql.Database = session.Database

		if !fiber.IsChild() {
			session.Database.AutoMigrate(&models.User{})
		}
	}

	// Debug db
	if !config.CF.App.Environment.Production() {
		if config.CF.Database.RelaySQL.Enable {
			sql.DebugDatabase()
		}
	}

	return nil
}

// newCache new cache
func newCache() error {
	if config.CF.Cache.Redis.Enable {
		configuration := &cache.Configuration{
			Host:     config.CF.Cache.Redis.Host,
			Port:     config.CF.Cache.Redis.Port,
			Password: config.CF.Cache.Redis.Password,
			DB:       config.CF.Cache.Redis.DB,
		}
		err := cache.Init(configuration)
		if err != nil {
			return err
		}
	}

	return nil
}
