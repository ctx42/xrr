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

<!-- gmdoceg:pkg/xrr/Example -->
```go
ts := time.Date(2026, 5, 9, 11, 23, 0, 0, time.UTC)
meta := xrr.Meta().
	Int("attempt", 3).
	Str("user_id", "u-123").
	Time("timestamp", ts)
err := xrr.New("user not found", "EC_USER_NOT_FOUND", meta.Option())

handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		return a
	},
})
slog.New(handler).Error(
	"request failed",
	"code", xrr.GetCode(err),
	"meta", xrr.GetMeta(err),
)

fmt.Printf("%s\n", must.Value(json.MarshalIndent(err, "", "  ")))
// Output:
// {"level":"ERROR","msg":"request failed","code":"EC_USER_NOT_FOUND","meta":{"attempt":3,"timestamp":"2026-05-09T11:23:00Z","user_id":"u-123"}}
// {
//   "code": "EC_USER_NOT_FOUND",
//   "error": "user not found",
//   "meta": {
//     "attempt": 3,
//     "timestamp": "2026-05-09T11:23:00Z",
//     "user_id": "u-123"
//   }
// }
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
	xrr.Meta().Float64("amount", 99.50).Option(),
)
```

Now you get full compile-time safety:

<!-- gmdoceg:pkg/xrr/ExampleIsDomain -->
```go
type edPayment struct{}
type PaymentError = xrr.GenericError[edPayment]
newPaymentError := xrr.ErrorFunc[edPayment]()

err := newPaymentError("insufficient funds", "EC_INSUFFICIENT_FUNDS")

var pe *PaymentError
fmt.Println(errors.As(err, &pe))
fmt.Println(xrr.IsDomain[edPayment](err))
// Output:
// true
// true
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

<!-- gmdoceg:pkg/xrr/ExampleGetMeta -->
```go
err := xrr.New("payment failed", "EC_PAYMENT_FAILED",
	xrr.Meta().
		Str("order_id", "ord-456").
		Int64("user_id", 12345).
		Float64("amount", 199.99).
		Bool("retryable", true).
		Option(),
)

fmt.Println(xrr.GetCode(err))

orderID, _ := xrr.GetStr(err, "order_id")
amount, _ := xrr.GetFloat64(err, "amount")
retryable, _ := xrr.GetBool(err, "retryable")
fmt.Println(orderID, amount, retryable)
// Output:
// EC_PAYMENT_FAILED
// ord-456 199.99 true
```

---

## Wrapping & Cause Chains

<!-- gmdoceg:pkg/xrr/ExampleWrap -->
```go
err := fmt.Errorf("connection refused")
wrapped := xrr.Wrap(err, xrr.WithCode("EC_CONN"))

fmt.Println(errors.Is(wrapped, err))
fmt.Println(xrr.GetCode(wrapped))
fmt.Println(wrapped.Error())
// Output:
// true
// EC_CONN
// connection refused
```

<!-- gmdoceg:pkg/xrr/ExampleNew_withCauseAndMsg -->
```go
err := fmt.Errorf("connection refused")
wrapped := xrr.New("dial failed", "EC_CONN", xrr.WithCause(err))

fmt.Println(errors.Is(wrapped, err))
fmt.Println(xrr.GetCode(wrapped))
fmt.Println(wrapped.Error())
// Output:
// true
// EC_CONN
// dial failed: connection refused
```

Full support for `errors.Is`, `errors.As`, and `%w`.

---

## JSON & Structured Logging

All errors implement `json.Marshaler` and provide `GetMeta()` for `slog`.

<!-- gmdoceg:pkg/xrr/ExampleNew_marshal_with_metadata -->
```go
meta := xrr.Meta().Int("attempt", 3).Str("user_id", "u-123")
err := xrr.New("user not found", "EC_USER_NOT_FOUND", meta.Option())

