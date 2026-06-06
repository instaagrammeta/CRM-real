<template>
  <div class="space-y-4">
    <PageHeader title="Хонаҳо" description="Феҳристи объектҳои амлок">
      <template #actions>
        <input v-model="search" class="input !py-2 !w-64 hidden sm:block" placeholder="Ҷустуҷӯ..." />
        <button class="btn-primary" @click="openEditor(null)">
          <i class="fa-solid fa-plus"></i> Илова
        </button>
      </template>
    </PageHeader>

    <div v-if="loading" class="card flex justify-center py-12">
      <Spinner class="text-brand-600" size="lg" />
    </div>

    <EmptyState v-else-if="filtered.length === 0" icon="fa-building" title="Хона нест" />

    <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3 lg:gap-4">
      <div
        v-for="h in filtered"
        :key="h.id"
        class="card-hover cursor-pointer relative group"
        @click="openEditor(h)"
      >
        <div class="flex items-start justify-between gap-2 mb-2">
          <div class="flex-1 min-w-0">
            <h3 class="font-semibold text-ink truncate">{{ h.title || 'Бе ном' }}</h3>
            <div class="text-xs text-ink-mute truncate">{{ h.district }} · {{ h.address }}</div>
          </div>
          <span class="badge-brand">{{ h.construction_type || '—' }}</span>
        </div>
        <div class="grid grid-cols-3 gap-2 my-3 text-center">
          <div class="bg-page rounded-lg py-2">
            <div class="text-xs text-ink-mute">Масоҳат</div>
            <div class="font-semibold">{{ h.area || '—' }} м²</div>
          </div>
          <div class="bg-page rounded-lg py-2">
            <div class="text-xs text-ink-mute">Ҳуҷра</div>
            <div class="font-semibold">{{ h.rooms || '—' }}</div>
          </div>
          <div class="bg-page rounded-lg py-2">
            <div class="text-xs text-ink-mute">Ошёна</div>
            <div class="font-semibold">{{ h.floor }}/{{ h.total_floors }}</div>
          </div>
        </div>
        <div class="flex items-center justify-between">
          <div class="text-lg font-bold text-brand-700">
            {{ formatNumber(h.total_price) }}
            <span class="text-xs font-normal text-ink-mute">USD</span>
          </div>
          <button
            class="opacity-0 group-hover:opacity-100 text-ink-mute hover:text-red-600 transition"
            @click.stop="onDelete(h)"
          >
            <i class="fa-solid fa-trash text-sm"></i>
          </button>
        </div>
      </div>
    </div>

    <Modal v-model="modal" :title="editing ? 'Тағйири хона' : 'Хонаи нав'" size="xl">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div class="md:col-span-2"><label class="label">Унвон</label><input v-model="form.title" class="input" /></div>
        <div><label class="label">Намуди бино</label><input v-model="form.construction_type" class="input" /></div>
        <div><label class="label">Ноҳия</label><input v-model="form.district" class="input" /></div>
        <div class="md:col-span-2"><label class="label">Манзил</label><input v-model="form.address" class="input" /></div>
        <div><label class="label">Масоҳат</label><input v-model.number="form.area" type="number" step="0.1" class="input" /></div>
        <div><label class="label">Ҳуҷраҳо</label><input v-model.number="form.rooms" type="number" class="input" /></div>
        <div><label class="label">Тиреза</label><input v-model.number="form.windows" type="number" class="input" /></div>
        <div><label class="label">Ошёна</label><input v-model.number="form.floor" type="number" class="input" /></div>
        <div><label class="label">Аз ҳамагӣ</label><input v-model.number="form.total_floors" type="number" class="input" /></div>
        <div><label class="label">Нарх / м²</label><input v-model.number="form.price_per_m2" type="number" class="input" /></div>
        <div><label class="label">Нархи умумӣ</label><input v-model.number="form.total_price" type="number" class="input" /></div>
        <div><label class="label">Bauer</label><input v-model="form.developer" class="input" /></div>
        <div><label class="label">Тел.</label><input v-model="form.contact_phone" class="input" /></div>
      </div>
      <template #footer>
        <button class="btn-ghost" @click="modal = false">Бекор</button>
        <button class="btn-primary" @click="save">Захира</button>
      </template>
    </Modal>

    <ConfirmDialog v-model="confirmOpen" message="Хона нест карда мешавад" danger
                   :busy="busyConfirm" @confirm="confirmDelete" />
  </div>
</template>

<script setup lang="ts">
import type { House } from '~/types/api'

const api = useApi()
const toast = useToast()
const { formatNumber } = useFormat()

const houses = ref<House[]>([])
const loading = ref(true)
const search = ref('')

const modal = ref(false)
const editing = ref<House | null>(null)
const form = reactive<Partial<House>>({})

const confirmOpen = ref(false)
const busyConfirm = ref(false)
let pending: House | null = null

const load = async () => {
  loading.value = true
  try { houses.value = await api.get<House[]>('/api/houses') }
  finally { loading.value = false }
}
onMounted(load)

const filtered = computed(() => {
  const s = search.value.trim().toLowerCase()
  if (!s) return houses.value
  return houses.value.filter(h =>
    [h.title, h.address, h.district, h.developer].some(v => (v || '').toLowerCase().includes(s)),
  )
})

const openEditor = (h: House | null) => {
  editing.value = h
  Object.keys(form).forEach(k => delete (form as any)[k])
  if (h) Object.assign(form, h)
  modal.value = true
}

const save = async () => {
  if (editing.value) await api.put(`/api/houses/${editing.value.id}`, form)
  else await api.post('/api/houses', form)
  modal.value = false
  await load()
  toast.success('Захира шуд')
}

const onDelete = (h: House) => { pending = h; confirmOpen.value = true }
const confirmDelete = async () => {
  if (!pending) return
  busyConfirm.value = true
  try { await api.delete(`/api/houses/${pending.id}`); await load(); toast.success('Нест шуд') }
  finally { busyConfirm.value = false; confirmOpen.value = false; pending = null }
}
</script>
