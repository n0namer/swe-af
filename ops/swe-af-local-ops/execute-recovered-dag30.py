"""Operator-approved native SWE-AF execute from preserved PM PRDs and architecture.
This is one-shot issue planning; native SWE-AF owns DAG, worktrees, QA, replans.
"""
import json
import subprocess
import sys
import urllib.request
from pathlib import Path
ROOT="http://127.0.0.1:18080"
OPS=Path(__file__).parent
PROFILES=json.loads((OPS/"project-profiles.json").read_text(encoding="utf-8"))
sys.path.insert(0,PROFILES["swe_af_code_root"])
from swe_af.reasoners.schemas import PlanResult
from swe_af.reasoners.pipeline import (
    _compute_levels, _serialize_file_conflicts, _validate_file_conflicts
)

def get(path):
    with urllib.request.urlopen(ROOT+path,timeout=20) as r:
        return json.load(r)

def post(path, body):
    req=urllib.request.Request(ROOT+path,method="POST",data=json.dumps(body).encode(),
                               headers={"Content-Type":"application/json"})
    with urllib.request.urlopen(req,timeout=25) as r:
        return json.load(r)

def issue(name,title,description,owned,new,tests,accept,needs=(),interfaces=False):
    return {
      "name":name,"title":title,"description":description,
      "acceptance_criteria":accept,"depends_on":list(needs),
      "provides":[name],"estimated_complexity":"medium",
      "files_to_modify":owned,"files_to_create":new,
      "testing_strategy":tests,
      "guidance":{"needs_new_tests":True,"estimated_scope":"medium",
          "touches_interfaces":interfaces,"needs_deeper_qa":True,
          "testing_guidance":tests,
          "review_focus":"Real Windows/local Go regression and exact owned-file scope",
          "risk_rationale":"MVP correctness and fail-closed state"}
    }