fmt.Printf("%s\n", must.Value(json.MarshalIndent(err, "", "  ")))
// Output:
// {
//   "code": "EC_USER_NOT_FOUND",
//   "error": "user not found",
//   "meta": {
//     "attempt": 3,
//     "user_id": "u-123"
//   }
// }
```

<!-- gmdoceg:pkg/xrr/ExampleNew_with_slog -->
```go
meta := xrr.Meta().Int("attempt", 3).Str("user_id", "u-123")
err := xrr.New("user not found", "EC_USER_NOT_FOUND", meta.Option())

handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		return a
	},
})

slog.New(handler).Error(
	err.Error(),
	"code", xrr.GetCode(err),
	"meta", xrr.GetMeta(err),
)
// Output:
// {"level":"ERROR","msg":"user not found","code":"EC_USER_NOT_FOUND","meta":{"attempt":3,"user_id":"u-123"}}
```

Perfect for REST APIs and observability platforms.

---

## Field Errors for Validation

<!-- gmdoceg:pkg/xrr/ExampleGenericFields -->
```go
fields := map[string]error{
	"username": errors.New("username not found"),
	"email": xrr.New(
		"invalid email",
		"EC_INVALID_EMAIL",
		xrr.Meta().Str("action", "context").Option(),
	),
}
err := xrr.NewFieldErrors(fields)

fmt.Printf("%s\n", must.Value(json.MarshalIndent(err, "", "  ")))
// Output:
// {
//   "email": {
//     "code": "EC_INVALID_EMAIL",
//     "error": "invalid email",
//     "meta": {
//       "action": "context"
//     }
//   },
//   "username": {
//     "code": "ECGeneric",
//     "error": "username not found"
//   }
// }
```

Serializes cleanly to JSON arrays/objects for frontend consumption.

---

## Domain-Specific Errors

<!-- gmdoceg:pkg/xrr/ExampleIsDomain -->
```go
type edPayment struct{}
type PaymentError = xrr.GenericError[edPayment]
newPaymentError := xrr.ErrorFunc[edPayment]()

err := newPaymentError("insufficient funds", "EC_INSUFFICIENT_FUNDS")

var pe *PaymentError
fmt.Println(errors.As(err, &pe))
fmt.Println(xrr.IsDomain[edPayment](err))
// Output:
// true
// true
```

---

## API Envelope Pattern

<!-- gmdoceg:pkg/xrr/ExampleEnclose -->
```go
cause := xrr.New("cause", "EC_CAUSE")
lead := xrr.New("lead", "EC_LEAD")

err := xrr.Enclose(cause, lead)

fmt.Printf("is lead error: %v\n", errors.Is(err, cause))
fmt.Printf("id db error: %v\n", errors.Is(err, lead))
fmt.Printf("unwrap: %v\n", errors.Unwrap(err))
fmt.Printf("message: %v\n", err.Error())
fmt.Printf("%s\n", must.Value(json.MarshalIndent(err, "", "  ")))
// Output:
// is lead error: true
// id db error: true
// unwrap: cause
// message: cause
// {
//   "code": "EC_LEAD",
//   "error": "lead",
//   "errors": [
//     {
//       "code": "EC_CAUSE",
//       "error": "cause"
//     }
//   ]
// }
```

---

## Testing with xrrtest

```go
import "github.com/ctx42/xrr/pkg/xrr/xrrtest"

xrrtest.AssertError[payment.EdPayment](t, err)
xrrtest.AssertCode(t, "EC_...", err)
xrrtest.AssertStr(t, "key", "expected", err)
```

---

## Comparison with Alternatives

| Feature                 | std `error` | pkg/errors | go-errors | **xrr**      |
|-------------------------|-------------|------------|-----------|--------------|
| Stable error codes      | No          | No         | No        | Yes          |
| Typed metadata          | No          | No         | Limited   | Yes (fluent) |
| Domain generics         | No          | No         | No        | Yes          |
| Built-in JSON           | No          | No         | No        | Yes          |
| slog integration        | Manual      | Manual     | Manual    | Native       |
| Field validation errors | No          | No         | No        | Yes          |
| Full errors.Is/As       | Yes         | Yes        | Yes       | Yes          |

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
