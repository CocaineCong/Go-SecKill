package service

import (
	"SecKill/cache"
	"SecKill/pkg/e"
	"SecKill/serializer"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	logging "github.com/sirupsen/logrus"
	//"go.etcd.io/etcd/clientv3"
	"math/rand"
	"strconv"
	"time"
)

type EtcdMutex struct {
	Ttl     int64              //租约时间
	Conf    clientv3.Config  //etcd集群配置
	Key     string             //etcd的key
	cancel  context.CancelFunc //关闭续租的func
	lease   clientv3.Lease
	leaseID clientv3.LeaseID
	txn     clientv3.Txn
}

func(em *EtcdMutex)initETCD()error {
	var err error
	var ctx context.Context
	client, err := clientv3.New(em.Conf)
	if err != nil {
		return err
	}
	em.txn = clientv3.NewKV(client).Txn(context.TODO())
	em.lease = clientv3.NewLease(client)
	leaseResp, err := em.lease.Grant(context.TODO(), em.Ttl)
	if err != nil {
		return err
	}
	ctx, em.cancel = context.WithCancel(context.TODO())
	em.leaseID = clientv3.LeaseID(leaseResp.ID)
	_, err = em.lease.KeepAlive(ctx, em.leaseID)
	return err
}

func(em *EtcdMutex)Lock()error{
	err := em.initETCD()
	if err != nil{
		return err
	}
	em.txn.If(clientv3.Compare(clientv3.CreateRevision(em.Key),"=",0)).
		Then(clientv3.OpPut(em.Key,"",clientv3.WithLease(em.leaseID))).
		Else()
	txnResp,err := em.txn.Commit()
	if err != nil{
		return err
	}
	if !txnResp.Succeeded{   //判断txn.if条件是否成立
		return fmt.Errorf("抢锁失败")
	}
	return nil
}

func(em *EtcdMutex)UnLock() {
	em.cancel()
	_, _ = em.lease.Revoke(context.TODO(), em.leaseID)
	fmt.Println("释放了锁")
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

func WithETCDSecKillGoods(gid, userID int) error {
	var conf = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	eMutex1 := &EtcdMutex{
		Conf: conf,
		Ttl:  10,
		Key:  "lock",
	}
	err := eMutex1.Lock()
	if err != nil {
		return err
	}
	err = WithoutLockSecKillGoods(gid, userID)
	eMutex1.UnLock()
	return err
}

func WithETCDSecKill(gid int) serializer.Response {
	code := e.SUCCESS
	seckillNum := 50
	wg.Add(seckillNum)
	InitializerSecKill(gid)

	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := WithETCDSecKillGoods(gid, userID)
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