package tests

import (
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func Test_QuickBite_Select_Class_Category_FindDrinks(t *testing.T) {
	wd, cleanup := startChrome(t)
	defer cleanup()

	loginAndOpenMenu(t, wd)

	waitUntil(t, 15*time.Second, func() (bool, error) {
		_, err := wd.FindElement(selenium.ByID, "categorySelect")
		return err == nil, nil
	})

	time.Sleep(2 * time.Second)

	sel, err := wd.FindElement(selenium.ByID, "categorySelect")
	if err != nil {
		t.Fatalf("QuickBite: categorySelect not found: %v", err)
	}
	_ = sel.Click()

	time.Sleep(1 * time.Second)

	opt, err := wd.FindElement(
		selenium.ByXPATH,
		"//select[@id='categorySelect']/option[@value='drinks']",
	)
	if err != nil {
		t.Fatalf("QuickBite: drinks option not found: %v", err)
	}
	_ = opt.Click()

	time.Sleep(2 * time.Second)

	val, err := wd.ExecuteScript(
		"return document.getElementById('categorySelect').value;",
		nil,
	)
	if err != nil {
		t.Fatalf("QuickBite: read categorySelect value: %v", err)
	}
	if val != "drinks" {
		t.Fatalf("QuickBite: expected selected value 'drinks', got %v", val)
	}

	time.Sleep(2 * time.Second)

	waitUntil(t, 10*time.Second, func() (bool, error) {
		cards, err := wd.FindElements(
			selenium.ByCSSSelector,
			"#recommendedGrid .menu-card",
		)
		return err == nil && len(cards) > 0, nil
	})

	time.Sleep(3 * time.Second)
}
