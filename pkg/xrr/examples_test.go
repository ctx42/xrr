// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package xrr_test

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ctx42/testing/pkg/must"

	"github.com/ctx42/xrr/pkg/xrr"
)

func ExampleNew() {
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
}

func ExampleNew_with_metadata() {
	meta := xrr.Meta().Str("action", "context")
	err := xrr.New("user not found", "EC_USER_NOT_FOUND", meta.Option())

	fmt.Printf("%s\n", must.Value(json.MarshalIndent(err, "", "  ")))
	// Output:
	// {
	//   "code": "EC_USER_NOT_FOUND",
	//   "error": "user not found",
	//   "meta": {
	//     "action": "context"
	//   }
	// }
}

func ExampleFields() {
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
}
