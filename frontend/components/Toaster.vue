<template>
  <Teleport to="body">
    <div class="fixed top-4 right-4 z-[60] space-y-2 max-w-sm w-full pointer-events-none">
      <transition-group name="toast" tag="div" class="space-y-2">
        <div
          v-for="t in toaster.toasts.value"
          :key="t.id"
          class="pointer-events-auto bg-white rounded-2xl shadow-card-lg border-l-4 px-4 py-3 flex items-start gap-3"
          :class="borderColor(t.type)"
        >
          <i :class="['fa-solid mt-0.5 text-base', iconClass(t.type), iconColor(t.type)]"></i>
          <div class="min-w-0 flex-1">
            <div class="font-semibold text-ink text-sm">{{ t.title }}</div>
            <div v-if="t.message" class="text-xs text-ink-soft mt-0.5 break-words">{{ t.message }}</div>
          </div>
          <button class="text-ink-mute hover:text-ink shrink-0" @click="toaster.remove(t.id)">
            <i class="fa-solid fa-xmark text-xs"></i>
          </button>
        </div>
      </transition-group>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import type { ToastType } from '~/composables/useToast'
const toaster = useToast()

const iconClass = (t: ToastType) => ({
  success: 'fa-circle-check',
  error:   'fa-circle-xmark',
  info:    'fa-circle-info',
  warning: 'fa-triangle-exclamation',
}[t])
const iconColor = (t: ToastType) => ({
  success: 'text-brand-600',
  error:   'text-red-600',
  info:    'text-blue-600',
  warning: 'text-amber-600',
}[t])
const borderColor = (t: ToastType) => ({
  success: 'border-brand-500',
  error:   'border-red-500',
  info:    'border-blue-500',
  warning: 'border-amber-500',
}[t])
</script>

<style scoped>
.toast-enter-active, .toast-leave-active { transition: all 0.25s ease; }
.toast-enter-from { opacity: 0; transform: translateX(40px); }
.toast-leave-to   { opacity: 0; transform: translateX(40px); }
</style>
