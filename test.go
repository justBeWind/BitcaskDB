package main

import (
	"sync/atomic"
)

type Foo struct {
	a int64
	b int32
	c int64
}

func main() {
	var f Foo
	//atomic.AddInt64(&f.a, 1) // 这里不会崩溃
	atomic.AddInt64(&f.c, 1) // 这里会崩溃
}
