// Copyright Built On Envoy
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

// Package waf implements the WAF HTTP filter plugin using Coraza.
package waf

import (
	"net"
	"strconv"

	"github.com/corazawaf/coraza/v3"
	ctypes "github.com/corazawaf/coraza/v3/types"
	"github.com/envoyproxy/envoy/source/extensions/dynamic_modules/sdk/go/shared"

	"github.com/tetratelabs/built-on-envoy/extensions/composer/pkg"
	waf "github.com/tetratelabs/built-on-envoy/extensions/composer/waf/coraza"
	"github.com/tetratelabs/built-on-envoy/extensions/composer/waf/logger"
)

const (
	defaultMetadataNamespace = "io.builtonenvoy.waf"
	metadataKeyBlockRule     = "block_rule"
	metadataKeyBlockPhase    = "block_phase"
)

type wafPluginFactory struct {
	shared.EmptyHttpFilterFactory
	config     coraza.WAF
	mode       waf.WAFMode
	headerMode waf.HeaderMode
	metrics    *metrics
}

type perRouteWafPluginConfig struct {
	config     coraza.WAF
	mode       waf.WAFMode
	headerMode waf.HeaderMode
}

func (f *wafPluginFactory) Create(handle shared.HttpFilterHandle) shared.HttpFilter {
	config := f.config
	mode := f.mode
	headerMode := f.headerMode

	// Check for per-route config and override if present.
	perRouteWafPluginConfig := pkg.GetMostSpecificConfig[*perRouteWafPluginConfig](handle)
	if perRouteWafPluginConfig != nil {
		config = perRouteWafPluginConfig.config
		mode = perRouteWafPluginConfig.mode
		headerMode = perRouteWafPluginConfig.headerMode
	}

	if config == nil {
		handle.Log(shared.LogLevelInfo, "waf: no config available and use empty filter")
		return &shared.EmptyHttpFilter{}
	}

	return &wafPlugin{
		logger:            logger.GetLogger(),
		handle:            handle,
		config:            config,
		mode:              mode,
		headerMode:        headerMode,
		metrics:           f.metrics,
		metadataNamespace: defaultMetadataNamespace,
	}
}

type wafPluginConfigFactory struct {
	shared.EmptyHttpFilterConfigFactory
}

func (f *wafPluginConfigFactory) Create(
	handle shared.HttpFilterConfigHandle,
	unparsedConfig []byte,
) (shared.HttpFilterFactory, error) {
	var wafConfig coraza.WAF
	var mode waf.WAFMode
	var headerMode waf.HeaderMode
	var err error

	if len(unparsedConfig) > 0 {
		wafConfig, mode, headerMode, err = waf.NewWAFConfigFromBytes(unparsedConfig, logger.GetLogger())
	}

	if err != nil {
		return nil, err
	}

	if wafConfig == nil {
		handle.Log(shared.LogLevelInfo, "waf: empty filter config")
	}

	return &wafPluginFactory{
		config:     wafConfig,
		mode:       mode,
		headerMode: headerMode,
		metrics:    newMetrics(handle),
	}, nil
}

func (f *wafPluginConfigFactory) CreatePerRoute(unparsedConfig []byte) (any, error) {
	wafConfig, mode, headerMode, err := waf.NewWAFConfigFromBytes(unparsedConfig, logger.GetLogger())
	if err != nil {
		return nil, err
	}
	return &perRouteWafPluginConfig{
		config:     wafConfig,
		mode:       mode,
		headerMode: headerMode,
	}, nil
}

// The plugin struct that implements the actual logic.
type wafPlugin struct {
	shared.EmptyHttpFilter
	logger            *logger.Logger
	handle            shared.HttpFilterHandle
	config            coraza.WAF
	mode              waf.WAFMode
	headerMode        waf.HeaderMode
	metrics           *metrics
	metadataNamespace string

	txContext             ctypes.Transaction
	protocol              string
	isUpgrade             bool
	isSSE                 bool
	requestBodyProcessed  bool
	responseBodyProcessed bool
	authority             string
}

