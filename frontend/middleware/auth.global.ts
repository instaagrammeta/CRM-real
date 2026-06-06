// Global middleware: redirects to /login when route requires auth.

const PUBLIC_ROUTES = new Set<string>([
  '/login',
  '/ipoteka',
  '/rasrochka',
])

const isPublic = (path: string): boolean => {
  if (PUBLIC_ROUTES.has(path)) return true
  // public sub-pages: /ipoteka/<slug>, /rasrochka/<slug>
  if (path.startsWith('/ipoteka/')) return true
  if (path.startsWith('/rasrochka/')) return true
  return false
}

export default defineNuxtRouteMiddleware(async (to) => {
  // Server-side rendering: skip; auth lives in localStorage on the client.
  if (process.server) return

  const auth = useAuthStore()
  if (!auth.booted) {
    await auth.restore()
  }

  // Already on /login but authenticated -> go home.
  if (to.path === '/login' && auth.isAuthenticated) {
    return navigateTo('/')
  }

  if (isPublic(to.path)) return

  if (!auth.isAuthenticated) {
    return navigateTo({
      path: '/login',
      query: to.fullPath !== '/' ? { redirect: to.fullPath } : undefined,
    })
  }

  // /admin/* requires admin role.
  if (to.path.startsWith('/admin') && !auth.isAdmin) {
    return navigateTo('/')
  }
})
