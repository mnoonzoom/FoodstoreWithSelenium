package tests

import (
	"testing"
	"time"
	"github.com/tebeka/selenium"
)

func Test_QuickBite_Actions_Class(t *testing.T) {
	wd, cleanup := startChrome(t)
	defer cleanup()

	loginAndOpenMenu(t, wd)

	searchInput, err := wd.FindElement(selenium.ByID, "filterByNameInput")
	if err != nil {
		t.Fatalf("QuickBite: search input not found: %v", err)
	}

	time.Sleep(2 * time.Second) 

	actionFocusClickType(t, wd, searchInput, "salad")

	time.Sleep(3 * time.Second) 
}
