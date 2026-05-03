# DEFENSA

Este proyecto incluye:

- `gestor-iot-backend/`: backend en Go
- `gestor-iot-frontend/`: frontend con Vite y Vue
- `gestor-iot-backend/simulador.py`: simulador de telemetría

## Configuración de la base de datos

El backend usa variables de entorno para la conexión a PostgreSQL.

1. Copia `gestor-iot-backend/.env.example` a `gestor-iot-backend/.env`
2. Ajusta los valores de `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` y `DB_SSLMODE`
3. No subas el archivo `gestor-iot-backend/.env` al repositorio: está ignorado por `.gitignore`

## Cómo ejecutar

### Backend

```bash
cd gestor-iot-backend
go run .
```

### Frontend

```bash
cd gestor-iot-frontend
npm install
npm run dev
```

## Subir al repositorio remoto

Después de inicializar git, agrega un remoto y sube el repositorio:

```bash
git remote add origin <URL_DE_TU_REPOSITORIO>
git push -u origin main
```
