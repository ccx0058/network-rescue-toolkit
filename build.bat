@echo off
echo 正在构建网络急救工具箱...

:: 检查 Go 是否安装
where go >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo 错误: 未找到 Go，请先安装 Go 1.21+
    exit /b 1
)

:: 检查 Wails 是否安装
where wails >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo 正在安装 Wails...
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
)

:: 下载依赖
echo 正在下载 Go 依赖...
go mod tidy

:: 安装前端依赖
echo 正在安装前端依赖...
cd frontend
call npm install
cd ..

:: 构建应用
echo 正在构建应用...
wails build

echo 构建完成！输出文件在 build/bin 目录
pause
