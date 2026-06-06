<template>
  <div>
    <NuxtLink to="/rasrochka" class="text-sm text-ink-soft hover:text-brand-700 mb-4 inline-flex items-center gap-1.5">
      <i class="fa-solid fa-arrow-left"></i> Бозгашт
    </NuxtLink>

    <div v-if="loading" class="text-center py-20"><Spinner size="lg" class="text-brand-600" /></div>

    <template v-else-if="obj">
      <div class="card flex flex-col md:flex-row gap-5 items-start mb-6">
        <div class="w-20 h-20 rounded-2xl bg-brand-50 grid place-items-center overflow-hidden shrink-0">
          <img v-if="obj.logo" :src="resolveImg(obj.logo)" :alt="obj.name" class="w-full h-full object-cover" />
          <i v-else class="fa-solid fa-handshake text-brand-700 text-3xl"></i>
        </div>
        <div class="flex-1 min-w-0">
          <h1 class="text-2xl font-bold text-ink">{{ obj.name }}</h1>
          <p class="text-ink-soft mt-1">{{ obj.description }}</p>
          <div class="mt-3 flex flex-wrap gap-3 text-sm">
            <span v-if="obj.phone"><i class="fa-solid fa-phone text-brand-600"></i> {{ obj.phone }}</span>
            <span v-if="obj.developer"><i class="fa-solid fa-helmet-safety text-brand-600"></i> {{ obj.developer }}</span>
            <span v-if="obj.address"><i class="fa-solid fa-location-dot text-brand-600"></i> {{ obj.address }}</span>
          </div>
        </div>
      </div>

      <h2 class="text-xl font-bold text-ink mb-3">Шартҳо</h2>
      <EmptyState v-if="conditions.length === 0" icon="fa-money-bill" title="Шартҳо нестанд" />
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div v-for="c in conditions" :key="c.id" class="card">
          <h3 class="font-bold text-ink mb-2">{{ c.title }}</h3>
          <div class="grid grid-cols-2 gap-2 mb-3">
            <div class="bg-page rounded-lg p-2.5">
              <div class="text-xs text-ink-mute">Аввалин пардохт</div>
              <div class="font-bold text-brand-700 text-lg">{{ c.min_down_payment }}%</div>
            </div>
            <div class="bg-page rounded-lg p-2.5">
              <div class="text-xs text-ink-mute">Муҳлат</div>
              <div class="font-bold text-ink text-lg">{{ c.max_term_months }} моҳ</div>
            </div>
          </div>
          <div v-if="c.additional_info" class="text-sm text-ink-soft">{{ c.additional_info }}</div>
        </div>
      </div>
    </template>

    <EmptyState v-else icon="fa-circle-question" title="Объект ёфт нашуд" />
  </div>
</template>

<script setup lang="ts">
import type { InstallmentObject } from '~/types/api'
definePageMeta({ layout: 'public' })

const route = useRoute()
const api = useApi()
const config = useRuntimeConfig()
const obj = ref<InstallmentObject | null>(null)
const conditions = ref<any[]>([])
const loading = ref(true)

const resolveImg = (s: string) => s.startsWith('http') ? s : `${config.public.apiBase}${s}`

onMounted(async () => {
  try {
    obj.value = await api.get<InstallmentObject>(`/api/installment-objects/by-slug/${route.params.slug}`)
    conditions.value = await api.get<any[]>(`/api/installment-objects/${obj.value.id}/conditions`)
  } catch {/* */} finally { loading.value = false }
})
</script>
