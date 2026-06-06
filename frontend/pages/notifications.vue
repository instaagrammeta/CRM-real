<template>
  <div class="space-y-4">
    <PageHeader title="Оғоҳномаҳо" :description="`${store.unread} нав, ҳамагӣ ${store.items.length}`">
      <template #actions>
        <button v-if="store.unread > 0" class="btn-soft" @click="store.markAllRead()">
          <i class="fa-solid fa-check-double"></i> Ҳамаро хонда шудагӣ кардан
        </button>
        <button class="btn-ghost" @click="store.fetch()">
          <i class="fa-solid fa-rotate"></i>
        </button>
      </template>
    </PageHeader>

    <div class="card !p-0">
      <div v-if="store.loading" class="py-12 text-center">
        <Spinner class="text-brand-600" size="lg" />
      </div>

      <EmptyState
        v-else-if="store.items.length === 0"
        icon="fa-bell-slash"
        title="Ҳоло оғоҳнома нест"
        description="Оғоҳномаҳои real-time дар ин ҷо пайдо мешаванд"
      />

      <div v-else>
        <div
          v-for="n in store.items"
          :key="n.id"
          class="px-5 py-4 border-b border-border/40 last:border-0 flex items-start gap-3 hover:bg-page transition-colors cursor-pointer"
          :class="!n.is_read && 'bg-brand-50/30'"
          @click="onClick(n)"
        >
          <div class="w-10 h-10 rounded-xl bg-brand-100 grid place-items-center text-brand-700 shrink-0">
            <i class="fa-solid" :class="iconFor(n.type)"></i>
          </div>
          <div class="min-w-0 flex-1">
            <div class="flex items-baseline justify-between gap-2 mb-0.5">
              <div class="font-semibold text-ink truncate">{{ n.title }}</div>
              <span class="text-[10px] text-ink-mute shrink-0">{{ timeAgo(n.created_at) }}</span>
            </div>
            <div class="text-sm text-ink-soft">{{ n.body }}</div>
            <div v-if="n.entity_type" class="mt-1.5">
              <span class="badge-gray">{{ n.entity_type }}</span>
            </div>
          </div>
          <div class="flex flex-col items-end gap-2 shrink-0">
            <span v-if="!n.is_read" class="w-2 h-2 rounded-full bg-brand-500"></span>
            <button class="text-ink-mute hover:text-red-600" @click.stop="store.remove(n.id)">
              <i class="fa-solid fa-trash text-xs"></i>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Notification } from '~/types/api'
const store = useNotificationsStore()
const { timeAgo } = useFormat()

const iconFor = (type: string): string => ({
  task: 'fa-list-check',
  lead: 'fa-bullhorn',
  request: 'fa-clipboard-list',
  chat: 'fa-comment',
  tariff: 'fa-sim-card',
}[type] || 'fa-bell')

const onClick = async (n: Notification) => {
  if (!n.is_read) await store.markRead(n.id)
  if (n.link) await navigateTo(n.link)
}

onMounted(() => { if (!store.loaded) store.fetch() })
</script>
