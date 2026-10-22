"""Conservative native AgentField/SWE-AF recovery watch; runs under Windows task.
Auto-resumes only timeout/failed transport executions with saved native checkpoint.
Never resumes a logical coding failure or crosses active owner / dirty tracked repo.
"""
import json
import os
import subprocess
import time
from pathlib import Path
from safe_resume_admission import guarded_payload
from urllib.request import Request, urlopen

BASE="http://127.0.0.1:18080"
OPS=Path(__file__).resolve().parent
LOG=OPS/"logs"/"full30-recovery.ndjson"
WORK=[
("codeaf", "exec_20261020_110826_spu6kmke", "full30-codeaf-r1", r"D:\Temp\swe-full-codeaf-20261020"),
("af-skill-intelligence", "exec_20261020_110826_4bf89ais", "full30-skill-r1", r"D:\Temp\swe-full-skill-intelligence-20261020"),
("pr-af", "exec_20261020_110826_db5taohy", "full30-pr-r1", r"D:\Temp\swe-full-pr-af-20261020"),
]
MARKER=OPS/"logs"/"full30-resume-attempts.json"
TERMINAL={"failed","timeout","cancelled","succeeded"}


def api(path, body=None):
    data=json.dumps(body).encode() if body is not None else None
    req=Request(BASE+path,data=data,
                headers={"Content-Type":"application/json"} if data else {},
                method="POST" if data else "GET")
    with urlopen(req,timeout=15) as r:return json.load(r)


def log(kind,**values):
    with LOG.open("a",encoding="utf-8") as f:
        f.write(json.dumps({"t":time.strftime("%Y-%m-%dT%H:%M:%S%z"),
                            "kind":kind,**values})+"\n")


def git(repo,*argv):
    return subprocess.check_output(
        ["git","-C",repo,*argv], text=True, timeout=10,
        creationflags=subprocess.CREATE_NO_WINDOW if os.name=="nt" else 0,
    ).strip()


def main():
    import msvcrt
    lockfile=OPS/"logs"/".full30-recovery.lock"
    with lockfile.open("a+b") as lock:
        try:msvcrt.locking(lock.fileno(),msvcrt.LK_NBLCK,1)
        except OSError:return
        config=json.loads((OPS/"project-profiles.json").read_text(encoding="utf-8"))
        runs=api("/api/v1/executions/active?limit=100").get("runs",[])
        active_ids={r.get("root_execution_id") for r in runs}
        state=json.loads(MARKER.read_text(encoding="utf-8")) if MARKER.is_file() else {}
        changed=False
        for project,root,build_id,repo in WORK:
            try:
                detail=api("/api/ui/v1/executions/"+root+"/details")
                status=detail.get("status","")
                if root in active_ids or status not in TERMINAL:
                    continue
                if status in ("succeeded","cancelled"):
                    continue
                error=str(detail.get("error_message") or "")
                # Not a reliable candidate for retry if the executor itself
                # returned a structured failure after handling coding logic.
                if status=="failed" and not any(
                    x in error.lower() for x in
                    ("timed out", "timeout", "cancelled by upstream client", "connection lost", "crashed")
                ):
                    log("manual_review_required",project=project,status=status,error=error[:240])
                    continue
                if state.get(root):
                    continue
                planfile=OPS/"plans"/("native30-"+project+".json")
                checkpoint=Path(repo)/".artifacts"/"native30-execute"/"execution"/"checkpoint.json"
                if not planfile.is_file() or not checkpoint.is_file():
                    log("no_recovery_artifacts",project=project)
                    continue
                plan=json.loads(planfile.read_text(encoding="utf-8"))
                ck=json.loads(checkpoint.read_text(encoding="utf-8"))
                if ck.get("build_id") != build_id or ck.get("repo_path") != repo:
                    log("checkpoint_mismatch",project=project)
                    continue
                if git(repo,"status","--porcelain","--untracked-files=no"):
                    log("tracked_dirty_manual_review",project=project)
                    continue
                # Query the complete original workflow, not just active roots.
                run_id=detail.get("workflow_id")
                if not run_id:
                    log("blocked_missing_workflow",project=project)
                    continue
                from urllib.parse import quote
                children=api("/api/ui/v1/agents/swe-factory-30/executions?pageSize=100&workflowId="+quote(run_id)).get("executions",[])
                # A failed process inventory query must never become "no writers".
                proc=subprocess.run(
                    ["powershell.exe","-NoProfile","-Command",
                     "Get-CimInstance Win32_Process | Where-Object { $_.Name -match 'codex|python|node' } | Select-Object -ExpandProperty CommandLine | ConvertTo-Json -Compress"],
                    capture_output=True,text=True,timeout=15,
                    creationflags=subprocess.CREATE_NO_WINDOW if os.name=="nt" else 0)
                if proc.returncode:
                    log("process_inventory_unavailable",project=project)
                    continue
                commands=json.loads(proc.stdout or "[]")
                if isinstance(commands,str): commands=[commands]
                payload,blocker=guarded_payload(
                    project,root,build_id,repo,detail,ck,plan,active_runs=runs,
                    child_executions=children,process_commands=commands)
                if blocker:
                    log("resume_blocked",project=project,root=root,reason=blocker)
                    continue
                # Write-ahead attempt record BEFORE dispatch. If a POST times
                # out ambiguously, never resend; require readback/manual QA.
                state[root]={"state":"DISPATCH_AMBIGUOUS_UNTIL_READBACK",
                             "at":time.strftime("%Y-%m-%dT%H:%M:%S%z"),
                             "build_id":build_id}
                tmp=MARKER.with_suffix(".tmp")
                tmp.write_text(json.dumps(state,indent=2),encoding="utf-8")
                os.replace(tmp,MARKER)
                changed=False
                log("resume_dispatch_start",project=project,root=root)
                response=api("/api/v1/execute/async/swe-factory-30.execute",payload)
                state[root]["response"]=response
                tmp.write_text(json.dumps(state,indent=2),encoding="utf-8")
                os.replace(tmp,MARKER)
                log("resume_dispatched",project=project,root=root,response=response)
            except Exception as exc:
                log("check_failed",project=project,error=str(exc)[:250])
        if changed:
            MARKER.write_text(json.dumps(state,indent=2),encoding="utf-8")
        try:msvcrt.locking(lock.fileno(),msvcrt.LK_UNLCK,1)
        except OSError:pass

if __name__=="__main__":
    main()
