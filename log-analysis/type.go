package main

import (
	"github.com/mgutz/str"
	"strconv"
	"time"
)

const HANDLE_DIG = " /dig?"
const HANDLE_MOVIE = "/movie/"
const HANDLE_LIST = "/list/"
const HANDLE_HTML = ".html"


type cmdParams struct {
	logFilePath string
	routineNum int
}
type digData struct{
	time   string
	url    string
	refer  string
	ua        string
}
type urlData struct {
	data   digData
	uid    string
	unode  urlNode
}
type urlNode struct {
	unType     string // 详情页 或者 列表页 或者 首页
	unRid  int        // Resource ID 资源ID
	unUrl  string     // 当前这个页面的url
	unTime  string    // 当前访问这个页面的时间
}
type storageBlock struct {
	counterType       string
	storageModel   string
	unode        urlNode
}

func formatUrl( url, t string ) urlNode{
	// 一定从量大的着手,  详情页>列表页≥首页
	pos1 := str.IndexOf( url, HANDLE_MOVIE, 0)
	if pos1!=-1 {
		pos1 += len( HANDLE_MOVIE )
		pos2 := str.IndexOf( url, HANDLE_HTML, 0 )
		idStr := str.Substr( url , pos1, pos2-pos1 )
		id, _ := strconv.Atoi( idStr )
		return urlNode{ "movie", id, url, t }
	} else {
		pos1 = str.IndexOf( url, HANDLE_LIST, 0 )
		if pos1!=-1 {
			pos1 += len( HANDLE_LIST )
			pos2 := str.IndexOf( url, HANDLE_HTML, 0 )
			idStr := str.Substr( url , pos1, pos2-pos1 )
			id, _ := strconv.Atoi( idStr )
			return urlNode{ "list", id, url, t }
		} else {
			return urlNode{ "home", 1, url, t}
		} // 如果页面url有很多种，就不断在这里扩展
	}
}

func getTime( logTime, timeType string ) string {
	var item string
	switch timeType {
	case "day":
		item = "2006-01-02"
		break
	case "hour":
		item = "2006-01-02 15"
		break
	case "min":
		item = "2006-01-02 15:04"
		break
	}
	t, _ := time.Parse( item, time.Now().Format(item) )
	return strconv.FormatInt( t.Unix(), 10 )
}