// useApi — обёртка над $fetch с автоматическим JWT и обработкой 401.

const TOKEN_KEY = 'auth_token'

export const getAuthToken = (): string | null => {
  if (process.server) return null
  return localStorage.getItem(TOKEN_KEY)
}

export const setAuthToken = (token: string | null) => {
  if (process.server) return
  if (token) localStorage.setItem(TOKEN_KEY, token)
  else localStorage.removeItem(TOKEN_KEY)
}

export interface ApiOptions extends Record<string, any> {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  body?: any
  query?: Record<string, any>
  headers?: Record<string, string>
}

export const useApi = () => {
  const config = useRuntimeConfig()
  const base = config.public.apiBase

  const request = async <T = any>(path: string, opts: ApiOptions = {}): Promise<T> => {
    const token = getAuthToken()
    const headers: Record<string, string> = { ...(opts.headers || {}) }

    if (token && !headers.Authorization) {
      headers.Authorization = `Bearer ${token}`
    }

    // Don't set Content-Type for FormData (browser sets the boundary).
    const isForm = typeof FormData !== 'undefined' && opts.body instanceof FormData
    if (!isForm && opts.body && typeof opts.body !== 'string' && !headers['Content-Type']) {
      headers['Content-Type'] = 'application/json'
    }

    try {
      const url = path.startsWith('http') ? path : `${base}${path}`
      const res = await $fetch<T>(url, {
        ...opts,
        headers,
        credentials: 'include',
      })
      return res
    } catch (err: any) {
      // 401 => logout & redirect to /login
      const status = err?.response?.status || err?.statusCode
      if (status === 401 && process.client) {
        setAuthToken(null)
        const auth = useAuthStore()
        auth.clear()
        const route = useRoute()
        if (route.path !== '/login') {
          await navigateTo({ path: '/login', query: { redirect: route.fullPath } })
        }
      }
      throw err
    }
  }

  return {
    get:    <T = any>(path: string, query?: Record<string, any>) => request<T>(path, { method: 'GET', query }),
    post:   <T = any>(path: string, body?: any) => request<T>(path, { method: 'POST', body }),
    put:    <T = any>(path: string, body?: any) => request<T>(path, { method: 'PUT', body }),
    patch:  <T = any>(path: string, body?: any) => request<T>(path, { method: 'PATCH', body }),
    delete: <T = any>(path: string) => request<T>(path, { method: 'DELETE' }),
    upload: <T = any>(path: string, form: FormData) => request<T>(path, { method: 'POST', body: form }),
    raw:    request,
    apiBase: base,
  }
}
