INSTRUCCIONES PARA AUTO-HOSPEDAJE (sin nube, solo red local)

Requisitos en el PC/servidor de la agencia (Windows, Linux o macOS):
- Docker Desktop (o Docker Engine + Docker Compose)
- Al menos 2GB RAM y espacio para MySQL + build

IMPORTANTE: este docker-compose.yml vive en la carpeta del API. Para que
funcione, el FRONT debe estar en una carpeta hermana llamada "front".

1. Estructura de carpetas requerida
   quinielas/
   ├── api/                  <-- repositorio del API (aqui esta este README y docker-compose.yml)
   │   ├── docker-compose.yml
   │   ├── bd_quinielas.sql
   │   ├── Dockerfile
   │   └── (codigo Go)
   └── front/                <-- repositorio del FRONT (carpeta hermana)
       ├── Dockerfile
       ├── docker/nginx.conf
       └── (codigo Angular)

   Para obtenerlas:
     git clone https://github.com/ra0214/quinielasApi.git api
     git clone https://github.com/ra0214/quinielasFront.git front

2. Configuracion
   - El API lee DB_* desde el entorno (ya definidos en docker-compose.yml).
   - El front usa base relativa /api (nginx hace el proxy). No hay que editar URLs.

3. Levantar los servicios
   Desde la carpeta api/ (donde esta docker-compose.yml):
     docker compose up -d --build

4. Acceder
   - Aplicacion: http://IP_DEL_SERVIDOR/   (ej. http://192.168.1.100/)

5. Datos
   - MySQL persiste en el volumen docker "db_data".
   - La BD se inicializa con bd_quinielas.sql solo la primera vez (si el volumen esta vacio).

6. Apagar / actualizar
   - Parar: docker compose down
   - Actualizar codigo: git pull en api/ y front/, y volver a ejecutar
     docker compose up -d --build

Notas:
- LAN-only: basta con usar la IP interna del equipo (ipconfig / ifconfig). No hace falta abrir puertos al exterior.
- Solo se publica el puerto 80 (web) hacia el host. La base de datos y el API quedan
  dentro de la red de Docker (el front los alcanza por nginx), asi que casi no hay
  riesgo de choque con otros programas. Si el puerto 80 esta ocupado, cambia el mapeo
  del servicio "web" en docker-compose.yml (ej. "8080:80" y accede por http://IP:8080).
