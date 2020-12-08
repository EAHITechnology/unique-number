package cache

import (
	"github.com/EAHITechnology/inf/golang/eredis"
	"github.com/EAHITechnology/inf/golang/log"
)

type CacheManager struct {
	lockRds *eredis.Redis
}

func New() (*CacheManager, error) {
	lockRds, err := eredis.GetClient(LOCK_REDIS)
	if err != nil {
		log.Errorf("LOCK_REDIS init error:%s", err.Error())
		return nil, err
	}
	return &CacheManager{lockRds: lockRds}, nil
}
