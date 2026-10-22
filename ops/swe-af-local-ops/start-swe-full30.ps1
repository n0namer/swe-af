# Native SWE-AF Full 30-slot node. Separate identity; no new scheduler.
$ErrorActionPreference = 'Stop'
$ops = Split-Path -Parent $MyInvocation.MyCommand.Path
$port = 8010
if (Get-NetTCPConnection -State Listen -LocalPort $port -ErrorAction SilentlyContinue) { exit 0 }
$src = 'D:\Users\NIKITA\Documents\DEV\swe-af-native-runtime-heartbeat-20261020'
$python = 'D:\Users\NIKITA\Documents\DEV\universal-solver-site\agentfield\sdk\python\.venv\Scripts\python.exe'
if (-not (Test-Path $python) -or -not (Test-Path (Join-Path $src 'swe_af\app.py'))) { throw 'SWE Full30 source or python missing' }
if ((git -C $src rev-parse HEAD).Trim() -ne 'beb9049e6d5003d193ab34cfee33ddf8a3024132') { throw 'Unverified SWE Full30 SHA' }
$env:AGENTFIELD_SERVER='http://127.0.0.1:18080'
$env:NO_PROXY='127.0.0.1,localhost,::1'
$env:NODE_ID='swe-factory-30'
$env:PORT='8010'
$env:AGENT_CALLBACK_URL='http://127.0.0.1:8010'
# Full DAG parents can legitimately run beyond 2h while children QA/repair.
# The SDK default (7200s) cancelled two healthy roots on 2026-10-20.
# Use its supported 6h ceiling; per-issue budgets stay independently bounded.
$env:AGENTFIELD_ASYNC_DEFAULT_EXECUTION_TIMEOUT='21600'
$env:SWE_DEFAULT_RUNTIME='codex'
$env:SWE_CODEX_AUTH_MODE='chatgpt'
$env:SWE_DEFAULT_MODEL='gpt-6-luna'
$env:SWE_MODEL_LOW='gpt-6-luna'
$env:SWE_MODEL_MED='gpt-6-luna'
$env:SWE_MODEL_HIGH='gpt-6.1-sol'
$env:TEMP='D:\Temp\codeaf-b1-go-temp-20261019'
$env:TMP=$env:TEMP
$env:GOTMPDIR=$env:TEMP
$env:GOTOOLCHAIN='local'
$env:GOPROXY='off'
$env:GOSUMDB='off'
$env:CODEAF_TELEMETRY='off'
$env:DO_NOT_TRACK='1'
New-Item -ItemType Directory -Force -Path $env:TEMP | Out-Null
Set-Location $src
$ErrorActionPreference='Continue'
& $python -m swe_af *>> (Join-Path $ops 'logs\full30-node.log')
exit $LASTEXITCODE
