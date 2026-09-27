<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Menu, ChevronDown, User, Lock, LogOut } from '@lucide/vue'
import { useAuthStore } from '@/stores/auth'

defineEmits<{
  'toggle-sidebar': []
}>()

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const isOpen = ref(false)

const pageTitle = computed(() => {
  return (route.meta.title as string) ?? 'Dashboard'
})

const admin = computed(() => ({
  name: authStore.admin?.name ?? 'Admin User',
  email: authStore.admin?.email ?? '',
  role: authStore.admin?.role ?? 'admin',
}))

/* Admin initials */
const initials = computed(() => {
  const name = admin.value.name?.trim()

  if (!name) {
    return 'AD'
  }

  const words = name.split(/\s+/)

  if (words.length >= 2) {
    return `${words[0][0]}${words[1][0]}`.toUpperCase()
  }

  return name.substring(0, 2).toUpperCase()
})

const displayRole = computed(() => {
  return admin.value.role?.toUpperCase() || 'ADMIN'
})

/* Toggle dropdown */
const toggleProfile = () => {
  isOpen.value = !isOpen.value
}

/* Close dropdown */
const closeProfile = () => {
  isOpen.value = false
}

/* Click outside */
const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement

  if (!target.closest('.admin-profile')) {
    closeProfile()
  }
}

/* Profile */
const goToProfile = () => {
  closeProfile()

  router.push('/admin/profile')
}

/* Change password */
const goToChangePassword = () => {
  closeProfile()

  router.push('/admin/change-password')
}

/* Logout */
const logout = () => {
  closeProfile()
  authStore.logoutAdmin()
  router.push('/admin/login')
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <header
    class="sticky top-0 z-30 flex h-[80px] items-center justify-between border-b border-white/10 bg-[#0d0d12] px-6 lg:px-8"
  >
    <!-- Left -->
    <div class="flex items-center gap-4">
      <!-- Mobile menu -->
      <button
        type="button"
        class="text-gray-400 transition hover:text-white lg:hidden"
        aria-label="เปิดเมนู"
        @click="$emit('toggle-sidebar')"
      >
        <Menu class="h-5 w-5" />
      </button>

      <!-- Page title -->
      <h1
        class="text-base font-semibold uppercase tracking-wide text-white"
      >
        {{ pageTitle }}
      </h1>
    </div>

    <!-- Admin Profile -->
    <div class="admin-profile relative">
      <!-- Profile Button -->
      <button
        type="button"
        class="flex items-center gap-3 rounded-xl px-2 py-1.5 transition hover:bg-white/5"
        @click.stop="toggleProfile"
      >
        <!-- Avatar -->
        <div
          class="flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-purple-500 to-pink-500 text-xs font-bold text-white shadow-lg shadow-purple-500/20"
        >
          {{ initials }}
        </div>

        <!-- Admin information -->
        <div class="hidden text-left sm:block">
          <p class="text-sm font-semibold text-white">
            {{ admin.name }}
          </p>

          <p class="text-[11px] uppercase tracking-wide text-gray-500">
            {{ displayRole }}
          </p>
        </div>

        <!-- Arrow -->
        <ChevronDown
          class="hidden h-4 w-4 text-gray-500 transition duration-200 sm:block"
          :class="{ 'rotate-180': isOpen }"
        />
      </button>

      <!-- Dropdown -->
      <Transition
        enter-active-class="transition duration-150 ease-out"
        enter-from-class="translate-y-1 opacity-0"
        enter-to-class="translate-y-0 opacity-100"
        leave-active-class="transition duration-100 ease-in"
        leave-from-class="translate-y-0 opacity-100"
        leave-to-class="translate-y-1 opacity-0"
      >
        <div
          v-if="isOpen"
          class="absolute right-0 top-[calc(100%+10px)] z-50 w-72 overflow-hidden rounded-xl border border-white/10 bg-[#15151c] shadow-2xl shadow-black/40"
          @click.stop
        >
          <!-- Admin information -->
          <div class="border-b border-white/10 p-4">
            <div class="flex items-center gap-3">
              <!-- Avatar -->
              <div
                class="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-purple-500 to-pink-500 text-sm font-bold text-white"
              >
                {{ initials }}
              </div>

              <!-- Name / Email -->
              <div class="min-w-0">
                <p class="truncate text-sm font-semibold text-white">
                  {{ admin.name }}
                </p>

                <p class="truncate text-xs text-gray-500">
                  {{ admin.email }}
                </p>

                <span
                  class="mt-1 inline-flex rounded-full bg-green-500/10 px-2 py-0.5 text-[10px] font-semibold text-green-400"
                >
                  {{ displayRole }}
                </span>
              </div>
            </div>
          </div>

          <!-- Menu -->
          <div class="p-2">
            <!-- Profile -->
            <button
              type="button"
              class="flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-sm text-gray-300 transition hover:bg-white/5 hover:text-white"
              @click="goToProfile"
            >
              <User class="h-4 w-4 text-gray-500" />

              <span>Profile</span>
            </button>

            <!-- Change Password -->
            <button
              type="button"
              class="flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-sm text-gray-300 transition hover:bg-white/5 hover:text-white"
              @click="goToChangePassword"
            >
              <Lock class="h-4 w-4 text-gray-500" />

              <span>Change Password</span>
            </button>
          </div>

          <!-- Logout -->
          <div class="border-t border-white/10 p-2">
            <button
              type="button"
              class="flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium text-red-400 transition hover:bg-red-500/10"
              @click="logout"
            >
              <LogOut class="h-4 w-4" />

              <span>Logout</span>
            </button>
          </div>
        </div>
      </Transition>
    </div>
  </header>
</template>