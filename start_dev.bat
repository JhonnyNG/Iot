@echo off
cd /d "%~dp0"
echo ==========================================
echo    Iniciando Sistema Gestor IoT
echo ==========================================

echo [1/2] Levantando el Backend (Go) en el puerto 8080...
start "Backend Gestor IoT" cmd /k "cd /d "%~dp0"gestor-iot-backend && go run ."

echo [2/2] Levantando el Frontend (Vue Vite) en el puerto 5173...
start "Frontend Gestor IoT" cmd /k "cd /d "%~dp0"gestor-iot-frontend && npm run dev"

echo.
echo Los servicios se estan iniciando en dos ventanas nuevas.
echo Puedes acceder a tu aplicacion en: http://localhost:5173/
echo.
pause
