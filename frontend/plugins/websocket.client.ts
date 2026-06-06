// Client-side plugin: подписывается на WS-eventы и роутит их в stores/toasts.

export default defineNuxtPlugin(() => {
  const ws = useWebSocket()
  const notifications = useNotificationsStore()
  const toast = useToast()

  // notification -> store + toast (toast уже встроен в pushFromWs)
  ws.on('notification', (ev) => {
    notifications.pushFromWs(ev as any)
  })

  // лид-эффекты (кроме автора)
  ws.on('lead:created', (ev) => {
    const auth = useAuthStore()
    const data: any = ev?.data || {}
    if (data.author_id && data.author_id === auth.user?.id) return
    toast.info('Лиди нав', `${data.client_name || ''} ${data.phone ? '— ' + data.phone : ''}`)
  })

  // chat — только глобальный signal; страница чата сама подписывается.

  // tariff:expiring
  ws.on('tariff:expiring', (ev) => {
    const data: any = ev?.data || {}
    toast.warning('Муҳлати тариф наздик аст', data.phone_number || '')
  })
})
