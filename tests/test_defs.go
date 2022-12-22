package test_defs

import (
  "logger"
  "runtime"
  "reflect"
)

type TestSuite struct {
  Name string
  Tests []func()
}

func catchFailingTest(test int) {
  if r := recover(); r != nil {
    logger.Error("Test %d Failed: %s", test, r)
  }
}

func (ts TestSuite) Run() {
  logger.Info("Running tests for %s", ts.Name)
  for ti, test := range ts.Tests {
    p := reflect.ValueOf(test).Pointer()
    logger.Info("Running test %d: %v", ti, runtime.FuncForPC(p).Name())
    func() {
      defer catchFailingTest(ti)
      test()
    }()
  }
}

