package main

import (
	"time"
)

func main() {
	// 获取参数
	params, l := parseParam()

	// 打日志
	writeLog(params, l)

	// 初始化一些channel，用于数据传递
	var logChannel = make(chan string, 3*params.routineNum)
	var pvChannel = make(chan urlData, params.routineNum)
	var uvChannel = make(chan urlData, params.routineNum)
	var storageChannel = make(chan storageBlock, params.routineNum)

	// Redis Pool
	redisPool := initRedisPool(params)

	// 日志消费者
	go readFileLinebyLine(params, logChannel)

	// 创建一组日志处理
	for i := 0; i < params.routineNum; i++ {
		go logConsumer(logChannel, pvChannel, uvChannel)
	}

	// 创建PV UV 统计器
	go pvCounter(pvChannel, storageChannel)
	go uvCounter(uvChannel, storageChannel, redisPool)
	// 可扩展的 xxxCounter

	// 创建 存储器
	go dataStorage(storageChannel, redisPool)

	time.Sleep(1000 * time.Second)
}