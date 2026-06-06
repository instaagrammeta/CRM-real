<template>
  <div class="card !p-0 flex flex-col h-[calc(100vh-130px)]">
    <div class="px-5 py-3 border-b border-border/40 flex items-center justify-between">
      <div>
        <h2 class="font-semibold text-ink flex items-center gap-2">
          <i class="fa-solid fa-comments text-brand-600"></i> Чати ҷамъ
        </h2>
        <div class="text-xs text-ink-mute mt-0.5">
          <span class="inline-block w-1.5 h-1.5 rounded-full mr-1" :class="ws.connected.value ? 'bg-brand-500' : 'bg-gray-300'" />
          {{ ws.connected.value ? 'Real-time фаъол' : 'Real-time қатъ' }}
        </div>
      </div>
      <button class="btn-ghost !p-2" @click="loadHistory" title="Нав">
        <i class="fa-solid fa-rotate"></i>
      </button>
    </div>

    <!-- Messages -->
    <div ref="listEl" class="flex-1 overflow-y-auto px-4 py-4 space-y-2">
      <div v-if="loading" class="text-center py-10"><Spinner class="text-brand-600" size="lg" /></div>
      <EmptyState v-else-if="messages.length === 0" icon="fa-comment-dots" title="Паём нест"
                  description="Аввалин паёмро шумо нависед" />

      <div
        v-for="m in messages"
        :key="m.id"
        class="flex gap-2.5 group"
        :class="m.user_id === auth.user?.id ? 'flex-row-reverse' : ''"
      >
        <Avatar :name="m.user_name" size="sm" />
        <div
          class="max-w-[75%] rounded-2xl px-3.5 py-2 text-sm relative"
          :class="m.user_id === auth.user?.id
            ? 'bg-brand-600 text-white rounded-br-sm'
            : 'bg-page text-ink rounded-bl-sm'"
        >
          <div class="flex items-baseline gap-2 mb-0.5">
            <span class="text-xs font-semibold" :class="m.user_id === auth.user?.id ? 'text-white/90' : 'text-brand-700'">
              {{ m.user_name }}
            </span>
            <span class="text-[10px] opacity-70">{{ timeAgo(m.created_date) }}</span>
          </div>
          <div v-if="m.message" class="whitespace-pre-wrap break-words">{{ m.message }}</div>
          <a
            v-if="m.file_path"
            :href="apiBase + m.file_path"
            target="_blank"
            class="flex items-center gap-2 mt-1 text-xs underline opacity-90"
          >
            <i class="fa-solid fa-paperclip"></i> {{ m.file_name || 'Файл' }}
          </a>
          <button
            v-if="m.user_id === auth.user?.id"
            class="absolute -top-2 right-1 opacity-0 group-hover:opacity-100 transition-opacity bg-white text-red-600 rounded-full w-5 h-5 grid place-items-center text-[10px] shadow-card"
            @click="onDelete(m.id)"
            title="Нест"
          >
            <i class="fa-solid fa-xmark"></i>
          </button>
        </div>
      </div>
    </div>

    <!-- Composer -->
    <form
      class="border-t border-border/40 p-3 flex items-end gap-2"
      @submit.prevent="send"
    >
      <label class="btn-ghost !p-2.5 cursor-pointer" title="Илова кардани файл">
        <i class="fa-solid fa-paperclip"></i>
        <input type="file" class="hidden" @change="onFile" />
      </label>
      <textarea
        v-model="text"
        rows="1"
        class="input flex-1 resize-none !py-2.5"
        placeholder="Паёмро нависед..."
        @keydown.enter.exact.prevent="send"
      />
      <button type="submit" class="btn-primary !px-4 !py-2.5" :disabled="!canSend">
        <i class="fa-solid fa-paper-plane"></i>
      </button>
    </form>

    <div v-if="pendingFile" class="px-3 pb-2 text-xs text-ink-soft flex items-center gap-2">
      <i class="fa-solid fa-paperclip"></i> {{ pendingFile.name }}
      <button class="text-red-500" @click="pendingFile = null">Хориҷ кардан</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Message, WsEvent } from '~/types/api'

const api = useApi()
const auth = useAuthStore()
const ws = useWebSocket()
const toast = useToast()
const { timeAgo } = useFormat()

const messages = ref<Message[]>([])
const loading = ref(true)
const text = ref('')
const pendingFile = ref<File | null>(null)
const listEl = ref<HTMLElement | null>(null)
const apiBase = api.apiBase

const canSend = computed(() => text.value.trim().length > 0 || !!pendingFile.value)

const scrollDown = () => {
  nextTick(() => {
    if (listEl.value) listEl.value.scrollTop = listEl.value.scrollHeight
  })
}

const loadHistory = async () => {
  loading.value = true
  try {
    messages.value = await api.get<Message[]>('/api/messages', { limit: 200 })
    scrollDown()
  } catch (e: any) {
    toast.error('Хато', e?.message || '')
  } finally {
    loading.value = false
  }
}

const onFile = (ev: Event) => {
  const input = ev.target as HTMLInputElement
  pendingFile.value = input.files?.[0] || null
  input.value = ''
}

const send = async () => {
  if (!canSend.value) return
  try {
    if (pendingFile.value) {
      const fd = new FormData()
      fd.append('message', text.value)
      fd.append('file', pendingFile.value)
      await api.upload('/api/messages', fd)
    } else {
      await api.post('/api/messages', { message: text.value })
    }
    text.value = ''
    pendingFile.value = null
  } catch (e: any) {
    toast.error('Хато', e?.message || '')
  }
}

const onDelete = async (id: number) => {
  await api.delete(`/api/messages/${id}`)
}

// WS
let unsubMsg: (() => void) | null = null
let unsubDel: (() => void) | null = null
let unsubUpd: (() => void) | null = null

onMounted(() => {
  loadHistory()
  unsubMsg = ws.on('chat:message', (ev: WsEvent<Message>) => {
    messages.value.push(ev.data)
    scrollDown()
  })
  unsubDel = ws.on('chat:delete', (ev) => {
    const id = (ev.data as any)?.id
    messages.value = messages.value.filter(m => m.id !== id)
  })
  unsubUpd = ws.on('chat:update', (ev) => {
    const data: any = ev.data
    const m = messages.value.find(x => x.id === data?.id)
    if (m) m.message = data.message
  })
})
onBeforeUnmount(() => { unsubMsg?.(); unsubDel?.(); unsubUpd?.() })
</script>
