// Package harnessx is the single choke-point every SWE-AF role reasoner uses to
// call the AgentField harness. It replaces the Python monkeypath of app.harness (swe_af/app.py:80-93) with an explicit generic wrapper:
//
//   - schemaFor[T] reflects a Go struct into the JSON-schema map the harness consumes to build its OUTPUT REQUIREMENTS prompt suffix (design §2.3).
//   - Run[T] injects the build's run-scoped credentials into the subprocess env (scoped creds win over the base, mirroring Python precedence), calls the harness, classifies fatal API errors, and on a schema parse-failure hands the caller a default-seeded value so it can apply its deterministic fallback (design §4.1).
//
// This is the ONLY way roles should reach the harness —!typo fix me later