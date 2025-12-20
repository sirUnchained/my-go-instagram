package main

import (
	"net/http"

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

	// init main database (postgres)
	postgres := pg_db{
		Addr:         cfg.Postgres.Addr,
		MaxOpenConns: cfg.Postgres.MaxOpenConns,
		MaxIdleConns: cfg.Postgres.MaxIdleConns,
		MaxIdleTime:  cfg.Postgres.MaxIdleTime,
	}
	db, err := database.NewPostgreSQL(postgres.Addr, postgres.MaxOpenConns, postgres.MaxIdleConns, postgres.MaxIdleTime)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	pgStorage := storage.NewPgStorage(db)

	// init cache database (redis)
	redis := redis_db{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DBNumber: cfg.Redis.DBNumber,
		Enabled:  cfg.Redis.Enabled,
	}
	redisClient := database.NewRedisClient(redis.Addr, redis.Password, redis.DBNumber)
	redisStorage := cache.NewRedisStorage(redisClient)

	// set server configs
	srvCfg := serverConfigs{
		addr:     cfg.Addr,
		database: postgres,
		cache:    redis,
	}

	// create server struct
	srv := server{
		serverConfigs:  srvCfg,
		postgreStorage: pgStorage,
		redisStorage:   redisStorage,
		logger:         sugar,
	}

	// start server
	mux := http.NewServeMux()
	srv.start(mux)

}
