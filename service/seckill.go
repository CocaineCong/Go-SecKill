package service

import (
	"SecKill/model"
	"SecKill/pkg/e"
	"SecKill/serializer"
	"fmt"
	logging "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var wg sync.WaitGroup

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

func NormallSecKillGoods(gid, userID int) error {
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

//
func GetKilledCount(gid int) (int64, error) {
	return model.GetKilledCountByGoodsId(gid)
}

func NormalSecKill(gid int) serializer.Response {
	code := e.SUCCESS
	seckillNum := 50
	wg.Add(seckillNum)
	InitializerSecKill(gid)
	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := NormallSecKillGoods(gid, userID)
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
