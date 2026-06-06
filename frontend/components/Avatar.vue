<template>
  <div
    :class="[
      'rounded-full overflow-hidden grid place-items-center font-semibold shrink-0',
      sizeClass,
      photoUrl ? '' : 'bg-brand-100 text-brand-700',
    ]"
    :title="name"
  >
    <img v-if="photoUrl" :src="photoUrl" :alt="name" class="w-full h-full object-cover" />
    <span v-else>{{ initials(name) }}</span>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  name?: string
  photo?: string
  size?: 'xs' | 'sm' | 'md' | 'lg'
}>()

const { initials } = useFormat()
const config = useRuntimeConfig()

const photoUrl = computed(() => {
  if (!props.photo) return ''
  if (props.photo.startsWith('http')) return props.photo
  return `${config.public.apiBase}${props.photo}`
})

const sizeClass = computed(() => ({
  xs: 'w-7 h-7 text-xs',
  sm: 'w-9 h-9 text-sm',
  md: 'w-11 h-11 text-base',
  lg: 'w-14 h-14 text-lg',
}[props.size || 'md']))
</script>
