package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func Test_QuickBite_PlaceFoodOrder_FullFlow(t *testing.T) {
	logStart(t)
	defer logEnd(t)

	wd, cleanup := startChrome(t)
	defer cleanup()

	logStep(t, "Open auth page")
	if err := wd.Get(authURL); err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: open auth page: %v", err)
	}

	title1, err := wd.Title()
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: read title: %v", err)
	}
	if !strings.Contains(title1, "Authentication") || !strings.Contains(title1, "QuickBite") {
		t.Fatalf("QuickBite: unexpected auth page title: %q", title1)
	}

	logStep(t, "Open Register tab")
	regTab, err := wd.FindElement(selenium.ByCSSSelector, "[data-tab='register']")
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: register tab not found: %v", err)
	}
	_ = regTab.Click()

	waitUntil(t, 10*time.Second, func() (bool, error) {
		_, e := wd.FindElement(selenium.ByID, "registerEmail")
		return e == nil, nil
	})

	email := quickbiteUniqueEmail()
	password := "QuickBite_123!"

	logStep(t, "Fill Register form")
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
	_ = nameEl.SendKeys("QuickBite Order Tester")
	_ = emailEl.Clear()
	_ = emailEl.SendKeys(email)
	_ = passEl.Clear()
	_ = passEl.SendKeys(password)
	_ = confEl.Clear()
	_ = confEl.SendKeys(password)

	logStep(t, "Click Register button")
	regBtn, err := wd.FindElement(selenium.ByCSSSelector, "#registerForm button[type='submit']")
	if err != nil {
	
		regBtn, err = wd.FindElement(selenium.ByCSSSelector, "#registerForm button")
	}
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: Register button not found: %v", err)
	}
	_ = regBtn.Click()

	logStep(t, "Wait registration alert")
	waitUntil(t, 10*time.Second, func() (bool, error) {
		_, e := wd.AlertText()
		return e == nil, nil
	})
	regAlert, _ := wd.AlertText()
	_ = wd.AcceptAlert()
	if regAlert != "" && !strings.Contains(strings.ToLower(regAlert), "success") {
		t.Logf("QuickBite: registration alert text: %q", regAlert)
	}

	logStep(t, "Open Login tab")
	loginTab, err := wd.FindElement(selenium.ByCSSSelector, "[data-tab='login']")
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: login tab not found: %v", err)
	}
	_ = loginTab.Click()

	waitUntil(t, 10*time.Second, func() (bool, error) {
		_, e := wd.FindElement(selenium.ByID, "loginEmail")
		return e == nil, nil
	})

	logStep(t, "Fill Login form")
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

	logStep(t, "Click Login button")
	loginBtn, err := wd.FindElement(selenium.ByCSSSelector, "#loginForm button[type='submit']")
	if err != nil {
		loginBtn, err = wd.FindElement(selenium.ByCSSSelector, "#loginForm button")
	}
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: Login button not found: %v", err)
	}
	_ = loginBtn.Click()

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

	logStep(t, "Wait Menu page title")
	waitUntil(t, 20*time.Second, func() (bool, error) {
		tt, e := wd.Title()
		return e == nil && strings.Contains(tt, "QuickBite Menu"), nil
	})

	logStep(t, "Wait menu cards visible")
	waitUntil(t, 25*time.Second, func() (bool, error) {
		cards, e := wd.FindElements(selenium.ByCSSSelector, "#recommendedGrid .menu-card")
		return e == nil && len(cards) > 0, nil
	})

	logStep(t, "Click first Order button")
	orderBtn, err := wd.FindElement(
		selenium.ByXPATH,
		"//div[@id='recommendedGrid']//div[contains(@class,'menu-card')][1]//button[contains(.,'Order')]",
	)
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: Order button not found: %v", err)
	}
	_ = orderBtn.Click()

	logStep(t, "Wait cart item appears")
	waitUntil(t, 12*time.Second, func() (bool, error) {
		items, e := wd.FindElements(selenium.ByXPATH, "//*[@id='cartItems']//*[@class='cart-item']")
		return e == nil && len(items) >= 1, nil
	})

	logStep(t, "Click Checkout")
	checkoutBtn, err := wd.FindElement(selenium.ByCSSSelector, "#checkoutButton")
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: checkoutButton not found: %v", err)
	}
	_ = checkoutBtn.Click()

	logStep(t, "Wait checkout form enabled")
	waitUntil(t, 12*time.Second, func() (bool, error) {
		el, e := wd.FindElement(selenium.ByID, "customerName")
		if e != nil {
			return false, nil
		}
		enabled, _ := el.IsEnabled()
		return enabled, nil
	})

	logStep(t, "Fill customer details")
	cName, err := wd.FindElement(selenium.ByID, "customerName")
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: customerName not found: %v", err)
	}
	cAddr, err := wd.FindElement(selenium.ByID, "customerAddress")
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: customerAddress not found: %v", err)
	}
	cPhone, err := wd.FindElement(selenium.ByID, "customerPhone")
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: customerPhone not found: %v", err)
	}

	_ = cName.Clear()
	_ = cName.SendKeys("QuickBite Student")
	_ = cAddr.Clear()
	_ = cAddr.SendKeys("Astana, Kazakhstan")
	_ = cPhone.Clear()
	_ = cPhone.SendKeys("87001234567")

	logStep(t, "Click Confirm Order")
	confirmBtn, err := wd.FindElement(selenium.ByCSSSelector, "#confirmOrderButton")
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: confirmOrderButton not found: %v", err)
	}
	_ = confirmBtn.Click()

	logStep(t, "Wait confirmation alert")
	waitUntil(t, 15*time.Second, func() (bool, error) {
		_, e := wd.AlertText()
		return e == nil, nil
	})

	alertText, _ := wd.AlertText()
	_ = wd.AcceptAlert()
txt := strings.ToLower(alertText)
if !strings.Contains(txt, "order placed") || !strings.Contains(txt, "order id") {
    t.Fatalf("QuickBite: unexpected order confirmation alert: %q", alertText)
}



	t.Logf("QuickBite: order placed successfully. Alert: %q", alertText)
	logStep(t, "Order placed successfully")
}
