import { ref } from 'vue'

export const toasts = ref([])
let nextId = 0

export const showToast = (message, type = 'success', duration = 3000) => {
  const id = nextId++
  toasts.value.push({ id, message, type })
  
  setTimeout(() => {
    removeToast(id)
  }, duration)
}

export const removeToast = (id) => {
  const index = toasts.value.findIndex(t => t.id === id)
  if (index > -1) {
    toasts.value.splice(index, 1)
  }
}
