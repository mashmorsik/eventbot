package app

type App struct {
	GRPCServer *app.App
}

func New(grpcPort int) *App {
	// init eventserver
	grpcApp := app.New(grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
