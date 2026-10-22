# Thin REST client for native SWE-AF implement_issue; not an orchestrator.
param(
  [Parameter(Mandatory=$true)][string]$Project,
  [Parameter(Mandatory=$true)][string]$IssueFile,
  [switch]$Submit
)
$ErrorActionPreference = 'Stop'
$ops = Split-Path -Parent $MyInvocation.MyCommand.Path
$cfg = Get-Content -Raw (Join-Path $ops 'project-profiles.json') | ConvertFrom-Json
$p = @($cfg.projects | Where-Object { $_.id -eq $Project })
if ($p.Count -ne 1) { throw "Unknown project: $Project" }
$p = $p[0]
if (-not (Test-Path $p.safe_repo_path)) { throw "Missing isolated repo: $($p.safe_repo_path)" }
$issue = Get-Content -Raw -LiteralPath $IssueFile | ConvertFrom-Json
if (-not $issue.name -or -not $issue.title -or -not $issue.description) {
  throw 'Issue requires name/title/description'
}
# Permit only evidence-gathering G0 before PR-AF permanent DEV identity is proven.
# Coding/redeploy remains forbidden until source/runtime/callback identity is verified.
if ($p.readiness -eq 'g0_runtime_unverified') {
  $allowed = ($p.id -eq 'pr-af') -and
    ($issue.name -eq 'pr-af-g0-permanent-dev-runtime-identity') -and
    (@($issue.files_to_modify).Count -eq 1) -and
    ($issue.files_to_modify[0] -eq 'PLAN.md') -and
    (@($issue.files_to_create).Count -eq 0)
  if (-not $allowed) {
    throw 'PR-AF runtime unverified: only exact PLAN.md G0 evidence task is allowed'
  }
}
$wanted = @(@($issue.files_to_modify) + @($issue.files_to_create) |
  Where-Object { $_ -and $_ -ne '' } | Select-Object -Unique)
if ($wanted.Count -eq 0) { throw 'Issue has no owned files; avoid unscoped coding' }
git -C $p.safe_repo_path rev-parse --verify $p.base_branch *> $null
if ($LASTEXITCODE -ne 0) { throw "Missing base branch $($p.base_branch)" }
$health = Invoke-RestMethod ($cfg.agentfield_url + '/health') -TimeoutSec 8
if ($health.status -ne 'healthy') { throw 'AgentField unhealthy' }
$active = Invoke-RestMethod ($cfg.agentfield_url + '/api/v1/executions/active?limit=100') -TimeoutSec 8
foreach ($run in @($active.runs)) {
  if ($run.target -notlike '*.implement_issue') { continue }
  $url = $cfg.agentfield_url + '/api/ui/v1/executions/' + $run.root_execution_id + '/details'
  try { $d = Invoke-RestMethod $url -TimeoutSec 8 }
  catch { throw "Cannot prove no duplicate for $($run.root_execution_id): $($_.Exception.Message)" }
  $prior = $d.input_data.input
  if (-not $prior -or -not $prior.repo_path) { continue }
  $pathA = [IO.Path]::GetFullPath([string]$prior.repo_path)
  $pathB = [IO.Path]::GetFullPath([string]$p.safe_repo_path)
  if (-not $pathA.Equals($pathB, [StringComparison]::OrdinalIgnoreCase)) { continue }
  if ($prior.issue.name -eq $issue.name) {
    throw "Same issue already running: $($run.root_execution_id)"
  }
  $occupied = @(@($prior.issue.files_to_modify) + @($prior.issue.files_to_create) |
    Where-Object { $_ -and $_ -ne '' })
  foreach ($f in $wanted) {
    if ($occupied -contains $f) {
      throw "File conflict for $f with active $($run.root_execution_id)"
    }
  }
}
if ($Submit -and [int]$active.count -ge [int]$cfg.initial_limits.max_parallel_issue_roots) {
  throw "Configured parallel cap reached ($($active.count) active issue roots / $($cfg.initial_limits.max_parallel_issue_roots) allowed); wait for a slot or adjust capacity explicitly"
}
$models = @{ default='gpt-6-luna'; coder='gpt-6-luna'
  qa='gpt-6.1-sol'; code_reviewer='gpt-6.1-sol'
  qa_synthesizer='gpt-6-luna'; verifier='gpt-6.1-sol'; git='gpt-6-luna' }
$payload = @{ input = @{
  issue=$issue; repo_path=$p.safe_repo_path; base_branch=$p.base_branch
  artifacts_dir='.artifacts'
  additional_context=("Project="+$p.id+"; source of truth="+$p.source_of_truth+
    "; original checkout is read-only; BMad Quick Dev, native AgentField; "+
    "no external sends, no deploy, no push, no self-accepted DONE")
  config=@{runtime='codex';models=$models
    max_coding_iterations=2;agent_max_turns=18;agent_timeout_seconds=1800
    verify=$true; enable_github_pr=$false;keep_worktree=$true
    branch_prefix=$p.branch_prefix }
}}
if (-not $Submit) {
  [pscustomobject]@{
    mode='DRY_RUN'; project=$p.id; issue=$issue.name
    repo=$p.safe_repo_path; base_branch=$p.base_branch
    target=$cfg.submission_target; currently_active=$active.count
    required_action='Repeat with -Submit to enqueue'
  } | Format-List
  exit 0
}
$url = $cfg.agentfield_url + '/api/v1/execute/async/' + $cfg.submission_target
$result = Invoke-RestMethod -Uri $url -Method Post -ContentType 'application/json' -TimeoutSec 20 -Body ($payload | ConvertTo-Json -Depth 12)
$result | ConvertTo-Json -Depth 6
