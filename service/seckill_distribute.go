package service

import (
	"SecKill/cache"
	"SecKill/pkg/e"
	"SecKill/serializer"
	"bytes"
	"errors"
	"fmt"
	logging "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"time"
)


func lockO(myfunc func()) {
	//lock
	//uuid := getUuid()
	//lockSuccess, err := cache.RedisClient.SetNX(key, uuid, time.Second*3).Result()
	//if err != nil || !lockSuccess {
	//	fmt.Println("get lock fail", err)
	//	return
	//} else {
	//	fmt.Println("get lock success")
	//}
	////run func
	//myfunc()
	////unlock
	//value, _ := redisclient .Get(key).Result()
	//if value == uuid { //compare value,if equal then del
	//	_, err := redisclient .Del(key).Result()
	//	if err != nil {
	//		fmt.Println("unlock fail")
	//	}  else {
	//		fmt.Println("unlock success")
	//	}
	//}
}


func WithRedissionSecKill(gid int) serializer.Response {
	code := e.SUCCESS
	seckillNum := 50
	wg.Add(seckillNum)
	InitializerSecKill(gid)
	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := WithRedssionSecKillGoods(gid, userID)
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

func WithRedssionSecKillGoods(gid , userID int) error {
	g := strconv.Itoa(gid)
	uuid := getUuid(g)
	lockSuccess, err := cache.RedisClient.SetNX(g, uuid, time.Second*3).Result()
	if err != nil || !lockSuccess {
		fmt.Println("get lock fail", err)
		return errors.New("get lock fail")
	} else {
		fmt.Println("get lock success")
	}
	err = WithoutLockSecKillGoods(gid, userID)
	if err != nil {
		return err
	}
	value, _ := cache.RedisClient.Get(g).Result()
	if value == uuid { //compare value,if equal then del
		_, err := cache.RedisClient.Del(g).Result()
		if err != nil {
			fmt.Println("unlock fail")
			return nil
		} else {
			fmt.Println("unlock success")
		}
	}
	return nil
}

func getUuid(gid string) string {
	codeLen := 8
	// 1. 定义原始字符串
	rawStr := "jkwangagDGFHGSERKILMJHSNOPQR546413890_"
	// 2. 定义一个buf，并且将buf交给bytes往buf中写数据
	buf := make([]byte, 0, codeLen)
	b := bytes.NewBuffer(buf)
	// 随机从中获取
	rand.Seed(time.Now().UnixNano())
	for rawStrLen := len(rawStr); codeLen > 0; codeLen-- {
		randNum := rand.Intn(rawStrLen)
		b.WriteByte(rawStr[randNum])
	}
	return b.String() + gid
}