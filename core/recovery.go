package core

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"neko_server_go/utils"
)

func recoveryHandlerPanic(w http.ResponseWriter, r *http.Request) {
	// 处理handler panic。返回500错误 Internal Error:
	if errorRecover := recover(); errorRecover != nil {
		message := fmt.Sprintf("%s", errorRecover)
		stackInfo := tracePanic(message)
		utils.LogError("Internal Error: ", errorRecover, stackInfo)
		InternalErrorHandler(w, r)
	}
}

func tracePanic(message string) string {
	var pcs [32]uintptr
	// Callers 用来返回调用栈的程序计数器,
	// 第 0 个 Caller 是 Callers 本身
	// 第 1 个是上一层 trace
	// 第 2 个是再上一层的 defer func。
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		// 获取调用函数
		fn := runtime.FuncForPC(pc)
		// 获取调用函数文件，行
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
