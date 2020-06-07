package main

import "github.com/mediocregopher/radix.v2/pool"

func pvCounter( pvChannel chan urlData, storageChannel chan storageBlock ) {
	for data := range pvChannel {
		sItem := storageBlock{ "pv", "ZINCRBY", data.unode }
		storageChannel <- sItem
	}
}

func uvCounter( uvChannel chan urlData, storageChannel chan storageBlock, redisPool *pool.Pool ) {
	for data := range uvChannel {
		// HyperLoglog redis
		hyperLogLogKey := "uv_hpll_"+getTime(data.data.time, "day")
		ret, err := redisPool.Cmd( "PFADD", hyperLogLogKey, data.uid, "EX", 86400 ).Int()
		if err!=nil {
			log.Warningln( "UvCounter check redis hyperloglog failed, ", err )
		}
		if ret!=1 {
			continue
		}

		sItem := storageBlock{ "uv", "ZINCRBY", data.unode }
		storageChannel <- sItem
	}
}
