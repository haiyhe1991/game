echo off & color 0A

set SOURCE_FOLDER=.


set GO_TARGET_PATH=.\


for /f "delims=" %%i in ('dir /b "%SOURCE_FOLDER%\*.proto"') do (
    echo protoc -I. -I=%GOPATH%\src --gogoslick_out=%GO_TARGET_PATH% %SOURCE_FOLDER%\%%i 
    protoc -I. -I=%GOPATH%\src --proto_path=%SOURCE_FOLDER% --gogoslick_out=%GO_TARGET_PATH% %SOURCE_FOLDER%\%%i 
)

echo SUCCESS.

pause



