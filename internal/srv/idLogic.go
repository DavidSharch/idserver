package srv

import (
	"context"
	"errors"
	"github.com/sharch/idserver/config"
	"github.com/sharch/idserver/internal/entity"
	"sync"
)

type IdTagMap struct {
	Mu     sync.RWMutex
	TagMap map[string]*Tag
}

// Tag 记录tag对应的ID分配情况
type Tag struct {
	Mu      sync.RWMutex
	TagName string
	// IdArray 记录当前使用的id，以及提前申请的下一个id段
	IdArray []*ID
}

// ID 左开右闭
type ID struct {
	Cur   int64
	Start int64
	End   int64
}

func (s *Service) NewIdTagMap() (*IdTagMap, error) {
	// 读取最近的id 写入map
	var err error
	var res []entity.Segments
	if res, err = s.r.SegmentsGetAll(); err != nil {
		return nil, err
	}
	tagMap := make(map[string]*Tag)
	for _, seg := range res {
		seg := seg
		tagMap[seg.BizTag] = &Tag{
			TagName: seg.BizTag,
			IdArray: make([]*ID, 0),
		}
		tagMap[seg.BizTag].IdArray = append(
			tagMap[seg.BizTag].IdArray, &ID{Start: seg.MaxId, End: seg.MaxId + seg.Step},
		)
	}
	return &IdTagMap{TagMap: tagMap}, nil
}

// GetId 获取指定tag的下一个新id
func (s *Service) GetId(tag string) (id int64, err error) {
	mu := s.idTagMap.Mu
	t, ok := s.idTagMap.TagMap[tag]
	if !ok && config.Conf.Biz.CreatWhenNotExists == 1 {
		// tag不存在，而且配置了不存在就新建
		newSeg := entity.Segments{
			BizTag: tag,
			MaxId:  1,
			Step:   config.Conf.Biz.DefaultStep,
		}
		mu.Lock()
		if err = s.CreateTag(&newSeg); err != nil {
			mu.Unlock()
			return 0, err
		}
		_ = s.CreateTag(&newSeg)
		// t重新赋值，此时一定有数据
		t, _ = s.idTagMap.TagMap[tag]
		mu.Unlock()
	} else if !ok && config.Conf.Biz.CreatWhenNotExists == 0 {
		// 不存在，同时没有配置不存在就新建
		return 0, errors.New("tag does not exist")
	}

	t.Mu.Lock()
	id = t.nextId()
	t.Mu.Unlock()
	return
	// todo 这里可以考虑：返回错误，还是等到拿到新id段， 返回新id
}

// RemainsId 查询tag的id是否有剩余
func (t *Tag) remainsId() bool {
	if len(t.IdArray) > 1 {
		// 超过一个号段，说明已经提前申请下一个号段了，一定存在id可用
		return true
	} else {
		// cur < end
		return t.IdArray[0].Cur < t.IdArray[0].End
	}
}

// RemainsIdCnt 查询tag的id剩余个数
func (t *Tag) remainsIdCnt() int64 {
	var cnt int64
	for _, id := range t.IdArray {
		id := id
		cnt += id.End - id.Cur + 1
	}
	return cnt
}

// NextId 获取下一个id
func (t *Tag) nextId() (id int64) {
	id = t.IdArray[0].Start + t.IdArray[0].Cur
	t.IdArray[0].Cur++
	// 移除已经使用完毕的id段
	if id+1 >= t.IdArray[0].End {
		t.IdArray = append(t.IdArray[:0], t.IdArray[1:]...)
	}
	return
}

// getNextIdStep 获取下一个号段
func (t *Tag) getNextIdStep(cancel context.CancelFunc, s *Service) {
	defer cancel()
	// todo 加入重试，已经各种告警手段
	t.Mu.RLock()
	// 已经有号段储备，就不处理
	if len(t.IdArray) < 2 {
		t.Mu.RUnlock()
		// todo service去数据库读取
		nextStep := &entity.Segments{}
		t.Mu.Lock()
		t.IdArray = append(t.IdArray, &ID{Start: nextStep.MaxId, End: nextStep.MaxId + nextStep.Step})
		t.Mu.Unlock()
	} else {
		t.Mu.RUnlock()
	}
}

// CreateTag 把seg转换成tag，然后写入map
func (s *Service) CreateTag(e *entity.Segments) error {
	data, err := s.r.SegmentsCreate(e)
	if err != nil {
		return err
	}
	b := &Tag{
		TagName: e.BizTag,
		IdArray: make([]*ID, 0),
	}
	b.IdArray = append(b.IdArray, &ID{Start: data.MaxId, End: data.MaxId + data.Step})
	s.idTagMap.TagMap[e.BizTag] = b
	return nil
}
