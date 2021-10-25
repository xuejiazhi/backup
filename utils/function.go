package utils

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

// 获取当前协程ID
func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func InArray(need interface{}, needArr []string) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

var timeLayout = "2006-01-02 15:04:05"

func GetFormatDateTime() string {
	timeUnix := time.Now().UnixNano() / 1e6 //已知的时间戳
	miroTime := timeUnix - timeUnix/1000*1000
	formatTimeStr := time.Unix(timeUnix/1000, 0).Format(timeLayout)
	//fmt.Println(formatTimeStr) //打印结果：2017-04-11 13:30:39
	return formatTimeStr + "." + ConvertInterfaceToString(miroTime)
}

func GetDate() string {
	t1 := time.Now().Year()  //年
	t2 := time.Now().Month() //月
	t3 := time.Now().Day()   //日
	return fmt.Sprintf("%d-%d-%d", t1, t2, t3)
}

func ConvertInterfaceToString(value interface{}) (result string) {
	typeOfA := reflect.TypeOf(value)
	switch typeOfA.Kind().String() {
	case "int":
		return strconv.Itoa(value.(int))
	case "int64":
		return strconv.FormatInt(value.(int64), 10)
	case "string":
		return value.(string)
	case "float32":
		return strconv.FormatFloat(value.(float64), 'G', -1, 32)
	case "float64":
		return strconv.FormatFloat(value.(float64), 'G', -1, 64)
	default:
		return ""
	}
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
