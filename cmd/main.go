package main

import (
	"github.com/sirUnchained/my-go-instagram/internal/configs"
	"go.uber.org/zap"
)

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer log.Sync()
	sugar := log.Sugar()

	cfg, err := configs.GetConfigs()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	postgres := pg_db{
		Addr:         cfg.Postgres.Addr,
		MaxOpenConns: cfg.Postgres.MaxOpenConns,
		MaxIdleConns: cfg.Postgres.MaxIdleConns,
		MaxIdleTime:  cfg.Postgres.MaxIdleTime,
	}
	// db, err := database.NewPostgreSQL(postgres.Addr, postgres.MaxOpenConns, postgres.MaxIdleConns, postgres.MaxIdleTime)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// 	return
	// }

	redis := redis_db{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DBNumber: cfg.Redis.DBNumber,
		Enabled:  cfg.Redis.Enabled,
	}
	// redisClient := database.NewRedisClient(redis.Addr, redis.Password, redis.DBNumber)

	srvCfg := serverConfigs{
		addr:     cfg.Addr,
		database: postgres,
		cache:    redis,
	}

	srv := server{
		serverConfigs: srvCfg,
		logger:        sugar,
	}

	srv.start()

}
