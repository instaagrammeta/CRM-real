<template>
  <div class="space-y-4">
    <PageHeader title="Постҳо (SMM)" description="Маркетинги шабакаҳои иҷтимоӣ">
      <template #actions>
        <button class="btn-primary" @click="openEditor(null)">
          <i class="fa-solid fa-plus"></i> Пости нав
        </button>
      </template>
    </PageHeader>

    <div v-if="loading" class="card flex justify-center py-12"><Spinner size="lg" class="text-brand-600" /></div>
    <EmptyState v-else-if="posts.length === 0" icon="fa-bullhorn" title="Пост нест" />
    <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3 lg:gap-4">
      <div v-for="p in posts" :key="p.id" class="card-hover cursor-pointer" @click="openEditor(p)">
        <div class="flex items-start justify-between gap-2 mb-2">
          <h3 class="font-semibold text-ink truncate">{{ p.title || 'Бе ном' }}</h3>
          <span class="badge-brand shrink-0">{{ p.category || '—' }}</span>
        </div>
        <p class="text-xs text-ink-soft line-clamp-3 mb-3">{{ p.description }}</p>
        <div class="flex items-center justify-between text-[11px] text-ink-mute">
          <span><i class="fa-regular fa-user"></i> {{ p.user_name }}</span>
          <span>{{ formatDate(p.post_date) || formatDate(p.created_date) }}</span>
        </div>
        <div class="grid grid-cols-4 gap-1.5 mt-3 text-center text-[11px]">
          <div class="bg-page rounded-lg py-1">
            <div class="font-bold text-ink">{{ p.likes }}</div>
            <div class="text-ink-mute"><i class="fa-regular fa-heart"></i></div>
          </div>
          <div class="bg-page rounded-lg py-1">
            <div class="font-bold text-ink">{{ p.comments }}</div>
            <div class="text-ink-mute"><i class="fa-regular fa-comment"></i></div>
          </div>
          <div class="bg-page rounded-lg py-1">
            <div class="font-bold text-ink">{{ p.views }}</div>
            <div class="text-ink-mute"><i class="fa-regular fa-eye"></i></div>
          </div>
          <div class="bg-page rounded-lg py-1">
            <div class="font-bold text-ink">{{ p.shares }}</div>
            <div class="text-ink-mute"><i class="fa-solid fa-share"></i></div>
          </div>
        </div>
      </div>
    </div>

    <Modal v-model="modal" :title="editing ? 'Тағйири пост' : 'Пости нав'" size="lg">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div class="md:col-span-2"><label class="label">Унвон</label><input v-model="form.title" class="input" /></div>
        <div><label class="label">Категория</label><input v-model="form.category" class="input" /></div>
        <div><label class="label">Намуди контент</label><input v-model="form.content_type" class="input" placeholder="reel, post, story..." /></div>
        <div><label class="label">Лоиҳа</label><input v-model="form.project" class="input" /></div>
        <div><label class="label">Санаи нашр</label><input v-model="form.post_date" type="date" class="input" /></div>
        <div class="md:col-span-2"><label class="label">Тавсиф</label><textarea v-model="form.description" class="input min-h-[100px]"></textarea></div>
        <div><label class="label">Likes</label><input v-model.number="form.likes" type="number" class="input" /></div>
        <div><label class="label">Comments</label><input v-model.number="form.comments" type="number" class="input" /></div>
        <div><label class="label">Views</label><input v-model.number="form.views" type="number" class="input" /></div>
        <div><label class="label">Shares</label><input v-model.number="form.shares" type="number" class="input" /></div>
        <div class="md:col-span-2"><label class="label">Линк</label><input v-model="form.link" class="input" /></div>
      </div>
      <template #footer>
        <button v-if="editing" class="btn-danger mr-auto" @click="onDelete">
          <i class="fa-solid fa-trash"></i> Нест
        </button>
        <button class="btn-ghost" @click="modal = false">Бекор</button>
        <button class="btn-primary" @click="save">Захира</button>
      </template>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import type { Post } from '~/types/api'
const api = useApi()
const toast = useToast()
const { formatDate } = useFormat()

const posts = ref<Post[]>([])
const loading = ref(true)
const modal = ref(false)
const editing = ref<Post | null>(null)
const form = reactive<Partial<Post>>({})

const load = async () => { loading.value = true; try { posts.value = await api.get<Post[]>('/api/posts') } finally { loading.value = false } }
onMounted(load)

const openEditor = (p: Post | null) => {
  editing.value = p
  Object.keys(form).forEach(k => delete (form as any)[k])
  if (p) Object.assign(form, p)
  modal.value = true
}

const save = async () => {
  if (editing.value) await api.put(`/api/posts/${editing.value.id}`, form)
  else await api.post('/api/posts', form)
  modal.value = false; await load(); toast.success('Захира шуд')
}

const onDelete = async () => {
  if (!editing.value) return
  await api.delete(`/api/posts/${editing.value.id}`)
  modal.value = false; await load(); toast.success('Нест шуд')
}
</script>
