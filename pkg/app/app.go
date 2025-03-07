package app

import (
	"context"
	"errors"
	"fmt"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	auth "firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type Authenticator interface {
	VerifyIDToken(token string) (string, error)
}

type App struct {
	firebase *firebase.App
	auth     *auth.Client
	store    *firestore.Client
	ctx      context.Context
}

func NewApp(ctx context.Context) (*App, error) {
	// the file is in the root directory
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, opt)
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

func (app *App) Authenticate(token string) (string, error) {
	user, err := app.auth.VerifyIDToken(app.ctx, token)
	if err != nil {
		return "", errors.New("failed to verify id token")
	}

	return user.UID, nil
}

func (app *App) VerifyIDToken(token string) (string, error) {
	user, err := app.auth.VerifyIDToken(app.ctx, token)
	if err != nil {
		return "", errors.New("failed to verify id token")
	}

	return user.UID, nil
}
