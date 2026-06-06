<template>
  <div class="space-y-4">
    <PageHeader title="SIM-кортҳо" description="Идоракунии SIM ва тарифҳо">
      <template #actions>
        <button class="btn-primary" @click="openEditor(null)">
          <i class="fa-solid fa-plus"></i> SIM нав
        </button>
      </template>
    </PageHeader>

    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
      <StatCard icon="fa-mobile-screen" label="Телефонҳо" :value="stats.phones || 0" variant="violet" />
      <StatCard icon="fa-sim-card" label="SIM кортҳо" :value="stats.sims || 0" variant="brand" />
      <StatCard icon="fa-circle-check" label="Фаъол" :value="stats.active_sims || 0" variant="blue" />
      <StatCard icon="fa-clock" label="Тарифи гузашта" :value="stats.expired_tariffs || 0" variant="amber" />
    </div>

    <div v-if="loading" class="card flex justify-center py-12"><Spinner size="lg" class="text-brand-600" /></div>

    <div v-else class="card !p-0 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-page text-ink-mute uppercase text-[11px]">
          <tr>
            <th class="text-left px-4 py-3">Рақам</th>
            <th class="text-left px-4 py-3">Оператор</th>
            <th class="text-left px-4 py-3">Ҳолат</th>
            <th class="text-left px-4 py-3">Tasvir</th>
            <th class="text-right px-4 py-3 w-32">Амалҳо</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in sims" :key="s.id" class="border-t border-border/40 hover:bg-page/60">
            <td class="px-4 py-3 font-medium">{{ s.phone_number }}</td>
            <td class="px-4 py-3">{{ s.operator }}</td>
            <td class="px-4 py-3">
              <span :class="s.status === 'active' ? 'badge-brand' : 'badge-gray'">{{ s.status }}</span>
            </td>
            <td class="px-4 py-3 text-ink-soft truncate max-w-xs">{{ s.description }}</td>
            <td class="px-4 py-3 text-right">
              <button class="text-brand-700 hover:underline text-xs mr-3" @click="openEditor(s)">
                <i class="fa-solid fa-pen"></i>
              </button>
              <button class="text-red-600 hover:underline text-xs" @click="onDelete(s)">
                <i class="fa-solid fa-trash"></i>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      <EmptyState v-if="sims.length === 0" icon="fa-sim-card" title="SIM нест" />
    </div>

    <Modal v-model="modal" :title="editing ? 'Тағйири SIM' : 'SIM-и нав'">
      <div class="space-y-3">
        <div><label class="label">Рақами телефон *</label><input v-model="form.phone_number" class="input" /></div>
        <div><label class="label">Оператор *</label>
          <select v-model="form.operator" class="input">
            <option>Megafon</option><option>Tcell</option><option>Beeline</option><option>Babilon</option><option>ZetMobile</option>
          </select>
        </div>
        <div><label class="label">Ҳолат</label>
          <select v-model="form.status" class="input">
            <option value="active">Фаъол</option><option value="inactive">Ғайри фаъол</option>
          </select>
        </div>
        <div><label class="label">Тавсиф</label><textarea v-model="form.description" class="input"></textarea></div>
      </div>
      <template #footer>
        <button class="btn-ghost" @click="modal = false">Бекор</button>
        <button class="btn-primary" @click="save">Захира</button>
      </template>
    </Modal>

    <ConfirmDialog v-model="confirmOpen" message="SIM нест карда мешавад" danger :busy="busy" @confirm="confirmDel" />
  </div>
</template>

<script setup lang="ts">
import type { SimCard } from '~/types/api'

const api = useApi()
const toast = useToast()

const sims = ref<SimCard[]>([])
const stats = ref<any>({})
const loading = ref(true)

const modal = ref(false)
const editing = ref<SimCard | null>(null)
const form = reactive<Partial<SimCard>>({})

const confirmOpen = ref(false)
const busy = ref(false)
let pending: SimCard | null = null

const load = async () => {
  loading.value = true
  try {
    sims.value = await api.get<SimCard[]>('/api/sim-cards')
    stats.value = await api.get<any>('/api/sim-phones/stats')
  } finally { loading.value = false }
}
onMounted(load)

const openEditor = (s: SimCard | null) => {
  editing.value = s
  Object.keys(form).forEach(k => delete (form as any)[k])
  if (s) Object.assign(form, s)
  else Object.assign(form, { operator: 'Megafon', status: 'active' })
  modal.value = true
}

const save = async () => {
  if (!form.phone_number?.trim()) return
  if (editing.value) await api.put(`/api/sim-cards/${editing.value.id}`, form)
  else await api.post('/api/sim-cards', form)
  modal.value = false; await load(); toast.success('Захира шуд')
}

const onDelete = (s: SimCard) => { pending = s; confirmOpen.value = true }
const confirmDel = async () => {
  if (!pending) return
  busy.value = true
  try { await api.delete(`/api/sim-cards/${pending.id}`); await load(); toast.success('Нест шуд') }
  finally { busy.value = false; confirmOpen.value = false; pending = null }
}
</script>
