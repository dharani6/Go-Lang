// Package netconf provides NETCONF protocol testing utilities.
// NETCONF (RFC 6241) uses SSH transport + XML RPC operations.
package netconf

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/nettest/framework/internal/common"
)

// ─── XML Structures ───────────────────────────────────────────────────────────

// RPCMessage represents a NETCONF <rpc> message.
type RPCMessage struct {
	XMLName   xml.Name `xml:"rpc"`
	MessageID string   `xml:"message-id,attr"`
	Xmlns     string   `xml:"xmlns,attr,omitempty"`
	Body      interface{}
}

// RPCReply represents a NETCONF <rpc-reply> message.
type RPCReply struct {
	XMLName   xml.Name    `xml:"rpc-reply"`
	MessageID string      `xml:"message-id,attr"`
	OK        *struct{}   `xml:"ok"`
	Data      string      `xml:"data"`
	Errors    []RPCError  `xml:"rpc-error"`
}

// RPCError holds a NETCONF error element.
type RPCError struct {
	Type     string `xml:"error-type"`
	Tag      string `xml:"error-tag"`
	Severity string `xml:"error-severity"`
	Message  string `xml:"error-message"`
}

// Datastore constants.
const (
	DatastoreRunning   = "running"
	DatastoreCandidate = "candidate"
	DatastoreStartup   = "startup"
)

// Operation constants.
const (
	OpGetConfig    = "get-config"
	OpEditConfig   = "edit-config"
	OpCopyConfig   = "copy-config"
	OpDeleteConfig = "delete-config"
	OpLock         = "lock"
	OpUnlock       = "unlock"
	OpGet          = "get"
	OpCloseSession = "close-session"
	OpCommit       = "commit"
)

// ─── Simulated Session ────────────────────────────────────────────────────────

// Session simulates a NETCONF SSH session (in a real implementation,
// this would use golang.org/x/crypto/ssh to open a real session).
type Session struct {
	host       string
	port       int
	user       string
	sessionID  int
	capabilities []string
	msgID      int
	connected  bool
	// In a real impl: sshClient *ssh.Client, stdin io.WriteCloser, stdout *bufio.Reader
}

// NewSession creates a simulated NETCONF session.
func NewSession(host string, port int, user, _ string) *Session {
	return &Session{
		host: host,
		port: port,
		user: user,
		capabilities: []string{
			"urn:ietf:params:netconf:base:1.0",
			"urn:ietf:params:netconf:base:1.1",
			"urn:ietf:params:netconf:capability:writable-running:1.0",
			"urn:ietf:params:netconf:capability:candidate:1.0",
			"urn:ietf:params:netconf:capability:rollback-on-error:1.0",
		},
	}
}

// Connect simulates opening an SSH connection and NETCONF hello exchange.
func (s *Session) Connect() error {
	// Real implementation would:
	//   config := &ssh.ClientConfig{User: s.user, ...}
	//   client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port), config)
	//   s.sshClient = client
	//   Then open a session, send <hello>, receive <hello>
	if s.host == "" {
		return fmt.Errorf("netconf: host is empty")
	}
	s.sessionID = 42 // simulated
	s.connected = true
	return nil
}

// Close simulates closing the NETCONF session.
func (s *Session) Close() error {
	if s.connected {
		_, _ = s.sendRPC(OpCloseSession, "")
		s.connected = false
	}
	return nil
}

// nextMsgID returns the next message-id.
func (s *Session) nextMsgID() string {
	s.msgID++
	return fmt.Sprintf("%d", s.msgID)
}

// sendRPC simulates sending a NETCONF RPC and receiving a reply.
func (s *Session) sendRPC(operation, filter string) (*RPCReply, error) {
	if !s.connected {
		return nil, fmt.Errorf("netconf: session not connected")
	}
	msgID := s.nextMsgID()

	// Build XML (simplified)
	var body string
	switch operation {
	case OpGetConfig:
		body = fmt.Sprintf(`<%s><source><%s/></source>%s</%s>`,
			operation, DatastoreRunning, filter, operation)
	case OpGet:
		body = fmt.Sprintf(`<%s>%s</%s>`, operation, filter, operation)
	case OpCommit:
		body = fmt.Sprintf(`<%s/>`, operation)
	default:
		body = fmt.Sprintf(`<%s/>`, operation)
	}
	_ = fmt.Sprintf(`<?xml version="1.0"?><rpc message-id="%s" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">%s</rpc>`,
		msgID, body)

	// Simulate a successful reply
	return &RPCReply{
		MessageID: msgID,
		OK:        &struct{}{},
		Data:      simulatedXMLData(operation),
	}, nil
}

