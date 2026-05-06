<script setup>
import { ref } from 'vue'
import Dashboard from './components/Dashboard.vue'
import Inventario from './components/Inventario.vue'
import Alertas from './components/Alertas.vue'
import Ajustes from './components/Ajustes.vue'
import SensoresCustom from './components/SensoresCustom.vue'
import { toasts, removeToast } from './utils/toast.js'

const isSidebarOpen = ref(true)
const currentView = ref('dashboard')

const setView = (view) => {
  currentView.value = view
}

const getIconForToast = (type) => {
  if (type === 'success') return 'ph-check-circle'
  if (type === 'error') return 'ph-warning-circle'
  return 'ph-info'
}

// Live Clock
import { onMounted, onUnmounted } from 'vue'
const currentTime = ref('')
let timeInterval
onMounted(() => {
  const updateTime = () => {
    const now = new Date()
    currentTime.value = now.toLocaleTimeString('es-ES', { hour12: false })
  }
  updateTime()
  timeInterval = setInterval(updateTime, 1000)
})
onUnmounted(() => clearInterval(timeInterval))
</script>

<template>
  <div class="layout-container">
    <!-- Sidebar -->
    <aside :class="['sidebar glass-panel', { 'sidebar-closed': !isSidebarOpen }]">
      <div class="sidebar-header">
        <i class="ph-fill ph-cpu text-primary" style="font-size: 2rem;"></i>
        <h2 v-if="isSidebarOpen">Gestor IoT</h2>
      </div>
      
      <nav class="sidebar-nav">
        <a href="#" class="nav-item" :class="{ active: currentView === 'dashboard' }" @click.prevent="setView('dashboard')">
          <i class="ph ph-squares-four"></i>
          <span v-if="isSidebarOpen">Dashboard</span>
        </a>
        <a href="#" class="nav-item" :class="{ active: currentView === 'inventario' }" @click.prevent="setView('inventario')">
          <i class="ph ph-package"></i>
          <span v-if="isSidebarOpen">Inventario</span>
        </a>
        <a href="#" class="nav-item" :class="{ active: currentView === 'alertas' }" @click.prevent="setView('alertas')">
          <i class="ph ph-warning-circle"></i>
          <span v-if="isSidebarOpen">Alertas</span>
        </a>
        <a href="#" class="nav-item" :class="{ active: currentView === 'sensores' }" @click.prevent="setView('sensores')">
          <i class="ph ph-broadcast"></i>
          <span v-if="isSidebarOpen">Sensores Custom</span>
        </a>
        <a href="#" class="nav-item" :class="{ active: currentView === 'ajustes' }" @click.prevent="setView('ajustes')">
          <i class="ph ph-gear"></i>
          <span v-if="isSidebarOpen">Ajustes</span>
        </a>
      </nav>

      <button class="toggle-sidebar" @click="isSidebarOpen = !isSidebarOpen">
        <i :class="isSidebarOpen ? 'ph ph-caret-left' : 'ph ph-caret-right'"></i>
      </button>
    </aside>

    <!-- Main Content -->
    <main class="main-content">
      <header class="topbar glass-panel">
        <div>
          <h1 style="font-size: 1.25rem; font-weight: 600;">
            <Transition name="page" mode="out-in">
              <span v-if="currentView === 'dashboard'" key="dash">Resumen del Sistema</span>
              <span v-else-if="currentView === 'inventario'" key="inv">Inventario</span>
              <span v-else-if="currentView === 'alertas'" key="alt">Historial de Alertas</span>
              <span v-else-if="currentView === 'sensores'" key="sen">Sensores NoSQL Custom</span>
              <span v-else-if="currentView === 'ajustes'" key="aju">Ajustes</span>
            </Transition>
          </h1>
          <p class="text-muted" style="font-size: 0.875rem;">Monitoreo en tiempo real</p>
        </div>
        <div class="system-indicators">
          <div class="clock-display">
            <i class="ph ph-clock"></i>
            {{ currentTime }}
          </div>
          <div class="status-indicator">
            <span class="ping-dot"></span>
            <span>Servidor Online</span>
          </div>
        </div>
        <div class="user-profile">
          <img src="https://ui-avatars.com/api/?name=Admin&background=3b82f6&color=fff" alt="User" class="avatar">
        </div>
      </header>

      <div class="scrollable-content">
        <!-- Page Transition Wrapper -->
        <Transition name="page" mode="out-in">
          <Dashboard v-if="currentView === 'dashboard'" key="view-dashboard" />
          <Inventario v-else-if="currentView === 'inventario'" key="view-inventario" />
          <Alertas v-else-if="currentView === 'alertas'" key="view-alertas" />
          <SensoresCustom v-else-if="currentView === 'sensores'" key="view-sensores" />
          <Ajustes v-else-if="currentView === 'ajustes'" key="view-ajustes" />
        </Transition>
      </div>
    </main>

    <!-- Global Toast Notifications Container -->
    <div class="toast-container">
      <TransitionGroup name="toast">
        <div v-for="toast in toasts" :key="toast.id" class="toast" :class="toast.type">
          <i class="ph" :class="getIconForToast(toast.type)" style="font-size: 1.5rem;"></i>
          <span>{{ toast.message }}</span>
          <button @click="removeToast(toast.id)" style="margin-left: auto; background: none; border: none; color: white; cursor: pointer;">
            <i class="ph ph-x"></i>
          </button>
        </div>
      </TransitionGroup>
    </div>
  </div>
