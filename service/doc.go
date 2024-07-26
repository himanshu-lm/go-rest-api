package service

import (
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -source=service.go -destination=./db_mock.go -package=service . Db
