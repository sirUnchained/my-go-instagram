package main

import (
	"time"

	"github.com/sirUnchained/my-go-instagram/internal/auth"
	"github.com/sirUnchained/my-go-instagram/internal/configs"
	"github.com/sirUnchained/my-go-instagram/internal/database"
	"github.com/sirUnchained/my-go-instagram/internal/storage"
	"github.com/sirUnchained/my-go-instagram/internal/storage/cache"
	"go.uber.org/zap"
)

func main() {
	// init logger
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer log.Sync()
	sugar := log.Sugar()

	// init configs from json file
	cfg, err := configs.GetConfigs()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// init postgres configs
	postgres := pg_db{
		Addr:         cfg.Postgres.Addr,
		MaxOpenConns: cfg.Postgres.MaxOpenConns,
		MaxIdleConns: cfg.Postgres.MaxIdleConns,
		MaxIdleTime:  cfg.Postgres.MaxIdleTime,
	}

	// init main database (postgres)
	db, err := database.NewPostgreSQL(postgres.Addr, postgres.MaxOpenConns, postgres.MaxIdleConns, postgres.MaxIdleTime)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	pgStorage := storage.NewPgStorage(db)

	// init redis configs
	redis := redis_db{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DBNumber: cfg.Redis.DBNumber,
		Enabled:  cfg.Redis.Enabled,
	}

	// init cache database (redis)
	redisClient := database.NewRedisClient(redis.Addr, redis.Password, redis.DBNumber)
	redisStorage := cache.NewRedisStorage(redisClient)

	// init token configs
	tokenConfigs := authConfig{
		secretKey: cfg.Auth.SecretKey,
		aud:       cfg.Auth.Aud,
		iss:       cfg.Auth.Iss,
		expMin:    time.Duration(cfg.Auth.ExpMin) * time.Minute,
	}

	// init authenticator
	authenticator := auth.NewJWTAuthenticator(cfg.Auth.SecretKey, cfg.Auth.Aud, cfg.Auth.Iss)

	// set server configs
	srvCfg := serverConfigs{
		addr:          cfg.Addr,
		isDevelopment: cfg.IsDevelopment,
		database:      postgres,
		cache:         redis,
		auth:          tokenConfigs,
	}

	// create server struct
	srv := server{
		serverConfigs:  srvCfg,
		postgreStorage: pgStorage,
		redisStorage:   redisStorage,
		auth:           authenticator,
		logger:         sugar,
	}

	// start server
	mux := srv.getRouter()
	srv.start(mux)

}
