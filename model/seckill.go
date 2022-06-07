package model

import (
	"time"
)

type PromotionSecKill struct {
	PsId         int64     `db:"ps_id"`
	GoodsId      int64     `db:"goods_id"`
	PsCount      int64     `db:"ps_count"`
	StartTime    time.Time `db:"start_time"`
	EndTime      time.Time `db:"end_time"`
	Status       int32     `db:"status"`
	CurrentPrice float64   `db:"current_price"`
	Version      int64     `db:"version"`
}

type SuccessKilled struct {
	GoodsId    int64     `db:"goods_id"`
	UserId     int64     `db:"user_id"`
	State      int16     `db:"state"`
	CreateTime time.Time `db:"create_time"`
}


// 查询全部信息
func SelectGoodByGoodsId(gid int) (PromotionSecKill, error) {
	var ps PromotionSecKill
	err := DB.Model(PromotionSecKill{}).Where("goods_id=?", gid).First(&ps).Error
	return ps, err
}

// 更新剩余数量
func UpdateCountByGoodsId(gid int) error {
	return DB.Model(PromotionSecKill{}).Where("goods_id=?", gid).
		Updates(map[string]interface{}{"ps_count": 100, "version": 0}).Error
}

// 查询剩余数量
func SelectCountByGoodsId(gid int) (int64, error) {
	var ps PromotionSecKill
	err := DB.Model(PromotionSecKill{}).Where("goods_id=?", gid).
		First(&ps).Error
	return ps.PsCount, err
}


// 减少特定数量
func ReduceStockByGoodsId(gid int, count int) error {
	return DB.Model(&PromotionSecKill{}).Where("goods_id=?", gid).
		Update("ps_count", count).Error
}

// 减少一个
func ReduceByGoodsId(gid int) (int64, error) {
	var count int64
	sqlStr := `UPDATE promotion_sec_kill SET ps_count = ps_count-1 WHERE ps_count>0 AND goods_id = ?`
	res := DB.Exec(sqlStr, gid)
	if err := res.Error; err != nil {
		return count, err
	}
	count = res.RowsAffected
	return count, nil
}

// 减少一个
func ReduceOneByGoodsId(gid int) error {
	sqlStr := `UPDATE promotion_sec_kill SET ps_count = ps_count-1 WHERE goods_id = ?`
	res := DB.Exec(sqlStr, gid)
	return res.Error
}


// 减少指定数量
func ReduceStockByOcc(gid int, num int, version int) (int64, error) {
	var count int64
	sqlStr := "UPDATE promotion_sec_kill SET ps_count = ps_count-?, version = version+1 " +
		"WHERE version = ? AND goods_id = ?"
	res := DB.Exec(sqlStr, num, version, gid)
	if err := res.Error; err != nil {
		return count, err
	}
	count = res.RowsAffected
	return count, nil
}

// 删除已经秒杀成功的
func DeleteByGoodsId(gid int) error {
	return DB.Where("goods_id=?", gid).Delete(SuccessKilled{}).Error
}

// 创建订单
func CreateOrder(k SuccessKilled) error {
	return DB.Model(SuccessKilled{}).Create(&k).Error
}

// 获取已经成功的
func GetKilledCountByGoodsId(gid int) (int64,error) {
	var count int64
	err := DB.Model(&SuccessKilled{}).Where("goods_id=?", gid).Count(&count).Error
	return count, err
}

// 加读锁
func SelectCountByGoodsIdPcc(gid int) (int64, error) {
	skGood:=PromotionSecKill{}
	err := DB.Model(PromotionSecKill{}).Set("gorm:query_option", "FOR UPDATE").
		Where("goods_id=?",gid).First(&skGood).Error
	return skGood.PsCount, err
}