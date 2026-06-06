<template>
  <div class="min-h-screen grid place-items-center p-6 bg-page">
    <div class="card max-w-md w-full text-center">
      <div class="text-7xl font-extrabold text-brand-600 mb-2">{{ error?.statusCode || '?' }}</div>
      <div class="text-lg text-ink mb-4">{{ error?.statusMessage || 'Хатои номаълум' }}</div>
      <p class="text-ink-soft text-sm mb-6">{{ error?.message || '' }}</p>
      <button class="btn-primary" @click="handleError">Бозгашт ба саҳифаи асосӣ</button>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ error: { statusCode?: number; statusMessage?: string; message?: string } }>()
const handleError = () => clearError({ redirect: '/' })
</script>
