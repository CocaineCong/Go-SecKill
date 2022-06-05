package api

import (
	"SecKill/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetGoodDetail(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.GetGoodDetailList(gid)
	c.JSON(res.Status, res)
}
