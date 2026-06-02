package controller

import ("gorm.io/gorm")

type Controller struct {
	db *gorm.DB
}

func Create(db *gorm.DB) *Controller {
	return &Controller{ db }
}