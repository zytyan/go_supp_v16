package main

import _ "github.com/mattn/go-sqlite3"

import (
	"github.com/BurntSushi/toml"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SuppConfig struct {
	BotToken string `toml:"bot_token"`
	BaseUrl  string `toml:"base_url"`

	ChannelId      int64 `toml:"channel_id"`
	GroupId        int64 `toml:"group_id"`
	VideoChannelId int64 `toml:"video_channel_id"`
	AdminId        int64 `toml:"admin_id"`
}

var config SuppConfig

func LoadCfg() {
	// read config.toml
	// if not exists, or error, print error and exit
	cfgFile := "config/config.toml"
	_, err := toml.DecodeFile(cfgFile, &config)
	if err != nil {
		panic(err)
	}
}

var db *gorm.DB

func InitDB() {
	dbFile := "data.db"
	var err error
	db, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}
