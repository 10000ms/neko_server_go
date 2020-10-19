package utils

import (
    "fmt"
    "time"
)

func log(level string, value []interface{}, color string) {
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
    str := fmt.Sprintf("%s %s %s", dateString, level, info)
    if color != "" {
        str = fmt.Sprintf("\u001B[1;0;%sm%s\u001B[0m\n", color, str)
    }
    fmt.Println(str)
}

func LogDebug(value ...interface{}) {
    log("[DEBUG]", value, "")
}

func LogInfo(value ...interface{}) {
    log("[INFO]", value, "")
}

func LogWarning(value ...interface{}) {
    log("[WARNING]", value, "33")
}

func LogError(value ...interface{}) {
    log("[ERROR]", value, "31")
}

func LogFatal(value ...interface{}) {
    log("[FATAL]", value, "31")
}

func LogSystem(value ...interface{}) {
    var info string
    for _, i := range value {
        info += fmt.Sprintf("%s", i)
    }
    fmt.Printf("\033[1;0;36m%s\033[0m\n", info)
}
