package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	_ "github.com/sirUnchained/my-go-instagram/docs"
	"github.com/sirUnchained/my-go-instagram/internal/auth"
	"github.com/sirUnchained/my-go-instagram/internal/storage"
	"github.com/sirUnchained/my-go-instagram/internal/storage/cache"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/unrolled/secure"
	"go.uber.org/zap"
)

type server struct {
	serverConfigs  serverConfigs
	postgreStorage *storage.PgStorage
	redisStorage   *cache.RedisStorage
	auth           auth.Authenticator
	logger         *zap.SugaredLogger
}

type serverConfigs struct {
	addr          string
	isDevelopment bool
	database      pg_db
	cache         redis_db
	auth          authConfig
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

type authConfig struct {
	secretKey string
	aud       string
	iss       string
	expMin    time.Duration
}

func (s *server) getRouter() http.Handler {
	r := chi.NewRouter()

	// set middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4000/*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(secure.New(secure.Options{
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		// currently we have no http/ssl host
		SSLRedirect: false,
		SSLHost:     "",
		// security policy
		ReferrerPolicy:    "strict-origin-when-cross-origin",
		PermissionsPolicy: "camera=(), microphone=(), geolocation=(), payment=()",
		// we are currently in development mode
		IsDevelopment: s.serverConfigs.isDevelopment,
	}).Handler)
	r.Use(httprate.Limit(
		100,
		time.Minute*1,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	))

	r.Get("/health-check", s.checkHealthHandler)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:4000/v1/swagger/doc.json")))

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", s.registerUserHandler)
			r.Post("/login", s.loginUserHandler)
		})

		r.Route("/posts", func(r chi.Router) {
			r.Use(s.checkUserTokenMiddleware)
			r.Use(s.checkIsUserVerifiedMiddleware)
			r.Post("/new", s.createPostHandler)
			r.Get("/{postid}", s.getPostHandler)
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(s.checkUserTokenMiddleware)
			r.Get("/me", s.getMeHandler)
			r.Group(func(r chi.Router) {
				r.Use(s.checkIsUserVerifiedMiddleware)
				r.Get("/{userid}", s.checkAccessToPageMiddleware(s.getUserHandler))
			})
		})
	})

	return r
}

func (s *server) start(mux http.Handler) {
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
