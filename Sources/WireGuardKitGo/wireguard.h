/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2018-2023 WireGuard LLC. All Rights Reserved.
 */

#ifndef WIREGUARD_H
#define WIREGUARD_H

#include <sys/types.h>
#include <stdint.h>
#include <stdbool.h>

typedef void(*logger_fn_t)(void *context, int level, const char *msg);
extern void wgSetLogger(void *context, logger_fn_t logger_fn);
extern int wgTurnOn(const char *settings, int32_t tun_fd);
extern void wgTurnOff(int handle);
extern int64_t wgSetConfig(int handle, const char *settings);
extern char *wgGetConfig(int handle);
extern void wgBumpSockets(int handle);
extern void wgDisableSomeRoamingForBrokenMobileSemantics(int handle);
extern const char *wgVersion();

// SingTun — sing-tun gVisor TUN stack with UoT UDP forwarding, compiled into the same Go binary.
// All functions return a base64-encoded response JSON. Caller must free() the result.
extern char *SingTunStart(int tunFd, int mtu, const char *proxyAddr);
extern char *SingTunStop();

// LibXray — upstream xtls/libxray, compiled into the same Go binary.
// LibXrayInvoke takes a plain JSON request envelope ({"apiVersion":3,"method":...,"payload":...})
// and returns a plain JSON response ({"success":...,"data":...,"error":...}).
// Caller must free() the result. The TUN fd goes into the Xray config's root `env` as "xray.tun.fd".
// libxray manages Go memory itself on iOS (libxray/memory), so there's no memory-limit export.
extern char *LibXrayInvoke(const char *requestJSON);

#endif
