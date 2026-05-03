package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// POST /api/telemetria - Recibir datos del sensor
func RegistrarTelemetria(c *gin.Context) {
	var req TelemetriaRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	// Obtener umbrales del activo
	var umbralTemp, umbralHum sql.NullFloat64
	err := db.QueryRow("SELECT umbral_temperatura, umbral_humedad FROM activo WHERE id_activo = $1",
		req.IDActivo).Scan(&umbralTemp, &umbralHum)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	// Construir JSONB para PostgreSQL
	datosJSON := map[string]interface{}{
		"temperatura": req.Temperatura,
		"humedad":     req.Humedad,
		"ubicacion":   req.Ubicacion,
	}
	datosBytes, _ := json.Marshal(datosJSON)

	// Insertar telemetria
	_, err = db.Exec(`INSERT INTO telemetria (id_activo, datos, timestamp) VALUES ($1, $2, NOW())`,
		req.IDActivo, datosBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar telemetria"})
		return
	}

	// Generar alerta si supera umbrales
	if umbralTemp.Valid && req.Temperatura > umbralTemp.Float64 {
		db.Exec(`INSERT INTO alerta (id_activo, tipo, mensaje) 
                 VALUES ($1, 'temperatura', $2)`,
			req.IDActivo, fmt.Sprintf("Temperatura alta: %.1f°C (umbral %.1f°C)", req.Temperatura, umbralTemp.Float64))
	}

	if umbralHum.Valid && float64(req.Humedad) > umbralHum.Float64 {
		db.Exec(`INSERT INTO alerta (id_activo, tipo, mensaje) 
                 VALUES ($1, 'humedad', $2)`,
			req.IDActivo, fmt.Sprintf("Humedad alta: %d%% (umbral %.0f%%)", req.Humedad, umbralHum.Float64))
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Telemetria registrada"})
}

// GET /api/activos/todos-reportes - Todos los activos con su ultima telemetria
func ObtenerTodosLosReportes(c *gin.Context) {
	query := `
        SELECT 
            a.id_activo,
            a.nombre,
            a.stock_actual,
            COALESCE(t.datos->>'temperatura', '---'),
            COALESCE(t.datos->>'humedad', '---'),
            COALESCE(t.datos->>'ubicacion', '---'),
            TO_CHAR(t.timestamp, 'YYYY-MM-DD HH24:MI:SS')
        FROM activo a
        LEFT JOIN LATERAL (
            SELECT datos, timestamp
            FROM telemetria
            WHERE id_activo = a.id_activo
            ORDER BY timestamp DESC
            LIMIT 1
        ) t ON true
        ORDER BY a.id_activo
    `

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error consultando reportes"})
		return
	}
	defer rows.Close()

	var resultados []ActivoConTelemetria
	for rows.Next() {
		var r ActivoConTelemetria
		var temp, hum, ubic, fecha sql.NullString
		rows.Scan(&r.IDActivo, &r.Nombre, &r.StockActual, &temp, &hum, &ubic, &fecha)

		if temp.Valid {
			r.UltimaTemp = temp.String
		}
		if hum.Valid {
			r.UltimaHum = hum.String
		}
		if ubic.Valid {
			r.UltimaUbic = ubic.String
		}
		if fecha.Valid {
			r.UltimaFecha = fecha.String
		}

		resultados = append(resultados, r)
	}

	if resultados == nil {
		resultados = make([]ActivoConTelemetria, 0)
	}

	c.JSON(http.StatusOK, resultados)
}

// GET /api/activo/:id/ultimo-reporte - Un activo especifico
func ObtenerUltimoReporte(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	query := `
        SELECT 
            a.id_activo,
            a.nombre,
            a.stock_actual,
            COALESCE(t.datos->>'temperatura', '---'),
            COALESCE(t.datos->>'humedad', '---'),
            COALESCE(t.datos->>'ubicacion', '---'),
            TO_CHAR(t.timestamp, 'YYYY-MM-DD HH24:MI:SS')
        FROM activo a
        LEFT JOIN LATERAL (
            SELECT datos, timestamp
            FROM telemetria
            WHERE id_activo = a.id_activo
            ORDER BY timestamp DESC
            LIMIT 1
        ) t ON true
        WHERE a.id_activo = $1
    `

	var r ActivoConTelemetria
	var temp, hum, ubic, fecha sql.NullString

	err = db.QueryRow(query, idInt).Scan(&r.IDActivo, &r.Nombre, &r.StockActual, &temp, &hum, &ubic, &fecha)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	if temp.Valid {
		r.UltimaTemp = temp.String
	}
	if hum.Valid {
		r.UltimaHum = hum.String
	}
	if ubic.Valid {
		r.UltimaUbic = ubic.String
	}
	if fecha.Valid {
		r.UltimaFecha = fecha.String
	}

	c.JSON(http.StatusOK, r)
}

