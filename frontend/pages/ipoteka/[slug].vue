<template>
  <div>
    <NuxtLink to="/ipoteka" class="text-sm text-ink-soft hover:text-brand-700 mb-4 inline-flex items-center gap-1.5">
      <i class="fa-solid fa-arrow-left"></i> Бозгашт
    </NuxtLink>

    <div v-if="loading" class="text-center py-20"><Spinner size="lg" class="text-brand-600" /></div>

    <template v-else-if="bank">
      <div class="card flex flex-col md:flex-row gap-5 items-start mb-6">
        <div class="w-20 h-20 rounded-2xl bg-brand-50 grid place-items-center overflow-hidden shrink-0">
          <img v-if="bank.logo" :src="resolveImg(bank.logo)" :alt="bank.name" class="w-full h-full object-cover" />
          <i v-else class="fa-solid fa-building-columns text-brand-700 text-3xl"></i>
        </div>
        <div class="flex-1 min-w-0">
          <h1 class="text-2xl font-bold text-ink">{{ bank.name }}</h1>
          <p class="text-ink-soft mt-1">{{ bank.description }}</p>
          <div class="mt-3 flex flex-wrap gap-3 text-sm">
            <span v-if="bank.phone"><i class="fa-solid fa-phone text-brand-600"></i> {{ bank.phone }}</span>
            <span v-if="bank.address"><i class="fa-solid fa-location-dot text-brand-600"></i> {{ bank.address }}</span>
            <a v-if="bank.website" :href="bank.website" target="_blank" class="text-brand-700 hover:underline">
              <i class="fa-solid fa-globe"></i> Сомона
            </a>
          </div>
        </div>
      </div>

      <h2 class="text-xl font-bold text-ink mb-3">Шартҳо</h2>
      <EmptyState v-if="conditions.length === 0" icon="fa-percent" title="Шартҳо нестанд" />
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div v-for="c in conditions" :key="c.id" class="card">
          <h3 class="font-bold text-ink mb-2">{{ c.title }}</h3>
          <div class="grid grid-cols-2 gap-2 mb-3">
            <div class="bg-page rounded-lg p-2.5">
              <div class="text-xs text-ink-mute">Ставка</div>
              <div class="font-bold text-brand-700 text-lg">{{ c.interest_rate }}%</div>
            </div>
            <div class="bg-page rounded-lg p-2.5">
              <div class="text-xs text-ink-mute">Муҳлат</div>
              <div class="font-bold text-ink text-lg">{{ c.max_term_years }} сол</div>
            </div>
            <div class="bg-page rounded-lg p-2.5">
              <div class="text-xs text-ink-mute">Аввалин пардохт</div>
              <div class="font-semibold">{{ c.min_down_payment }}%</div>
            </div>
            <div class="bg-page rounded-lg p-2.5">
              <div class="text-xs text-ink-mute">Маблағи макс.</div>
              <div class="font-semibold">{{ c.max_amount }} {{ c.currency }}</div>
            </div>
          </div>
          <div v-if="c.requirements" class="text-sm text-ink-soft mb-2">
            <strong>Талабот:</strong> {{ c.requirements }}
          </div>
          <div v-if="c.documents" class="text-sm text-ink-soft">
            <strong>Ҳуҷҷатҳо:</strong> {{ c.documents }}
          </div>
        </div>
      </div>
    </template>

    <EmptyState v-else icon="fa-circle-question" title="Бонк ёфт нашуд" />
  </div>
</template>

<script setup lang="ts">
import type { Bank } from '~/types/api'

definePageMeta({ layout: 'public' })

const route = useRoute()
const api = useApi()
const config = useRuntimeConfig()

const bank = ref<Bank | null>(null)
const conditions = ref<any[]>([])
const loading = ref(true)

const resolveImg = (s: string) => s.startsWith('http') ? s : `${config.public.apiBase}${s}`

onMounted(async () => {
  try {
    bank.value = await api.get<Bank>(`/api/banks/by-slug/${route.params.slug}`)
    conditions.value = await api.get<any[]>(`/api/banks/${bank.value.id}/conditions`)
  } catch {/* */} finally { loading.value = false }
})
</script>
