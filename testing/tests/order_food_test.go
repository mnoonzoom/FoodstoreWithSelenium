package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func Test_QuickBite_PlaceFoodOrder_FullFlow(t *testing.T) {
	wd, cleanup := startChrome(t)
	defer cleanup()

	if err := wd.Get(authURL); err != nil {
		t.Fatalf("QuickBite: open auth page: %v", err)
	}

	title1, _ := wd.Title()
	if !strings.Contains(title1, "Authentication") || !strings.Contains(title1, "QuickBite") {
		t.Fatalf("QuickBite: unexpected auth page title: %q", title1)
	}

	regTab, err := wd.FindElement(selenium.ByCSSSelector, "[data-tab='register']")
	if err != nil {
		t.Fatalf("QuickBite: register tab not found: %v", err)
	}
	_ = regTab.Click()

	email := quickbiteUniqueEmail()
	password := "QuickBite_123!"

	nameEl, _ := wd.FindElement(selenium.ByID, "registerName")
	emailEl, _ := wd.FindElement(selenium.ByID, "registerEmail")
	passEl, _ := wd.FindElement(selenium.ByID, "registerPassword")
	confEl, _ := wd.FindElement(selenium.ByID, "registerConfirmPassword")

	_ = nameEl.Clear()
	_ = nameEl.SendKeys("QuickBite Order Tester")
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

	waitUntil(t, 10*time.Second, func() (bool, error) {
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
		tt, err := wd.Title()
		return err == nil && strings.Contains(tt, "QuickBite Menu"), nil
	})

	waitUntil(t, 25*time.Second, func() (bool, error) {
		cards, err := wd.FindElements(selenium.ByCSSSelector, "#recommendedGrid .menu-card")
		if err != nil {
			return false, nil
		}
		return len(cards) > 0, nil
	})

	orderBtn, err := wd.FindElement(
		selenium.ByXPATH,
		"//div[@id='recommendedGrid']//div[contains(@class,'menu-card')][1]//button[contains(.,'Order')]",
	)
	if err != nil {
		t.Fatalf("QuickBite: Order button not found (XPath): %v", err)
	}
	_ = orderBtn.Click()

	waitUntil(t, 10*time.Second, func() (bool, error) {
		items, err := wd.FindElements(selenium.ByXPATH, "//*[@id='cartItems']//*[@class='cart-item']")
		if err != nil {
			return false, nil
		}
		return len(items) >= 1, nil
	})

	checkoutBtn, err := wd.FindElement(selenium.ByCSSSelector, "#checkoutButton")
	if err != nil {
		t.Fatalf("QuickBite: checkoutButton not found (CSS): %v", err)
	}
	_ = checkoutBtn.Click()

	waitUntil(t, 10*time.Second, func() (bool, error) {
		el, err := wd.FindElement(selenium.ByID, "customerName")
		if err != nil {
			return false, nil
		}
		enabled, _ := el.IsEnabled()
		return enabled, nil
	})

	cName, _ := wd.FindElement(selenium.ByID, "customerName")
	cAddr, _ := wd.FindElement(selenium.ByID, "customerAddress")
	cPhone, _ := wd.FindElement(selenium.ByID, "customerPhone")

	_ = cName.Clear()
	_ = cName.SendKeys("QuickBite Student")
	_ = cAddr.Clear()
	_ = cAddr.SendKeys("Astana, Kazakhstan")
	_ = cPhone.Clear()
	_ = cPhone.SendKeys("87001234567")

	confirmBtn, err := wd.FindElement(selenium.ByCSSSelector, "#confirmOrderButton")
	if err != nil {
		t.Fatalf("QuickBite: confirmOrderButton not found (CSS): %v", err)
	}
	_ = confirmBtn.Click()

	waitUntil(t, 15*time.Second, func() (bool, error) {
		_, err := wd.AlertText()
		return err == nil, nil
	})

	alertText, _ := wd.AlertText()
	_ = wd.AcceptAlert()

	if !strings.Contains(alertText, "Order placed") || !strings.Contains(alertText, "Order ID") {
		t.Fatalf("QuickBite: unexpected order confirmation alert: %q", alertText)
	}

	t.Logf("QuickBite: order placed successfully. Alert: %q", alertText)
}
