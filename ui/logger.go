package logger

import (
  "fmt"
  "os"
  "runtime"
  "strings"
)

func escapeCode(args string) string {
  return fmt.Sprintf("\033[%s", args)
}

func changeColor(color int8) string {
  return escapeCode(fmt.Sprintf("%dm", color))
}

func resetColor() string {
  return changeColor(0)
}

func TruncateCallingFile(file string) string {
  mudIndex := strings.Index(file, "mudServer")
  if mudIndex >= 0 {
    file = file[mudIndex + 10:]
  }

  return file
}

func LogFileLine(color int8, file string, line int, prefix string, args ...interface{}) {
  file = TruncateCallingFile(file)

  fmt.Fprintf(os.Stderr, "%s%s:%d: %s", changeColor(97), file, line, changeColor(color))
  fmt.Fprintf(os.Stderr, "[%s]: %s", prefix, resetColor())
  if len(args) == 1 {
    fmt.Fprint(os.Stderr, args[0])
  } else {
    fmt.Fprintf(os.Stderr, args[0].(string), args[1:]...)
  }
  fmt.Fprint(os.Stderr, "\n")
}

func LogCustomCaller(color int8, prefix string, caller int, args ...interface{}) {
  _, file, line, _ := runtime.Caller(caller + 1)
  LogFileLine(color, file, line, prefix, args...)
}

func Log(color int8, prefix string, args ...interface{}) {
  LogCustomCaller(color, prefix, 2, args...)
}

func Info(args ...interface{}) {
  Log(96, "INFO", args...)
}

func Warn(args ...interface{}) {
  Log(93, "WARN", args...)
}

func Error(args ...interface{}) {
  Log(91, "ERROR", args...)
}

func ErrorFileLine(file string, line int, args ...interface{}) {
  LogFileLine(91, file, line, "ERROR", args...)
}

func ErrorCustomCaller(caller int, args ...interface{}) {
  LogCustomCaller(91, "ERROR", caller + 1, args...)
}

