package main

import (
    "errors"
    "fmt"
    
    //"github.com/wanmei002/goutil/threading"
    "sync"
    "sync/atomic"
)

// 此包主要用来并发执行方法

// 第一个方法用来往管道里存储要执行的事情
// 第二个方法用来执行方法
// 第三个方法用来归并结果集

var (
    // 英语不好, 请各位看官谅解
    cancelWithNil = errors.New("reduce end with nil error")
)

type (
    // 传入一个管道; 主要用于把要执行的数据放入管道里
    GenerateFunc func(source chan<- interface{})
    // item 是要处理的数据, Writer 要把结果集写入管道中, cancel 如果执行有错误, 可以调用它,停止后续执行
    MapperFunc   func(item interface{}, writer Writer, cancel func(err error))
    // pipe 是用来保存处理传入数据的结果的, writer 用于把合并的结果集合并, cancel 如果出现错误了请调用它
    ReducerFunc  func(pipe <-chan interface{}, writer Writer, cancel func(err error))
    
    Writer interface {
        Writer(val interface{})
    }
)

func newWriteChan(write chan interface{}, done chan struct{}) writeChan {
    return writeChan{
        write: write,
        done:  done,
    }
}

type writeChan struct {
    write chan interface{}
    done chan struct{}
}

func (w writeChan) Writer(val interface{}) {
    select {
    case <-w.done:
        return
    default:
        w.write <- val
    }
    
}

func (w writeChan) Load() interface{} {
    return <-w.write
}


// 用于把要处理的数据传进chan里
func buildSource(generate GenerateFunc) chan interface{} {
    source := make(chan interface{})
    go func(){
        // 在这里关闭管道
        defer func(){
            close(source)
        }()
        generate(source)
    }()
    
    return source
}

func MapReduce(generate GenerateFunc, mapper MapperFunc, reducer ReducerFunc) (interface{}, error) {
    source := buildSource(generate)
    // 启动一个协程用于保存错误
    // 在启动一个协程用于保存执行结果
    // 在启动一个协程用于保存处理结果集
    errChan := make(chan interface{})
    // 原子存储错误, 避免资源竞争
    var errVal atomic.Value
    // 执行一次关闭
    var cancelOnce sync.Once
    cancel := func(err error){
        if err != nil {
            errVal.Store(err)
        } else {
            errVal.Store(cancelWithNil)
        }
    
        cancelOnce.Do(func(){
            close(errChan)
            
        })
    }
    // 创建写入实例 reducer 这个变量把结果汇总给这个chan
    reduceChan := make(chan interface{})
    done := make(chan struct{})
    write := newWriteChan(reduceChan, done)
    // 存放处理结果, mapper 往这个里面写入， reduce 从这个里面读出
    resChan := make(chan interface{})
    // 启动协程 执行结果集归并方法
    go func(){
        reducer(resChan, write, cancel)
    }()
    
    
   // 开始处理  现在一个方法里写 最后再整理下
   // 现在开始从执行管道里读取数据处理
   // 上面那个 write 是用来把结果合并的。这里再来一个 write 用来把处理出来的结果保存
   execWrite := newWriteChan(resChan, done)
   
   go func (){
       wg := sync.WaitGroup{}
       defer func(){
           wg.Wait()
           fmt.Println("start close reschan")
           close(resChan)
       }()
       // 在这里建一个管道 控制开启的协程数量
       pool := make(chan struct{}, 16)
       for {
           select {
           case <-done:
               return
           // 读取要处理的数据
           case pool <- struct{}{}:
               // 在这里判断管道是否关闭了
               item, ok := <-source
               if !ok {// 说明管道已经关闭了
                   fmt.Println("管道已经关闭")
                   return
               }
               wg.Add(1)
               go func(){
                   defer func(){
                       wg.Done()
                   }()
    
                   // 运行自定义的处理函数
                   mapper(item, execWrite, cancel)
               }()
           }
       }
   }()
   
   
   // 此时我们应该取出错误 和 结果
   res, ok := <-reduceChan
   errIntF := errVal.Load()
   var err error
   if errIntF != nil {
       err = errIntF.(error)
   }
   if !ok {
       return nil, err
   } else {
       return res, nil
   }
}

func main(){
    
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
    
    b := func(item interface{}, writer Writer, cancel func(err error)){
        fmt.Println("b: item:", item)
        tmp := item.(int) + 1
        fmt.Println("b:", tmp)
        writer.Writer(tmp)
    }
    
    c := func(pipe <-chan interface{}, writer Writer, cancel func(err error)){
        fmt.Println("c:", writer)
        var uid []int
        for v := range pipe {
            uid = append(uid, v.(int))
        }
        fmt.Println(uid)
        writer.Writer(uid)
    }
    
    res, err := MapReduce(a, b, c)
    
    fmt.Printf("res:{%v}; err:{%v}\n", res, err)
}
