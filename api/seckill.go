package api

import (
	"SecKill/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func WithoutLock(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithoutLockSecKill(gid)
	c.JSON(res.Status, res)
}

func WithLock(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithLockSecKill(gid)
	c.JSON(res.Status, res)
}