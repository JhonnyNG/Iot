<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import SensorCard from './SensorCard.vue'

// Importaciones para el gráfico (Chart.js + vue-chartjs)
import { Line } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, LineElement, PointElement, CategoryScale, LinearScale, Filler } from 'chart.js'

ChartJS.register(Title, Tooltip, Legend, LineElement, PointElement, CategoryScale, LinearScale, Filler)

const activos = ref([])
const alertas = ref([])
const historialTemp = ref([])
const isError = ref(false)
const isLoading = ref(true)

const fetchData = async () => {
  try {
    const resActivos = await fetch('http://localhost:8080/api/activos/todos-reportes')
    if (resActivos.ok) activos.value = await resActivos.json()
    
    const resAlertas = await fetch('http://localhost:8080/api/alertas')
    if (resAlertas.ok) alertas.value = (await resAlertas.json()).slice(0, 5) // Top 5

    // Obtener historial del primer activo (si existe) para el gráfico en vivo
    if (activos.value.length > 0) {
      const firstId = activos.value[0].id_activo
      const resHist = await fetch(`http://localhost:8080/api/telemetria/historial/${firstId}`)
      if (resHist.ok) {
        historialTemp.value = await resHist.json()
      }
    }
    
    isError.value = false
    isLoading.value = false
  } catch (error) {
    console.error("Error fetching data:", error)
    isError.value = true
  }
}

let interval
let visualInterval

onMounted(() => {
  fetchData()
  interval = setInterval(fetchData, 5000) // Poll every 5s

  // Efecto Visual Dinámico: Simula que los sensores son ultra-sensibles leyendo en tiempo real
  visualInterval = setInterval(() => {
    activos.value.forEach(a => {
      if (a.ultima_temperatura && a.ultima_temperatura !== '---') {
        let temp = parseFloat(a.ultima_temperatura)
        // Variación mínima de +/- 0.2 grados
        temp += (Math.random() - 0.5) * 0.4
        a.ultima_temperatura = temp.toFixed(1)
      }
      
      if (a.ultima_humedad && a.ultima_humedad !== '---') {
        let hum = parseInt(a.ultima_humedad)
        // Variación de +/- 1%
        hum += Math.random() > 0.5 ? 1 : -1
        // Limitar entre 0 y 100
        if (hum < 0) hum = 0
        if (hum > 100) hum = 100
        a.ultima_humedad = hum.toString()
      }
    })
  }, 1500) // Cambia visualmente cada segundo y medio
})

onUnmounted(() => {
  clearInterval(interval)
  clearInterval(visualInterval)
})

// Configuraciones del Gráfico de Líneas
const chartData = computed(() => {
  return {
    labels: historialTemp.value.map(h => h.hora),
    datasets: [
      {
        label: 'Temperatura (°C)',
        backgroundColor: 'rgba(239, 68, 68, 0.1)', // Red transparent
        borderColor: '#ef4444', // Red danger/heat
        borderWidth: 2,
        pointBackgroundColor: '#ef4444',
        pointRadius: 4,
        tension: 0.4, // Curva suave (Electrocardiograma)
        fill: true,
        data: historialTemp.value.map(h => parseFloat(h.temperatura))
      }
    ]
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false }
  },
  scales: {
    y: {
      beginAtZero: false,
      grid: { color: 'rgba(255, 255, 255, 0.05)' },
      ticks: { color: '#94a3b8' },
      suggestedMin: 15,
      suggestedMax: 40
    },
    x: {
      grid: { display: false },
      ticks: { color: '#94a3b8', maxTicksLimit: 7 }
    }
  },
  animation: {
    duration: 500
  }
}
</script>

