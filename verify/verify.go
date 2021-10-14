package verify

import (
	"errors"
	"log"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
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
					// 剩下的就是大于等于比较了
				case strings.Index(v, "=") >= 0:
					CompareVerify(fV, v)
				}
			}
		}

	}

	return nil
}

// 比较下
func CompareVerify(val reflect.Value, exp string) bool {
	tp := val.Kind()
	// 根据不同的类型进行不同的比较
	tmp := strings.Split(exp, "=")
	if len(tmp) < 2 {
		pc, file, line, _ := runtime.Caller(0)
		fnN := runtime.FuncForPC(pc)
		log.Printf("file[%v] line[%v]; func name[%v]\n", file, line, fnN.Name())
		return false
	}

	cp := tmp[0]
	cV := tmp[1]
	intV, _ := strconv.Atoi(cV)
	switch tp {
	case reflect.Slice, reflect.Map, reflect.Array, reflect.String:
		return compare(val.Len(), cp, intV)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compare(int(val.Int()), cp, intV)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return compare(int(val.Uint()), cp, intV)
	}

	return true
}

// 对数据进行长度比较
func compare(len int, tp string, val int) bool {
	switch tp {
	case "lt":
		return len < val
	case "le":
		return len <= val
	case "eq":
		return len == val
	case "gt":
		return len > val
	case "ge":
		return len >= val
	}

	return false
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

	// 零值比较
	return reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())
}
