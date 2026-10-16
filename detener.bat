@echo off
title Detener Sistema de Quinielas
cd /d "%~dp0"

echo ==========================================
echo    Deteniendo el Sistema de Quinielas...
echo ==========================================
echo.

docker --version >nul 2>&1
if errorlevel 1 (
  echo ERROR: Docker no esta instalado o no esta en ejecucion.
  echo.
  pause
  exit /b 1
)

docker compose down
if errorlevel 1 (
  echo.
  echo Ocurrio un error al detener. Revisa los mensajes de arriba.
  echo.
  pause
  exit /b 1
)

echo.
echo ==========================================
echo   El sistema se ha detenido.
echo   Tus datos (clientes, quinielas, etc.) NO se borraron.
echo ==========================================
echo.
pause
