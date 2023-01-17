package main

import (
	"core-server/src/controller/storage"
	"core-server/src/controller/web"
	"core-server/src/pkg/config"
	"core-server/src/pkg/jwt"
	"core-server/src/pkg/logger"
	"core-server/src/pkg/turnstile"
)

func main() {
	// init logger
	l := logger.NewLogger()

	// init config
	cfg, err := config.NewConfig()
	if err != nil {
		l.Fatal(err)
	}

	// init database
	st := storage.NewStorage()
	db, err := st.Connect(&cfg.Database)
	if err != nil {
		l.Fatal(err)
	}

	// init cf turnstile
	c := turnstile.NewService(cfg.CF.Secret)

	// init jwt service
	j := jwt.NewService(&jwt.Service{
		AccessSecret:  cfg.JWT.AccessSecret,
		RefreshSecret: cfg.JWT.RefreshSecret,
	})

	// init and start web server
	s := web.NewServer(l, j, db, c)
	s.Start(cfg.HTTP.Port)
}
