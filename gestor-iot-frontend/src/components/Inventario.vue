<script setup>
import { ref, computed, onMounted } from 'vue'
import { showToast } from '../utils/toast.js'
import jsPDF from 'jspdf'
import 'jspdf-autotable'

const activos = ref([])
const isLoading = ref(true)
const showModal = ref(false)
const searchQuery = ref('')
const modalData = ref({
  id_activo: null,
  nombre: '',
  tipo: 'entrada',
  cantidad: 1,
  id_usuario: 1 // Usuario por defecto para prueba
})

const fetchActivos = async () => {
  isLoading.value = true
  try {
    const res = await fetch('http://localhost:8080/api/activos')
    if (res.ok) {
      activos.value = await res.json()
    }
  } catch (error) {
    console.error("Error fetching activos:", error)
    showToast("Error al conectar con el servidor", "error")
  } finally {
    isLoading.value = false
  }
}

const filteredActivos = computed(() => {
  if (!searchQuery.value) return activos.value
  const query = searchQuery.value.toLowerCase()
  return activos.value.filter(a => a.nombre.toLowerCase().includes(query) || a.id_activo.toString().includes(query))
})

const exportToCSV = () => {
  if (activos.value.length === 0) {
    showToast("No hay datos para exportar", "warning")
    return
  }
  
  const headers = ["ID Activo", "Nombre", "Stock Actual", "Stock Mínimo", "Umbral Temp", "Umbral Hum"]
  const rows = activos.value.map(a => [
    a.id_activo,
    `"${a.nombre}"`,
    a.stock_actual,
    a.stock_minimo,
    a.umbral_temperatura || 'N/A',
    a.umbral_humedad || 'N/A'
  ])
  
  const csvContent = [headers.join(","), ...rows.map(r => r.join(","))].join("\n")
  
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement("a")
  link.setAttribute("href", url)
  link.setAttribute("download", `inventario_export_${new Date().toISOString().split('T')[0]}.csv`)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  showToast("Inventario exportado a Excel exitosamente", "success")
}

const exportToPDF = () => {
  if (activos.value.length === 0) {
    showToast("No hay datos para exportar", "warning")
    return
  }
  
  const doc = new jsPDF()
  
  // Título
  doc.setFontSize(18)
  doc.setTextColor(40, 40, 40)
  doc.text("Reporte de Inventario IoT", 14, 22)
  
  // Subtítulo / Fecha
  doc.setFontSize(11)
  doc.setTextColor(100, 100, 100)
  doc.text(`Generado el: ${new Date().toLocaleString()}`, 14, 30)
  
  // Tabla
  const headers = [["ID", "Nombre", "Stock Actual", "Mínimo", "Umbrales"]]
  const data = activos.value.map(a => [
    `#${a.id_activo}`,
    a.nombre,
    a.stock_actual.toString(),
    a.stock_minimo.toString(),
    `${a.umbral_temperatura || '-'}°C / ${a.umbral_humedad || '-'}%`
  ])
  
  doc.autoTable({
    startY: 36,
    head: headers,
    body: data,
    theme: 'grid',
    headStyles: { fillColor: [59, 130, 246] },
    alternateRowStyles: { fillColor: [245, 247, 250] },
    margin: { top: 36 }
  })
  
  doc.save(`reporte_inventario_${new Date().getTime()}.pdf`)
  showToast("Reporte PDF generado exitosamente", "success")
}

const openModal = (activo, tipo) => {
  modalData.value.id_activo = activo.id_activo
  modalData.value.nombre = activo.nombre
  modalData.value.tipo = tipo
  modalData.value.cantidad = 1
  showModal.value = true
}

const registrarMovimiento = async () => {
  try {
    const res = await fetch('http://localhost:8080/api/movimiento', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        id_activo: modalData.value.id_activo,
        tipo: modalData.value.tipo,
        cantidad: modalData.value.cantidad,
        id_usuario: modalData.value.id_usuario
      })
    })
    
    if (res.ok) {
      showModal.value = false
      showToast(`Se registró la ${modalData.value.tipo} de ${modalData.value.cantidad} unidades`, "success")
      fetchActivos()
    } else {
      const errorData = await res.json()
      showToast(errorData.error || "No hay stock suficiente", "warning")
    }
  } catch (error) {
    console.error("Error al registrar movimiento:", error)
    showToast("Error de conexión al servidor", "error")
  }
}

onMounted(() => {
  fetchActivos()
})
</script>

