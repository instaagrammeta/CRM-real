<template>
  <div class="w-full max-w-md">
    <div class="text-center mb-8">
      <div class="inline-flex w-14 h-14 rounded-2xl bg-brand-100 items-center justify-center mb-3">
        <i class="fa-solid fa-house-chimney text-2xl text-brand-700"></i>
      </div>
      <h1 class="text-2xl font-bold text-ink">Real Estate <span class="text-brand-600">CRM</span></h1>
      <p class="text-sm text-ink-soft mt-1">Хуш омадед! Барои идома ворид шавед.</p>
    </div>

    <form class="card space-y-4" @submit.prevent="submit">
      <div>
        <label class="label">Логин</label>
        <div class="relative">
          <i class="fa-solid fa-user absolute left-4 top-1/2 -translate-y-1/2 text-ink-mute text-sm"></i>
          <input v-model="form.login" class="input pl-10" placeholder="admin" autocomplete="username" required />
        </div>
      </div>

      <div>
        <label class="label">Парол</label>
        <div class="relative">
          <i class="fa-solid fa-lock absolute left-4 top-1/2 -translate-y-1/2 text-ink-mute text-sm"></i>
          <input
            v-model="form.password"
            :type="showPwd ? 'text' : 'password'"
            class="input pl-10 pr-10"
            placeholder="••••••••"
            autocomplete="current-password"
            required
          />
          <button type="button" class="absolute right-3 top-1/2 -translate-y-1/2 text-ink-mute hover:text-ink"
                  @click="showPwd = !showPwd">
            <i class="fa-solid" :class="showPwd ? 'fa-eye-slash' : 'fa-eye'"></i>
          </button>
        </div>
      </div>

      <div v-if="error" class="bg-red-50 text-red-700 text-sm rounded-xl px-4 py-2.5 flex items-start gap-2">
        <i class="fa-solid fa-circle-exclamation mt-0.5"></i>
        <span>{{ error }}</span>
      </div>

      <button type="submit" class="btn-primary w-full !py-3" :disabled="busy">
        <Spinner v-if="busy" size="sm" />
        <span v-else>Воридшавӣ</span>
      </button>

      <div class="text-center pt-2 border-t border-border/50 -mx-5 -mb-5 px-5 pb-3">
        <NuxtLink to="/ipoteka" class="text-xs text-ink-mute hover:text-brand-700">Ипотека</NuxtLink>
        <span class="text-ink-mute mx-2">·</span>
        <NuxtLink to="/rasrochka" class="text-xs text-ink-mute hover:text-brand-700">Рассрочка</NuxtLink>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'auth' })

const auth = useAuthStore()
const route = useRoute()
const toast = useToast()

const form = reactive({ login: '', password: '' })
const showPwd = ref(false)
const busy = ref(false)
const error = ref('')

const submit = async () => {
  busy.value = true
  error.value = ''
  try {
    await auth.login(form.login.trim(), form.password)
    toast.success('Хуш омадед!', auth.fullName)
    const redirect = (route.query.redirect as string) || '/'
    await navigateTo(redirect)
  } catch (e: any) {
    error.value = e?.data?.error || e?.message || 'Воридшавӣ ноком шуд'
  } finally {
    busy.value = false
  }
}
</script>
