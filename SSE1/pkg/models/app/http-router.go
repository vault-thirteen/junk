package app

// A Configurator of an HTTP Server's Router.
func (app *Application) initHttpRouter() (err error) {

	// Actions with a single User.
	app.httpRouter.POST(
		"/api/user/register",
		app.httpProtocolCheck(
			app.httpAuthentication(
				app.httpHandlerApiUserRegister, false,
			),
		),
	)

	app.httpRouter.POST(
		"/api/user/disable",
		app.httpProtocolCheck(
			app.httpAuthentication(
				app.httpHandlerApiUserDisable, true,
			),
		),
	)

	app.httpRouter.POST(
		"/api/user/log-in",
		app.httpProtocolCheck(
			app.httpAuthentication(
				app.httpHandlerApiUserLogIn, false,
			),
		),
	)

	app.httpRouter.POST(
		"/api/user/log-out",
		app.httpProtocolCheck(
			app.httpAuthentication(
				app.httpHandlerApiUserLogOut, true,
			),
		),
	)

	// Actions with multiple Users.
	app.httpRouter.POST(
		"/api/users/list",
		app.httpProtocolCheck(
			app.httpAuthentication(
				app.httpHandlerApiUsersList, true,
			),
		),
	)

	return
}
