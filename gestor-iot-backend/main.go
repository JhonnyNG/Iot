package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializar conexión a BD usando godotenv
	initDB()

	r := gin.Default()

	// Middleware CORS para que Vue.js pueda llamar a la API
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Rutas API
	api := r.Group("/api")
	{
		api.POST("/telemetria", RegistrarTelemetria)
		api.GET("/activos/todos-reportes", ObtenerTodosLosReportes)
		api.GET("/activo/:id/ultimo-reporte", ObtenerUltimoReporte)
		api.GET("/telemetria/historial/:id", ObtenerHistorialTelemetria)
		api.GET("/alertas", ObtenerAlertas)
		api.POST("/movimiento", RegistrarMovimiento)
		api.GET("/activos", ListarActivos)
		api.POST("/activos", CrearActivo)
		api.DELETE("/activos/:id", EliminarActivo)
		// Sensores custom NoSQL (JSONB flexible)
		api.POST("/sensores/custom", RegistrarTelemetriaCustom)
		api.GET("/sensores/custom/:id", ObtenerTelemetriaCustom)
	}

	fmt.Println(" Servidor corriendo en http://localhost:8080")
	r.Run(":8080")
}
