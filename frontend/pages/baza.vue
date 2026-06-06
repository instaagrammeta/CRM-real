<template>
  <div class="space-y-4">
    <PageHeader title="База файлҳо" description="Папкаҳо ва файлҳо">
      <template #actions>
        <button class="btn-soft" @click="newFolderModal = true">
          <i class="fa-solid fa-folder-plus"></i> Папка
        </button>
        <label class="btn-primary cursor-pointer">
          <i class="fa-solid fa-upload"></i> Илова кардани файл
          <input type="file" class="hidden" multiple @change="onUpload" />
        </label>
      </template>
    </PageHeader>

    <!-- Breadcrumb -->
    <div class="flex items-center gap-2 text-sm">
      <button class="hover:text-brand-700" @click="goTo(0, [])">
        <i class="fa-solid fa-home"></i> Решаи
      </button>
      <template v-for="(b, i) in breadcrumbs" :key="b.id">
        <i class="fa-solid fa-chevron-right text-xs text-ink-mute"></i>
        <button class="hover:text-brand-700" @click="goTo(b.id, breadcrumbs.slice(0, i + 1))">
          {{ b.name }}
        </button>
      </template>
    </div>

    <div v-if="loading" class="card flex justify-center py-12"><Spinner size="lg" class="text-brand-600" /></div>

    <div v-else class="space-y-3">
      <!-- Folders -->
      <div v-if="folders.length > 0" class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3">
        <div
          v-for="f in folders"
          :key="f.id"
          class="card-hover cursor-pointer text-center !p-4"
          @click="enterFolder(f)"
        >
          <i class="fa-solid fa-folder text-3xl text-brand-600 mb-2"></i>
          <div class="text-sm font-medium truncate">{{ f.name }}</div>
        </div>
      </div>

      <!-- Files -->
      <div v-if="files.length > 0" class="card !p-0 overflow-hidden">
        <table class="w-full text-sm">
          <thead class="bg-page text-ink-mute uppercase text-[11px]">
            <tr>
              <th class="text-left px-4 py-2.5">Ном</th>
              <th class="text-left px-4 py-2.5">Намуд</th>
              <th class="text-left px-4 py-2.5">Андоза</th>
              <th class="text-left px-4 py-2.5">Муаллиф</th>
              <th class="text-right px-4 py-2.5 w-24">Амалҳо</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="f in files" :key="f.id" class="border-t border-border/40 hover:bg-page/60">
              <td class="px-4 py-2.5 truncate"><i class="fa-regular fa-file mr-2 text-brand-600"></i> {{ f.original_name }}</td>
              <td class="px-4 py-2.5 text-xs text-ink-mute">{{ f.filetype || '—' }}</td>
              <td class="px-4 py-2.5 text-xs">{{ humanSize(f.filesize) }}</td>
              <td class="px-4 py-2.5 text-xs">{{ f.author_name }}</td>
              <td class="px-4 py-2.5 text-right">
                <a :href="downloadUrl(f.id)" class="text-brand-700 hover:underline text-xs mr-3">
                  <i class="fa-solid fa-download"></i>
                </a>
                <button class="text-ink-mute hover:text-red-600 text-xs" @click="deleteFile(f.id)">
                  <i class="fa-solid fa-trash"></i>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <EmptyState v-if="folders.length === 0 && files.length === 0" icon="fa-folder-open" title="Холӣ" />
    </div>

    <Modal v-model="newFolderModal" title="Папкаи нав">
      <input v-model="newFolderName" class="input" placeholder="Номи папка" />
      <template #footer>
        <button class="btn-ghost" @click="newFolderModal = false">Бекор</button>
        <button class="btn-primary" @click="createFolder">Сохтан</button>
      </template>
    </Modal>
  </div>
</template>

<script setup lang="ts">
const api = useApi()
const toast = useToast()
const config = useRuntimeConfig()

const currentFolder = ref(0)
const breadcrumbs = ref<{ id: number; name: string }[]>([])
const folders = ref<any[]>([])
const files = ref<any[]>([])
const loading = ref(true)

const newFolderModal = ref(false)
const newFolderName = ref('')

const downloadUrl = (id: number) => `${config.public.apiBase}/api/download/${id}`

const humanSize = (s?: number) => {
  if (!s) return '0 B'
  const u = ['B', 'KB', 'MB', 'GB']
  let i = 0; let n = s
  while (n >= 1024 && i < u.length - 1) { n /= 1024; i++ }
  return n.toFixed(1) + ' ' + u[i]
}

const load = async () => {
  loading.value = true
  try {
    folders.value = await api.get<any[]>('/api/folders', { parent_id: currentFolder.value })
    if (currentFolder.value > 0) {
      files.value = await api.get<any[]>(`/api/folders/${currentFolder.value}/files`)
    } else {
      files.value = []
    }
  } catch (e: any) {
    toast.error('Хато', e?.message || '')
  } finally { loading.value = false }
}

const enterFolder = (f: any) => {
  currentFolder.value = f.id
  breadcrumbs.value.push({ id: f.id, name: f.name })
  load()
}

const goTo = (id: number, crumbs: any[]) => {
  currentFolder.value = id
  breadcrumbs.value = crumbs
  load()
}

const createFolder = async () => {
  if (!newFolderName.value.trim()) return
  await api.post('/api/folders', { name: newFolderName.value, parent_id: currentFolder.value })
  newFolderModal.value = false
  newFolderName.value = ''
  await load()
  toast.success('Папка сохта шуд')
}

const onUpload = async (ev: Event) => {
  if (!currentFolder.value) {
    toast.warning('Огоҳӣ', 'Аввал ба папка ворид шавед')
    return
  }
  const input = ev.target as HTMLInputElement
  const files = Array.from(input.files || [])
  for (const file of files) {
    const fd = new FormData()
    fd.append('file', file)
    await api.upload(`/api/folders/${currentFolder.value}/files`, fd)
  }
  input.value = ''
  await load()
  toast.success(`${files.length} файл бор шуд`)
}

const deleteFile = async (id: number) => {
  await api.delete(`/api/files/${id}`)
  await load()
  toast.success('Файл нест шуд')
}

onMounted(load)
</script>
