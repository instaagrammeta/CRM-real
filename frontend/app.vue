<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
  <Toaster />
</template>

<script setup lang="ts">
// Bootstrap auth + WebSocket once on the client side.
const auth = useAuthStore()
const notifications = useNotificationsStore()

if (process.client) {
  // Try to restore session from localStorage / cookie on app boot.
  await auth.restore()

  if (auth.isAuthenticated) {
    notifications.fetch().catch(() => {})
    useWebSocket().connect()
  }
}
</script>