func (p *wafPlugin) getSourceAddress() (string, int) {
	addressAttr, _ := p.handle.GetAttributeString(shared.AttributeIDSourceAddress)
	address := addressAttr.ToUnsafeString()
	if address == "" {
		p.handle.Log(shared.LogLevelDebug, "No source.address attribute")
		// Use a default value if the attribute is not set.
		return "127.0.0.1", 80
	}

	targetIP, targetPortStr, err := net.SplitHostPort(address)
	if err != nil {
		p.handle.Log(shared.LogLevelDebug, "Invalid source.address attribute format")
		return "127.0.0.1", 80
	}
	targetPort, err := strconv.Atoi(targetPortStr)
	if err != nil {
		p.handle.Log(shared.LogLevelDebug, "Invalid source.address port")
		return "127.0.0.1", 80
	}
	return targetIP, targetPort
}

func (p *wafPlugin) getRequestProtocol() string {
	protocolAttr, _ := p.handle.GetAttributeString(shared.AttributeIDRequestProtocol)
	protocol := protocolAttr.ToString()
	if protocol == "" {
		p.handle.Log(shared.LogLevelDebug, "No request.protocol attribute")
		return "HTTP/1.1"
	}
	return protocol
}

// minimalRequestHeaders is the set of security-relevant headers forwarded to Coraza
// when header_mode is MINIMAL. Host is always forwarded separately.
var minimalRequestHeaders = []string{
	"user-agent",
	"accept",
	"content-type",
	"content-length",
	"cookie",
	"authorization",
	"referer",
	"origin",
	"x-forwarded-for",
	"x-real-ip",
}

// containsIgnoreASCII reports whether substr appears in s using ASCII case-insensitive
// comparison. It does not allocate.
func containsIgnoreASCII(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
outer:
	for i := 0; i <= len(s)-len(substr); i++ {
		for j := 0; j < len(substr); j++ {
			sc := s[i+j]
			tc := substr[j]
			if sc == tc {
				continue
			}
			// Fold ASCII letters to lowercase and compare.
			if sc|0x20 != tc|0x20 || sc|0x20 < 'a' || sc|0x20 > 'z' {
				continue outer
			}
		}
		return true
	}
	return false
}

func (p *wafPlugin) mayInitializeTransaction(headers shared.HeaderMap) {
	if p.txContext != nil {
		return
	}
	id := headers.GetOne("x-request-id").ToString()
	p.txContext = p.config.NewTransactionWithID(id)
	p.isUpgrade = p.checkUpgrade(headers)
}

func (p *wafPlugin) checkUpgrade(headers shared.HeaderMap) bool {
	connectionHeader := headers.GetOne("connection").ToUnsafeString()
	upgrade := headers.GetOne("upgrade").ToUnsafeString()
	return containsIgnoreASCII(connectionHeader, "upgrade") && upgrade != ""
}

func (p *wafPlugin) checkSSE(headers shared.HeaderMap) bool {
	return containsIgnoreASCII(headers.GetOne("content-type").ToUnsafeString(), "text/event-stream")
}

// noBodyExpected returns true when the request is guaranteed to carry no body,
// allowing the request-body phase to be run immediately without buffering.
func noBodyExpected(method string, headers shared.HeaderMap) bool {
	if method == "GET" || method == "HEAD" {
		return true
	}
	cl := headers.GetOne("content-length").ToUnsafeString()
	return cl == "0"
}

