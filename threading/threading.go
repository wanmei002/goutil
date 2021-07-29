package threading

import "log"

func SafeGoroutine(fn func()){
    go RunFn(fn)
}

func RunFn(fn func()) {
    defer func(){
        if r := recover(); r != nil {
            log.Println("err:", r)
        }
    }()
    
    fn()
}
