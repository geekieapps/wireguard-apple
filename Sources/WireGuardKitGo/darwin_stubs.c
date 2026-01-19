// SPDX-License-Identifier: MIT
//
// Stub implementations for Darwin/Mach functions not available in iOS Simulator

#include <TargetConditionals.h>

#if TARGET_OS_SIMULATOR

void _darwin_arm_init_mach_exception_handler(void) {
    // No-op for simulator
}

void _darwin_arm_init_thread_exception_port(void) {
    // No-op for simulator
}

#endif
