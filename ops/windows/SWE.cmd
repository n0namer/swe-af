@echo off
setlocal
title SWE - Full / Fast
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%USERPROFILE%\.agentfield\ops\SWE.ps1" %*
if errorlevel 1 (
  echo SWE command returned an error.
  pause
)
