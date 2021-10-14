package verify

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
)

// 包的主要作用是校验客户端传进来的参数
// 1. 首先我们先建一个类型 保存要怎样校验数据
type Rules map[string][]string

// ======定义怎样检查参数
// 检查参数是否为空
func CheckEmpty() string {
	return "checkEmpty"
}

// 正则匹配
func RegExp(exp string) string {
	return "regexp=" + exp
}

func Lt(c string) string {
	return "lt=" + c
}

func Le(c string) string {
	return "le=" + c
}

func Eq(c string) string {
	return "eq=" + c
}

func Gt(c string) string {
	return "gt=" + c
}

func Ge(c string) string {
	return "ge=" + c
}

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
				switch {
				case v == "checkEmpty":
					if IsBlank(fV) {
						return errors.New(fieldN.Name + " is empty")
					}
				case strings.Split(v, "=")[0] == "regexp":
					if Exp(strings.Split(v, "=")[1], fV.String()) == false {
						return errors.New(fieldN.Name + "regexp match falied")
					}
				}
			}
		}

	}
}

// 正则匹配
func Exp(rule, c string) bool {
	return regexp.MustCompile(rule).MatchString(c)
}

// 判断值是否是空白值
// 现在主要判断是否是 bool "" 0 ptr==nil intface==nil
func IsBlank(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.String:
		return val.String() == ""
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0
	case reflect.Bool:
		return !val.Bool()
	case reflect.Ptr, reflect.Interface:
		return val.IsNil()
	}

	return false
}
