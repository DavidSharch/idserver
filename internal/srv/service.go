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

func NewService(c *config.Config) (s *Service) {
	var err error
	s = &Service{
		c: c,
		r: orm.NewRepository(c),
	}
	if s.idTagMap, err = NewIdTagMap(); err != nil {
		log.GetLogger().Error("[NewService] NewAllocId ", zap.Error(err))
		panic(err)
	}
	return s
}
