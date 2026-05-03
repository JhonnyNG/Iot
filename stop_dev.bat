@echo off
echo ==========================================
echo    Apagando Sistema Gestor IoT
echo ==========================================

echo [1/2] Cerrando Backend (Go)...
taskkill /F /IM go.exe /T >nul 2>&1
taskkill /F /IM gestor-iot.exe /T >nul 2>&1
taskkill /F /IM main.exe /T >nul 2>&1

echo [2/2] Cerrando Frontend (Node/Vite)...
taskkill /F /IM node.exe /T >nul 2>&1

echo.
echo Todos los servidores han sido apagados y desconectados.
echo Los puertos 8080 y 5173 ahora estan libres.
echo.
pause
