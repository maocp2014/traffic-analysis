package main

import (
	"flag"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var log = logrus.New()

func parseParam() (cmdParams, *string) {
	// 获取参数
	logFilePath := flag.String("logFilePath", "/Users/pangee/Public/nginx/logs/dig.log", "log file path")
	routineNum := flag.Int("routineNum", 5, "consumer numble by goroutine")
	l := flag.String("l", "/tmp/log", "this programe runtime log target file path")
	flag.Parse()

	params := cmdParams{*logFilePath, *routineNum}
	return params, l
}

func writeLog(params cmdParams, l *string) {
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)

	logFd, err := os.OpenFile(*l, os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		log.Out = logFd
		defer logFd.Close()
	}
	log.Infof("Exec start.")
	log.Infof("Params: logFilePath=%s, routineNum=%d", params.logFilePath, params.routineNum)
}

func initRedisPool(params cmdParams) *pool.Pool {
	// Redis Pool
	redisPool, err := pool.New("tcp", "localhost:6379", 2*params.routineNum)
	if err != nil {
		log.Fatalln("Redis pool created failed.")
		panic(err)
	} else {
		go func() {
			for {
				redisPool.Cmd("PING")
				time.Sleep(3 * time.Second)
			}
		}()
	}
	return redisPool
}