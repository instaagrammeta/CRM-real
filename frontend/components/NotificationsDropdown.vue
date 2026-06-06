<template>
  <div class="relative">
    <button
      class="btn-ghost !p-2 relative"
      @click="open = !open"
      :title="`${store.unread} нав`"
    >
      <i class="fa-solid fa-bell"></i>
      <span
        v-if="store.unread > 0"
        class="absolute -top-0.5 -right-0.5 min-w-[18px] h-[18px] px-1 bg-red-500 text-white text-[10px] font-bold rounded-full grid place-items-center"
      >
        {{ store.unread > 99 ? '99+' : store.unread }}
      </span>
    </button>

    <Transition name="fade">
      <div
        v-if="open"
        ref="dropEl"
        class="absolute right-0 top-full mt-2 w-[360px] max-h-[80vh] bg-white rounded-2xl shadow-card-lg border border-border/40 z-50 flex flex-col overflow-hidden animate-fade-in"
      >
        <div class="px-4 py-3 border-b border-border/40 flex items-center justify-between gap-2">
          <div class="font-semibold text-ink">Оғоҳномаҳо</div>
          <button
            v-if="store.unread > 0"
            class="text-xs text-brand-700 hover:underline"
            @click="store.markAllRead()"
          >
            Ҳамаро хонда шудагӣ кардан
          </button>
        </div>

        <div class="overflow-y-auto flex-1">
          <div v-if="store.items.length === 0" class="py-12 text-center text-ink-mute text-sm">
            <i class="fa-solid fa-bell-slash text-3xl mb-3 opacity-50"></i>
            <div>Оғоҳнома нест</div>
          </div>

          <button
            v-for="n in store.items.slice(0, 20)"
            :key="n.id"
            class="w-full text-left px-4 py-3 border-b border-border/40 hover:bg-brand-50/50 transition-colors flex items-start gap-3"
            :class="!n.is_read && 'bg-brand-50/30'"
            @click="onClick(n)"
          >
            <span class="w-2 h-2 rounded-full mt-2 shrink-0"
                  :class="n.is_read ? 'bg-transparent' : 'bg-brand-500'"></span>
            <div class="min-w-0 flex-1">
              <div class="text-sm font-semibold text-ink truncate">{{ n.title }}</div>
              <div class="text-xs text-ink-soft mt-0.5 line-clamp-2">{{ n.body }}</div>
              <div class="text-[10px] text-ink-mute mt-1">{{ timeAgo(n.created_at) }}</div>
            </div>
          </button>
        </div>

        <NuxtLink
          to="/notifications"
          class="block text-center text-sm text-brand-700 hover:bg-brand-50 py-3 border-t border-border/40 font-medium"
          @click="open = false"
        >
          Ҳамаи оғоҳномаҳо
        </NuxtLink>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import type { Notification } from '~/types/api'

const store = useNotificationsStore()
const { timeAgo } = useFormat()
const open = ref(false)
const dropEl = ref<HTMLElement | null>(null)

onClickOutside(dropEl, () => { open.value = false })

const onClick = async (n: Notification) => {
  if (!n.is_read) await store.markRead(n.id)
  if (n.link) {
    open.value = false
    await navigateTo(n.link)
  }
}
</script>
