package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func Test_QuickBite_Select_Class_Category_FindDrinks(t *testing.T) {
	logStart(t)
	defer logEnd(t)

	wd, cleanup := startChrome(t)
	defer cleanup()

	logStep(t, "Login and open menu page")
	loginAndOpenMenu(t, wd)

	logStep(t, "Wait for categorySelect to be present")
	waitUntil(t, 15*time.Second, func() (bool, error) {
		_, err := wd.FindElement(selenium.ByID, "categorySelect")
		return err == nil, nil
	})

_, err := wd.FindElement(selenium.ByID, "categorySelect")
if err != nil {
    t.Fatalf("QuickBite: categorySelect not found: %v", err)
}


	logStep(t, "Select category value = drinks (JS select-equivalent)")
	_, err = wd.ExecuteScript(`
		var s = document.getElementById('categorySelect');
		if(!s) return "NO_SELECT";
		s.value = "drinks";
		s.dispatchEvent(new Event('change', {bubbles:true}));
		return s.value;
	`, nil)
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: select drinks via JS failed: %v", err)
	}

	logStep(t, "Verify selected value is 'drinks'")
	val, err := wd.ExecuteScript(
		"return document.getElementById('categorySelect') && document.getElementById('categorySelect').value;",
		nil,
	)
	if err != nil {
		logError(t, err)
		t.Fatalf("QuickBite: read categorySelect value: %v", err)
	}
	if s, ok := val.(string); !ok || s != "drinks" {
		t.Fatalf("QuickBite: expected selected value 'drinks', got %v", val)
	}

	logStep(t, "Wait for filtered menu cards in #recommendedGrid")
	waitUntil(t, 12*time.Second, func() (bool, error) {
		cards, err := wd.FindElements(selenium.ByCSSSelector, "#recommendedGrid .menu-card")
		return err == nil && len(cards) > 0, nil
	})

	logStep(t, "Optional: check at least one card text contains drink-related keyword")
	cards, _ := wd.FindElements(selenium.ByCSSSelector, "#recommendedGrid .menu-card")
	if len(cards) > 0 {
		txt, _ := cards[0].Text()

		if !strings.Contains(strings.ToLower(txt), "drink") && !strings.Contains(strings.ToLower(txt), "cola") {
			t.Logf("QuickBite: first card text after drinks filter: %q", txt)
		}
	}

	logStep(t, "Done")
}
