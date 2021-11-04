package selectx

import (
	"reflect"
)

type KindError struct{
	message string
}
func (k KindError) Error() string {
    return k.message
}

func (k KindError) Set(str string) {
	k.message = str
}

func newKindErr(str string) KindError {
	var ret KindError
	ret.Set(str)
	return ret
}

// GetInt 返回int 类型的切片集合
func GetInt(value interface{}, fieldName string) ([]int64, error) {
	res, err := SelectField(value, fieldName)
	if err != nil {
		return nil, err
	}
	var ret []int64
	for _, v := range res {
		tmp, _ := v.(int64)
		ret = append(ret, tmp)
	}
	return ret, err
}

// GetString 返回string 集合
func GetString(value interface{}, fieldName string) ([]string, error) {
	res, err := SelectField(value, fieldName)
	if err != nil {
		return nil, err
	}

	var ret []string
	for _, v := range res {
		tmp, _ := v.(string)
		ret = append(ret, tmp)
	}

	return ret, err
}



// SelectField 包的功能是获取 []struct 中的某个字段组成切片返回
func SelectField(value interface{}, fieldName string) ([]interface{}, error) {
	// 先检查下传入的参数是不是切片
	refV := reflect.ValueOf(value)
	if refV.Kind() != reflect.Slice {
        return nil, newKindErr("传入的类型必需是切片类型")
	}
	elemLen := refV.Len()
	if elemLen < 1 {
		return []interface{}{}, nil
	}
	var ret []interface{}
	for i:=0; i< elemLen; i++ {
		var store interface{}
		val := refV.Index(i)
		switch val.Kind() {
		case reflect.Struct:
			store = getStructFieldValue(val, fieldName)
		case reflect.Ptr:
			if val.Elem().Kind() == reflect.Struct {
				store = getStructFieldValue(val.Elem(), fieldName)
			} else {
				return nil, newKindErr("目前只支持 struct 和 *struct")
			}
		default:
			return nil, newKindErr("目前只支持 struct 和 *struct")
		}
		ret = append(ret, store)
	}

	return ret, nil

}


func getStructFieldValue(value reflect.Value, fieldName string) interface{} {
	val := value.FieldByName(fieldName)
	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64:
		return val.Int()
	case reflect.Uint,reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint()
	}

	return nil

}