// GetConfig retrieves configuration from a datastore.
func (s *Session) GetConfig(datastore string) (*RPCReply, error) {
	return s.sendRPC(OpGetConfig, "")
}

// EditConfig modifies configuration in a datastore.
func (s *Session) EditConfig(datastore, config string) (*RPCReply, error) {
	return s.sendRPC(OpEditConfig, config)
}

// Commit commits the candidate datastore to running.
func (s *Session) Commit() (*RPCReply, error) {
	return s.sendRPC(OpCommit, "")
}

// Lock locks a datastore.
func (s *Session) Lock(datastore string) (*RPCReply, error) {
	return s.sendRPC(OpLock, "")
}

// Unlock unlocks a datastore.
func (s *Session) Unlock(datastore string) (*RPCReply, error) {
	return s.sendRPC(OpUnlock, "")
}

// HasCapability checks whether the remote peer advertised a capability.
func (s *Session) HasCapability(cap string) bool {
	for _, c := range s.capabilities {
		if strings.Contains(c, cap) {
			return true
		}
	}
	return false
}

func simulatedXMLData(op string) string {
	switch op {
	case OpGetConfig:
		return `<configuration><interfaces><interface><name>eth0</name><type>ethernetCsmacd</type><enabled>true</enabled></interface></interfaces></configuration>`
	default:
		return ""
	}
}

// ─── Test Cases ───────────────────────────────────────────────────────────────

// ConnectivityTest verifies NETCONF session establishment.
type ConnectivityTest struct {
	session *Session
}

func NewConnectivityTest(s *Session) *ConnectivityTest { return &ConnectivityTest{session: s} }
func (t *ConnectivityTest) Name() string               { return "NETCONF: Session Connectivity" }
func (t *ConnectivityTest) Severity() common.Severity  { return common.SeverityCritical }
func (t *ConnectivityTest) Tags() []string             { return []string{"netconf", "connectivity", "smoke"} }

func (t *ConnectivityTest) Run() common.TestResult {
	start := time.Now()
	steps := []common.TestStep{}
	assert := common.NewAssert(&steps)

	err := t.session.Connect()
	if err != nil {
		return common.TestResult{
			Name: t.Name(), Suite: "NETCONF", Status: common.StatusError,
			Severity: t.Severity(), Error: err, StartTime: start, EndTime: time.Now(), Steps: steps,
		}
	}
	steps = append(steps, common.TestStep{Name: "SSH connect + hello exchange",
		Status: common.StatusPass, Detail: fmt.Sprintf("sessionID=%d", t.session.sessionID)})

	var finalErr error
	if err := assert.True("has base capability",
		t.session.HasCapability("base"), "urn:ietf:params:netconf:base"); err != nil {
		finalErr = err
	}
	if err := assert.True("session ID assigned",
		t.session.sessionID > 0, fmt.Sprintf("sessionID=%d", t.session.sessionID)); err != nil && finalErr == nil {
		finalErr = err
	}

	status := common.StatusPass
	if finalErr != nil {
		status = common.StatusFail
	}
	return common.TestResult{
		Name: t.Name(), Suite: "NETCONF", Status: status, Severity: t.Severity(),
		Duration: time.Since(start), StartTime: start, EndTime: time.Now(),
		Steps: steps, Error: finalErr, Tags: t.Tags(),
	}
}

// GetConfigTest verifies <get-config> operation.
type GetConfigTest struct{ session *Session }

func NewGetConfigTest(s *Session) *GetConfigTest { return &GetConfigTest{session: s} }
func (t *GetConfigTest) Name() string            { return "NETCONF: get-config (running)" }
func (t *GetConfigTest) Severity() common.Severity { return common.SeverityHigh }
func (t *GetConfigTest) Tags() []string          { return []string{"netconf", "get-config"} }

func (t *GetConfigTest) Run() common.TestResult {
	start := time.Now()
	steps := []common.TestStep{}
	assert := common.NewAssert(&steps)
	var finalErr error
	check := func(err error) {
		if err != nil && finalErr == nil {
			finalErr = err
		}
	}

	if err := t.session.Connect(); err != nil {
		return errorResult(t.Name(), "NETCONF", t.Severity(), err, start, steps, t.Tags())
	}
	reply, err := t.session.GetConfig(DatastoreRunning)
	if err != nil {
		return errorResult(t.Name(), "NETCONF", t.Severity(), err, start, steps, t.Tags())
	}
	steps = append(steps, common.TestStep{Name: "send get-config RPC",
		Status: common.StatusPass, Detail: "RPC sent successfully"})

	check(assert.True("reply has OK", reply.OK != nil, "rpc-reply contains <ok/>"))
	check(assert.True("no RPC errors", len(reply.Errors) == 0,
		fmt.Sprintf("%d errors returned", len(reply.Errors))))
	check(assert.Contains("data contains interface", reply.Data, "interface"))

	status := common.StatusPass
	if finalErr != nil {
		status = common.StatusFail
	}
	return common.TestResult{
		Name: t.Name(), Suite: "NETCONF", Status: status, Severity: t.Severity(),
		Duration: time.Since(start), StartTime: start, EndTime: time.Now(),
		Steps: steps, Error: finalErr, Tags: t.Tags(),
	}
}

