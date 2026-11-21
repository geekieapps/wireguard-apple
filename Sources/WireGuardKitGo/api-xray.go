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
// Call this before invoking runXray. limitBytes <= 0 disables the limit.
//
//export LibXraySetMemoryLimit
func LibXraySetMemoryLimit(limitBytes C.int64_t) {
	debug.SetMemoryLimit(int64(limitBytes))
}

// LibXrayInvoke forwards to libxray's single structured entrypoint.
//
// requestJSON is a plain (not base64) UTF-8 JSON envelope:
//
//	{"apiVersion": 3, "method": "runXray", "payload": {"xrayJson": "..."}}
//
// Returns a plain JSON response `{"success": bool, "data": ..., "error": "..."}`.
// The caller must free() the returned string.
//
// The TUN file descriptor is no longer set through a separate call: write it
// into the Xray config's root `env` object as `"xray.tun.fd"` before runXray.
//
//export LibXrayInvoke
func LibXrayInvoke(requestJSON *C.char) *C.char {
	return C.CString(libxray.Invoke(C.GoString(requestJSON)))
}
