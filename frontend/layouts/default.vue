<template>
  <div class="min-h-screen flex bg-page p-3 lg:p-4 gap-3 lg:gap-4">
    <!-- Sidebar -->
    <aside
      :class="[
        'shrink-0 bg-white rounded-3xl shadow-card flex flex-col overflow-hidden transition-all duration-200',
        sidebarCollapsed ? 'w-[72px]' : 'w-[260px]',
        sidebarOpen ? 'fixed inset-y-3 left-3 z-40 lg:static' : 'hidden lg:flex',
      ]"
    >
      <!-- Brand -->
      <div class="px-5 py-5 border-b border-border/40 flex items-center justify-between">
        <NuxtLink to="/" class="flex items-center gap-2.5">
          <div class="w-9 h-9 rounded-xl bg-brand-100 grid place-items-center">
            <i class="fa-solid fa-house-chimney text-brand-700 text-base"></i>
          </div>
          <span v-if="!sidebarCollapsed" class="font-bold text-ink whitespace-nowrap">
            Real Estate <span class="text-brand-600">CRM</span>
          </span>
        </NuxtLink>
        <button
          class="lg:hidden text-ink-mute hover:text-ink"
          @click="sidebarOpen = false"
        >
          <i class="fa-solid fa-xmark"></i>
        </button>
      </div>

      <!-- Nav -->
      <nav class="flex-1 px-3 py-4 space-y-0.5 overflow-y-auto">
        <NavItem to="/"          icon="fa-chart-line"      label="Дашборд"     :collapsed="sidebarCollapsed" />
        <NavItem to="/lids"      icon="fa-bullhorn"        label="Лидҳо"       :collapsed="sidebarCollapsed" />
        <NavItem to="/zayavka"   icon="fa-clipboard-list"  label="Заявкаҳо"    :collapsed="sidebarCollapsed" />
        <NavItem to="/zadacha"   icon="fa-list-check"      label="Вазифаҳо"    :collapsed="sidebarCollapsed" />
        <NavItem to="/houses"    icon="fa-building"        label="Хонаҳо"      :collapsed="sidebarCollapsed" />
        <NavItem to="/karta"     icon="fa-map-location-dot" label="Карта"      :collapsed="sidebarCollapsed" />
        <NavItem to="/obiekt"    icon="fa-table-cells"     label="Шахматка"    :collapsed="sidebarCollapsed" />
        <NavItem to="/posts"     icon="fa-bullhorn"        label="Постҳо"      :collapsed="sidebarCollapsed" />
        <NavItem to="/baza"      icon="fa-folder-open"     label="База"        :collapsed="sidebarCollapsed" />
        <NavItem to="/chat"      icon="fa-comments"        label="Чат"         :collapsed="sidebarCollapsed" :badge="onlineBadge" />
        <NavItem to="/ipoteka"   icon="fa-percent"         label="Ипотека"     :collapsed="sidebarCollapsed" />
        <NavItem to="/rasrochka" icon="fa-money-bill"      label="Рассрочка"   :collapsed="sidebarCollapsed" />

        <template v-if="auth.isAdmin">
          <div v-if="!sidebarCollapsed" class="mt-5 mb-1 px-3 text-[11px] uppercase font-semibold text-ink-mute tracking-wider">
            Админпанел
          </div>
          <NavItem to="/admin/users"     icon="fa-users-gear" label="Кормандон"  :collapsed="sidebarCollapsed" />
          <NavItem to="/sim-cards"       icon="fa-sim-card"   label="SIM-кортҳо" :collapsed="sidebarCollapsed" />
          <NavItem to="/admin/ipoteka"   icon="fa-building-columns" label="Ипотекаи админ" :collapsed="sidebarCollapsed" />
          <NavItem to="/admin/rasrochka" icon="fa-handshake" label="Рассрочкаи админ"   :collapsed="sidebarCollapsed" />
        </template>
      </nav>

      <!-- User card -->
      <div class="border-t border-border/40 p-3">
        <NuxtLink
          to="/profile"
          class="flex items-center gap-3 px-2 py-2 rounded-xl hover:bg-brand-50 transition-colors"
        >
          <Avatar :name="auth.fullName" :photo="auth.avatar" size="sm" />
          <div v-if="!sidebarCollapsed" class="min-w-0 flex-1">
            <div class="text-sm font-semibold text-ink truncate">{{ auth.fullName || '—' }}</div>
            <div class="text-xs text-ink-mute truncate">{{ auth.user?.category || auth.role }}</div>
          </div>
        </NuxtLink>
      </div>
    </aside>

    <!-- Mobile overlay -->
    <div
      v-if="sidebarOpen"
      class="fixed inset-0 bg-black/30 z-30 lg:hidden"
      @click="sidebarOpen = false"
    />

    <!-- Main -->
    <main class="flex-1 min-w-0 flex flex-col gap-3 lg:gap-4">
      <!-- Top bar -->
      <header class="bg-white rounded-2xl shadow-card px-4 lg:px-6 py-3 flex items-center justify-between gap-3">
        <div class="flex items-center gap-3 min-w-0">
          <button
            class="lg:hidden btn-ghost !p-2"
            @click="sidebarOpen = true"
          >
            <i class="fa-solid fa-bars"></i>
          </button>
          <button
            class="hidden lg:inline-flex btn-ghost !p-2"
            @click="sidebarCollapsed = !sidebarCollapsed"
            :title="sidebarCollapsed ? 'Кушодан' : 'Пӯшидан'"
          >
            <i class="fa-solid" :class="sidebarCollapsed ? 'fa-angles-right' : 'fa-angles-left'"></i>
          </button>
          <div class="min-w-0">
            <div class="text-base font-semibold text-ink truncate">{{ pageTitle }}</div>
            <div class="text-xs text-ink-mute hidden sm:block">{{ today }}</div>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <div
            class="hidden sm:flex items-center gap-1.5 text-xs"
            :title="ws.connected.value ? 'Real-time фаъол' : 'Real-time қатъ'"
          >
            <span
              class="w-2 h-2 rounded-full"
              :class="ws.connected.value ? 'bg-brand-500 animate-pulse-soft' : 'bg-gray-300'"
            />
            <span class="text-ink-mute">{{ ws.connected.value ? 'Онлайн' : 'Офлайн' }}</span>
          </div>

          <NotificationsDropdown />

          <div class="relative">
            <button class="btn-ghost !p-2" @click="userMenu = !userMenu">
              <Avatar :name="auth.fullName" :photo="auth.avatar" size="xs" />
            </button>
            <div
              v-if="userMenu"
              v-click-outside="() => (userMenu = false)"
              class="absolute right-0 top-full mt-2 w-56 bg-white rounded-2xl shadow-card-lg p-2 z-50 animate-fade-in border border-border/40"
            >
              <div class="px-3 py-2.5 border-b border-border/40 mb-1">
                <div class="text-sm font-semibold text-ink">{{ auth.fullName }}</div>
                <div class="text-xs text-ink-mute">{{ auth.user?.login }}</div>
              </div>
              <NuxtLink
                to="/profile"
                class="flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm text-ink-soft hover:bg-brand-50 hover:text-brand-700"
                @click="userMenu = false"
              >
                <i class="fa-solid fa-user w-4"></i> Профил
              </NuxtLink>
              <NuxtLink
                to="/notifications"
                class="flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm text-ink-soft hover:bg-brand-50 hover:text-brand-700"
                @click="userMenu = false"
              >
                <i class="fa-solid fa-bell w-4"></i> Оғоҳномаҳо
              </NuxtLink>
              <button
                class="w-full flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm text-red-600 hover:bg-red-50"
                @click="logout"
              >
                <i class="fa-solid fa-right-from-bracket w-4"></i> Баромад
              </button>
            </div>
          </div>
        </div>
      </header>

      <!-- Page content -->
      <div class="flex-1 min-h-0 animate-fade-in">
        <slot />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
