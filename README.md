[![Go Report Card](https://goreportcard.com/badge/github.com/ctx42/xrr)](https://goreportcard.com/report/github.com/ctx42/xrr)
[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/ctx42/xrr)
![Tests](https://github.com/ctx42/xrr/actions/workflows/go.yml/badge.svg?branch=master)

`xrr` extends standard Go errors with optional string codes and structured
metadata, building on the standard `error` interface without replacing it.

```bash
go get github.com/ctx42/xrr
```

<!-- TOC -->
* [Quick Start](#quick-start)
* [Error Codes](#error-codes)
* [Metadata and Logging](#metadata-and-logging)
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

```go
import "github.com/ctx42/xrr/pkg/xrr"

meta := xrr.Meta().Str("user_id", "u-123")
err := xrr.New("user not found", "EC_USER_NOT_FOUND", meta.Option())

xrr.GetCode(err)                     // EC_USER_NOT_FOUND
val, _ := xrr.GetStr(err, "user_id") // val == "u-123"
```

# Error Codes

```go
err := xrr.New("user not found", "EC_USER_NOT_FOUND")

fmt.Printf("%v\n", err)              // Print message.
fmt.Printf("%+v\n", err)             // Print message and error code.
fmt.Printf("%s\n", xrr.GetCode(err)) // Print error code.
fmt.Printf("%s\n", must.Value(json.MarshalIndent(err, "", "  ")))

// Output:
// user not found
// user not found (EC_USER_NOT_FOUND)
// EC_USER_NOT_FOUND
// {
//   "code": "EC_USER_NOT_FOUND",
//   "error": "user not found"
// }
```

The `GetCode` function returns `ECGeneric` when no error in the chain
implements `Coder`.

Use `IsCode` to check whether any error in a chain carries a specific code:

```go
if xrr.IsCode(err, "EC_USER_NOT_FOUND") {
    // handle not found
}
```

# Metadata and Logging

The `Metadata` builder supports `bool`, `string`, `int`, `int64`, `float64`,
`time.Time`, and `time.Duration` values:

```go
meta := xrr.Meta().
    Str("user_id", "u-123").
    Int("attempt", 3).
    Bool("retryable", true)

err := xrr.New("user not found", "EC_USER_NOT_FOUND", meta.Option())

```

Metadata integrates naturally with structured loggers:

```go
log.Info().
    Fields(xrr.GetMeta(err)).
    Str("error_code", xrr.GetCode(err)).
    Err(err).
    Send()

fmt.Printf("%s\n", buf.String())
// Output:
// {"level":"info","action":"context","error_code":"EC_USER_NOT_FOUND","error":"user not found"}
```

To copy metadata from an existing error, use `WithMetaFrom`:

```go
wrapped := xrr.New("request failed", "EC_REQUEST", xrr.WithMetaFrom(original))
```

# Wrapping Errors

`Wrap` adds a code or metadata to an existing error, preserving it in the
chain and retaining its code unless overridden:

```go
err := fmt.Errorf("connection refused")
wrapped := xrr.Wrap[xrr.EDGeneric](err, xrr.WithCode("EC_CONN"))

fmt.Println(errors.Is(wrapped, err)) // true
fmt.Println(xrr.GetCode(wrapped))    // EC_CONN
```

# Inspecting Error Trees

The `Get*` functions walk the full error tree (single-wrapped, joined, and
field errors) recursively:

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

To check whether an error belongs to a specific domain, use `IsDomain`:

```go
if xrr.IsDomain[EDPayment](err) {
    // err is a *GenericError[EDPayment]
}
```

# Field Errors

`GenericFields[T]` is a `map[string]error` for associating errors with named
fields — most commonly used for validation:

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

Inspect field errors programmatically:

```go
names := xrr.FieldNames(err)                 // Sorted field names.
fe := xrr.GetFieldError(err, "email")        // Gets a specific field's error.
ok := xrr.FieldErrorIs(err, "email", target) // Checks a field's error chain.
```

Build and manipulate field error maps:

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

Use `FieldError` to create a single-field error:

```go
err := xrr.FieldError("email", xrr.New("invalid email", "EC_INVALID_EMAIL"))
```

# Domain-Specific Errors

Use Go generics to create error types scoped to your domain, giving
compile-time separation between errors from different subsystems:

```go
type EDPayment string

var (
    NewPaymentError = xrr.NewErrorFor[EDPayment]()
    PaymentFieldErr = xrr.NewFieldErrorFor[EDPayment]()
)

err := NewPaymentError("charge failed", "EC_CHARGE_FAILED")
```

Domain-specific errors carry the same functionality as `xrr.New` — codes,
metadata, JSON marshaling, and error tree traversal — but their Go types are
distinct, enabling type-based routing or assertions.

# Error Utilities

The following functions work with joined errors:

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

A few additional helpers are also available:

```go
// Return the first non-empty code from the list, or the fallback.
code := xrr.DefaultCode("ECGeneric", codeA, codeB)
```

# Sentinel Errors

```go
xrr.ErrInvJSON      // invalid JSON (code: ECInvJSON)
xrr.ErrInvJSONError // JSON is valid but not a GenericError representation
// (code: ECInvJSONError)
xrr.ErrFields // generic field error (code: ECFields)
```

These are returned or wrapped by the library itself and can be used with
`errors.Is` for targeted handling.

# Envelope

`Envelope` wraps errors into a structured JSON envelope for API responses.

## Regular Error

```go
cause := xrr.New("cause", "EC_CAUSE")
lead := xrr.New("lead", "EC_LEAD")

err := xrr.Enclose(cause, lead)

fmt.Printf("is cause: %v\n", errors.Is(err, cause))
fmt.Printf("is lead: %v\n", errors.Is(err, lead))
fmt.Printf("unwrap: %v\n", errors.Unwrap(err))
fmt.Printf("message: %v\n", err.Error())
fmt.Printf("%s\n", must.Value(json.MarshalIndent(err, "", "  ")))

// Output:
// is cause: true
// is lead: true
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

`Errors` is a simple `[]error` slice for accumulating errors:

```go
errs := xrr.NewErrors()
errs.Add(fmt.Errorf("first"))
errs.Add(fmt.Errorf("second"))
fmt.Println(errs.First()) // first
```

`SyncErrors` is the thread-safe version, suitable for use across goroutines:

```go
errs := xrr.NewSyncErrors()
errs.Add(fmt.Errorf("from goroutine"))
collected := errs.Collect() // drains and returns all errors
```

# Test Helpers

The `xrrtest` subpackage provides assertion helpers built on
[`github.com/ctx42/testing`](https://github.com/ctx42/testing):

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
