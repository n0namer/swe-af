### go/internal/harnessx/schema.go (showing first 400 of 429)
```
1: // Package harnessx is the single choke-point every SWE-AF role reasoner uses to
2: // call the AgentField harness. It replaces the Python monkeypatch of
3: // app.harness (swe_af/app.py:80-93) with an explicit generic wrapper:
4: //
5: //   - schemaFor[T] reflects a Go struct into the JSON-schema map the harness
6: //     consumes to build its OUTPUT REQUIREMENTS prompt suffix (design §2.3).
7: //   - Run[T] injects the build's run-scoped credentials into the subprocess env
8: //     (scoped creds win over the base, mirroring Python precedence), calls the
9: //     harness, classifies fatal API errors, and on a schema parse-failure hands
10: //     the caller a default-seeded value so it can apply its deterministic
11: //     fallback (design §4.1).
12: //
13: // This is the ONLY way roles should reach the harness — it guarantees uniform
14: // credential injection and fatal-error handling across all 22 role reasoners.
15: package harnessx
16: 
17: import (
18: 	"bytes"
19: 	"context"
20: 	"encoding/json"
21: 	"fmt"
22: 	"os"
23: 	"path/filepath"
24: 	"reflect"
25: 	"sort"
26: 	"strings"
27: 	"sync"
28: 	"time"
29: 
30: 	"github.com/Agent-Field/agentfield/sdk/go/harness"
31: 	invjsonschema "github.com/invopop/jsonschema"
32: 	tekjsonschema "github.com/santhosh-tekuri/jsonschema/v5"
33: 
34: 	"github.com/Agent-Field/SWE-AF/go/internal/fatal"
35: 	"github.com/Agent-Field/SWE-AF/go/internal/schemas"
36: )
37: 
38: // schemaCache memoizes the reflected JSON-schema map per concrete type T so the
39: // (non-trivial) reflection + marshal round-trip runs once per type. Keyed by
40: // reflect.Type; the stored map[string]any is treated as immutable by callers
41: // (the harness only ever marshals/reads it, never mutates), so sharing the
42: // cached value across goroutines is safe.
43: var schemaCache sync.Map // reflect.Type -> map[string]any
44: 
45: // schemaFor reflects T into the JSON-schema map the Go SDK harness consumes.
46: //
47: // How the SDK consumes this map (verified against sdk/go/harness/schema.go):
48: // the map is NOT used for programmatic validation — validity is defined purely
49: // by json.Unmarshal into the dest struct succeeding. The harness uses the map
50: // only to (a) embed a pretty-printed schema in the BuildPromptSuffix /
51: // BuildFollowupPrompt OUTPUT REQUIREMENTS instruction (harness/schema.go:36-74,
52: // :322-354) and (b) list expected top-level keys in DiagnoseOutputFailure
53: // (schema.go:306-318, which reads map["properties"]). Keys are alphabetized by
54: // json.MarshalIndent, so field ordering and `required` completeness are
55: // cosmetic. This means invopop's output — with $defs, items, and enum — is more
56: // than sufficient, and far richer than the SDK's own shallow StructToJSONSchema
57: // (which drops nested props/items/enums; design §2.3 says do NOT use it).
58: //
59: // Reflector configuration:
60: //   - ExpandedStruct: inline the root type's own properties at the top level so
61: //     map["properties"] is populated for DiagnoseOutputFailure (rather than a
62: //     bare $ref to $defs).
63: //   - DoNotReference=false (default): emit a $defs map for nested struct types.
64: //   - Anonymous: suppress the auto-generated $id derived from the package path.
65: func schemaFor[T any]() map[string]any {
66: 	t := reflect.TypeOf((*T)(nil)).Elem()
67: 	if cached, ok := schemaCache.Load(t); ok {
68: 		return cached.(map[string]any)
69: 	}
70: 
71: 	r := &invjsonschema.Reflector{
72: 		ExpandedStruct: true,  // root properties inline at top level
73: 		DoNotReference: false, // emit $defs for nested types
74: 		Anonymous:      true,  // no auto-generated $id from PkgPath
75: 	}
76: 	schema := r.ReflectFromType(t)
77: 
78: 	b, err := json.Marshal(schema)
79: 	if err != nil {
80: 		return map[string]any{}
81: 	}
82: 	var m map[string]any
83: 	if err := json.Unmarshal(b, &m); err != nil {
84: 		return map[string]any{}
85: 	}
86: 
87: 	// Make the reflected schema faithful to the source Pydantic model: the SDK
88: 	// now validates harness output against this map, so invopop's all-fields
89: 	// `required`, `additionalProperties:false`, and non-null Optional subschemas
90: 	// would OVER-reject valid output. (Enums are added on the enum types via
91: 	// JSONSchemaExtend, already present in `m`.)
92: 	m = schemas.MakePydanticFaithful(m, t.Name())
93: 
94: 	schemaCache.Store(t, m)
95: 	return m
96: }
97: 
98: // executeStructured owns SWE's structured-output reliability policy around the
99: // upstream AgentField harness call. Run[T] intentionally stays a thin
100: // integration seam so recovery behavior can evolve without spreading changes
101: // across role code or increasing rebase pressure on the base call path.
102: func executeStructured[T any](ctx context.Context, app HarnessCaller, prompt string, schema map[string]any, opts harness.Options) (*T, *harness.Result, error) {
103: 	// Weak OpenCode-backed models are materially more reliable when the SDK
104: 	// builds the structured envelope incrementally. Keep that policy here, next
105: 	// to recovery/validation, so role code and the base Run seam do not know how
106: 	// structured output is repaired.
107: 	if opts.Provider == "opencode" && opts.SchemaMode == "" {
108: 		opts.SchemaMode = "incremental"
109: 	}
110: 
111: 	var dest T
112: 	stopOutputCapture := startStructuredOutputCapture[T](ctx, opts.ProjectDir, schema)
113: 	result, err := app.Harness(ctx, prompt, schema, &dest, opts)
114: 	capturedOutput := stopOutputCapture()
115: 	if err != nil {
116: 		// A no-progress watchdog is a liveness signal, not proof that an already
117: 		// completed exact-schema result is unusable. Generic/provider errors remain
118: 		// fail-closed.
119: 		if strings.Contains(err.Error(), "CLI command made no progress") {
120: 			if recovered, ok := recoverStructuredResult[T](result, schema); ok {
121: 				return recovered, result, nil
122: 			}
123: 			if len(capturedOutput) > 0 {
124: 				var recovered T
125: 				if recoverErr := recoverStructuredText(string(capturedOutput), schema, &recovered); recoverErr == nil {
126: 					schemas.EmptyForNilSlices(&recovered)
127: 					result = normalizeRecoveredResult(result, &recovered)
128: 					return &recovered, result, nil
129: 				}
130: 			}
131: 		}
132: 		return nil, result, err
133: 	}
134: 
135: 	// Fatal API errors outrank schema fallback so billing/auth/provider failures
136: 	// cannot be hidden behind a default-seeded result.
137: 	if fErr := fatal.CheckFatalHarnessError(result); fErr != nil {
138: 		return nil, result, fErr
139: 	}
140: 
141: 	if result != nil && result.Parsed == nil &&
142: 		(result.FailureType == harness.FailureSchema || result.FailureType == harness.FailureNoOutput || result.FailureType == harness.FailureNone) {
143: 		if recovered, ok := recoverStructuredResult[T](result, schema); ok {
144: 			return recovered, result, nil
145: 		}
146: 	}
147: 
148: 	if result == nil || result.Parsed == nil {
149: 		seeded := seedDefaults[T]()
150: 		schemas.EmptyForNilSlices(&seeded)
151: 		return &seeded, result, nil
152: 	}
153: 
154: 	schemas.EmptyForNilSlices(&dest)
155: 	return &dest, result, nil
156: }
157: 
158: func recoverStructuredResult[T any](result *harness.Result, schema map[string]any) (*T, bool) {
159: 	if result == nil {
160: 		return nil, false
161: 	}
162: 	for _, text := range structuredResultCandidates(result) {
163: 		var recovered T
164: 		if err := recoverStructuredText(text, schema, &recovered); err != nil {
165: 			continue
166: 		}
167: 		schemas.EmptyForNilSlices(&recovered)
168: 		normalizeRecoveredResult(result, &recovered)
169: 		return &recovered, true
170: 	}
171: 	return nil, false
172: }
173: 
174: func normalizeRecoveredResult[T any](result *harness.Result, recovered *T) *harness.Result {
175: 	if result == nil {
176: 		result = &harness.Result{}
177: 	}
178: 	result.Parsed = recovered
179: 	result.IsError = false
180: 	result.ErrorMessage = ""
181: 	result.FailureType = harness.FailureNone
182: 	return result
183: }
184: 
185: // structuredResultCandidates returns only assistant text surfaces that can
186: // legitimately contain the final structured result. Tool outputs are excluded
187: // so repository/file JSON cannot be mistaken for the orchestration envelope.
188: func structuredResultCandidates(result *harness.Result) []string {
189: 	if result == nil {
190: 		return nil
191: 	}
192: 	seen := make(map[string]struct{})
193: 	out := make([]string, 0, 4)
194: 	appendUnique := func(text string) {
195: 		if text == "" {
196: 			return
197: 		}
198: 		if _, ok := seen[text]; ok {
199: 			return
200: 		}
201: 		seen[text] = struct{}{}
202: 		out = append(out, text)
203: 	}
204: 
205: 	appendUnique(result.Result)
206: 	for i := len(result.Messages) - 1; i >= 0; i-- {
207: 		msg := result.Messages[i]
208: 		kind, _ := msg["type"].(string)
209: 		if kind != "text" && kind != "assistant" && kind != "result" {
210: 			continue
211: 		}
212: 		if text, ok := msg["text"].(string); ok {
213: 			appendUnique(text)
214: 		}
215: 		if part, ok := msg["part"].(map[string]any); ok {
216: 			if text, ok := part["text"].(string); ok {
217: 				appendUnique(text)
218: 			}
219: 		}
220: 	}
221: 	return out
222: }
223: 
224: // startStructuredOutputCapture watches only output directories created after
225: // this invocation starts and keeps the latest exact-schema-valid output bytes
226: // in memory. The AgentField SDK removes its temporary output directory before
227: // returning a CLI no-progress error, so this narrow monitor lets SWE salvage a
228: // completed result without weakening validation. If more than one new output
229: // directory appears, capture becomes ambiguous and fails closed.
230: func startStructuredOutputCapture[T any](ctx context.Context, projectDir string, schema map[string]any) func() []byte {
231: 	if projectDir == "" || schema == nil {
232: 		return func() []byte { return nil }
233: 	}
234: 
235: 	pattern := filepath.Join(projectDir, ".agentfield-out-*")
236: 	existingDirs, _ := filepath.Glob(pattern)
237: 	existing := make(map[string]struct{}, len(existingDirs))
238: 	for _, dir := range existingDirs {
239: 		existing[dir] = struct{}{}
240: 	}
241: 
242: 	watchCtx, cancel := context.WithCancel(ctx)
243: 	var mu sync.Mutex
244: 	var latest []byte
245: 	ambiguous := false
246: 	done := make(chan struct{})
247: 
248: 	go func() {
249: 		defer close(done)
250: 		ticker := time.NewTicker(20 * time.Millisecond)
251: 		defer ticker.Stop()
252: 
253: 		scan := func() {
254: 			dirs, _ := filepath.Glob(pattern)
255: 			newDirs := make([]string, 0, 1)
256: 			for _, dir := range dirs {
257: 				if _, ok := existing[dir]; !ok {
258: 					newDirs = append(newDirs, dir)
259: 				}
260: 			}
261: 			if len(newDirs) > 1 {
262: 				mu.Lock()
263: 				ambiguous = true
264: 				latest = nil
265: 				mu.Unlock()
266: 				return
267: 			}
268: 			if len(newDirs) != 1 {
269: 				return
270: 			}
271: 
272: 			mu.Lock()
273: 			isAmbiguous := ambiguous
274: 			mu.Unlock()
275: 			if isAmbiguous {
276: 				return
277: 			}
278: 
279: 			b, err := os.ReadFile(filepath.Join(newDirs[0], ".agentfield_output.json"))
280: 			if err != nil || len(b) == 0 {
281: 				return
282: 			}
283: 			var recovered T
284: 			if err := recoverStructuredText(string(b), schema, &recovered); err != nil {
285: 				return
286: 			}
287: 			mu.Lock()
288: 			latest = append(latest[:0], b...)
289: 			mu.Unlock()
290: 		}
291: 
292: 		for {
293: 			scan()
294: 			select {
295: 			case <-watchCtx.Done():
296: 				return
297: 			case <-ticker.C:
298: 			}
299: 		}
300: 	}()
301: 
302: 	return func() []byte {
303: 		cancel()
304: 		<-done
305: 		mu.Lock()
306: 		defer mu.Unlock()
307: 		if ambiguous || len(latest) == 0 {
308: 			return nil
309: 		}
310: 		return append([]byte(nil), latest...)
311: 	}
312: }
313: 
314: // recoverStructuredText extracts candidate JSON objects with a string-aware
315: // scanner, validates them against the exact generated schema, and only then
316: // unmarshals into dest. Malformed or schema-invalid text stays a failure.
317: func recoverStructuredText[T any](text string, schema map[string]any, dest *T) error {
318: 	schemaBytes, err := json.Marshal(schema)
319: 	if err != nil {
320: 		return fmt.Errorf("marshal schema: %w", err)
321: 	}
322: 	compiler := tekjsonschema.NewCompiler()
323: 	if err := compiler.AddResource("mem://swe/schema.json", bytes.NewReader(schemaBytes)); err != nil {
324: 		return fmt.Errorf("add schema: %w", err)
325: 	}
326: 	compiled, err := compiler.Compile("mem://swe/schema.json")
327: 	if err != nil {
328: 		return fmt.Errorf("compile schema: %w", err)
329: 	}
330: 
331: 	for _, candidate := range extractJSONObjectCandidates(text) {
332: 		var data any
333: 		if err := json.Unmarshal([]byte(candidate), &data); err != nil {
334: 			continue
335: 		}
336: 		if err := compiled.Validate(data); err == nil {
337: 			if err := json.Unmarshal([]byte(candidate), dest); err == nil {
338: 				return nil
339: 			}
340: 		}
341: 
342: 		// CoderResult mirrors a Pydantic model where every field has a default.
343: 		// Let its custom UnmarshalJSON normalize its one known weak-model shape,
344: 		// then validate the fully materialized typed object. Other schemas stay
345: 		// strict; required fields are never synthesized generically.
346: 		if reflect.TypeOf((*T)(nil)).Elem().Name() == "CoderResult" {
347: 			var normalized T
348: 			if err := json.Unmarshal([]byte(candidate), &normalized); err != nil {
349: 				continue
350: 			}
351: 			normalizedBytes, err := json.Marshal(normalized)
352: 			if err != nil {
353: 				continue
354: 			}
355: 			var normalizedData any
356: 			if err := json.Unmarshal(normalizedBytes, &normalizedData); err != nil {
357: 				continue
358: 			}
359: 			if err := compiled.Validate(normalizedData); err != nil {
360: 				continue
361: 			}
362: 			*dest = normalized
363: 			return nil
364: 		}
365: 	}
366: 	return fmt.Errorf("no schema-valid JSON object found in final text")
367: }
368: 
369: // extractJSONObjectCandidates finds balanced top-level JSON objects while
370: // ignoring braces inside quoted strings. Larger candidates win over nested
371: // examples so the orchestration envelope is preferred.
372: func extractJSONObjectCandidates(text string) []string {
373: 	candidates := make([]string, 0, 2)
374: 	depth := 0
375: 	start := -1
376: 	inString := false
377: 	escaped := false
378: 
379: 	for i := 0; i < len(text); i++ {
380: 		ch := text[i]
381: 		if inString {
382: 			if escaped {
383: 				escaped = false
384: 				continue
385: 			}
386: 			if ch == '\\' {
387: 				escaped = true
388: 				continue
389: 			}
390: 			if ch == '"' {
391: 				inString = false
392: 			}
393: 			continue
394: 		}
395: 
396: 		switch ch {
397: 		case '"':
398: 			inString = true
399: 		case '{':
400: 			if depth == 0 {
```
_import/usage context:_ IMPORTS: import (
IMPORTED BY: none

### go/internal/harnessx/harnessx_test.go (showing first 400 of 510)
```
1: package harnessx
2: 
3: import (
4: 	"context"
5: 	"encoding/json"
6: 	"errors"
7: 	"os"
8: 	"path/filepath"
9: 	"strings"
10: 	"testing"
11: 	"time"
12: 
13: 	"github.com/Agent-Field/agentfield/sdk/go/harness"
14: 
15: 	"github.com/Agent-Field/SWE-AF/go/internal/fatal"
16: 	"github.com/Agent-Field/SWE-AF/go/internal/hitl"
17: 	"github.com/Agent-Field/SWE-AF/go/internal/schemas"
18: )
19: 
20: // --- test fixtures ----------------------------------------------------------
21: 
22: // mockHarness is the HarnessCaller seam the Python tests get by patching
23: // router.harness. It records what Run passed and returns a scripted result.
24: type mockHarness struct {
25: 	fn        func(ctx context.Context, prompt string, schema map[string]any, dest any, opts harness.Options) (*harness.Result, error)
26: 	gotOpts   harness.Options
27: 	gotSchema map[string]any
28: