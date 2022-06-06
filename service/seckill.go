package service

import (
	"SecKill/model"
	"SecKill/pkg/e"
	"SecKill/serializer"
	"errors"
	"fmt"
	logging "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var wg sync.WaitGroup
var lock sync.Mutex

func InitializerSecKill(gid int) {
	tx := model.DB.Begin()            // 开启事务
	err := model.DeleteByGoodsId(gid) // 删除事务
	if err != nil { // 发生错误的话就进行回滚
		tx.Rollback()
	}
	err = model.UpdateCountByGoodsId(gid) // 更新事务
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
}


// 获取总共秒杀了多少商品
func GetKilledCount(gid int) (int64, error) {
	return model.GetKilledCountByGoodsId(gid)
}


func WithoutLockSecKillGoods(gid, userID int) error {
	tx := model.DB.Begin()
	// 检查库存
	count, err := model.SelectCountByGoodsId(gid)
	if err != nil {
		return err
	}
	if count > 0 {
		// 1. 扣库存
		err = model.ReduceStockByGoodsId(gid, int(count-1))
		if err != nil {
			tx.Rollback()
			return err
		}
		// 2. 创建订单
		kill := model.SuccessKilled{
			GoodsId:    int64(gid),
			UserId:     int64(userID),
			State:      0,
			CreateTime: time.Now(),
		}
		err = model.CreateOrder(kill)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func WithoutLockSecKill(gid int) serializer.Response {
	code := e.SUCCESS
	seckillNum := 50
	wg.Add(seckillNum)
	InitializerSecKill(gid)
	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := WithoutLockSecKillGoods(gid, userID)
			if err != nil {
				fmt.Println("Error",err)
			} else {
				fmt.Printf("User: %d seckill successfully.\n", userID)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	killedCount, err := GetKilledCount(gid)
	if err != nil {
		code = e.ERROR
		logging.Error("Seckill System Error")
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	fmt.Println(killedCount)
	logging.Infof("kill %v product", killedCount)
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func WithLockSecKillGoods(gid,userID int) error {
	lock.Lock()
	err := WithoutLockSecKillGoods(gid, userID)
	lock.Unlock()
	return err
}

func WithLockSecKill(gid int) serializer.Response {
	code := e.SUCCESS
	seckillNum := 50
	wg.Add(seckillNum)
	InitializerSecKill(gid)
	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := WithLockSecKillGoods(gid, userID)
			if err != nil {
				code = e.ERROR
				logging.Errorln("Error", err)
			} else {
				logging.Infof("User: %d seckill successfully.\n", userID)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	kCount, err := GetKilledCount(gid)
	if err != nil {
		code = e.ERROR
		logging.Infoln("Error")
	}
	logging.Infof("Total %v goods", kCount)
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func WithPccReadSecKillGoods(gid, userID int) error {
	tx := model.DB.Begin()
	count, err := model.SelectCountByGoodsIdPcc(gid)
	// 先读后更新的数据竞争场景
	if err != nil {
		return err
	}
	if count > 0 {
		// 1. 扣库存
		err = model.ReduceStockByGoodsId(gid, int(count-1))
		if err != nil {
			tx.Rollback()
			return err
		}
		// 2. 创建订单
		kill := model.SuccessKilled{
			GoodsId:    int64(gid),
			UserId:     int64(userID),
			State:      0,
			CreateTime: time.Now(),
		}
		err = model.CreateOrder(kill)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func WithPccReadSecKill(gid int) serializer.Response {
	code := e.SUCCESS
	seckillNum := 50
	wg.Add(seckillNum)
	InitializerSecKill(gid)
	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := WithPccReadSecKillGoods(gid, userID)
			if err != nil {
				code = e.ERROR
				logging.Errorln("Error", err)
			} else {
				logging.Infof("User: %d seckill successfully.\n", userID)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	kCount, err := GetKilledCount(gid)
	if err != nil {
		code = e.ERROR
		logging.Infoln("Error")
	}
	logging.Infof("Total %v goods", kCount)
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func WithPccUpdateSecKillGoods(gid, userID int) error {
	tx := model.DB.Begin()
	// 1. 扣库存
	count, err := model.ReduceByGoodsId(gid)
	if err != nil {
		return err
	}
	if count > 0 {
		// 2. 创建订单
		kill := model.SuccessKilled{
			GoodsId:    int64(gid),
			UserId:     int64(userID),
			State:      0,
			CreateTime: time.Now(),
		}
		err = model.CreateOrder(kill)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}


func WithPccUpdateSecKill(gid int) serializer.Response {
	code := e.SUCCESS
	seckillNum := 50
	wg.Add(seckillNum)
	InitializerSecKill(gid)
	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := WithPccUpdateSecKillGoods(gid, userID)
			if err != nil {
				code = e.ERROR
				logging.Errorln("Error", err)
			} else {
				logging.Infof("User: %d seckill successfully.\n", userID)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	kCount, err := GetKilledCount(gid)
	if err != nil {
		code = e.ERROR
		logging.Infoln("Error")
	}
	logging.Infof("Total %v goods", kCount)
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func WithOccSecKillGoods(gid, userID,num int) error {
	tx := model.DB.Begin()
	good, err := model.SelectGoodByGoodsId(gid)
	if err != nil {
		return err
	}
	if int(good.PsCount) >= num {
		// 1. 扣库存
		count, err := model.ReduceStockByOcc(gid, num, int(good.Version))
		if err != nil {
			tx.Rollback()
			return err
		}
		if count > 0 {
			// 2. 创建订单
			kill := model.SuccessKilled{
				GoodsId:    int64(gid),
				UserId:     int64(userID),
				State:      0,
				CreateTime: time.Now(),
			}
			err = model.CreateOrder(kill)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
		}
	} else {
		tx.Rollback()
		return errors.New("库存不够了")
	}
	tx.Commit()
	return nil
}

func WithOccSecKill(gid int) serializer.Response {
	code := e.SUCCESS
	seckillNum := 50
	wg.Add(seckillNum)
	InitializerSecKill(gid)
	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := WithOccSecKillGoods(gid, userID, 1)
			if err != nil {
				code = e.ERROR
				logging.Errorln("Error", err)
			} else {
				logging.Infof("User: %d seckill successfully.\n", userID)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	kCount, err := GetKilledCount(gid)
	if err != nil {
		code = e.ERROR
		logging.Infoln("Error")
	}
	logging.Infof("Total %v goods", kCount)
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func WithChannelSecKillGoods(gid, userID int) error {
	kill := [2]int{gid, userID}
	kChan := GetInstance()
	*kChan <- kill
	return nil
}

func ChannelConsumer() {
	for {
		kill, ok := <-(*GetInstance())
		if !ok {
			continue
		}
		err := WithoutLockSecKillGoods(kill[0], kill[1])
		if err != nil {
			logging.Error("Error")
		} else {
			logging.Infof("User:%v SecKill Successfully", kill[1])
		}
	}
}

func WithChannelSecKill(gid int) serializer.Response {
	code := e.SUCCESS
	seckillNum := 50
	wg.Add(seckillNum)
	InitializerSecKill(gid)
	go ChannelConsumer()
	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := WithChannelSecKillGoods(gid, userID)
			if err != nil {
				code = e.ERROR
				logging.Errorln("Error", err)
			} else {
				logging.Infof("User: %d seckill successfully.\n", userID)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	kCount, err := GetKilledCount(gid)
	if err != nil {
		code = e.ERROR
		logging.Infoln("Error")
	}
	logging.Infof("Total %v goods", kCount)
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}