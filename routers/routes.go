package routers

import (
	"SecKill/api"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})
	// 商品信息展示页面获取数据
	r.GET("/good", api.GetGoodDetail)

	// 单机锁
	skGroup := r.Group("/api/v1")
	{
		// 不加锁,出现超卖现象
		skGroup.GET("/without-lock", api.WithoutLock)
		// 加锁(sync包中的Mutex类型的互斥锁),没有问题
		skGroup.GET("/with-lock", api.WithLock)
		// 加锁(数据库悲观锁，查询加锁)
		skGroup.GET("/with-pcc-read",api.WithPccRead)
		// 加锁(数据库悲观锁，更新限定)
		skGroup.GET("/with-pcc-update",api.WithPccUpdate)
		// 加锁(数据库乐观锁，正常)
		skGroup.GET("/with-occ",api.WithOcc)
		// channel 限制，正常
		skGroup.GET("/with-channel",api.WithChannel)
	}
	return r
}
