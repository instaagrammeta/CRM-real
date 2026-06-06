<template>
  <Modal v-model="open" size="sm" :close-on-backdrop="!busy">
    <template #header>
      <div class="flex items-center gap-2.5">
        <div class="w-9 h-9 rounded-xl grid place-items-center"
             :class="iconWrapClass">
          <i class="fa-solid" :class="iconClass"></i>
        </div>
        <div class="font-semibold text-ink">{{ title || 'Тасдиқ' }}</div>
      </div>
    </template>

    <p class="text-sm text-ink-soft">{{ message }}</p>

    <template #footer>
      <button class="btn-ghost" :disabled="busy" @click="open = false">{{ cancelText || 'Бекор' }}</button>
      <button
        :class="danger ? 'btn-danger' : 'btn-primary'"
        :disabled="busy"
        @click="onConfirm"
      >
        <Spinner v-if="busy" size="xs" class="mr-1" />
        {{ confirmText || 'Тасдиқ' }}
      </button>
    </template>
  </Modal>
</template>

<script setup lang="ts">
const props = defineProps<{
  modelValue: boolean
  title?: string
  message: string
  confirmText?: string
  cancelText?: string
  danger?: boolean
  busy?: boolean
}>()
const emit = defineEmits<{
  (e: 'update:modelValue', v: boolean): void
  (e: 'confirm'): void
}>()

const open = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v),
})

const iconClass = computed(() => props.danger ? 'fa-triangle-exclamation' : 'fa-circle-question')
const iconWrapClass = computed(() => props.danger ? 'bg-red-100 text-red-600' : 'bg-brand-100 text-brand-700')

const onConfirm = () => emit('confirm')
</script>
