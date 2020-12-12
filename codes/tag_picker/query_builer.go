package tag_picker

import (
	"context"
	"sync"
	"time"
)

const (
  TAG_ORDER_NUM = "order_num"
)

// 对外提供的批量查询用户标签值的方法
func BulkQueryUserTagValue(tagNames []string, userId int64, queryArgs ...interface{}) (tagValuePairs []*TagValuePair) {
	tagCount := len(tagNames)
	if tagCount < 1 {
		return
	}
	wg := &sync.WaitGroup{}
	wg.Add(tagCount)
	tagValueCh := make(chan *TagValuePair, tagCount) // 用于接收所有Picker查到的标签值的Channel
	ctx, _ := context.WithTimeout(context.Background(), time.Minute) // 设置执行标签值查找的超时时间
	for _, tagName := range tagNames {
		go queryTagValue(ctx, wg, tagName, userId, tagValueCh, queryArgs...)
	}
	wg.Wait()
	close(tagValueCh) // 先关闭通道 方便下面for range不发生阻塞, 从channel中读完值即退出
	tagValuePairs = make([]*TagValuePair, 0)
	for tagValue := range tagValueCh {
		if tagValue.Value != nil {
			tagValuePairs = append(tagValuePairs, tagValue)
		}
	}

	return tagValuePairs
}

type TagValuePair struct {
	Name string `json:"tag_name"`
	Value interface{} `json:"tag_value"`
}


func queryTagValue(ctx context.Context, wg *sync.WaitGroup, tagName string, userId int64, tagValueCh chan *TagValuePair, queryArgs ...interface{}) {
	defer wg.Done()
	tagPicker := resolveTagPicker(tagName)
	if tagPicker == nil {
		log.Error("未识别的业务标签", common.ErrUnknownBusinessTag)
		return
	}
	go tagPicker.PickTagValueForUser(userId, queryArgs...)
	select {
	case <- ctx.Done(): // 超时返回
		return
	case tagValue := <- tagPicker.Notify(): // 接收标签值
		TagValuePair := &TagValuePair{
			Name:  tagName,
			Value: tagValue,
		}
		tagValueCh <- TagValuePair

		return
	}
}

// 根据标签名解析出对应的TagPicker
func resolveTagPicker (tagName string) Picker {
	switch tagName {
	case TAG_ORDER_NUM:
		return &OrderNumTagPicker{
			TagName: tagName,
			ValueCh: make(chan interface{}),
		}
	default:
		return nil
	}
}