JOBS=[
  {"id":"codeaf","repo":r"D:\Temp\swe-full-codeaf-20261020",
   "pm":"exec_20261020_020557_1qlkfcw7","build_id":"full30-codeaf-r1",
   "architecture":"Use existing Go modular monolith and Windows LockFileEx filelock boundary. Existing dual-range design from .artifacts/plan/architecture.md is a candidate, not accepted. Legacy-byte-zero and new-range simultaneous ownership must be impossible; preserve exact resident caller semantics and headless provider->tools->journal E2E.",
   "issues":[
     issue("win32-mixed-version-lock-compat",
       "Repair Windows mixed-version LockFileEx exclusion",
       "At candidate 45b41f2c existing filelock_windows.go locks 1<<40, whereas historical binaries lock byte 0. RED before GREEN: use real Windows LockFileEx byte-zero holder in a new test; new Lock must return busy on another handle. Repair with minimal compatible protocol that retains both old/new exclusion, handles partial acquisition rollback, unlock, shared/exclusive semantics, and payload read safety. Preserve existing code and never invent new lock service; do not change schema or provider. Architecture document includes dual-range candidate and Windows mapping caveats; implement only proven safe design.",
       ["internal/filelock/filelock_windows.go","internal/filelock/filelock_windows_test.go"],
       ["internal/filelock/filelock_windows_compat_test.go"],
       "go test ./internal/filelock -count=1 -timeout 120s; mixed-version actual lock red/green on Windows; go vet ./internal/filelock.",
       ["Real byte-zero legacy process lock excludes new exclusive owner","New-vs-new reader/writer semantics preserved","No partial lock leak; payload/metadata semantics correctly handled","Only owned files in commit"],interfaces=True),
     issue("resident-lease-exact-compat",
       "Validate and minimally repair resident lifecycle with interoperable filelock",
       "After compatible filelock passes, use existing resident.go protocol to confirm metadata observation, active lease ownership, tick/release and mixed-version admission behavior. Fail closed on unknown holder, never permit dual ownership. Keep deterministic cancellation and no stale-owner guessing. Add regression only when missing; preserve existing source and no production deploy.",
       ["internal/lease/resident.go","internal/lease/resident_test.go"],
       ["cmd/codeaf/mission_windows_resident_test.go"],
       "go test ./internal/filelock ./internal/lease -count=1 -timeout 180s; go test ./cmd/codeaf -run '^TestDoFailsFastWhenAnotherProcessHoldsTheResidentLock$|^TestMissionWindowsHeadlessCLI$' -count=1 -timeout 240s.",
       ["Mixed-version lease never creates two live owners","Resident readback and release stable under Win32 locking","Actual headless Windows E2E remains green"],needs=["win32-mixed-version-lock-compat"],interfaces=True)
  ]},
  {"id":"af-skill-intelligence","repo":r"D:\Temp\swe-full-skill-intelligence-20261020",
   "pm":"exec_20261020_020557_w58bvzv8","build_id":"full30-skill-r1",
   "architecture":"Preserve current Go package boundaries and exact trusted skill provenance. Existing WIP c821400 has hot adapter, cloning, replay-key patch. Architecture file explicitly marks external C1 Registry/Profile and C2 AgentField target contracts blocked. Independently implement only local in-process, version-exact testable pieces; no invented external routes or Registry authority. Durable reservation/commit is local host only.",
   "issues":[
     issue("hot-receipt-durable-local-fence",
       "Implement local durable idempotent receipt reservation",
       "Build a small local file-backed, crash-restart recoverable reservation and terminal receipt mechanism for existing skillindex AgentField hot execution seam, using local Go stdlib only. Must preserve exact binding-key+request-key, distinguish in-progress from completed and fail closed on corruption/conflict. Do not invent external Registry/AgentField contract, new scheduler or promoted capability; use existing hot_agentfield types, keep separate new files and tests, and no external calls.",
       [],
       ["skillindex/hot_receipt_store.go","skillindex/hot_receipt_store_test.go"],
       "go test ./skillindex -count=1 -timeout 120s; add subprocess/restart and key-mismatch negative regression; targeted go vet.",
       ["Exact same binding key/request key does not cause second effect","Different binding or conflicting receipt rejected","Crash/restart preserves terminal and in-progress fences","No external side effects; clear limitations"],interfaces=True),
     issue("hot-exact-profile-contract-check",
       "Pin local profile/skill identity and fail-closed input contracts",
       "In local skillindex library use existing ExactReference, TrustedSkillPin and query/profile types. Provide minimal local validation and regression for exact immutable profile-skill binding, owner/source/trust evidence and stale/forged input rejection. Do not create a fake authoritative Registry, invent live remote endpoint, or promote skill. Other new files must be disjoint from receipt store; if exact owner contract absent, return typed BLOCKED_CONTRACT, not fabricated allowed.",
       [],
       ["skillindex/hot_profile_contract.go","skillindex/hot_profile_contract_test.go"],
       "go test ./skillindex -count=1 -timeout 120s; tests fail closed on nonexact/mismatched/missing owner/profile evidence.",
       ["No mutable aliases or unverified source becomes eligible","Existing hot adapter/index behavior preserved","Block unproven external membership claims"],interfaces=True),
     issue("hot-native-integration-qa",
       "Integrate verified local fences with existing AgentField seam",
       "After local durable fence and profile contract artifacts pass QA, minimally wire their existing interfaces into hot_agentfield.go while preserving c821400 exact execution receipt key and clone isolation. One permitted local in-process or loopback fixture verifies selection -> exact resolution -> one submission -> validated terminal receipt and restart replay. Fail closed for missing C1/C2 remote contracts; no external HTTP calls without permission, no Registry promotion.",
       ["skillindex/hot_agentfield.go","skillindex/hot_agentfield_native_mvp_test.go"],
       ["skillindex/hot_native_integration_test.go"],
       "go test ./skillindex ./authoring ./acquisition -count=1 -timeout 180s; independent adversarial replay-key QA.",
       ["At most one local effect on same exact binding","Invalid/stale receipt refused","No fake success on unknown registry or deployment","Evidence identifies local fixture not real external DEV"],needs=["hot-receipt-durable-local-fence","hot-exact-profile-contract-check"],interfaces=True)
   ]},
  {"id":"pr-af","repo":r"D:\Temp\swe-full-pr-af-20261020",
   "pm":"exec_20261020_020556_ava789lu","build_id":"full30-pr-r1",
   "architecture":"PR-AF current North Star depends on Go permanent DEV runtime identity, not local source alone. PLAN.md is canonical and G0 remains blocked, no deploy or external messages. Local Go node is 17 reasoners, with typed SDK review and trace. Implement independent local negative-path regression tasks only; preserve existing source and fail closed on missing SDK cache or CURRENT runtime evidence.",
   "issues":[
     issue("pr-go-node-local-runtime-contract",
       "Check maintained PR-AF Go node registration and callback contract",
       "Using existing Go node registration and package manifest, add/repair smallest local deterministic regressions proving maintained Go 17-reasoner registration and callback wiring. Distinguish Go from legacy Python local-path install; do not alter shared runtime or fake installed source identity. If SDK v0.1.130 unavailable, record typed dependency blocker instead of modifying go.mod or downloading without authorization. Only edit files in go/internal/node.",
       ["go/internal/node/register_test.go","go/internal/node/node_test.go"],
       ["go/internal/node/runtime_identity_test.go"],
       "cd go && go test ./internal/node -count=1 -timeout 120s; report missing external module as BLOCKED_INFRA, no fake green.",
       ["Expected reasoner registration and routed callback exercised locally","Go/Python selector ambiguity clearly represented","No service mutation or promoted runtime identity"],interfaces=True),
     issue("pr-go-review-local-failclosed",
       "Strengthen local PR review result fail-closed behavior",
       "Inspect existing orchestrator tests and only repair source-level genuine defects in failure propagation, empty/no evidence result and bounded partial-review finalization. Preserve typed API and existing orchestrator, no new provider or SDK; add focused adversarial unit regression without posting PR comments or remote GitHub mutations.",
       ["go/internal/orch/output.go","go/internal/orch/output_test.go"],
       ["go/internal/orch/review_failure_gate_test.go"],
       "cd go && go test ./internal/orch -count=1 -timeout 120s; report offline module limitation separately.",
       ["Empty or contradictory evidence cannot be emitted as positive useful review","Actual partial/failure result is truthful and bounded","No external outbound call or changed public schema"],interfaces=True)
   ]}
]