func (p *wafPlugin) OnRequestHeaders(headers shared.HeaderMap, endOfStream bool) shared.HeadersStatus {
	// Save for later use in response processing.
	host := headers.GetOne(":authority").ToString()
	p.protocol = p.getRequestProtocol()
	if authority, _, err := net.SplitHostPort(host); err == nil {
		p.authority = authority
	} else {
		p.authority = host
	}
	p.mayInitializeTransaction(headers)

	if p.txContext.IsRuleEngineOff() || (p.mode == waf.ModeResponseOnly) {
		return shared.HeadersStatusContinue
	}

	srcIP, srcPort := p.getSourceAddress()
	// Destination is not known in this context. Use placeholders.
	dstIP, dstPort := "127.0.0.1", 80

	scheme := headers.GetOne(":scheme").ToString()
	if scheme == "" {
		scheme = "http"
	}
	path := headers.GetOne(":path").ToString()
	method := headers.GetOne(":method").ToString()
	uri := scheme + "://" + host + path

	// CRS rules tend to expect Host even with HTTP/2
	p.txContext.AddRequestHeader("Host", host)
	p.txContext.SetServerName(host)
	if p.headerMode == waf.HeaderModeMinimal {
		for _, name := range minimalRequestHeaders {
			if v := headers.GetOne(name).ToString(); v != "" {
				p.txContext.AddRequestHeader(name, v)
			}
		}
	} else {
		headerMap := headers.GetAll()
		for _, header := range headerMap {
			p.txContext.AddRequestHeader(header[0].ToString(), header[1].ToString())
		}
	}

	p.txContext.ProcessConnection(srcIP, srcPort, dstIP, dstPort)
	p.txContext.ProcessURI(uri, method, p.protocol)
	interruption := p.txContext.ProcessRequestHeaders()
	if interruption != nil {
		p.blockRequest(interruption, ctypes.PhaseRequestHeaders, "waf_request_headers_blocked")
		return shared.HeadersStatusStop
	}

	// If endOfStream is true, if we won't buffer the body (upgrade), or if no body is
	// expected (GET/HEAD or Content-Length: 0), run phase-2 rules immediately with an
	// empty body. This avoids HeadersStatusStop for the common no-body case.
	if endOfStream || p.isUpgrade || noBodyExpected(method, headers) {
		if !p.handleRequestBody() {
			return shared.HeadersStatusStop
		}
		return shared.HeadersStatusContinue
	}
	return shared.HeadersStatusStop
}

func (p *wafPlugin) OnRequestBody(body shared.BodyBuffer, endOfStream bool) shared.BodyStatus {
	if p.txContext.IsRuleEngineOff() || (p.mode == waf.ModeResponseOnly ||
		p.requestBodyProcessed) {
		return shared.BodyStatusContinue
	}

	if !p.txContext.IsRequestBodyAccessible() {
		if !p.handleRequestBody() {
			return shared.BodyStatusStopNoBuffer
		}
		return shared.BodyStatusContinue
	}

	// Write the body chunks to the WAF and handle possible interruptions.
	if !p.writeRequestBody(body) {
		return shared.BodyStatusStopNoBuffer
	}

	// If endOfStream is true, process the body now.
	if endOfStream {
		if !p.handleRequestBody() {
			return shared.BodyStatusStopNoBuffer
		}
		return shared.BodyStatusContinue
	}

	if p.isUpgrade {
		// In case of upgrade, we cannot buffer the body anyway.
		return shared.BodyStatusContinue
	}

	// In other cases, buffer the body until end of stream.
	return shared.BodyStatusStopAndBuffer
}

func (p *wafPlugin) OnRequestTrailers(_ shared.HeaderMap) shared.TrailersStatus {
	if p.txContext.IsRuleEngineOff() || (p.mode == waf.ModeResponseOnly ||
		p.requestBodyProcessed) {
		return shared.TrailersStatusContinue
	}

	// Has trailers means we never had endOfStream in OnRequestBody. If the body
	// is not yet processed, process it now.
	if !p.handleRequestBody() {
		return shared.TrailersStatusStop
	}
	return shared.TrailersStatusContinue
}

