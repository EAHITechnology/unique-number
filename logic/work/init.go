package work

import (
	"gitlab.sftcwl.com/tc-inf/unique-number/cache"
	"golang.org/x/net/context"
)

var (
	cacheManager *cache.CacheManager
)

func InitManager(ctx context.Context) (err error) {
	cacheManager, err = cache.New()
	if err != nil {
		return
	}
	return
}
