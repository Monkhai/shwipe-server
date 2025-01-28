package app

import auth "firebase.google.com/go/auth"

func (app *App) AuthenticateUser(token string) (string, error) {
	user, err := app.auth.VerifyIDToken(app.ctx, token)
	if err != nil {
		return "", err
	}

	return user.UID, nil
}

func (app *App) GetUserRecord(id string) (*auth.UserRecord, error) {
	usr, err := app.auth.GetUser(app.ctx, id)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
