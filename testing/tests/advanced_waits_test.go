package tests

import (
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func Test_QuickBite_Waits_Implicit_Explicit_Fluent(t *testing.T) {
	wd, cleanup := startChrome(t)
	defer cleanup()

	_ = wd.SetImplicitWaitTimeout(10 * time.Second)

	if err := wd.Get(authURL); err != nil {
		t.Fatalf("QuickBite: open auth page: %v", err)
	}

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

	nameEl, err := wd.FindElement(selenium.ByID, "registerName")
	if err != nil {
		t.Fatalf("QuickBite: registerName not found: %v", err)
	}
	emailEl, err := wd.FindElement(selenium.ByID, "registerEmail")
	if err != nil {
		t.Fatalf("QuickBite: registerEmail not found: %v", err)
	}
	passEl, err := wd.FindElement(selenium.ByID, "registerPassword")
	if err != nil {
		t.Fatalf("QuickBite: registerPassword not found: %v", err)
	}
	confEl, err := wd.FindElement(selenium.ByID, "registerConfirmPassword")
	if err != nil {
		t.Fatalf("QuickBite: registerConfirmPassword not found: %v", err)
	}

	_ = nameEl.Clear()
	_ = nameEl.SendKeys("Waits Tester")
	_ = emailEl.Clear()
	_ = emailEl.SendKeys(email)
	_ = passEl.Clear()
	_ = passEl.SendKeys(password)
	_ = confEl.Clear()
	_ = confEl.SendKeys(password)

	btn, err := wd.FindElement(selenium.ByCSSSelector, "#registerForm button")
	if err != nil {
		t.Fatalf("QuickBite: register button not found: %v", err)
	}
	_ = btn.Click()

	func() {
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			if _, err := wd.AlertText(); err == nil {
				_ = wd.AcceptAlert()
				return
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	waitUntil(t, 15*time.Second, func() (bool, error) {
		_, err := wd.FindElement(selenium.ByID, "loginEmail")
		return err == nil, nil
	})

	time.Sleep(2 * time.Second)

	fluentWait(t, 15*time.Second, 300*time.Millisecond, func() (bool, error) {
	el, err := wd.FindElement(selenium.ByID, "loginEmail")
	if err != nil {
		return false, nil
	}

	enabled, _ := el.IsEnabled()
	if !enabled {
		return false, nil
	}

	_ = el.Clear()
	err = el.SendKeys("waits@test.com")

	return err == nil, nil
})


	time.Sleep(2 * time.Second)
}
