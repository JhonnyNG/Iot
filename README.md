# Gestor IoT - Sistema de Gestión de Inventario con Telemetría

Sistema full-stack para gestión de activos/productos con monitoreo de sensores en tiempo real.

## 📋 Estructura del Proyecto

```
.
├── gestor-iot-backend/       # Backend en Go (Gin Framework)
│   ├── main.go              # Punto de entrada, rutas API
│   ├── handlers.go          # Lógica de endpoints
│   ├── models.go            # Estructuras de datos
│   ├── db.go                # Conexión PostgreSQL
│   └── simulador.py         # Generador de datos de prueba
├── gestor-iot-frontend/      # Frontend en Vue 3 + Vite
│   ├── src/
│   │   ├── components/      # Componentes (Dashboard, Inventario, Alertas, etc.)
│   │   ├── utils/           # Toast notifications
│   │   └── style.css        # Estilos glassmorphism
│   └── vite.config.js
├── SCHEMA.sql               # Esquema de base de datos
└── README.md
```

## 🏗️ Arquitectura Relacional (PostgreSQL)

### Modelo de Datos

```
┌─────────────────────────────────────────────────────────────┐
│                      ACTIVO (Principal)                     │
├─────────────────┬───────────────────┬─────────────────┬─────┤
│ id_activo (PK)  │ nombre            │ stock_actual    │ ... │
└─────────────────┴───────────────────┴─────────────────┴─────┘
         ▲                    ▲                    ▲
         │ 1:N                │ 1:N                │ 1:N
         │                    │                    │
    ┌────┴─────────┐    ┌─────┴──────────┐  ┌────┴─────────┐
    │  TELEMETRIA   │    │    ALERTA      │  │  MOVIMIENTO  │
    │               │    │                │  │              │
    │ - id_activo   │    │ - id_activo    │  │ - id_activo  │
    │ - datos       │    │ - tipo         │  │ - tipo       │
    │   (JSONB)     │    │ - mensaje      │  │ - cantidad   │
    │ - timestamp   │    │ - leida        │  │ - timestamp  │
    └───────────────┘    │ - timestamp    │  └──────────────┘
                         └────────────────┘
```

### Tablas y Relaciones

**1. `activo` - Tabla Principal (Equipos/Productos)**
```sql
CREATE TABLE activo (
    id_activo SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL UNIQUE,
    stock_actual INT DEFAULT 0,
    stock_minimo INT DEFAULT 0,
    umbral_temperatura FLOAT,        -- Para alertas automáticas
    umbral_humedad FLOAT,            -- Para alertas automáticas
    fecha_creacion TIMESTAMP DEFAULT NOW()
);
```

**2. `telemetria` - Datos de Sensores (1:N)**
- Un activo = múltiples lecturas de sensores
- Almacena JSONB para **flexibilidad futura**: soporta sensores desconocidos
- Ejemplo: `{"temperatura": 25.3, "humedad": 65, "ubicacion": "Almacén A"}`

```sql
CREATE TABLE telemetria (
    id_telemetria SERIAL PRIMARY KEY,
    id_activo INT NOT NULL REFERENCES activo(id_activo) ON DELETE CASCADE,
    datos JSONB NOT NULL,            -- ← Flexible para nuevos sensores
    timestamp TIMESTAMP DEFAULT NOW()
);
```

**3. `alerta` - Notificaciones Automáticas (1:N)**
- Se generan cuando telemetría excede umbrales
- Relacionadas por `id_activo`

```sql
CREATE TABLE alerta (
    id_alerta SERIAL PRIMARY KEY,
    id_activo INT NOT NULL REFERENCES activo(id_activo) ON DELETE CASCADE,
    tipo VARCHAR(50),                -- 'temperatura', 'humedad', 'stock_bajo'
    mensaje VARCHAR(500),
    leida BOOLEAN DEFAULT FALSE,
    fecha TIMESTAMP DEFAULT NOW()
);
```

**4. `movimiento` - Auditoría de Stock (1:N)**
- Historial de entrada/salida de productos
- Usa transacciones con locks para evitar race conditions

```sql
CREATE TABLE movimiento (
    id_movimiento SERIAL PRIMARY KEY,
    id_activo INT NOT NULL REFERENCES activo(id_activo) ON DELETE CASCADE,
    tipo VARCHAR(20),                -- 'entrada' o 'salida'
    cantidad INT NOT NULL,
    id_usuario INT,
    timestamp TIMESTAMP DEFAULT NOW()
);
```

### Características Arquitectónicas

**JSONB para Flexibilidad (NoSQL híbrido):**
- La columna `telemetria.datos` es JSONB
- Permite almacenar sensores futuros sin cambiar esquema
- Ejemplo: Hoy temperatura/humedad, mañana GPS o sensor personalizado

