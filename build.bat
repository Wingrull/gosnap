@echo off
pushd "%~dp0"

:: Создаём папку bin
echo Creating bin directory...
rmdir bin /S /Q
mkdir bin
if %errorlevel% neq 0 (
  echo Failed to create bin directory
  popd
  exit /b %errorlevel%
)

set GOOS=windows
set GO_BUILD_CMD=go build -o "bin/gosnap.exe" -a -gcflags=all="-l -B" -ldflags="-w -s" -trimpath -buildvcs=false

echo %GOOS% build ...
echo %GO_BUILD_CMD%
%GO_BUILD_CMD%
if %errorlevel% neq 0 (
  echo Build for %GOOS% failed
  popd
  exit /b %errorlevel%
)
if not exist "bin\gosnap.exe" (
  echo Error: bin/gosnap.exe not created
  popd
  exit /b 1
)

set GOOS=linux
set GO_BUILD_CMD=go build -o "bin/gosnap" -a -gcflags=all="-l -B" -ldflags="-w -s" -trimpath -buildvcs=false

echo %GOOS% build ...
echo %GO_BUILD_CMD%
%GO_BUILD_CMD%
if %errorlevel% neq 0 (
  echo Build for %GOOS% failed
  popd
  exit /b %errorlevel%
)
if not exist "bin\gosnap" (
  echo Error: bin/gosnap not created
  popd
  exit /b 1
)

popd
echo Well done.
exit /b %errorlevel%