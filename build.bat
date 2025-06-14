@echo off
chcp 65001>nul
echo "正在构建护眼提示器应用..."

REM 确保依赖已安装
echo "检查并更新依赖..."
go mod tidy

REM 生成Windows资源文件
echo "生成资源文件..."
if not exist "rsrc.syso" (
    rsrc -manifest app.manifest -o rsrc.syso
    if %ERRORLEVEL% NEQ 0 (
        echo 生成资源文件失败，请检查错误信息。
        pause
        exit /b 1
    )
)

REM 构建Windows可执行文件
echo "编译应用..."
go build -tags walk_use_cgo -ldflags="-H windowsgui -s -w" -o EyeProtectionTimer.exe

if %ERRORLEVEL% EQU 0 (
    echo "构建成功！"
    echo "可执行文件: %CD%\EyeProtectionTimer.exe"
) else (
    echo "构建失败，请检查错误信息。"
)

pause