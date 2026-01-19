package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func Test_QuickBite_Register_Login_Logout_Flow(t *testing.T) {
	logStart(t)
	defer logEnd(t)

	wd, cleanup := startChrome(t)
	defer cleanup()

	logStep(t, "Open auth page")
	if err := wd.Get(authURL); err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: open auth page: %v", err)
	}

	logStep(t, "Open Register tab")
	regTab, err := wd.FindElement(selenium.ByCSSSelector, "[data-tab='register']")
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: register tab not found: %v", err)
	}
	if err := regTab.Click(); err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: register tab click failed: %v", err)
	}

	waitUntil(t, 8*time.Second, func() (bool, error) {
		_, e := wd.FindElement(selenium.ByID, "registerEmail")
		return e == nil, nil
	})

	email := quickbiteUniqueEmail()
	password := "QuickBitePass_123!"

	logStep(t, "Fill Register form fields")
	nameEl, err := wd.FindElement(selenium.ByID, "registerName")
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: registerName not found: %v", err)
	}
	emailEl, err := wd.FindElement(selenium.ByID, "registerEmail")
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: registerEmail not found: %v", err)
	}
	passEl, err := wd.FindElement(selenium.ByID, "registerPassword")
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: registerPassword not found: %v", err)
	}
	confEl, err := wd.FindElement(selenium.ByID, "registerConfirmPassword")
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: registerConfirmPassword not found: %v", err)
	}

	_ = nameEl.Clear()
	_ = nameEl.SendKeys("QuickBite Tester")
	_ = emailEl.Clear()
	_ = emailEl.SendKeys(email)
	_ = passEl.Clear()
	_ = passEl.SendKeys(password)
	_ = confEl.Clear()
	_ = confEl.SendKeys(password)

	registerBtnSel := []struct {
		by   string
		sel  string
		name string
	}{
		{selenium.ByCSSSelector, "#registerForm button[type='submit']", "register submit button"},
		{selenium.ByCSSSelector, "#registerForm button", "register form button"},
		{selenium.ByID, "registerButton", "registerButton id"},
		{selenium.ByCSSSelector, "button[data-action='register']", "data-action=register"},
	}

	logStep(t, "Click Register button")
	var regBtn selenium.WebElement
	for _, s := range registerBtnSel {
		btn, e := wd.FindElement(s.by, s.sel)
		if e == nil {
			regBtn = btn
			break
		}
	}
	if regBtn == nil {
		t.Fatalf("QuickBite: Register button not found (update selector)")
	}
	if err := regBtn.Click(); err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: Register button click failed: %v", err)
	}

	logStep(t, "Wait for registration alert ")
	waitUntil(t, 8*time.Second, func() (bool, error) {
		_, e := wd.AlertText()
		return e == nil, nil
	})
	alertText, _ := wd.AlertText()
	_ = wd.AcceptAlert()

	if alertText != "" && !strings.Contains(strings.ToLower(alertText), "success") {
		t.Logf("QuickBite: registration alert text: %q", alertText)
	}

	logStep(t, "Open Login tab")
	loginTab, err := wd.FindElement(selenium.ByCSSSelector, "[data-tab='login']")
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: login tab not found: %v", err)
	}
	if err := loginTab.Click(); err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: login tab click failed: %v", err)
	}

	waitUntil(t, 8*time.Second, func() (bool, error) {
		_, e := wd.FindElement(selenium.ByID, "loginEmail")
		return e == nil, nil
	})

	logStep(t, "Fill Login form fields")
	loginEmail, err := wd.FindElement(selenium.ByID, "loginEmail")
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: loginEmail not found: %v", err)
	}
	loginPass, err := wd.FindElement(selenium.ByID, "loginPassword")
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: loginPassword not found: %v", err)
	}
	_ = loginEmail.Clear()
	_ = loginEmail.SendKeys(email)
	_ = loginPass.Clear()
	_ = loginPass.SendKeys(password)

	loginBtnSel := []struct {
		by   string
		sel  string
		name string
	}{
		{selenium.ByCSSSelector, "#loginForm button[type='submit']", "login submit button"},
		{selenium.ByCSSSelector, "#loginForm button", "login form button"},
		{selenium.ByID, "loginButton", "loginButton id"},
		{selenium.ByCSSSelector, "button[data-action='login']", "data-action=login"},
	}

	logStep(t, "Click Login button")
	var loginBtn selenium.WebElement
	for _, s := range loginBtnSel {
		btn, e := wd.FindElement(s.by, s.sel)
		if e == nil {
			loginBtn = btn
			break
		}
	}
	if loginBtn == nil {
		t.Fatalf("QuickBite: Login button not found (update selector)")
	}
	if err := loginBtn.Click(); err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: Login button click failed: %v", err)
	}

	logStep(t, "Handle login alert (optional)")
	func() {
		deadline := time.Now().Add(2 * time.Second)
		for time.Now().Before(deadline) {
			if txt, e := wd.AlertText(); e == nil {
				_ = wd.AcceptAlert()
				t.Logf("QuickBite: login alert text: %q", txt)
				return
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	logStep(t, "Wait token in localStorage")
	waitUntil(t, 20*time.Second, func() (bool, error) {
		v, e := wd.ExecuteScript("return localStorage.getItem('token');", nil)
		if e != nil {
			return false, nil
		}
		s, _ := v.(string)
		return s != "", nil
	})

	logStep(t, "Wait redirect to index or logout button visible")
	waitUntil(t, 20*time.Second, func() (bool, error) {
		u, e := wd.CurrentURL()
		if e == nil && strings.Contains(u, "index.html") {
			return true, nil
		}
		_, e = wd.FindElement(selenium.ByID, "logoutButton")
		return e == nil, nil
	})

	tokenRaw, err := wd.ExecuteScript("return localStorage.getItem('token');", nil)
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: read token from localStorage: %v", err)
	}
	userIDRaw, err := wd.ExecuteScript("return localStorage.getItem('userId');", nil)
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: read userId from localStorage: %v", err)
	}

	token, _ := tokenRaw.(string)
	userID, _ := userIDRaw.(string)

	if token == "" || userID == "" {
		t.Fatalf("QuickBite: expected token and userId after login, got token=%q userId=%q", token, userID)
	}

	logStep(t, "Click Logout")
	logoutBtn, err := wd.FindElement(selenium.ByID, "logoutButton")
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: logoutButton not found: %v", err)
	}
	_ = logoutBtn.Click()

	logStep(t, "Wait redirect to auth.html")
	waitUntil(t, 12*time.Second, func() (bool, error) {
		u, e := wd.CurrentURL()
		if e != nil {
			return false, e
		}
		return strings.Contains(u, "auth.html"), nil
	})

	logStep(t, "Verify localStorage cleared")
	token2, _ := wd.ExecuteScript("return localStorage.getItem('token');", nil)
	user2, _ := wd.ExecuteScript("return localStorage.getItem('userId');", nil)
	if token2 != nil || user2 != nil {
		t.Fatalf("QuickBite: expected localStorage cleared after logout, got token=%v userId=%v", token2, user2)
	}
}
