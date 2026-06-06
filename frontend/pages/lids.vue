<template>
  <div class="space-y-4">
    <PageHeader title="Канбани лидҳо" description="Ҳама лидҳо аз рӯи ҳолат">
      <template #actions>
        <select v-model="boardId" class="input !py-2 !px-3 !w-auto" @change="loadAll">
          <option v-for="b in boards" :key="b.id" :value="b.id">{{ b.title }}</option>
        </select>
        <button class="btn-soft" @click="newBoardModal = true">
          <i class="fa-solid fa-plus"></i> Доска
        </button>
      </template>
    </PageHeader>

    <div v-if="loading" class="card flex items-center justify-center h-60">
      <Spinner size="lg" class="text-brand-600" />
    </div>

    <EmptyState
      v-else-if="boards.length === 0"
      icon="fa-clipboard"
      title="Доска нест"
      description="Аввалин доскаатонро созед то лидҳоро бо канбан идора кунед"
    >
      <template #action>
        <button class="btn-primary" @click="newBoardModal = true">
          <i class="fa-solid fa-plus"></i> Доскаи нав
        </button>
      </template>
    </EmptyState>

    <KanbanBoard
      v-else
      :columns="columns"
      :leads="leads"
      @lead-move="onLeadMove"
      @lead-create="(p) => openLeadEditor(null, p.column_id)"
      @lead-edit="(l) => openLeadEditor(l as KanbanLead)"
      @lead-delete="onLeadDelete"
      @columns-reorder="onColumnsReorder"
      @column-create="newColumnModal = true"
      @column-edit="(c) => openColumnEditor(c as KanbanColumn)"
      @column-delete="onColumnDelete"
    />

    <!-- New / edit board -->
    <Modal v-model="newBoardModal" :title="'Доскаи нав'">
      <div class="space-y-3">
        <div>
          <label class="label">Унвон</label>
          <input v-model="boardForm.title" class="input" placeholder="Лидҳои моҳи рав" />
        </div>
        <div>
          <label class="label">Ранг</label>
          <input v-model="boardForm.color" type="color" class="input !p-1 h-11 w-20" />
        </div>
      </div>
      <template #footer>
        <button class="btn-ghost" @click="newBoardModal = false">Бекор</button>
        <button class="btn-primary" @click="createBoard">Сохтан</button>
      </template>
    </Modal>

    <!-- New / edit column -->
    <Modal v-model="newColumnModal" :title="editingColumn ? 'Тағйири колонна' : 'Колоннаи нав'">
      <div class="space-y-3">
        <div>
          <label class="label">Унвон</label>
          <input v-model="columnForm.title" class="input" placeholder="Дар ҷараён" />
        </div>
        <div>
          <label class="label">Ранг</label>
          <input v-model="columnForm.color" type="color" class="input !p-1 h-11 w-20" />
        </div>
      </div>
      <template #footer>
        <button class="btn-ghost" @click="closeColumnEditor">Бекор</button>
        <button class="btn-primary" @click="saveColumn">Захира</button>
      </template>
    </Modal>

    <!-- New / edit lead -->
    <Modal v-model="leadModal" :title="editingLead ? 'Тағйири лид' : 'Лиди нав'" size="lg">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <div class="sm:col-span-2">
          <label class="label">Номи мизоҷ</label>
          <input v-model="leadForm.client_name" class="input" />
        </div>
        <div>
          <label class="label">Телефон</label>
          <input v-model="leadForm.phone" class="input" placeholder="+992..." />
        </div>
        <div>
          <label class="label">Манбаъ</label>
          <input v-model="leadForm.source" class="input" placeholder="Instagram, тамос..." />
        </div>
        <div class="sm:col-span-2">
          <label class="label">Мавзӯъ</label>
          <input v-model="leadForm.topic" class="input" />
        </div>
        <div class="sm:col-span-2">
          <label class="label">Эзоҳ</label>
          <textarea v-model="leadForm.comment" class="input min-h-[90px]"></textarea>
        </div>
        <label class="flex items-center gap-2 text-sm text-ink-soft">
          <input v-model="leadForm.mortgage" type="checkbox" class="w-4 h-4 accent-brand-600" />
          Ипотека
        </label>
        <label class="flex items-center gap-2 text-sm text-ink-soft">
          <input v-model="leadForm.box" type="checkbox" class="w-4 h-4 accent-brand-600" />
          Box
        </label>
      </div>
      <template #footer>
        <button class="btn-ghost" @click="leadModal = false">Бекор</button>
        <button class="btn-primary" @click="saveLead">Захира</button>
      </template>
    </Modal>

    <ConfirmDialog
      v-model="confirmOpen"
      :title="confirmTitle"
      :message="confirmMessage"
      danger
      :busy="busyConfirm"
      @confirm="onConfirm"
    />
  </div>
</template>

<script setup lang="ts">
import type { KanbanBoard, KanbanColumn, KanbanLead } from '~/types/api'

const api = useApi()
const toast = useToast()
const ws = useWebSocket()

const boards = ref<KanbanBoard[]>([])
const boardId = ref<number | null>(null)
const columns = ref<KanbanColumn[]>([])
const leads = ref<KanbanLead[]>([])
const loading = ref(true)

const newBoardModal = ref(false)
const newColumnModal = ref(false)
const leadModal = ref(false)
const editingColumn = ref<KanbanColumn | null>(null)
const editingLead = ref<KanbanLead | null>(null)

