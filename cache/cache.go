package cache

import (
	"errors"
	"strconv"
)

/*
通过步长设置appname的自增边界,并返回目前最大id值
*/
func (this *CacheManager) SetMaxId(appname string, step int64) (int64, error) {
	resIf, err := this.lockRds.Exec("INCRBY", appname, step)
	if err != nil {
		return 0, err
	}
	if resIf == nil {
		return 0, errors.New("SetMaxId INCRBY resIf nil")
	}
	return resIf.(int64), nil
}

/*
获得目前最大id值
*/
func (this *CacheManager) GetMaxId(appname string, step int64) (int64, error) {
	resIf, err := this.lockRds.Exec("GET", appname)
	if err != nil {
		return 0, err
	}
	if resIf == nil {
		return 0, errors.New("GetMaxId GET resIf nil")
	}
	resByt := resIf.([]byte)
	res, err := strconv.ParseInt(Bytes2str(resByt), 10, 64)
	if err != nil {
		return 0, err
	}
	return res, nil
}
