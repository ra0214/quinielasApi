@echo off
title Sistema de Quinielas
cd /d "%~dp0"

echo ==========================================
echo    Iniciando el Sistema de Quinielas...
echo ==========================================
echo.

docker --version >nul 2>&1
if errorlevel 1 (
  echo ERROR: Docker no esta instalado o no esta en ejecucion.
  echo.
  echo Instala "Docker Desktop" y vuelve a abrir este archivo.
  echo   https://www.docker.com/products/docker-desktop/
  echo.
  pause
  exit /b 1
)

echo Construyendo y levantando los servicios (puede tardar unos minutos la primera vez)...
echo.
docker compose up -d --build
if errorlevel 1 (
  echo.
  echo Ocurrio un error al iniciar. Revisa los mensajes de arriba.
  echo.
  pause
  exit /b 1
)

echo.
echo Listo. La aplicacion se abrira en el navegador en unos segundos...
timeout /t 5 /nobreak >nul
start "" http://localhost

echo.
echo ==========================================
echo   La aplicacion esta corriendo.
echo.
echo   En este equipo:      http://localhost
echo   En otros equipos:    http://IP_DE_ESTE_PC   (ver con "ipconfig")
echo ==========================================
echo.
echo Esta ventana se puede cerrar. El sistema seguira funcionando.
echo Para apagarlo, usa el comando: docker compose down
echo.
pause
