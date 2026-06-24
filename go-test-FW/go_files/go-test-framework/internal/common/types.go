// Package common provides shared types and utilities for the testing framework.
package common

import (
	"fmt"
	"time"
)

// TestStatus represents the outcome of a test.
type TestStatus string

const (
	StatusPass    TestStatus = "PASS"
	StatusFail    TestStatus = "FAIL"
	StatusSkip    TestStatus = "SKIP"
	StatusError   TestStatus = "ERROR"
	StatusPending TestStatus = "PENDING"
)

// Severity indicates the criticality of a test.
type Severity string

const (
	SeverityCritical Severity = "CRITICAL"
	SeverityHigh     Severity = "HIGH"
	SeverityMedium   Severity = "MEDIUM"
	SeverityLow      Severity = "LOW"
)

// TestResult holds the outcome of a single test case.
type TestResult struct {
	Name        string
	Suite       string
	Status      TestStatus
	Severity    Severity
	Duration    time.Duration
	Message     string
	Error       error
	StartTime   time.Time
	EndTime     time.Time
	Tags        []string
	Steps       []TestStep
	Attachments []Attachment
}

// TestStep represents a single step within a test case.
type TestStep struct {
	Name     string
	Status   TestStatus
	Duration time.Duration
	Detail   string
}

// Attachment holds extra diagnostic data (logs, payloads, etc.).
type Attachment struct {
	Name    string
	MIMEType string
	Content []byte
}

// SuiteResult aggregates results for an entire test suite.
type SuiteResult struct {
	Name      string
	Protocol  string
	StartTime time.Time
	EndTime   time.Time
	Results   []TestResult
	Tags      []string
}

// Stats computes summary statistics for a SuiteResult.
func (s *SuiteResult) Stats() (pass, fail, skip, errCount int) {
	for _, r := range s.Results {
		switch r.Status {
		case StatusPass:
			pass++
		case StatusFail:
			fail++
		case StatusSkip:
			skip++
		case StatusError:
			errCount++
		}
	}
	return
}

// Duration returns total wall-clock duration of the suite.
func (s *SuiteResult) Duration() time.Duration {
	return s.EndTime.Sub(s.StartTime)
}

// PassRate returns percentage of passing tests.
func (s *SuiteResult) PassRate() float64 {
	if len(s.Results) == 0 {
		return 0
	}
	pass, _, _, _ := s.Stats()
	return float64(pass) / float64(len(s.Results)) * 100
}

// TestCase is the interface every test case must implement.
type TestCase interface {
	Name() string
	Severity() Severity
	Tags() []string
	Run() TestResult
}

// Config holds connection parameters for each protocol.
type Config struct {
	// REST / UI
	BaseURL  string
	Username string
	Password string
	Timeout  time.Duration

	// NETCONF
	NetconfHost     string
	NetconfPort     int
	NetconfUser     string
	NetconfPassword string
	NetconfKeyFile  string

	// SNMP
	SNMPHost      string
	SNMPPort      int
	SNMPCommunity string
	SNMPVersion   string // "v1", "v2c", "v3"
	SNMPUser      string // v3
	SNMPAuthProto string // MD5, SHA
	SNMPAuthKey   string
	SNMPPrivProto string // DES, AES
	SNMPPrivKey   string

	// UI / Browser
	BrowserType    string // "chrome", "firefox", "headless"
	ScreenshotDir  string
	UIBaseURL      string
}

// ErrTestFailed is a sentinel error for failed assertions.
type ErrTestFailed struct {
	Assertion string
	Expected  interface{}
	Got       interface{}
}

func (e *ErrTestFailed) Error() string {
	return fmt.Sprintf("assertion failed [%s]: expected=%v got=%v",
		e.Assertion, e.Expected, e.Got)
}
