# Read-only Windows AgentField / SWE-AF monitor. No agent dispatch, retries or task mutation.
param([switch]$Once)
$ErrorActionPreference = 'Stop'
$ops = Split-Path -Parent $MyInvocation.MyCommand.Path
$logs = Join-Path $ops 'logs'
New-Item -ItemType Directory -Path $logs -Force | Out-Null
$lockPath = Join-Path $ops '.monitor.lock'
$lock = $null
try {
  $lock = [System.IO.File]::Open($lockPath, [System.IO.FileMode]::OpenOrCreate,
    [System.IO.FileAccess]::ReadWrite, [System.IO.FileShare]::None)
} catch { exit 0 } # Another invocation still running.
function Emit-Event([string]$severity, [string]$kind, [object]$details) {
  $event = [pscustomobject]@{ timestamp=(Get-Date).ToUniversalTime().ToString('o')
    level=$severity; kind=$kind; details=$details }
  $file = Join-Path $logs 'events.ndjson'
  if ((Test-Path $file) -and (Get-Item $file).Length -gt 4194304) {
    Move-Item $file (Join-Path $logs 'events.previous.ndjson') -Force
  }
  Add-Content -LiteralPath $file -Value ($event | ConvertTo-Json -Depth 8 -Compress) -Encoding UTF8
}
try {
  $server = 'http://127.0.0.1:18080'
  $health = Invoke-RestMethod -Uri ($server + '/health') -TimeoutSec 7
  $ns = Invoke-RestMethod -Uri ($server + '/api/ui/v1/nodes/summary') -TimeoutSec 7
  $runs = Invoke-RestMethod -Uri ($server + '/api/v1/executions/active?limit=100') -TimeoutSec 8
  $nodes = @($ns.nodes | Where-Object { $_.id -like 'swe-*' } |
    ForEach-Object { [pscustomobject]@{ id=$_.id; health=$_.health_status
      lifecycle=$_.lifecycle_status; heartbeat=$_.last_heartbeat } })
  $active = @($runs.runs | ForEach-Object { [pscustomobject]@{
    run_id=$_.run_id; execution_id=$_.root_execution_id
    target=$_.target; status=$_.root_status; active_executions=$_.active_executions
  } })
  $nodeFingerprint = @($nodes | Sort-Object id | ForEach-Object { $_.id+':'+$_.health+':'+$_.lifecycle }) -join '|'
  $runFingerprint = @($active | Sort-Object execution_id | ForEach-Object { $_.execution_id+':'+$_.status }) -join '|'
  $signature = 'health='+$health.status+'|nodes='+$nodeFingerprint+'|runs='+$runFingerprint
  $snapshotPath = Join-Path $logs 'latest.json'
  $old = $null
  if (Test-Path $snapshotPath) {
    try { $old = Get-Content -Raw -LiteralPath $snapshotPath | ConvertFrom-Json }
    catch { Emit-Event 'warn' 'previous_snapshot_invalid' $_.Exception.Message }
  }
  $now = (Get-Date).ToUniversalTime()
  foreach ($node in $nodes) {
    try {
      $age = ($now - ([DateTimeOffset]::Parse($node.heartbeat)).UtcDateTime).TotalSeconds
      if ($age -gt 30) { Emit-Event 'error' 'stale_node_heartbeat' @{
        node=$node.id; age_seconds=[math]::Round($age,0) } }
    } catch { Emit-Event 'warn' 'heartbeat_unparseable' $node.id }
  }
  if ($old -and $old.signature -ne $signature) {
    Emit-Event 'info' 'fleet_state_changed' @{ previous=$old.signature; current=$signature }
    $currentIDs = @($active | ForEach-Object { $_.execution_id })
    foreach ($before in @($old.active)) {
      if ($before.execution_id -and $before.execution_id -notin $currentIDs) {
        try {
          $url = $server + '/api/ui/v1/executions/' + $before.execution_id + '/details'
          $detail = Invoke-RestMethod -Uri $url -TimeoutSec 8
          Emit-Event 'info' 'execution_terminal_readback' @{
            id=$before.execution_id; status=$detail.status
            result_success=$detail.output_data.success
            outcome=$detail.output_data.outcome; error=$detail.error_message
          }
        } catch { Emit-Event 'warn' 'terminal_readback_failed' $before.execution_id }
      }
    }
  } elseif (-not $old) { Emit-Event 'info' 'fleet_monitor_started' $signature }
  $result = [pscustomobject]@{
    timestamp=$now.ToString('o'); health=$health.status; signature=$signature
    nodes=$nodes; active=$active; active_count=$runs.count
  }
  $result | ConvertTo-Json -Depth 9 | Set-Content -LiteralPath $snapshotPath -Encoding UTF8
  if ($health.status -ne 'healthy') { Emit-Event 'error' 'agentfield_unhealthy' $health.status }
} catch {
  Emit-Event 'error' 'agentfield_monitor_failed' $_.Exception.Message
  exit 1
} finally {
  if ($lock) { $lock.Dispose() }
}
