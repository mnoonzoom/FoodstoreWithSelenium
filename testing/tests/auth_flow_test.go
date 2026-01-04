package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func Test_QuickBite_Register_Login_Logout_Flow(t *testing.T) {
	wd, cleanup := startChrome(t)
	defer cleanup()

	if err := wd.Get(authURL); err != nil {
		t.Fatalf("QuickBite: open auth page: %v", err)
	}

	regTab, err := wd.FindElement(selenium.ByCSSSelector, "[data-tab='register']")
	if err != nil {
		t.Fatalf("QuickBite: register tab not found: %v", err)
	}
	_ = regTab.Click()

	email := quickbiteUniqueEmail()
	password := "QuickBitePass_123!"

	nameEl, _ := wd.FindElement(selenium.ByID, "registerName")
	emailEl, _ := wd.FindElement(selenium.ByID, "registerEmail")
	passEl, _ := wd.FindElement(selenium.ByID, "registerPassword")
	confEl, _ := wd.FindElement(selenium.ByID, "registerConfirmPassword")

	_ = nameEl.Clear()
	_ = nameEl.SendKeys("QuickBite Tester")
	_ = emailEl.Clear()
	_ = emailEl.SendKeys(email)
	_ = passEl.Clear()
	_ = passEl.SendKeys(password)
	_ = confEl.Clear()
	_ = confEl.SendKeys(password)

	regForm, err := wd.FindElement(selenium.ByID, "registerForm")
	if err != nil {
		t.Fatalf("QuickBite: registerForm not found: %v", err)
	}
	_ = regForm.Submit()

	waitUntil(t, 8*time.Second, func() (bool, error) {
		_, err := wd.AlertText()
		return err == nil, nil
	})
	alertText, _ := wd.AlertText()
	_ = wd.AcceptAlert()

	if !strings.Contains(strings.ToLower(alertText), "success") {
		t.Logf("QuickBite: registration alert text: %q", alertText)
	}

	loginTab, err := wd.FindElement(selenium.ByCSSSelector, "[data-tab='login']")
	if err != nil {
		t.Fatalf("QuickBite: login tab not found: %v", err)
	}
	_ = loginTab.Click()

	loginEmail, _ := wd.FindElement(selenium.ByID, "loginEmail")
	loginPass, _ := wd.FindElement(selenium.ByID, "loginPassword")
	_ = loginEmail.Clear()
	_ = loginEmail.SendKeys(email)
	_ = loginPass.Clear()
	_ = loginPass.SendKeys(password)

	loginForm, err := wd.FindElement(selenium.ByID, "loginForm")
	if err != nil {
		t.Fatalf("QuickBite: loginForm not found: %v", err)
	}
	_ = loginForm.Submit()

	func() {
		deadline := time.Now().Add(2 * time.Second)
		for time.Now().Before(deadline) {
			if txt, err := wd.AlertText(); err == nil {
				_ = wd.AcceptAlert()
				t.Logf("QuickBite: login alert text: %q", txt)
				return
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	waitUntil(t, 20*time.Second, func() (bool, error) {
		v, err := wd.ExecuteScript("return localStorage.getItem('token');", nil)
		if err != nil {
			return false, nil
		}
		s, _ := v.(string)
		return s != "", nil
	})

	waitUntil(t, 20*time.Second, func() (bool, error) {
		u, err := wd.CurrentURL()
		if err == nil && strings.Contains(u, "index.html") {
			return true, nil
		}
		_, err = wd.FindElement(selenium.ByID, "logoutButton")
		return err == nil, nil
	})

	tokenRaw, err := wd.ExecuteScript("return localStorage.getItem('token');", nil)
	if err != nil {
		t.Fatalf("QuickBite: read token from localStorage: %v", err)
	}
	userIDRaw, err := wd.ExecuteScript("return localStorage.getItem('userId');", nil)
	if err != nil {
		t.Fatalf("QuickBite: read userId from localStorage: %v", err)
	}

	token, _ := tokenRaw.(string)
	userID, _ := userIDRaw.(string)

	if token == "" || userID == "" {
		t.Fatalf(
			"QuickBite: expected token and userId after login, got token=%q userId=%q",
			token, userID,
		)
	}

	logoutBtn, err := wd.FindElement(selenium.ByID, "logoutButton")
	if err != nil {
		t.Fatalf("QuickBite: logoutButton not found: %v", err)
	}
	_ = logoutBtn.Click()

	waitUntil(t, 12*time.Second, func() (bool, error) {
		u, err := wd.CurrentURL()
		if err != nil {
			return false, err
		}
		return strings.Contains(u, "auth.html"), nil
	})

	token2, _ := wd.ExecuteScript("return localStorage.getItem('token');", nil)
	user2, _ := wd.ExecuteScript("return localStorage.getItem('userId');", nil)
	if token2 != nil || user2 != nil {
		t.Fatalf(
			"QuickBite: expected localStorage cleared after logout, got token=%v userId=%v",
			token2, user2,
		)
	}
}
