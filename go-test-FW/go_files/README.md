# 🧪 Go Network Test Framework

A production-ready, extensible Go testing framework for network protocols —
**UI**, **REST API**, **NETCONF**, and **SNMP** — built as a learning demo
you can enhance incrementally.

---

## Project Structure

```
go-test-framework/
├── cmd/
│   └── nettest/
│       └── main.go              ← Entry point & suite wiring
├── internal/
│   ├── common/
│   │   ├── types.go             ← TestResult, TestCase, Config, Severity
│   │   └── assert.go            ← Fluent assertion helpers
│   ├── ui/
│   │   └── browser.go           ← Browser simulation (Playwright-style)
│   ├── rest/
│   │   └── client.go            ← HTTP client + test cases
│   ├── netconf/
│   │   └── session.go           ← NETCONF SSH session + test cases
│   └── snmpconf/
│       └── client.go            ← SNMP v1/v2c/v3 client + test cases
├── internal/runner/
│   └── runner.go                ← Test runner, parallel execution, reporting
├── reports/                     ← JSON reports (auto-created)
├── screenshots/                 ← UI screenshots (auto-created)
└── go.mod
```

---

## Quick Start

```bash
# Run all suites (serial, verbose)
go run ./cmd/nettest/ -verbose

# Run only NETCONF tests
go run ./cmd/nettest/ -suite netconf -verbose

# Run only smoke-tagged tests across all suites
go run ./cmd/nettest/ -tags smoke -verbose

# Run in parallel with fail-fast
go run ./cmd/nettest/ -parallel -fail-fast -verbose

# Custom report path
go run ./cmd/nettest/ -report /tmp/my-report.json
```

---

## CLI Flags

| Flag         | Default              | Description                                     |
|--------------|----------------------|-------------------------------------------------|
| `-suite`     | `all`                | `all\|ui\|rest\|netconf\|snmp`                  |
| `-verbose`   | `false`              | Print individual test steps                     |
| `-parallel`  | `false`              | Run suites with goroutine pool (MaxWorkers=4)   |
| `-fail-fast` | `false`              | Stop on first failure within a suite            |
| `-tags`      | _(all)_              | Comma-separated tag filter e.g. `smoke,rest`    |
| `-report`    | `reports/results.json` | Path for JSON report output                  |

---

## Protocol Coverage

### UI Testing
- **Browser simulation** (Playwright/Selenium-style API)
- Page load & navigation, element presence checks
- Login flow (fill → click → redirect assert)
- Accessibility checks (WCAG structural landmarks)
- Screenshot capture

**Real implementation upgrade:**
```go
import "github.com/playwright-community/playwright-go"
// Replace Browser.Navigate() with playwright.Launch() → page.Goto()
```

### REST API Testing
- `GET / POST / PUT / PATCH / DELETE` with JSON body
- Status code assertions, body content checks
- Bearer token & Basic Auth helpers
- Full CRUD lifecycle test
- Auth flow (valid + invalid credentials)

**Real implementation:** Wire to any live REST API.
Point `NewClient()` at your service URL.

### NETCONF (RFC 6241)
- SSH session establishment + hello exchange
- Capability negotiation (base, candidate, rollback)
- `<get-config>` / `<edit-config>` / `<commit>`
- Lock → Edit → Commit → Unlock lifecycle
- XML payload validation

**Real implementation upgrade:**
```go
import "golang.org/x/crypto/ssh"
// Replace Session.Connect() with real SSH dial + NETCONF framing
```

### SNMP (RFC 3416 / RFC 3414)
- SNMPv2c community-string authentication
- SNMPv3 USM (SHA/MD5 auth + AES/DES privacy)
- GET, GETNEXT, GETBULK, SET operations
- MIB walk (subtree traversal)
- Standard MIB-II system group validation (1.3.6.1.2.1.1)

**Real implementation upgrade:**
```go
import "github.com/gosnmp/gosnmp"
// Replace Client.Get() with gosnmp.Default.Get(oids)
```

---

## Adding a New Test Case

1. Implement the `common.TestCase` interface:

```go
type MyTest struct { client *rest.Client }

func (t *MyTest) Name()           string          { return "REST: My Custom Test" }
func (t *MyTest) Severity()       common.Severity { return common.SeverityHigh }
func (t *MyTest) Tags()           []string        { return []string{"rest", "custom"} }
func (t *MyTest) Run()            common.TestResult {
    start := time.Now()
    steps := []common.TestStep{}
    assert := common.NewAssert(&steps)

    resp, err := t.client.Get("/my-endpoint")
    if err != nil { /* return error result */ }

    assert.StatusCode(resp.Status, 200)
    assert.Contains("response body", resp.String(), "expected-value")

    return common.TestResult{
        Name: t.Name(), Suite: "REST",
        Status: common.StatusPass, /* or Fail */
        Steps: steps, Duration: time.Since(start),
        StartTime: start, EndTime: time.Now(), Tags: t.Tags(),
    }
}
```

2. Register it in `main.go`:
```go
r.AddSuite("My Suite", "REST", []common.TestCase{
    &MyTest{client: myClient},
})
```

---

## Enhancement Roadmap (Learning Path)

### Stage 1 — Current (Demo)
- ✅ Simulated NETCONF + SNMP engines
- ✅ HTTP REST client with live network
- ✅ Browser page simulation
- ✅ Serial/parallel runner
- ✅ JSON reports

### Stage 2 — Real Protocol Connections
- [ ] Real SSH-based NETCONF (`golang.org/x/crypto/ssh`)
- [ ] Real SNMP (`github.com/gosnmp/gosnmp`)
- [ ] Real browser (`github.com/playwright-community/playwright-go`)
- [ ] YANG model validation for NETCONF payloads

### Stage 3 — Test Infrastructure
- [ ] YAML/JSON config-driven test suites
- [ ] JUnit XML report (CI/CD integration)
- [ ] HTML dashboard report
- [ ] Test retry on flaky assertions
- [ ] Parameterised/table-driven tests

### Stage 4 — Advanced
- [ ] RESTCONF (RFC 8040) — REST wrapper for NETCONF
- [ ] gNMI / gRPC device management
- [ ] BGP/OSPF route verification via SNMP MIBs
- [ ] Chaos/fault injection testing
- [ ] Distributed test execution
- [ ] Metrics export (Prometheus)

---

## Testing Philosophy

Each test case is:
- **Self-contained** — no shared mutable state
- **Step-tracked** — every assertion records a named step
- **Severity-tagged** — CRITICAL / HIGH / MEDIUM / LOW
- **Tag-filterable** — `smoke`, `regression`, `auth`, etc.
- **Parallel-safe** — no global variables in test logic
