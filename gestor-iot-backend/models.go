package main

type TelemetriaRequest struct {
	IDActivo    int     `json:"id_activo"`
	Temperatura float64 `json:"temperatura"`
	Humedad     int     `json:"humedad"`
	Ubicacion   string  `json:"ubicacion"`
}

type ActivoConTelemetria struct {
	IDActivo    int    `json:"id_activo"`
	Nombre      string `json:"nombre"`
	StockActual int    `json:"stock_actual"`
	UltimaTemp  string `json:"ultima_temperatura"`
	UltimaHum   string `json:"ultima_humedad"`
	UltimaUbic  string `json:"ultima_ubicacion"`
	UltimaFecha string `json:"ultima_fecha"`
}

type Alerta struct {
	IDAlerta int    `json:"id_alerta"`
	IDActivo int    `json:"id_activo"`
	Tipo     string `json:"tipo"`
	Mensaje  string `json:"mensaje"`
	Fecha    string `json:"fecha"`
	Leida    bool   `json:"leida"`
}
