package verify

import (
	"errors"
	"reflect"
)

// 包的主要作用是校验客户端传进来的参数
// 1. 首先我们先建一个类型 保存要怎样校验数据
type Rules map[string][]string

// ======定义怎样检查参数
// 检查参数是否为空
func CheckEmpty() string {
	return "CheckEmpty"
}

// 正则匹配
func RegExp(exp string) string {
	return "RegExp="+exp
}

func 

// 1. 首先我们传进来一个结构体和一个规则(Rules)
// 2. 获取对应结构体的字段 去看规则里有没有，如果有则根据规则校验
// 3. 如果校验不通过，则
func Verify(st interface{}, ruleMap Rules) error {
	typ := reflect.TypeOf(st)
	val := reflect.ValueOf(st)
	if val.Kind() != reflect.Struct {
		return errors.New("expect struct")
	}

	// 获取结构体字段的数量
	numF := typ.NumField()
	for i := 0; i < numF; i++ {
		// 获取字段名和字段值
		fieldN := typ.Field(i)
		fV := val.Field(i)

		if len(ruleMap[fieldN.Name]) > 0 { // 说明有规则要处理
			for _, v := range ruleMap[fieldN.Name] {

			}
		}

	}
}
