<template>
  <div class="kanban-board flex gap-4 overflow-x-auto pb-4 -mx-1 px-1 min-h-[60vh]">
    <draggable
      v-model="columnsModel"
      :item-key="(c: any) => c.id"
      handle=".kanban-col-handle"
      group="columns"
      class="flex gap-4"
      :animation="200"
      @end="onColumnDragEnd"
    >
      <template #item="{ element: column }">
        <KanbanColumn
          :column="column"
          :leads="leadsByColumn[column.id] || []"
          @lead-move="onLeadMove"
          @lead-create="(payload) => $emit('lead-create', { ...payload, column_id: column.id })"
          @lead-edit="(lead) => $emit('lead-edit', lead)"
          @lead-delete="(lead) => $emit('lead-delete', lead)"
          @column-edit="(col) => $emit('column-edit', col)"
          @column-delete="(col) => $emit('column-delete', col)"
        />
      </template>
    </draggable>

    <!-- Add column trigger -->
    <button
      class="shrink-0 w-[300px] h-fit border-2 border-dashed border-brand-200 hover:border-brand-400 hover:bg-brand-50/40 rounded-2xl p-4 text-brand-700 font-medium text-sm flex items-center justify-center gap-2 transition-colors"
      @click="$emit('column-create')"
    >
      <i class="fa-solid fa-plus"></i> Колоннаи нав
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
  columns: AnyColumn[]
  leads: AnyLead[]
}>()

const emit = defineEmits<{
  (e: 'columns-reorder', cols: AnyColumn[]): void
  (e: 'lead-move', payload: { lead: AnyLead; columnId: number; orderIndex: number }): void
  (e: 'lead-create', payload: { column_id: number }): void
  (e: 'lead-edit', lead: AnyLead): void
  (e: 'lead-delete', lead: AnyLead): void
  (e: 'column-create'): void
  (e: 'column-edit', col: AnyColumn): void
  (e: 'column-delete', col: AnyColumn): void
}>()

const columnsModel = ref<AnyColumn[]>([])

watch(() => props.columns, (cols) => {
  columnsModel.value = [...cols].sort((a, b) => a.order_index - b.order_index)
}, { immediate: true })

const leadsByColumn = computed<Record<number, AnyLead[]>>(() => {
  const map: Record<number, AnyLead[]> = {}
  for (const l of props.leads) {
    const cid = (l as any).column_id
    if (!map[cid]) map[cid] = []
    map[cid].push(l)
  }
  for (const k of Object.keys(map)) {
    map[+k].sort((a, b) => (a as any).order_index - (b as any).order_index)
  }
  return map
})

const onColumnDragEnd = () => {
  emit('columns-reorder', columnsModel.value.map((c, i) => ({ ...c, order_index: i })))
}

const onLeadMove = (payload: { lead: AnyLead; columnId: number; orderIndex: number }) => {
  emit('lead-move', payload)
}
</script>

<style scoped>
.kanban-board { scroll-snap-type: x proximity; }
.kanban-board > * { scroll-snap-align: start; }
</style>
