<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { showToast } from '../utils/toast.js'

const activos = ref([])
const isLoading = ref(true)

// Datos del simulador
const simData = ref({
  id_activo: '',
  temperatura: 25.0,
  humedad: 50,
  ubicacion: 'Almacén Principal'
})

const isAutoSimulating = ref(false)
let autoSimInterval = null

const fetchActivos = async () => {
  try {
    const res = await fetch('http://localhost:8080/api/activos')
    if (res.ok) {
      activos.value = await res.json()
      if (activos.value.length > 0) {
        simData.value.id_activo = activos.value[0].id_activo
      }
    }
  } catch (error) {
    console.error("Error fetching activos:", error)
  } finally {
    isLoading.value = false
  }
}

const sendSimulatedData = async (isAuto = false) => {
  if (!simData.value.id_activo) {
    if (!isAuto) showToast("Seleccione un activo válido", "warning")
    return
  }
  
  try {
    const res = await fetch('http://localhost:8080/api/telemetria', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        id_activo: parseInt(simData.value.id_activo),
        temperatura: parseFloat(simData.value.temperatura.toFixed(1)),
        humedad: parseInt(simData.value.humedad),
        ubicacion: simData.value.ubicacion
      })
    })
    
    if (res.ok) {
      if (!isAuto) showToast("Datos de telemetría enviados simuladamente", "success")
      
      // Chequear si podría haber generado una alerta (incluso en modo auto para que el profe lo vea)
      const activo = activos.value.find(a => a.id_activo == simData.value.id_activo)
      if (activo) {
        if (activo.umbral_temperatura && simData.value.temperatura > activo.umbral_temperatura) {
          showToast(`¡Alerta generada por ALTA TEMPERATURA!`, "error")
        }
        if (activo.umbral_humedad && simData.value.humedad > activo.umbral_humedad) {
          setTimeout(() => showToast(`¡Alerta generada por ALTA HUMEDAD!`, "error"), 500)
        }
      }
    } else {
      if (!isAuto) showToast("Error al enviar la telemetría", "error")
    }
  } catch (error) {
    console.error("Error sending telemetry:", error)
    if (!isAuto) showToast("Error de conexión con el backend", "error")
    
    // Si hay error en auto, lo detenemos por seguridad
    if (isAuto) stopAutoSimulation()
  }
}

const toggleAutoSimulation = () => {
  if (isAutoSimulating.value) {
    stopAutoSimulation()
  } else {
    startAutoSimulation()
  }
}

const startAutoSimulation = () => {
  isAutoSimulating.value = true
  showToast("Simulación Automática Iniciada", "success")
  
  // Enviar un dato inmediatamente, luego cada 3 segundos
  generarFluctuacion()
  sendSimulatedData(true)
  
  autoSimInterval = setInterval(() => {
    generarFluctuacion()
    sendSimulatedData(true)
  }, 3000)
}

const stopAutoSimulation = () => {
  isAutoSimulating.value = false
  if (autoSimInterval) {
    clearInterval(autoSimInterval)
    autoSimInterval = null
  }
  showToast("Simulación Automática Detenida", "warning")
}

const generarFluctuacion = () => {
  // Fluctuar temperatura entre -1.5 y +1.5 grados (Random Walk)
  const tempDelta = (Math.random() - 0.5) * 3
  let nuevaTemp = parseFloat(simData.value.temperatura) + tempDelta
  
  // Mantenerla en un rango realista (ej. 15 a 45) para que el gráfico no se vuelva loco
  if (nuevaTemp < 15) nuevaTemp = 15
  if (nuevaTemp > 45) nuevaTemp = 45
  
  simData.value.temperatura = parseFloat(nuevaTemp.toFixed(1))
  
  // Fluctuar humedad levemente
  const humDelta = Math.floor((Math.random() - 0.5) * 4)
  let nuevaHum = parseInt(simData.value.humedad) + humDelta
  if (nuevaHum < 30) nuevaHum = 30
  if (nuevaHum > 90) nuevaHum = 90
  simData.value.humedad = nuevaHum
}

onMounted(() => {
  fetchActivos()
})

onUnmounted(() => {
  if (autoSimInterval) clearInterval(autoSimInterval)
})
</script>

