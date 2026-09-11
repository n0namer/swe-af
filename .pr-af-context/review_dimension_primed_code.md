### swe_af/execution/schemas.py (showing first 400 of 1376)
```
1: """Pydantic schemas for DAG execution state and replanning."""
2: 
3: from __future__ import annotations
4: 
5: import logging
6: import os
7: import re
8: import tempfile
9: from enum import Enum
10: from typing import Any, Literal
11: 
12: from pydantic import (
13:     BaseModel,
14:     ConfigDict,
15:     Field,
16:     PrivateAttr,
17:     field_validator,
18:     model_validator,
19: )
20: from swe_af.hitl.ask_user import AskUserForm
21: from swe_af.runtime.providers import (
22:     RUNTIME_VALUES,
23:     normalize_runtime_provider,
24:     runtime_to_harness_provider,
25: )
26: 
27: # Global default for all agent max_turns. Change this one value to adjust everywhere.
28: DEFAULT_AGENT_MAX_TURNS: int = 150
29: 
30: 
31: def ensure_str_list(value: Any) -> Any:
32:     """Coerce LLM-shaped scalars into ``list[str]`` (str → [str], None → []).
33: 
34:     Weaker models sometimes emit a single criterion/filename as a bare string
35:     where the schema wants a list. Anything else passes through unchanged so
36:     genuine type errors still surface via normal validation.
37:     """
38:     if value is None:
39:         return []
40:     if isinstance(value, str):
41:         return [value] if value.strip() else []
42:     return value
43: 
44: 
45: # ---------------------------------------------------------------------------
46: # Provider normalization
47: # ---------------------------------------------------------------------------
48: 
49: 
50: def _normalize_provider(ai_provider: str) -> str:
51:     """Map legacy provider names to AgentField native names.
52: 
53:     Ensures backward compatibility between old "claude" provider name
54:     and AgentField's native "claude-code" provider name.
55:     """
56:     return {"claude": "claude-code"}.get(ai_provider, ai_provider)
57: 
58: 
59: # ---------------------------------------------------------------------------
60: # Multi-repo helper
61: # ---------------------------------------------------------------------------
62: 
63: 
64: def _derive_repo_name(url: str) -> str:
65:     """Extract repo name from a git URL.
66: 
67:     Examples:
68:         'https://github.com/org/my-project.git' -> 'my-project'
69:         'git@github.com:org/repo.git'           -> 'repo'
70:         'https://github.com/org/repo'           -> 'repo'
71:     """
72:     if not url:
73:         return ""
74:     # Strip trailing .git, then take last path component
75:     stripped = re.sub(r"\.git$", "", url.rstrip("/"))
76:     # Handle both HTTPS and SSH URLs
77:     name = re.split(r"[/:]", stripped)[-1]
78:     return name
79: 
80: 
81: def _workspace_root() -> str:
82:     """Base directory into which builds clone repositories by default.
83: 
84:     Resolution order:
85:       1. ``SWE_WORKSPACE_ROOT`` env var, when set, on every platform.
86:       2. On Windows (``os.name == "nt"``):
87:          ``%LOCALAPPDATA%\\agentfield\\workspaces`` — an absolute drive-letter
88:          path. Falls back to ``<tempdir>\\agentfield\\workspaces`` when
89:          LOCALAPPDATA is unset.
90:       3. Everywhere else: exactly ``/workspaces`` (Docker parity).
91: 
92:     A hardcoded ``/workspaces`` base is *drive-relative* on Windows (no drive
93:     letter), which the node's spawn context resolves unpredictably: makedirs
94:     appears to succeed but ``git clone`` then fails with "destination path ...
95:     already exists and is not an empty directory". Rooting the default under an
96:     absolute base avoids that.
97: 
98:     Ref: https://github.com/Agent-Field/SWE-AF/issues/107
99:     """
100:     root = os.environ.get("SWE_WORKSPACE_ROOT")
101:     if root:
102:         return root
103:     if os.name == "nt":
104:         base = os.environ.get("LOCALAPPDATA") or tempfile.gettempdir()
105:         return os.path.join(base, "agentfield", "workspaces")
106:     return "/workspaces"
107: 
108: 
109: # ---------------------------------------------------------------------------
110: # Multi-repo models
111: # ---------------------------------------------------------------------------
112: 
113: 
114: class RepoSpec(BaseModel):
115:     """Specification for a single repository in a multi-repo build."""
116: 
117:     repo_url: str = ""  # GitHub/git URL (required if repo_path empty)
118:     repo_path: str = ""  # Absolute path to an existing local repo
119:     role: str  # 'primary' or 'dependency'
120:     branch: str = ""  # Branch to checkout (empty = default branch)
121:     sparse_paths: list[str] = []  # For sparse checkout; empty = full checkout
122:     mount_point: str = ""  # Workspace subdirectory override
123:     create_pr: bool = True  # Whether to create a PR for this repo
124: 
125:     @field_validator("role")
126:     @classmethod
127:     def _validate_role(cls, v: str) -> str:
128:         if v not in ("primary", "dependency"):
129:             raise ValueError(f"role must be 'primary' or 'dependency', got {v!r}")
130:         return v
131: 
132:     @field_validator("repo_url")
133:     @classmethod
134:     def _validate_repo_url(cls, v: str) -> str:
135:         if v and not (
136:             v.startswith("http://") or v.startswith("https://") or v.startswith("git@")
137:         ):
138:             raise ValueError(f"repo_url must be an HTTP(S) or SSH git URL, got {v!r}")
139:         return v
140: 
141: 
142: class WorkspaceRepo(BaseModel):
143:     """A repository that has been cloned into the workspace."""
144: 
145:     model_config = ConfigDict(
146:         frozen=False
147:     )  # Mutable: git_init_result assigned post-clone
148: 
149:     repo_name: str  # Derived name (from _derive_repo_name)
150:     repo_url: str  # Original git URL
151:     role: str  # 'primary' or 'dependency'
152:     absolute_path: str  # Path where the repo was cloned
153:     branch: str  # Actual checked-out branch
154:     sparse_paths: list[str] = []
155:     create_pr: bool = True
156:     git_init_result: dict | None = None  # Populated by _init_all_repos after cloning
157: 
158: 
159: class WorkspaceManifest(BaseModel):
160:     """Snapshot of all repositories cloned for a multi-repo build."""
161: 
162:     workspace_root: str  # Parent directory containing all repos
163:     repos: list[WorkspaceRepo]  # All cloned repos
164:     primary_repo_name: str  # Name of the primary repo
165: 
166:     @property
167:     def primary_repo(self) -> WorkspaceRepo | None:
168:         """Return the primary WorkspaceRepo, or None if not found."""
169:         for repo in self.repos:
170:             if repo.repo_name == self.primary_repo_name:
171:                 return repo
172:         return None
173: 
174: 
175: class RepoPRResult(BaseModel):
176:     """Result of creating a PR for a single repository."""
177: 
178:     repo_name: str
179:     repo_url: str
180:     success: bool
181:     pr_url: str = ""
182:     pr_number: int = 0
183:     error_message: str = ""
184: 
185: 
186: class AdvisorAction(str, Enum):
187:     """What the Issue Advisor decided to do after a coding loop failure."""
188: 
189:     RETRY_MODIFIED = "retry_modified"  # Relax ACs, retry coding loop
190:     RETRY_APPROACH = "retry_approach"  # Keep ACs, different strategy
191:     SPLIT = "split"  # Break into sub-issues
192:     ACCEPT_WITH_DEBT = "accept_with_debt"  # Close enough, record gaps
193:     ESCALATE_TO_REPLAN = "escalate_to_replan"  # Flag for outer loop
194: 
195: 
196: class IssueOutcome(str, Enum):
197:     """Outcome of executing a single issue."""
198: 
199:     COMPLETED = "completed"
200:     COMPLETED_WITH_DEBT = "completed_with_debt"  # Accepted via ACCEPT_WITH_DEBT
201:     FAILED_RETRYABLE = "failed_retryable"
202:     FAILED_UNRECOVERABLE = "failed_unrecoverable"
203:     FAILED_NEEDS_SPLIT = "failed_needs_split"  # Advisor wants to split
204:     FAILED_ESCALATED = "failed_escalated"  # Advisor escalated to replanner
205:     SKIPPED = "skipped"
206: 
207: 
208: class IssueAdaptation(BaseModel):
209:     """Records one AC/scope modification. Accumulated as technical debt."""
210: 
211:     adaptation_type: AdvisorAction
212:     original_acceptance_criteria: list[str] = []
213:     modified_acceptance_criteria: list[str] = []
214:     dropped_criteria: list[str] = []
215:     failure_diagnosis: str = ""
216:     rationale: str = ""
217:     new_approach: str = ""
218:     missing_functionality: list[str] = []
219:     downstream_impact: str = ""
220:     severity: str = "medium"
221: 
222: 
223: class SplitIssueSpec(BaseModel):
224:     """Sub-issue spec when advisor decides to SPLIT."""
225: 
226:     name: str
227:     title: str
228:     description: str
229:     acceptance_criteria: list[str]
230:     depends_on: list[str] = []
231:     provides: list[str] = []
232:     files_to_create: list[str] = []
233:     files_to_modify: list[str] = []
234:     parent_issue_name: str = ""
235: 
236:     @field_validator(
237:         "acceptance_criteria", "depends_on", "provides",
238:         "files_to_create", "files_to_modify",
239:         mode="before",
240:     )
241:     @classmethod
242:     def _coerce_str_list(cls, v: Any) -> Any:
243:         return ensure_str_list(v)
244: 
245: 
246: class IssueAdvisorDecision(BaseModel):
247:     """Structured output from the Issue Advisor agent."""
248: 
249:     action: AdvisorAction
250:     failure_diagnosis: str
251:     failure_category: str = ""  # environment|logic|dependency|approach|scope
252:     rationale: str
253:     confidence: float = 0.5
254:     # RETRY_MODIFIED
255:     modified_acceptance_criteria: list[str] = []
256:     dropped_criteria: list[str] = []
257:     modification_justification: str = ""
258:     # RETRY_APPROACH
259:     new_approach: str = ""
260:     approach_changes: list[str] = []
261:     # SPLIT
262:     sub_issues: list[SplitIssueSpec] = []
263:     split_rationale: str = ""
264:     # ACCEPT_WITH_DEBT
265:     missing_functionality: list[str] = []
266:     debt_severity: str = "medium"
267:     # ESCALATE_TO_REPLAN
268:     escalation_reason: str = ""
269:     dag_impact: str = ""
270:     suggested_restructuring: str = ""
271:     # Always
272:     downstream_impact: str = ""
273:     summary: str = ""
274:     # HITL — see swe_af/hitl/ for semantics
275:     ask_user_form: AskUserForm | None = None
276: 
277: 
278: class IssueResult(BaseModel):
279:     """Result of executing a single issue."""
280: 
281:     issue_name: str
282:     outcome: IssueOutcome
283:     result_summary: str = ""
284:     error_message: str = ""
285:     error_context: str = ""  # traceback/logs for replanner
286:     attempts: int = 1
287:     files_changed: list[str] = []
288:     branch_name: str = ""
289:     repo_name: str = ""  # Repo where this issue was coded (propagated from CoderResult)
290:     # Advisor fields
291:     advisor_invocations: int = 0
292:     adaptations: list[IssueAdaptation] = []
293:     debt_items: list[dict] = []
294:     split_request: list[SplitIssueSpec] | None = None
295:     escalation_context: str = ""
296:     final_acceptance_criteria: list[str] = []
297:     iteration_history: list[dict] = []
298: 
299:     @field_validator("final_acceptance_criteria", mode="before")
300:     @classmethod
301:     def _coerce_final_acceptance_criteria(cls, v: Any) -> Any:
302:         # LLM-generated fix issues have carried a bare-string criterion here;
303:         # without coercion a checkpoint reload (or the replanner's DAGState
304:         # re-validation) kills the whole build. See PR for the incident.
305:         return ensure_str_list(v)
306: 
307: 
308: class LevelResult(BaseModel):
309:     """Aggregated result of executing all issues in a single level."""
310: 
311:     level_index: int
312:     completed: list[IssueResult] = []
313:     failed: list[IssueResult] = []
314:     skipped: list[IssueResult] = []
315: 
316: 
317: class ReplanAction(str, Enum):
318:     """What the replanner decided to do."""
319: 
320:     CONTINUE = "continue"  # proceed unchanged
321:     MODIFY_DAG = "modify_dag"  # restructured
322:     REDUCE_SCOPE = "reduce_scope"  # dropped non-essential issues
323:     ABORT = "abort"  # cannot recover
324: 
325: 
326: class ReplanDecision(BaseModel):
327:     """Structured output from the replanner agent."""
328: 
329:     action: ReplanAction
330:     rationale: str
331:     updated_issues: list[dict] = []  # modified remaining issues
332:     removed_issue_names: list[str] = []
333:     skipped_issue_names: list[str] = []
334:     new_issues: list[dict] = []
335:     summary: str = ""
336:     # HITL — see swe_af/hitl/ for semantics
337:     ask_user_form: AskUserForm | None = None
338: 
339: 
340: class DAGState(BaseModel):
341:     """Full execution state of the DAG — passed to replanner for context."""
342: 
343:     # --- Artifact paths (so any agent can read the full context) ---
344:     repo_path: str = ""
345:     artifacts_dir: str = ""
346:     prd_path: str = ""
347:     architecture_path: str = ""
348:     issues_dir: str = ""
349: 
350:     # --- Plan context (summaries for quick reference by replanner) ---
351:     original_plan_summary: str = ""
352:     prd_summary: str = ""
353:     architecture_summary: str = ""
354: 
355:     # --- Issue tracking ---
356:     all_issues: list[dict] = []  # full PlannedIssue dicts
357:     levels: list[list[str]] = []  # parallel execution levels
358: 
359:     # --- Execution progress ---
360:     completed_issues: list[IssueResult] = []
361:     failed_issues: list[IssueResult] = []
362:     skipped_issues: list[str] = []
363:     in_flight_issues: list[str] = []  # names of issues currently executing
364:     current_level: int = 0
365: 
366:     # --- Replan tracking ---
367:     replan_count: int = 0
368:     replan_history: list[ReplanDecision] = []
369:     max_replans: int = 2
370: 
371:     # --- Git branch tracking ---
372:     git_integration_branch: str = ""
373:     git_original_branch: str = ""
374:     git_initial_commit: str = ""
375:     git_mode: str = ""  # "fresh" or "existing"
376:     pending_merge_branches: list[str] = []
377:     merged_branches: list[str] = []
378:     unmerged_branches: list[str] = []  # branches that failed to merge
379:     worktrees_dir: str = ""  # e.g. repo_path/.worktrees
380:     build_id: str = ""  # unique per build() call; namespaces git branches/worktrees
381: 
382:     # --- Merge/test history ---
383:     merge_results: list[dict] = []
384:     integration_test_results: list[dict] = []
385: 
386:     # --- Debt tracking ---
387:     accumulated_debt: list[dict] = []
388:     adaptation_history: list[dict] = []
389: 
390:     # --- Multi-repo workspace ---
391:     workspace_manifest: dict | None = (
392:         None  # Serialised WorkspaceManifest (dict for JSON compat)
393:     )
394: 
395: 
396: class GitInitResult(BaseModel):
397:     """Result of git initialization."""
398: 
399:     mode: str  # "fresh" or "existing"
400:     original_branch: str  # "" for fresh, e.g. "main" for existing
```
_import/usage context:_ IMPORTS: from __future__ import annotations, import logging, import os, import re, import tempfile, from enum import Enum, from typing import Any, Literal, from pydantic import (, from swe_af.hitl.ask_user import AskUserForm, from swe_af.runtime.providers import (
IMPORTED BY: none

### go/internal/config/resolve.go (showing first 400 of 434)
```
1: // Package config ports the SWE-AF configuration models (BuildConfig,
2: // ExecutionConfig, FastBuildConfig) together with their V2 validators and the
3: // runtime/model resolution logic, from swe_af/execution/schemas.py and
4: // swe_af/fast/schemas.py (design §4.7, §6).
5: //
6: // Environment variables are read via os.Getenv at call time (not cached at
7: // package init) so that tests using t.Setenv and deployers changing env after
8: // import see the current value — matching the Python