### 本包的作用
从切片中找到某个字段，组合一个新切片，example:
```go
有一个切片: var value []map[string]interface{}
想获得 map 中 key 为 aaa 值
则执行：
GetInterface(value, "aaa")
```

### 支持的数据结构
```go
type User struct {
    Name string
    Age  int
}
var value []User
var value []*User
var value []map[string]interface{}
var value []map[string]int
var value []map[string]string
```