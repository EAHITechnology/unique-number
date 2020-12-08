package work

import (
	"errors"
	"strings"
	"sync/atomic"

	"github.com/EAHITechnology/inf/unique-number/logic/work/mist"
	"github.com/EAHITechnology/inf/unique-number/utils"
	"github.com/coocood/freecache"
)

/*
唯一号生成器管理器
*/
type Worker struct {
	Workers *freecache.Cache
}

/*
唯一号生成器
*/
type workInfo struct {
	Topic           string     // 队列名
	Channel         chan int64 // 队列存储管道
	Max             int64      // 本地当前最大值
	Offset          int64      // 消费位置
	ExpansionStatus int32      // 扩充状态
}

/*
用于转换存储结对象地址
*/
type tempAddr struct {
	addr uintptr
	len  int
	cap  int
}

/*
内部方法，初始加载唯一号
*/
func (this *workInfo) loadIncreasNum(step int64) error {
	maxId, err := cacheManager.SetMaxId(this.Topic, step)
	if err != nil {
		return err
	}
	head := maxId - utils.Step
	end := maxId
	for head <= end {
		uniqueNum := mist.Generate(head)
		this.Channel <- uniqueNum
		head++
	}
	this.Max = maxId
	this.Offset = 0
	return nil
}

/*
初始化函数,初始化会提供一个默认全集群公用的唯一数队列,名字为base
*/
func InitWorker() (*Worker, error) {
	//初始化工作对象管理器
	cacheSize := 100 * 1024 * 1024
	cache := freecache.NewCache(cacheSize)
	//填充配置
	channel := make(chan int64, 10000)
	baseWorkInfo := workInfo{
		Topic:   "base",
		Channel: channel,
	}
	//加载唯一号信息
	if err := baseWorkInfo.loadIncreasNum(utils.Step); err != nil {
		return nil, err
	}
	//存储生成器信息
	data := workInfo2Byte(baseWorkInfo)
	if err := cache.Set(utils.Str2bytes(baseWorkInfo.Topic), data, 0); err != nil {
		return nil, err
	}
	return &Worker{Workers: cache}, nil
}

/*
获得或者添加一个工作者，判断本地有目标工作者则返回工作者指针，没有则创建工作者并返回nil
*/
func (this *Worker) GetOrAddWorker(topic string) (*workInfo, error) {
	if this == nil || this.Workers == nil {
		return nil, errors.New("Worker nil")
	}
	//优先拿取
	res, err := this.Workers.Get(utils.Str2bytes(topic))
	if err == nil {
		return byte2WorkInfo(res), nil
	}
	if !strings.Contains(err.Error(), "not found") {
		return nil, err
	}
	//填充配置
	channel := make(chan int64, 10000)
	otherWorkInfo := workInfo{
		Topic:   topic,
		Channel: channel,
	}
	//加载唯一号信息
	if err := otherWorkInfo.loadIncreasNum(utils.Step); err != nil {
		return nil, err
	}
	//存储生成器信息
	data := workInfo2Byte(otherWorkInfo)
	if err := this.Workers.Set(utils.Str2bytes(otherWorkInfo.Topic), data, 0); err != nil {
		return nil, err
	}
	return nil, nil
}

/*
获得唯一号
*/
func (this *workInfo) GetUniqueNum() (int64, error) {
	//先判断剩余条件
	remaining := this.Max - this.Offset
	if remaining < utils.ExpansionRemaining {
		if swapOk := atomic.CompareAndSwapInt32(&this.ExpansionStatus, 0, 1); swapOk {
			this.expansionUniqueNum()
		}
	}
	uniqueNum := <-this.Channel
	this.Offset++
	return uniqueNum, nil
}

/*
扩容唯一数，当判断唯一数数量到达阈值则打开写锁阻塞读请求直到扩充完毕
原因是因为消费太快在扩充时将剩余唯一数消费完.
*/
func (this *workInfo) expansionUniqueNum() error {
	if err := this.loadIncreasNum(this.Offset); err != nil {
		return err
	}

	return nil
}
