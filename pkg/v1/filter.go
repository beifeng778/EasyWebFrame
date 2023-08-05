package v1

import (
	"fmt"
	"time"
)

type handlerFuc func(c *Context)

type Filter func(c *Context)

type FilterBuilder func(next Filter) Filter

var _ FilterBuilder = MetricFilter

func MetricFilter(next Filter) Filter {
	return func(c *Context) {
		start := time.Now().Nanosecond()
		next(c)
		end := time.Now().Nanosecond()
		fmt.Printf("执行时间为: %d 纳秒\n", end-start)
	}
}
