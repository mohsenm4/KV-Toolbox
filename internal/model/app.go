package model

type App struct {
	Name       string
	MainWindow MainWindow
}

func NewApp(name string) *App {
	return &App{
		Name: name,
	}
}
