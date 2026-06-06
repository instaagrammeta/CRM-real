<template>
  <div>
    <div class="text-center mb-8 lg:mb-10">
      <h1 class="text-3xl lg:text-4xl font-bold text-ink mb-2">Рассрочка</h1>
      <p class="text-ink-soft">Объектҳо ва шартҳои онҳо</p>
    </div>

    <div v-if="loading" class="text-center py-20"><Spinner size="lg" class="text-brand-600" /></div>
    <EmptyState v-else-if="objects.length === 0" icon="fa-handshake" title="Объектҳо нестанд" />

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 lg:gap-6">
      <NuxtLink v-for="o in objects" :key="o.id" :to="`/rasrochka/${o.slug}`" class="card-hover group">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-14 h-14 rounded-2xl bg-brand-50 grid place-items-center overflow-hidden shrink-0">
            <img v-if="o.logo" :src="resolveImg(o.logo)" :alt="o.name" class="w-full h-full object-cover" />
            <i v-else class="fa-solid fa-handshake text-brand-700 text-xl"></i>
          </div>
          <div class="min-w-0 flex-1">
            <h3 class="font-bold text-ink truncate group-hover:text-brand-700">{{ o.name }}</h3>
            <p v-if="o.developer" class="text-xs text-ink-mute"><i class="fa-solid fa-helmet-safety"></i> {{ o.developer }}</p>
          </div>
        </div>
        <p v-if="o.description" class="text-sm text-ink-soft line-clamp-3 mb-3">{{ o.description }}</p>
        <div class="text-brand-700 text-sm font-semibold flex items-center gap-1.5">
          Шартҳоро бинед <i class="fa-solid fa-arrow-right text-xs"></i>
        </div>
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { InstallmentObject } from '~/types/api'
definePageMeta({ layout: 'public' })

const api = useApi()
const config = useRuntimeConfig()
const objects = ref<InstallmentObject[]>([])
const loading = ref(true)

const resolveImg = (s: string) => s.startsWith('http') ? s : `${config.public.apiBase}${s}`

onMounted(async () => {
  try { objects.value = await api.get<InstallmentObject[]>('/api/installment-objects', { active: 'true' }) }
  finally { loading.value = false }
})
</script>
