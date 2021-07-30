package mr

import (
    "errors"
    "fmt"
    "github.com/wanmei002/goutil/threading"
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
    // 用来扩展上面的 MapperFunc
    MapFunc func(item interface{}, writer Writer)
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

func drain(channel <-chan interface{}) {
    for  range channel {
    
    }
}


func Finish(fns  ...func()error) error {
    
    _, err := MapReduce(func(source chan<- interface{}){
        for _, fn := range fns {
            source <- fn
        }
    }, func(item interface{}, writer Writer, cancel func(err error)){
        fmt.Println(item)
        fn := item.(func()error)
        if err := fn(); err != nil {
            cancel(err)
        }
    }, func(pipe <-chan interface{}, write Writer, cancel func(err error)){
        drain(pipe)
    })
    
    return err
}


// MapReduce 逻辑简单介绍下:
//1. 要处理的数据放到一个无缓冲的管道里，再来一个协程从这个无缓冲的管道里读取要处理的数据，然后把读取出来的数据开一个协程用 传入的方法处理;
//2. 创建一个无缓冲的管道, 用来保存执行的结果, 让合并结果的协程从这个管道里读取数据, 然后合并数据, 写入合并数据的管道里
//2. 创建一个用来停止其它协程的管道, 如果执行中有什么错误就关闭这个管道里，同时停止执行其它协程，返回失败
//4. 最后要把没有关闭的管道关闭了, 让 gc 回收
func MapReduce(generate GenerateFunc, mapper MapperFunc, reducer ReducerFunc) (interface{}, error) {
    source := buildSource(generate)
    // 启动一个协程用于保存错误
    // 在启动一个协程用于保存执行结果
    // 在启动一个协程用于保存处理结果集
    // 原子存储错误, 避免资源竞争
    var errVal atomic.Value
    // 执行一次关闭
    // 如果有任何异常都把这个关闭了，让其他的接收到通知，停止运行
    done := make(chan struct{})
    // 创建写入实例 reducer 这个变量把结果汇总给这个chan
    reduceChan := make(chan interface{})
    var (
    	cancelOnce sync.Once
    	reduceClose sync.Once
    )
    finish := func(){
        cancelOnce.Do(func(){
            close(done)
        })
        reduceClose.Do(func(){
            close(reduceChan)
        })
    }
    
    cancel := func(err error){
        
        if err != nil {
            errVal.Store(err)
        } else {
            errVal.Store(cancelWithNil)
        }
        defer func(){
            // 把资源管道里的数据清空
            drain(source)
        }()
    
        finish()
    }
    write := newWriteChan(reduceChan, done)
    // 存放处理结果, mapper 往这个里面写入， reduce 从这个里面读出
    resChan := make(chan interface{})
    // 启动协程 执行结果集归并方法
    go func(){
        
        defer func(){
            finish()
            drain(resChan)
        }()
        // 在这里可能遇到错误就结束运行了, reschan 可能还有数据, 所以要在 defer 中把数据都给读取完
        reducer(resChan, write, cancel)
    }()
    
    
   // 现在开始从执行管道里读取数据处理
   go executeMappers(func(item interface{}, writer Writer) {
       mapper(item, writer, cancel)
   }, resChan, done, source)
   
   
   // 此时我们应该取出错误 和 结果
   res, ok := <-reduceChan
   errIntF := errVal.Load()
   var err error
   if errIntF != nil {
       err = errIntF.(error)
       return nil, err
   }
   if !ok {
       return nil, err
   } else {
       return res, nil
   }
}

// executeMappers
func executeMappers(mapper MapFunc,resChan chan interface{}, done chan struct{},source <-chan interface{}){
    wg := sync.WaitGroup{}
    defer func(){
        wg.Wait()
        // 在这里关闭 resChan 管道，不用关闭done管道, 在这里 done可能关闭了, 或者在上一层关闭
        close(resChan)
    }()
    writer := newWriteChan(resChan, done)
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
                <-pool
                return
            }
            wg.Add(1)
            threading.SafeGoroutine(func(){
                defer func(){
                    wg.Done()
                    // 在这里关闭, 以保证最多有 16 个在进行
                    <-pool
                }()
    
                // 运行自定义的处理函数
                mapper(item, writer)
            })
        }
    }
}
