@echo off
pushd "%~dp0"

set GOOS=windows
set GO_BUILD_CMD=go build -o "./bin" -a -gcflags=all="-l -B" -ldflags="-w -s" -trimpath -buildvcs=false

mkdir -p bin
echo %GOOS% build ...
echo %GO_BUILD_CMD%
%GO_BUILD_CMD%
if %errorlevel% neq 0 (
  popd
  echo exit code: %errorlevel%
  exit /b %errorlevel%
)

set GOOS=linux
echo %GOOS% build ...
echo %GO_BUILD_CMD%
%GO_BUILD_CMD%
if %errorlevel% neq 0 (
  popd
  echo exit code: %errorlevel%
  exit /b %errorlevel%
)

popd
echo well done.
exit /b %errorlevel%