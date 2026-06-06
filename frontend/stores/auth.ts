import { defineStore } from 'pinia'
import type { User, LoginResponse } from '~/types/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    role: '' as string,
    token: '' as string,
    booted: false,
  }),

  getters: {
    isAuthenticated: (s) => !!s.token && !!s.user,
    isAdmin: (s) => s.role === 'admin',
    fullName: (s) => s.user?.full_name || '',
    avatar: (s) => s.user?.photo || '',
  },

  actions: {
    async restore() {
      if (process.server || this.booted) return
      const token = getAuthToken()
      if (!token) {
        this.booted = true
        return
      }
      this.token = token
      try {
        const api = useApi()
        const me = await api.get<User>('/api/me')
        this.user = me
        this.role = me.role
      } catch {
        this.clear()
      } finally {
        this.booted = true
      }
    },

    async login(login: string, password: string) {
      const api = useApi()
      const res = await api.post<LoginResponse>('/api/login', { login, password })
      if (!res.success) throw new Error(res.error || 'Login failed')
      setAuthToken(res.token)
      this.token = res.token
      this.role = res.role
      // also fetch full user so we have all fields
      try {
        this.user = await api.get<User>('/api/me')
      } catch {
        // fallback
        this.user = res.user as User
      }
      // start ws & notifications
      useWebSocket().connect()
      const notif = useNotificationsStore()
      notif.fetch().catch(() => {})
      return res
    },

    async logout() {
      try {
        const api = useApi()
        await api.post('/api/logout')
      } catch {
        // ignore
      }
      this.clear()
      useWebSocket().disconnect()
      await navigateTo('/login')
    },

    clear() {
      this.user = null
      this.role = ''
      this.token = ''
      setAuthToken(null)
      const notif = useNotificationsStore()
      notif.$reset?.()
    },
  },
})
