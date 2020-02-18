package utils

import (
    "fmt"
    "time"
)

func log(level string, value []interface{}) {
    now := time.Now()
    dateString := fmt.Sprintf(
        "%d-%d-%d %d:%d:%d.%d",
        now.Year(),
        now.Month(),
        now.Day(),
        now.Hour(),
        now.Minute(),
        now.Second(),
        now.Nanosecond(),
    )
    var info string
    for _, i := range value {
        info += fmt.Sprintf("%s", i)
    }
    fmt.Println(dateString, level, info)
}


func LogDebug(value ...interface{}) {
    log("-DEBUG-", value)
}


func LogInfo(value ...interface{}) {
    log("-INFO-", value)
}


func LogWarning(value ...interface{}) {
    log("-WARNING-", value)
}


func LogError(value ...interface{}) {
    log("-ERROR-", value)
}


func LogFatal(value ...interface{}) {
    log("-FATAL-", value)
}
