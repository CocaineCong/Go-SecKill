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


type GoodsCover struct {
	GcId       int64  `db:"gc_id"`
	GoodsId    int64  `db:"goods_id"`
	GcPicUrl   string `db:"gc_pic_url"`
	GcThumbUrl string `db:"gc_thumb_url"`
	GcOrder    int64  `db:"gc_order"`
}

type GoodsDetail struct {
	GdId     int64  `db:"gd_id"`
	GoodsId  int64  `db:"goods_id"`
	GdPicUrl string `db:"gd_pic_url"`
	GdOrder  int32  `db:"gd_order"`
}

type GoodsParam struct {
	GdId         int64  `db:"gp_id"`
	GpParamName  string `db:"gp_param_name"`
	GpParamValue string `db:"gp_param_value"`
	GoodsId      int64  `db:"goods_id"`
	GdOrder      int32  `db:"gp_order"`
}

func FindGoodsById(goodsID int) (Goods,error) {
	var good Goods
	err := DB.Where("goods_id=?", goodsID).First(&good).Error
	if err != nil {
		return good, err
	}
	return good, nil
}