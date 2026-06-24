// Package snmpconf provides SNMP protocol testing utilities.
// Supports SNMPv1, SNMPv2c, and SNMPv3 (simulated; real impl would use
// a library like github.com/gosnmp/gosnmp).
package snmpconf

import (
	"fmt"
	"time"

	"github.com/nettest/framework/internal/common"
)

// ─── SNMP Types ───────────────────────────────────────────────────────────────

// Version represents the SNMP version.
type Version int

const (
	Version1  Version = 1
	Version2c Version = 2
	Version3  Version = 3
)

// PDUType represents SNMP PDU types.
type PDUType int

const (
	GetRequest     PDUType = 0xA0
	GetNextRequest PDUType = 0xA1
	GetResponse    PDUType = 0xA2
	SetRequest     PDUType = 0xA3
	GetBulkRequest PDUType = 0xA5
)

// VarBind represents an SNMP variable binding (OID → value).
type VarBind struct {
	OID   string
	Type  string // "Integer", "OctetString", "Counter32", "Gauge32", "TimeTicks", "OID"
	Value interface{}
}

// PDU represents an SNMP Protocol Data Unit.
type PDU struct {
	Type      PDUType
	RequestID int32
	ErrorCode int
	ErrorIdx  int
	VarBinds  []VarBind
}

// Standard OIDs (MIB-II)
const (
	OIDSysDescr    = "1.3.6.1.2.1.1.1.0"
	OIDSysObjectID = "1.3.6.1.2.1.1.2.0"
	OIDSysUpTime   = "1.3.6.1.2.1.1.3.0"
	OIDSysContact  = "1.3.6.1.2.1.1.4.0"
	OIDSysName     = "1.3.6.1.2.1.1.5.0"
	OIDSysLocation = "1.3.6.1.2.1.1.6.0"
	OIDIfNumber    = "1.3.6.1.2.1.2.1.0"
	OIDIfTable     = "1.3.6.1.2.1.2.2"
)

// ─── Simulated SNMP Client ────────────────────────────────────────────────────

// Client simulates an SNMP client.
// Real implementation: use github.com/gosnmp/gosnmp
type Client struct {
	Host      string
	Port      int
	Community string
	Version   Version
	Timeout   time.Duration
	// SNMPv3 fields
	User      string
	AuthProto string
	AuthKey   string
	PrivProto string
	PrivKey   string
}

// NewV2cClient creates an SNMPv2c client.
func NewV2cClient(host string, port int, community string) *Client {
	return &Client{
		Host:      host,
		Port:      port,
		Community: community,
		Version:   Version2c,
		Timeout:   5 * time.Second,
	}
}

// NewV3Client creates an SNMPv3 client.
func NewV3Client(host string, port int, user, authProto, authKey, privProto, privKey string) *Client {
	return &Client{
		Host:      host,
		Port:      port,
		Version:   Version3,
		Timeout:   5 * time.Second,
		User:      user,
		AuthProto: authProto,
		AuthKey:   authKey,
		PrivProto: privProto,
		PrivKey:   privKey,
	}
}

// Get performs an SNMP GET for the given OIDs.
func (c *Client) Get(oids []string) (*PDU, error) {
	if c.Host == "" {
		return nil, fmt.Errorf("snmp: host not configured")
	}
	// Simulate response
	vars := make([]VarBind, 0, len(oids))
	for _, oid := range oids {
		vars = append(vars, simulatedVarBind(oid))
	}
	return &PDU{
		Type:      GetResponse,
		RequestID: 1001,
		VarBinds:  vars,
	}, nil
}

// GetNext performs an SNMP GETNEXT walk starting from oid.
func (c *Client) GetNext(oid string) (*PDU, error) {
	if c.Host == "" {
		return nil, fmt.Errorf("snmp: host not configured")
	}
	return &PDU{
		Type:      GetResponse,
		RequestID: 1002,
		VarBinds:  []VarBind{simulatedVarBind(oid + ".0")},
	}, nil
}

// GetBulk performs an SNMP GETBULK (v2c/v3 only).
func (c *Client) GetBulk(oid string, maxRep int) (*PDU, error) {
	if c.Version < Version2c {
		return nil, fmt.Errorf("snmp: GETBULK requires v2c or v3")
	}
	vars := make([]VarBind, maxRep)
	for i := 0; i < maxRep; i++ {
		vars[i] = simulatedVarBind(fmt.Sprintf("%s.%d", oid, i+1))
	}
	return &PDU{Type: GetResponse, RequestID: 1003, VarBinds: vars}, nil
}

