package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"runtime/debug"

	libxray "github.com/xtls/libxray"
)

// SetMemoryLimit sets the Go runtime's soft heap limit.
// Call this before RunXray. limitBytes <= 0 disables the limit.
//
//export LibXraySetMemoryLimit
func LibXraySetMemoryLimit(limitBytes C.int64_t) {
	debug.SetMemoryLimit(int64(limitBytes))
}

// Run Xray instance from a config file on disk.
// base64Text is a base64-encoded RunXrayRequest JSON.
// Returns a base64-encoded response JSON. Caller must free the returned string.
//
//export LibXrayRunXray
func LibXrayRunXray(base64Text *C.char) *C.char {
	return C.CString(libxray.RunXray(C.GoString(base64Text)))
}

// Run Xray instance with inline JSON config.
// base64Text is a base64-encoded RunXrayFromJSONRequest JSON.
// Returns a base64-encoded response JSON. Caller must free the returned string.
//
//export LibXrayRunXrayFromJSON
func LibXrayRunXrayFromJSON(base64Text *C.char) *C.char {
	return C.CString(libxray.RunXrayFromJSON(C.GoString(base64Text)))
}

// Stop the running Xray instance.
// Returns a base64-encoded response JSON. Caller must free the returned string.
//
//export LibXrayStopXray
func LibXrayStopXray() *C.char {
	return C.CString(libxray.StopXray())
}

// Returns 1 if an Xray instance is currently running, 0 otherwise.
//
//export LibXrayGetXrayState
func LibXrayGetXrayState() C.int {
	if libxray.GetXrayState() {
		return 1
	}
	return 0
}

// Validate an Xray config without starting an instance.
// base64Text is a base64-encoded TestXrayRequest JSON.
// Returns a base64-encoded response JSON. Caller must free the returned string.
//
//export LibXrayTestXray
func LibXrayTestXray(base64Text *C.char) *C.char {
	return C.CString(libxray.TestXray(C.GoString(base64Text)))
}

// Measure latency by booting a temporary Xray instance and making one HTTP request.
// base64Text is a base64-encoded PingRequest JSON.
// Returns a base64-encoded response JSON. Caller must free the returned string.
//
//export LibXrayPing
func LibXrayPing(base64Text *C.char) *C.char {
	return C.CString(libxray.Ping(C.GoString(base64Text)))
}

// Returns the xray-core version string. Caller must free the returned string.
//
//export LibXrayXrayVersion
func LibXrayXrayVersion() *C.char {
	return C.CString(libxray.XrayVersion())
}
