package main

import (
  "test_defs"
  "logger"
  "csv_tests"
  "db_tests"
)

var tests []test_defs.TestSuite = []test_defs.TestSuite {
  csv_tests.TEST_SUITE,
  db_tests.TEST_SUITE,
}

func main() {
  logger.Info("Running tests...")
  for _, ts := range tests {
    ts.Run()
  }
  logger.Info("Finished running tests...")
}

