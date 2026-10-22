# SWE single operator entrypoint. Calls existing AgentField FULL/FAST reasoners.
# Does not create a supervisor, scheduler, or alternate execution state.
[CmdletBinding()]
param(
  [ValidateSet('full','fast')][string]$Mode='full',
  [string]$Goal='',
  [string]$RepoPath='',
  [string]$RepoUrl='',
  [string]$ExecutionId='',
  [switch]$List,
  [switch]$NoPR,
  [int]$AgentTimeoutMinutes=90,
  [int]$MaxConcurrentIssues=10,
  [int]$MaxReplans=3
)
$ErrorActionPreference='Stop'
$endpoint='http://127.0.0.1:18080'
$ops=Join-Path $env:USERPROFILE '.agentfield\ops'
$receipts=Join-Path $ops 'swe-submissions'
if ($ExecutionId) {
    $u=$endpoint+'/api/ui/v1/executions/'+[uri]::EscapeDataString($ExecutionId)+'/details'
    $resp=Invoke-RestMethod -Uri $u -TimeoutSec 15
    $logical=if ($null -ne $resp.output_data) {$resp.output_data.success} else {$null}
    $effective=$resp.status
    if ($resp.status -eq 'succeeded' -and $logical -eq $false) { $effective='LOGICAL_FAILED' }
    [pscustomobject]@{
        execution_id=$resp.execution_id; effective_status=$effective; agentfield_status=$resp.status
        workflow_id=$resp.workflow_id; node=$resp.agent_node_id; notes=$resp.notes_count
        logical_success=$logical
        summary=if ($null -ne $resp.output_data) {$resp.output_data.summary} else {$null}
    } | Format-List
    exit 0
}
if ($List) {
    $r=Invoke-RestMethod ($endpoint+'/api/v1/executions/active') -TimeoutSec 15
    $r|ConvertTo-Json -Depth 6
    exit 0
}
if (-not $Goal) { $Goal=Read-Host 'SWE task (goal)' }
if (-not $RepoPath -and -not $RepoUrl) { $RepoPath=Read-Host 'Absolute local Git repository path (or press Enter to use repo URL)';if (-not $RepoPath) {$RepoUrl=Read-Host 'Git repository URL'} }
if (-not $Goal -or (-not $RepoPath -and -not $RepoUrl)) { throw 'A goal and one repository source are required.' }
if ($RepoPath -and $RepoUrl) { throw 'Specify RepoPath OR RepoUrl, not both.' }
if ($RepoPath) {
    if (-not (Test-Path -LiteralPath $RepoPath -PathType Container)) { throw "RepoPath does not exist: $RepoPath" }
    $RepoPath=(Resolve-Path -LiteralPath $RepoPath).Path
    & git -C $RepoPath rev-parse --is-inside-work-tree *> $null
    if ($LASTEXITCODE -ne 0) { throw "Not a Git repository: $RepoPath" }
}
$agentTimeout=[Math]::Max(30,[Math]::Min($AgentTimeoutMinutes,180))*60
if ($Mode -eq 'full') {
   $target='swe-planner.build'
   $config=@{
     runtime='codex'; permission_mode='workspace-write'
     agent_timeout_seconds=$agentTimeout
     max_replans=[Math]::Max(1,[Math]::Min($MaxReplans,5))
     max_retries_per_issue=2; enable_replanning=$true
     enable_issue_advisor=$true; max_verify_fix_cycles=2
     max_concurrent_issues=[Math]::Max(1,[Math]::Min($MaxConcurrentIssues,10))
     enable_github_pr=(-not $NoPR)
   }
} else {
   $target='swe-fast.build'
   $config=@{
     runtime='codex'; permission_mode='workspace-write'
     task_timeout_seconds=$agentTimeout
     build_timeout_seconds=([Math]::Max(60,[Math]::Min($AgentTimeoutMinutes*3,360))*60)
     max_tasks=10; enable_github_pr=(-not $NoPR)
   }
}
$payload=@{input=@{goal=$Goal;config=$config}}
if ($RepoPath) {$payload.input.repo_path=$RepoPath} else {$payload.input.repo_url=$RepoUrl}
# Fail before dispatch if the control plane cannot be reached.
$health=Invoke-RestMethod ($endpoint+'/health') -TimeoutSec 10
if ($health.status -ne 'healthy') { throw 'AgentField control plane is not healthy.' }
$body=$payload|ConvertTo-Json -Depth 10 -Compress
$url=$endpoint+'/api/v1/execute/async/'+$target
$resp=Invoke-RestMethod -Method Post -Uri $url -ContentType 'application/json; charset=utf-8' -Body ([Text.Encoding]::UTF8.GetBytes($body)) -TimeoutSec 20
# Record the returned actual identity. Do not retry an ambiguous POST automatically.
if ($resp.execution_id) {
    New-Item -ItemType Directory -Path $receipts -Force|Out-Null
    $record=[pscustomobject]@{
        saved_at=(Get-Date).ToString('o')
        target=$target; mode=$Mode; goal=$Goal; repo_path=$RepoPath; repo_url=$RepoUrl
        execution_id=$resp.execution_id; workflow_id=$resp.workflow_id
        run_id=$resp.run_id; status=$resp.status
    }
    $record|ConvertTo-Json -Depth 5|Set-Content -Encoding UTF8 (Join-Path $receipts ($resp.execution_id+'.json'))
}
Write-Output "SWE ${Mode}: $target"
$resp|ConvertTo-Json -Depth 8
Write-Output "Readback: powershell -NoProfile -File `"$PSCommandPath`" -ExecutionId $($resp.execution_id)"
