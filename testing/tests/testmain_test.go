package tests

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

var suiteLogger *log.Logger
var logFile *os.File
func TestMain(m *testing.M) {
	_ = os.MkdirAll("tests/target/logs", 0755)
	_ = os.MkdirAll("tests/target/screenshots", 0755)
	_ = os.MkdirAll("tests/target/reports", 0755)

	var err error
	logFile, err = os.Create(filepath.Join("tests", "target", "logs", "test-run.log"))
	if err != nil {
		suiteLogger = log.New(os.Stdout, "[tests] ", log.LstdFlags)
		code := m.Run()
		os.Exit(code)
	}

	suiteLogger = log.New(logFile, "[tests] ", log.LstdFlags)
	suiteLogger.Println("=== SUITE START ===")

	code := m.Run()

	suiteLogger.Println("=== SUITE END ===")
	_ = logFile.Close()
	os.Exit(code)
}
