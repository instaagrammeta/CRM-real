<template>
  <div>
    <div class="text-center mb-8 lg:mb-10">
      <h1 class="text-3xl lg:text-4xl font-bold text-ink mb-2">Ипотека дар Тоҷикистон</h1>
      <p class="text-ink-soft">Бонкҳои шарик ва шартҳои онҳо</p>
    </div>

    <div v-if="loading" class="text-center py-20"><Spinner size="lg" class="text-brand-600" /></div>
    <EmptyState v-else-if="banks.length === 0" icon="fa-building-columns" title="Бонкҳо илова нашудаанд" />

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 lg:gap-6">
      <NuxtLink
        v-for="b in banks"
        :key="b.id"
        :to="`/ipoteka/${b.slug}`"
        class="card-hover group"
      >
        <div class="flex items-center gap-3 mb-3">
          <div class="w-14 h-14 rounded-2xl bg-brand-50 grid place-items-center overflow-hidden shrink-0">
            <img v-if="b.logo" :src="resolveImg(b.logo)" :alt="b.name" class="w-full h-full object-cover" />
            <i v-else class="fa-solid fa-building-columns text-brand-700 text-xl"></i>
          </div>
          <div class="min-w-0 flex-1">
            <h3 class="font-bold text-ink truncate group-hover:text-brand-700 transition-colors">{{ b.name }}</h3>
            <p v-if="b.phone" class="text-xs text-ink-mute"><i class="fa-solid fa-phone"></i> {{ b.phone }}</p>
          </div>
        </div>
        <p v-if="b.description" class="text-sm text-ink-soft line-clamp-3 mb-3">{{ b.description }}</p>
        <div class="text-brand-700 text-sm font-semibold flex items-center gap-1.5">
          Шартҳоро бинед <i class="fa-solid fa-arrow-right text-xs group-hover:translate-x-0.5 transition-transform"></i>
        </div>
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Bank } from '~/types/api'

definePageMeta({ layout: 'public' })

const api = useApi()
const config = useRuntimeConfig()

const banks = ref<Bank[]>([])
const loading = ref(true)

const resolveImg = (s: string) => s.startsWith('http') ? s : `${config.public.apiBase}${s}`

onMounted(async () => {
  try { banks.value = await api.get<Bank[]>('/api/banks', { active: 'true' }) }
  finally { loading.value = false }
})
</script>
