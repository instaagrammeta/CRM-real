<template>
  <div class="space-y-4">
    <PageHeader title="Идоракунии бонкҳо" description="CRUD-и бонкҳои ипотека">
      <template #actions>
        <button class="btn-primary" @click="openEditor(null)">
          <i class="fa-solid fa-plus"></i> Бонк
        </button>
      </template>
    </PageHeader>

    <div v-if="loading" class="card flex justify-center py-12"><Spinner size="lg" class="text-brand-600" /></div>
    <EmptyState v-else-if="banks.length === 0" icon="fa-building-columns" title="Бонк нест" />

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-3 lg:gap-4">
      <div v-for="b in banks" :key="b.id" class="card-hover">
        <div class="flex items-start gap-3">
          <div class="w-12 h-12 rounded-xl bg-brand-50 grid place-items-center text-brand-700">
            <i class="fa-solid fa-building-columns"></i>
          </div>
          <div class="flex-1 min-w-0">
            <h3 class="font-bold text-ink truncate">{{ b.name }}</h3>
            <div class="text-xs text-ink-mute">/{{ b.slug }}</div>
            <p class="text-sm text-ink-soft mt-1 line-clamp-2">{{ b.description }}</p>
          </div>
          <div class="flex flex-col gap-1.5">
            <button class="btn-ghost !p-2" @click="openEditor(b)" title="Тағйир">
              <i class="fa-solid fa-pen text-xs"></i>
            </button>
            <button class="btn-ghost !p-2 hover:!text-red-600" @click="onDelete(b)" title="Нест">
              <i class="fa-solid fa-trash text-xs"></i>
            </button>
          </div>
        </div>
      </div>
    </div>

    <Modal v-model="modal" :title="editing ? 'Тағйири бонк' : 'Бонки нав'" size="lg">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div class="md:col-span-2"><label class="label">Ном *</label><input v-model="form.name" class="input" /></div>
        <div><label class="label">Slug</label><input v-model="form.slug" class="input" placeholder="auto" /></div>
        <div><label class="label">Тел.</label><input v-model="form.phone" class="input" /></div>
        <div class="md:col-span-2"><label class="label">Логотип (URL)</label><input v-model="form.logo" class="input" /></div>
        <div class="md:col-span-2"><label class="label">Тавсиф</label><textarea v-model="form.description" class="input min-h-[100px]"></textarea></div>
        <div><label class="label">Сомона</label><input v-model="form.website" class="input" /></div>
        <div><label class="label">Манзил</label><input v-model="form.address" class="input" /></div>
        <label class="flex items-center gap-2 text-sm md:col-span-2">
          <input v-model="form.is_active" type="checkbox" class="w-4 h-4 accent-brand-600" />
          Фаъол (дар саҳифаи ҷамъиятӣ намоиш дода шавад)
        </label>
      </div>
      <template #footer>
        <button class="btn-ghost" @click="modal = false">Бекор</button>
        <button class="btn-primary" @click="save">Захира</button>
      </template>
    </Modal>

    <ConfirmDialog v-model="confirmOpen" message="Бонк ва ҳама шартҳои он нест мешаванд" danger
                   :busy="busy" @confirm="confirmDel" />
  </div>
</template>

<script setup lang="ts">
import type { Bank } from '~/types/api'

const api = useApi()
const toast = useToast()

const banks = ref<Bank[]>([])
const loading = ref(true)
const modal = ref(false)
const editing = ref<Bank | null>(null)
const form = reactive<Partial<Bank>>({})

const confirmOpen = ref(false)
const busy = ref(false)
let pending: Bank | null = null

const load = async () => { loading.value = true; try { banks.value = await api.get<Bank[]>('/api/banks') } finally { loading.value = false } }
onMounted(load)

const openEditor = (b: Bank | null) => {
  editing.value = b
  Object.keys(form).forEach(k => delete (form as any)[k])
  if (b) Object.assign(form, b)
  else Object.assign(form, { is_active: true })
  modal.value = true
}

const save = async () => {
  if (!form.name?.trim()) return
  if (editing.value) await api.put(`/api/banks/${editing.value.id}`, form)
  else await api.post('/api/banks', form)
  modal.value = false; await load(); toast.success('Захира шуд')
}

const onDelete = (b: Bank) => { pending = b; confirmOpen.value = true }
const confirmDel = async () => {
  if (!pending) return
  busy.value = true
  try { await api.delete(`/api/banks/${pending.id}`); await load(); toast.success('Нест шуд') }
  finally { busy.value = false; confirmOpen.value = false; pending = null }
}
</script>
