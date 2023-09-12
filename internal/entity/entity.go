package entity

type Segments struct {
	BizTag     string `json:"biz_tag" gorm:"'biz_tag'"`
	MaxId      int64  `json:"max_id" gorm:"'max_id'"`
	Step       int64  `json:"step" gorm:"'step'"`
	CreateTime int64  `json:"create_time" gorm:"'create_time'"`
	UpdateTime int64  `json:"update_time" gorm:"'update_time'"`
	Version    int    `json:"version" gorm:"version"`
}

func TableName() string {
	return "segments"
}
