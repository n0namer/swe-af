# Read-only multi-project SWE-AF status. Does not claim DONE from transport status.
param([switch]$Json)
$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$cfg = Get-Content -Raw (Join-Path $root 'project-profiles.json') | ConvertFrom-Json
$base = $cfg.agentfield_url
$health = Invoke-RestMethod ($base + '/health') -TimeoutSec 8
$nodes = Invoke-RestMethod ($base + '/api/ui/v1/nodes/summary') -TimeoutSec 8
$active = Invoke-RestMethod ($base + '/api/v1/executions/active?limit=100') -TimeoutSec 10
$items = @()
foreach ($run in @($active.runs)) {
  $project = 'unmapped'
  $issue = '?'
  try {
    $details = Invoke-RestMethod ($base + '/api/ui/v1/executions/' + $run.root_execution_id + '/details') -TimeoutSec 8
    $args = $details.input_data.input
    if ($args) {
      $issue = [string]$args.issue.name
      $repo = [string]$args.repo_path
      foreach ($p in @($cfg.projects)) {
        if ([string]::Equals($repo, [string]$p.safe_repo_path,
            [StringComparison]::OrdinalIgnoreCase)) { $project = $p.id; break }
        if ($p.id -eq 'codeaf' -and $repo -like '*codeaf-swe-full-mission*') {
          $project = 'codeaf'; break
        }
      }
    }
  } catch { $issue = 'details unreadable'; }
  $items += [pscustomobject]@{
    project=$project; issue=$issue; execution=$run.root_execution_id
    target=$run.target; status=$run.root_status
    active_children=$run.active_executions
    latest_activity=$run.latest_activity
  }
}
$snapshot = [pscustomobject]@{
  checked_at=(Get-Date).ToUniversalTime().ToString('o')
  agentfield=$health.status
  swe_nodes=@($nodes.nodes | Where-Object { $_.id -like 'swe-*' } |
    Select-Object id,health_status,last_heartbeat)
  active=$items
  active_count=$active.count
  max_parallel_issue_roots=[int]$cfg.initial_limits.max_parallel_issue_roots
  available_root_slots=[Math]::Max(0,([int]$cfg.initial_limits.max_parallel_issue_roots - [int]$active.count))
  reminder='Running or succeeded transport status is not CrewSync SERVER_DONE.'
}
if ($Json) { $snapshot | ConvertTo-Json -Depth 9 }
else {
  $snapshot | Select-Object checked_at,agentfield,active_count,max_parallel_issue_roots,available_root_slots |Format-List
  $items|Format-Table project,issue,status,active_children,execution -AutoSize
  $snapshot.swe_nodes|Format-Table id,health_status,last_heartbeat -AutoSize
  Write-Output $snapshot.reminder
}
