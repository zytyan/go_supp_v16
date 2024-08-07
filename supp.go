package main

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"gorm.io/gorm"
	"html"
	"log"
	"main/crawler"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sync"
)

var tgClient *gotgbot.Bot

type Msg struct {
	ChatId int64
	Id     int64
}

type Supp struct {
	gorm.Model
	ArticleUrlPath string `gorm:"primaryKey"`
	ChannelMsg     Msg    `gorm:"embedded,embeddedPrefix:channel_"`
	LinkedGroupMsg Msg    `gorm:"embedded,embeddedPrefix:linked_group_"`
	Magnets        []Magnet
	Status         string
}

var runningSuppByMsg = make(map[Msg]*Supp)
var mu sync.Mutex

func fileSchema(filename string) string {
	if !filepath.IsAbs(filename) {
		var err error
		filename, err = filepath.Abs(filename)
		if err != nil {
			log.Println(err)
			return ""
		}
	}
	return "file://" + url.QueryEscape(filename)
}

func prepareSuppMsg(article *crawler.Article) string {
	return fmt.Sprintf(`<a href="%s">%s</a>\n`+
		"由 %s 发表于 %s\n"+
		"分类：%s\n"+
		"标签：%s\n"+
		"%s",
		html.EscapeString(article.Url),
		html.EscapeString(article.Title),
		html.EscapeString(article.Author),
		html.EscapeString(article.PostTime),
		html.EscapeString(article.Category),
		html.EscapeString(article.HashTags()),
		html.EscapeString(article.IdTag()))
}

func SendSuppMsg(article *crawler.Article) error {
	urlPath := article.UrlPath()
	if urlPath == "" {
		return fmt.Errorf("article url path is empty, url is %s", article.Url)
	}
	supp := &Supp{ArticleUrlPath: urlPath}
	if db.Take(&supp).Error == nil {
		if supp.Status == "done" {
			return nil
		}
		log.Printf("supp %s already exists, status is %s", urlPath, supp.Status)
		/* TODO: 应该有其他的功能，用于从异常退出时恢复 */
	}
	text := prepareSuppMsg(article)
	ext := path.Ext(article.Url)
	f, err := os.CreateTemp("", "supp*"+ext)
	if err != nil {
		return err
	}
	defer f.Close()
	buf, err := article.DownloadImg()
	if err != nil {
		return err
	}
	_, err = f.Write(buf)
	if err != nil {
		return err
	}
	_ = f.Close()
	photo := gotgbot.InputFileByURL(fileSchema(f.Name()))
	msg, err := tgClient.SendPhoto(config.ChannelId, photo, &gotgbot.SendPhotoOpts{
		Caption:   text,
		ParseMode: "HTML",
	})
	if err != nil {
		return err
	}
	mu.Lock()
	defer mu.Unlock()
	key := Msg{ChatId: msg.Chat.Id, Id: msg.MessageId}
	supp, ok := runningSuppByMsg[key]
	if !ok {
		supp = &Supp{
			ArticleUrlPath: urlPath,
			ChannelMsg:     Msg{ChatId: msg.Chat.Id, Id: msg.MessageId},
		}
		runningSuppByMsg[key] = supp
	} else {
		supp.ArticleUrlPath = urlPath
		supp.ChannelMsg = Msg{ChatId: msg.Chat.Id, Id: msg.MessageId}
		db.Save(supp)
	}
	return nil
}

func IsAutoForwardedSuppMsg(msg *gotgbot.Message) bool {
	if msg.Chat.Id != config.GroupId {
		return false
	}
	if !msg.IsAutomaticForward || msg.ForwardOrigin == nil {
		return false
	}
	ori := msg.ForwardOrigin.MergeMessageOrigin()
	if ori.Chat == nil || ori.Chat.Id != config.ChannelId {
		return false
	}
	return true
}

func OnLinkedGroupMsg(bot *gotgbot.Bot, msg *gotgbot.Message) {
	mu.Lock()
	defer mu.Unlock()
	cid := msg.ForwardOrigin.MergeMessageOrigin().Chat.Id
	mid := msg.ForwardOrigin.MergeMessageOrigin().MessageId
	key := Msg{ChatId: cid, Id: mid}
	supp, ok := runningSuppByMsg[key]
	if !ok {
		supp = &Supp{
			LinkedGroupMsg: Msg{ChatId: msg.Chat.Id, Id: msg.MessageId},
		}
		runningSuppByMsg[key] = supp
		return
	}
	supp.LinkedGroupMsg = Msg{ChatId: msg.Chat.Id, Id: msg.MessageId}
	db.Save(supp)
}
