### swe_af/execution/schemas.py
```diff
 )
 
 
-def _default_model_from_env() -> str | None:
-    """Pick a single model id from deployer env vars.
+def _default_model_from_env(runtime: str) -> str | None:
+    """Pick a single model id from deployer env vars, for ``runtime``.
 
     Cascades through the well-known env-var names this stack uses for model
     selection so the same Railway / docker-compose variable that points

 
         SWE_DEFAULT_MODEL  →  AI_MODEL  →  HARNESS_MODEL
 
+    ``HARNESS_MODEL`` is an OpenCode-ecosystem variable — it also feeds
+    OpenCode's ``small_model`` via config interpolation, and the Docker image
+    bakes a default value precisely so that interpolation always has one — so
+    it only participates in the cascade for the ``open_code`` runtime. Letting
+    it steer ``claude_code`` / ``codex`` pushed the image's baked
+    ``openrouter/…`` id into CLIs that cannot consume it, breaking every
+    non-OpenCode Docker deployment that didn't also set ``SWE_DEFAULT_MODEL``.
+
     Caller-supplied ``models={"default": …}`` and per-role overrides still
     beat the env value (see ``resolve_runtime_models`` precedence). All
     unset / empty → ``None``, which means "use the runtime base defaults".
     """
     for var in _DEFAULT_MODEL_ENV_VARS:
+        if var == "HARNESS_MODEL" and runtime != "open_code":
+            continue
         value = os.getenv(var, "").strip()
         if value:
             return value

     Precedence is inherited from ``resolve_runtime_models`` (highest first):
 
         1. ``SWE_MODEL_HIGH`` (planning reasoners are high-tier)
-        2. deployer env (``SWE_DEFAULT_MODEL`` → ``AI_MODEL`` → ``HARNESS_MODEL``)
+        2. deployer env (``SWE_DEFAULT_MODEL`` → ``AI_MODEL`` →
+           ``HARNESS_MODEL``, the latter only on ``open_code``)
         3. the runtime's own auto/base default:
              - ``codex``       → a codex-native model (never ``openrouter/…``)
              - ``open_code``   → the OpenRouter auto default (OpenRouter-only

     Resolution order (lowest → highest precedence):
         1. runtime base defaults (``_RUNTIME_BASE_MODELS[runtime]``)
         2. env-var cascade: ``SWE_DEFAULT_MODEL`` → ``AI_MODEL`` →
-           ``HARNESS_MODEL`` (first non-empty wins, applies to all roles)
+           ``HARNESS_MODEL`` (first non-empty wins, applies to all roles;
+           ``HARNESS_MODEL`` is consulted only on the ``open_code`` runtime —
+           see ``_default_model_from_env``)
         3. tier env vars: ``SWE_MODEL_LOW`` / ``SWE_MODEL_MED`` /
            ``SWE_MODEL_HIGH``, each applying to the roles in its tier
            (see ``ROLE_TO_TIER``)

         base = {field: _OPENROUTER_AUTO_DEFAULT_MODEL for field in base}
     resolved: dict[str, str] = {field: base[field] for field in field_names}
 
-    env_default = _default_model_from_env()
+    env_default = _default_model_from_env(runtime)
     if env_default:
         for field in field_names:
             resolved[field] = env_default
```

### go/internal/config/resolve.go
```diff
 
 // defaultModelFromEnv ports _default_model_from_env: first non-empty (stripped)
 // of SWE_DEFAULT_MODEL → AI_MODEL → HARNESS_MODEL, else "" (meaning None).
-func defaultModelFromEnv() string {
+//
+// HARNESS_MODEL is an OpenCode-ecosystem variable — it also feeds OpenCode's
+// small_model via config interpolation, and the Docker image bakes a default
+// value precisely so that interpolation always has one — so it is consulted
+// only for the open_code runtime. Letting it steer claude_code / codex pushed
+// the image's baked openrouter/… id into CLIs that cannot consume it.
+func defaultModelFromEnv(runtime string) string {
 	for _, v := range defaultModelEnvVars {
+		if v == "HARNESS_MODEL" && runtime != "open_code" {
+			continue
+		}
 		if value := envStripped(v); value != "" {
 			return value
 		}

 	if highModel := tierModelsFromEnv()["high"]; highModel != "" {
 		return highModel
 	}
-	if envModel := defaultModelFromEnv(); envModel != "" {
+	if envModel := defaultModelFromEnv(DefaultRuntime()); envModel != "" {
 		return envModel
 	}
 	if openRouterOnlyEnv() {

 		resolved[field] = base[field]
 	}
 
-	if envDefault := defaultModelFromEnv(); envDefault != "" {
+	if envDefault := defaultModelFromEnv(runtime); envDefault != "" {
 		for _, field := range fieldNames {
 			resolved[field] = envDefault
 		}
```

### swe_af/fast/schemas.py
```diff
     Resolution order (last wins):
       1. Runtime default (haiku or the shared open_code default, per runtime)
       2. Env cascade: ``SWE_DEFAULT_MODEL`` → ``AI_MODEL`` → ``HARNESS_MODEL``
+         (``HARNESS_MODEL`` only on the ``open_code`` runtime)
       3. ``models["default"]`` — overrides all roles
       4. ``models["<role>"]`` — overrides a specific role (pm, coder, verifier, git)
 

 
     resolved: dict[str, str] = {role: runtime_default for role in _FAST_ROLES}
 
-    # Deployer env cascade (SWE_DEFAULT_MODEL → AI_MODEL → HARNESS_MODEL), same
-    # as the main path — lets the variable that selects a model for the main
-    # node select it for fast builds too. Caller-supplied models still win.
+    # Deployer env cascade (SWE_DEFAULT_MODEL → AI_MODEL → HARNESS_MODEL, the
+    # latter only on open_code), same as the main path — lets the variable that
+    # selects a model for the main node select it for fast builds too.
+    # Caller-supplied models still win.
     from swe_af.execution.schemas import _default_model_from_env  # noqa: PLC0415
 
-    env_model = _default_model_from_env()
+    env_model = _default_model_from_env(config.runtime)
     if env_model:
         resolved = {role: env_model for role in _FAST_ROLES}
 
```

### go/internal/config/fastconfig.go
```diff
 
 // FastResolveModels ports fast_resolve_models — resolves the four role model
 // strings. Resolution order (last wins): runtime default → env cascade
-// (SWE_DEFAULT_MODEL → AI_MODEL → HARNESS_MODEL, same as the main path) →
+// (SWE_DEFAULT_MODEL → AI_MODEL → HARNESS_MODEL, the latter only on
+// open_code — same as the main path) →
 // models["default"] → models["<role>"]. An unknown key yields the verbatim
 // "Unknown role key" error.
 func FastResolveModels(config *FastBuildConfig) (map[string]string, error) {

 	// Deployer env cascade: lets the same variable that selects a model for
 	// the main node select it for fast builds too. Caller-supplied models
 	// (below) still win.
-	if envModel := defaultModelFromEnv(); envModel != "" {
+	if envModel := defaultModelFromEnv(config.Runtime); envModel != "" {
 		for _, role := range fastRoles {
 			resolved[role] = envModel
 		}
```