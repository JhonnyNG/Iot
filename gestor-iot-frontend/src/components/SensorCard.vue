<script setup>
import { computed } from 'vue'

const props = defineProps({
  activo: {
    type: Object,
    required: true
  }
})

// Dynamic classes based on temperature
const tempStatus = computed(() => {
  const temp = parseFloat(props.activo.ultima_temperatura)
  if (isNaN(temp)) return 'neutral'
  if (temp > 25) return 'danger'
  if (temp > 20) return 'warning'
  return 'success'
})

const humStatus = computed(() => {
  const hum = parseFloat(props.activo.ultima_humedad)
  if (isNaN(hum)) return 'neutral'
  if (hum > 60) return 'danger'
  if (hum < 30) return 'warning'
  return 'success'
})
</script>

<template>
  <div class="sensor-card glass-panel">
    <div class="card-header">
      <div class="header-info">
        <h4>{{ activo.nombre }}</h4>
        <span class="text-muted text-xs">ID: {{ activo.id_activo }}</span>
      </div>
      <div class="stock-badge">
        <i class="ph ph-package"></i>
        <span>{{ activo.stock_actual }}</span>
      </div>
    </div>
    
    <div class="telemetry-grid">
      <!-- Temperatura -->
      <div class="telemetry-item">
        <div class="telemetry-icon" :class="`icon-${tempStatus}`">
          <i class="ph ph-thermometer"></i>
        </div>
        <div class="telemetry-data">
          <span class="value">{{ activo.ultima_temperatura }}<small v-if="activo.ultima_temperatura !== '---'">°C</small></span>
          <span class="label">Temperatura</span>
        </div>
      </div>

      <!-- Humedad -->
      <div class="telemetry-item">
        <div class="telemetry-icon" :class="`icon-${humStatus}`">
          <i class="ph ph-drop"></i>
        </div>
        <div class="telemetry-data">
          <span class="value">{{ activo.ultima_humedad }}<small v-if="activo.ultima_humedad !== '---'">%</small></span>
          <span class="label">Humedad</span>
        </div>
      </div>
    </div>

    <div class="card-footer">
      <div class="location">
        <i class="ph ph-map-pin"></i>
        <span>{{ activo.ultima_ubicacion !== '---' ? activo.ultima_ubicacion : 'Sin ubicar' }}</span>
      </div>
      <div class="timestamp">
        <i class="ph ph-clock"></i>
        <span>{{ activo.ultima_fecha ? activo.ultima_fecha.split(' ')[1] : '--:--' }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.sensor-card {
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  transition: transform 0.2s, box-shadow 0.2s;
}

.sensor-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.2);
  border-color: rgba(255, 255, 255, 0.2);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.header-info h4 {
  font-size: 1.125rem;
  font-weight: 600;
  margin-bottom: 0.125rem;
}

.text-xs { font-size: 0.75rem; }
.text-muted { color: var(--text-muted); }

.stock-badge {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  background: rgba(255, 255, 255, 0.1);
  padding: 0.25rem 0.625rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 500;
}

.telemetry-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  background: rgba(0, 0, 0, 0.2);
  padding: 1rem;
  border-radius: var(--radius-md);
}

.telemetry-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.telemetry-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.25rem;
}

.icon-danger { background: rgba(239, 68, 68, 0.15); color: var(--danger); }
.icon-warning { background: rgba(245, 158, 11, 0.15); color: var(--warning); }
.icon-success { background: rgba(16, 185, 129, 0.15); color: var(--success); }
.icon-neutral { background: rgba(148, 163, 184, 0.15); color: var(--text-muted); }

.telemetry-data {
  display: flex;
  flex-direction: column;
}

.value {
  font-size: 1.125rem;
  font-weight: 700;
  line-height: 1.2;
}

.value small {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--text-muted);
  margin-left: 1px;
}

.label {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.75rem;
  color: var(--text-muted);
  border-top: 1px solid var(--border-color);
  padding-top: 0.75rem;
  margin-top: -0.25rem;
}

.location, .timestamp {
  display: flex;
  align-items: center;
  gap: 0.375rem;
}
</style>