// Set performs an SNMP SET for a single OID.
func (c *Client) Set(oid string, val interface{}) (*PDU, error) {
	if c.Host == "" {
		return nil, fmt.Errorf("snmp: host not configured")
	}
	return &PDU{
		Type:      GetResponse,
		RequestID: 1004,
		VarBinds:  []VarBind{{OID: oid, Value: val}},
	}, nil
}

// Walk performs an SNMP walk under a subtree OID.
func (c *Client) Walk(subtree string) ([]VarBind, error) {
	simulated := []VarBind{
		{OID: subtree + ".1", Type: "OctetString", Value: "eth0"},
		{OID: subtree + ".2", Type: "OctetString", Value: "eth1"},
		{OID: subtree + ".3", Type: "OctetString", Value: "lo"},
	}
	return simulated, nil
}

func simulatedVarBind(oid string) VarBind {
	switch oid {
	case OIDSysDescr, OIDSysDescr + ".0":
		return VarBind{OID: oid, Type: "OctetString", Value: "Linux router 5.15.0 #1 SMP"}
	case OIDSysName, OIDSysName + ".0":
		return VarBind{OID: oid, Type: "OctetString", Value: "core-router-01"}
	case OIDSysUpTime, OIDSysUpTime + ".0":
		return VarBind{OID: oid, Type: "TimeTicks", Value: uint32(123456)}
	case OIDIfNumber, OIDIfNumber + ".0":
		return VarBind{OID: oid, Type: "Integer", Value: 4}
	case OIDSysLocation:
		return VarBind{OID: oid, Type: "OctetString", Value: "Data Center Rack A1"}
	default:
		return VarBind{OID: oid, Type: "OctetString", Value: "simulated-value"}
	}
}

// ─── Test Cases ───────────────────────────────────────────────────────────────

// ReachabilityTest verifies basic SNMP reachability via sysDescr GET.
type ReachabilityTest struct{ client *Client }

func NewReachabilityTest(c *Client) *ReachabilityTest { return &ReachabilityTest{client: c} }
func (t *ReachabilityTest) Name() string              { return "SNMP: Reachability (sysDescr)" }
func (t *ReachabilityTest) Severity() common.Severity { return common.SeverityCritical }
func (t *ReachabilityTest) Tags() []string            { return []string{"snmp", "reachability", "smoke"} }

func (t *ReachabilityTest) Run() common.TestResult {
	start := time.Now()
	steps := []common.TestStep{}
	assert := common.NewAssert(&steps)
	var finalErr error
	check := func(err error) {
		if err != nil && finalErr == nil {
			finalErr = err
		}
	}

	pdu, err := t.client.Get([]string{OIDSysDescr})
	if err != nil {
		return errorResult(t.Name(), "SNMP", t.Severity(), err, start, steps, t.Tags())
	}
	steps = append(steps, common.TestStep{Name: "GET sysDescr", Status: common.StatusPass,
		Detail: fmt.Sprintf("received %d varbind(s)", len(pdu.VarBinds))})

	check(assert.True("PDU has varbinds", len(pdu.VarBinds) > 0,
		fmt.Sprintf("%d varbinds", len(pdu.VarBinds))))
	if len(pdu.VarBinds) > 0 {
		val := fmt.Sprintf("%v", pdu.VarBinds[0].Value)
		check(assert.NotEmpty("sysDescr is non-empty", val))
		steps = append(steps, common.TestStep{Name: "sysDescr value",
			Status: common.StatusPass, Detail: val})
	}
	check(assert.True("no SNMP error", pdu.ErrorCode == 0,
		fmt.Sprintf("errorCode=%d", pdu.ErrorCode)))

	status := common.StatusPass
	if finalErr != nil {
		status = common.StatusFail
	}
	return common.TestResult{
		Name: t.Name(), Suite: "SNMP", Status: status, Severity: t.Severity(),
		Duration: time.Since(start), StartTime: start, EndTime: time.Now(),
		Steps: steps, Error: finalErr, Tags: t.Tags(),
	}
}

// MIBWalkTest walks a MIB subtree and checks minimum number of objects.
type MIBWalkTest struct {
	client  *Client
	subtree string
	minObjs int
}

func NewMIBWalkTest(c *Client, subtree string, minObjs int) *MIBWalkTest {
	return &MIBWalkTest{client: c, subtree: subtree, minObjs: minObjs}
}
func (t *MIBWalkTest) Name() string            { return "SNMP: MIB Walk — " + t.subtree }
func (t *MIBWalkTest) Severity() common.Severity { return common.SeverityMedium }
func (t *MIBWalkTest) Tags() []string          { return []string{"snmp", "mib-walk"} }

