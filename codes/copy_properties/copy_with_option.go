package main

import (
	"errors"
	"github.com/jinzhu/copier"
	"regexp"
	"time"
)

func CopyProperties(dst, src interface{}) error {
	err := copier.CopyWithOption(dst, src, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			{ // time.Time 转换成字符串
				SrcType: time.Time{},
				DstType: copier.String,
				Fn: func(src interface{}) (dst interface{}, err error) {
					s, ok := src.(time.Time)
					if !ok {
						return nil, errors.New("src type is not time.Time")
					}
					return s.Format("2006-01-02 15:04:05"), nil
				},
			},
			{ // 字符串转成time.Time
				SrcType: copier.String,
				DstType: time.Time{},
				Fn: func(src interface{}) (dst interface{}, err error) {
					s, ok := src.(string)

					if !ok {
						return nil, errors.New("src type is not time format string")
					}
					pattern := `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$` // YYYY-MM-DD HH:MM:SS
					matched, _ := regexp.MatchString(pattern, s)
					if matched {
						return time.Parse("2006-01-02 15:04:05", s)
					}
					return nil, errors.New("src type is not time format string")
				},
			},
		},
		
	})

	return err
}
