package service

import (
	"SecKill/model"
	"SecKill/pkg/e"
	"SecKill/serializer"
)

// 获取商品的详细信息
func GetGoodDetailList(gid int) serializer.Response {
	code := e.SUCCESS
	good, err := model.FindGoodsById(gid)
	if err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   good,
		Msg:    e.GetMsg(code),
	}
}