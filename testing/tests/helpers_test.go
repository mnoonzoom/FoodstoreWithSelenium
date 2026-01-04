package tests

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

const (
	baseURL  = "http://localhost:8082"
	authURL  = baseURL + "/auth.html"
	indexURL = baseURL + "/index.html"
	chromeDriverPath = "chromedriver"
)

func startChrome(t *testing.T) (selenium.WebDriver, func()) {
	t.Helper()

	const port = 9515

	service, err := selenium.NewChromeDriverService(chromeDriverPath, port)
	if err != nil {
		t.Fatalf("QuickBite: failed to start ChromeDriver: %v", err)
	}

	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := map[string]interface{}{
		"args": []string{
			"--start-maximized",
		},
	}
	caps["goog:chromeOptions"] = chromeCaps

	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		_ = service.Stop()
		t.Fatalf("QuickBite: failed to start WebDriver: %v", err)
	}

	cleanup := func() {
		_ = driver.Quit()
		_ = service.Stop()
	}
	return driver, cleanup
}

func quickbiteUniqueEmail() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("quickbite_%d@mail.test", rand.Intn(1_000_000_000))
}

func waitUntil(t *testing.T, timeout time.Duration, cond func() (bool, error)) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		ok, err := cond()
		if err == nil && ok {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("QuickBite: condition not met within %v", timeout)
}