// EditConfigTest verifies <edit-config> + <commit> operations.
type EditConfigTest struct{ session *Session }

func NewEditConfigTest(s *Session) *EditConfigTest { return &EditConfigTest{session: s} }
func (t *EditConfigTest) Name() string             { return "NETCONF: edit-config + commit" }
func (t *EditConfigTest) Severity() common.Severity { return common.SeverityHigh }
func (t *EditConfigTest) Tags() []string           { return []string{"netconf", "edit-config", "commit"} }

func (t *EditConfigTest) Run() common.TestResult {
	start := time.Now()
	steps := []common.TestStep{}
	assert := common.NewAssert(&steps)
	var finalErr error
	check := func(err error) {
		if err != nil && finalErr == nil {
			finalErr = err
		}
	}

	if err := t.session.Connect(); err != nil {
		return errorResult(t.Name(), "NETCONF", t.Severity(), err, start, steps, t.Tags())
	}

	cfg := `<config><interfaces><interface><name>lo0</name><enabled>true</enabled></interface></interfaces></config>`

	// Lock → Edit → Commit → Unlock
	lockReply, _ := t.session.Lock(DatastoreCandidate)
	check(assert.True("lock candidate", lockReply != nil && lockReply.OK != nil, "lock OK"))

	editReply, _ := t.session.EditConfig(DatastoreCandidate, cfg)
	check(assert.True("edit-config OK", editReply != nil && editReply.OK != nil, "edit OK"))

	commitReply, _ := t.session.Commit()
	check(assert.True("commit OK", commitReply != nil && commitReply.OK != nil, "commit OK"))

	unlockReply, _ := t.session.Unlock(DatastoreCandidate)
	check(assert.True("unlock candidate", unlockReply != nil && unlockReply.OK != nil, "unlock OK"))

	steps = append(steps, common.TestStep{Name: "lock → edit → commit → unlock",
		Status: common.StatusPass, Detail: "full edit-config lifecycle completed"})

	status := common.StatusPass
	if finalErr != nil {
		status = common.StatusFail
	}
	return common.TestResult{
		Name: t.Name(), Suite: "NETCONF", Status: status, Severity: t.Severity(),
		Duration: time.Since(start), StartTime: start, EndTime: time.Now(),
		Steps: steps, Error: finalErr, Tags: t.Tags(),
	}
}

// CapabilityTest verifies that required capabilities are advertised.
type CapabilityTest struct {
	session  *Session
	required []string
}

func NewCapabilityTest(s *Session, caps []string) *CapabilityTest {
	return &CapabilityTest{session: s, required: caps}
}
func (t *CapabilityTest) Name() string            { return "NETCONF: Capability Negotiation" }
func (t *CapabilityTest) Severity() common.Severity { return common.SeverityMedium }
func (t *CapabilityTest) Tags() []string          { return []string{"netconf", "capability"} }

func (t *CapabilityTest) Run() common.TestResult {
	start := time.Now()
	steps := []common.TestStep{}
	assert := common.NewAssert(&steps)
	var finalErr error

	if err := t.session.Connect(); err != nil {
		return errorResult(t.Name(), "NETCONF", t.Severity(), err, start, steps, t.Tags())
	}
	for _, cap := range t.required {
		if err := assert.True("capability: "+cap,
			t.session.HasCapability(cap), "advertised by remote"); err != nil && finalErr == nil {
			finalErr = err
		}
	}

	status := common.StatusPass
	if finalErr != nil {
		status = common.StatusFail
	}
	return common.TestResult{
		Name: t.Name(), Suite: "NETCONF", Status: status, Severity: t.Severity(),
		Duration: time.Since(start), StartTime: start, EndTime: time.Now(),
		Steps: steps, Error: finalErr, Tags: t.Tags(),
	}
}

func errorResult(name, suite string, sev common.Severity, err error,
	start time.Time, steps []common.TestStep, tags []string) common.TestResult {
	return common.TestResult{
		Name: name, Suite: suite, Status: common.StatusError, Severity: sev,
		Error: err, StartTime: start, EndTime: time.Now(), Steps: steps, Tags: tags,
	}
}
