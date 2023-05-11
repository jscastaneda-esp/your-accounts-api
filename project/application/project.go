package application

//go:generate mockery --name IProjectApp --filename project-app.go
type IProjectApp interface {
}

type projectApp struct {
}

func NewProjectApp() IProjectApp {
	return &projectApp{}
}
