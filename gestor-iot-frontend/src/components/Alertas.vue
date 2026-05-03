<script setup>
import { ref, computed, onMounted } from 'vue'
import { showToast } from '../utils/toast.js'
import jsPDF from 'jspdf'
import 'jspdf-autotable'

const alertas = ref([])
const isLoading = ref(true)
const searchQuery = ref('')

const fetchAlertas = async () => {
  isLoading.value = true
  try {
    const res = await fetch('http://localhost:8080/api/alertas')
    if (res.ok) {
      alertas.value = await res.json()
    }
  } catch (error) {
    console.error("Error fetching alertas:", error)
    showToast("Error al conectar con el servidor", "error")
  } finally {
    isLoading.value = false
  }
}

const filteredAlertas = computed(() => {
  if (!searchQuery.value) return alertas.value
  const query = searchQuery.value.toLowerCase()
  return alertas.value.filter(a => 
    a.mensaje.toLowerCase().includes(query) || 
    a.tipo.toLowerCase().includes(query) ||
    a.id_activo.toString().includes(query)
  )
})

const exportToCSV = () => {
  if (alertas.value.length === 0) {
    showToast("No hay datos para exportar", "warning")
    return
  }
  
  const headers = ["ID Alerta", "ID Activo", "Tipo", "Mensaje", "Fecha"]
  const rows = alertas.value.map(a => [
    a.id_alerta,
    a.id_activo,
    a.tipo,
    `"${a.mensaje}"`,
    `"${a.fecha}"`
  ])
  
  const csvContent = [headers.join(","), ...rows.map(r => r.join(","))].join("\n")
  
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement("a")
  link.setAttribute("href", url)
  link.setAttribute("download", `alertas_export_${new Date().toISOString().split('T')[0]}.csv`)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  showToast("Historial de alertas exportado a Excel exitosamente", "success")
}

const exportToPDF = () => {
  if (alertas.value.length === 0) {
    showToast("No hay datos para exportar", "warning")
    return
  }
  
  const doc = new jsPDF()
  
  doc.setFontSize(18)
  doc.setTextColor(40, 40, 40)
  doc.text("Historial de Alertas IoT", 14, 22)
  
  doc.setFontSize(11)
  doc.setTextColor(100, 100, 100)
  doc.text(`Generado el: ${new Date().toLocaleString()}`, 14, 30)
  
  const headers = [["ID Alerta", "Activo", "Tipo", "Mensaje", "Fecha"]]
  const data = alertas.value.map(a => [
    `#${a.id_alerta}`,
    `Activo #${a.id_activo}`,
    a.tipo.toUpperCase(),
    a.mensaje,
    a.fecha
  ])
  
  doc.autoTable({
    startY: 36,
    head: headers,
    body: data,
    theme: 'grid',
    headStyles: { fillColor: [239, 68, 68] }, // Red header for alerts
    alternateRowStyles: { fillColor: [254, 242, 242] },
    margin: { top: 36 }
  })
  
  doc.save(`alertas_historial_${new Date().getTime()}.pdf`)
  showToast("Reporte PDF generado exitosamente", "success")
}

onMounted(() => {
  fetchAlertas()
})
</script>

<template>
  <div class="alertas">
    <div class="header-section mb-4">
      <h2 class="section-title">Historial de Alertas</h2>
      
      <div class="header-actions">
        <div class="search-box">
          <i class="ph ph-magnifying-glass text-muted"></i>
          <input type="text" v-model="searchQuery" placeholder="Buscar mensaje, tipo, id..." class="search-input">
        </div>
        
        <button class="btn btn-outline" @click="exportToCSV" title="Exportar a CSV">
          <i class="ph ph-file-csv"></i> Excel
        </button>
        <button class="btn btn-outline" @click="exportToPDF" title="Exportar a PDF" style="border-color: #ef4444; color: #ef4444;">
          <i class="ph ph-file-pdf"></i> PDF
        </button>
        <button class="btn btn-primary" @click="fetchAlertas">
          <i class="ph ph-arrows-clockwise"></i> Actualizar
        </button>
      </div>
    </div>

    <div class="glass-panel table-container">
      <table class="w-full">
        <thead>
          <tr>
            <th>ID</th>
            <th>Activo Asociado</th>
            <th>Tipo</th>
            <th>Mensaje</th>
            <th>Fecha</th>
          </tr>
        </thead>
        <tbody>
          <!-- Skeleton Loading State -->
          <template v-if="isLoading">
            <tr v-for="i in 5" :key="'skel'+i">
              <td><div class="skeleton skeleton-text" style="width: 40px"></div></td>
              <td><div class="skeleton skeleton-text short"></div></td>
              <td><div class="skeleton skeleton-text" style="width: 60px; border-radius: 99px;"></div></td>
              <td><div class="skeleton skeleton-text"></div></td>
              <td><div class="skeleton skeleton-text" style="width: 120px"></div></td>
            </tr>
          </template>

          <!-- Data State -->
          <template v-else>
            <tr v-for="alerta in filteredAlertas" :key="alerta.id_alerta">
              <td class="text-muted">#{{ alerta.id_alerta }}</td>
              <td class="font-medium">Activo #{{ alerta.id_activo }}</td>
              <td>
                <span class="badge" :class="alerta.tipo === 'temperatura' ? 'badge-danger' : 'badge-warning'">
                  {{ alerta.tipo }}
                </span>
              </td>
              <td>{{ alerta.mensaje }}</td>
              <td class="text-muted">{{ alerta.fecha }}</td>
            </tr>
            <tr v-if="filteredAlertas.length === 0">
              <td colspan="5" class="text-center p-4 text-muted">No se encontraron alertas.</td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.mb-4 { margin-bottom: 1.5rem; }
.p-4 { padding: 1.5rem; }
.text-center { text-align: center; }
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
</style>
