package server

import (
	"context"
	"log"
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
	DB             *db.DB
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

	sessionStorage := session.NewSessionMangerDbOps(db)

	return &Server{
		ctx:            ctx,
		wg:             wg,
		mux:            &sync.RWMutex{},
		SessionManager: session.NewSessionManager(ctx, sessionStorage),
		UserManager:    user.NewUserManager(),
		app:            a,
		DB:             db,
	}, nil
}

func (s *Server) Shutdown() error {
	log.Println("Shutting down server")
	defer s.DB.Close()
	err := s.DB.DeleteAllSessions()
	if err != nil {
		log.Printf("Error deleting sessions: %v", err)
		return err
	}

	log.Println("Sessions deleted")
	return nil
}
