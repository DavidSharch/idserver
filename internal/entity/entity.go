package entity

type Segments struct {
	Id         int64  `json:"id",gorm:"primarykey"`
	BizTag     string `json:"biz_tag",gorm:"column:biz_tag"`
	MaxId      int64  `json:"max_id",gorm:"column:max_id"`
	Step       int64  `json:"step",gorm:"column:step"`
	CreateTime int64  `json:"create_time",gorm:"column:create_time"`
	UpdateTime int64  `json:"update_time",gorm:"column:update_time"`
}

func (s Segments) TableName() string {
	return "segments"
}
