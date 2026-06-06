<template>
  <div class="space-y-4">
    <PageHeader title="Канбани заявкаҳо" description="Идоракунии заявкаҳои амлок">
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
      icon="fa-clipboard-list"
      title="Доска нест"
      description="Аввалин доскаатонро созед"
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
      :leads="items as any"
      @lead-move="onMove"
      @lead-create="(p) => openEditor(null, p.column_id)"
      @lead-edit="(it) => openEditor(it as RequestsItem)"
      @lead-delete="onDelete"
      @column-create="newColumnModal = true"
      @column-edit="(c) => openColumnEditor(c as RequestsColumn)"
      @column-delete="onColumnDelete"
    />

    <Modal v-model="newBoardModal" title="Доскаи нав">
      <div class="space-y-3">
        <div><label class="label">Унвон</label><input v-model="boardForm.title" class="input" /></div>
        <div><label class="label">Ранг</label><input v-model="boardForm.color" type="color" class="input !p-1 h-11 w-20" /></div>
      </div>
      <template #footer>
        <button class="btn-ghost" @click="newBoardModal = false">Бекор</button>
        <button class="btn-primary" @click="createBoard">Сохтан</button>
      </template>
    </Modal>

    <Modal v-model="newColumnModal" :title="editingColumn ? 'Тағйири колонна' : 'Колоннаи нав'">
      <div class="space-y-3">
        <div><label class="label">Унвон</label><input v-model="columnForm.title" class="input" /></div>
        <div><label class="label">Ранг</label><input v-model="columnForm.color" type="color" class="input !p-1 h-11 w-20" /></div>
      </div>
      <template #footer>
        <button class="btn-ghost" @click="closeColumnEditor">Бекор</button>
        <button class="btn-primary" @click="saveColumn">Захира</button>
      </template>
    </Modal>

    <Modal v-model="itemModal" :title="editingItem ? 'Тағйири заявка' : 'Заявкаи нав'" size="xl">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div><label class="label">Намуди амлок</label><input v-model="itemForm.property_type" class="input" /></div>
        <div><label class="label">Манзил</label><input v-model="itemForm.address" class="input" /></div>
        <div><label class="label">Масоҳат (м²)</label><input v-model.number="itemForm.area" type="number" step="0.1" class="input" /></div>
        <div><label class="label">Ҳуҷраҳо</label><input v-model.number="itemForm.rooms" type="number" class="input" /></div>
        <div><label class="label">Тиреза</label><input v-model.number="itemForm.windows" type="number" class="input" /></div>
        <div><label class="label">Ошёна</label><input v-model.number="itemForm.floor" type="number" class="input" /></div>
        <div><label class="label">Аз ҳамаги ошёна</label><input v-model.number="itemForm.total_floors" type="number" class="input" /></div>
        <div><label class="label">Нархи умумӣ</label><input v-model.number="itemForm.total_price" type="number" class="input" /></div>
        <div><label class="label">Нарх / м²</label><input v-model.number="itemForm.price_per_m2" type="number" class="input" /></div>
        <div><label class="label">Телефон</label><input v-model="itemForm.phone" class="input" /></div>
        <div><label class="label">Номи мизоҷ</label><input v-model="itemForm.client_name" class="input" /></div>
        <div class="md:col-span-2">
          <label class="label">Эзоҳ</label>
          <textarea v-model="itemForm.comment" class="input min-h-[80px]"></textarea>
        </div>
      </div>
      <template #footer>
        <button class="btn-ghost" @click="itemModal = false">Бекор</button>
        <button class="btn-primary" @click="saveItem">Захира</button>
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
import type { RequestsBoard, RequestsColumn, RequestsItem } from '~/types/api'

const api = useApi()
const toast = useToast()
const ws = useWebSocket()

const boards = ref<RequestsBoard[]>([])
const boardId = ref<number | null>(null)
const columns = ref<RequestsColumn[]>([])
const items = ref<RequestsItem[]>([])
const loading = ref(true)

const newBoardModal = ref(false)
const newColumnModal = ref(false)
const itemModal = ref(false)
const editingColumn = ref<RequestsColumn | null>(null)
const editingItem = ref<RequestsItem | null>(null)

const boardForm = reactive({ title: '', color: '#1f7a4d' })
const columnForm = reactive({ title: '', color: '#1f7a4d' })
const itemForm = reactive<Partial<RequestsItem>>({})

