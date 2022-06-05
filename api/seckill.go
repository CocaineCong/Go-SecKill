package api

import (
	"SecKill/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Normal(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.NormalSecKill(gid)
	c.JSON(res.Status, res)
}
