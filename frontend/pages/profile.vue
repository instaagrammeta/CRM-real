<template>
  <div class="space-y-4 lg:space-y-6">
    <PageHeader title="Профил" description="Маълумоти ҳисоб ва Telegram" />

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Profile -->
      <div class="card text-center lg:col-span-1">
        <div class="flex justify-center mb-3">
          <Avatar :name="auth.fullName" :photo="auth.avatar" size="lg" />
        </div>
        <h2 class="text-lg font-bold text-ink">{{ auth.fullName }}</h2>
        <div class="text-sm text-ink-soft">{{ auth.user?.login }}</div>
        <span class="badge-brand mt-2 inline-flex">{{ auth.role }}</span>
        <div v-if="auth.user?.category" class="text-xs text-ink-mute mt-1">{{ auth.user.category }}</div>

        <button class="btn-outline w-full mt-5" @click="auth.logout()">
          <i class="fa-solid fa-right-from-bracket"></i> Баромад
        </button>
      </div>

      <!-- Telegram -->
      <div class="card lg:col-span-2 space-y-4">
        <div class="flex items-center gap-3">
          <div class="w-12 h-12 rounded-2xl bg-blue-100 grid place-items-center text-blue-600 text-xl">
            <i class="fa-brands fa-telegram"></i>
          </div>
          <div>
            <h3 class="font-semibold text-ink">Telegram оғоҳномаҳо</h3>
            <p class="text-xs text-ink-soft">
              Ҳама лидҳои нав, заявкаҳо ва паёмҳо ба Telegram-и шумо мерасанд.
            </p>
          </div>
        </div>

        <div v-if="loading" class="text-center py-6"><Spinner class="text-brand-600" /></div>

        <template v-else>
          <div v-if="subs.length > 0" class="space-y-2">
            <div class="text-xs text-ink-mute uppercase font-semibold tracking-wider">Пайвастҳои фаъол</div>
            <div
              v-for="s in subs"
              :key="s.id"
              class="flex items-center gap-3 p-3 rounded-xl bg-page"
            >
              <div class="w-10 h-10 rounded-xl bg-blue-100 grid place-items-center text-blue-600">
                <i class="fa-brands fa-telegram"></i>
              </div>
              <div class="min-w-0 flex-1">
                <div class="text-sm font-semibold">
                  {{ s.first_name }} {{ s.last_name }}
                  <span v-if="s.username" class="text-ink-mute font-normal">@{{ s.username }}</span>
                </div>
                <div class="text-xs text-ink-mute">chat_id: {{ s.chat_id }}</div>
              </div>
              <span v-if="s.is_active" class="badge-brand">Фаъол</span>
              <span v-else class="badge-gray">Ғайри фаъол</span>
              <button class="text-ink-mute hover:text-red-600" @click="unlink(s.id)">
                <i class="fa-solid fa-trash text-sm"></i>
              </button>
            </div>
          </div>

          <div v-if="linkData" class="rounded-xl border-2 border-dashed border-brand-200 bg-brand-50/40 p-4">
            <div class="text-sm text-ink-soft mb-2">Барои pуйвастан:</div>
            <ol class="text-sm space-y-1.5 list-decimal list-inside text-ink mb-3">
              <li>Telegram-ро кушоед ва ботро ёбед</li>
              <li>Командаи зеринро равон кунед:</li>
            </ol>
            <div class="bg-white border border-brand-200 rounded-xl p-3 flex items-center justify-between gap-2 mb-3">
              <code class="text-sm text-brand-700 font-mono">{{ linkData.command }}</code>
              <button class="btn-ghost !p-2" @click="copy(linkData.command)" title="Нусха">
                <i class="fa-regular fa-copy"></i>
              </button>
            </div>
            <a
              v-if="linkData.link_url"
              :href="linkData.link_url"
              target="_blank"
              class="btn-primary w-full justify-center"
            >
              <i class="fa-brands fa-telegram"></i> Кушодани ботро дар Telegram
            </a>
          </div>

          <button v-else class="btn-soft" @click="generate" :disabled="busy">
            <Spinner v-if="busy" size="sm" />
            <i v-else class="fa-solid fa-link"></i>
            Эҷоди токени pуйвастан
          </button>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { TelegramSubscriber, TelegramLinkResponse } from '~/types/api'

const api = useApi()
const auth = useAuthStore()
const toast = useToast()

const subs = ref<TelegramSubscriber[]>([])
const linkData = ref<TelegramLinkResponse | null>(null)
const loading = ref(true)
const busy = ref(false)

const load = async () => {
  loading.value = true
  try {
    subs.value = await api.get<TelegramSubscriber[]>('/api/me/telegram')
  } catch {/* */} finally { loading.value = false }
}

const generate = async () => {
  busy.value = true
  try {
    linkData.value = await api.post<TelegramLinkResponse>('/api/me/telegram/link')
    toast.success('Токен сохта шуд')
  } catch (e: any) {
    toast.error('Хато', e?.data?.error || e?.message || '')
  } finally { busy.value = false }
}

const unlink = async (id: number) => {
  await api.delete(`/api/me/telegram/${id}`)
  toast.success('Узв партофта шуд')
  await load()
}

const copy = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    toast.success('Нусха шуд')
  } catch {/* */}
}

onMounted(load)
</script>
