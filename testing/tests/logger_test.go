package tests

import (
	"fmt"
	"testing"
	"time"
)

func logStart(t *testing.T) {
	msg := fmt.Sprintf("START: %s", t.Name())
	if suiteLogger != nil {
		suiteLogger.Println(msg)
	}
	t.Log(msg)
}

func logEnd(t *testing.T) {
	msg := fmt.Sprintf("END: %s", t.Name())
	if suiteLogger != nil {
		suiteLogger.Println(msg)
	}
	t.Log(msg)
}

func logStep(t *testing.T, step string) {
	msg := fmt.Sprintf("STEP: %s | %s", t.Name(), step)
	if suiteLogger != nil {
		suiteLogger.Println(msg)
	}
	t.Log(msg)
}

func logError(t *testing.T, err error) {
	msg := fmt.Sprintf("ERROR: %s | %v", t.Name(), err)
	if suiteLogger != nil {
		suiteLogger.Println(msg)
	}
	t.Log(msg)
}

func ts() string {
	return time.Now().Format("20060102_150405")
}