**Transacciones ACID:**
- `RegistrarMovimiento` usa `FOR UPDATE` para locking pesimista
- Garantiza que no haya race conditions en stock

**Cascada ON DELETE:**
- Eliminar un activo elimina automáticamente su telemetría, alertas y movimientos
- Mantiene integridad referencial

**Índices Optimizados:**
- `idx_telemetria_activo`: Consultas rápidas por sensor
- `idx_telemetria_timestamp`: Historial temporal
- `idx_alerta_fecha`: Alertas recientes

---

## 🔌 API Endpoints

### Backend (Go + Gin)

**Telemetría:**
- `POST /api/telemetria` - Registrar lectura de sensor
- `GET /api/telemetria/historial/:id` - Historial de sensores

**Activos:**
- `GET /api/activos` - Listar todos
- `POST /api/activos` - Crear nuevo **(NUEVO)**
- `DELETE /api/activos/:id` - Eliminar **(NUEVO)**
- `GET /api/activos/todos-reportes` - Reportes consolidados
- `GET /api/activo/:id/ultimo-reporte` - Último reporte

**Movimientos:**
- `POST /api/movimiento` - Registrar entrada/salida
- `GET /api/alertas` - Listar alertas

---

## 🎨 Frontend (Vue 3 + Vite)

**Componentes:**
- `Dashboard` - Gráficos y métricas
- `Inventario` - Gestión de stock **(CON: Crear/Editar/Eliminar)**
- `Alertas` - Notificaciones
- `Ajustes` - Configuración

**Diseño:** Dark mode con glassmorphism, responsive

---

## 🚀 Cómo Ejecutar

### Requisitos
- Go 1.20+
- Node.js 18+
- PostgreSQL 14+

### 1. Configurar Base de Datos

```bash
# Conectarse a PostgreSQL y ejecutar:
psql -U postgres -d gestor_iot < SCHEMA.sql

# O desde la CLI:
psql -U $DB_USER -d $DB_NAME -f SCHEMA.sql
```

### 2. Configurar Backend

```bash
cd gestor-iot-backend

# Crear archivo .env
cat > .env << EOF
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=tu_password
DB_NAME=gestor_iot
DB_SSLMODE=disable
EOF

# Ejecutar
go run main.go handlers.go models.go db.go
# Servidor: http://localhost:8080
```

### 3. Ejecutar Frontend

```bash
cd gestor-iot-frontend
npm install
npm run dev
# Acceso: http://localhost:5173
```

---

## 📝 Ejemplo: Crear Producto

**Request:**
```bash
curl -X POST http://localhost:8080/api/activos \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Sensor DHT22",
    "stock_actual": 50,
    "stock_minimo": 10,
    "umbral_temperatura": 30,
    "umbral_humedad": 80
  }'
```

**Response:**
```json
{
  "message": "Activo creado exitosamente",
  "id_activo": 5,
  "nombre": "Sensor DHT22"
}
```

---

## 🔄 Flujo de Datos

1. **Sensor envía datos** → `POST /api/telemetria`
2. **Backend valida y almacena** en tabla `telemetria` (JSONB)
3. **Backend compara con umbrales** de tabla `activo`
4. **Si supera umbral** → Inserta en tabla `alerta`
5. **Frontend consulta** → `GET /api/alertas`
6. **Dashboard muestra** notificación en tiempo real

---

## 📚 Para Responder: "Cómo Manejas Relaciones"

> **Respuesta Completa:**

"Utilizamos PostgreSQL con un modelo relacional normalizado donde `activo` es la tabla principal (PK: id_activo) con relaciones 1:N hacia:
- **telemetria**: Almacena datos de sensores en JSONB para flexibilidad futura
- **alerta**: Notificaciones cuando telemetría excede umbrales
- **movimiento**: Auditoría de stock con transacciones ACID

La capa de conexión es nativa con `lib/pq` (driver PostgreSQL en Go). El esquema SQL está en `SCHEMA.sql` y se inicializa en la BD. Usamos locks pesimistas (FOR UPDATE) para garantizar integridad en movimientos concurrentes."

---

## 🛠️ Desarrollo

```bash
# Backend en modo debug:
export GIN_MODE=debug
go run main.go handlers.go models.go db.go

# Frontend con hot-reload:
npm run dev

# Generar datos de prueba:
python3 gestor-iot-backend/simulador.py
```

---

## 📦 Deploy

```bash
# Backend:
cd gestor-iot-backend
go build -o gestor-iot-server
./gestor-iot-server

# Frontend:
cd gestor-iot-frontend
npm run build
# Servir dist/ con nginx/apache
```

---

## 📄 Licencia

MIT
