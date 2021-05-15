package data_base

import (
	redisDb "aapanavyapar-service-buying/configurations/redis"
	"github.com/go-redis/redis/v8"
	"sync"
)

type CashDataBase struct {
	mutex sync.RWMutex
	Cash  *redis.Client
}

func NewDataBase() *CashDataBase {

	rdb := redisDb.InitRedis()

	return &CashDataBase{
		Cash: rdb,
	}
}
