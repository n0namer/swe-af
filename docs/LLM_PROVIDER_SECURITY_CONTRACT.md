# SWE-AF — LLM Provider & Security Profile

Status: canonical repo-local provider profile
Cross-component contract: `n0namer/universal-solver:main/docs/architecture/llm-provider-security-contract.md`

## Scope

This profile defines SWE-AF-specific provider/model/runtime behavior. It does not replace SWE-AF `PLAN.md` and does not contain secret values.

## Current provider surfaces

SWE-AF supports coding-runtime/model selection separately from raw credential presence. Relevant runtime/provider variables include:

- `OPENAI_API_KEY`
- `OPENAI_BASE_URL`
- `OPENROUTER_API_KEY`
- `ANTHROPIC_API_KEY`
- `SWE_DEFAULT_RUNTIME`
- `SWE_DEFAULT_MODEL`

The permanent AgentField DEV topology may inject OpenAI-compatible Gonka key/base values, but source/env presence alone is not CURRENT proof that a coding call used Gonka.

## Required OpenAI-compatible / Gonka pattern

When SWE-AF is intentionally routed through an OpenAI-compatible Gonka provider, the intended shape is:

```text
SWE_DEFAULT_RUNTIME = open_code
SWE_DEFAULT_MODEL = openai/<Gonka model>
OPENAI_API_KEY = <Gonka secret from runtime secret store>
OPENAI_BASE_URL = <Gonka OpenAI-compatible endpoint>
```

The model/runtime selection MUST remain explicit. Merely setting `OPENAI_API_KEY` is insufficient if runtime resolution can choose another harness/provider.

## Bootstrap and manifest contract

Bootstrap admission and actual provider routing are separate gates.

- The package/startup contract MUST admit the provider paths SWE-AF actually supports.
- A legacy preflight that requires Anthropic/OpenRouter while OpenAI-compatible routing is supported is `BOOTSTRAP_ADMISSION` drift.
- Satisfying a startup gate with a compatibility credential is not proof that the coding call used that provider.

## Harness propagation

The selected coding runtime must receive the intended provider environment end-to-end.

For OpenAI-compatible routing, verify that the actual coding subprocess receives both:

- `OPENAI_API_KEY`
- `OPENAI_BASE_URL`

and that the selected model remains `openai/<...>` through the final coding call.

## Security requirements

- Never commit or log credential values.
- Runtime logs/evidence may record variable presence, provider/model identity, and redacted endpoint identity only.
- Do not use a real fallback-provider token merely to satisfy a legacy preflight unless that provider is intentionally part of policy.
- If a compatibility shim is ever used for bootstrap, semantic acceptance MUST prove no unintended fallback consumed it.
- Coding task repository credentials are separate from LLM provider credentials and must be evaluated separately.

## Acceptance ladder

1. Exact SWE-AF source/runtime identity known.
2. Package/startup admits the intended provider path.
3. `swe-planner` active/ready.
4. Runtime resolves to the intended coding harness.
5. Model resolves to the intended OpenAI-compatible model namespace.
6. `OPENAI_API_KEY` + `OPENAI_BASE_URL` reach the actual coding subprocess.
7. Minimal reasoner/coding call succeeds.
8. Execution evidence shows no unintended OpenRouter/Anthropic fallback.
9. `implement_issue` (or current canonical implementation canary) succeeds on a bounded task with result inspection.

Node health alone is not provider or semantic PASS.

## Failure classes to use

- `BOOTSTRAP_ADMISSION`
- `MODEL_RESOLUTION`
- `ENV_PROPAGATION`
- `BASE_URL_LOSS`
- `AUTH`
- `FALLBACK`
- `TRANSPORT`
- `SEMANTIC`

Patch the first failing layer only and preserve source/runtime identity evidence.
