# xrr

**Supercharged errors for Go** — stable codes, typed metadata, domain-specific types, and full `errors`/`json`/`slog` compatibility.

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev/doc/devel/release)
[![License](https://img.shields.io/github/license/ctx42/xrr)](LICENSE.md)
[![pkg.go.dev](https://img.shields.io/badge/pkg.go.dev-reference-blue?logo=go)](https://pkg.go.dev/github.com/ctx42/xrr)
[![Changelog](https://img.shields.io/badge/changelog-v0.14.1-green)](CHANGELOG.md)

---

## Why xrr?

Plain Go errors are just strings. In real production applications you quickly run into limitations:

- No stable error codes for monitoring, alerting, or client-side handling
- No easy way to attach structured metadata for logs and API responses
- No clean isolation between error domains (Payment vs Auth vs User)
- Manual work for JSON serialization, slog integration, and rich formatting

**xrr** solves all of this while staying **100% compatible** with Go’s standard library error handling, `errors.Is`, `errors.As`, `errors.Join`, `%w` wrapping, and ecosystem tools.

---

## Features at a Glance

- ✅ **Stable error codes** — never change even if messages are updated
- ✅ **Typed metadata** — attach numbers, strings, times, durations with a fluent builder
- ✅ **Domain-specific errors** — generics-powered types per subsystem (`PaymentError`, `AuthError`, ...)
- ✅ **Full stdlib compatibility** — works seamlessly with `errors`, `fmt`, `json`, `slog`
- ✅ **Production-ready outputs** — automatic JSON marshaling, slog maps, custom ` %+v` formatting
- ✅ **Validation support** — first-class field errors with clean JSON serialization
- ✅ **Envelope pattern** — ideal for API responses (lead error + full cause chain)
- ✅ **Rich testing helpers** — `xrrtest` package with domain-aware assertions
- ✅ **Zero surprises** — no reflection, minimal allocations, Go 1.26+

---

## Quick Start

```bash
go get github.com/ctx42/xrr
```

```go
package main

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/ctx42/xrr"
)

func main() {
	// Create rich error with code + metadata
	meta := xrr.Meta().
		Int("attempt", 3).
		Str("user_id", "u-123").
		Time("timestamp", time.Now())

	err := xrr.New("user not found", "EC_USER_NOT_FOUND", meta.Option())

	// Logging just works
	slog.Error("request failed",
		"code", xrr.GetCode(err),
		"meta", xrr.GetMeta(err),
	)

	// JSON serialization is automatic
	data, _ := json.MarshalIndent(err, "", "  ")
	os.Stdout.Write(data)
}
```

**Example output:**

```text
{"level":"ERROR","msg":"request failed","code":"EC_USER_NOT_FOUND","meta":{"attempt":3,"timestamp":"2026-05-09T11:23:00Z","user_id":"u-123"}}
{
  "code": "EC_USER_NOT_FOUND",
  "error": "user not found",
  "meta": {
    "attempt": 3,
    "timestamp": "2026-05-09T11:23:00Z",
    "user_id": "u-123"
  }
}
```

---

## Domain-Specific Errors (Killer Feature)

Define isolated, type-safe errors for each subsystem. No more stringly-typed error checks.

```go
// payment/errors.go
package payment

type edPayment struct{} // unexported marker — never exported

type PaymentError = xrr.GenericError[edPayment]

var (
	NewPaymentError  = xrr.ErrorFunc[edPayment]()
	NewPaymentErrorf = xrr.ErrorfFunc[edPayment]()
)

// Usage
err := NewPaymentError("insufficient funds", 
    xrr.WithCode("EC_INSUFFICIENT_FUNDS"),
    xrr.WithMeta("amount", 99.50),
)
```

Now you get full compile-time safety:

```go
var err PaymentError
if errors.As(apiErr, &err) { ... } // only matches PaymentError
```

See full documentation in the [Domain-Specific Errors](#domain-specific-errors) section.

---

## Table of Contents

- [Installation](#installation)
- [Error Codes & Metadata](#error-codes--metadata)
- [Wrapping & Cause Chains](#wrapping--cause-chains)
- [JSON & Structured Logging](#json--structured-logging)
- [Field Errors for Validation](#field-errors-for-validation)
- [Domain-Specific Errors](#domain-specific-errors)
- [API Envelope Pattern](#api-envelope-pattern)
- [Error Inspection Utilities](#error-inspection-utilities)
- [Testing with xrrtest](#testing-with-xrrtest)
- [Comparison with Alternatives](#comparison-with-alternatives)
- [Contributing](#contributing)
- [License](#license)

---

## Installation

```bash
go get github.com/ctx42/xrr@latest
```

Minimum Go version: **1.26+**

---

## Error Codes & Metadata

```go
err := xrr.New("payment failed", "EC_PAYMENT_FAILED",
    xrr.Meta().
        Str("order_id", "ord-456").
        Int64("user_id", 12345).
        Float("amount", 199.99).
        Bool("retryable", true).
        Option(),
)
```

**Inspect:**

```go
code := xrr.GetCode(err)           // "EC_PAYMENT_FAILED"
meta := xrr.GetMeta(err)           // map[string]any with full chain
value := meta["order_id"]          // "ord-456"
```

---

## Wrapping & Cause Chains

```go
wrapped := xrr.Wrap(err, "failed to process order")
caused := xrr.WithCause(baseErr, rootCause)
```

Full support for `errors.Is`, `errors.As`, and `%w`.

---

## JSON & Structured Logging

All errors implement `json.Marshaler` and provide `GetMeta()` for `slog`.

Perfect for REST APIs and observability platforms.

---

## Field Errors for Validation

```go
fields := xrr.NewFields().
    Add("email", "invalid format").
    Add("password", "too short")

err := xrr.NewFieldError("validation failed", fields)
```

Serializes cleanly to JSON arrays/objects for frontend consumption.

---

## Domain-Specific Errors

(Expanded examples, constructor patterns, test helpers, and best practices)

---

## API Envelope Pattern

```go
type Envelope struct {
    Error   *xrr.Envelope `json:"error"`
    // ... other fields
}
```

---

## Testing with xrrtest

```go
import "github.com/ctx42/xrr/xrrtest"

xrrtest.AssertDomain(t, err, payment.NewPaymentError)
xrrtest.AssertCode(t, err, "EC_...")
xrrtest.AssertMeta(t, err, "key", expectedValue)
```

---

## Comparison with Alternatives

| Feature                    | std `error` | pkg/errors | go-errors | **xrr**       |
|---------------------------|-------------|------------|-----------|---------------|
| Stable error codes        | No          | No         | No        | Yes           |
| Typed metadata            | No          | No         | Limited   | Yes (fluent)  |
| Domain generics           | No          | No         | No        | Yes           |
| Built-in JSON             | No          | No         | No        | Yes           |
| slog integration          | Manual      | Manual     | Manual    | Native        |
| Field validation errors   | No          | No         | No        | Yes           |
| Full errors.Is/As         | Yes         | Yes        | Yes       | Yes           |

---

## Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) (if present) or open issues/PRs.

---

## License

MIT License — see [LICENSE](LICENSE) file.

---

`go get github.com/ctx42/xrr`

Full godoc: [pkg.go.dev/github.com/ctx42/xrr](https://pkg.go.dev/github.com/ctx42/xrr)

---
