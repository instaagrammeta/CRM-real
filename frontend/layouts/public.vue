<template>
  <div class="min-h-screen bg-page">
    <header class="bg-white border-b border-border/60">
      <div class="max-w-6xl mx-auto px-4 lg:px-6 py-4 flex items-center justify-between">
        <NuxtLink to="/" class="flex items-center gap-2">
          <div class="w-9 h-9 rounded-xl bg-brand-100 grid place-items-center">
            <i class="fa-solid fa-house-chimney text-brand-700"></i>
          </div>
          <span class="font-bold text-ink">Real Estate <span class="text-brand-600">CRM</span></span>
        </NuxtLink>
        <nav class="flex items-center gap-1 text-sm">
          <NuxtLink to="/ipoteka"   class="nav-link" active-class="nav-link-active">Ипотека</NuxtLink>
          <NuxtLink to="/rasrochka" class="nav-link" active-class="nav-link-active">Рассрочка</NuxtLink>
          <NuxtLink to="/login"     class="btn-soft ml-2">Воридшавӣ</NuxtLink>
        </nav>
      </div>
    </header>

    <main class="max-w-6xl mx-auto px-4 lg:px-6 py-6 lg:py-10">
      <slot />
    </main>

    <footer class="text-center text-xs text-ink-mute py-6">
      © {{ new Date().getFullYear() }} Real Estate CRM. Ҳамаи ҳуқуқ ҳифз карда шудааст.
    </footer>
  </div>
</template>
