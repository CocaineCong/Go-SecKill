package routers

import (
	"SecKill/api"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("ping", func(c *gin.Context) {
		c.JSON(200,"pong")
	})
	// 商品信息展示页面获取数据
	r.GET("/good", api.GetGoodDetail)

	// 单机锁
	//seckillGroup := r.Group("/seckill")
	//{
	//	// case1:不加锁,出现超卖现象
	//	seckillGroup.GET("/handle", api.Handle)
	//	// case2:使用sync包中的Mutex类型的互斥锁,秒杀正常
	//	seckillGroup.GET("/handleWithLock", api.HandleWithLock)
	//	// case4:数据库悲观锁(查询加锁),不能
	//	seckillGroup.GET("/handleWithPccOne", api.HandleWithPccOne)
	//	// case5:数据库悲观锁(更新加锁),正常
	//	seckillGroup.GET("/handleWithPccTwo", api.HandleWithPccTwo)
	//	// case6:数据库乐观锁，正常
	//	seckillGroup.GET("/handleWithOcc", api.HandleWithOcc)
	//	// case7:GoLang中的channel，正常
	//	seckillGroup.GET("/handleWithChannel", api.HandleWithChannel)
	//}
	return r
}

