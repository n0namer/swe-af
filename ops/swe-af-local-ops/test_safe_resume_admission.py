"""Offline fail-closed acceptance probes for SWE-AF recovery admission."""
import sys
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch
sys.path.insert(0, str(Path(__file__).resolve().parent))
from safe_resume_admission import guarded_payload

REPO=r"D:\Temp\swe-full-fixture"
PLAN={"issues":[{"name":"alpha"}],"levels":[["alpha"]]}
CK={"build_id":"b1","repo_path":REPO,"all_issues":[{"name":"alpha"}],"levels":[["alpha"]],"failed_issues":[],"completed_issues":[],"in_flight_issues":[]}
DETAIL={"status":"failed","error_message":"Execution cancelled by upstream client","input_data":{"input":{"plan_result":PLAN,"repo_path":REPO,"build_id":"b1","git_config":{"integration_branch":"mission","initial_commit_sha":"abc"},"config":{}}}}
def go(**overrides):
    kwargs=dict(project="fixture",root="exec_1",build_id="b1",repo=REPO,detail=DETAIL,
                checkpoint=CK,plan=PLAN,active_runs=[],child_executions=[],process_commands=[])
    kwargs.update(overrides)
    return guarded_payload(**kwargs)

class Admission(unittest.TestCase):
    @patch("safe_resume_admission.git", side_effect=["","", "mission"])
    def test_valid_recoverable_is_exact_resume(self,_):
        payload,reason=go()
        self.assertIsNone(reason)
        self.assertTrue(payload["input"]["resume"])
        self.assertEqual(payload["input"]["plan_result"],PLAN)
        self.assertEqual(payload["input"]["build_id"],"b1")
    def test_noop_healthy(self):
        _,why=go(detail={**DETAIL,"status":"running"})
        self.assertEqual(why,"not_failed_transport")
    def test_noop_logical_failure(self):
        _,why=go(detail={**DETAIL,"error_message":"QA failed"})
        self.assertEqual(why,"logical_failure_not_recoverable")
    def test_active_owner(self):
        _,why=go(active_runs=[{"root_execution_id":"other"}])
        self.assertEqual(why,"active_agentfield_execution")
    def test_child_writer(self):
        _,why=go(child_executions=[{"status":"running"}])
        self.assertEqual(why,"live_or_ambiguous_child")
    def test_process_writer(self):
        _,why=go(process_commands=["codex exec -C "+REPO])
        self.assertEqual(why,"live_or_ambiguous_process")
    def test_dirty_failed_checkpoint(self):
        _,why=go(checkpoint={**CK,"failed_issues":[{"issue_name":"bad"}]})
        self.assertEqual(why,"logical_issue_failures")
    def test_stale_plan(self):
        _,why=go(plan={"issues":[{"name":"other"}],"levels":[["other"]]})
        self.assertEqual(why,"plan_checkpoint_drift")
    def test_inflight_checkpoint(self):
        _,why=go(checkpoint={**CK,"in_flight_issues":["alpha"]})
        self.assertEqual(why,"inflight_checkpoint_ambiguous")
    @patch("safe_resume_admission.git",return_value=" M dirty.txt")
    def test_dirty_git(self,_):
        _,why=go()
        self.assertEqual(why,"dirty_integration_worktree")

if __name__=="__main__": unittest.main()
