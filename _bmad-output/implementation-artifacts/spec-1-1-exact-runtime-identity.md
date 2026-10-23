---
title: 'VP49 SWE-AF exact loaded runtime identity'
type: 'bugfix'
created: '2026-10-23'
status: 'draft'
review_loop_iteration: 0
context:
  - 'AGENTS.md'
  - 'PLAN.md'
  - '_bmad-output/planning-artifacts/epics-sweaf-mvp.md'
---

<frozen-after-approval reason="human-owned intent — do not modify unless human renegotiates">

## Intent

**Проблема:** На VP49 исходный `swe-fast` был успешно протестирован в интеграционном worktree и опубликован как `52b3fddc844fc3b5d7e5f0943bf1b19aaf78201b`, но изменение пути Windows Scheduled Task ещё не доказывает, что именно этот код загружен в живой Python-процесс. Предыдущий Fast E2E создал правильный файл без Git-коммита; `output_data.success=false`. Запуск нового E2E до подтверждения идентичности runtime даёт недостоверную диагностику.

**Подход:** Использовать существующие Windows Scheduled Task `VP49-AgentField-swe-fast`, launcher `~/.agentfield/ops/start-swe-fast.ps1`, `ops/windows/SWE.ps1` и AgentField API. Поставить один безопасный read-only preflight, затем при необходимости мягко перезапустить **только существующий** владеющий процесс, подтвердить его загруженный каталог/коммит/среду/регистрацию и записать доказательства. Не создавать вторую службу или альтернативный runner.

## Boundaries & Constraints

**Always:** Гарантировать отсутствие активных AgentField execution перед рестартом, чистый целевой Git worktree и точный SHA; сохранять текущую конфигурацию моделей (основная `gpt-6.1-sol`, агенты `gpt-6-luna`) и retry-политику; получать реальные PID, callback/node identity, связку launched-path→Python process и post-restart health/readback. В случае неуспеха вернуть прежний launcher и исходную Windows task без одновременных писателей. Защитить чужие dirty/untracked файлы.

**Ask First:** Только если потребуется остановить чужую задачу, заменить секреты или удалить существующие полезные worktree/ветки.

**Never:** Новый supervisor или контроллер, новый сервис, force push/reset, PR, deploy, CI/CD, изменение чужих репозиториев, вывод токенов, запуск повторного E2E при неоднозначном результате прежнего execution.

## I/O & Edge-Case Matrix

| Scenario | Input / State | Expected Output / Behavior | Error Handling |
|---|---|---|---|
| Активный узел на нужном SHA | Один `swe-fast`, источник `52b3fdd` | Не перезапускать; сохранить PID, SHA, регистрацию | При недоказуемой идентичности: `EVIDENCE_MISSING` |
| Устаревший код, нет jobs | `swe-fast` из старого worktree, `active.count=0` | Перезапустить ту же Scheduled Task; один listener и нужный SHA | Post-state readback, затем rollback при неготовности |
| Есть активные jobs | `active.count>0` | Не останавливать процесс | Gate BLOCKED с реальными execution ID |
| Restart/callback failure | Нет listener или реестрового узла | Не допустить второго процесса, восстановить старый launcher | FAIL с PID, Task state и diagnostic receipt |

</frozen-after-approval>

## Code Map

- `C:\Users\Рафик\.agentfield\ops\start-swe-fast.ps1` — существующий launcher: путь checkout, модель, порт и Python venv.
- `ops/windows/SWE.ps1` — единственный операторский CLI для execute/status; не заменять.
- `swe_af/fast/app.py` — функция build и fail-closed success gate, текущий точный кандидат `52b3fdd`.
- `docs/RECOVERY-A55-VP49-20261023.md` — подтверждения восстановления и известный красный E2E.
- AgentField `/health`, `/api/v1/executions/active`, `/api/ui/v1/nodes/summary` — существующие readback endpoints.

## Tasks & Acceptance

**Execution:**
- [ ] Verify local integration SHA `52b3fdd`, clean worktree, launcher path, Scheduled Task owner PID, OpenCodex `10100` and AgentField `18080`.
- [ ] Query active executions. When zero and loaded code differs, restart only `VP49-AgentField-swe-fast`, bounded wait, identify new listener and node registration.
- [ ] Save minimal evidence into existing `_bmad-output/implementation-artifacts` and update sprint/CrewSync only after DoD; no parallel PLAN.md.

**Acceptance Criteria:**
- Given no active jobs and a clean checkout, when the existing SWE node starts, then one registered `swe-fast` maps to the intended source SHA and existing callback with no competing writer.
- Given a pending execution or inconclusive owner identity, when reconciliation runs, then no mutation occurs and the gate returns BLOCKED/EVIDENCE_MISSING.
- Given an unsuccessful restart, when verifying post-state, then the old launcher/owner can be restored without data loss.

## Spec Change Log

## Verification

**Commands:**
- Existing `C:\Temp\SWEAF-Before-Switch-Audit.ps1` or equivalent read-only process + AgentField probe — exact ownership and zero active jobs.
- Existing `C:\Temp\SWEAF-Restart-Existing-Fast.ps1` **only after independent preflight** — status/PID/node readback.
- `python -m pytest -q tests/fast/test_app.py tests/test_native_git_init.py` on Python 3.12, using existing venv — tests must remain GREEN (latest recorded 30/30 PASS).

**Manual checks:** Confirm new process loaded the intended checkout; a matching file on disk alone is insufficient.