</template>

<style scoped>
.layout-container {
  display: flex;
  height: 100vh;
  overflow: hidden;
  padding: 1rem;
  gap: 1rem;
}

.sidebar {
  width: 250px;
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease;
  position: relative;
  border-radius: var(--radius-lg);
}

.sidebar-closed {
  width: 80px;
}

.sidebar-header {
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  border-bottom: 1px solid var(--border-color);
}

.sidebar-header h2 {
  font-size: 1.25rem;
  white-space: nowrap;
}

.sidebar-nav {
  padding: 1rem 0;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem 1.5rem;
  color: var(--text-muted);
  text-decoration: none;
  transition: all 0.2s;
  white-space: nowrap;
}

.nav-item:hover, .nav-item.active {
  color: var(--text-main);
  background: rgba(255, 255, 255, 0.05);
  border-left: 3px solid var(--primary);
}

.nav-item i {
  font-size: 1.5rem;
}

.toggle-sidebar {
  position: absolute;
  bottom: 1rem;
  right: -1rem;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 50%;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  z-index: 10;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  min-width: 0; /* Prevents overflow */
}

.topbar {
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-radius: var(--radius-lg);
}

.system-indicators {
  display: flex;
  align-items: center;
  gap: 2rem;
  margin-left: auto;
  margin-right: 2rem;
}

.clock-display {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-family: monospace;
  font-size: 1.125rem;
  color: var(--primary);
  background: rgba(59, 130, 246, 0.1);
  padding: 0.5rem 1rem;
  border-radius: var(--radius-md);
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--success);
  font-size: 0.875rem;
  font-weight: 500;
  background: rgba(16, 185, 129, 0.1);
  padding: 0.5rem 1rem;
  border-radius: var(--radius-md);
}

.ping-dot {
  width: 10px;
  height: 10px;
  background-color: var(--success);
  border-radius: 50%;
  box-shadow: 0 0 8px var(--success);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.7); }
  70% { transform: scale(1); box-shadow: 0 0 0 6px rgba(16, 185, 129, 0); }
  100% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(16, 185, 129, 0); }
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
}

.scrollable-content {
  flex: 1;
  overflow-y: auto;
  padding-right: 0.5rem;
}

.text-primary {
  color: var(--primary);
}

.text-muted {
  color: var(--text-muted);
}
</style>
