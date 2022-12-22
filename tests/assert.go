package assert

import (
  "logger"
  "runtime"
  "fmt"
)

func Assert(condition bool, args ...interface{}) {
  internalAssert(condition, 1, args...)
}

func internalAssert(condition bool, caller int, args ...interface{}) {
  if !condition {
    _, file, line, _ := runtime.Caller(caller)
    logger.ErrorFileLine(file, line, args...)

    var panicString string
    if len(args) == 1 {
      panicString = args[0].(string)
    } else if len(args) > 1 {
      panicString = fmt.Sprintf(args[0].(string), args[1:]...)
    } else {
      panicString = "Test Failed"
    }

    panic(panicString)
  }
}

func ToEqual(actual interface{}, expected interface{}, message string) {
  message = fmt.Sprintf("%s: (expected %v, found %v)", message, expected, actual)
  internalAssert(actual == expected, 2, message)
}
