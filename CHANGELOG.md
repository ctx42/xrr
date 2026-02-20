## v0.5.2 (Fri, 20 Feb 2026 16:39:33 UTC)
- chore: Update dependencies.

## v0.5.1 (Mon, 24 Nov 2025 09:28:17 UTC)
- chore: Update dependencies.

## v0.5.0 (Sat, 20 Sep 2025 09:30:06 UTC)
- feat: Only the supported types are used for `xrr.Error` metadata when using `xrr.WithMeta`.
- fix: Error code was set to empty string instead od `xrr.ECGeneric`.
- chore: More realistic test fixtures.
- test: Add tests for `xrr/Error.Error` method.
- feat: Implement custom error string marshaling where joined error messages are concatenated with ";" instead of "\n".
- feat!: The `xrr.Split` and `xrr.IsJoined` no longer use `errors.As` and test the passed instance directly.

## v0.4.0 (Wed, 17 Sep 2025 20:06:55 UTC)
- doc: Update documentation.
- feat!: Rename `vrr/Envelope.ErrCode` method to `vrr/Envelope.ErrorCode` to match the `vrr.Coder` interface.

## v0.3.0 (Mon, 15 Sep 2025 13:39:41 UTC)
- chore: Update dependencies.

## v0.2.0 (Sun, 14 Sep 2025 18:45:18 UTC)
- feat: Add `xrrtest.AssertMsg` asserting just the error message.
- test: Improve the `xrrtest` package tests readability.
- feat: Implement `xrrtest.AssertFieldsEqual` asserting the `xrr.Fields` string representation equals the provided one.

## v0.1.0 (Sat, 13 Sep 2025 08:11:00 UTC)
- feat: Initial commit.
- feat: Add `xrr.DefaultCode` helper.
- feat: Define `xrr` package interfaces and generic error code.
- feat: Define `xrr.Error` structure with a basic constructor and core methods: `xrr/Error.ErrorCode`, `xrr.Error.Unwrap`.
- misc: Add a ` dev ` directory with development-related configuration: Idea test runner configuration.
- feat: Add the `meta` struct with setters for basic types to build Error metadata.
- feat!: Rename `vrr.meta` to `vrr.Metadata`.
- feat: Add `vrr.GetCode` inspection helper.
- feat: Implement `xrr.GetMeta` helper.
- feat: Implement `xrr.WithCode` configuration option and `xrr/Error.MetaAll`.
- chore: Update dependencies.
- feat: Implement `xrr/Error.MarshalJSON` method.
- feat: Add sentinel errors dealing with `xrr.Error` JSON representation.
- feat: Implement `xrr/Error.UnmarshalJSON` method.
- feat: Implement `xrr/Error.Format` method implementing `fmt.Formatter` interface.
- feat: Implement `xrr.GetCodes` helper recursively returning all error codes.
- feat: Implement recursive error chain (tree) walk.
- chore: Remove dead code.
- feat: Implement `xrr.Split` helper.
- feat: Implement thread safe error collection `xrr.SyncErrors`.
- feat: Implement error collection `xrr.Errors`.
- feat: Implement `xrr.Fields`.
- feat: Implement `xrr.Fields`.
- feat: Implement `xrr.IsJoin` helper returning true when an error is implementing `Unwrap() []error`.
- feat: Implement `xrr.Envelope`.
- test: Add `xrr.Wrap` tests.
- doc: Update code documentation.
- doc: Add README documentation.
- chore: Update dependencies.
- chore: Add MIT license.
- feat: Add `xrrtest` package with custom assertions related to `xrr` types.
- doc: Add `Envelope` documentation to README.md file.
- doc: Improve the `Envelope` code documentation.
- test: Improve `Fields` tests.
- test: Add github workflow.
- doc: Update README.md file.

