package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	sep             = string(os.PathSeparator)
	skipForThisFunc = 3
)

// callerLine 获取调用者的代码行号
func callerLine(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	filePath, _ := filepath.Abs(file)
	filePathSplit := strings.Split(filePath, sep)
	pathLength := 4 // 路径层数 默认四层
	if len(filePathSplit) < pathLength {
		pathLength = len(filePathSplit)
	}
	return fmt.Sprintf("%s:%d", filepath.Join(filePathSplit[len(filePathSplit)-pathLength:]...), line)
}

// funcName current func name, use 'F' for short, if you want to change the skip, input one
func funcName(skip ...int) string {
	useSkip := 1
	if len(skip) > 0 {
		useSkip = skip[0]
	}
	pc, _, _, ok := runtime.Caller(useSkip)
	if !ok {
		return ""
	}
	dotVec := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	return dotVec[len(dotVec)-1]
}
