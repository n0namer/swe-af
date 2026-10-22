# Native SWE-AF node for multi-project implement_issue. No custom scheduler.
$ErrorActionPreference = 'Stop'
$ops = Split-Path -Parent $MyInvocation.MyCommand.Path
$port = 8009
$busy = Get-NetTCPConnection -State Listen -LocalPort $port -ErrorAction SilentlyContinue
if ($busy) { Add-Content -Path "$ops\logs\factory-node.log" -Value "Port $port busy; refused duplicate node"; exit 0 }
$src = 'D:\Users\NIKITA\Documents\DEV\swe-af-native-runtime-30-20261020'
$python = 'D:\Users\NIKITA\Documents\DEV\universal-solver-site\agentfield\sdk\python\.venv\Scripts\python.exe'
if (-not (Test-Path $python) -or -not (Test-Path (Join-Path $src 'swe_af\app.py'))) {
  throw 'Expected native SWE-AF source or Python environment missing'
}
$env:AGENTFIELD_SERVER = 'http://127.0.0.1:18080'
$env:NO_PROXY = '127.0.0.1,localhost,::1'
$env:NODE_ID = 'swe-factory'
$env:PORT = '8009'
$env:AGENT_CALLBACK_URL = 'http://127.0.0.1:8009'
$env:SWE_DEFAULT_RUNTIME = 'codex'
$env:SWE_CODEX_AUTH_MODE = 'chatgpt'
$env:SWE_DEFAULT_MODEL = 'gpt-6-luna'
$env:SWE_MODEL_LOW = 'gpt-6-luna'
$env:SWE_MODEL_MED = 'gpt-6-luna'
$env:SWE_MODEL_HIGH = 'gpt-6.1-sol'
$env:TEMP = 'D:\Temp\codeaf-b1-go-temp-20261019'
$env:TMP = $env:TEMP
$env:GOTMPDIR = $env:TEMP
$env:GOTOOLCHAIN = 'local'
$env:GOPROXY = 'off'
$env:GOSUMDB = 'off'
$env:CODEAF_TELEMETRY = 'off'
$env:DO_NOT_TRACK = '1'
New-Item -ItemType Directory -Path $env:TEMP -Force | Out-Null
Set-Location $src
$ErrorActionPreference = 'Continue' # Native Uvicorn INFO lines arrive on stderr.
& $python -m swe_af *>> (Join-Path $ops 'logs\factory-node.log')
exit $LASTEXITCODE
