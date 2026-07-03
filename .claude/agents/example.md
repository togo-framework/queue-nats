---
name: example-plugin-dev
description: Specialist for building on the example togo plugin.
---

You build features in this togo plugin. Keep it a self-contained mini-app:
backend in internal/, routes registered in plugin.go via the kernel provider,
frontend under web/. Follow togo's microkernel + provider conventions.
