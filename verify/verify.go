package verify

// 包的主要作用是校验客户端传进来的参数
// 1. 首先我们先建一个类型 保存要怎样校验数据
type Rules map[string][]string

// 1. 首先我们传进来一个结构体和一个规则(Rules)
// 2. 获取对应结构体的字段 去看规则里有没有，如果有则根据规则校验