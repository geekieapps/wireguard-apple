package main

/*
#include <stdlib.h>
*/
import "C"

import (
	libxray "github.com/xtls/libxray"
)

// LibXrayInvoke forwards to libxray's single structured entrypoint.
//
// Memory is libxray's job on iOS: the first runXray in a process sets a 30 MB
// soft heap limit and GC percent 10, and returns freed memory to the OS every
// second (libxray/memory). Don't add a competing SetMemoryLimit export.
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
