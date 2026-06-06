<template>
  <div class="space-y-4 lg:space-y-6">
    <PageHeader title="Дашборд" :description="greeting">
      <template #actions>
        <button class="btn-soft" @click="refresh">
          <i class="fa-solid fa-rotate"></i>
          Нав
        </button>
      </template>
    </PageHeader>

    <!-- Stats grid -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 lg:gap-4">
      <StatCard icon="fa-bullhorn"       label="Лидҳо"     :value="stats.leads"     variant="brand"  />
      <StatCard icon="fa-clipboard-list" label="Заявкаҳо"  :value="stats.requests"  variant="blue"   />
      <StatCard icon="fa-list-check"     label="Вазифаҳо"  :value="stats.tasks"     variant="amber"  />
      <StatCard icon="fa-building"       label="Хонаҳо"    :value="stats.houses"    variant="violet" />
      <StatCard icon="fa-bullhorn"       label="Постҳо"    :value="stats.posts"     variant="rose"   />
      <StatCard icon="fa-sim-card"       label="SIM-кортҳо" :value="stats.sim_cards" variant="brand"  />
      <StatCard icon="fa-users"          label="Кормандон" :value="stats.users"     variant="blue"   />
      <StatCard icon="fa-bell"           label="Оғоҳномаҳои нав" :value="notifStore.unread" variant="amber" />
    </div>

    <!-- Quick links -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-3 lg:gap-4">
      <div class="card">
        <h3 class="font-semibold text-ink mb-3 flex items-center gap-2">
          <i class="fa-solid fa-bolt text-brand-600"></i> Амалҳои зуд
        </h3>
        <div class="space-y-2">
          <NuxtLink to="/lids" class="flex items-center gap-3 px-3 py-2.5 rounded-xl hover:bg-brand-50 transition-colors">
            <i class="fa-solid fa-bullhorn w-5 text-brand-600"></i>
            <span class="flex-1 text-sm font-medium">Канбани лидҳо</span>
            <i class="fa-solid fa-chevron-right text-ink-mute text-xs"></i>
          </NuxtLink>
          <NuxtLink to="/zayavka" class="flex items-center gap-3 px-3 py-2.5 rounded-xl hover:bg-brand-50 transition-colors">
            <i class="fa-solid fa-clipboard-list w-5 text-brand-600"></i>
            <span class="flex-1 text-sm font-medium">Канбани заявкаҳо</span>
            <i class="fa-solid fa-chevron-right text-ink-mute text-xs"></i>
          </NuxtLink>
          <NuxtLink to="/chat" class="flex items-center gap-3 px-3 py-2.5 rounded-xl hover:bg-brand-50 transition-colors">
            <i class="fa-solid fa-comments w-5 text-brand-600"></i>
            <span class="flex-1 text-sm font-medium">Чати ҷамъ</span>
            <i class="fa-solid fa-chevron-right text-ink-mute text-xs"></i>
          </NuxtLink>
          <NuxtLink to="/profile" class="flex items-center gap-3 px-3 py-2.5 rounded-xl hover:bg-brand-50 transition-colors">
            <i class="fa-brands fa-telegram w-5 text-brand-600"></i>
            <span class="flex-1 text-sm font-medium">Пайвасти Telegram</span>
            <i class="fa-solid fa-chevron-right text-ink-mute text-xs"></i>
          </NuxtLink>
        </div>
      </div>

      <div class="card lg:col-span-2">
        <div class="flex items-center justify-between mb-3">
          <h3 class="font-semibold text-ink flex items-center gap-2">
            <i class="fa-solid fa-bell text-brand-600"></i> Оғоҳномаҳои охирин
          </h3>
          <NuxtLink to="/notifications" class="text-xs text-brand-700 hover:underline">Ҳама</NuxtLink>
        </div>
        <div v-if="notifStore.recent.length === 0">
          <EmptyState icon="fa-bell-slash" title="Оғоҳнома нест" description="Ҳангоми расидани лиди нав, заявка ё паём шумо хабардор мешавед" />
        </div>
        <div v-else class="space-y-2">
          <div
            v-for="n in notifStore.recent.slice(0, 5)"
            :key="n.id"
            class="flex items-start gap-3 p-3 rounded-xl hover:bg-page transition-colors"
            :class="!n.is_read && 'bg-brand-50/40'"
          >
            <div class="w-9 h-9 rounded-xl bg-brand-100 grid place-items-center text-brand-700">
              <i class="fa-solid fa-bell text-sm"></i>
            </div>
            <div class="min-w-0 flex-1">
              <div class="text-sm font-semibold text-ink truncate">{{ n.title }}</div>
              <div class="text-xs text-ink-soft truncate">{{ n.body }}</div>
            </div>
            <span class="text-[10px] text-ink-mute shrink-0">{{ timeAgo(n.created_at) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { DashboardStats } from '~/types/api'

const auth = useAuthStore()
const notifStore = useNotificationsStore()
const api = useApi()
const { timeAgo } = useFormat()

const stats = ref<DashboardStats>({
  users: 0, tasks: 0, leads: 0, requests: 0, houses: 0, posts: 0, sim_cards: 0,
})

const greeting = computed(() => {
  const h = new Date().getHours()
  const part = h < 6 ? 'Шаби хайр' : h < 12 ? 'Субҳ ба хайр' : h < 18 ? 'Рӯзатон хуш' : 'Бегоҳатон хуш'
  return `${part}, ${auth.fullName}!`
})

const refresh = async () => {
  try {
    stats.value = await api.get<DashboardStats>('/api/dashboard/stats')
    await notifStore.fetch()
  } catch {/* */}
}

onMounted(refresh)
</script>
