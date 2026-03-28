[![Go Report Card](https://goreportcard.com/badge/github.com/ctx42/xrr)](https://goreportcard.com/report/github.com/ctx42/xrr)
[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/ctx42/xrr)
![Tests](https://github.com/ctx42/xrr/actions/workflows/go.yml/badge.svg?branch=master)

# xrr

Extend standard Go errors with optional string codes and structured metadata.
Use codes for clear API responses and monitoring. Add metadata to enrich logs
without extra boilerplate.

Requires **Go 1.24+**.

```bash
go get github.com/ctx42/xrr
```

<!-- TOC -->
* [xrr](#xrr)
  * [Quick Start](#quick-start)
  * [Error Codes](#error-codes)
  * [Metadata and Logging](#metadata-and-logging)
  * [Wrapping Errors](#wrapping-errors)
  * [Inspecting Error Trees](#inspecting-error-trees)
  * [Field Errors](#field-errors)
  * [Domain-Specific Errors](#domain-specific-errors)
  * [Envelope](#envelope)
    * [Regular Error](#regular-error)
    * [Joined Errors](#joined-errors)
    * [Fields Error](#fields-error)
  * [Error Collections](#error-collections)
  * [Test Helpers](#test-helpers)
<!-- TOC -->

## Quick Start

```go
import "github.com/ctx42/xrr/pkg/xrr"

// Create an error with a code.
err := xrr.New("user not found", "EC_USER_NOT_FOUND")

// Add structured metadata.
meta := xrr.Meta().Str("user_id", "u-123")
err = xrr.New("user not found", "EC_USER_NOT_FOUND", meta.Option())

// Inspect the error.
fmt.Println(xrr.GetCode(err))            // EC_USER_NOT_FOUND
fmt.Println(xrr.GetStr(err, "user_id"))  // u-123 true
```

## Error Codes

Error handling in Go is straightforward: functions return an `error` type, you
check if it's `nil`, and that's enough in most cases. But when returning the
error through an API, just a string is often not enough. The bare minimum is to
return an error message and its code.

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

Use `IsCode` to check whether any error in a chain carries a specific code:

```go
if xrr.IsCode(err, "EC_USER_NOT_FOUND") {
    // handle not found
}
```

## Metadata and Logging

In most server applications, errors aren't just handled — they're logged for
diagnostics, auditing, and post-mortem analysis. Structured metadata bridges
the worlds of error handling and logging:

```go
meta := xrr.Meta().Str("action", "context")
err := xrr.New("user not found", "EC_USER_NOT_FOUND", meta.Option())

buf := &bytes.Buffer{}
log := zerolog.New(buf)
log.Info().
    Fields(xrr.GetMeta(err)).
    Str("error_code", xrr.GetCode(err)).
    Err(err).
    Send()

fmt.Printf("%s\n", buf.String())
// Output:
// {"level":"info","action":"context","error_code":"EC_USER_NOT_FOUND","error":"user not found"}
```

The `Metadata` builder supports typed values — `bool`, `string`, `int`, `int64`,
`float64`, `time.Time`, `time.Duration`:

```go
meta := xrr.Meta().
    Str("user_id", "u-123").
    Int("attempt", 3).
    Bool("retryable", true)
```

## Wrapping Errors

Use `Wrap` to add a code or metadata to an existing error. It preserves the
original error in the chain and retains its code unless overridden:

```go
err := fmt.Errorf("connection refused")
wrapped := xrr.Wrap[xrr.EDGeneric](err, xrr.WithCode("EC_CONN"))

fmt.Println(errors.Is(wrapped, err)) // true
fmt.Println(xrr.GetCode(wrapped))    // EC_CONN
```

## Inspecting Error Trees

The `Get*` functions walk the full error tree (chains, joined errors, and
field errors) using BFS. For metadata, errors closer to the root override
deeper ones.

```go
// Retrieve the first error code in the chain.
code := xrr.GetCode(err)

// Collect all unique codes.
codes := xrr.GetCodes(err)

// Retrieve merged metadata from the entire tree.
meta := xrr.GetMeta(err)

// Retrieve typed metadata values.
userID, ok := xrr.GetStr(err, "user_id")
count, ok  := xrr.GetInt(err, "attempt")
flag, ok   := xrr.GetBool(err, "retryable")
ts, ok     := xrr.GetTime(err, "created_at")
dur, ok    := xrr.GetDuration(err, "elapsed")
```

## Field Errors

A common requirement in validation is to associate errors with specific
fields. `Fields` is a `map[string]error` underneath:

```go
err := xrr.Fields{
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

Inspect field errors programmatically:

```go
names := xrr.FieldNames(err)                 // sorted field names
fe := xrr.GetFieldError(err, "email")        // get a specific field
ok := xrr.FieldErrorIs(err, "email", target) // check field error chain
```

Merge and flatten nested field errors:

```go
merged := xrr.MergeFields[xrr.EDGeneric](fieldsA, fieldsB)
flat := xrr.Flatten[xrr.EDGeneric](nested) // dot-notation keys
```

## Domain-Specific Errors

Use Go generics to create error types scoped to your domain. This gives you
compile-time separation between error types from different parts of your
system:

```go
type EDPayment string

var (
    NewPaymentError = xrr.NewErrorFor[EDPayment]()
    PaymentFieldErr = xrr.NewFieldErrorFor[EDPayment]()
)

err := NewPaymentError("charge failed", "EC_CHARGE_FAILED")
```

Domain-specific errors carry the same functionality as the default `xrr.New`
— codes, metadata, JSON marshaling, and error tree traversal — but their
Go types are distinct, enabling type-based routing or assertions.

## Envelope

The `Envelope` provides facilities to create a JSON envelope for errors of
different kinds — useful for API responses.

### Regular Error

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

### Joined Errors

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

### Fields Error

```go
cause := xrr.Fields{
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

## Error Collections

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
collected := errs.Collect()
```

## Test Helpers

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
