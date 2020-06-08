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

// 日志消费
func logConsumer(logChannel chan string, pvChannel, uvChannel chan urlData) error {
	for logStr := range logChannel {
		// 切割日志字符串，扣出打点上报的数据
		data := cutLogFetchData(logStr)

		// uid
		// 说明：课程中模拟生成uid，md5(refer + ua)
		harsher := md5.New()
		harsher.Write([]byte(data.refer + data.ua))
		uid := hex.EncodeToString(harsher.Sum(nil))

		// 很多解析的工作都可以放到这里完成
		// ...
		// ...

		uData := urlData{data, uid, formatUrl( data.url, data.time )}
        // 数据塞入 channel
		pvChannel <- uData
		uvChannel <- uData
	}
	return nil
}

// 日志切割
func cutLogFetchData(logStr string) digData {
	// 去除空格
	logStr = strings.TrimSpace(logStr)
	// 查找索引
	pos1 := str.IndexOf(logStr,  HANDLE_DIG, 0)
	if pos1 == -1 {
		return digData{}
	}
	pos1 += len(HANDLE_DIG)
	pos2 := str.IndexOf(logStr, " HTTP/", pos1)
	// 返回子串
	d := str.Substr(logStr, pos1, pos2-pos1)

	// url.Parse 必须构造成 url
	urlInfo, err := url.Parse("http://localhost/?" + d)
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

// 读取日志文件
func readFileLineByLine(params cmdParams, logChannel chan string) error {
	fd, err := os.Open(params.logFilePath)
	if err != nil {
		log.Warningf("ReadFileLineByLine can't open file: %s", params.logFilePath)
		return err
	}
	defer fd.Close()

	count := 0
	// 缓冲读
	bufferRead := bufio.NewReader(fd)

	for {
		// 读取字符串，以"\n"分割
		line, err := bufferRead.ReadString('\n')
		// 塞入channel
		logChannel <- line
		count++
		// 每 1000 行打印日志
		if count % (1000*params.routineNum) == 0 {
			log.Infof("ReadFileLineByLine line: %d", count)
		}
		if err != nil {
			// 文件结尾
			if err == io.EOF {
				time.Sleep( 3*time.Second )
				log.Infof( "ReadFileLineByLine wait, readline: %d", count )
			} else {
				log.Warningf( "ReadFileLineByLine read log error" )
			}
		}
	}
}
