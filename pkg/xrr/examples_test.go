// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/ctx42/testing/pkg/must"

	"github.com/ctx42/xrr/pkg/xrr"
)

func ExampleNew() {
	err := xrr.New("user not found", "EC_USER_NOT_FOUND")

	fmt.Printf("%v\n", err)              // Print message.
	fmt.Printf("%+v\n", err)             // Print message and error code.
	fmt.Printf("%s\n", xrr.GetCode(err)) // Print error code.

	// Output:
	// user not found
	// user not found (EC_USER_NOT_FOUND)
	// EC_USER_NOT_FOUND
}

func ExampleNew_with_metadata() {
	meta := xrr.Meta().Int("attempt", 3).Str("user_id", "u-123")
	err := xrr.New("user not found", "EC_USER_NOT_FOUND", meta.Option())

	fmt.Println(xrr.GetMeta(err))
	// Output:
	// map[attempt:3 user_id:u-123]
}

func ExampleNew_marshal() {
	err := xrr.New("user not found", "EC_USER_NOT_FOUND")

	fmt.Printf("%s\n", must.Value(json.MarshalIndent(err, "", "  ")))
	// Output:
	// {
	//   "code": "EC_USER_NOT_FOUND",
	//   "error": "user not found"
	// }
}

func ExampleNew_marshal_with_metadata() {
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
}

func ExampleNew_with_slog() {
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
}

func ExampleWrap() {
	type edMyDomain struct{}
	err := fmt.Errorf("connection refused") // Some action error.
	wrapped := xrr.Wrap[edMyDomain](err, xrr.WithCode("EC_CONN"))

	fmt.Println(errors.Is(wrapped, err))
	fmt.Println(xrr.GetCode(wrapped))
	fmt.Println(wrapped.Error())
	// Output:
	// true
	// EC_CONN
	// connection refused
}

func ExampleNew_withCause() {
	err := fmt.Errorf("connection refused")
	wrapped := xrr.New("", "EC_CONN", xrr.WithCause(err))

	fmt.Println(errors.Is(wrapped, err))
	fmt.Println(xrr.GetCode(wrapped))
	fmt.Println(wrapped.Error())
	// Output:
	// true
	// EC_CONN
	// connection refused
}

func ExampleNew_withCauseAndMsg() {
	err := fmt.Errorf("connection refused")
	wrapped := xrr.New("dial failed", "EC_CONN", xrr.WithCause(err))

	fmt.Println(errors.Is(wrapped, err))
	fmt.Println(xrr.GetCode(wrapped))
	fmt.Println(wrapped.Error())
	// Output:
	// true
	// EC_CONN
	// dial failed: connection refused
}

func ExampleGenericFields() {
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
}

func ExampleEnclose() {
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
}

func ExampleEnclose_joined_errors() {
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
}

func ExampleEnclose_fields_error() {
	fields := map[string]error{
		"a": xrr.New("cause A", "EC_A"),
		"b": xrr.New("cause B", "EC_B"),
	}
	cause := xrr.NewFieldErrors(fields)
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
}

func ExampleSplit() {
	joined := errors.Join(
		xrr.New("first", "EC_FIRST"),
		xrr.New("second", "EC_SECOND"),
	)

	for _, p := range xrr.Split(joined) {
		fmt.Println(p)
	}
	// Output:
	// first
	// second
}

func ExampleJoin() {
	combined := xrr.Join(
		xrr.New("first", "EC_FIRST"),
		nil,
		xrr.New("second", "EC_SECOND"),
	)

	fmt.Println(xrr.IsJoined(combined))
	for _, p := range xrr.Split(combined) {
		fmt.Println(p)
	}
	// Output:
	// true
	// first
	// second
}

func ExampleIsJoined() {
	single := xrr.New("single error", "EC_SINGLE")
	joined := errors.Join(
		xrr.New("first", "EC_FIRST"),
		xrr.New("second", "EC_SECOND"),
	)

	fmt.Println(xrr.IsJoined(single))
	fmt.Println(xrr.IsJoined(joined))
	// Output:
	// false
	// true
}

func ExampleDefaultCode() {
	code := xrr.DefaultCode("ECFallback", "", "EC_FOUND", "EC_IGNORED")

	fmt.Println(code)
	// Output:
	// EC_FOUND
}

func ExampleErrInvJSONError() {
	var e xrr.Error
	err := json.Unmarshal([]byte(`{"status": "ok"}`), &e)

	fmt.Println(errors.Is(err, xrr.ErrInvJSONError))
	// Output:
	// true
}

func ExampleErrFields() {
	cause := xrr.NewFieldErrors(map[string]error{
		"email": xrr.New("invalid email", "EC_INVALID_EMAIL"),
	})
	err := xrr.Enclose(cause)

	fmt.Printf("%s\n", must.Value(json.MarshalIndent(err, "", "  ")))
	// Output:
	// {
	//   "code": "ECFields",
	//   "error": "fields error",
	//   "fields": {
	//     "email": {
	//       "code": "EC_INVALID_EMAIL",
	//       "error": "invalid email"
	//     }
	//   }
	// }
}
