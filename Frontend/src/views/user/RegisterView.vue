<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../../services/api'

const router = useRouter()

const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')

const loading = ref(false)
const googleLoading = ref(false)

const errorMessage = ref('')
const successMessage = ref('')

const showPassword = ref(false)
const showConfirmPassword = ref(false)

async function handleRegister() {
  errorMessage.value = ''
  successMessage.value = ''

  // ตรวจสอบข้อมูลเบื้องต้น
  if (
    !name.value.trim() ||
    !email.value.trim() ||
    !password.value ||
    !confirmPassword.value
  ) {
    errorMessage.value = 'Please fill in all fields.'
    return
  }

  // ตรวจสอบ password
  if (password.value.length < 8) {
    errorMessage.value = 'Password must be at least 8 characters.'
    return
  }

  // ตรวจสอบ password ตรงกัน
  if (password.value !== confirmPassword.value) {
    errorMessage.value = 'Passwords do not match.'
    return
  }

  loading.value = true

  try {
    const response = await api.post('/auth/register', {
      name: name.value.trim(),
      email: email.value.trim(),
      password: password.value,
    })

    const data = response.data

    // เก็บ JWT Token
    if (data.token) {
      localStorage.setItem(
        'user_auth',
        JSON.stringify({
          token: data.token,
          user: data.user ?? null,
        }),
      )
    }

    // เก็บข้อมูล User
    if (data.user) {
      localStorage.setItem('user', JSON.stringify(data.user))
    }

    successMessage.value = 'Registration successful.'

    // รอเล็กน้อยเพื่อให้ผู้ใช้เห็นข้อความสำเร็จ
    setTimeout(() => {
      router.push('/')
    }, 800)
  } catch (error: any) {
    console.error('Register error:', error)

    if (error.response?.data?.error) {
      errorMessage.value = error.response.data.error
    } else {
      errorMessage.value =
        'Unable to create your account. Please try again.'
    }
  } finally {
    loading.value = false
  }
}

// --------------------------------------------------
// Google Sign Up
// --------------------------------------------------

