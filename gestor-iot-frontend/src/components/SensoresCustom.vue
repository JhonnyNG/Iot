<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { showToast } from '../utils/toast.js'

const activos = ref([])
const activoSeleccionado = ref(null)
const lecturas = ref([])
const isLoadingLecturas = ref(false)
const isEnviando = ref(false)

// Campos dinámicos del sensor custom
const campos = ref([{ clave: '', valor: '', tipo: 'number' }])

const tipoOpciones = ['number', 'text', 'boolean']

const agregarCampo = () => {
  campos.value.push({ clave: '', valor: '', tipo: 'number' })
}

const eliminarCampo = (i) => {
  if (campos.value.length > 1) campos.value.splice(i, 1)
}

const fetchActivos = async () => {
  try {
    const res = await fetch('http://localhost:8080/api/activos')
    if (res.ok) {
      activos.value = await res.json()
      if (activos.value.length > 0 && !activoSeleccionado.value) {
        activoSeleccionado.value = activos.value[0].id_activo
        await fetchLecturas()
      }
    }
  } catch (e) {
    console.error(e)
  }
}

const fetchLecturas = async () => {
  if (!activoSeleccionado.value) return
  isLoadingLecturas.value = true
  try {
    const res = await fetch(`http://localhost:8080/api/sensores/custom/${activoSeleccionado.value}`)
    if (res.ok) lecturas.value = await res.json()
  } catch (e) {
    console.error(e)
  } finally {
    isLoadingLecturas.value = false
  }
}

