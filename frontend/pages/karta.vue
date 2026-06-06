<template>
  <div class="space-y-4">
    <PageHeader title="Картаи объектҳо" description="Объектҳои амлок дар харита" />

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <div class="card !p-0 overflow-hidden lg:col-span-1 max-h-[70vh] flex flex-col">
        <div class="px-4 py-3 border-b border-border/40">
          <input v-model="search" class="input !py-2" placeholder="Ҷустуҷӯ..." />
        </div>
        <div class="flex-1 overflow-y-auto">
          <div
            v-for="o in filtered"
            :key="o.id"
            class="px-4 py-3 border-b border-border/40 cursor-pointer hover:bg-page/60"
            :class="selected?.id === o.id && 'bg-brand-50/60'"
            @click="selected = o"
          >
            <div class="font-semibold text-sm">{{ o.name }}</div>
            <div class="text-xs text-ink-mute truncate">{{ o.address }}</div>
            <div class="text-xs text-brand-700 mt-1">
              {{ o.lat.toFixed(4) }}, {{ o.lng.toFixed(4) }}
            </div>
          </div>
          <EmptyState v-if="filtered.length === 0" icon="fa-map" title="Объект нест" />
        </div>
      </div>

      <div class="card lg:col-span-2 min-h-[70vh]">
        <div v-if="!selected" class="grid place-items-center h-full text-ink-mute">
          <div class="text-center">
            <i class="fa-solid fa-map-location-dot text-5xl mb-3 text-brand-200"></i>
            <p>Объектро аз рӯйхат интихоб кунед</p>
          </div>
        </div>

        <div v-else>
          <h2 class="text-xl font-bold text-ink mb-1">{{ selected.name }}</h2>
          <p class="text-ink-soft text-sm mb-4">{{ selected.address }}</p>

          <div class="grid grid-cols-2 gap-3 mb-4">
            <div class="bg-page rounded-xl p-3"><div class="text-xs text-ink-mute">Намуди бино</div>
              <div class="font-semibold">{{ selected.construction_type }}</div></div>
            <div class="bg-page rounded-xl p-3"><div class="text-xs text-ink-mute">Ноҳия</div>
              <div class="font-semibold">{{ selected.district }}</div></div>
            <div class="bg-page rounded-xl p-3"><div class="text-xs text-ink-mute">Ҳуҷра</div>
              <div class="font-semibold">{{ selected.rooms }}</div></div>
            <div class="bg-page rounded-xl p-3"><div class="text-xs text-ink-mute">Масоҳат</div>
              <div class="font-semibold">{{ selected.area }} м²</div></div>
            <div class="bg-page rounded-xl p-3"><div class="text-xs text-ink-mute">Ошёна</div>
              <div class="font-semibold">{{ selected.floor }}/{{ selected.total_floors }}</div></div>
            <div class="bg-page rounded-xl p-3"><div class="text-xs text-ink-mute">Нархи умумӣ</div>
              <div class="font-semibold text-brand-700">{{ formatNumber(selected.total_price) }}</div></div>
          </div>

          <a
            :href="`https://www.google.com/maps/search/?api=1&query=${selected.lat},${selected.lng}`"
            target="_blank"
            class="btn-primary w-full justify-center"
          >
            <i class="fa-solid fa-map-pin"></i> Дар Google Maps бинед
          </a>

          <p v-if="selected.description" class="text-sm text-ink-soft mt-4">{{ selected.description }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const api = useApi()
const { formatNumber } = useFormat()

const objects = ref<any[]>([])
const selected = ref<any | null>(null)
const search = ref('')

const filtered = computed(() => {
  const s = search.value.trim().toLowerCase()
  if (!s) return objects.value
  return objects.value.filter(o =>
    [o.name, o.address, o.district].some(v => (v || '').toLowerCase().includes(s)),
  )
})

onMounted(async () => {
  objects.value = await api.get<any[]>('/api/realty/objects')
})
</script>
