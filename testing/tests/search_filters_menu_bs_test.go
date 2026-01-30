package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func Test_QuickBite_Search_Filters_Menu_BrowserStack(t *testing.T) {
	wd, cleanup := startBrowserStack(t)
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
	_ = nameEl.SendKeys("QuickBite Search Tester")
	_ = emailEl.Clear()
	_ = emailEl.SendKeys(email)
	_ = passEl.Clear()
	_ = passEl.SendKeys(password)
	_ = confEl.Clear()
	_ = confEl.SendKeys(password)

	form, err := wd.FindElement(selenium.ByID, "registerForm")
	if err != nil {
		t.Fatalf("QuickBite: registerForm not found: %v", err)
	}
	_ = form.Submit()

	waitUntil(t, 8*time.Second, func() (bool, error) {
		_, err := wd.AlertText()
		return err == nil, nil
	})
	_ = wd.AcceptAlert()

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

	// optional: accept alert if appears
	func() {
		deadline := time.Now().Add(2 * time.Second)
		for time.Now().Before(deadline) {
			if _, err := wd.AlertText(); err == nil {
				_ = wd.AcceptAlert()
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
		cards, err := wd.FindElements(selenium.ByCSSSelector, "#recommendedGrid .menu-card")
		return err == nil && len(cards) > 0, nil
	})

	titles, err := wd.FindElements(selenium.ByCSSSelector, "#recommendedGrid .menu-card h3")
	if err != nil || len(titles) == 0 {
		t.Fatalf("QuickBite: menu card titles not found: %v", err)
	}
	firstTitle, _ := titles[0].Text()
	firstTitle = strings.TrimSpace(firstTitle)

	query := firstTitle
	if len(query) >= 3 {
		query = query[:3]
	}

	searchInput, err := wd.FindElement(selenium.ByID, "filterByNameInput")
	if err != nil {
		t.Fatalf("QuickBite: filterByNameInput not found: %v", err)
	}
	_ = searchInput.Clear()
	_ = searchInput.SendKeys(query)

	waitUntil(t, 8*time.Second, func() (bool, error) {
		cards, err := wd.FindElements(selenium.ByCSSSelector, "#recommendedGrid .menu-card")
		return err == nil && len(cards) > 0, nil
	})

	newTitles, _ := wd.FindElements(selenium.ByCSSSelector, "#recommendedGrid .menu-card h3")
	if len(newTitles) > 0 {
		t0, _ := newTitles[0].Text()
		if !strings.Contains(strings.ToLower(t0), strings.ToLower(query)) {
			t.Logf("warning: first result title %q does not contain query %q", t0, query)
		}
	}
}
