<template>
  <NuxtLink
    :to="to"
    class="nav-link group"
    active-class="nav-link-active"
    :exact-active-class="''"
  >
    <i class="fa-solid w-5 text-center text-base shrink-0" :class="icon"></i>
    <span v-if="!collapsed" class="flex-1 truncate">{{ label }}</span>
    <span
      v-if="!collapsed && badge"
      class="ml-auto bg-red-500 text-white text-[10px] font-bold px-1.5 py-0.5 rounded-full min-w-[18px] text-center"
    >
      {{ badge }}
    </span>
  </NuxtLink>
</template>

<script setup lang="ts">
defineProps<{
  to: string
  icon: string
  label: string
  collapsed?: boolean
  badge?: string | number
}>()
</script>
