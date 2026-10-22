"""Manually authorized one-shot native SWE-AF Full launch, with readback."""
import json
import urllib.request
from pathlib import Path

root="http://127.0.0.1:18080"
def get(path):
    with urllib.request.urlopen(root+path,timeout=15) as r:
        return json.loads(r.read())
active=get("/api/v1/executions/active?limit=100")
if active.get("count") != 0:
    raise RuntimeError("Refusing duplicate launch: active root executions exist: "+
                       str([(r.get("root_execution_id"),r.get("target")) for r in active.get("runs",[])]))

models={
 "default":"gpt-6-luna",
 "pm":"gpt-6.1-sol", "architect":"gpt-6.1-sol",
 "tech_lead":"gpt-6.1-sol", "sprint_planner":"gpt-6-luna",
 "issue_writer":"gpt-6-luna", "coder":"gpt-6-luna",
 "qa":"gpt-6.1-sol", "code_reviewer":"gpt-6.1-sol",
 "qa_synthesizer":"gpt-6-luna", "issue_advisor":"gpt-6.1-sol",
 "replan":"gpt-6.1-sol", "retry_advisor":"gpt-6.1-sol",
 "verifier":"gpt-6.1-sol", "git":"gpt-6-luna",
 "merger":"gpt-6.1-sol", "integration_tester":"gpt-6.1-sol", "ci_fixer":"gpt-6.1-sol"
}
base_cfg={
 "runtime":"codex", "models":models,
 "max_coding_iterations":6,
 "agent_max_turns":40,
 "agent_timeout_seconds":3600,
 "max_retries_per_issue":2, "max_advisor_invocations":2,
 "enable_issue_advisor":True, "enable_replanning":True,
 "max_replans":3, "max_review_iterations":2,
 "deterministic_git":True,
 "enable_github_pr":False,
 "check_ci":False,
 "enable_learning":False,
 "enable_integration_testing":True,
 "max_verify_fix_cycles":2
}
jobs=[
 ("PR-AF",
  r"D:\Temp\swe-full-pr-af-20261020",
  10,
  "Complete the existing PR-AF project from its canonical PLAN.md and existing source, using BMad Quick Dev and no new architecture. The top goal is a real useful semantic PR review by the maintained Go PR-AF 17-reasoner implementation, with exact current AgentField DEV source/runtime identity, callback and SourceLoop readback. First inspect the live permanent DEV runtime and distinguish Go node from legacy Python node. PLAN.md states G0 BLOCKED until authorized CURRENT runtime readback; do not claim G0 passed or change production/deploy anything without that observation. Independently implement and test any purely local source/test improvements available now and capture precise blockers. Build a dependency-aware plan with independent tasks parallelized (up to 10), use native issue advisor and replan after failures, preserve all current project code, prioritize MVP. Real Git commits with scoped changes and independent QA. Do NOT push, publish a PR, post comments, send outbound messages, restart shared systems, or deploy. At the final result separate genuinely accepted functional E2E from blocked runtime gates; no fictitious success."),
 ("AF Skill Intelligence",
  r"D:\Temp\swe-full-skill-intelligence-20261020",
  10,
  "Bring existing AF Skill Intelligence from actual Go repository sources to MVP, with BMad Quick Dev, evidence-led tests, native independent QA. Preserve WIP hot retrieval AgentField adapter, its exact provenance/fail-closed boundary, and cherry-picked source fixes at current Git HEAD. Project is intelligence and trusted skill retrieval, NOT Registry semantic authority or a new orchestrator. Plan only missing dependency-ordered steps: repeatable Go hot-path indexed retrieval -> exact stable provenance/profile pin -> resolver -> durable AgentField execution -> correctly fenced replay-key and receipt, cold acquisition only when necessary to MVP. Use existing source, enable native parallel DAG, reviewer and advisor/replanner on errors. Run targeted and package Go tests offline, fail closed for missing Registry/AgentField runtime rather than fake acceptance. Commit only candidate-owned files, no cache, never mutate original dirty source, never promote Registry or call external providers without explicit authorization; no push or deployment."),
 ("CodeAF",
  r"D:\Temp\swe-full-codeaf-20261020",
  10,
  "Complete the CodeAF Windows Native MVP using existing Go source and latest preserved candidate WIP HEAD 45b41f2c, not a new factory/orchestrator. Existing TestMissionWindowsHeadlessCLI passed on A55: local provider HTTP -> observed edits[].old/new -> Go fixture source 41->42 -> go test -> JSON ok=true stop=done -> SQLite journal persisted. Blocking defect: windows filelock LockFileEx moved lock from byte 0 to 1<<40; this violates exclusion with OLD running binaries. Implement a Windows mixed-version byte-zero contention red test and minimal safe compatibility repair, preserve payload readability and active resident semantics, fail closed. Then rerun native E2E, Go filelock/lease/standing tests and affected package regression, classify failures against baseline. Prioritize real provider/tool/journal E2E, independent QA and correct accepted durable branch, not all ancillary modernization. Use native parallel DAG only on independent files; Issue Advisor and replanner on RED. Isolated workspace only. Do not touch original dirty CodeAF, do not push/deploy or send outside. Never assert CrewSync SERVER_DONE without actual authorized server readback.")
]
profiles=json.loads((Path(__file__).resolve().parent/"project-profiles.json").read_text(encoding="utf-8"))
project_names={"PR-AF":"pr-af","AF Skill Intelligence":"af-skill-intelligence","CodeAF":"codeaf"}
wanted=profiles["full_dag_issue_limits"]
for name,_,slots,_ in jobs:
    if slots != wanted[project_names[name]]:
        raise RuntimeError(f"Full DAG slot mismatch for {name}: script={slots} config={wanted[project_names[name]]}")
if sum(job[2] for job in jobs) > profiles["initial_limits"]["max_parallel_issue_roots"]:
    raise RuntimeError("Configured concurrent Full DAG width exceeds site limit")
result=[]
for name,repo,slots,goal in jobs:
    if not Path(repo).is_dir():
        result.append({"project":name,"error":"worktree missing"})
        continue
    cfg=dict(base_cfg)
    cfg["max_concurrent_issues"]=slots
    payload={"input":{
      "goal":goal,"repo_path":repo,
      "artifacts_dir":".artifacts",
      "additional_context":(
        "A55, native Full Plan -> DAG Execute -> QA/Verify. "
        "Before issuing parallel work, identify and freeze public API and data contracts "
        "(function signatures, input/output types, wire schemas, migration/order constraints) "
        "with compile-time or focused contract tests. "
        "Use the built-in PlannedIssue depends_on/provides, target_repo, "
        "files_to_create/files_to_modify, acceptance_criteria and touches_interfaces. "
        "Partition into independent file/package owners, maximize READY DAG width up to "
        f"{slots} for this project, serialize shared interfaces, schema migrations, "
        "generated code and any touching files; do not promise 30 READY tasks if "
        "dependencies prohibit them. Existing accepted Git SHA and BMad/CrewSync "
        "sources win, no rewrite and no duplicate work. No external side effects, "
        "no auto-push/deploy/PR, only independently verified useful progress counts."
      ),
      "config":cfg
    }}
    req=urllib.request.Request(root+"/api/v1/execute/async/swe-factory-30.build",
        data=json.dumps(payload).encode("utf-8"),
        headers={"Content-Type":"application/json"},method="POST")
    try:
        with urllib.request.urlopen(req,timeout=25) as r:
            data=json.loads(r.read())
        result.append({"project":name,"repo":repo,**data})
    except Exception as exc:
        result.append({"project":name,"repo":repo,"error":str(exc)})
print(json.dumps(result,indent=2))