const boardForm = reactive({ title: '', color: '#1f7a4d' })
const columnForm = reactive({ title: '', color: '#1f7a4d' })
const leadForm = reactive<Partial<KanbanLead>>({
  client_name: '', phone: '', source: '', topic: '', comment: '', mortgage: false, box: false,
  column_id: 0, board_id: 0,
})

const confirmOpen = ref(false)
const confirmTitle = ref('')
const confirmMessage = ref('')
const busyConfirm = ref(false)
let confirmAction: (() => Promise<void>) | null = null

const askConfirm = (title: string, message: string, action: () => Promise<void>) => {
  confirmTitle.value = title
  confirmMessage.value = message
  confirmAction = action
  confirmOpen.value = true
}

const onConfirm = async () => {
  if (!confirmAction) return
  busyConfirm.value = true
  try { await confirmAction() } finally { busyConfirm.value = false; confirmOpen.value = false; confirmAction = null }
}

const loadBoards = async () => {
  boards.value = await api.get<KanbanBoard[]>('/api/kanban/boards')
  if (boards.value.length && (!boardId.value || !boards.value.find(b => b.id === boardId.value))) {
    boardId.value = boards.value[0].id
  }
}

const loadAll = async () => {
  if (!boardId.value) return
  loading.value = true
  try {
    const [cols, leadsRes] = await Promise.all([
      api.get<KanbanColumn[]>(`/api/kanban/boards/${boardId.value}/columns`),
      api.get<KanbanLead[]>(`/api/kanban/leads`, { board_id: boardId.value }),
    ])
    columns.value = cols
    leads.value = leadsRes
  } catch (e: any) {
    toast.error('Хато', e?.message || '')
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadBoards()
  await loadAll()
})

const createBoard = async () => {
  if (!boardForm.title.trim()) return
  const res = await api.post<{ id: number }>('/api/kanban/boards', boardForm)
  toast.success('Доска сохта шуд')
  newBoardModal.value = false
  boardForm.title = ''
  await loadBoards()
  boardId.value = res.id
  await loadAll()
}

const saveColumn = async () => {
  if (!columnForm.title.trim() || !boardId.value) return
  if (editingColumn.value) {
    await api.put(`/api/kanban/columns/${editingColumn.value.id}`, columnForm)
    toast.success('Колонна нав карда шуд')
  } else {
    await api.post('/api/kanban/columns', { ...columnForm, board_id: boardId.value })
    toast.success('Колонна сохта шуд')
  }
  newColumnModal.value = false
  editingColumn.value = null
  columnForm.title = ''
  await loadAll()
}

const closeColumnEditor = () => { newColumnModal.value = false; editingColumn.value = null; columnForm.title = '' }
const openColumnEditor = (col: KanbanColumn) => {
  editingColumn.value = col
  columnForm.title = col.title
  columnForm.color = col.color || '#1f7a4d'
  newColumnModal.value = true
}

const onColumnDelete = (col: any) => {
  askConfirm('Колоннаро нест кунем?', `Колонна "${col.title}" ва ҳама лидҳои он нест мешаванд.`, async () => {
    await api.delete(`/api/kanban/columns/${col.id}`)
    toast.success('Колонна нест шуд')
    await loadAll()
  })
}

const onColumnsReorder = async (cols: any[]) => {
  await Promise.all(cols.map(c => api.put(`/api/kanban/columns/${c.id}/order`, { order_index: c.order_index })))
}

const openLeadEditor = (lead: KanbanLead | null, columnId?: number) => {
  editingLead.value = lead
  if (lead) {
    Object.assign(leadForm, lead)
  } else {
    Object.assign(leadForm, {
      client_name: '', phone: '', source: '', topic: '', comment: '',
      mortgage: false, box: false,
      column_id: columnId || columns.value[0]?.id || 0,
      board_id: boardId.value || 0,
    })
  }
  leadModal.value = true
}

const saveLead = async () => {
  if (!leadForm.client_name?.trim()) return
  if (editingLead.value) {
    await api.put(`/api/kanban/leads/${editingLead.value.id}`, leadForm)
    toast.success('Лид нав карда шуд')
  } else {
    await api.post('/api/kanban/leads', leadForm)
    toast.success('Лид сохта шуд')
  }
  leadModal.value = false
  await loadAll()
}

const onLeadDelete = (lead: any) => {
  askConfirm('Лидро нест кунем?', `Лиди "${lead.client_name}" нест карда мешавад.`, async () => {
    await api.delete(`/api/kanban/leads/${lead.id}`)
    toast.success('Лид нест шуд')
    await loadAll()
  })
}

const onLeadMove = async ({ lead, columnId, orderIndex }: any) => {
  // optimistic update
  const found = leads.value.find(l => l.id === lead.id)
  if (found) { found.column_id = columnId; found.order_index = orderIndex }
  try {
    await api.post('/api/kanban/leads/move', { lead_id: lead.id, column_id: columnId, order_index: orderIndex })
  } catch (e: any) {
    toast.error('Хато ҳангоми кӯчонидан', e?.message || '')
    await loadAll()
  }
}

// Real-time updates
let unsubMoved: (() => void) | null = null
let unsubCreated: (() => void) | null = null
let unsubUpdated: (() => void) | null = null
onMounted(() => {
  unsubMoved   = ws.on('lead:moved',   () => loadAll())
  unsubCreated = ws.on('lead:created', () => loadAll())
  unsubUpdated = ws.on('lead:updated', () => loadAll())
})
onBeforeUnmount(() => {
  unsubMoved?.(); unsubCreated?.(); unsubUpdated?.()
})
</script>
