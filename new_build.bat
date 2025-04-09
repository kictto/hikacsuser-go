@echo off
chcp 65001
echo Compiling HikVision ACS SDK Wrapper

REM Set Go environment variables
set GOOS=windows
set CGO_ENABLED=1

REM Create bin directory if it doesn't exist
if not exist bin mkdir bin

echo Copying DLL files and resources...

if exist lib\*.dll (
    xcopy /Y lib\*.dll bin\
    echo Root directory DLLs copied
)

if exist lib\win64\*.dll (
    xcopy /Y lib\win64\*.dll bin\
    echo Win64 directory DLLs copied
)
