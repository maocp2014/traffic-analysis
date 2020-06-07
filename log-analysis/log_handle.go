package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"github.com/mgutz/str"
	"io"
	"net/url"
	"os"
	"strings"
	"time"
)

func logConsumer( logChannel chan string, pvChannel, uvChannel chan urlData ) error {
	for logStr := range logChannel {
		// 切割日志字符串，扣出打点上报的数据
		data := cutLogFetchData( logStr )

		// uid
		// 说明： 课程中模拟生成uid， md5(refer+ua)
		hasher := md5.New()
		hasher.Write( []byte( data.refer+data.ua ) )
		uid := hex.EncodeToString( hasher.Sum(nil) )

		// 很多解析的工作都可以放到这里完成
		// ...
		// ...

		uData := urlData{ data, uid, formatUrl( data.url, data.time ) }

		pvChannel <- uData
		uvChannel <- uData
	}
	return nil
}
func cutLogFetchData( logStr string ) digData {
	logStr = strings.TrimSpace( logStr )
	pos1 := str.IndexOf( logStr,  HANDLE_DIG, 0)
	if pos1==-1 {
		return digData{}
	}
	pos1 += len( HANDLE_DIG )
	pos2 := str.IndexOf( logStr, " HTTP/", pos1 )
	d := str.Substr( logStr, pos1, pos2-pos1 )

	urlInfo, err := url.Parse( "http://localhost/?"+d )
	if err != nil {
		return digData{}
	}
	data := urlInfo.Query()
	return digData{
		data.Get("time"),
		data.Get("refer"),
		data.Get("url"),
		data.Get("ua"),
	}
}
func readFileLinebyLine( params cmdParams, logChannel chan string ) error {
	fd, err := os.Open( params.logFilePath )
	if err != nil {
		log.Warningf( "ReadFileLinebyLine can't open file:%s", params.logFilePath )
		return err
	}
	defer fd.Close()

	count := 0
	bufferRead := bufio.NewReader( fd )
	for {
		line, err := bufferRead.ReadString( '\n' )
		logChannel <- line
		count++

		if count%(1000*params.routineNum) == 0 {
			log.Infof( "ReadFileLinebyLine line: %d", count )
		}
		if err != nil {
			if err == io.EOF {
				time.Sleep( 3*time.Second )
				log.Infof( "ReadFileLinebyLine wait, raedline:%d", count )
			} else {
				log.Warningf( "ReadFileLinebyLine read log error" )
			}
		}
	}
}
