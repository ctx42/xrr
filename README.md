[![Go Report Card](https://goreportcard.com/badge/github.com/ctx42/xrr)](https://goreportcard.com/report/github.com/ctx42/xrr)
[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/ctx42/xrr)
![Tests](https://github.com/ctx42/xrr/actions/workflows/go.yml/badge.svg?branch=master)

<!-- TOC -->
* [Errors with Codes and Metadata](#errors-with-codes-and-metadata)
  * [Error Codes](#error-codes)
  * [Bridge Errors and Logging](#bridge-errors-and-logging)
<!-- TOC -->

# Errors with Codes and Metadata

Extend standard Go errors with optional string codes and structured metadata.
Use codes for clear API responses and monitoring. Add metadata to enrich logs
without an extra boilerplate.

**Codes**: String identifiers (e.g., `EC_USER_NOT_FOUND`) for APIs, monitoring, 
and consistency.

**Metadata**: Structured key-value data that sticks with errors, perfect for
rich logging in Zerolog, Zap, or beyond.

Keep Go's simplicity. Use only what you need.

Go's error handling is one of its most beloved features. It encourages
developers to handle failures directly without the overhead of exceptions or
complex mechanisms. This library is merely an extension of that philosophy,
building on the standard `error` interface. It doesn't replace or complicate 
the basics; instead, it adds optional enhancements for scenarios where more
structure is needed, particularly in larger applications or services.

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

## Bridge Errors and Logging

In most server applications — be it web services, APIs, or backend systems —
errors aren't just handled; they're logged for diagnostics, auditing, and
post-mortem analysis. Traditional Go errors are flat strings, which work fine
for simple logging but lose context in structured loggers (e.g., Zerolog). This
is where structured metadata bridges the worlds of error handling and logging:

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

# Other Error Types

This module includes a handful of commonly used error types, built for fast
adoption in everyday cases like validation failures. They stick closely to Go's 
error conventions, so they feel like a natural fit without extra ceremony.

## Fields

A common requirement, for example, in validation errors, is to associate an 
error with a specific field. The `Fields` type handles this, and it is just a 
map of errors underneath.

```go
type Fields map[string]error
```

This makes it straightforward to return detailed responses in for example APIs 
or process / handle field-specific errors.

```go
type Fields map[string]error
```

Example:

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

## Envelope

The `Envelope` provides facilities to create JSON envelope for errors of 
different kinds.

### Regular Error

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