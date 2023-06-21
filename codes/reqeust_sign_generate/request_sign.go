package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
)

func GenRequestSign(request map[string]interface{}) (string, string) {

	delete(request, "sign")

	var keys []string
	for k := range request {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var reString bytes.Buffer
	for _, kk := range keys {
		if kk == "" {
			continue
		}
		if reString.Len() > 0 {
			reString.WriteByte('&')
		}

		reString.WriteString(kk)
		reString.WriteByte('=')
		switch vv := request[kk].(type) {
		case string:
			reString.WriteString(vv)
		case int, int8, int16, int32, int64:
			reString.WriteString(fmt.Sprintf("%d", vv))
		case float64:
			reString.WriteString(strconv.FormatInt(int64(vv), 10))
		default:
			continue
		}
	}
	secret := "37y4uxXZXeWtCDRq3z14dEhUhCawb2tM"
	reString.WriteString(secret)

	return reString.String(), Md5(reString.String())
}

func Md5(str string) string {
	md5ctx := md5.New()
	md5ctx.Write([]byte(str))
	return hex.EncodeToString(md5ctx.Sum(nil))
}

func main()  {
	req := map[string]interface{} {
		"request_no":"230621175012994051431130267648serial",
		"company_sign":"aaa",
	}

	_, b := GenRequestSign(req)
	fmt.Println(b)

}
