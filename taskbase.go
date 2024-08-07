package main

import (
	"gorm.io/gorm"
	"main/crawler"
)

type Magnet struct {
	gorm.Model
	Hash   string
	Status string
	Tasks  []TaskDb
}

type TaskDb struct {
	gorm.Model
	Name    string
	Status  string
	CtxData []byte
}

type Task struct {
	TaskDb
	Article  *crawler.Article
	Callback func(*Supp, *Magnet, *Task)
}
