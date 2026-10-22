# SWE-AF / AgentField MVP Critical Path — 2026-10-20

Authority: proposed operator backlog, NOT confirmed CrewSync records.
Source: AgentField root executions run_20261020_110826_9ue9ofsh, run_20261020_110826_n9zh5oyq, run_20261020_110826_4t42h1x5; Git and QA readbacks.
North Star: three native MVPs with independent evidence, accepted CrewSync SERVER_DONE for tasks registered there. No push/deploy/publication.

## Priority work DAG (12 goals)
| Key | Dependency | Goal / DoD | Starting evidence |
|---|---|---|---|
| F1 | — | SWE-AF output bridge independent QA; context isolation, schema retries, no duplicate effects; accepted SHA and same-code replay | candidate a28969d246b2c155c739dbec35f147be08681e92, 19 pytest pass; external reviewer failed 429 |
| F2 | F1 | Apply accepted bridge to Full30, verify live CLI flags and structured result; root max duration configured 21600s; no orphan writes | launch env 21600 tested, new node pid 320728; old root errors 'Execution cancelled by upstream client' at 7200s |
| F3 | F2 | Resume/recover existing plans without re-running already accepted work; validated checkpoint + same build ID / Git + no old writers; prove terminate/restart/readback | all 3 former roots terminal, active=0 at 13:31 |
| C1 | — | Independent native Windows mixed-version filelock QA; historical byte-zero contention + release/shared/rollback; exact SHA | recovered 7571105caa460b3ddeb06d0a88bd295ada992519, filelock go test/vet PASS |
| C2 | C1 | Resident lifecycle/native headless end-to-end codeaf do -> provider -> 41 to 42 -> tests -> SQLite journal, negative controls; terminal evidence; approved SHA | candidate existing Windows E2E, requires revalidation |
| C3 | C2,F2 | Merge allowed CodeAF scopes, rerun real Windows E2E on final integrated SHA, independent QA and CrewSync acceptance where available | integration still 45b41f2c |
| S1 | — | Independently adversarial QA durable local receipt: binding, malformed UTF8/duplicate keys, foreign commit, restart; tests on exact commit | 2d5be92b476714c7813d44120e4182b6bb7c74f0, full go tests + vet pass |
| S2 | — | Resolve actual C1 Registry/Profile authority and C2 AgentField target contract; typed BLOCKED_CONTRACT until verified; no mock acceptance | 4 failing contract regressions on untrusted owner/profile |
| S3 | S1,S2 | Integrate skill lookup/exact pin/real AgentField submit/validated terminal receipt/replay, no double effect; independent native E2E and SHA | external authority still missing |
| P1 | — | Pin Go SDK v0.1.130 and JWT v5.3.1 via authorized isolated cache; run node/orch/full Go tests, no version substitution | JWT acquired, SDK fetch TLS timed out |
| P2 | P1 | Fix recursive review coverage counters, truthful partial failure, focused adversarial tests; independent QA on integrated SHA | blocking review failure Attempted > planned |
| P3 | P1,P2 | Verify installed persistent DEV Go runtime identity, 17 reasoners, callbacks, SourceLoop and true nonpublishing PR review E2E; independent QA + SHA | current DEV identity unverified, G0 blocked |

Cross-cutting acceptance: L1 terminal actual output_data; L2 durable scoped Git SHA; L3 tests integrated final SHA with independent QA; L4 exact DoD and CrewSync SERVER_DONE readback (when CrewSync tasks exist). No fabricated metrics; remaining external blockers stay BLOCKED.

Execution order: W0 F1,F2,F3 and concurrent C1,S1,P1,S2; W1 C2,P2; W2 C3,S3,P3. W0 contracts are local gates, not global barrier. Keep tasks in original project attempts and worktrees. Do not repeat build/PM/Architect.

Citable practice references: https://dora.dev/capabilities/continuous-integration/ ; https://dora.dev/capabilities/working-in-small-batches/ ; https://github.com/Agent-Field/SWE-AF/blob/main/docs/ARCHITECTURE.md

## Current state
- F1 candidate, not independent QA accepted.
- F2 config implemented/restarted, CLI schema flags not live-E2E verified.
- F3 not yet verified.
- C1 recovered and local tests pass, independent QA missing.
- S1 local tests pass, independent QA missing.
- P1 partially blocked by SDK download.
- CrewSync network/readback unverified; do not count server DONE.