<template>
  <div class="ajustes">
    <div class="header-section mb-4">
      <h2 class="section-title">Herramientas de Desarrollador</h2>
    </div>

    <div class="main-grid">
      <!-- Simulador Panel -->
      <div class="glass-panel p-4 simulator-panel">
        <div class="simulator-header mb-4">
          <div class="simulator-icon" :class="isAutoSimulating ? 'bg-success-light' : ''">
            <i class="ph text-primary" :class="isAutoSimulating ? 'ph-waveform text-success' : 'ph-broadcast'"></i>
          </div>
          <div>
            <h3 class="font-medium" style="font-size: 1.125rem;">Simulador de Sensores IoT</h3>
            <p class="text-muted" style="font-size: 0.875rem;">
              {{ isAutoSimulating ? 'Generando variaciones de temperatura en tiempo real...' : 'Envía datos virtuales para probar alertas y gráficos.' }}
            </p>
          </div>
        </div>

        <div v-if="isLoading" class="text-muted text-center py-4">Cargando sistema...</div>

        <form v-else @submit.prevent="sendSimulatedData(false)" class="simulator-form">
          <div class="form-group">
            <label>Seleccionar Activo (ID)</label>
            <select v-model="simData.id_activo" class="input-field" :disabled="isAutoSimulating">
              <option v-for="a in activos" :key="a.id_activo" :value="a.id_activo">
                #{{ a.id_activo }} - {{ a.nombre }} (Max: {{ a.umbral_temperatura || 'N/A' }}°C / {{ a.umbral_humedad || 'N/A' }}%)
              </option>
            </select>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>Temperatura simulada (°C)</label>
              <div class="input-with-icon">
                <i class="ph ph-thermometer" :class="isAutoSimulating ? 'text-danger pulse-icon' : 'text-muted'"></i>
                <input type="number" step="0.1" v-model="simData.temperatura" class="input-field" :disabled="isAutoSimulating" />
              </div>
            </div>

            <div class="form-group">
              <label>Humedad simulada (%)</label>
              <div class="input-with-icon">
                <i class="ph ph-drop" :class="isAutoSimulating ? 'text-info pulse-icon' : 'text-muted'"></i>
                <input type="number" v-model="simData.humedad" class="input-field" :disabled="isAutoSimulating" />
              </div>
            </div>
          </div>

          <div class="form-group mb-4">
            <label>Ubicación reportada</label>
            <div class="input-with-icon">
              <i class="ph ph-map-pin text-muted"></i>
              <input type="text" v-model="simData.ubicacion" class="input-field" :disabled="isAutoSimulating" />
            </div>
          </div>

          <div class="button-group">
            <button type="submit" class="btn btn-outline flex-1" style="height: 3rem;" :disabled="isAutoSimulating">
              <i class="ph ph-paper-plane-tilt"></i> Enviar Manual
            </button>
            <button type="button" class="btn flex-2" :class="isAutoSimulating ? 'btn-danger-solid' : 'btn-primary'" style="height: 3rem;" @click="toggleAutoSimulation">
              <i class="ph" :class="isAutoSimulating ? 'ph-stop-circle' : 'ph-play-circle'"></i> 
              {{ isAutoSimulating ? 'Detener Simulación Automática' : 'Iniciar Simulación Automática' }}
            </button>
          </div>
        </form>
      </div>

      <!-- General Settings Placeholder -->
      <div class="glass-panel p-4 text-center settings-placeholder">
        <i class="ph ph-wrench text-muted mb-3" style="font-size: 3rem;"></i>
        <h3 class="font-medium mb-2">Ajustes Globales</h3>
        <p class="text-muted" style="font-size: 0.875rem;">Sección de configuración de la cuenta y preferencias de tema en construcción.</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.mb-4 { margin-bottom: 1.5rem; }
.mb-3 { margin-bottom: 1rem; }
.mb-2 { margin-bottom: 0.5rem; }
.p-4 { padding: 1.5rem; }
.py-4 { padding-top: 1.5rem; padding-bottom: 1.5rem; }
.text-center { text-align: center; }
.font-medium { font-weight: 500; }
.text-muted { color: var(--text-muted); }
.text-primary { color: var(--primary); }
.text-success { color: var(--success); }
.text-danger { color: #ef4444; }
.text-info { color: #3b82f6; }
.w-full { width: 100%; }

.flex-1 { flex: 1; }
.flex-2 { flex: 2; }

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-main);
}

.main-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 1.5rem;
}

@media (max-width: 1024px) {
  .main-grid {
    grid-template-columns: 1fr;
  }
}

.simulator-header {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.simulator-icon {
  width: 48px;
  height: 48px;
  background: rgba(59, 130, 246, 0.1);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.75rem;
  transition: all 0.3s;
}

.bg-success-light {
  background: rgba(16, 185, 129, 0.1);
}

.simulator-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-muted);
}

.input-with-icon {
  position: relative;
  display: flex;
  align-items: center;
}

.input-with-icon i {
  position: absolute;
  left: 0.75rem;
  font-size: 1.25rem;
}

.pulse-icon {
  animation: beat 1s infinite;
}

@keyframes beat {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.2); }
}

.input-field {
  width: 100%;
  padding: 0.75rem;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  background: rgba(0, 0, 0, 0.2);
  color: white;
  font-family: inherit;
  transition: border-color 0.2s;
}

.input-field:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

select.input-field {
  appearance: none;
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2394a3b8' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
  background-repeat: no-repeat;
  background-position: right 0.75rem center;
  background-size: 1em;
  padding-right: 2.5rem;
}

.input-with-icon .input-field {
  padding-left: 2.5rem;
}

.input-field:focus {
  border-color: var(--primary);
  outline: none;
}

.button-group {
  display: flex;
  gap: 1rem;
}

.btn-danger-solid {
  background: var(--danger);
  color: white;
}
.btn-danger-solid:hover {
  background: #dc2626;
}

.settings-placeholder {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
</style>
