package tag_picker
// TagPicker 接口定义
type Picker interface {
	// 用于查询用户的标签值
	PickTagValueForUser (userId int64, args ... interface{})
	// 通知查询到的标签值
	Notify () <-chan interface{}
}
