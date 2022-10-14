package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// Recovery 的实现，使用 defer 挂载上错误恢复的函数，在这个函数中调用 *recover()*，
// 捕获 panic，并且将堆栈信息打印在日志中，向用户返回 Internal Server Error。
// trace() 函数用来获取触发 panic 的堆栈信息
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}

// print stack trace for debug
// Recovery 的实现非常简单，使用 defer 挂载上错误恢复的函数，
// 在这个函数中调用 *recover()*，捕获 panic，并且将堆栈信息打印在日志中，
// 向用户返回 Internal Server Error。
// trace() 函数，这个函数是用来获取触发 panic 的堆栈信息，
// 在 trace() 中，调用了 runtime.Callers(3, pcs[:])，Callers 用来返回调用栈的程序计数器,
// 第 0 个 Caller 是 Callers 本身，第 1 个是上一层 trace，第 2 个是再上一层的 defer func。
// 因此，为了日志简洁一点，我们跳过了前 3 个 Caller。
// 接下来，通过 runtime.FuncForPC(pc) 获取对应的函数，
// 在通过 fn.FileLine(pc) 获取到调用该函数的文件名和行号，打印在日志中。
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
