package main

import "github.com/mediocregopher/radix.v2/pool"

// HBase 劣势：列簇需要声明清楚
func dataStorage( storageChannel chan storageBlock, redisPool *pool.Pool) {
	for block := range storageChannel {
		prefix := block.counterType + "_"

		// 逐层添加，加洋葱皮的过程
		// 维度： 天-小时-分钟
		// 层级： 定级-大分类-小分类-终极页面
		// 存储模型： Redis  SortedSet
		setKeys := []string{
			prefix+"day_"+getTime(block.unode.unTime, "day"),
			prefix+"hour_"+getTime(block.unode.unTime, "hour"),
			prefix+"min_"+getTime(block.unode.unTime, "min"),
			prefix+block.unode.unType+"_day_"+getTime(block.unode.unTime, "day"),
			prefix+block.unode.unType+"_hour_"+getTime(block.unode.unTime, "hour"),
			prefix+block.unode.unType+"_min_"+getTime(block.unode.unTime, "min"),
		}

		rowId := block.unode.unRid

		for _,key := range setKeys {
			ret, err := redisPool.Cmd( block.storageModel, key, 1, rowId ).Int()
			if ret<=0 || err!=nil {
				log.Errorln( "DataStorage redis storage error.", block.storageModel, key, rowId )
			}
		}
	}
}
