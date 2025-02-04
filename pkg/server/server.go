package server

import (
	"context"
	"sync"

	"github.com/Monkhai/shwipe-server.git/pkg/app"
	"github.com/Monkhai/shwipe-server.git/pkg/db"
	"github.com/Monkhai/shwipe-server.git/pkg/session"
	"github.com/Monkhai/shwipe-server.git/pkg/user"
)

type Server struct {
	ctx            context.Context
	mux            *sync.RWMutex
	wg             *sync.WaitGroup
	SessionManager *session.SessionManager
	UserManager    *user.UserManager
	app            *app.App
	db             *db.DB
}

func NewServer(ctx context.Context, wg *sync.WaitGroup) (*Server, error) {
	a, err := app.NewApp(ctx)
	if err != nil {
		return &Server{}, nil
	}
	db, err := db.NewDB(ctx)
	if err != nil {
		return &Server{}, err
	}

	return &Server{
		ctx:            ctx,
		wg:             wg,
		mux:            &sync.RWMutex{},
		SessionManager: session.NewSessionManager(ctx),
		UserManager:    user.NewUserManager(),
		app:            a,
		db:             db,
	}, nil
}
