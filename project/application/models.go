package application

import "api-your-accounts/project/domain"

type CreateData struct {
	Name   string
	Type   domain.ProjectType
	UserId uint
}

type FindByUserRecord struct {
	ID   uint
	Name string
	Type domain.ProjectType
	Data map[string]any
}
