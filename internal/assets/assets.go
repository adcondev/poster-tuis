package assets

import _ "embed"

// ══════════════════════════════════════════════════════════════
// Embedded Service Binaries
// ══════════════════════════════════════════════════════════════
// These binaries are built by Taskfile and embedded at compile time.
// Paths are relative to this file (internal/assets/).
//
// Build with: task installer:build
// ══════════════════════════════════════════════════════════════

// BasculaLocalBinary contains the Scale Service (Local) executable
//
//go:embed bin/R2k_BasculaServicio_Local.exe
var BasculaLocalBinary []byte

// BasculaRemoteBinary contains the Scale Service (Remote) executable
//
//go:embed bin/R2k_BasculaServicio_Remote.exe
var BasculaRemoteBinary []byte

// TicketLocalBinary contains the Ticket Service (Local) executable
//
//go:embed bin/R2k_TicketServicio_Local.exe
var TicketLocalBinary []byte

// TicketRemoteBinary contains the Ticket Service (Remote) executable
//
//go:embed bin/R2k_TicketServicio_Remote.exe
var TicketRemoteBinary []byte
