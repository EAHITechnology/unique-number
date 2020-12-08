package logic

import (
	"golang.org/x/net/context"
)

func GetUniqueNumberLogic(ctx context.Context, req GetUniqueNumberRequest) error {
	//检测本地是否有队列

	//没有阻塞加载

	//有直接取

	//判断剩余，不够异步加载（以后改成增长因子）
	return nil
}
