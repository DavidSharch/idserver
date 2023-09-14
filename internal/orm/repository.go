package orm

import (
	"errors"
	"github.com/sharch/idserver/internal/entity"
	"github.com/sharch/idserver/internal/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	c  *MysqlConfig
	db *gorm.DB
}

func NewRepository(c *MysqlConfig) (r *Repository) {
	r = &Repository{c: c, db: NewMysql(c)}
	return
}

// SegmentsGetAll 获取最近6小时的数据，预热
func (r *Repository) SegmentsGetAll() (res []entity.Segments, err error) {
	if errs := r.db.Where("update_time >= ?", time.Now().Unix()-21600).Find(&res); errs != nil {
		log.GetLogger().Error("[Repository] SegmentsGetAll Find", zap.Error(err))
	}
	return
}

// SegmentsCreate 新建
func (r *Repository) SegmentsCreate(s *entity.Segments) (data *entity.Segments, err error) {
	var cnt int64
	r.db.Where("biz_tag = ?", s.BizTag).Count(&cnt)
	if cnt > 0 {
		return nil, errors.New("已经存在")
	}
	s.CreateTime = time.Now().Unix()
	s.UpdateTime = time.Now().Unix()
	err = r.db.Create(s).Error
	if err != nil {
		return
	}
	return s, nil
}

// SegmentsIdNext 获取下一个号段
func (r *Repository) SegmentsIdNext(tag string) (id *entity.Segments, err error) {
	updateSQL := "update segments set max_id=max_id+step,update_time = ? where biz_tag = ?"
	r.db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Exec(updateSQL, time.Now().Unix(), tag).Error; err != nil {
			log.GetLogger().Error("update failed", zap.Error(err), zap.String("tag", tag))
			// return err 自动回滚
			return err
		}
		if err = tx.Where("biz_tag = ?", tag).Find(id).Error; err != nil {
			log.GetLogger().Error("[Repository] SegmentsIdNext Get", zap.String("tag", tag), zap.Error(err))
			return err
		}
		// return nil 提交事务
		return nil
	})
	return
}
