@echo off
setlocal DisableDelayedExpansion

for /F "tokens=*" %%a in ('findstr /n $') do (
  set "line=%%a"
  setlocal EnableDelayedExpansion
  set "line=!line:*:=!"
  echo(!line!
  endlocal
)

echo output

>&2 echo error