<template>
  <div class="dashboard">
    <!-- Error notification -->
    <div v-if="isError" class="glass-panel error-panel mb-4">
      <i class="ph-fill ph-warning-circle text-danger"></i>
      <span>No se pudo conectar con el servidor. ¿Está el backend corriendo?</span>
    </div>

    <!-- Quick Stats Skeletons -->
    <div v-if="isLoading && activos.length === 0" class="stats-grid mb-4">
      <div v-for="i in 3" :key="i" class="glass-panel stat-card">
        <div class="skeleton skeleton-circle"></div>
        <div style="flex: 1">
          <div class="skeleton skeleton-text short"></div>
          <div class="skeleton skeleton-text" style="width: 30%"></div>
        </div>
      </div>
    </div>

    <!-- Quick Stats -->
    <div v-else class="stats-grid mb-4">
      <div class="glass-panel stat-card">
        <div class="stat-icon bg-primary-light">
          <i class="ph ph-cpu text-primary"></i>
        </div>
        <div>
          <p class="text-muted">Total Activos</p>
          <h3 class="stat-value animate-counter">{{ activos.length }}</h3>
        </div>
      </div>
      <div class="glass-panel stat-card">
        <div class="stat-icon bg-warning-light">
          <i class="ph ph-warning text-warning"></i>
        </div>
        <div>
          <p class="text-muted">Alertas Recientes</p>
          <h3 class="stat-value animate-counter">{{ alertas.length }}</h3>
        </div>
      </div>
      <div class="glass-panel stat-card">
        <div class="stat-icon bg-success-light">
          <i class="ph ph-thermometer text-success"></i>
        </div>
        <div>
          <p class="text-muted">Sensores Activos</p>
          <h3 class="stat-value animate-counter">{{ activos.filter(a => a.ultima_temperatura !== '---').length }}</h3>
        </div>
      </div>
    </div>

    <div class="main-grid">
      <!-- Sección Izquierda: Sensores y Gráfico -->
      <div class="left-section">
        <!-- Gráfico de Líneas "En Vivo" -->
        <div class="glass-panel p-4 mb-4 relative" style="height: 320px;">
          <div class="flex-between mb-3">
            <h3 class="section-title" style="font-size: 1rem;">
              <span class="live-dot"></span> Temperatura en Vivo (Activo #1)
            </h3>
          </div>
          
          <div v-if="isLoading && historialTemp.length === 0" class="skeleton" style="height: 240px; width: 100%;"></div>
          <Line v-else-if="historialTemp.length > 0" :data="chartData" :options="chartOptions" />
          <div v-else class="text-center text-muted mt-4">Sin historial para graficar. Usa el simulador.</div>
        </div>

        <h3 class="section-title mb-3">Monitoreo de Sensores</h3>
        <div class="sensors-grid">
          <!-- Skeletons para sensores -->
          <template v-if="isLoading && activos.length === 0">
            <div v-for="i in 4" :key="i" class="glass-panel p-4" style="height: 180px;">
              <div class="skeleton skeleton-text short mb-4"></div>
              <div class="skeleton skeleton-text"></div>
              <div class="skeleton skeleton-text"></div>
            </div>
          </template>
          
          <!-- Sensores Reales -->
          <template v-else>
            <SensorCard v-for="activo in activos" :key="activo.id_activo" :activo="activo" />
            <div v-if="activos.length === 0 && !isError" class="glass-panel p-4 text-center text-muted col-span-full">
              No hay activos registrados.
            </div>
          </template>
        </div>
      </div>

      <!-- Sección Derecha: Alertas Panel -->
      <div class="alerts-section">
        <h3 class="section-title mb-3">Últimas Alertas</h3>
        <div class="glass-panel p-4 alerts-container">
          <template v-if="isLoading && alertas.length === 0">
            <div v-for="i in 5" :key="i" class="mb-4">
              <div class="skeleton skeleton-text"></div>
              <div class="skeleton skeleton-text short"></div>
            </div>
          </template>

          <template v-else>
            <div v-if="alertas.length === 0" class="text-center text-muted py-4">
              No hay alertas recientes.
            </div>
            <div v-for="alerta in alertas" :key="alerta.id_alerta" class="alert-item">
              <div class="alert-icon" :class="alerta.tipo === 'temperatura' ? 'bg-danger-light' : 'bg-warning-light'">
                <i :class="alerta.tipo === 'temperatura' ? 'ph ph-thermometer-hot text-danger' : 'ph ph-drop text-warning'"></i>
              </div>
              <div class="alert-content">
                <p class="alert-msg">{{ alerta.mensaje }}</p>
                <span class="alert-time">{{ alerta.fecha }}</span>
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.mb-4 { margin-bottom: 1.5rem; }
.mb-3 { margin-bottom: 1rem; }
.p-4 { padding: 1.5rem; }
.py-4 { padding-top: 1.5rem; padding-bottom: 1.5rem; }
.text-center { text-align: center; }
.mt-4 { margin-top: 1.5rem; }
.relative { position: relative; }
.flex-between { display: flex; justify-content: space-between; align-items: center; }

.error-panel {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 1.5rem;
  border-left: 4px solid var(--danger);
  color: #f87171;
}
.error-panel i { font-size: 1.5rem; }

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 1.25rem;
  padding: 1.25rem;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.75rem;
}

.bg-primary-light { background: rgba(59, 130, 246, 0.1); }
.bg-warning-light { background: rgba(245, 158, 11, 0.1); }
.bg-success-light { background: rgba(16, 185, 129, 0.1); }
.bg-danger-light { background: rgba(239, 68, 68, 0.1); }

.text-primary { color: var(--primary); }
.text-warning { color: var(--warning); }
.text-success { color: var(--success); }
.text-danger { color: var(--danger); }
.text-muted { color: var(--text-muted); font-size: 0.875rem; margin-bottom: 0.25rem; }

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  line-height: 1;
}

.animate-counter {
  animation: slideUpFade 0.5s ease-out forwards;
}

@keyframes slideUpFade {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
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

.left-section {
  display: flex;
  flex-direction: column;
}

.section-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-main);
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.live-dot {
  width: 8px;
  height: 8px;
  background-color: var(--danger);
  border-radius: 50%;
  box-shadow: 0 0 6px var(--danger);
  animation: pulse-red 1.5s infinite;
}

@keyframes pulse-red {
  0% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.7); }
  70% { transform: scale(1); box-shadow: 0 0 0 6px rgba(239, 68, 68, 0); }
  100% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(239, 68, 68, 0); }
}

.sensors-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1rem;
}

.col-span-full { grid-column: 1 / -1; }

.alerts-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  max-height: 600px;
  overflow-y: auto;
}

.alert-item {
  display: flex;
  gap: 1rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--border-color);
}
.alert-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.alert-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.25rem;
  flex-shrink: 0;
}

.alert-content {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.alert-msg {
  font-size: 0.875rem;
  line-height: 1.4;
}

.alert-time {
  font-size: 0.75rem;
  color: var(--text-muted);
}
</style>
