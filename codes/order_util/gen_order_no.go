package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println(GenOrderNo(12344566256))
}

func GenOrderNo(userId int64) string {
	const TimeFormatPlainDate = "20060102"
	day := time.Now().Format(TimeFormatPlainDate)

	rand.Seed(time.Now().UnixNano())
	seqStr := fmt.Sprintf("%014d", rand.Intn(99999999999999))

	subId := fmt.Sprintf("%04d", 12344566)
	if len(subId) > 4 {
		subId = subId[len(subId)-5 : len(subId)-1]
	}
	return day + seqStr + subId
}