const auth = useAuthStore()
const route = useRoute()
const ws = useWebSocket()

const sidebarCollapsed = ref(false)
const sidebarOpen = ref(false)
const userMenu = ref(false)

watch(() => route.fullPath, () => { sidebarOpen.value = false; userMenu.value = false })

const pageTitle = computed(() => {
  const map: Record<string, string> = {
    '/':            'Дашборд',
    '/lids':        'Лидҳо',
    '/zayavka':     'Заявкаҳо',
    '/zadacha':     'Вазифаҳо',
    '/houses':      'Хонаҳо',
    '/karta':       'Карта',
    '/obiekt':      'Шахматка',
    '/posts':       'Постҳо',
    '/baza':        'База',
    '/chat':        'Чат',
    '/ipoteka':     'Ипотека',
    '/rasrochka':   'Рассрочка',
    '/profile':     'Профил',
    '/notifications': 'Оғоҳномаҳо',
    '/sim-cards':   'SIM-кортҳо',
    '/admin/users':     'Идоракунии кормандон',
    '/admin/ipoteka':   'Идоракунии бонкҳо',
    '/admin/rasrochka': 'Идоракунии рассрочка',
  }
  return map[route.path] || 'Real Estate CRM'
})

const today = computed(() => new Date().toLocaleDateString('ru-RU', {
  day: '2-digit', month: 'long', year: 'numeric',
}))

const onlineBadge = computed(() => '')

const logout = async () => {
  userMenu.value = false
  await auth.logout()
}

// click-outside директива (минимальная инлайн-реализация)
const vClickOutside = {
  mounted(el: HTMLElement, binding: any) {
    el._clickOutside = (e: Event) => {
      if (!el.contains(e.target as Node)) binding.value(e)
    }
    setTimeout(() => document.addEventListener('click', el._clickOutside), 0)
  },
  unmounted(el: any) {
    document.removeEventListener('click', el._clickOutside)
  },
}
</script>
