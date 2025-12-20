package main

import (
	"net/http"
	"time"

	"github.com/sirUnchained/my-go-instagram/internal/storage"
	"github.com/sirUnchained/my-go-instagram/internal/storage/cache"
	"go.uber.org/zap"
)

type server struct {
	serverConfigs  serverConfigs
	postgreStorage *storage.PgStorage
	redisStorage   *cache.RedisStorage
	logger         *zap.SugaredLogger
}

type serverConfigs struct {
	addr     string
	database pg_db
	cache    redis_db
}

type pg_db struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type redis_db struct {
	Addr     string
	Password string
	DBNumber int
	Enabled  bool
}

func (s *server) start(mux *http.ServeMux) {
	server := &http.Server{
		Addr:         s.serverConfigs.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Minute * 2,
	}

	s.logger.Infoln("starting server at", s.serverConfigs.addr)
	err := server.ListenAndServe()
	if err != nil {
		s.logger.Fatalln(err.Error())
	}
}
