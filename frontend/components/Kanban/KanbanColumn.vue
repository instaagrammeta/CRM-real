<template>
  <div class="kanban-col shrink-0 w-[300px] bg-white/60 backdrop-blur rounded-2xl flex flex-col max-h-[calc(100vh-180px)]">
    <!-- Header -->
    <div
      class="kanban-col-handle px-4 py-3 rounded-t-2xl flex items-center justify-between gap-2 cursor-grab active:cursor-grabbing"
      :style="{ background: column.color ? column.color + '14' : 'var(--brand-50)' }"
    >
      <div class="flex items-center gap-2 min-w-0">
        <span class="w-2.5 h-2.5 rounded-full" :style="{ background: column.color || '#1f7a4d' }" />
        <div class="font-semibold text-ink truncate">{{ column.title }}</div>
        <span class="badge-gray !py-0.5">{{ leads.length }}</span>
      </div>
      <div class="flex items-center gap-1 shrink-0">
        <button class="text-ink-mute hover:text-brand-700 p-1" @click="$emit('column-edit', column)" title="Вироиш">
          <i class="fa-solid fa-pen text-xs"></i>
        </button>
        <button class="text-ink-mute hover:text-red-600 p-1" @click="$emit('column-delete', column)" title="Нест кардан">
          <i class="fa-solid fa-trash text-xs"></i>
        </button>
      </div>
    </div>

    <!-- Cards -->
    <div class="flex-1 overflow-y-auto p-2 space-y-2">
      <draggable
        :model-value="leads"
        :item-key="(l: any) => l.id"
        group="leads"
        class="space-y-2 min-h-[20px]"
        :animation="180"
        @end="onEnd"
        @update:model-value="onListUpdate"
      >
        <template #item="{ element: lead }">
          <KanbanCard
            :lead="lead"
            @click="$emit('lead-edit', lead)"
            @delete="$emit('lead-delete', lead)"
          />
        </template>
      </draggable>
    </div>

    <!-- Add card -->
    <button
      class="m-2 mt-0 px-3 py-2 rounded-xl text-sm text-brand-700 hover:bg-brand-50 font-medium flex items-center gap-2 justify-center transition-colors"
      @click="$emit('lead-create', { column_id: column.id })"
    >
      <i class="fa-solid fa-plus"></i> Картаи нав
    </button>
  </div>
</template>

<script setup lang="ts">
// @ts-ignore - vuedraggable lacks types in v4 ESM
import draggable from 'vuedraggable'
import type { KanbanColumn as Col, KanbanLead, RequestsColumn, RequestsItem } from '~/types/api'

type AnyColumn = Col | RequestsColumn
type AnyLead = KanbanLead | RequestsItem

const props = defineProps<{
  column: AnyColumn
  leads: AnyLead[]
}>()
const emit = defineEmits<{
  (e: 'lead-move', payload: { lead: AnyLead; columnId: number; orderIndex: number }): void
  (e: 'lead-create', payload: { column_id: number }): void
  (e: 'lead-edit', lead: AnyLead): void
  (e: 'lead-delete', lead: AnyLead): void
  (e: 'column-edit', col: AnyColumn): void
  (e: 'column-delete', col: AnyColumn): void
}>()

// vuedraggable v4 fires `end` after move; the new list shape is reported via update:modelValue.
let pendingList: AnyLead[] = []

const onListUpdate = (newList: AnyLead[]) => {
  pendingList = newList
}

const onEnd = (e: any) => {
  // We compare positions to detect movement.
  const list = pendingList.length ? pendingList : props.leads
  list.forEach((lead, idx) => {
    const orig = props.leads.find(x => x.id === lead.id)
    const moved = !orig || (orig as any).column_id !== props.column.id || (orig as any).order_index !== idx
    if (moved) {
      emit('lead-move', { lead, columnId: props.column.id, orderIndex: idx })
    }
  })
  pendingList = []
}
</script>