function signUpWithGoogle() {
  errorMessage.value = ''
  successMessage.value = ''
  googleLoading.value = true

  // เริ่ม Google OAuth flow
  window.location.href =
    `${import.meta.env.VITE_API_BASE_URL}/auth/google`
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

      <!-- Register Card -->
      <div
        class="rounded-2xl border border-white/10 bg-[#15151d] p-5 shadow-2xl sm:p-8"
      >

        <!-- Heading -->
        <div class="mb-6 sm:mb-7">
          <h2 class="text-xl font-bold text-white sm:text-2xl">
            Create an account
          </h2>

          <p class="mt-2 text-sm text-gray-400">
            Sign up to book your favorite movies
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

        <!-- Google Sign Up -->
        <button
          type="button"
          :disabled="loading || googleLoading"
          @click="signUpWithGoogle"
          class="flex w-full items-center justify-center gap-2.5 rounded-xl border border-white/10 bg-white px-3 py-3.5 text-sm font-semibold text-gray-800 transition hover:bg-gray-100 disabled:cursor-not-allowed disabled:opacity-60 sm:gap-3 sm:text-base"
        >
          <!-- Google Logo -->
          <svg
            v-if="!googleLoading"
            class="h-5 w-5 shrink-0"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fill="#4285F4"
              d="M21.35 12.27c0-.79-.07-1.55-.2-2.27H12v4.3h5.24a4.48 4.48 0 0 1-1.94 2.94v2.45h3.14c1.84-1.69 2.91-4.18 2.91-7.42z"
            />

            <path
              fill="#34A853"
              d="M12 21.5c2.63 0 4.84-.87 6.45-2.36l-3.14-2.45c-.87.58-1.98.93-3.31.93-2.54 0-4.7-1.72-5.47-4.03H3.28v2.53A9.74 9.74 0 0 0 12 21.5z"
            />

            <path
              fill="#FBBC05"
              d="M6.53 13.59A5.85 5.85 0 0 1 6.23 12c0-.55.1-1.09.3-1.59V7.88H3.28A9.75 9.75 0 0 0 2.25 12c0 1.57.38 3.05 1.03 4.12l3.25-2.53z"
            />

            <path
              fill="#EA4335"
              d="M12 6.38c1.43 0 2.72.49 3.73 1.45l2.8-2.8C16.84 3.43 14.63 2.5 12 2.5a9.74 9.74 0 0 0-8.72 5.38l3.25 2.53C7.3 8.1 9.46 6.38 12 6.38z"
            />
          </svg>

          <span
            v-if="!googleLoading"
            class="truncate"
          >
            Sign up with Google
          </span>

          <span v-else>
            Connecting to Google...
          </span>
        </button>

        <!-- Divider -->
        <div class="my-5 flex items-center gap-3 sm:my-6 sm:gap-4">
          <div class="h-px flex-1 bg-white/10"></div>

          <span class="shrink-0 text-xs text-gray-500">
            OR
          </span>

          <div class="h-px flex-1 bg-white/10"></div>
        </div>

        <!-- Register Form -->
        <form
          @submit.prevent="handleRegister"
          class="space-y-5"
        >

          <!-- Name -->
          <div>
            <label
              for="name"
              class="mb-2 block text-sm font-medium text-gray-300"
            >
              Name
            </label>

            <input
              id="name"
              v-model="name"
              type="text"
              autocomplete="name"
              placeholder="Enter your name"
              class="w-full rounded-xl border border-white/10 bg-[#0d0d10] px-4 py-3 text-base text-white outline-none transition placeholder:text-gray-600 focus:border-[#ff3158] focus:ring-2 focus:ring-[#ff3158]/20 sm:text-sm"
            />
          </div>

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
            <label
              for="password"
              class="mb-2 block text-sm font-medium text-gray-300"
            >
              Password
            </label>

            <div class="relative">
              <input
                id="password"
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="new-password"
                placeholder="Create a password"
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

            <p class="mt-2 text-xs leading-5 text-gray-500">
              Password must be at least 8 characters.
            </p>
          </div>

          <!-- Confirm Password -->
          <div>
            <label
              for="confirmPassword"
              class="mb-2 block text-sm font-medium text-gray-300"
            >
              Confirm password
            </label>

            <div class="relative">
              <input
                id="confirmPassword"
                v-model="confirmPassword"
                :type="showConfirmPassword ? 'text' : 'password'"
                autocomplete="new-password"
                placeholder="Confirm your password"
                class="w-full rounded-xl border border-white/10 bg-[#0d0d10] px-4 py-3 pr-16 text-base text-white outline-none transition placeholder:text-gray-600 focus:border-[#ff3158] focus:ring-2 focus:ring-[#ff3158]/20 sm:text-sm"
              />

              <button
                type="button"
                class="absolute right-3 top-1/2 -translate-y-1/2 px-1 text-sm text-gray-500 transition hover:text-white"
                @click="showConfirmPassword = !showConfirmPassword"
              >
                {{ showConfirmPassword ? 'Hide' : 'Show' }}
              </button>
            </div>
          </div>

          <!-- Register Button -->
          <button
            type="submit"
            :disabled="loading || googleLoading"
            class="w-full rounded-xl bg-[#ff3158] py-3.5 text-sm font-semibold text-white transition hover:bg-[#ff1748] disabled:cursor-not-allowed disabled:opacity-60 sm:text-base"
          >
            <span v-if="!loading">
              Create Account
            </span>

            <span v-else>
              Creating account...
            </span>
          </button>
        </form>

        <!-- Login -->
        <div
          class="mt-6 border-t border-white/10 pt-5 text-center sm:mt-7 sm:pt-6"
        >
          <p class="text-sm leading-6 text-gray-400">
            Already have an account?

            <router-link
              to="/login"
              class="ml-1 font-medium text-[#ff3158] transition hover:text-[#ff5273]"
            >
              Log in
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