const enviarTelemetria = async () => {
  if (!activoSeleccionado.value) {
    showToast('Selecciona un activo primero', 'error')
    return
  }
  const camposValidos = campos.value.filter(c => c.clave.trim() !== '')
  if (camposValidos.length === 0) {
    showToast('Agrega al menos un campo con nombre', 'error')
    return
  }

  const payload = { id_activo: Number(activoSeleccionado.value) }
  for (const c of camposValidos) {
    if (c.tipo === 'number') payload[c.clave] = parseFloat(c.valor) || 0
    else if (c.tipo === 'boolean') payload[c.clave] = c.valor === 'true'
    else payload[c.clave] = c.valor
  }

  isEnviando.value = true
  try {
    const res = await fetch('http://localhost:8080/api/sensores/custom', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    if (res.ok) {
      showToast(`✓ Sensor enviado con ${camposValidos.length} campo(s)`, 'success')
      await fetchLecturas()
    } else {
      const err = await res.json()
      showToast(err.error || 'Error al enviar', 'error')
    }
  } catch (e) {
    showToast('No se pudo conectar al servidor', 'error')
  } finally {
    isEnviando.value = false
  }
}

const camposDeUltimaLectura = (lectura) => {
  return Object.entries(lectura).filter(([k]) => k !== '_timestamp')
}

const getIconForKey = (key) => {
  const k = key.toLowerCase()
  if (k.includes('temp')) return 'ph-thermometer'
  if (k.includes('hum')) return 'ph-drop'
  if (k.includes('lat') || k.includes('lon') || k.includes('gps')) return 'ph-map-pin'
  if (k.includes('vel') || k.includes('speed')) return 'ph-gauge'
  if (k.includes('pres')) return 'ph-wind'
  if (k.includes('co2') || k.includes('gas')) return 'ph-cloud'
  if (k.includes('luz') || k.includes('light') || k.includes('lux')) return 'ph-sun'
  if (k.includes('volt') || k.includes('bat')) return 'ph-battery-charging'
  return 'ph-waves'
}

const getColorForKey = (key) => {
  const k = key.toLowerCase()
  if (k.includes('temp')) return '#ef4444'
  if (k.includes('hum')) return '#3b82f6'
  if (k.includes('lat') || k.includes('lon')) return '#10b981'
  if (k.includes('vel') || k.includes('speed')) return '#f59e0b'
  if (k.includes('pres')) return '#8b5cf6'
  if (k.includes('co2') || k.includes('gas')) return '#6366f1'
  if (k.includes('luz') || k.includes('light')) return '#fbbf24'
  if (k.includes('volt') || k.includes('bat')) return '#22d3ee'
  return '#94a3b8'
}

let pollInterval
onMounted(async () => {
  await fetchActivos()
  pollInterval = setInterval(fetchLecturas, 6000)
})
onUnmounted(() => clearInterval(pollInterval))
</script>

<template>
  <div class="sensores-custom">

    <!-- Header banner -->
    <div class="nosql-banner glass-panel mb-4">
      <div class="banner-icon">
        <i class="ph-fill ph-broadcast"></i>
      </div>
      <div>
        <h2 class="banner-title">Sensores NoSQL Dinámicos</h2>
        <p class="banner-sub">Define cualquier sensor con campos libres — la capa JSONB de PostgreSQL almacena todo sin modificar el esquema.</p>
      </div>
      <div class="badge-nosql">JSONB</div>
    </div>

    <div class="two-col">
      <!-- Panel izquierdo: Constructor -->
      <div class="builder-panel glass-panel p-4">
        <h3 class="section-title mb-3">
          <i class="ph ph-plus-circle text-primary"></i> Constructor de Sensor
        </h3>

        <!-- Selector de activo -->
        <div class="field-group mb-3">
          <label class="field-label">Activo destino</label>
          <select v-model="activoSeleccionado" class="custom-select" @change="fetchLecturas">
            <option v-for="a in activos" :key="a.id_activo" :value="a.id_activo">
              {{ a.nombre }} (ID: {{ a.id_activo }})
            </option>
          </select>
        </div>

        <div class="divider mb-3"></div>

        <!-- Campos dinámicos -->
        <label class="field-label mb-2">Campos del sensor</label>
        <TransitionGroup name="field" tag="div" class="campos-list">
          <div v-for="(campo, i) in campos" :key="i" class="campo-row">
            <input
              v-model="campo.clave"
              class="campo-input"
              placeholder="nombre (ej: presion)"
              :id="`campo-clave-${i}`"
            />
            <input
              v-model="campo.valor"
              class="campo-input"
              :placeholder="campo.tipo === 'boolean' ? 'true / false' : 'valor'"
              :type="campo.tipo === 'number' ? 'number' : 'text'"
              :id="`campo-valor-${i}`"
            />
            <select v-model="campo.tipo" class="tipo-select">
              <option v-for="t in tipoOpciones" :key="t" :value="t">{{ t }}</option>
            </select>
            <button class="btn-remove" @click="eliminarCampo(i)" :disabled="campos.length === 1">
              <i class="ph ph-trash"></i>
            </button>
          </div>
        </TransitionGroup>

        <button class="btn-add-field mt-3 mb-4" @click="agregarCampo">
          <i class="ph ph-plus"></i> Agregar campo
        </button>

        <!-- Preview JSON -->
        <div class="json-preview mb-4">
          <span class="json-tag">Vista previa JSON (NoSQL)</span>
          <pre class="json-body">{{ JSON.stringify(
            Object.fromEntries([
              ['id_activo', activoSeleccionado],
              ...campos.filter(c=>c.clave).map(c=>[c.clave, c.tipo==='number'? (parseFloat(c.valor)||0) : c.valor])
            ]), null, 2) }}</pre>
        </div>

        <button
          class="btn-enviar"
          :class="{ 'is-loading': isEnviando }"
          @click="enviarTelemetria"
          :disabled="isEnviando"
          id="btn-enviar-sensor"
        >
          <i :class="isEnviando ? 'ph ph-spinner' : 'ph-fill ph-paper-plane-tilt'"></i>
          {{ isEnviando ? 'Enviando...' : 'Enviar al servidor' }}
        </button>
      </div>

      <!-- Panel derecho: Lecturas -->
      <div class="readings-panel">
        <h3 class="section-title mb-3">
          <span class="live-dot"></span> Lecturas Recibidas
          <span class="count-badge">{{ lecturas.length }}</span>
        </h3>

        <div v-if="isLoadingLecturas && lecturas.length === 0" class="glass-panel p-4">
          <div v-for="i in 3" :key="i" class="skeleton-reading mb-3"></div>
        </div>

        <div v-else-if="lecturas.length === 0" class="glass-panel p-4 text-center text-muted">
          <i class="ph ph-database" style="font-size:3rem; display:block; margin-bottom:1rem; opacity:0.3;"></i>
          Sin lecturas. Usa el constructor para enviar datos.
        </div>

        <TransitionGroup name="reading" tag="div" class="readings-list">
          <div v-for="(lectura, idx) in lecturas" :key="idx" class="reading-card glass-panel">
            <div class="reading-header">
              <span class="reading-index">#{{ lecturas.length - idx }}</span>
              <span class="reading-time">
                <i class="ph ph-clock"></i>
                {{ lectura._timestamp?.split(' ')[1] || '--:--' }}
              </span>
            </div>
            <div class="reading-fields">
              <div
                v-for="[key, val] in camposDeUltimaLectura(lectura)"
                :key="key"
                class="field-chip"
                :style="{ borderColor: getColorForKey(key) + '40', background: getColorForKey(key) + '12' }"
              >
                <i :class="`ph ${getIconForKey(key)}`" :style="{ color: getColorForKey(key) }"></i>
                <div class="chip-content">
                  <span class="chip-key">{{ key }}</span>
                  <span class="chip-val" :style="{ color: getColorForKey(key) }">{{ val }}</span>
                </div>
              </div>
            </div>
          </div>
        </TransitionGroup>
      </div>
    </div>
  </div>
</template>

<style scoped>
.sensores-custom { padding-bottom: 2rem; }
.mb-4 { margin-bottom: 1.5rem; }
.mb-3 { margin-bottom: 1rem; }
.mb-2 { margin-bottom: 0.5rem; }
.mt-3 { margin-top: 1rem; }
.p-4 { padding: 1.5rem; }
.text-center { text-align: center; }
.text-muted { color: var(--text-muted); font-size: 0.9rem; }
.text-primary { color: var(--primary); }

/* Banner */
.nosql-banner {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  padding: 1.25rem 1.75rem;
  border-left: 4px solid #8b5cf6;
  position: relative;
}
.banner-icon {
  width: 52px; height: 52px;
  border-radius: 14px;
  background: rgba(139,92,246,0.15);
  color: #8b5cf6;
  display: flex; align-items: center; justify-content: center;
  font-size: 1.75rem;
  flex-shrink: 0;
}
.banner-title { font-size: 1.1rem; font-weight: 700; margin-bottom: 0.2rem; }
.banner-sub { font-size: 0.82rem; color: var(--text-muted); line-height: 1.5; }
.badge-nosql {
  margin-left: auto;
  background: linear-gradient(135deg, #8b5cf6, #6366f1);
  color: white;
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 1.5px;
  padding: 0.4rem 0.8rem;
  border-radius: 8px;
  flex-shrink: 0;
}

/* Layout */
.two-col {
  display: grid;
  grid-template-columns: 420px 1fr;
  gap: 1.5rem;
  align-items: start;
}
@media (max-width: 900px) {
  .two-col { grid-template-columns: 1fr; }
}

/* Section title */
.section-title {
  font-size: 1rem; font-weight: 600;
  display: flex; align-items: center; gap: 0.5rem;
  color: var(--text-main);
}

/* Form */
.field-group { display: flex; flex-direction: column; gap: 0.4rem; }
.field-label { font-size: 0.8rem; color: var(--text-muted); font-weight: 500; text-transform: uppercase; letter-spacing: 0.5px; }

.custom-select {
  background: rgba(255,255,255,0.05);
  border: 1px solid var(--border-color);
  color: var(--text-main);
  padding: 0.6rem 0.75rem;
  border-radius: 8px;
  font-size: 0.9rem;
  width: 100%;
}
.custom-select:focus { outline: none; border-color: var(--primary); }

.divider { border: none; border-top: 1px solid var(--border-color); }

.campos-list { display: flex; flex-direction: column; gap: 0.5rem; }

.campo-row {
  display: grid;
  grid-template-columns: 1fr 1fr 80px 36px;
  gap: 0.5rem;
  align-items: center;
}
.campo-input {
  background: rgba(255,255,255,0.05);
  border: 1px solid var(--border-color);
  color: var(--text-main);
  padding: 0.5rem 0.6rem;
  border-radius: 8px;
  font-size: 0.85rem;
  width: 100%;
}
.campo-input:focus { outline: none; border-color: #8b5cf6; }
.tipo-select {
  background: rgba(255,255,255,0.05);
  border: 1px solid var(--border-color);
  color: var(--text-muted);
  padding: 0.5rem 0.3rem;
  border-radius: 8px;
  font-size: 0.75rem;
}
.btn-remove {
  background: rgba(239,68,68,0.1); border: none; color: #ef4444;
  border-radius: 8px; width: 36px; height: 36px;
  display: flex; align-items: center; justify-content: center;
  cursor: pointer; transition: background 0.2s;
}
.btn-remove:hover { background: rgba(239,68,68,0.25); }
.btn-remove:disabled { opacity: 0.3; cursor: not-allowed; }

.btn-add-field {
  display: flex; align-items: center; gap: 0.5rem;
  background: rgba(139,92,246,0.1);
  border: 1px dashed rgba(139,92,246,0.4);
  color: #8b5cf6;
  padding: 0.6rem 1rem; border-radius: 8px;
  font-size: 0.85rem; cursor: pointer;
  transition: all 0.2s; width: 100%; justify-content: center;
}
.btn-add-field:hover { background: rgba(139,92,246,0.2); }

/* JSON Preview */
.json-preview {
  background: rgba(0,0,0,0.3);
  border-radius: 10px;
  border: 1px solid rgba(139,92,246,0.2);
  overflow: hidden;
}
.json-tag {
  display: block;
  font-size: 0.7rem; letter-spacing: 1px; text-transform: uppercase;
  color: #8b5cf6; padding: 0.4rem 0.75rem;
  background: rgba(139,92,246,0.1);
  border-bottom: 1px solid rgba(139,92,246,0.15);
}
.json-body {
  padding: 0.75rem; margin: 0;
  font-size: 0.78rem; color: #94a3b8;
  font-family: 'Fira Code', monospace;
  max-height: 140px; overflow-y: auto;
}

/* Send button */
.btn-enviar {
  width: 100%; padding: 0.85rem;
  background: linear-gradient(135deg, #8b5cf6, #6366f1);
  border: none; color: white;
  border-radius: 10px; font-size: 0.95rem; font-weight: 600;
  cursor: pointer; display: flex; align-items: center;
  justify-content: center; gap: 0.6rem;
  transition: all 0.2s; box-shadow: 0 4px 15px rgba(139,92,246,0.3);
}
.btn-enviar:hover:not(:disabled) { transform: translateY(-1px); box-shadow: 0 6px 20px rgba(139,92,246,0.45); }
.btn-enviar:disabled { opacity: 0.6; cursor: not-allowed; }
.btn-enviar.is-loading i { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

/* Live dot */
.live-dot {
  width: 8px; height: 8px;
  background: #8b5cf6; border-radius: 50%;
  box-shadow: 0 0 6px #8b5cf6;
  animation: pulse-purple 1.5s infinite;
}
@keyframes pulse-purple {
  0% { box-shadow: 0 0 0 0 rgba(139,92,246,0.7); }
  70% { box-shadow: 0 0 0 6px rgba(139,92,246,0); }
  100% { box-shadow: 0 0 0 0 rgba(139,92,246,0); }
}
.count-badge {
  margin-left: auto;
  background: rgba(139,92,246,0.15);
  color: #8b5cf6; font-size: 0.75rem; font-weight: 700;
  padding: 0.15rem 0.6rem; border-radius: 9999px;
}

/* Readings */
.readings-list { display: flex; flex-direction: column; gap: 0.75rem; }
.reading-card { padding: 1rem 1.25rem; }
.reading-header {
  display: flex; justify-content: space-between;
  align-items: center; margin-bottom: 0.75rem;
}
.reading-index {
  font-size: 0.75rem; font-weight: 700;
  color: #8b5cf6; background: rgba(139,92,246,0.12);
  padding: 0.1rem 0.5rem; border-radius: 4px;
}
.reading-time {
  display: flex; align-items: center; gap: 0.3rem;
  font-size: 0.75rem; color: var(--text-muted);
}
.reading-fields {
  display: flex; flex-wrap: wrap; gap: 0.5rem;
}
.field-chip {
  display: flex; align-items: center; gap: 0.5rem;
  padding: 0.4rem 0.75rem;
  border-radius: 8px; border: 1px solid transparent;
  min-width: 100px;
}
.chip-content { display: flex; flex-direction: column; }
.chip-key { font-size: 0.68rem; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.5px; }
.chip-val { font-size: 0.9rem; font-weight: 700; }

/* Skeleton */
.skeleton-reading {
  height: 80px; border-radius: 10px;
  background: linear-gradient(90deg, rgba(255,255,255,0.04) 0%, rgba(255,255,255,0.08) 50%, rgba(255,255,255,0.04) 100%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}
@keyframes shimmer { to { background-position: -200% 0; } }

/* Transitions */
.field-enter-active, .field-leave-active { transition: all 0.25s ease; }
.field-enter-from { opacity: 0; transform: translateX(-10px); }
.field-leave-to { opacity: 0; transform: translateX(10px); }

.reading-enter-active { transition: all 0.35s ease; }
.reading-enter-from { opacity: 0; transform: translateY(-8px); }
</style>