<template>
  <div class="inventario">
    <div class="header-section mb-4">
      <h2 class="section-title">Gestión de Inventario</h2>
      
      <div class="header-actions">
        <div class="search-box">
          <i class="ph ph-magnifying-glass text-muted"></i>
          <input type="text" v-model="searchQuery" placeholder="Buscar por nombre o ID..." class="search-input">
        </div>
        
        <button class="btn btn-outline" @click="exportToCSV" title="Exportar a CSV">
          <i class="ph ph-file-csv"></i> Excel
        </button>
        <button class="btn btn-outline" @click="exportToPDF" title="Exportar a PDF" style="border-color: #ef4444; color: #ef4444;">
          <i class="ph ph-file-pdf"></i> PDF
        </button>
        <button class="btn btn-primary" @click="fetchActivos">
          <i class="ph ph-arrows-clockwise"></i> Actualizar
        </button>
      </div>
    </div>

    <div class="glass-panel table-container">
      <table class="w-full">
        <thead>
          <tr>
            <th>ID</th>
            <th>Nombre del Activo</th>
            <th>Stock Actual</th>
            <th>Stock Mínimo</th>
            <th>Umbrales (Temp/Hum)</th>
            <th class="text-right">Acciones</th>
          </tr>
        </thead>
        <tbody>
          <!-- Skeleton Loading State -->
          <template v-if="isLoading">
            <tr v-for="i in 5" :key="'skel'+i">
              <td><div class="skeleton skeleton-text" style="width: 40px"></div></td>
              <td><div class="skeleton skeleton-text short"></div></td>
              <td><div class="skeleton skeleton-text" style="width: 30px; border-radius: 99px;"></div></td>
              <td><div class="skeleton skeleton-text" style="width: 40px"></div></td>
              <td><div class="skeleton skeleton-text" style="width: 80px"></div></td>
              <td><div class="skeleton skeleton-text" style="width: 60px; float: right;"></div></td>
            </tr>
          </template>

          <!-- Data State -->
          <template v-else>
            <tr v-for="activo in filteredActivos" :key="activo.id_activo">
              <td class="text-muted">#{{ activo.id_activo }}</td>
              <td class="font-medium">{{ activo.nombre }}</td>
              <td>
                <span class="badge" :class="activo.stock_actual <= activo.stock_minimo ? 'badge-danger' : 'badge-success'">
                  {{ activo.stock_actual }}
                </span>
              </td>
              <td class="text-muted">{{ activo.stock_minimo }}</td>
              <td class="text-muted">
                {{ activo.umbral_temperatura ? activo.umbral_temperatura + '°C' : 'N/A' }} / 
                {{ activo.umbral_humedad ? activo.umbral_humedad + '%' : 'N/A' }}
              </td>
              <td class="actions">
                <button class="btn-action btn-entrada" @click="openModal(activo, 'entrada')" title="Registrar Entrada">
                  <i class="ph ph-plus"></i>
                </button>
                <button class="btn-action btn-salida" @click="openModal(activo, 'salida')" title="Registrar Salida">
                  <i class="ph ph-minus"></i>
                </button>
              </td>
            </tr>
            <tr v-if="filteredActivos.length === 0">
              <td colspan="6" class="text-center p-4 text-muted">No se encontraron activos.</td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>

    <!-- Modal para Entradas/Salidas con transition -->
    <Transition name="page">
      <div v-if="showModal" class="modal-backdrop">
        <div class="modal glass-panel">
          <div class="modal-header">
            <h3>Registrar {{ modalData.tipo === 'entrada' ? 'Entrada' : 'Salida' }}</h3>
            <button class="close-btn" @click="showModal = false"><i class="ph ph-x"></i></button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label>Activo</label>
              <input type="text" :value="modalData.nombre" disabled class="input-field disabled" />
            </div>
            <div class="form-group">
              <label>Cantidad a {{ modalData.tipo === 'entrada' ? 'añadir' : 'retirar' }}</label>
              <input type="number" v-model="modalData.cantidad" min="1" class="input-field" />
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-outline" @click="showModal = false">Cancelar</button>
            <button class="btn" :class="modalData.tipo === 'entrada' ? 'btn-primary' : 'btn-danger-solid'" @click="registrarMovimiento">
              Confirmar
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.mb-4 { margin-bottom: 1.5rem; }
.p-4 { padding: 1.5rem; }
.text-center { text-align: center; }
.text-right { text-align: right; }
.font-medium { font-weight: 500; }
.text-muted { color: var(--text-muted); }
.w-full { width: 100%; }

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-main);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.search-box {
  position: relative;
  display: flex;
  align-items: center;
}

.search-box i {
  position: absolute;
  left: 0.75rem;
  font-size: 1.25rem;
}

.search-input {
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid var(--border-color);
  color: white;
  border-radius: var(--radius-md);
  padding: 0.5rem 0.5rem 0.5rem 2.5rem;
  font-family: inherit;
  font-size: 0.875rem;
  outline: none;
  transition: border-color 0.2s;
  width: 250px;
}

.search-input:focus {
  border-color: var(--primary);
}

.table-container {
  overflow-x: auto;
}

table {
  border-collapse: collapse;
}

th, td {
  padding: 1rem 1.5rem;
  text-align: left;
  border-bottom: 1px solid var(--border-color);
}

th {
  font-weight: 600;
  color: var(--text-muted);
  font-size: 0.875rem;
  background: rgba(0, 0, 0, 0.2);
}

tr:last-child td {
  border-bottom: none;
}

tr:hover td {
  background: rgba(255, 255, 255, 0.02);
}

.actions {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
}

.btn-action {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  color: white;
}

.btn-entrada { background: var(--success); }
.btn-entrada:hover { background: #059669; }
.btn-salida { background: var(--warning); }
.btn-salida:hover { background: #d97706; }

/* Modal Styles */
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 50;
}

.modal {
  width: 100%;
  max-width: 400px;
  background: #1e293b;
}

.modal-header {
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  font-size: 1.125rem;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  font-size: 1.25rem;
}

.close-btn:hover { color: white; }

.modal-body {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-size: 0.875rem;
  color: var(--text-muted);
}

.input-field {
  padding: 0.75rem;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  background: rgba(0, 0, 0, 0.2);
  color: white;
  font-family: inherit;
}

.input-field.disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.modal-footer {
  padding: 1.25rem 1.5rem;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

.btn-danger-solid {
  background: var(--danger);
  color: white;
}
.btn-danger-solid:hover {
  background: #dc2626;
}
</style>
