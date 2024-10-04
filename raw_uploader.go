package main

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"log"
	"main/archive_proc"
	"main/helper"
	"main/qbit"
	"main/strnum"
	"os"
	"path/filepath"
)

func rarFiles(file string) (p *preparedFiles, err error) {
	stem := helper.Truncate(helper.Stem(file), 55)
	dir := filepath.Join(os.TempDir(), stem)
	if _, err = os.Stat(dir); err == nil {
		os.RemoveAll(dir)
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}
	p = &preparedFiles{isRar: true}
	rarName := stem + ".rar"
	err = archive_proc.PackToRar(file, dir, rarName)
	p.cleanupFn = func() { os.RemoveAll(dir) }
	var entries []os.DirEntry
	entries, err = os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, f := range entries {
		p.files = append(p.files, filepath.Join(dir, f.Name()))
	}
	return
}

type preparedFiles struct {
	files     []string
	isRar     bool
	cleanupFn func()
}

func (p *preparedFiles) cleanup() {
	if p.cleanupFn != nil {
		p.cleanupFn()
	}
}

func prepareUploadFiles(path string) (p *preparedFiles, err error) {
	if helper.NeedRar(path) {
		return rarFiles(path)
	}
	p = &preparedFiles{isRar: false}
	stat, err := os.Stat(path)
	if err != nil {
		return
	}
	if !stat.IsDir() {
		p.files = []string{path}
		return
	}
	files, err := os.ReadDir(path)
	if err != nil {
		return
	}
	for _, file := range files {
		p.files = append(p.files, filepath.Join(path, file.Name()))
	}
	return
}

func UploadRawFiles(t *qbit.Torrent, supp *Supp) error {
	path := t.ContentPath
	files, err := prepareUploadFiles(path)
	defer files.cleanup()
	if err != nil {
		log.Println(err)
		return err
	}
	var newFiles []string
	for _, file := range files.files {
		if helper.IsVideoFile(file) {
			// skip video files
			continue
		}
		newFiles = append(newFiles, file)
	}
	newFiles = strnum.SortedStrings(newFiles)
	log.Printf("prepare to upload %d files\n", len(newFiles))
	for chunkIdx, fileChunk := range helper.Chunk(newFiles, 10) {
		inputMedia := make([]gotgbot.InputMedia, 0, len(fileChunk))
		for _, f := range fileChunk {
			log.Printf("send file to: %d, %s\n", supp.LinkedGroupMsg.ChatId, f)
			inputMedia = append(inputMedia, gotgbot.InputMediaDocument{
				Media:     gotgbot.InputFileByURL(fileSchema(f)),
				Caption:   "",
				ParseMode: "",
			})
		}
		if files.isRar {
			captionText := fmt.Sprintf("这是一个分卷压缩包文件，共%d个分卷。\n"+
				"本条消息中的压缩包为分卷%d-%d\n"+
				"您需要下载总共所有%d个分卷才可以正确解压这批压缩包。\n"+
				"如果您使用的软件并非WinRAR，在解压本消息中的压缩包时遇到需要输入密码等意外情况，请考虑您是否下载了所有的压缩包分卷并放置于同一目录下。\n"+
				"如果您使用的是WinRAR这款软件，其应能正确提示您需要更多分卷才能完整解压。\n"+
				"由本程序创建的压缩文件不会有任何密码，如果内部还有压缩文件+密码，请检查原始帖文中是否有对应的密码。",
				len(newFiles), chunkIdx*10+1, chunkIdx*10+len(fileChunk), len(newFiles))
			inputMedia[len(inputMedia)-1].(*gotgbot.InputMediaDocument).Caption = captionText
		}
		_, err = bot.SendMediaGroup(supp.LinkedGroupMsg.ChatId, inputMedia, &gotgbot.SendMediaGroupOpts{
			ReplyParameters: &gotgbot.ReplyParameters{
				MessageId:                supp.LinkedGroupMsg.Id,
				ChatId:                   supp.LinkedGroupMsg.ChatId,
				AllowSendingWithoutReply: false,
				Quote:                    "",
				QuoteParseMode:           "",
				QuoteEntities:            nil,
				QuotePosition:            0,
			},
		})
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}
