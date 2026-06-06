<template>
  <div class="space-y-4">
    <PageHeader title="Кормандон" description="Идоракунии корбарон ва ҳуқуқҳо">
      <template #actions>
        <input v-model="search" class="input !py-2 !w-64 hidden sm:block" placeholder="Ҷустуҷӯ..." />
        <button class="btn-primary" @click="openEditor(null)">
          <i class="fa-solid fa-plus"></i> Корманди нав
        </button>
      </template>
    </PageHeader>

    <div v-if="loading" class="card flex justify-center py-12"><Spinner size="lg" class="text-brand-600" /></div>

    <div v-else class="card !p-0 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-page text-ink-mute uppercase text-[11px]">
          <tr>
            <th class="text-left px-4 py-3">Корманд</th>
            <th class="text-left px-4 py-3">Логин</th>
            <th class="text-left px-4 py-3">Категория</th>
            <th class="text-left px-4 py-3">Нақш</th>
            <th class="text-right px-4 py-3 w-32">Амалҳо</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in filtered" :key="u.id" class="border-t border-border/40 hover:bg-page/60">
            <td class="px-4 py-3 flex items-center gap-3">
              <Avatar :name="u.full_name" :photo="u.photo" size="sm" />
              <div>
                <div class="font-medium">{{ u.full_name }}</div>
                <div class="text-xs text-ink-mute">{{ (u.personal_phones || []).join(', ') }}</div>
              </div>
            </td>
            <td class="px-4 py-3 text-ink-soft">{{ u.login }}</td>
            <td class="px-4 py-3 text-ink-soft">{{ u.category }}</td>
            <td class="px-4 py-3">
              <span :class="u.role === 'admin' ? 'badge-brand' : 'badge-gray'">{{ u.role }}</span>
            </td>
            <td class="px-4 py-3 text-right space-x-2">
              <button class="text-brand-700 hover:underline text-xs" @click="openEditor(u)">
                <i class="fa-solid fa-pen"></i>
              </button>
              <button class="text-red-600 hover:underline text-xs" @click="onDelete(u)" :disabled="u.id === me">
                <i class="fa-solid fa-trash"></i>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      <EmptyState v-if="filtered.length === 0" icon="fa-users" title="Корманд нест" />
    </div>

    <Modal v-model="modal" :title="editing ? 'Тағйири корманд' : 'Корманди нав'" size="lg">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div class="md:col-span-2"><label class="label">Номи пурра *</label><input v-model="form.full_name" class="input" /></div>
        <div><label class="label">Логин *</label><input v-model="form.login" class="input" /></div>
        <div><label class="label">{{ editing ? 'Парол (агар тағйир дода шавад)' : 'Парол *' }}</label>
          <input v-model="form.password" type="password" class="input" />
        </div>
        <div><label class="label">Синну сол</label><input v-model.number="form.age" type="number" class="input" /></div>
        <div>
          <label class="label">Нақш</label>
          <select v-model="form.role" class="input">
            <option value="employee">Корманд</option>
            <option value="admin">Админ</option>
          </select>
        </div>
        <div class="md:col-span-2"><label class="label">Категория</label><input v-model="form.category" class="input" placeholder="Менеҷер, SMM, агент..." /></div>
        <div class="md:col-span-2"><label class="label">Тел. шахсӣ (бо вергул)</label><input v-model="form.personal_phones_csv" class="input" /></div>
        <div class="md:col-span-2"><label class="label">Тел. кор (бо вергул)</label><input v-model="form.work_phones_csv" class="input" /></div>
      </div>
      <template #footer>
        <button class="btn-ghost" @click="modal = false">Бекор</button>
        <button class="btn-primary" @click="save">Захира</button>
      </template>
    </Modal>

    <ConfirmDialog v-model="confirmOpen" message="Корманд нест карда мешавад" danger
                   :busy="busy" @confirm="confirmDel" />
  </div>
</template>

<script setup lang="ts">
import type { User } from '~/types/api'

definePageMeta({ middleware: [] }) // global middleware already enforces /admin/*

const api = useApi()
const auth = useAuthStore()
const toast = useToast()

const users = ref<(User & any)[]>([])
const loading = ref(true)
const search = ref('')
const me = computed(() => auth.user?.id || 0)

const modal = ref(false)
const editing = ref<User | null>(null)
const form = reactive<any>({})

const confirmOpen = ref(false)
const busy = ref(false)
let pending: User | null = null

const load = async () => { loading.value = true; try { users.value = await api.get<User[]>('/api/users') } finally { loading.value = false } }
onMounted(load)

const filtered = computed(() => {
  const s = search.value.trim().toLowerCase()
  if (!s) return users.value
  return users.value.filter(u =>
    [u.full_name, u.login, u.category, u.role].some(v => (v || '').toLowerCase().includes(s)),
  )
})

const openEditor = (u: User | null) => {
  editing.value = u
  Object.keys(form).forEach(k => delete form[k])
  if (u) {
    Object.assign(form, u, {
      password: '',
      personal_phones_csv: (u.personal_phones || []).join(', '),
      work_phones_csv: (u.work_phones || []).join(', '),
    })
  } else {
    Object.assign(form, { full_name: '', login: '', password: '', role: 'employee', category: '', age: 0,
      personal_phones_csv: '', work_phones_csv: '' })
  }
  modal.value = true
}

const save = async () => {
  if (!form.full_name?.trim() || !form.login?.trim()) return
  if (!editing.value && !form.password) {
    toast.warning('Парол лозим аст'); return
  }
  const payload: any = {
    full_name: form.full_name,
    login: form.login,
    role: form.role,
    category: form.category,
    age: form.age || 0,
    personal_phones: (form.personal_phones_csv || '').split(',').map((s: string) => s.trim()).filter(Boolean),
    work_phones: (form.work_phones_csv || '').split(',').map((s: string) => s.trim()).filter(Boolean),
  }
  if (form.password) payload.password = form.password

  if (editing.value) await api.put(`/api/users/${editing.value.id}`, payload)
  else await api.post('/api/users', payload)
  modal.value = false
  await load()
  toast.success('Захира шуд')
}

const onDelete = (u: User) => { pending = u; confirmOpen.value = true }
const confirmDel = async () => {
  if (!pending) return
  busy.value = true
  try { await api.delete(`/api/users/${pending.id}`); await load(); toast.success('Нест шуд') }
  finally { busy.value = false; confirmOpen.value = false; pending = null }
}
</script>
