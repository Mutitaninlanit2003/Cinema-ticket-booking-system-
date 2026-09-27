<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const errorMessage = ref('')

// ล้างข้อมูลในฟอร์มทุกครั้งที่เปิดหน้า Admin Login
onMounted(() => {
  email.value = ''
  password.value = ''
  errorMessage.value = ''
})

async function handleSubmit() {
  loading.value = true
  errorMessage.value = ''

  try {
    // Admin Login
    await authStore.loginAdmin(
      email.value.trim(),
      password.value,
    )

    // Login สำเร็จ → เข้า Admin Dashboard
    await router.push('/admin')
  } catch (err: unknown) {
    console.error('Admin login error:', err)

    errorMessage.value =
      err instanceof Error
        ? err.message
        : 'Login failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div
    class="flex min-h-screen items-center justify-center bg-[#0d0d12] px-4"
  >
    <div
      class="w-full max-w-md rounded-2xl border border-white/10 bg-[#15151d] p-8"
    >
      <h1 class="text-2xl font-bold text-white">
        Admin Login
      </h1>

      <p class="mt-1 text-sm text-[#8190ad]">
        Cinema Ticket Booking — Admin Panel
      </p>

      <form
        class="mt-6 space-y-4"
        autocomplete="off"
        @submit.prevent="handleSubmit"
      >
        <!-- Error -->
        <div
          v-if="errorMessage"
          class="rounded-lg border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400"
        >
          {{ errorMessage }}
        </div>

        <!-- Email -->
        <div>
          <label
            class="mb-1.5 block text-xs font-medium text-[#8190ad]"
          >
            Email
          </label>

          <input
            v-model="email"
            type="email"
            required
            autocomplete="off"
            class="w-full rounded-lg border border-white/10 bg-[#0d0d12] px-4 py-2.5 text-sm text-white outline-none focus:border-[#ff3158]"
          />
        </div>

        <!-- Password -->
        <div>
          <label
            class="mb-1.5 block text-xs font-medium text-[#8190ad]"
          >
            Password
          </label>

          <input
            v-model="password"
            type="password"
            required
            autocomplete="new-password"
            class="w-full rounded-lg border border-white/10 bg-[#0d0d12] px-4 py-2.5 text-sm text-white outline-none focus:border-[#ff3158]"
          />
        </div>

        <!-- Login Button -->
        <button
          type="submit"
          :disabled="loading"
          class="w-full rounded-lg bg-[#ff3158] py-2.5 text-sm font-semibold text-white transition hover:bg-[#e02a4d] disabled:opacity-50"
        >
          {{ loading ? 'Logging in...' : 'Log In' }}
        </button>
      </form>
    </div>
  </div>
</template>