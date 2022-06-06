package api

import (
	"SecKill/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func WithRedission(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithRedissionSecKill(gid)
	c.JSON(res.Status, res)
}