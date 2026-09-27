<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const showPassword = ref(false)

// ==========================================
// Google OAuth
// ==========================================

// ลิงก์นี้พาไปที่ backend โดยตรง (ไม่ใช่ frontend route)
// backend จะ redirect ต่อไปยัง Google เอง
const googleLoginUrl =
  `${import.meta.env.VITE_API_BASE_URL}/auth/google`

async function handleLogin() {
  errorMessage.value = ''
  successMessage.value = ''

  if (!email.value.trim() || !password.value) {
    errorMessage.value = 'Please enter your email and password.'
    return
  }

  loading.value = true

  try {
    // ==========================================
    // User Login
    // ==========================================
    // ใช้ loginUser() โดยเฉพาะ
    // เพื่อเก็บ session ลง user_auth
    await authStore.loginUser(
      email.value.trim(),
      password.value,
    )

    successMessage.value = 'Login successful.'

    // กลับไปหน้า User Home
    await router.push('/')
  } catch (error: unknown) {
    console.error('User login error:', error)

    errorMessage.value =
      error instanceof Error
        ? error.message
        : 'Unable to login. Please check your email and password.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div
    class="flex min-h-screen items-center justify-center bg-[#0d0d10] px-4 py-8 text-white sm:px-6 sm:py-12"
  >
    <div class="w-full max-w-md">

      <!-- Logo / Brand -->
      <div class="mb-6 text-center sm:mb-8">
        <div
          class="mx-auto mb-3 flex h-12 w-12 items-center justify-center rounded-2xl bg-[#ff3158] shadow-lg shadow-[#ff3158]/20 sm:mb-4 sm:h-14 sm:w-14"
        >
          <span class="text-xl sm:text-2xl">🎬</span>
        </div>

        <h1
          class="px-2 text-2xl font-bold leading-tight text-white sm:text-3xl"
        >
          Cinema Ticket Booking
        </h1>

        <p class="mt-2 px-2 text-sm leading-5 text-gray-400">
          Book your favorite movies with us
        </p>
      </div>

      <!-- Login Card -->
      <div
        class="rounded-2xl border border-white/10 bg-[#15151d] p-5 shadow-2xl sm:p-8"
      >
        <!-- Heading -->
        <div class="mb-6 sm:mb-7">
          <h2 class="text-xl font-bold text-white sm:text-2xl">
            Welcome back
          </h2>

          <p class="mt-2 text-sm text-gray-400">
            Log in to your account
          </p>
        </div>

        <!-- Error -->
        <div
          v-if="errorMessage"
          class="mb-5 rounded-lg border border-red-500/20 bg-red-500/10 px-3 py-3 text-sm leading-5 text-red-400 break-words sm:px-4"
        >
          {{ errorMessage }}
        </div>

        <!-- Success -->
        <div
          v-if="successMessage"
          class="mb-5 rounded-lg border border-green-500/20 bg-green-500/10 px-3 py-3 text-sm leading-5 text-green-400 break-words sm:px-4"
        >
          {{ successMessage }}
        </div>

        <!-- Login Form -->
        <form
          @submit.prevent="handleLogin"
          class="space-y-5"
        >

          <!-- Email -->
          <div>
            <label
              for="email"
              class="mb-2 block text-sm font-medium text-gray-300"
            >
              Email
            </label>

            <input
              id="email"
              v-model="email"
              type="email"
              autocomplete="email"
              placeholder="Enter your email"
              class="w-full rounded-xl border border-white/10 bg-[#0d0d10] px-4 py-3 text-base text-white outline-none transition placeholder:text-gray-600 focus:border-[#ff3158] focus:ring-2 focus:ring-[#ff3158]/20 sm:text-sm"
            />
          </div>

          <!-- Password -->
          <div>
            <div
              class="mb-2 flex items-center justify-between gap-3"
            >
              <label
                for="password"
                class="block text-sm font-medium text-gray-300"
              >
                Password
              </label>

              <router-link
                to="/forgot-password"
                class="shrink-0 text-xs text-[#ff3158] transition hover:text-[#ff5273] sm:text-sm"
              >
                Forgot password?
              </router-link>
            </div>

            <div class="relative">
              <input
                id="password"
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                placeholder="Enter your password"
                class="w-full rounded-xl border border-white/10 bg-[#0d0d10] px-4 py-3 pr-16 text-base text-white outline-none transition placeholder:text-gray-600 focus:border-[#ff3158] focus:ring-2 focus:ring-[#ff3158]/20 sm:text-sm"
              />

              <button
                type="button"
                class="absolute right-3 top-1/2 -translate-y-1/2 px-1 text-sm text-gray-500 transition hover:text-white"
                @click="showPassword = !showPassword"
              >
                {{ showPassword ? 'Hide' : 'Show' }}
              </button>
            </div>
          </div>

          <!-- Login Button -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full rounded-xl bg-[#ff3158] py-3.5 text-sm font-semibold text-white transition hover:bg-[#ff1748] disabled:cursor-not-allowed disabled:opacity-60 sm:text-base"
          >
            <span v-if="!loading">
              Log In
            </span>

            <span v-else>
              Logging in...
            </span>
          </button>
        </form>

        <!-- Divider -->
        <div class="my-5 flex items-center gap-3 sm:my-6">
          <div class="h-px flex-1 bg-white/10"></div>

          <span class="shrink-0 text-xs text-gray-500">
            OR
          </span>

          <div class="h-px flex-1 bg-white/10"></div>
        </div>

        <!-- Google Login -->
        <a
          :href="googleLoginUrl"
          class="flex w-full items-center justify-center gap-2.5 rounded-xl border border-white/10 bg-white px-3 py-3.5 text-sm font-semibold text-gray-800 transition hover:bg-gray-100 sm:gap-3 sm:text-base"
        >
          <svg
            class="h-5 w-5 shrink-0"
            viewBox="0 0 24 24"
          >
            <path
              fill="#4285F4"
              d="M23.52 12.27c0-.85-.08-1.67-.22-2.45H12v4.64h6.47c-.28 1.5-1.13 2.77-2.4 3.62v3.01h3.89c2.28-2.1 3.56-5.2 3.56-8.82z"
            />

            <path
              fill="#34A853"
              d="M12 24c3.24 0 5.95-1.07 7.94-2.91l-3.89-3.01c-1.08.72-2.45 1.15-4.05 1.15-3.11 0-5.75-2.1-6.69-4.92H1.29v3.1C3.26 21.3 7.31 24 12 24z"
            />

            <path
              fill="#FBBC05"
              d="M5.31 14.31A7.2 7.2 0 0 1 4.93 12c0-.8.14-1.58.38-2.31v-3.1H1.29A11.96 11.96 0 0 0 0 12c0 1.93.46 3.76 1.29 5.41l4.02-3.1z"
            />

            <path
              fill="#EA4335"
              d="M12 4.77c1.76 0 3.35.6 4.59 1.79l3.44-3.44C17.94 1.19 15.24 0 12 0 7.31 0 3.26 2.7 1.29 6.59l4.02 3.1C6.25 6.87 8.89 4.77 12 4.77z"
            />
          </svg>

          <span class="truncate">
            Continue with Google
          </span>
        </a>

        <!-- Register -->
        <div
          class="mt-6 border-t border-white/10 pt-5 text-center sm:mt-7 sm:pt-6"
        >
          <p class="text-sm leading-6 text-gray-400">
            Don't have an account?

            <router-link
              to="/register"
              class="ml-1 font-medium text-[#ff3158] transition hover:text-[#ff5273]"
            >
              Create an account
            </router-link>
          </p>
        </div>
      </div>

      <!-- Back -->
      <div class="mt-5 text-center sm:mt-6">
        <router-link
          to="/"
          class="text-sm text-gray-500 transition hover:text-white"
        >
          ← Back to Home
        </router-link>
      </div>
    </div>
  </div>
</template>