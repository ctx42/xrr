// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package xrr_test

import (
	"encoding/json"
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
