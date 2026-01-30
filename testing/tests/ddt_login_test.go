package tests

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/tebeka/selenium"
	"github.com/xuri/excelize/v2"
)

type LoginCase struct {
	Row      int
	CaseID   string
	Email    string
	Password string
	Expected string 
}

func readLoginCasesXLSX(f *excelize.File, sheet string) ([]LoginCase, error) {
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("no data rows in %s", sheet)
	}

	get := func(r []string, idx int) string {
		if idx < len(r) {
			return strings.TrimSpace(r[idx])
		}
		return ""
	}

	var cases []LoginCase
	for i := 1; i < len(rows); i++ { 
		r := rows[i]

		c := LoginCase{
			Row:      i + 1, 
			CaseID:   get(r, 0),
			Email:    get(r, 1),
			Password: get(r, 2),
			Expected: strings.ToUpper(get(r, 3)),
		}

		if c.CaseID == "" && c.Email == "" && c.Password == "" && c.Expected == "" {
			continue
		}
		if c.CaseID == "" {
			c.CaseID = fmt.Sprintf("ROW-%d", c.Row)
		}
		cases = append(cases, c)
	}
	return cases, nil
}

func Test_DDT_Login_FromExcel(t *testing.T) {
	logStart(t)
	defer logEnd(t)

	const (
		xlsxPath = "login_data.xlsx"
		sheet    = "LoginData"
	)
	f, err := excelize.OpenFile(xlsxPath)
	if err != nil {
		logError(t, err)
		t.Fatalf("open xlsx: %v", err)
	}
	defer func() {
		_ = f.Save()
		_ = f.Close()
	}()

	cases, err := readLoginCasesXLSX(f, sheet)
	if err != nil {
		logError(t, err)
		t.Fatalf("read xlsx: %v", err)
	}
	if len(cases) == 0 {
		t.Fatalf("no test cases found in %s", xlsxPath)
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.CaseID, func(t *testing.T) {
			logStep(t, "Start dataset: "+tc.CaseID)

			wd, cleanup := startChrome(t)
			defer cleanup()

			if err := wd.Get(authURL); err != nil {
				t.Fatalf("open auth: %v", err)
			}

			loginTab, err := wd.FindElement(selenium.ByCSSSelector, "[data-tab='login']")
			if err != nil {
				t.Fatalf("login tab not found: %v", err)
			}
			_ = loginTab.Click()

			waitUntil(t, 10*time.Second, func() (bool, error) {
				_, e := wd.FindElement(selenium.ByID, "loginEmail")
				return e == nil, nil
			})

			emailEl, _ := wd.FindElement(selenium.ByID, "loginEmail")
			passEl, _ := wd.FindElement(selenium.ByID, "loginPassword")
			_ = emailEl.Clear()
			_ = emailEl.SendKeys(tc.Email)
			_ = passEl.Clear()
			_ = passEl.SendKeys(tc.Password)

			btn, err := wd.FindElement(selenium.ByCSSSelector, "#loginForm button")
			if err != nil {
				t.Fatalf("login button not found: %v", err)
			}
			_, _ = wd.ExecuteScript("arguments[0].click();", []interface{}{btn})

			deadline := time.Now().Add(8 * time.Second)
			gotToken := false
			for time.Now().Before(deadline) {
				v, _ := wd.ExecuteScript("return localStorage.getItem('token');", nil)
				if s, ok := v.(string); ok && s != "" {
					gotToken = true
					break
				}
				time.Sleep(250 * time.Millisecond)
			}

			alertText := ""
			if txt, e := wd.AlertText(); e == nil {
				alertText = txt
				_ = wd.AcceptAlert()
			}

			actual := "FAIL"
			if gotToken {
				actual = "PASS"
			}

			cellActual := fmt.Sprintf("E%d", tc.Row)
			cellMsg := fmt.Sprintf("F%d", tc.Row)
			_ = f.SetCellValue(sheet, cellActual, actual)
			_ = f.SetCellValue(sheet, cellMsg, alertText)

			t.Logf("Expected=%s | Actual=%s | Token=%v | Alert=%q", tc.Expected, actual, gotToken, alertText)

			if tc.Expected != "PASS" && tc.Expected != "FAIL" {
				t.Fatalf("Invalid expected value in Excel (row %d). Expected must be PASS or FAIL, got=%q", tc.Row, tc.Expected)
			}
			if tc.Expected != actual {
				t.Fatalf("Mismatch for %s: expected=%s actual=%s (row %d). Alert=%q", tc.CaseID, tc.Expected, actual, tc.Row, alertText)
			}

			logStep(t, "Dataset finished: "+tc.CaseID)
		})
	}

	if err := f.Save(); err != nil {
		t.Logf("WARN: failed to save updated Excel file: %v", err)
	}
}
