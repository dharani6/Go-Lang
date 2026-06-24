// Package runner provides the core test execution engine and reporting.
package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/nettest/framework/internal/common"
)

// ─── Runner ───────────────────────────────────────────────────────────────────

// Options configures the test runner.
type Options struct {
	Parallel     bool
	MaxWorkers   int
	Tags         []string // filter: run only tests with these tags
	FailFast     bool
	Verbose      bool
	OutputFormat string // "console", "json", "junit"
	ReportPath   string
}

// Runner executes a collection of test suites.
type Runner struct {
	suites  []SuiteEntry
	opts    Options
	results []common.SuiteResult
	mu      sync.Mutex
}

// SuiteEntry pairs a suite name + protocol label with its test cases.
type SuiteEntry struct {
	Name     string
	Protocol string
	Cases    []common.TestCase
}

// New creates a Runner.
func New(opts Options) *Runner {
	if opts.MaxWorkers <= 0 {
		opts.MaxWorkers = 4
	}
	return &Runner{opts: opts}
}

// AddSuite registers a test suite.
func (r *Runner) AddSuite(name, protocol string, cases []common.TestCase) {
	r.suites = append(r.suites, SuiteEntry{Name: name, Protocol: protocol, Cases: cases})
}

// Run executes all registered suites and returns the aggregated results.
func (r *Runner) Run() []common.SuiteResult {
	fmt.Println(banner())
	for _, suite := range r.suites {
		result := r.runSuite(suite)
		r.mu.Lock()
		r.results = append(r.results, result)
		r.mu.Unlock()
	}
	r.printSummary()
	if r.opts.ReportPath != "" {
		_ = r.writeReport()
	}
	return r.results
}

func (r *Runner) runSuite(entry SuiteEntry) common.SuiteResult {
	printSuiteHeader(entry.Name, entry.Protocol)
	suite := common.SuiteResult{
		Name:      entry.Name,
		Protocol:  entry.Protocol,
		StartTime: time.Now(),
		Tags:      entry.Cases[0].Tags(), // rough
	}

	cases := r.filterCases(entry.Cases)
	if r.opts.Parallel {
		suite.Results = r.runParallel(cases)
	} else {
		suite.Results = r.runSerial(cases)
	}
	suite.EndTime = time.Now()

	pass, fail, skip, errCount := suite.Stats()
	fmt.Printf("   Suite finished in %v — PASS:%d FAIL:%d SKIP:%d ERROR:%d (%.0f%% pass rate)\n\n",
		suite.Duration().Round(time.Millisecond), pass, fail, skip, errCount, suite.PassRate())

	return suite
}

func (r *Runner) filterCases(cases []common.TestCase) []common.TestCase {
	if len(r.opts.Tags) == 0 {
		return cases
	}
	var filtered []common.TestCase
	for _, c := range cases {
		for _, tag := range c.Tags() {
			for _, want := range r.opts.Tags {
				if tag == want {
					filtered = append(filtered, c)
					goto next
				}
			}
		}
	next:
	}
	return filtered
}

func (r *Runner) runSerial(cases []common.TestCase) []common.TestResult {
	results := make([]common.TestResult, 0, len(cases))
	for _, tc := range cases {
		result := r.runOne(tc)
		results = append(results, result)
		if r.opts.FailFast && result.Status == common.StatusFail {
			break
		}
	}
	return results
}

func (r *Runner) runParallel(cases []common.TestCase) []common.TestResult {
	type indexed struct {
		idx int
		res common.TestResult
	}
	ch := make(chan indexed, len(cases))
	sem := make(chan struct{}, r.opts.MaxWorkers)

	var wg sync.WaitGroup
	for i, tc := range cases {
		wg.Add(1)
		go func(idx int, c common.TestCase) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			ch <- indexed{idx: idx, res: r.runOne(c)}
		}(i, tc)
	}
	wg.Wait()
	close(ch)

	indexed_results := make([]indexed, 0, len(cases))
	for ir := range ch {
		indexed_results = append(indexed_results, ir)
	}
	sort.Slice(indexed_results, func(i, j int) bool {
		return indexed_results[i].idx < indexed_results[j].idx
	})
	results := make([]common.TestResult, len(indexed_results))
	for i, ir := range indexed_results {
		results[i] = ir.res
	}
	return results
}

func (r *Runner) runOne(tc common.TestCase) common.TestResult {
	start := time.Now()
	result := tc.Run()
	result.Duration = time.Since(start)

	icon := statusIcon(result.Status)
	sev := string(result.Severity)
	fmt.Printf("   %s %-55s [%s] %v\n",
		icon, truncate(result.Name, 55), padRight(sev, 8), result.Duration.Round(time.Millisecond))

	if r.opts.Verbose && len(result.Steps) > 0 {
		for _, step := range result.Steps {
			fmt.Printf("      %s %s — %s\n", statusIcon(step.Status), step.Name, step.Detail)
		}
	}
	if result.Error != nil {
		fmt.Printf("      ⚠  %v\n", result.Error)
	}
	return result
}

func (r *Runner) printSummary() {
	var totalPass, totalFail, totalSkip, totalErr, totalTests int
	var totalDur time.Duration
	for _, s := range r.results {
		p, f, sk, e := s.Stats()
		totalPass += p
		totalFail += f
		totalSkip += sk
		totalErr += e
		totalTests += len(s.Results)
		totalDur += s.Duration()
	}

	fmt.Println(strings.Repeat("─", 68))
	fmt.Println("  FINAL SUMMARY")
	fmt.Println(strings.Repeat("─", 68))
	fmt.Printf("  %-20s %d\n", "Total Suites:", len(r.results))
	fmt.Printf("  %-20s %d\n", "Total Tests:", totalTests)
	fmt.Printf("  %-20s %d\n", "✅ Passed:", totalPass)
	fmt.Printf("  %-20s %d\n", "❌ Failed:", totalFail)
	fmt.Printf("  %-20s %d\n", "⏭  Skipped:", totalSkip)
	fmt.Printf("  %-20s %d\n", "💥 Errors:", totalErr)
	if totalTests > 0 {
		rate := float64(totalPass) / float64(totalTests) * 100
		fmt.Printf("  %-20s %.1f%%\n", "Pass Rate:", rate)
	}
	fmt.Printf("  %-20s %v\n", "Total Duration:", totalDur.Round(time.Millisecond))
	fmt.Println(strings.Repeat("─", 68))
}

func (r *Runner) writeReport() error {
	type Report struct {
		GeneratedAt time.Time            `json:"generated_at"`
		Suites      []common.SuiteResult `json:"suites"`
	}
	rpt := Report{GeneratedAt: time.Now(), Suites: r.results}
	data, err := json.MarshalIndent(rpt, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.opts.ReportPath, data, 0644)
}

// ─── Formatting helpers ───────────────────────────────────────────────────────

func statusIcon(s common.TestStatus) string {
	switch s {
	case common.StatusPass:
		return "✅"
	case common.StatusFail:
		return "❌"
	case common.StatusSkip:
		return "⏭ "
	case common.StatusError:
		return "💥"
	default:
		return "⬜"
	}
}

func banner() string {
	return `
╔══════════════════════════════════════════════════════════════════╗
║          Go Network Test Framework — Demo Run                    ║
║    Protocols: UI · REST API · NETCONF · SNMP                     ║
╚══════════════════════════════════════════════════════════════════╝
`
}

func printSuiteHeader(name, protocol string) {
	fmt.Printf("\n%s\n  ▶  %s  [%s]\n%s\n",
		strings.Repeat("─", 68), name, protocol, strings.Repeat("─", 68))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}
