package work

import (
	"time"

	"github.com/EAHITechnology/inf/golang/log"
	"github.com/EAHITechnology/inf/unique-number/utils"
	"golang.org/x/net/context"
)

type expansionErrorReturn struct {
	Topic string
	err   error
}

type expansionSuccessReturn struct {
	Topic  string
	Offset int64
	Max    int64
}

/*
并发扩容，令牌控制并发
*/
func (this *Worker) expansionBackend(ctx context.Context) error {
	workInfos := []*workInfo{}
	//获得工作者
	rangeKey := this.Workers.NewIterator()
	for {
		key := rangeKey.Next()
		if key == nil {
			break
		}
		rangeWorkInfo := byte2WorkInfo(key.Value)
		workInfos = append(workInfos, rangeWorkInfo)
	}
	//设置并发参数
	concurrencyTaskNum := make(chan struct{}, utils.LoadFrequency)
	successTaskNum := make(chan expansionSuccessReturn, len(workInfos))
	errorChan := make(chan expansionErrorReturn, len(workInfos))
	defer close(concurrencyTaskNum)
	defer close(successTaskNum)
	defer close(errorChan)
	//开始并发
	for _, w := range workInfos {
		log.DebugfCtx(ctx, "expansion start:%s", w.Topic)
		concurrencyTaskNum <- struct{}{}
		go func(w *workInfo) {
			if err := w.expansionUniqueNum(); err != nil {
				errorChan <- expansionErrorReturn{Topic: w.Topic, err: err}
				<-concurrencyTaskNum
				return
			}
			successTaskNum <- expansionSuccessReturn{w.Topic, w.Offset, w.Max}
			<-concurrencyTaskNum
		}(w)
	}
	//等待结果
	times := 0
	for {
		select {
		case <-time.After(time.Duration(3) * time.Second):
			log.Warn("同步超时，请关注")
		case err := <-errorChan:
			log.Errorf("topic:%s expansion err:%s", err.Topic, err)
			//计算结束
			times++
			if times >= len(workInfos) {
				return nil
			}
		case successRes := <-successTaskNum:
			log.Infof("topic:%s expansion offset:%s max:%d", successRes.Topic, successRes.Offset, successRes.Max)
			//计算结束
			times++
			if times >= len(workInfos) {
				return nil
			}
		}
	}
}