func (t *MIBWalkTest) Run() common.TestResult {
	start := time.Now()
	steps := []common.TestStep{}
	assert := common.NewAssert(&steps)
	var finalErr error

	vars, err := t.client.Walk(t.subtree)
	if err != nil {
		return errorResult(t.Name(), "SNMP", t.Severity(), err, start, steps, t.Tags())
	}
	steps = append(steps, common.TestStep{Name: "walk " + t.subtree, Status: common.StatusPass,
		Detail: fmt.Sprintf("retrieved %d objects", len(vars))})

	if err := assert.True(fmt.Sprintf("at least %d objects", t.minObjs),
		len(vars) >= t.minObjs, fmt.Sprintf("got %d", len(vars))); err != nil {
		finalErr = err
	}

	status := common.StatusPass
	if finalErr != nil {
		status = common.StatusFail
	}
	return common.TestResult{
		Name: t.Name(), Suite: "SNMP", Status: status, Severity: t.Severity(),
		Duration: time.Since(start), StartTime: start, EndTime: time.Now(),
		Steps: steps, Error: finalErr, Tags: t.Tags(),
	}
}

// SystemInfoTest retrieves standard MIB-II system objects.
type SystemInfoTest struct{ client *Client }

func NewSystemInfoTest(c *Client) *SystemInfoTest { return &SystemInfoTest{client: c} }
func (t *SystemInfoTest) Name() string            { return "SNMP: System Info (MIB-II system group)" }
func (t *SystemInfoTest) Severity() common.Severity { return common.SeverityHigh }
func (t *SystemInfoTest) Tags() []string          { return []string{"snmp", "system-info", "mib2"} }

func (t *SystemInfoTest) Run() common.TestResult {
	start := time.Now()
	steps := []common.TestStep{}
	assert := common.NewAssert(&steps)
	var finalErr error
	check := func(err error) {
		if err != nil && finalErr == nil {
			finalErr = err
		}
	}

	oids := []string{OIDSysDescr, OIDSysName, OIDSysUpTime, OIDSysLocation}
	pdu, err := t.client.Get(oids)
	if err != nil {
		return errorResult(t.Name(), "SNMP", t.Severity(), err, start, steps, t.Tags())
	}
	check(assert.Equal("varbind count", len(oids), len(pdu.VarBinds)))

	for _, vb := range pdu.VarBinds {
		val := fmt.Sprintf("%v", vb.Value)
		check(assert.NotEmpty("OID "+vb.OID, val))
		steps = append(steps, common.TestStep{Name: vb.OID, Status: common.StatusPass,
			Detail: fmt.Sprintf("type=%s value=%v", vb.Type, vb.Value)})
	}

	status := common.StatusPass
	if finalErr != nil {
		status = common.StatusFail
	}
	return common.TestResult{
		Name: t.Name(), Suite: "SNMP", Status: status, Severity: t.Severity(),
		Duration: time.Since(start), StartTime: start, EndTime: time.Now(),
		Steps: steps, Error: finalErr, Tags: t.Tags(),
	}
}

// GetBulkTest exercises SNMP GETBULK (v2c/v3).
type GetBulkTest struct{ client *Client }

func NewGetBulkTest(c *Client) *GetBulkTest { return &GetBulkTest{client: c} }
func (t *GetBulkTest) Name() string         { return "SNMP: GETBULK (interface table)" }
func (t *GetBulkTest) Severity() common.Severity { return common.SeverityMedium }
func (t *GetBulkTest) Tags() []string       { return []string{"snmp", "getbulk", "ifTable"} }

func (t *GetBulkTest) Run() common.TestResult {
	start := time.Now()
	steps := []common.TestStep{}
	assert := common.NewAssert(&steps)
	var finalErr error

	pdu, err := t.client.GetBulk(OIDIfTable, 5)
	if err != nil {
		return errorResult(t.Name(), "SNMP", t.Severity(), err, start, steps, t.Tags())
	}
	if err := assert.True("GETBULK returns objects", len(pdu.VarBinds) > 0,
		fmt.Sprintf("%d objects", len(pdu.VarBinds))); err != nil {
		finalErr = err
	}
	steps = append(steps, common.TestStep{Name: "GETBULK " + OIDIfTable, Status: common.StatusPass,
		Detail: fmt.Sprintf("retrieved %d entries", len(pdu.VarBinds))})

	status := common.StatusPass
	if finalErr != nil {
		status = common.StatusFail
	}
	return common.TestResult{
		Name: t.Name(), Suite: "SNMP", Status: status, Severity: t.Severity(),
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
