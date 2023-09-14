package srv

import (
	"github.com/sharch/idserver/config"
	"github.com/sharch/idserver/internal/log"
	"github.com/sharch/idserver/internal/orm"
	"go.uber.org/zap"
)

type Service struct {
	// config
	c *config.Config
	// db
	r *orm.Repository
	// idTagMap id分配
	idTagMap *IdTagMap
}

var ServiceInstance *Service

// NewService 初始化MySQL，缓存预热
func NewService(c *config.Config) (s *Service) {
	var err error
	s = &Service{
		c: c,
		r: orm.NewRepository(c.Mysql),
	}
	if s.idTagMap, err = s.NewIdTagMap(); err != nil {
		log.GetLogger().Error("service.NewService 初始化失败", zap.Error(err))
		panic(err)
	}
	ServiceInstance = s
	return s
}
