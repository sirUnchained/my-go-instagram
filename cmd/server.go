package main

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type server struct {
	serverConfigs serverConfigs
	logger        *zap.SugaredLogger
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

func (s *server) start() {
	server := &http.Server{
		Addr:         s.serverConfigs.addr,
		Handler:      nil,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Minute * 2,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
