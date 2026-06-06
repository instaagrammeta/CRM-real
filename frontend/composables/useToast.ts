// useToast — мини сервиси toast notifications (без сторонних либ).

export type ToastType = 'success' | 'error' | 'info' | 'warning'

export interface Toast {
  id: number
  type: ToastType
  title: string
  message?: string
  duration?: number
}

const _toasts = ref<Toast[]>([])
let _seq = 0

const push = (t: Omit<Toast, 'id'>) => {
  const id = ++_seq
  const toast: Toast = { duration: 4000, ...t, id }
  _toasts.value.push(toast)
  if (toast.duration && toast.duration > 0) {
    setTimeout(() => remove(id), toast.duration)
  }
  return id
}

const remove = (id: number) => {
  _toasts.value = _toasts.value.filter(t => t.id !== id)
}

export const useToast = () => ({
  toasts: _toasts,
  remove,
  success: (title: string, message?: string) => push({ type: 'success', title, message }),
  error:   (title: string, message?: string) => push({ type: 'error',   title, message, duration: 6000 }),
  info:    (title: string, message?: string) => push({ type: 'info',    title, message }),
  warning: (title: string, message?: string) => push({ type: 'warning', title, message }),
})
