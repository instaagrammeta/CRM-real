<template>
  <div class="flex flex-col items-center justify-center py-12 px-4 text-center">
    <div class="w-16 h-16 rounded-2xl bg-brand-50 grid place-items-center text-brand-600 text-2xl mb-4">
      <i class="fa-solid" :class="icon || 'fa-inbox'"></i>
    </div>
    <div class="text-ink font-semibold mb-1">{{ title || 'Маълумот нест' }}</div>
    <div v-if="description" class="text-sm text-ink-soft max-w-xs">{{ description }}</div>
    <div v-if="$slots.action" class="mt-5">
      <slot name="action" />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ icon?: string; title?: string; description?: string }>()
</script>
