@echo off
title Sistema de Quinielas
cd /d "%~dp0"

echo ==========================================
echo    Iniciando el Sistema de Quinielas...
echo ==========================================
echo.

docker --version >nul 2>&1
if errorlevel 1 (
  echo ERROR: Docker no esta instalado.
  echo.
  echo Instala "Docker Desktop" y vuelve a abrir este archivo.
  echo   https://www.docker.com/products/docker-desktop/
  goto :error
)

rem Verificar si el motor de Docker ya esta corriendo
docker info >nul 2>&1
if not errorlevel 1 goto :run

echo Docker Desktop no esta iniciado. Abriendolo...
if exist "%ProgramFiles%\Docker\Docker\Docker Desktop.exe" (
  start "" "%ProgramFiles%\Docker\Docker\Docker Desktop.exe"
) else if exist "%LOCALAPPDATA%\Docker\Docker Desktop.exe" (
  start "" "%LOCALAPPDATA%\Docker\Docker Desktop.exe"
) else (
  echo No se encontro Docker Desktop automaticamente.
  echo Abrelo manualmente desde el menu Inicio, espera a que diga "Engine running"
  echo y vuelve a ejecutar este archivo.
  goto :error
)

echo Esperando a que Docker termine de iniciar (puede tardar hasta 2 minutos)...
set /a intentos=0
:waitloop
timeout /t 3 /nobreak >nul
docker info >nul 2>&1
if not errorlevel 1 goto :run
set /a intentos+=1
if %intentos% geq 40 (
  echo.
  echo Docker tardo demasiado en iniciar.
  echo Abre Docker Desktop manualmente, espera a que diga "Engine running"
  echo y vuelve a ejecutar este archivo.
  goto :error
)
goto :waitloop

:run
echo Construyendo y levantando los servicios (puede tardar unos minutos la primera vez)...
echo.
docker compose up -d --build
if errorlevel 1 (
  echo.
  echo Ocurrio un error al iniciar. Revisa los mensajes de arriba.
  goto :error
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
echo Para apagarlo, ejecuta "detener.bat".
echo.
pause
exit /b 0

:error
echo.
pause
exit /b 1
