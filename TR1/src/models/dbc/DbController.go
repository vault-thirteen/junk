package dbc

import "gorm.io/gorm"

type DbController struct {
	db       *gorm.DB
	pageSize int
}

func NewDbController(db *gorm.DB) *DbController {
	return &DbController{
		db: db,
	}
}

func NewDbControllerWithPageSize(db *gorm.DB, pageSize int) *DbController {
	return &DbController{
		db:       db,
		pageSize: pageSize,
	}
}
