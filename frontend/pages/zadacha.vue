<template>
  <div class="space-y-4">
    <PageHeader title="Вазифаҳо" description="Идоракунии вазифаҳои корӣ">
      <template #actions>
        <button class="btn-primary" @click="openEditor(null)">
          <i class="fa-solid fa-plus"></i> Вазифа нав
        </button>
      </template>
    </PageHeader>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-3 lg:gap-4">
      <div v-for="col in cols" :key="col.key" class="card !p-3">
        <div class="flex items-center justify-between mb-3 px-2">
          <h3 class="font-semibold text-ink flex items-center gap-2">
            <span class="w-2 h-2 rounded-full" :class="col.dot" />
            {{ col.label }}
          </h3>
          <span class="badge-gray">{{ tasksByStatus[col.key]?.length || 0 }}</span>
        </div>
        <div class="space-y-2 max-h-[70vh] overflow-y-auto pr-1">
          <div
            v-for="t in tasksByStatus[col.key] || []"
            :key="t.id"
            class="bg-white border border-border/50 rounded-xl p-3 hover:shadow-card transition-shadow cursor-pointer"
            @click="openEditor(t)"
          >
            <div class="flex items-start justify-between gap-2 mb-1">
              <div class="font-semibold text-ink text-sm">{{ t.title }}</div>
              <button class="text-ink-mute hover:text-red-600 shrink-0" @click.stop="onDelete(t)">
                <i class="fa-solid fa-trash text-xs"></i>
              </button>
            </div>
            <p v-if="t.description" class="text-xs text-ink-soft line-clamp-2 mb-2">{{ t.description }}</p>
            <div class="flex items-center justify-between text-[11px] text-ink-mute">
              <span><i class="fa-regular fa-clock"></i> {{ formatDate(t.created_date) }}</span>
              <select
                :value="t.status"
                class="text-xs bg-page rounded-lg px-2 py-0.5 border-0 focus:ring-1 focus:ring-brand-300"
                @click.stop
                @change="onStatus(t, ($event.target as HTMLSelectElement).value)"
              >
                <option value="new">Нав</option>
                <option value="progress">Дар ҷараён</option>
                <option value="done">Анҷомёфта</option>
              </select>
            </div>
          </div>
          <div v-if="!tasksByStatus[col.key]?.length" class="text-center text-xs text-ink-mute py-6">холӣ</div>
        </div>
      </div>
    </div>

    <Modal v-model="modal" :title="editing ? 'Тағйири вазифа' : 'Вазифаи нав'" size="lg">
      <div class="space-y-3">
        <div>
          <label class="label">Унвон *</label>
          <input v-model="form.title" class="input" />
        </div>
        <div>
          <label class="label">Тавсиф</label>
          <textarea v-model="form.description" class="input min-h-[120px]"></textarea>
        </div>
        <div>
          <label class="label">Иҷрокунанда</label>
          <select v-model="form.executor_id" class="input">
            <option :value="0">— Интихоб —</option>
            <option v-for="u in users" :key="u.id" :value="u.id">{{ u.full_name }}</option>
          </select>
        </div>
        <div>
          <label class="label">Ҳолат</label>
          <select v-model="form.status" class="input">
            <option value="new">Нав</option>
            <option value="progress">Дар ҷараён</option>
            <option value="done">Анҷомёфта</option>
          </select>
        </div>
      </div>
      <template #footer>
        <button class="btn-ghost" @click="modal = false">Бекор</button>
        <button class="btn-primary" @click="save">Захира</button>
      </template>
    </Modal>

    <ConfirmDialog v-model="confirmOpen" title="Нест кардан?" message="Вазифа нест карда мешавад" danger
                   :busy="busy" @confirm="confirmDelete" />
  </div>
</template>

<script setup lang="ts">
import type { Task, User } from '~/types/api'

const api = useApi()
const toast = useToast()
const { formatDate } = useFormat()

const tasks = ref<Task[]>([])
const users = ref<Pick<User, 'id' | 'full_name'>[]>([])

const cols = [
  { key: 'new',      label: 'Нав',         dot: 'bg-blue-500' },
  { key: 'progress', label: 'Дар ҷараён',  dot: 'bg-amber-500' },
  { key: 'done',     label: 'Анҷомёфта',   dot: 'bg-brand-500' },
]

const tasksByStatus = computed<Record<string, Task[]>>(() => {
  const m: Record<string, Task[]> = {}
  for (const t of tasks.value) {
    const k = t.status || 'new'
    if (!m[k]) m[k] = []
    m[k].push(t)
  }
  return m
})

const modal = ref(false)
const editing = ref<Task | null>(null)
const form = reactive<Partial<Task>>({ title: '', description: '', executor_id: 0, status: 'new' })

const confirmOpen = ref(false)
const busy = ref(false)
let pendingDelete: Task | null = null

const load = async () => {
  tasks.value = await api.get<Task[]>('/api/tasks')
  users.value = await api.get<any[]>('/api/users/list')
}
onMounted(load)

const openEditor = (t: Task | null) => {
  editing.value = t
  if (t) Object.assign(form, t)
  else { form.title = ''; form.description = ''; form.executor_id = 0; form.status = 'new' }
  modal.value = true
}

const save = async () => {
  if (!form.title?.trim()) return
  if (editing.value) {
    await api.put(`/api/tasks/${editing.value.id}`, form)
  } else {
    await api.post('/api/tasks', form)
  }
  modal.value = false
  await load()
  toast.success('Захира шуд')
}

const onStatus = async (t: Task, status: string) => {
  await api.put(`/api/tasks/${t.id}`, { status })
  await load()
}

const onDelete = (t: Task) => { pendingDelete = t; confirmOpen.value = true }
const confirmDelete = async () => {
  if (!pendingDelete) return
  busy.value = true
  try {
    await api.delete(`/api/tasks/${pendingDelete.id}`)
    await load()
    toast.success('Нест шуд')
  } finally { busy.value = false; confirmOpen.value = false; pendingDelete = null }
}
</script>