func (p *wafPlugin) OnResponseHeaders(headers shared.HeaderMap, endOfStream bool) shared.HeadersStatus {
	p.mayInitializeTransaction(p.handle.RequestHeaders())
	p.isSSE = p.checkSSE(headers)

	if p.txContext.IsRuleEngineOff() || (p.mode == waf.ModeRequestOnly) {
		return shared.HeadersStatusContinue
	}

	for _, header := range headers.GetAll() {
		p.txContext.AddResponseHeader(header[0].ToString(), header[1].ToString())
	}
	if p.protocol == "" {
		p.protocol = p.getRequestProtocol()
	}
	codeStr := headers.GetOne(":status").ToUnsafeString()
	if codeStr == "" {
		codeStr = "500"
	}
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		p.handle.Log(shared.LogLevelInfo, "Invalid response status code: %s", codeStr)
		p.blockRequest(nil, ctypes.PhaseResponseHeaders, "waf_internal_error")
		return shared.HeadersStatusStop
	}

	interruption := p.txContext.ProcessResponseHeaders(code, p.protocol)
	if interruption != nil {
		p.blockRequest(interruption, ctypes.PhaseResponseHeaders, "waf_response_headers_blocked")
		return shared.HeadersStatusStop
	}

	// If endOfStream is true or if we are not going to buffer the response body anyway (upgrade or SSE),
	// run handleResponseBody to enforce phase 4 rules. Doing so, response variables already populated
	// are checked, before the headers are sent back.
	if endOfStream || p.isUpgrade || p.isSSE {
		if !p.handleResponseBody() {
			return shared.HeadersStatusStop
		}
		return shared.HeadersStatusContinue
	}
	return shared.HeadersStatusStop
}

func (p *wafPlugin) OnResponseBody(body shared.BodyBuffer, endOfStream bool) shared.BodyStatus {
	if p.txContext.IsRuleEngineOff() || (p.mode == waf.ModeRequestOnly) ||
		p.responseBodyProcessed {
		return shared.BodyStatusContinue
	}

	if !p.txContext.IsResponseBodyAccessible() || !p.txContext.IsResponseBodyProcessable() {
		if !p.handleResponseBody() {
			return shared.BodyStatusStopNoBuffer
		}
		return shared.BodyStatusContinue
	}

	// Write the body chunks to the WAF and handle possible interruptions.
	if !p.writeResponseBody(body) {
		return shared.BodyStatusStopNoBuffer
	}

	// If endOfStream is true, process the body now.
	if endOfStream {
		if !p.handleResponseBody() {
			return shared.BodyStatusStopNoBuffer
		}
		return shared.BodyStatusContinue
	}

	if p.isUpgrade || p.isSSE {
		// In case of upgrade or SSE, we cannot buffer the body anyway.
		return shared.BodyStatusContinue
	}

	// In other cases, buffer the body until end of stream.
	return shared.BodyStatusStopAndBuffer
}

func (p *wafPlugin) OnResponseTrailers(_ shared.HeaderMap) shared.TrailersStatus {
	if p.txContext.IsRuleEngineOff() || (p.mode == waf.ModeRequestOnly) ||
		p.responseBodyProcessed {
		return shared.TrailersStatusContinue
	}

	// Has trailers means we never had endOfStream in OnResponseBody. If the body
	// is not yet processed, process it now.
	if !p.handleResponseBody() {
		return shared.TrailersStatusStop
	}

	return shared.TrailersStatusContinue
}

func (p *wafPlugin) OnStreamComplete() {
	if p.txContext != nil {
		if !p.txContext.IsRuleEngineOff() {
			p.metrics.RecordTx(p.handle)
		}
		p.txContext.ProcessLogging()
		err := p.txContext.Close()
		if err != nil {
			p.handle.Log(shared.LogLevelDebug, "Failed to close WAF transaction: %v", err.Error())
		}
	}
}

