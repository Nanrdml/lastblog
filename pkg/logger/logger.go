package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"
)

type Level int8

type Fields map[string]interface{}

//定义日志输出等级
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

//返回方法接受者的日志等级,在不同的使用场景中记录不同级别的日志
func (l Level)String()string{
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "Info"
	case LevelWarn:
		return "Warn"
	case LevelError:
		return "Error"
	case LevelFatal:
		return "Fatal"
	case LevelPanic:
		return "Panic"
	}
	return ""
}

//我们完成了日志的分级方法后，
//开始编写具体的方法去进行日志的实例初始化和标准化参数绑定，继续写入如下代码：

type Logger struct {
	newLogger *log.Logger
	ctx       context.Context
	fields    Fields
	callers   []string
}

// NewLogger 底下这个其实可以依靠依赖注入，这样写耦合度高
func NewLogger(w io.Writer, prefix string, flag int) *Logger {
	//New创建一个新的Logger。out变量设置将写入日志数据的目的地。该前缀出现在生成的每一行日志的开头，
	//如果提供了Lmsgprefix标志，则出现在日志头之后。flag参数定义日志记录属性。
	l := log.New(w, prefix, flag)
	return &Logger{newLogger: l}
}
//因为是复制品，所以他们除了地址以外完全一样，这样调用后面的方法的时侯给”和原来一样的“增加方法
//得到了一个有了新功能，其他和原来一样的复制品
//这下面的设置方法都是为了设置logger的属性

//创造一个内容和现在logger一样的复制品
func (l *Logger) clone() *Logger {
	nl := *l
	return &nl
}

// WithFields 设置日志公共字段
func (l *Logger) WithFields(f Fields) *Logger {
	ll := l.clone()
	if ll.fields == nil {
		ll.fields = make(Fields)
	}
	for k, v := range f {
		ll.fields[k] = v
	}
	return ll
}


//WithContext 设置日志上下文属性
func (l *Logger) WithContext(ctx context.Context) *Logger {
	ll := l.clone()
	ll.ctx = ctx
	return ll
}

//WithCaller 设置当前某一层调用栈的信息（程序计数器、文件信息、行号）
func (l *Logger) WithCaller(skip int) *Logger {
	ll := l.clone()

	//调用者报告关于调用goroutine堆栈上函数调用的文件和行号信息。参数skip是要递增的堆栈帧数，0表示caller的调用者。
	//(由于历史原因，跳过的含义在Caller和Caller之间是不同的。)
	//返回值报告相应调用的文件中的程序计数器、文件名和行号。如果无法恢复信息，则boolean ok为false。
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		f := runtime.FuncForPC(pc)
		ll.callers = []string{fmt.Sprintf("%s: %d %s", file, line, f.Name())}
	}

	return ll
}

//WithCallersFrames 设置当前的整个调用栈信息
func (l *Logger) WithCallersFrames() *Logger {
	maxCallerDepth := 25
	minCallerDepth := 1
	callers := []string{}
	pcs := make([]uintptr, maxCallerDepth)
	depth := runtime.Callers(minCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		callers = append(callers, fmt.Sprintf("%s: %d %s", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}
	ll := l.clone()
	ll.callers = callers
	return ll
}


//日志内容格式化和日志内容输出相关方法

// JSONFormat 他是output的内部方法
func (l *Logger) JSONFormat(level Level, message string) map[string]interface{} {
	//加4确保容量一定够
	data := make(Fields, len(l.fields)+4)
	//下面四行添加的是基本属性
	data["level"] = level.String()
	data["time"] = time.Now().Local().UnixNano()
	data["message"] = message
	data["callers"] = l.callers
	//如果这个接受者的Fields字段里本来就有，那么加到data里
	if len(l.fields) > 0 {
		for k, v := range l.fields {
			if _, ok := data[k]; !ok {
				data[k] = v
			}
		}
	}
	return data
}

// Output
//Level参数是选择日志输出等级，
func (l *Logger) Output(level Level, message string) {
	//序列化data
	ll := l.WithCaller(2)
	body, _ := json.Marshal(ll.JSONFormat(level, message))
	content := string(body)
	//选择日志输出等级
	switch level {
	case LevelDebug:
		ll.newLogger.Print(content)
	case LevelInfo:
		ll.newLogger.Print(content)
	case LevelWarn:
		ll.newLogger.Print(content)
	case LevelError:
		ll.newLogger.Print(content)
	case LevelFatal:
		ll.newLogger.Fatal(content)
	case LevelPanic:
		ll.newLogger.Panic(content)
	}
}




//Sprint使用其操作数的默认格式进行格式化，
//并返回结果字符串。当两个操作数都不是字符串时，将在操作数之间添加空格。


//根据先前定义的日志分级，编写对应的日志输出的外部方法

func (l *Logger) Info(v ...interface{}) {
	l.Output(LevelInfo, fmt.Sprint(v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Output(LevelInfo, fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(v ...interface{}) {
	l.Output(LevelFatal, fmt.Sprint(v...))
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Output(LevelFatal, fmt.Sprintf(format, v...))
}

func (l *Logger) Debug(v ...interface{}){
	l.Output(LevelDebug,fmt.Sprint(v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Output(LevelDebug, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(v ...interface{}){
	l.Output(LevelWarn,fmt.Sprint(v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Output(LevelWarn, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(v ...interface{}){
	l.Output(LevelError,fmt.Sprint(v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Output(LevelError, fmt.Sprintf(format, v...))
}

func (l *Logger) Panic(v ...interface{}){
	l.Output(LevelPanic,fmt.Sprint(v...))
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Output(LevelPanic, fmt.Sprintf(format, v...))
}





