const confirmOpen = ref(false)
const confirmTitle = ref(''); const confirmMessage = ref('')
const busyConfirm = ref(false)
let confirmAction: (() => Promise<void>) | null = null

const askConfirm = (t: string, m: string, a: () => Promise<void>) => {
  confirmTitle.value = t; confirmMessage.value = m; confirmAction = a; confirmOpen.value = true
}
const onConfirm = async () => {
  if (!confirmAction) return
  busyConfirm.value = true
  try { await confirmAction() } finally { busyConfirm.value = false; confirmOpen.value = false; confirmAction = null }
}

const loadBoards = async () => {
  boards.value = await api.get<RequestsBoard[]>('/api/requests-boards')
  if (boards.value.length && (!boardId.value || !boards.value.find(b => b.id === boardId.value))) {
    boardId.value = boards.value[0].id
  }
}

const loadAll = async () => {
  if (!boardId.value) return
  loading.value = true
  try {
    const [cols, itemsRes] = await Promise.all([
      api.get<RequestsColumn[]>(`/api/requests-boards/${boardId.value}/columns`),
      api.get<RequestsItem[]>(`/api/requests-board/${boardId.value}/requests`),
    ])
    columns.value = cols
    items.value = itemsRes
  } catch (e: any) {
    toast.error('Хато', e?.message || '')
  } finally { loading.value = false }
}

onMounted(async () => { await loadBoards(); await loadAll() })

const createBoard = async () => {
  if (!boardForm.title.trim()) return
  const r = await api.post<{ id: number }>('/api/requests-boards', boardForm)
  newBoardModal.value = false
  boardForm.title = ''
  await loadBoards(); boardId.value = r.id; await loadAll()
  toast.success('Доска сохта шуд')
}

const closeColumnEditor = () => { newColumnModal.value = false; editingColumn.value = null; columnForm.title = '' }
const openColumnEditor = (c: RequestsColumn) => {
  editingColumn.value = c
  columnForm.title = c.title
  columnForm.color = c.color || '#1f7a4d'
  newColumnModal.value = true
}
const saveColumn = async () => {
  if (!columnForm.title.trim() || !boardId.value) return
  if (editingColumn.value) {
    await api.put(`/api/requests-columns/${editingColumn.value.id}`, columnForm)
  } else {
    await api.post('/api/requests-columns', { ...columnForm, board_id: boardId.value })
  }
  closeColumnEditor()
  await loadAll()
  toast.success('Захира шуд')
}

const onColumnDelete = (c: any) => {
  askConfirm('Колоннаро нест кунем?', `"${c.title}" ва заявкаҳои он нест мешаванд.`, async () => {
    await api.delete(`/api/requests-columns/${c.id}`)
    await loadAll()
    toast.success('Нест шуд')
  })
}

const openEditor = (item: RequestsItem | null, columnId?: number) => {
  editingItem.value = item
  if (item) Object.assign(itemForm, item)
  else {
    Object.keys(itemForm).forEach(k => delete (itemForm as any)[k])
    Object.assign(itemForm, {
      column_id: columnId || columns.value[0]?.id,
      board_id: boardId.value,
    })
  }
  itemModal.value = true
}

const saveItem = async () => {
  if (editingItem.value) {
    await api.put(`/api/requests-update/${editingItem.value.id}`, itemForm)
  } else {
    await api.post('/api/requests-new', itemForm)
  }
  itemModal.value = false
  await loadAll()
  toast.success('Захира шуд')
}

const onDelete = (it: any) => {
  askConfirm('Заявкаро нест кунем?', it.address || '', async () => {
    await api.delete(`/api/requests-delete/${it.id}`)
    await loadAll()
    toast.success('Нест шуд')
  })
}

const onMove = async ({ lead, columnId, orderIndex }: any) => {
  const found = items.value.find(i => i.id === lead.id)
  if (found) { found.column_id = columnId; found.order_index = orderIndex }
  try {
    await api.post('/api/requests-move', { request_id: lead.id, column_id: columnId, order_index: orderIndex })
  } catch (e: any) {
    toast.error('Хато', e?.message || ''); await loadAll()
  }
}

let unsubA: (() => void) | null = null
let unsubB: (() => void) | null = null
onMounted(() => {
  unsubA = ws.on('request:created', () => loadAll())
  unsubB = ws.on('request:moved',   () => loadAll())
})
onBeforeUnmount(() => { unsubA?.(); unsubB?.() })
</script>
