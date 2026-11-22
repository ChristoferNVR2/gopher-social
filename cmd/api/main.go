package main

import (
	"time"

	"github.com/ChristoferNVR2/gopher-social/internal/auth"
	"github.com/ChristoferNVR2/gopher-social/internal/db"
	"github.com/ChristoferNVR2/gopher-social/internal/env"
	"github.com/ChristoferNVR2/gopher-social/internal/mailer"
	"github.com/ChristoferNVR2/gopher-social/internal/ratelimiter"
	"github.com/ChristoferNVR2/gopher-social/internal/store"
	"github.com/ChristoferNVR2/gopher-social/internal/store/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const version = "0.0.2"

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social media network for Gophers.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("API_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/gopher-social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PASSWORD", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		env:     env.GetString("ENV", "development"),
		version: version,
		mail: mailConfig{
			resend: resendConfig{
				apiKey: env.GetString("RESEND_API_KEY", ""),
			},
			fromEmail: env.GetString("FROM_EMAIL", ""),
			exp:       time.Hour * 24 * 3,
		},
		auth: authConfig{
			basic: basicAuthConfig{
				user: env.GetString("BASIC_AUTH_USER", "admin"),
				pass: env.GetString("BASIC_AUTH_PASS", "password"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "mysecretkey"),
				exp:    time.Hour * 24 * 3,
				issuer: "gopher-social",
			},
		},
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATE_LIMIT_REQUESTS", 100),
			TimeFrame:            time.Minute,
			Enabled:              env.GetBool("RATE_LIMIT_ENABLED", true),
		},
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		defer rdb.Close()
		logger.Info("redis client established")
	}

	store := store.NewStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)

	resend, err := mailer.NewResend(cfg.mail.resend.apiKey, cfg.mail.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.issuer,
		cfg.auth.token.issuer,
	)

	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        resend,
		authenticator: jwtAuthenticator,
		rateLimiter:   rateLimiter,
	}

	mux := app.mount()

	if err := app.run(mux); err != nil {
		logger.Errorw("server stopped with error", "error", err)
	}
}
