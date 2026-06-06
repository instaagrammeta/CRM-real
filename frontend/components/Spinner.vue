<template>
  <span
    class="inline-block animate-spin rounded-full border-2 border-current border-t-transparent"
    :class="sizeClass"
    role="status"
    aria-label="loading"
  />
</template>

<script setup lang="ts">
const props = defineProps<{ size?: 'xs' | 'sm' | 'md' | 'lg' }>()
const sizeClass = computed(() => ({
  xs: 'w-3 h-3',
  sm: 'w-4 h-4',
  md: 'w-6 h-6',
  lg: 'w-8 h-8',
}[props.size || 'sm']))
</script>