// GET /api/alertas - Todas las alertas
func ObtenerAlertas(c *gin.Context) {
	rows, err := db.Query(`
        SELECT id_alerta, id_activo, tipo, mensaje, TO_CHAR(fecha, 'YYYY-MM-DD HH24:MI:SS'), leida
        FROM alerta
        ORDER BY fecha DESC
    `)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error consultando alertas"})
		return
	}
	defer rows.Close()

	var alertas []Alerta
	for rows.Next() {
		var a Alerta
		rows.Scan(&a.IDAlerta, &a.IDActivo, &a.Tipo, &a.Mensaje, &a.Fecha, &a.Leida)
		alertas = append(alertas, a)
	}

	if alertas == nil {
		alertas = make([]Alerta, 0)
	}

	c.JSON(http.StatusOK, alertas)
}

// POST /api/movimiento - Registrar entrada/salida de stock
func RegistrarMovimiento(c *gin.Context) {
	var req struct {
		IDActivo  int    `json:"id_activo"`
		Tipo      string `json:"tipo"`
		Cantidad  int    `json:"cantidad"`
		IDUsuario int    `json:"id_usuario"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	if req.Tipo != "entrada" && req.Tipo != "salida" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo debe ser 'entrada' o 'salida'"})
		return
	}
	if req.Cantidad <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La cantidad debe ser mayor a 0"})
		return
	}

	// Iniciar transacción
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al iniciar transacción"})
		return
	}

	// Asegurarnos de hacer rollback si algo falla (si hay commit, el rollback es no-op)
	defer tx.Rollback()

	// 1. Validar que la salida no deje stock negativo usando FOR UPDATE
	var stockActual int
	err = tx.QueryRow(`SELECT stock_actual FROM activo WHERE id_activo = $1 FOR UPDATE`, req.IDActivo).Scan(&stockActual)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	if req.Tipo == "salida" && stockActual < req.Cantidad {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuficiente para realizar la salida"})
		return
	}

	// 2. Registrar el movimiento
	_, err = tx.Exec(`INSERT INTO movimiento (tipo, cantidad, id_usuario, id_activo) 
                       VALUES ($1, $2, $3, $4)`,
		req.Tipo, req.Cantidad, req.IDUsuario, req.IDActivo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar movimiento"})
		return
	}

	// 3. Actualizar el stock
	var operacion string
	if req.Tipo == "entrada" {
		operacion = "+"
	} else {
		operacion = "-"
	}

	_, err = tx.Exec(`UPDATE activo SET stock_actual = stock_actual `+operacion+` $1 WHERE id_activo = $2`,
		req.Cantidad, req.IDActivo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar stock"})
		return
	}

	// Hacer commit de la transacción
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al confirmar transacción"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Movimiento registrado exitosamente"})
}

// GET /api/activos - Listar todos los activos
func ListarActivos(c *gin.Context) {
	rows, err := db.Query(`SELECT id_activo, nombre, stock_actual, stock_minimo, umbral_temperatura, umbral_humedad FROM activo ORDER BY id_activo`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error consultando activos"})
		return
	}
	defer rows.Close()

	var activos []struct {
		IDActivo    int      `json:"id_activo"`
		Nombre      string   `json:"nombre"`
		StockActual int      `json:"stock_actual"`
		StockMinimo int      `json:"stock_minimo"`
		UmbralTemp  *float64 `json:"umbral_temperatura"`
		UmbralHum   *float64 `json:"umbral_humedad"`
	}

	for rows.Next() {
		var a struct {
			IDActivo    int      `json:"id_activo"`
			Nombre      string   `json:"nombre"`
			StockActual int      `json:"stock_actual"`
			StockMinimo int      `json:"stock_minimo"`
			UmbralTemp  *float64 `json:"umbral_temperatura"`
			UmbralHum   *float64 `json:"umbral_humedad"`
		}
		rows.Scan(&a.IDActivo, &a.Nombre, &a.StockActual, &a.StockMinimo, &a.UmbralTemp, &a.UmbralHum)
		activos = append(activos, a)
	}

	if activos == nil {
		activos = make([]struct {
			IDActivo    int      `json:"id_activo"`
			Nombre      string   `json:"nombre"`
			StockActual int      `json:"stock_actual"`
			StockMinimo int      `json:"stock_minimo"`
			UmbralTemp  *float64 `json:"umbral_temperatura"`
			UmbralHum   *float64 `json:"umbral_humedad"`
		}, 0)
	}

	c.JSON(http.StatusOK, activos)
}

// GET /api/telemetria/historial/:id
func ObtenerHistorialTelemetria(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	query := `
        SELECT 
            COALESCE(datos->>'temperatura', '0'),
            TO_CHAR(timestamp, 'HH24:MI:SS')
        FROM telemetria
        WHERE id_activo = $1
        ORDER BY timestamp DESC
        LIMIT 15
    `
	rows, err := db.Query(query, idInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error consultando historial"})
		return
	}
	defer rows.Close()

	var resultados []map[string]interface{}
	for rows.Next() {
		var temp, hora string
		rows.Scan(&temp, &hora)
		resultados = append(resultados, map[string]interface{}{
			"temperatura": temp,
			"hora":        hora,
		})
	}

	// Invertir para que quede orden cronológico para el gráfico (izq a der)
	for i, j := 0, len(resultados)-1; i < j; i, j = i+1, j-1 {
		resultados[i], resultados[j] = resultados[j], resultados[i]
	}

	if resultados == nil {
		resultados = make([]map[string]interface{}, 0)
	}

	c.JSON(http.StatusOK, resultados)
}
