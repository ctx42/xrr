[![Go Report Card](https://goreportcard.com/badge/github.com/ctx42/xrr)](https://goreportcard.com/report/github.com/ctx42/xrr)
[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/ctx42/xrr)
![Tests](https://github.com/ctx42/xrr/actions/workflows/go.yml/badge.svg?branch=master)

Standard Go errors carry only a message string. The `xrr` module creates
errors with string codes and optional typed metadata, while remaining
fully compatible with the standard `error` interface — `errors.Is`,
`errors.As`, and wrapping all work as expected.

```bash
go get github.com/ctx42/xrr
```

<!-- TOC -->
* [Quick Start](#quick-start)
  * [Error Codes](#error-codes)
  * [Error Metadata](#error-metadata)
  * [Error Marshaling](#error-marshaling)
  * [Structured Logging](#structured-logging)
* [Wrapping Errors](#wrapping-errors)
* [Inspecting Error Trees](#inspecting-error-trees)
* [Field Errors](#field-errors)
* [Domain-Specific Errors](#domain-specific-errors)
* [Error Utilities](#error-utilities)
* [Sentinel Errors](#sentinel-errors)
* [Envelope](#envelope)
  * [Regular Error](#regular-error)
  * [Joined Errors](#joined-errors)
  * [Fields Error](#fields-error)
* [Error Collections](#error-collections)
* [Test Helpers](#test-helpers)
<!-- TOC -->

# Quick Start

The following examples introduce the two core concepts: _error codes_ and
_structured metadata_.

## Error Codes

Every `xrr` error carries a string code alongside its message, giving
callers a stable identifier that does not depend on the message wording:

<!-- gmdoceg:pkg/xrr/ExampleNew -->
```go
err := xrr.New("user not found", "EC_USER_NOT_FOUND")

fmt.Printf("%v\n", err)              // Print message.
fmt.Printf("%+v\n", err)             // Print message and error code.
fmt.Printf("%s\n", xrr.GetCode(err)) // Print error code.

// Output:
// user not found
// user not found (EC_USER_NOT_FOUND)
// EC_USER_NOT_FOUND
```

## Error Metadata

Attach typed key-value metadata to any error using the `Meta` builder.
Metadata is retrieved by key using typed getters:

<!-- gmdoceg:pkg/xrr/ExampleNew_with_metadata -->
```go
meta := xrr.Meta().Int("attempt", 3).Str("user_id", "u-123")
err := xrr.New("user not found", "EC_USER_NOT_FOUND", meta.Option())

fmt.Println(xrr.GetMeta(err))
// Output:
// map[attempt:3 user_id:u-123]
```

The supported value types are `bool`, `string`, `int`, `int64`, `float64`,
`time.Time`, and `time.Duration`.

When creating a new error from an existing one, use `WithMetaFrom` to
carry its metadata forward without copying the map manually:

```go
original := xrr.New("db timeout", "EC_TIMEOUT",
    xrr.Meta().Str("query", "SELECT ...").Option())

wrapped := xrr.New("request failed", "EC_REQUEST",
    xrr.WithMetaFrom(original))
```

## Error Marshaling

Every `xrr` error implements `json.Marshaler`. A plain error serializes to
its code and message:

<!-- gmdoceg:pkg/xrr/ExampleNew_marshal -->
```go
err := xrr.New("user not found", "EC_USER_NOT_FOUND")

fmt.Printf("%s\n", must.Value(json.MarshalIndent(err, "", "  ")))
// Output:
// {
//   "code": "EC_USER_NOT_FOUND",
//   "error": "user not found"
// }
```

When metadata is present it is included under the `meta` key:

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

## Structured Logging

Metadata is designed to be passed directly to structured loggers.
`GetMeta` retrieves all metadata accumulated across the error chain as
a `map[string]any`, ready for use with `log/slog` or any structured
logger:

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

# Wrapping Errors

When an error originates outside your code — from the standard library or
a third-party package — `Wrap` lets you attach a code or metadata without
losing the original error in the chain:

<!-- gmdoceg:pkg/xrr/ExampleWrap -->
```go
err := fmt.Errorf("connection refused")
wrapped := xrr.Wrap[xrr.EDGeneric](err, xrr.WithCode("EC_CONN"))

fmt.Println(errors.Is(wrapped, err))
fmt.Println(xrr.GetCode(wrapped))
// Output:
// true
// EC_CONN
```

# Inspecting Error Trees

Go errors compose into trees through wrapping and `errors.Join`. The
`Get*` functions traverse the entire tree — including field errors — to
extract codes and metadata without needing to know its shape:

```go
code := xrr.GetCode(err)   // The first code in the chain.
codes := xrr.GetCodes(err) // All unique codes.
meta := xrr.GetMeta(err)   // Merged metadata; root overrides deeper values.

// Typed metadata lookup.
userID, ok := xrr.GetStr(err, "user_id")
count, ok := xrr.GetInt(err, "attempt")
flag, ok := xrr.GetBool(err, "retryable")
ts, ok := xrr.GetTime(err, "created_at")
dur, ok := xrr.GetDuration(err, "elapsed")
```

# Field Errors

`GenericFields[T]` is a `map[string]error` for associating errors with named
fields — most commonly used for validation:

<!-- gmdoceg:pkg/xrr/ExampleGenericFields -->
```go
err := xrr.GenericFields[xrr.EDGeneric]{
	"username": errors.New("username not found"),
	"email": xrr.New(
		"invalid email",
		"EC_INVALID_EMAIL",
		xrr.Meta().Str("action", "context").Option(),
	),
}

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

Once you have a field error map, several functions let you query it
without asserting types manually:

```go
names := xrr.FieldNames(err)                 // Sorted field names.
fe := xrr.GetFieldError(err, "email")        // Gets a specific field's error.
ok := xrr.FieldErrorIs(err, "email", target) // Checks a field's error chain.
```

When building or combining field maps — for example when merging validation
results from multiple sources — the following helpers are available:

```go
// Merge two field maps; existing non-nil keys are not overwritten.
merged := xrr.MergeFields[xrr.EDGeneric](fieldsA, fieldsB)

// Flatten nested field maps to dot-notation keys.
flat := xrr.Flatten[xrr.EDGeneric](nested)

// Remove nil entries.
filtered := fs.Filter()

// Look up a field, including dot-notation paths.
fieldErr := fs.Get("address.city")
```

To wrap a single error under a field name, use `FieldError`:

```go
err := xrr.FieldError("email", xrr.New("invalid email", "EC_INVALID_EMAIL"))
```

# Domain-Specific Errors

By default all `xrr` errors share the same Go type. Using generics you
can create a distinct error type per domain, so callers can identify
which subsystem an error originated from:

```go
type EDPayment string

var (
    NewPaymentError = xrr.NewErrorFor[EDPayment]()
    PaymentFieldErr = xrr.NewFieldErrorFor[EDPayment]()
)

err := NewPaymentError("charge failed", "EC_CHARGE_FAILED")
```

Use `IsDomain[EDPayment](err)` to check at runtime whether an error
originated from the payment domain:

```go
if xrr.IsDomain[EDPayment](err) {
    // err is a *GenericError[EDPayment]
}
```

# Error Utilities

`xrr` provides several helpers that complement the standard `errors`
package when working with joined errors and codes:

```go
// Split a joined error into its constituent errors.
// Returns []error{err} for non-joined errors, nil for nil.
parts := xrr.Split(err)

// Join errors, skipping nils. Returns the single error directly
// when only one non-nil error is present.
combined := xrr.Join(err1, err2, err3)

// Check whether err was created by errors.Join.
if xrr.IsJoined(err) { ... }
```

When selecting a code from multiple candidates, `DefaultCode` returns
the first non-empty value, falling back to the provided default:

```go
// Return the first non-empty code from the list, or the fallback.
code := xrr.DefaultCode("ECGeneric", codeA, codeB)
```

# Sentinel Errors

The library defines sentinel errors for conditions it detects internally.
Use them with `errors.Is` to handle specific failure cases:

```go
xrr.ErrInvJSON      // invalid JSON (code: ECInvJSON)
xrr.ErrInvJSONError // JSON is valid but not a GenericError representation
// (code: ECInvJSONError)
xrr.ErrFields // generic field error (code: ECFields)
```

# Envelope

An `Envelope` combines two errors: a *cause* — the underlying error that
triggered the failure — and a *lead* — a higher-level error describing
the outcome to the caller. The lead's code and message are serialized at
the top level of the JSON response, while the cause is nested inside.
Both remain reachable via `errors.Is`.

Use `Enclose` to create an envelope:

## Regular Error

When the cause is a single error, it is nested under the `errors` key:

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

## Joined Errors

When the cause is a joined error, each constituent error is listed
separately under the `errors` key:

<!-- gmdoceg:pkg/xrr/ExampleEnclose_joined_errors -->
```go
cause := errors.Join(xrr.New("cause A", "EC_A"), xrr.New("cause B", "EC_B"))
lead := xrr.New("lead", "EC_LEAD")

err := xrr.Enclose(cause, lead)

fmt.Printf("%s\n", must.Value(json.MarshalIndent(err, "", "  ")))
// Output:
// {
//   "code": "EC_LEAD",
//   "error": "lead",
//   "errors": [
//     {
//       "code": "EC_A",
//       "error": "cause A"
//     },
//     {
//       "code": "EC_B",
//       "error": "cause B"
//     }
//   ]
// }
```

## Fields Error

When the cause is a `GenericFields` map, the errors are placed under
the `fields` key, keyed by field name:

<!-- gmdoceg:pkg/xrr/ExampleEnclose_fields_error -->
```go
cause := xrr.GenericFields[xrr.EDGeneric]{
	"a": xrr.New("cause A", "EC_A"),
	"b": xrr.New("cause B", "EC_B"),
}
lead := xrr.New("lead", "EC_LEAD")

err := xrr.Enclose(cause, lead)

fmt.Printf("%s\n", must.Value(json.MarshalIndent(err, "", "  ")))
// Output:
// {
//   "code": "EC_LEAD",
//   "error": "lead",
//   "fields": {
//     "a": {
//       "code": "EC_A",
//       "error": "cause A"
//     },
//     "b": {
//       "code": "EC_B",
//       "error": "cause B"
//     }
//   }
// }
```

# Error Collections

When processing multiple independent operations — iterating over a list,
running goroutines in parallel — errors need to be collected and reported
together rather than short-circuiting on the first one.

`Errors` is a simple `[]error` slice for sequential use:

```go
errs := xrr.NewErrors()
errs.Add(fmt.Errorf("first"))
errs.Add(fmt.Errorf("second"))
fmt.Println(errs.First()) // first
```

`SyncErrors` is the thread-safe variant, safe to use concurrently across
goroutines:

```go
errs := xrr.NewSyncErrors()
errs.Add(fmt.Errorf("from goroutine"))
collected := errs.Collect() // drains and returns all errors
```

# Test Helpers

The `xrrtest` subpackage provides assertion helpers for testing `xrr`
errors. They produce clear failure messages and avoid manual type
assertions in test code:

```go
import "github.com/ctx42/xrr/pkg/xrr/xrrtest"

// Assert error type and code.
ge, ok := xrrtest.AssertError[xrr.EDGeneric](t, err)
xrrtest.AssertCode(t, "EC_USER_NOT_FOUND", err)
xrrtest.AssertEqual(t, "user not found (EC_USER_NOT_FOUND)", err)

// Assert metadata.
xrrtest.AssertStr(t, "user_id", "u-123", err)
xrrtest.AssertInt(t, "attempt", 3, err)
xrrtest.AssertBool(t, "retryable", true, err)

// Assert field errors.
xrrtest.AssertFieldCnt(t, 2, err)
xrrtest.AssertFieldEqual(t, "email", "invalid email", err)
xrrtest.AssertFieldCode(t, "email", "EC_INVALID_EMAIL", err)
```
