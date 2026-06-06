<template>
  <Teleport to="body">
    <transition name="modal">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm"
        @click.self="closeOnBackdrop && $emit('update:modelValue', false)"
      >
        <div
          class="bg-white rounded-2xl shadow-card-lg w-full max-h-[90vh] flex flex-col animate-fade-in"
          :class="sizeClass"
        >
          <div v-if="title || $slots.header" class="px-5 py-4 border-b border-border/50 flex items-center justify-between gap-3">
            <slot name="header">
              <div class="font-semibold text-ink">{{ title }}</div>
            </slot>
            <button
              class="text-ink-mute hover:text-ink p-1.5 rounded-lg hover:bg-page transition-colors"
              @click="$emit('update:modelValue', false)"
            >
              <i class="fa-solid fa-xmark"></i>
            </button>
          </div>

          <div class="p-5 overflow-y-auto flex-1">
            <slot />
          </div>

          <div v-if="$slots.footer" class="px-5 py-4 border-t border-border/50 flex items-center justify-end gap-2">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<script setup lang="ts">
const props = defineProps<{
  modelValue: boolean
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
  closeOnBackdrop?: boolean
}>()
defineEmits<{ (e: 'update:modelValue', v: boolean): void }>()

const sizeClass = computed(() => ({
  sm: 'max-w-sm',
  md: 'max-w-lg',
  lg: 'max-w-2xl',
  xl: 'max-w-4xl',
}[props.size || 'md']))

// ESC closes
onMounted(() => {
  const onKey = (e: KeyboardEvent) => { if (e.key === 'Escape' && props.modelValue) emit?.('update:modelValue', false) }
  window.addEventListener('keydown', onKey)
  onBeforeUnmount(() => window.removeEventListener('keydown', onKey))
})

const emit = defineEmits<{ (e: 'update:modelValue', v: boolean): void }>()
</script>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.2s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>
