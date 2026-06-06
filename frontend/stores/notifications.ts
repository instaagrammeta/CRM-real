import { defineStore } from 'pinia'
import type { Notification, WsEvent } from '~/types/api'

export const useNotificationsStore = defineStore('notifications', {
  state: () => ({
    items: [] as Notification[],
    unread: 0,
    loaded: false,
    loading: false,
  }),

  getters: {
    hasUnread: (s) => s.unread > 0,
    recent: (s) => s.items.slice(0, 10),
  },

  actions: {
    async fetch() {
      this.loading = true
      try {
        const api = useApi()
        const items = await api.get<Notification[]>('/api/notifications', { limit: 50 })
        this.items = items || []
        this.unread = this.items.filter(n => !n.is_read).length
        this.loaded = true
      } catch {
        // ignore
      } finally {
        this.loading = false
      }
    },

    pushFromWs(ev: WsEvent<Notification>) {
      if (!ev?.data) return
      const n = ev.data
      this.items.unshift(n)
      if (!n.is_read) this.unread++

      // toast
      try {
        useToast().info(n.title, n.body)
      } catch {/* */}
    },

    async markRead(id: number) {
      try {
        const api = useApi()
        await api.post(`/api/notifications/${id}/read`)
        const item = this.items.find(n => n.id === id)
        if (item && !item.is_read) {
          item.is_read = true
          this.unread = Math.max(0, this.unread - 1)
        }
      } catch {/* */}
    },

    async markAllRead() {
      try {
        const api = useApi()
        await api.post('/api/notifications/read-all')
        this.items.forEach(n => (n.is_read = true))
        this.unread = 0
      } catch {/* */}
    },

    async remove(id: number) {
      try {
        const api = useApi()
        await api.delete(`/api/notifications/${id}`)
        const idx = this.items.findIndex(n => n.id === id)
        if (idx >= 0) {
          if (!this.items[idx].is_read) this.unread = Math.max(0, this.unread - 1)
          this.items.splice(idx, 1)
        }
      } catch {/* */}
    },
  },
})
