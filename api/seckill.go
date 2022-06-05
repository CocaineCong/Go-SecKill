package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func Normal(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	seckillNum := 50

}
