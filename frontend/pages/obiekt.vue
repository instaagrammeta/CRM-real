<template>
  <div class="space-y-4">
    <PageHeader title="Шахматка" description="Лоиҳаҳои сохтмон ва квартираҳо">
      <template #actions>
        <select v-model="projectId" class="input !py-2 !px-3 !w-auto" @change="loadProject">
          <option v-for="p in projects" :key="p.id" :value="p.id">{{ p.name }}</option>
        </select>
        <button class="btn-soft" @click="newProjectModal = true">
          <i class="fa-solid fa-plus"></i> Лоиҳа
        </button>
      </template>
    </PageHeader>

    <div v-if="loading" class="card flex justify-center py-12"><Spinner size="lg" class="text-brand-600" /></div>

    <EmptyState v-else-if="projects.length === 0" icon="fa-table-cells" title="Лоиҳа нест" description="Лоиҳаи нав созед">
      <template #action><button class="btn-primary" @click="newProjectModal = true"><i class="fa-solid fa-plus"></i> Лоиҳаи нав</button></template>
    </EmptyState>

    <div v-else-if="blocks.length === 0">
      <EmptyState icon="fa-cubes" title="Блок нест" description="Барои оғоз блок илова кунед" />
    </div>

    <!-- Blocks grid (read-only matrix) -->
    <div v-else class="space-y-4">
      <div v-for="b in blocks" :key="b.id" class="card">
        <h3 class="font-bold text-ink mb-3">{{ b.name }} <span class="text-sm text-ink-mute font-normal">{{ b.floor_from }}—{{ b.floor_to }} ошёна</span></h3>
        <div class="grid grid-cols-3 sm:grid-cols-5 lg:grid-cols-8 gap-2">
          <button
            v-for="apt in apartmentsByBlock[b.id] || []"
            :key="apt.id"
            class="rounded-lg p-2 text-center text-xs font-medium hover:opacity-80 transition-opacity"
            :class="statusClass(apt.status)"
            @click="onApartmentClick(apt)"
          >
            <div class="text-base font-bold">{{ apt.floor }}</div>
            <div class="opacity-80">{{ apt.rooms }}к / {{ apt.area }}м²</div>
          </button>
        </div>
      </div>
    </div>

    <Modal v-model="newProjectModal" title="Лоиҳаи нав">
      <div class="space-y-3">
        <div><label class="label">Ном *</label><input v-model="projectForm.name" class="input" /></div>
        <div><label class="label">Манзил</label><input v-model="projectForm.address" class="input" /></div>
        <div><label class="label">Сохтмонгар</label><input v-model="projectForm.developer" class="input" /></div>
        <div><label class="label">Тавсиф</label><textarea v-model="projectForm.description" class="input"></textarea></div>
      </div>
      <template #footer>
        <button class="btn-ghost" @click="newProjectModal = false">Бекор</button>
        <button class="btn-primary" @click="createProject">Сохтан</button>
      </template>
    </Modal>

    <Modal v-model="aptModal" :title="`Хонаи №${selectedApt?.floor || ''}`" size="lg">
      <div v-if="selectedApt" class="grid grid-cols-2 gap-3">
        <div><label class="label">Ҳолат</label>
          <select v-model="aptForm.status" class="input">
            <option value="free">Озод</option>
            <option value="reserved">Захираи мизоҷ</option>
            <option value="sold">Фурӯхта</option>
          </select>
        </div>
        <div><label class="label">Ҳуҷраҳо</label><input v-model.number="aptForm.rooms" type="number" class="input" /></div>
        <div><label class="label">Масоҳат</label><input v-model.number="aptForm.area" type="number" step="0.1" class="input" /></div>
        <div><label class="label">Тиреза</label><input v-model.number="aptForm.windows" type="number" class="input" /></div>
        <div><label class="label">Нарх / м²</label><input v-model.number="aptForm.price_per_m2" type="number" class="input" /></div>
        <div><label class="label">Нархи умумӣ</label><input v-model.number="aptForm.total_price" type="number" class="input" /></div>
      </div>
      <template #footer>
        <button class="btn-ghost" @click="aptModal = false">Бекор</button>
        <button class="btn-primary" @click="saveApt">Захира</button>
      </template>
    </Modal>
  </div>
</template>

<script setup lang="ts">
const api = useApi()
const toast = useToast()

const projects = ref<any[]>([])
const projectId = ref<number | null>(null)
const blocks = ref<any[]>([])
const apartments = ref<any[]>([])
const loading = ref(true)

const newProjectModal = ref(false)
const projectForm = reactive({ name: '', address: '', developer: '', description: '' })

const aptModal = ref(false)
const selectedApt = ref<any | null>(null)
const aptForm = reactive<any>({})

const apartmentsByBlock = computed<Record<number, any[]>>(() => {
  const m: Record<number, any[]> = {}
  for (const a of apartments.value) {
    if (!m[a.block_id]) m[a.block_id] = []
    m[a.block_id].push(a)
  }
  for (const k of Object.keys(m)) m[+k].sort((x, y) => x.floor - y.floor)
  return m
})

const statusClass = (s: string) => ({
  free:     'bg-brand-100 text-brand-800',
  reserved: 'bg-amber-100 text-amber-800',
  sold:     'bg-gray-200 text-gray-600',
}[s] || 'bg-gray-100')

const loadProjects = async () => {
  projects.value = await api.get<any[]>('/api/objekt/projects')
  if (projects.value.length && !projectId.value) projectId.value = projects.value[0].id
}

const loadProject = async () => {
  if (!projectId.value) return
  loading.value = true
  try {
    blocks.value = await api.get<any[]>(`/api/objekt/projects/${projectId.value}/blocks`)
    apartments.value = await api.get<any[]>(`/api/objekt/projects/${projectId.value}/apartments`)
  } finally { loading.value = false }
}

onMounted(async () => { await loadProjects(); await loadProject() })

const createProject = async () => {
  if (!projectForm.name.trim()) return
  const r = await api.post<{ id: number }>('/api/objekt/projects', projectForm)
  newProjectModal.value = false
  Object.assign(projectForm, { name: '', address: '', developer: '', description: '' })
  await loadProjects()
  projectId.value = r.id
  await loadProject()
  toast.success('Лоиҳа сохта шуд')
}

const onApartmentClick = (apt: any) => {
  selectedApt.value = apt
  Object.keys(aptForm).forEach(k => delete aptForm[k])
  Object.assign(aptForm, apt)
  aptModal.value = true
}

const saveApt = async () => {
  if (!selectedApt.value) return
  await api.put(`/api/objekt/apartments/${selectedApt.value.id}`, aptForm)
  aptModal.value = false
  await loadProject()
  toast.success('Захира шуд')
}
</script>
