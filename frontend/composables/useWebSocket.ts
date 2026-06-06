// useWebSocket — singleton-обёртка для WS-канала с авто-reconnect.

import type { WsEvent, WsEventType } from '~/types/api'

type Listener = (ev: WsEvent) => void

const _listeners = new Map<string, Set<Listener>>()
const _socket = ref<WebSocket | null>(null)
const _connected = ref(false)
let _reconnectTimer: ReturnType<typeof setTimeout> | null = null
let _reconnectAttempts = 0

const emit = (ev: WsEvent) => {
  const all = _listeners.get('*')
  if (all) all.forEach(fn => fn(ev))
  const typed = _listeners.get(ev.type)
  if (typed) typed.forEach(fn => fn(ev))
}

const buildUrl = (): string | null => {
  const config = useRuntimeConfig()
  const token = getAuthToken()
  if (!token) return null
  return `${config.public.wsBase}/ws?token=${encodeURIComponent(token)}`
}

const connect = () => {
  if (process.server) return
  if (_socket.value && _socket.value.readyState <= 1) return // already connected/ing

  const url = buildUrl()
  if (!url) return

  const ws = new WebSocket(url)
  _socket.value = ws

  ws.onopen = () => {
    _connected.value = true
    _reconnectAttempts = 0
  }

  ws.onmessage = (msg) => {
    try {
      const ev = JSON.parse(msg.data) as WsEvent
      emit(ev)
    } catch {
      // ignore
    }
  }

  ws.onerror = () => {
    _connected.value = false
  }

  ws.onclose = () => {
    _connected.value = false
    _socket.value = null
    // exponential backoff: 1s, 2s, 4s, 8s, max 30s
    const delay = Math.min(30000, 1000 * Math.pow(2, _reconnectAttempts))
    _reconnectAttempts++
    if (_reconnectTimer) clearTimeout(_reconnectTimer)
    _reconnectTimer = setTimeout(() => connect(), delay)
  }
}

const disconnect = () => {
  if (_reconnectTimer) {
    clearTimeout(_reconnectTimer)
    _reconnectTimer = null
  }
  if (_socket.value) {
    _socket.value.close()
    _socket.value = null
  }
  _connected.value = false
  _reconnectAttempts = 0
}

const on = (type: WsEventType | '*', fn: Listener): (() => void) => {
  if (!_listeners.has(type)) _listeners.set(type, new Set())
  _listeners.get(type)!.add(fn)
  return () => off(type, fn)
}

const off = (type: WsEventType | '*', fn: Listener) => {
  _listeners.get(type)?.delete(fn)
}

export const useWebSocket = () => ({
  connect,
  disconnect,
  on,
  off,
  connected: _connected,
  socket:    _socket,
})
