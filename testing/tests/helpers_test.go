package tests

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

const (
	baseURL          = "http://localhost:8082"
	authURL          = baseURL + "/auth.html"
	indexURL         = baseURL + "/index.html"
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
	caps["goog:chromeOptions"] = map[string]interface{}{
		"args": []string{
			"--start-maximized",
		},
	}

	driver, err := selenium.NewRemote(
		caps,
		fmt.Sprintf("http://localhost:%d/wd/hub", port),
	)
	if err != nil {
		_ = service.Stop()
		t.Fatalf("QuickBite: failed to start WebDriver: %v", err)
	}

	_ = driver.SetImplicitWaitTimeout(10 * time.Second)

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
	t.Fatalf("QuickBite: explicit wait condition not met within %v", timeout)
}

func fluentWait(t *testing.T, timeout, poll time.Duration, cond func() (bool, error)) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		ok, err := cond()
		if err == nil && ok {
			return
		}
		time.Sleep(poll)
	}
	t.Fatalf("QuickBite: fluent wait condition not met within %v", timeout)
}

func actionFocusClickType(
	t *testing.T,
	wd selenium.WebDriver,
	el selenium.WebElement,
	text string,
) {
	t.Helper()

	_, err := wd.ExecuteScript(`
arguments[0].scrollIntoView({block:"center"});
arguments[0].focus();
arguments[0].click();
return true;
`, []interface{}{el})
	if err != nil {
		t.Fatalf("QuickBite: action focus/click failed: %v", err)
	}

	_ = el.Clear()
	if err := el.SendKeys(text); err != nil {
		t.Fatalf("QuickBite: action sendKeys failed: %v", err)
	}
}

func loginAndOpenMenu(t *testing.T, wd selenium.WebDriver) {
	t.Helper()

	if err := wd.Get(authURL); err != nil {
		t.Fatalf("QuickBite: open auth page: %v", err)
	}

	waitUntil(t, 15*time.Second, func() (bool, error) {
		_, err := wd.FindElement(selenium.ByTagName, "body")
		return err == nil, nil
	})

	regTab, err := wd.FindElement(selenium.ByCSSSelector, "[data-tab='register']")
	if err != nil {
		t.Fatalf("QuickBite: register tab not found: %v", err)
	}
	_ = regTab.Click()

	waitUntil(t, 15*time.Second, func() (bool, error) {
		_, err := wd.FindElement(selenium.ByID, "registerEmail")
		return err == nil, nil
	})

	email := quickbiteUniqueEmail()
	password := "QuickBitePass_123!"

	nameEl, _ := wd.FindElement(selenium.ByID, "registerName")
	emailEl, _ := wd.FindElement(selenium.ByID, "registerEmail")
	passEl, _ := wd.FindElement(selenium.ByID, "registerPassword")
	confEl, _ := wd.FindElement(selenium.ByID, "registerConfirmPassword")

	_ = nameEl.Clear()
	_ = nameEl.SendKeys("QuickBite Advanced Tester")
	_ = emailEl.Clear()
	_ = emailEl.SendKeys(email)
	_ = passEl.Clear()
	_ = passEl.SendKeys(password)
	_ = confEl.Clear()
	_ = confEl.SendKeys(password)

	clicked := false

	if btn, err := wd.FindElement(selenium.ByCSSSelector, "#registerForm button"); err == nil {
		_, _ = wd.ExecuteScript("arguments[0].click();", []interface{}{btn})
		clicked = true
	}

	if !clicked {
		if btn, err := wd.FindElement(selenium.ByXPATH, "//form[@id='registerForm']//button[contains(.,'Register') or contains(.,'Sign Up')]"); err == nil {
			_, _ = wd.ExecuteScript("arguments[0].click();", []interface{}{btn})
			clicked = true
		}
	}

	if !clicked {
		t.Fatalf("QuickBite: could not find Register button to click")
	}

	func() {
		deadline := time.Now().Add(6 * time.Second)
		for time.Now().Before(deadline) {
			if _, err := wd.AlertText(); err == nil {
				_ = wd.AcceptAlert()
				return
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	loginTab, err := wd.FindElement(selenium.ByCSSSelector, "[data-tab='login']")
	if err != nil {
		t.Fatalf("QuickBite: login tab not found: %v", err)
	}
	_ = loginTab.Click()

	waitUntil(t, 15*time.Second, func() (bool, error) {
		_, err := wd.FindElement(selenium.ByID, "loginEmail")
		return err == nil, nil
	})

	loginEmail, _ := wd.FindElement(selenium.ByID, "loginEmail")
	loginPass, _ := wd.FindElement(selenium.ByID, "loginPassword")
	_ = loginEmail.Clear()
	_ = loginEmail.SendKeys(email)
	_ = loginPass.Clear()
	_ = loginPass.SendKeys(password)

	clicked = false

	if btn, err := wd.FindElement(selenium.ByCSSSelector, "#loginForm button"); err == nil {
		_, _ = wd.ExecuteScript("arguments[0].click();", []interface{}{btn})
		clicked = true
	}
	if !clicked {
		if btn, err := wd.FindElement(selenium.ByXPATH, "//form[@id='loginForm']//button[contains(.,'Login') or contains(.,'Sign In')]"); err == nil {
			_, _ = wd.ExecuteScript("arguments[0].click();", []interface{}{btn})
			clicked = true
		}
	}

	if !clicked {
		t.Fatalf("QuickBite: could not find Login button to click")
	}

	func() {
		deadline := time.Now().Add(4 * time.Second)
		for time.Now().Before(deadline) {
			if _, err := wd.AlertText(); err == nil {
				_ = wd.AcceptAlert()
				return
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	fluentWait(t, 40*time.Second, 250*time.Millisecond, func() (bool, error) {
		v, err := wd.ExecuteScript("return localStorage.getItem('token');", nil)
		if err != nil {
			return false, nil
		}
		s, _ := v.(string)
		return s != "", nil
	})

	waitUntil(t, 25*time.Second, func() (bool, error) {
		_, err := wd.FindElement(selenium.ByID, "recommendedGrid")
		return err == nil, nil
	})
}
