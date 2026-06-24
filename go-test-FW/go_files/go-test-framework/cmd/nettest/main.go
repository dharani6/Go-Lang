// Command nettest is the entry point for the Go Network Test Framework demo.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nettest/framework/internal/common"
	"github.com/nettest/framework/internal/netconf"
	"github.com/nettest/framework/internal/rest"
	"github.com/nettest/framework/internal/runner"
	"github.com/nettest/framework/internal/snmpconf"
	"github.com/nettest/framework/internal/ui"
)

func main() {
	// ── CLI flags ──────────────────────────────────────────────────────────
	parallel  := flag.Bool("parallel", false, "Run test suites in parallel")
	verbose   := flag.Bool("verbose", false, "Show individual test steps")
	failFast  := flag.Bool("fail-fast", false, "Stop suite on first failure")
	tags      := flag.String("tags", "", "Comma-separated tag filter (e.g. smoke,rest)")
	report    := flag.String("report", "reports/results.json", "Path to JSON report output")
	suiteFlag := flag.String("suite", "all", "Suites to run: all|ui|rest|netconf|snmp")
	flag.Parse()

	var tagList []string
	if *tags != "" {
		tagList = strings.Split(*tags, ",")
	}

	opts := runner.Options{
		Parallel:     *parallel,
		MaxWorkers:   4,
		Tags:         tagList,
		FailFast:     *failFast,
		Verbose:      *verbose,
		OutputFormat: "console",
		ReportPath:   *report,
	}

	// Ensure reports directory exists
	_ = os.MkdirAll("reports", 0755)
	_ = os.MkdirAll("screenshots", 0755)

	r := runner.New(opts)

	// ── Register suites based on -suite flag ───────────────────────────────
	suite := strings.ToLower(*suiteFlag)

	if suite == "all" || suite == "ui" {
		r.AddSuite("Browser UI Tests", "UI/HTML", buildUISuite())
	}
	if suite == "all" || suite == "rest" {
		r.AddSuite("REST API Tests", "HTTP/REST", buildRESTSuite())
	}
	if suite == "all" || suite == "netconf" {
		r.AddSuite("NETCONF Tests", "NETCONF/SSH", buildNetconfSuite())
	}
	if suite == "all" || suite == "snmp" {
		r.AddSuite("SNMP Tests", "SNMP/UDP", buildSNMPSuite())
	}

	results := r.Run()

	// Exit non-zero if any test failed
	for _, suite := range results {
		_, fail, _, errCount := suite.Stats()
		if fail > 0 || errCount > 0 {
			fmt.Fprintf(os.Stderr, "\n⚠  Some tests failed — check report at %s\n", *report)
			os.Exit(1)
		}
	}
	fmt.Printf("\n✅  All tests passed! Report saved to %s\n", *report)
}
d
// ─── Suite builders ───────────────────────────────────────────────────────────

func buildUISuite() []common.TestCase {
	b := ui.NewBrowser("chromium", true, 10*time.Second)
	// Simulated URL — Navigate gracefully handles non-reachable hosts in demo mode
	baseURL := "http://localhost:8080"

	return []common.TestCase{
		ui.NewPageLoadTest(b, baseURL, []string{"body", "h1", "nav"}),
		ui.NewLoginFlowTest(b, baseURL+"/login", "admin", "admin", baseURL+"/dashboard"),
		ui.NewAccessibilityTest(b, baseURL),
		ui.NewScreenshotTest(b, baseURL, "screenshots/homepage.png"),
	}
}

func buildRESTSuite() []common.TestCase {
	// JSONPlaceholder — free public REST API for testing
	client := rest.NewClient("https://jsonplaceholder.typicode.com", 10*time.Second)
	client.SetHeader("Accept", "application/json")

	return []common.TestCase{
		rest.NewHealthCheckTest(rest.NewClient("https://jsonplaceholder.typicode.com", 10*time.Second)),
		rest.NewCRUDTest(client, "/posts", map[string]interface{}{
			"title":  "Go Test Framework Demo",
			"body":   "Testing CRUD operations",
			"userId": 1,
		}),
		rest.NewAuthTest(
			rest.NewClient("https://reqres.in/api", 10*time.Second),
			"/login",
			"eve.holt@reqres.in",
			"cityslicka",
		),
	}
}

func buildNetconfSuite() []common.TestCase {
	// Simulated — in production, point at a real NETCONF device
	session := netconf.NewSession("192.168.1.1", 830, "admin", "admin")

	return []common.TestCase{
		netconf.NewConnectivityTest(session),
		netconf.NewGetConfigTest(session),
		netconf.NewEditConfigTest(session),
		netconf.NewCapabilityTest(session, []string{
			"base:1.0",
			"candidate:1.0",
			"rollback-on-error",
		}),
	}
}

func buildSNMPSuite() []common.TestCase {
	// Simulated — in production, point at a real SNMP agent
	v2c := snmpconf.NewV2cClient("192.168.1.1", 161, "public")
	v3  := snmpconf.NewV3Client("192.168.1.1", 161, "snmpv3user", "SHA", "authkey123", "AES", "privkey123")
	_ = v3 // v3 client ready for extension

	return []common.TestCase{
		snmpconf.NewReachabilityTest(v2c),
		snmpconf.NewSystemInfoTest(v2c),
		snmpconf.NewMIBWalkTest(v2c, snmpconf.OIDIfTable, 2),
		snmpconf.NewGetBulkTest(v2c),
	}
}
