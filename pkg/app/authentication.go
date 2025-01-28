package app

func (app *App) AuthenticateUser(token string) (bool, error) {
	_, err := app.auth.VerifyIDToken(app.ctx, token)
	if err != nil {
		return false, err
	}
	return true, nil
}
