<template>
  <div class="flex items-center justify-between gap-4 mb-4 lg:mb-6">
    <div class="min-w-0">
      <h1 class="text-xl lg:text-2xl font-bold text-ink truncate">{{ title }}</h1>
      <p v-if="description" class="text-sm text-ink-soft mt-0.5">{{ description }}</p>
    </div>
    <div class="flex items-center gap-2 shrink-0">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ title: string; description?: string }>()
</script>
