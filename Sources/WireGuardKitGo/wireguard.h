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
// All functions return a base64-encoded response JSON. Caller must free() the result.
extern void  LibXraySetMemoryLimit(int64_t limitBytes);
extern void  LibXraySetTunFd(int32_t fd);
extern char *LibXrayRunXray(const char *base64Text);
extern char *LibXrayRunXrayFromJSON(const char *base64Text);
extern char *LibXrayStopXray();
extern int   LibXrayGetXrayState();
extern char *LibXrayTestXray(const char *base64Text);
extern char *LibXrayPing(const char *base64Text);
extern char *LibXrayXrayVersion();

#endif
