package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
"github.com/joho/godotenv"
	"github.com/tebeka/selenium"
)

func startBrowserStack(t *testing.T) (selenium.WebDriver, func()) {
	_ = godotenv.Load()

	t.Helper()

	user := os.Getenv("BROWSERSTACK_USERNAME")
	key := os.Getenv("BROWSERSTACK_ACCESS_KEY")
	if user == "" || key == "" {
		t.Fatalf("BrowserStack credentials not set (BROWSERSTACK_USERNAME / BROWSERSTACK_ACCESS_KEY)")
	}

	browser := os.Getenv("BROWSER")
	if browser == "" {
		browser = "chrome"
	}

	// ensure screenshot dir exists
	_ = os.MkdirAll(filepath.Join("tests", "target", "screenshots"), 0755)

	caps := selenium.Capabilities{
		"browserName": browser,
	}

	caps["bstack:options"] = map[string]interface{}{
		"os":          "Windows",
		"osVersion":   "11",
		"buildName":   "QuickBite Cross Browser",
		"sessionName": t.Name(),
		"local":       true, // required when your app is on localhost
	}

	remoteURL := fmt.Sprintf("https://%s:%s@hub.browserstack.com/wd/hub", user, key)

	driver, err := selenium.NewRemote(caps, remoteURL)
	if err != nil {
		t.Fatalf("BrowserStack: cannot start driver: %v", err)
	}

	_ = driver.SetImplicitWaitTimeout(10 * time.Second)

	cleanup := func() {
		if t.Failed() {
			if img, err := driver.Screenshot(); err == nil {
				p := filepath.Join("tests", "target", "screenshots", fmt.Sprintf("%s_bs.png", t.Name()))
				_ = os.WriteFile(p, img, 0644)
				t.Logf("Screenshot saved: %s", p)
			}
		}
		_ = driver.Quit()
	}

	return driver, cleanup
}
