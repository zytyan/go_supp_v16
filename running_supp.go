package main

import (
	"fmt"
	"sync"
)

type runningSuppType struct {
	mu        sync.Mutex
	byUrlPath map[string]*Supp
	byMsgId   map[Msg]*Supp
}

var runningSupp = runningSuppType{
	mu:        sync.Mutex{},
	byUrlPath: make(map[string]*Supp),
	byMsgId:   make(map[Msg]*Supp),
}

func (*runningSuppType) Add(supp *Supp) {
	if supp.ChannelMsg.Id == 0 || supp.ArticleUrlPath == "" || supp.ChannelMsg.ChatId == 0 {
		panic(fmt.Sprintf("invalid supp: %+v", supp))
	}
	runningSupp.mu.Lock()
	defer runningSupp.mu.Unlock()
	runningSupp.byUrlPath[supp.ArticleUrlPath] = supp
	runningSupp.byMsgId[supp.ChannelMsg] = supp
}

func (*runningSuppType) Remove(supp *Supp) {
	runningSupp.mu.Lock()
	defer runningSupp.mu.Unlock()
	delete(runningSupp.byUrlPath, supp.ArticleUrlPath)
	delete(runningSupp.byMsgId, supp.ChannelMsg)
}

func (*runningSuppType) GetByMsg(msg Msg) (*Supp, bool) {
	runningSupp.mu.Lock()
	defer runningSupp.mu.Unlock()
	res, ok := runningSupp.byMsgId[msg]
	return res, ok
}

func (*runningSuppType) GetByUrlPath(urlPath string) (*Supp, bool) {
	runningSupp.mu.Lock()
	defer runningSupp.mu.Unlock()
	res, ok := runningSupp.byUrlPath[urlPath]
	return res, ok
}

func (*runningSuppType) Size() int {
	runningSupp.mu.Lock()
	defer runningSupp.mu.Unlock()
	return len(runningSupp.byUrlPath)
}
