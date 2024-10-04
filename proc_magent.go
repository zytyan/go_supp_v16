package main

import (
	"context"
	"log"
	"main/qbit"
	"sync"
	"time"
)

var cachedTorrents struct {
	mu         sync.RWMutex
	torrents   map[string]qbit.Torrent
	lastUpdate time.Time
}

func getTorrentsCached() (map[string]qbit.Torrent, error) {
	cachedTorrents.mu.RLock()
	if time.Since(cachedTorrents.lastUpdate) < 5*time.Second {
		cachedTorrents.mu.RUnlock()
		return cachedTorrents.torrents, nil
	}
	cachedTorrents.mu.RUnlock()
	torrents, err := qClient.GetTorrents()
	if err != nil {
		return nil, err
	}
	cachedTorrents.mu.Lock()
	defer cachedTorrents.mu.Unlock()
	cachedTorrents.lastUpdate = time.Now()
	var res = make(map[string]qbit.Torrent, len(torrents))
	for _, t := range torrents {
		res[t.Hash] = t
	}
	cachedTorrents.torrents = res
	return res, nil
}

func DownloadMagnet(hash []string) error {
	torrents, err := qClient.GetTorrents()
	if err != nil {
		return err
	}
	// 删除已经存在的hash，否则会报错
	ts := make(map[string]struct{}, len(torrents))
	for _, t := range torrents {
		ts[t.Hash] = struct{}{}
	}
	newHash := make([]string, 0, len(hash))
	for _, h := range hash {
		if _, ok := ts[h]; !ok {
			newHash = append(newHash, h)
		}
	}
	if len(newHash) == 0 {
		return nil
	}
	return qClient.DownloadMagnetUrls(newHash)
}

func WaitAndProcMagnet(ctx context.Context, supp *Supp, hash string) error {
	countNotInTorrents := 0
	for {
		torrents, err := getTorrentsCached()
		if err != nil {
			log.Println(err)
			return err
		}
		countNotInTorrents++
		torrent, ok := torrents[hash]
		if !ok {
			log.Printf("magnet %s not in torrents, countNotInTorrents: %d", hash, countNotInTorrents)
			if countNotInTorrents > 20 {
				return nil
			}
		}
		if torrent.Progress == 1 {
			uploadVideos(&torrent, supp)
			return UploadRawFiles(&torrent, supp)
		}
		countNotInTorrents = 0
		time.Sleep(10 * time.Second)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}
