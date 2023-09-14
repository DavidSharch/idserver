package srv

import (
	"context"
	"errors"
	"fmt"
	"github.com/sharch/idserver/internal/entity"
	"sync"
	"time"
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
func (t *Tag) GetId(s *Service) (id int64, err error) {
	canGetId := false
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)

	t.Mu.Lock()

	if t.remainsId() == true {
		id = t.nextId()
		canGetId = true
	}

	// 判断是否需要准备下一个id段 使用超过一半
	if t.IdArray[0].End/t.IdArray[0].Cur < 2 {
		t.Mu.Unlock()
		go t.getNextIdStep(cancel, s)
	} else {
		// 不需要申请新id段，就解锁，defer调用cancel
		t.Mu.Unlock()
		defer cancel()
	}
	// 已经拿到id，就返回
	if canGetId {
		return
	}

	// 等待获取到下一段
	select {
	case <-ctx.Done():
		fmt.Printf("等待获取到下一段，已经到时\n")
	}
	t.Mu.Lock()
	if t.remainsId() {
		id = t.nextId()
		canGetId = true
	} else {
		err = errors.New("get id failed")
	}
	t.Mu.Unlock()

	return
}

// RemainsId 查询tag的id是否有剩余
func (t *Tag) remainsId() bool {
	t.Mu.RLock()
	defer t.Mu.RUnlock()
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
	t.Mu.RLock()
	defer t.Mu.RUnlock()
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
