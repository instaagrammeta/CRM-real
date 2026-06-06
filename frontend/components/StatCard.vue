<template>
  <div class="card flex items-center gap-4">
    <div
      class="w-12 h-12 rounded-2xl grid place-items-center text-xl shrink-0"
      :class="bgClass"
    >
      <i class="fa-solid" :class="icon"></i>
    </div>
    <div class="min-w-0 flex-1">
      <div class="text-xs text-ink-mute uppercase tracking-wider mb-0.5">{{ label }}</div>
      <div class="text-2xl font-bold text-ink leading-tight">{{ value }}</div>
      <div v-if="hint" class="text-[11px] text-ink-mute mt-0.5">{{ hint }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  icon: string
  label: string
  value: string | number
  hint?: string
  variant?: 'brand' | 'amber' | 'blue' | 'rose' | 'violet'
}>()

const bgClass = computed(() => ({
  brand:  'bg-brand-100 text-brand-700',
  amber:  'bg-amber-100 text-amber-700',
  blue:   'bg-blue-100 text-blue-700',
  rose:   'bg-rose-100 text-rose-700',
  violet: 'bg-violet-100 text-violet-700',
}[props.variant || 'brand']))
</script>