func (p *wafPlugin) writeRequestBody(body shared.BodyBuffer) bool {
	if body == nil {
		return true
	}
	for _, chunk := range body.GetChunks() {
		interruption, _, err := p.txContext.WriteRequestBody(chunk.ToUnsafeBytes())
		if err != nil {
			p.handle.Log(shared.LogLevelInfo,
				"Failed to write partial request body to WAF: %v", err.Error())
			p.blockRequest(nil, ctypes.PhaseRequestBody, "waf_internal_error")
			return false
		}
		// Write*Body triggers Process*Body if the bodylimit (Sec*BodyLimit) is reached.
		if interruption != nil {
			p.blockRequest(interruption, ctypes.PhaseRequestBody, "waf_request_body_overflow")
			return false
		}
	}
	return true
}

func (p *wafPlugin) handleRequestBody() bool {
	p.requestBodyProcessed = true

	interruption, err := p.txContext.ProcessRequestBody()
	if err != nil {
		p.handle.Log(shared.LogLevelInfo, "Failed to process request body in WAF: %v", err.Error())
		p.blockRequest(nil, ctypes.PhaseRequestBody, "waf_internal_error")
		return false
	}
	if interruption != nil {
		p.blockRequest(interruption, ctypes.PhaseRequestBody, "waf_request_body_blocked")
		return false
	}
	return true
}

func (p *wafPlugin) writeResponseBody(body shared.BodyBuffer) bool {
	if body == nil {
		return true
	}
	for _, chunk := range body.GetChunks() {
		interruption, _, err := p.txContext.WriteResponseBody(chunk.ToUnsafeBytes())
		if err != nil {
			p.handle.Log(shared.LogLevelInfo, "Failed to write partial response body to WAF: %v", err.Error())
			p.blockRequest(nil, ctypes.PhaseResponseBody, "waf_internal_error")
			return false
		}
		// Write*Body triggers Process*Body if the bodylimit (Sec*BodyLimit) is reached.
		if interruption != nil {
			p.blockRequest(interruption, ctypes.PhaseResponseBody, "waf_response_body_overflow")
			return false
		}
	}
	return true
}

func (p *wafPlugin) handleResponseBody() bool {
	p.responseBodyProcessed = true

	interruption, err := p.txContext.ProcessResponseBody()
	if err != nil {
		p.handle.Log(shared.LogLevelInfo, "Failed to process response body in WAF: %v", err.Error())
		p.blockRequest(nil, ctypes.PhaseResponseBody, "waf_internal_error")
		return false
	}
	if interruption != nil {
		p.blockRequest(interruption, ctypes.PhaseResponseBody, "waf_response_body_blocked")
		return false
	}

	return true
}

// blockRequest is a helper method to send a local response with the appropriate status and body when a request is blocked by the WAF.
func (p *wafPlugin) blockRequest(interruption *ctypes.Interruption, phase ctypes.RulePhase, reason string) {
	var status int

	if interruption == nil {
		// If we have to block the request without a WAF interruption, that means some internal error happened.
		// In this case we record the metrics without the ruleID.
		status = 500
		p.metrics.RecordBlockInternal(p.handle, p.authority, phase)
	} else {
		status = interruption.Status
		if status == 0 {
			status = 403
		}
		p.metrics.RecordBlockedByRule(p.handle, p.authority, phase, interruption.RuleID)
		p.handle.SetMetadata(p.metadataNamespace, metadataKeyBlockRule, interruption.RuleID)
	}
	p.handle.SetMetadata(p.metadataNamespace, metadataKeyBlockPhase, int(phase))

	p.handle.SendLocalResponse(
		uint32(status), //nolint:gosec // status is validated to be non-zero
		nil,
		[]byte("Blocked by WAF"),
		reason,
	)
}

// ExtensionName is the name of the extension that will be used in the `run` command to refer to this embedded plugin.
const ExtensionName = "coraza-waf"

var wellKnownHTTPFilterConfigFactories = map[string]shared.HttpFilterConfigFactory{
	ExtensionName: &wafPluginConfigFactory{},
}

// WellKnownHttpFilterConfigFactories returns the map of well-known HTTP filter config factories.
func WellKnownHttpFilterConfigFactories() map[string]shared.HttpFilterConfigFactory { // nolint:revive
	return wellKnownHTTPFilterConfigFactories
}
