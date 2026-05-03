import requests
import random
import time
from datetime import datetime

BACKEND_URL = "http://localhost:8080/api/telemetria"

# Lista de activos con sus rangos
activos = [
    {"id": 1, "nombre": "Vacuna Antigripal", "temp_min": 2, "temp_max": 10, "hum_min": 40, "hum_max": 75},
    {"id": 2, "nombre": "Leche Pasteurizada", "temp_min": 2, "temp_max": 8, "hum_min": 50, "hum_max": 90},
    {"id": 3, "nombre": "Componentes Electronicos", "temp_min": 20, "temp_max": 35, "hum_min": 30, "hum_max": 65},
    {"id": 4, "nombre": "Pallet de Madera", "temp_min": 15, "temp_max": 30, "hum_min": 40, "hum_max": 80},
    {"id": 5, "nombre": "Suero Intravenoso", "temp_min": 15, "temp_max": 25, "hum_min": 40, "hum_max": 70},
    {"id": 6, "nombre": "Queso Fresco", "temp_min": 4, "temp_max": 8, "hum_min": 50, "hum_max": 85},
]

ubicaciones = ["Estante A-12", "Estante B-04", "Refrigerador B-03", "Anaquel C-07", "Zona de carga"]

print("=== SIMULADOR DE SENSORES IoT ===")
print("Enviando datos a:", BACKEND_URL)
print("Presiona Ctrl+C para detener")
print("=" * 50)

while True:
    for activo in activos:
        temperatura = round(random.uniform(activo["temp_min"], activo["temp_max"]), 1)
        humedad = random.randint(activo["hum_min"], activo["hum_max"])
        ubicacion = random.choice(ubicaciones)
        
        datos = {
            "id_activo": activo["id"],
            "temperatura": temperatura,
            "humedad": humedad,
            "ubicacion": ubicacion
        }
        
        try:
            response = requests.post(BACKEND_URL, json=datos, timeout=3)
            if response.status_code == 200:
                print(f"[{datetime.now().strftime('%H:%M:%S')}] {activo['nombre']} -> Temp:{temperatura}C Hum:{humedad}% Ubic:{ubicacion} [OK]")
            else:
                print(f"[{datetime.now().strftime('%H:%M:%S')}] Error: {response.status_code}")
        except Exception as e:
            print(f"[{datetime.now().strftime('%H:%M:%S')}] Error de conexion: {e}")
        
        time.sleep(2)  # Espera 2 segundos entre sensores
    
    print("\n--- Ciclo completado. Esperando 5 segundos ---\n")
    time.sleep(5)