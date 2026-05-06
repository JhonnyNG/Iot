-- ============================================================================
-- GESTOR IOT - ESQUEMA DE BASE DE DATOS
-- Base de datos: PostgreSQL
-- ============================================================================

-- Tabla principal: Activos (equipos/productos)
CREATE TABLE IF NOT EXISTS activo (
    id_activo SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL UNIQUE,
    stock_actual INT DEFAULT 0,
    stock_minimo INT DEFAULT 0,
    umbral_temperatura FLOAT,
    umbral_humedad FLOAT,
    fecha_creacion TIMESTAMP DEFAULT NOW()
);

-- Tabla: Telemetría (datos de sensores)
-- Almacena datos flexibles en JSONB para soportar sensores futuros desconocidos
CREATE TABLE IF NOT EXISTS telemetria (
    id_telemetria SERIAL PRIMARY KEY,
    id_activo INT NOT NULL REFERENCES activo(id_activo) ON DELETE CASCADE,
    datos JSONB NOT NULL,  -- Ej: {"temperatura": 25.3, "humedad": 65, "ubicacion": "Almacén A"}
    timestamp TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_telemetria_activo ON telemetria(id_activo);
CREATE INDEX idx_telemetria_timestamp ON telemetria(timestamp DESC);

-- Tabla: Alertas (notificaciones de umbral)
CREATE TABLE IF NOT EXISTS alerta (
    id_alerta SERIAL PRIMARY KEY,
    id_activo INT NOT NULL REFERENCES activo(id_activo) ON DELETE CASCADE,
    tipo VARCHAR(50) NOT NULL,  -- 'temperatura', 'humedad', 'stock_bajo', etc.
    mensaje VARCHAR(500),
    fecha TIMESTAMP DEFAULT NOW(),
    leida BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_alerta_activo ON alerta(id_activo);
CREATE INDEX idx_alerta_fecha ON alerta(fecha DESC);

-- Tabla: Movimientos de stock (auditoría)
CREATE TABLE IF NOT EXISTS movimiento (
    id_movimiento SERIAL PRIMARY KEY,
    id_activo INT NOT NULL REFERENCES activo(id_activo) ON DELETE CASCADE,
    tipo VARCHAR(20) NOT NULL,  -- 'entrada' o 'salida'
    cantidad INT NOT NULL,
    id_usuario INT,  -- Usuario que realizó el movimiento
    timestamp TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_movimiento_activo ON movimiento(id_activo);
CREATE INDEX idx_movimiento_timestamp ON movimiento(timestamp DESC);

-- ============================================================================
-- DATOS DE EJEMPLO
-- ============================================================================

INSERT INTO activo (nombre, stock_actual, stock_minimo, umbral_temperatura, umbral_humedad) 
VALUES 
  ('Sensor DHT22', 50, 10, 30.0, 80.0),
  ('Arduino Uno', 25, 5, NULL, NULL),
  ('Módulo WiFi ESP32', 100, 20, 35.0, 85.0);

-- ============================================================================
-- NOTAS ARQUITECTURA
-- ============================================================================
/*
RELACIONES 1:N:
- activo → telemetria: Un activo tiene muchos registros de telemetría
- activo → alerta: Un activo puede generar múltiples alertas
- activo → movimiento: Un activo tiene historial de movimientos

ESTRATEGIA NOSQL (JSONB en telemetria):
- Soporta sensores futuros sin cambiar esquema
- Ejemplo de datos flexibles:
  * Sensor tradicional: {"temperatura": 25, "humedad": 60}
  * Sensor GPS: {"latitud": -34.123, "longitud": -58.456}
  * Sensor personalizado: {"presion": 1013, "velocidad": 45, "tipo": "custom"}

TRANSACCIONES:
- RegistrarMovimiento usa FOR UPDATE para evitar race conditions
- Cascada ON DELETE CASCADE para mantener integridad

ÍNDICES:
- Optimizados para queries por id_activo y timestamps frecuentes
*/
