package app

import (
	"context"
	"fmt"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	auth "firebase.google.com/go/auth"
)

type App struct {
	firebase *firebase.App
	auth     *auth.Client
	store    *firestore.Client
	ctx      context.Context
}

func NewApp(ctx context.Context) (*App, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	auth, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %v", err)
	}

	store, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Firestore client: %v", err)
	}

	return &App{
		firebase: app,
		auth:     auth,
		store:    store,
		ctx:      ctx,
	}, nil
}
