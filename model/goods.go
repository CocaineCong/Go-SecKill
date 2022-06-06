package model

import "time"

type Goods struct {
	GoodsId        int64     `db:"goods_id"`
	Title          string    `db:"title"`
	SubTitle       string    `db:"sub_title"`
	OriginalCost   float64   `db:"original_cost"`
	CurrentPrice   float64   `db:"current_price"`
	Discount       float64   `db:"discount"`
	IsFreeDelivery int32     `db:"is_free_delivery"`
	CategoryId     int64     `db:"category_id"`
	LastUpdateTime time.Time `db:"last_update_time"`
}

func FindGoodsById(goodsID int) (Goods,error) {
	var good Goods
	err := DB.Where("goods_id=?", goodsID).First(&good).Error
	if err != nil {
		return good, err
	}
	return good, nil
}