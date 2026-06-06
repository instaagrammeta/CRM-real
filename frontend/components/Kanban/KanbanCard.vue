<template>
  <div
    class="group bg-white rounded-xl border border-border/60 hover:border-brand-300 hover:shadow-card cursor-pointer transition-all p-3 select-none"
    @click="$emit('click')"
  >
    <!-- Lid card -->
    <template v-if="kind === 'lid'">
      <div class="flex items-start justify-between gap-2 mb-1.5">
        <div class="font-semibold text-ink text-sm truncate">{{ (lead as any).client_name || 'Бе ном' }}</div>
        <button
          class="opacity-0 group-hover:opacity-100 text-ink-mute hover:text-red-600 transition-opacity"
          @click.stop="$emit('delete')"
        >
          <i class="fa-solid fa-trash text-xs"></i>
        </button>
      </div>
      <div v-if="(lead as any).phone" class="text-xs text-ink-soft flex items-center gap-1.5 mb-1">
        <i class="fa-solid fa-phone text-[10px] text-brand-600"></i>
        {{ (lead as any).phone }}
      </div>
      <div v-if="(lead as any).topic" class="text-xs text-ink-soft truncate">{{ (lead as any).topic }}</div>

      <div class="flex items-center gap-1.5 mt-2 flex-wrap">
        <span v-if="(lead as any).source" class="badge-brand">{{ (lead as any).source }}</span>
        <span v-if="(lead as any).mortgage" class="badge-blue">Ипотека</span>
        <span v-if="(lead as any).box" class="badge-amber">Box</span>
      </div>

      <div v-if="(lead as any).author_name" class="flex items-center gap-1.5 mt-2 pt-2 border-t border-border/40">
        <Avatar :name="(lead as any).author_name" size="xs" />
        <span class="text-[11px] text-ink-mute truncate">{{ (lead as any).author_name }}</span>
        <span class="ml-auto text-[10px] text-ink-mute">{{ timeAgo((lead as any).created_date) }}</span>
      </div>
    </template>

    <!-- Request card -->
    <template v-else>
      <div class="flex items-start justify-between gap-2 mb-1.5">
        <div class="font-semibold text-ink text-sm truncate">
          {{ (lead as any).property_type || 'Заявка' }} · {{ (lead as any).rooms || '?' }} ҳуҷра
        </div>
        <button
          class="opacity-0 group-hover:opacity-100 text-ink-mute hover:text-red-600 transition-opacity"
          @click.stop="$emit('delete')"
        >
          <i class="fa-solid fa-trash text-xs"></i>
        </button>
      </div>
      <div v-if="(lead as any).address" class="text-xs text-ink-soft flex items-center gap-1.5 mb-1 truncate">
        <i class="fa-solid fa-location-dot text-[10px] text-brand-600"></i>
        {{ (lead as any).address }}
      </div>
      <div class="text-xs text-ink-soft flex items-center gap-3 flex-wrap">
        <span v-if="(lead as any).area"><i class="fa-solid fa-vector-square text-[10px] text-brand-600"></i> {{ (lead as any).area }} м²</span>
        <span v-if="(lead as any).floor"><i class="fa-solid fa-layer-group text-[10px] text-brand-600"></i> {{ (lead as any).floor }}/{{ (lead as any).total_floors || '?' }}</span>
      </div>
      <div v-if="(lead as any).total_price" class="text-sm font-bold text-brand-700 mt-2">
        {{ formatNumber((lead as any).total_price) }} <span class="text-xs font-normal text-ink-mute">USD</span>
      </div>
      <div v-if="(lead as any).client_name || (lead as any).phone" class="text-[11px] text-ink-mute mt-1.5 truncate">
        {{ (lead as any).client_name }} <span v-if="(lead as any).phone">· {{ (lead as any).phone }}</span>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import type { KanbanLead, RequestsItem } from '~/types/api'

const props = defineProps<{
  lead: KanbanLead | RequestsItem
  kind?: 'lid' | 'request'
}>()
defineEmits<{ (e: 'click'): void; (e: 'delete'): void }>()

const { timeAgo, formatNumber } = useFormat()

const kind = computed(() => props.kind || ('client_name' in props.lead && !('property_type' in props.lead) ? 'lid' : 'request'))
</script>
