package main

import (
    "fmt"
    "github.com/wanmei002/goutil/mr"
    "os"
    "runtime/trace"
)

func main(){
    trace.Start(os.Stderr)
    defer trace.Stop()
    
    // 传入一个管道; 主要用于把要执行的数据放入管道里
    //GenerateFunc func(source <-chan interface{})
    //// item 是要处理的数据, Writer 要把结果集写入管道中, cancel 如果执行有错误, 可以调用它,停止后续执行
    //MapperFunc   func(item interface{}, writer Writer, cancel func(err error))
    //// pipe 是用来保存处理传入数据的结果的, writer 用于把合并的结果集合并, cancel 如果出现错误了请调用它
    //ReducerFunc  func(pipe <-chan interface{}, writer Writer, cancel func(err error))
    uid := []int{1,2,3,4,5,6}
    a := func(source chan<- interface{}){
        for _,v := range uid {
            source <- v
        }
    }
    
    b := func(item interface{}, writer mr.Writer, cancel func(err error)){
        fmt.Println("b: item:", item)
        tmp := item.(int) + 1
        fmt.Println("b:", tmp)
        writer.Writer(tmp)
    }
    
    c := func(pipe <-chan interface{}, writer mr.Writer, cancel func(err error)){
        fmt.Println("c:", writer)
        var uid []int
        for v := range pipe {
            uid = append(uid, v.(int))
        }
        fmt.Println(uid)
        writer.Writer(uid)
    }
    
    res, err := mr.MapReduce(a, b, c)
    
    fmt.Printf("res:{%v}; err:{%v}\n", res, err)
}
