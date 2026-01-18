package main

import (
	"bufio"
	"encoding/json"
	"html/template"
	"os"
	"sort"
	"strings"
)

type Event struct {
	Action  string `json:"Action"`
	Package string `json:"Package"`
	Test    string `json:"Test"`
	Elapsed float64 `json:"Elapsed"`
	Output  string `json:"Output"`
}

type Row struct {
	Name    string
	Status  string
	Output  string
	Elapsed float64
	TC      TestCase
}

type TestCase struct {
	ID            string
	RequirementID string
	Title         string
	Preconditions []string
	Steps         []string
	TestData      []string
	Expected      []string
}

func main() {
	in := "tests/target/reports/test.json"
	out := "tests/target/reports/report.html"
	_ = os.MkdirAll("tests/target/reports", 0755)


	testCases := map[string]TestCase{
		"Test_QuickBite_Register_Login_Logout_Flow": {
			ID:            "TC-AUTH-001",
			RequirementID: "REQ-AUTH-001",
			Title:         "Register → Login → Logout flow creates and clears session",
			Preconditions: []string{
				"QuickBite is running at http://localhost:8082",
				"Auth page is доступна: /auth.html",
				"Chrome + chromedriver configured",
			},
			Steps: []string{
				"Open /auth.html",
				"Open Register tab",
				"Fill registration form (Name, Email, Password, Confirm Password)",
				"Click Register button",
				"Verify registration feedback (alert/message)",
				"Open Login tab",
				"Enter Email + Password",
				"Click Login button",
				"Verify token in localStorage OR logout button visible",
				"Click Logout",
				"Verify redirect to auth.html and localStorage cleared",
			},
			TestData: []string{
				"Name: QuickBite Tester",
				"Email: auto-generated unique (example: quickbite_123456@mail.test)",
				"Password: QuickBitePass_123!",
			},
			Expected: []string{
				"Registration succeeds (success confirmation shown)",
				"Login succeeds (token stored in localStorage, user is authenticated)",
				"Logout succeeds (redirect to auth, token/userId removed)",
			},
		},

		"Test_QuickBite_Select_Class_Category_FindDrinks": {
			ID:            "TC-MENU-001",
			RequirementID: "REQ-MENU-002",
			Title:         "Filter menu items by category 'Drinks' using dropdown",
			Preconditions: []string{
				"QuickBite is running at http://localhost:8082",
				"Valid user can login",
				"Menu page available after login: /index.html",
			},
			Steps: []string{
				"Login to the system",
				"Wait until categorySelect is visible",
				"Select category value 'drinks' in categorySelect",
				"Verify selected value == drinks",
				"Verify filtered cards appear in #recommendedGrid .menu-card",
			},
			TestData: []string{
				"Category value: drinks",
				"Valid credentials (created during test run or pre-existing test user)",
			},
			Expected: []string{
				"Dropdown value becomes 'drinks'",
				"At least 1 menu card is displayed after filtering",
			},
		},

		"Test_QuickBite_PlaceFoodOrder_FullFlow": {
			ID:            "TC-ORDER-001",
			RequirementID: "REQ-ORDER-001",
			Title:         "Place order end-to-end (add item → checkout → confirm)",
			Preconditions: []string{
				"QuickBite is running at http://localhost:8082",
				"Menu contains at least one item with an Order button",
				"Checkout form is available after adding item to cart",
			},
			Steps: []string{
				"Register and login",
				"Wait menu cards appear",
				"Click Order on first menu card",
				"Verify cart has at least 1 item",
				"Click Checkout",
				"Fill customer details (name, address, phone)",
				"Click Confirm Order",
				"Verify confirmation alert contains 'Order placed' and 'Order ID'",
			},
			TestData: []string{
				"Customer Name: QuickBite Student",
				"Address: Astana, Kazakhstan",
				"Phone: 87001234567",
			},
			Expected: []string{
				"Order confirmation alert is shown",
				"Alert contains 'Order placed' and 'Order ID' (with actual ID value)",
			},
		},
	}


	f, err := os.Open(in)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	tests := map[string]*Row{}

	sc := bufio.NewScanner(f)

	buf := make([]byte, 0, 1024*1024)
	sc.Buffer(buf, 10*1024*1024)

	for sc.Scan() {
		var e Event
		_ = json.Unmarshal(sc.Bytes(), &e)

		if e.Test == "" {
			continue
		}

		r, ok := tests[e.Test]
		if !ok {
			r = &Row{Name: e.Test, Status: "RUNNING"}
	
			if tc, ok2 := testCases[e.Test]; ok2 {
				r.TC = tc
			}
			tests[e.Test] = r
		}

		switch e.Action {
		case "output":
	
			r.Output += e.Output
		case "pass":
			r.Status = "PASS"
			r.Elapsed = e.Elapsed
		case "fail":
			r.Status = "FAIL"
			r.Elapsed = e.Elapsed
		case "skip":
			r.Status = "SKIP"
			r.Elapsed = e.Elapsed
		}
	}


	var rows []Row
	pass, fail, skip := 0, 0, 0
	for _, r := range tests {

		if r.TC.ID == "" {
			r.TC = TestCase{
				ID:            "TC-UNKNOWN",
				RequirementID: "REQ-UNKNOWN",
				Title:         "No test case metadata mapped for this test",
				Preconditions: []string{"N/A"},
				Steps:         []string{"N/A"},
				TestData:      []string{"N/A"},
				Expected:      []string{"N/A"},
			}
		}
		rows = append(rows, *r)
		switch r.Status {
		case "PASS":
			pass++
		case "FAIL":
			fail++
		case "SKIP":
			skip++
		}
	}

	
	sort.Slice(rows, func(i, j int) bool { return rows[i].Name < rows[j].Name })

	funcMap := template.FuncMap{
		"nl2br": func(s string) template.HTML {
			escaped := template.HTMLEscapeString(s)
			escaped = strings.ReplaceAll(escaped, "\n", "<br/>")
			return template.HTML(escaped)
		},
	}

	tpl := template.Must(template.New("r").Funcs(funcMap).Parse(`
<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>Foodstore Go Test Report (Test Cases)</title>
  <style>
    body { font-family: Arial, sans-serif; padding: 16px; }
    .sum { margin: 12px 0 18px; padding: 10px; background:#f7f7f7; border:1px solid #ddd; }
    .PASS { color: green; font-weight: bold; }
    .FAIL { color: red; font-weight: bold; }
    .SKIP { color: orange; font-weight: bold; }
    .card { border:1px solid #ddd; border-radius: 8px; padding: 12px; margin: 12px 0; }
    .muted { color:#666; }
    .kv { margin: 6px 0; }
    .kv b { display:inline-block; width: 160px; }
    ol, ul { margin: 6px 0 6px 22px; }
    .out { margin-top:10px; padding:10px; background:#fafafa; border:1px dashed #ccc; }
    .out h4 { margin:0 0 6px; }
    .small { font-size: 13px; }
  </style>
</head>
<body>
  <h2>Foodstore Automation Report (Go testing + Selenium)</h2>

  <div class="sum">
    <b>Execution summary:</b>
    Pass: {{.Pass}} | Fail: {{.Fail}} | Skip: {{.Skip}} <br/>
    <span class="muted small">
      Logs: tests/target/logs/test-run.log • Screenshots: tests/target/screenshots/
    </span>
  </div>

  {{range .Rows}}
  <div class="card">
    <div class="kv"><b>1. ID:</b> {{.TC.ID}}</div>
    <div class="kv"><b>2. Requirement ID:</b> {{.TC.RequirementID}}</div>
    <div class="kv"><b>3. Test Case Title:</b> {{.TC.Title}}</div>
    <div class="kv"><b>8. Status:</b> <span class="{{.Status}}">{{.Status}}</span> <span class="muted small">(Elapsed: {{printf "%.2f" .Elapsed}}s)</span></div>

    <div class="kv"><b>4. Preconditions:</b></div>
    <ul>
      {{range .TC.Preconditions}}<li>{{.}}</li>{{end}}
    </ul>

    <div class="kv"><b>5. Test Steps:</b></div>
    <ol>
      {{range .TC.Steps}}<li>{{.}}</li>{{end}}
    </ol>

    <div class="kv"><b>6. Test Data:</b></div>
    <ul>
      {{range .TC.TestData}}<li>{{.}}</li>{{end}}
    </ul>

    <div class="kv"><b>7. Expected Result:</b></div>
    <ul>
      {{range .TC.Expected}}<li>{{.}}</li>{{end}}
    </ul>

    <div class="out">
      <h4>Actual Output / Logs (from go test):</h4>
      <div class="small">{{nl2br .Output}}</div>
    </div>

    <div class="muted small">Go Test Name: {{.Name}}</div>
  </div>
  {{end}}

</body>
</html>
`))

	of, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer of.Close()

	_ = tpl.Execute(of, map[string]any{
		"Rows": rows,
		"Pass": pass,
		"Fail": fail,
		"Skip": skip,
	})
}