models={"default":"gpt-6-luna","coder":"gpt-6-luna","qa":"gpt-6.1-sol",
  "code_reviewer":"gpt-6.1-sol","qa_synthesizer":"gpt-6-luna",
  "issue_advisor":"gpt-6.1-sol","replan":"gpt-6.1-sol",
  "retry_advisor":"gpt-6.1-sol","git":"gpt-6-luna",
  "verifier":"gpt-6.1-sol","merger":"gpt-6.1-sol",
  "integration_tester":"gpt-6.1-sol"}

active=get("/api/v1/executions/active?limit=100")
if active.get("count"):
    raise RuntimeError("Existing active runs: refuse duplicate tasks")
results=[]
for job in JOBS:
    repo=job["repo"]
    if not Path(repo).is_dir(): raise RuntimeError("Missing worktree "+repo)
    baseline=subprocess.check_output(["git","-C",repo,"rev-parse","HEAD"],text=True).strip()
    branch=subprocess.check_output(["git","-C",repo,"symbolic-ref","--short","HEAD"],text=True).strip()
    pm=get("/api/ui/v1/executions/"+job["pm"]+"/details")
    if pm.get("status")!="succeeded": raise RuntimeError("PM not succeeded")
    prd=pm.get("output_data")
    arch_text=(Path(repo)/".artifacts"/"plan"/"architecture.md").read_text(encoding="utf-8")
    architecture={
      "summary":job["architecture"]+" Source architecture: "+str(Path(repo)/".artifacts"/"plan"/"architecture.md"),
      "components":[{"name":"Existing library/runtime","responsibility":"Implement bounded local issue-owned code using existing project modules","touches_files":[],"depends_on":[]}],
      "interfaces":["Use exact existing Go source signature and schema; no invented external provider or registry contracts"],
      "decisions":[{"decision":"MVP first: bounded independent local proof","rationale":"External runtime and contract gates remain unverified; no fake acceptance"}],
      "file_changes_overview":", ".join(sorted({f for iss in job["issues"] for f in iss["files_to_modify"]+iss["files_to_create"]})),
    }
    issues=json.loads(json.dumps(job["issues"]))
    edges=_serialize_file_conflicts(issues)
    levels=_compute_levels(issues)
    conflicts=_validate_file_conflicts(issues,levels)
    if conflicts:raise RuntimeError("Unsafe conflicts "+str(conflicts))
    plan=PlanResult.model_validate({
      "prd":prd,"architecture":architecture,
      "review":{"approved":False,
        "feedback":"Architecture draft is unaccepted. Operator explicitly authorized bounded local test-led implementation only; external live gates remain blocked.",
        "summary":"Local bounded implementation authorized; do not claim whole-project approval"},
      "issues":issues,"levels":levels,"file_conflicts":conflicts,
      "artifacts_dir":str(Path(repo)/".artifacts"/"native30-execute"),
      "rationale":"Recovered from successful actual AgentField PM outputs, checked architecture.md, and explicit local issue file ownership. No repeated PRD or Architect generation; no external side effects."
    }).model_dump(mode="json")
    records=OPS/"plans"
    records.mkdir(exist_ok=True)
    (records/f"native30-{job['id']}.json").write_text(json.dumps(plan,indent=2),encoding="utf-8")
    config={
      "runtime":"codex","models":models,"max_concurrent_issues":PROFILES["full_dag_issue_limits"][job["id"]],
      "max_coding_iterations":6,"agent_max_turns":40,"agent_timeout_seconds":3600,
      "max_retries_per_issue":2,"max_advisor_invocations":2,
      "enable_issue_advisor":True,"enable_replanning":True,"max_replans":3,
      "deterministic_git":True,"enable_integration_testing":True,
      "check_ci":False,"enable_learning":False
    }
    payload={"input":{
      "plan_result":plan,"repo_path":repo,
      "git_config":{"integration_branch":branch,"original_branch":branch,
                    "initial_commit_sha":baseline,"mode":"existing",
                    "remote_url":"","remote_default_branch":""},
      "config":config,"build_id":job["build_id"],"resume":False
    }}
    results.append({"project":job["id"],"tasks":len(issues),
      "ready":len(levels[0]),"levels":levels,"edges":len(edges),
      "sha":baseline,"payload":payload})
print("READY",json.dumps([{k:v for k,v in row.items() if k!="payload"} for row in results]))
if "--submit" in sys.argv:
    for row in results:
        data=post("/api/v1/execute/async/swe-factory-30.execute",row["payload"])
        print("DISPATCH",json.dumps({"project":row["project"],**data}))
