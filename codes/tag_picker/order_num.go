package tag_picker


type OrderNumTagPicker struct {
	TagName string
	ValueCh chan interface{}
}

func (picker *OrderNumTagPicker) PickTagValueForUser(userId int64, args ...interface{}) {
	// 用类型转换得到交易号
	tradeNo, ok := args[0].(string)
	if !ok {
                // log.Error自己实现 
		log.Error("PayTotalTagPickerError", "Invalid arg", args[0])
		// 结束执行并通知外部
		picker.ValueCh <- nil
		return
	}
        // 这里就打印下参数值，标签查询的具体逻辑自己实现
        fmt.Println(userId)
        fmt.Println(TradeNo)
        ......
	// 查询到的用户的标签
	picker.ValueCh <- 10
}

func (picker *OrderNumTagPicker) Notify () <-chan interface{} {
	return picker.ValueCh
}
