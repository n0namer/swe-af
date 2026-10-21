# SWE — единственная операторская команда на VP49

SWE.cmd на рабочем столе вызывает SWE.ps1. Это только маршрутизация; источником состояния исполнений остаётся AgentField на 127.0.0.1:18080.

## Запуск
- По умолчанию FULL: SWE.cmd -Goal "Доведи проект до GREEN" -RepoPath "C:\path\to\git-repo"
- Явный FULL: SWE.cmd -Mode full -Goal "..." -RepoPath "..."
- Быстрый режим: SWE.cmd -Mode fast -Goal "..." -RepoPath "..."
- Из GitHub: SWE.cmd -Mode full -Goal "..." -RepoUrl "https://github.com/org/repo"
- Для локального теста без PR: добавить -NoPR
- Состояние: SWE.cmd -ExecutionId exec_...
- Активные: SWE.cmd -List

FULL (swe-planner.build): PM -> architecture -> tech lead -> sprint planner -> DAG -> issue advisor / replanner -> review/QA -> merge -> verify. Таймаут агента 90 минут по умолчанию (передаётся в build config); max concurrent issues 4 на уровне и replans 3, ограничены параметрами. Вызовы Codex остаются в workspace-write sandbox.

FAST (swe-fast.build): плоский короткий план и последовательное исполнение. Это альтернативный режим, а не другой глобальный контроллер. Внутренние node IDs не переименовывать: на них ссылаются AgentField lineage, history, restart и SDK.

Нельзя автоматически повторять POST при timeout: сначала искать существующее execution по идентичности. Receipts отправленных через SWE команд хранятся в .agentfield/ops/swe-submissions.

Обязательный финальный readback: смотреть и AgentField status, и output_data.success, и Git/test receipts. HTTP 200 / succeeded не равны принятым результатам.

## Автовосстановление
Штатный Windows Task Scheduler запускает Control Plane и зарегистрированные узлы, рестарт процесса не является resume DAG. Отдельный safe recovery watch, если отсутствует, должен настраивать владелец SWE-AF Runtime/Recovery; не создавать второй supervisor. AgentField -- источник состояния execution IDs и replay. Go DAG владеет checkpoint/READY/lease/fencing.
