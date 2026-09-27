<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../../services/api'

const route = useRoute()
const router = useRouter()

const token = ref((route.query.token as string) || '')
const newPassword = ref('')
const confirmPassword = ref('')

const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const showPassword = ref(false)
const showConfirmPassword = ref(false)

async function handleResetPassword() {
  errorMessage.value = ''
  successMessage.value = ''

  if (!token.value) {
    errorMessage.value = 'Invalid or missing reset token.'
    return
  }

  if (!newPassword.value) {
    errorMessage.value = 'Please enter your new password.'
    return
  }

  if (newPassword.value.length < 8) {
    errorMessage.value = 'Password must be at least 8 characters.'
    return
  }

  if (!confirmPassword.value) {
    errorMessage.value = 'Please confirm your new password.'
    return
  }

  if (newPassword.value !== confirmPassword.value) {
    errorMessage.value = 'Passwords do not match.'
    return
  }

  loading.value = true

  try {
    const response = await api.post('/auth/reset-password', {
      token: token.value,
      new_password: newPassword.value,
    })

    successMessage.value =
      response.data?.message || 'Password reset successful.'

    newPassword.value = ''
    confirmPassword.value = ''

    setTimeout(() => {
      router.push('/login')
    }, 1500)
  } catch (error: any) {
    console.error('Reset password error:', error)

    if (error.response?.data?.error) {
      errorMessage.value = error.response.data.error
    } else {
      errorMessage.value =
        'Unable to reset your password. Please try again.'
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div
    class="min-h-screen bg-[#0d0d10] flex items-center justify-center px-6 py-12"
  >
    <div class="w-full max-w-md">

      <!-- Logo / Brand -->
      <div class="text-center mb-8">
        <div
          class="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-[#ff3158] shadow-lg shadow-[#ff3158]/20"
        >
          <span class="text-2xl">🎬</span>
        </div>

        <h1 class="text-3xl font-bold text-white">
          Cinema Ticket Booking
        </h1>

        <p class="mt-2 text-sm text-gray-400">
          Book your favorite movies with us
        </p>
      </div>

      <!-- Reset Password Card -->
      <div
        class="rounded-2xl border border-white/10 bg-[#15151d] p-8 shadow-2xl"
      >

        <!-- Heading -->
        <div class="mb-7">
          <h2 class="text-2xl font-bold text-white">
            Reset password
          </h2>

          <p class="mt-2 text-sm leading-6 text-gray-400">
            Enter a new password for your account.
          </p>
        </div>

        <!-- Error -->
        <div
          v-if="errorMessage"
          class="mb-5 rounded-lg border border-red-500/20 bg-red-500/10 px-4 py-3 text-sm text-red-400"
        >
          {{ errorMessage }}
        </div>

        <!-- Success -->
        <div
          v-if="successMessage"
          class="mb-5 rounded-lg border border-green-500/20 bg-green-500/10 px-4 py-3 text-sm text-green-400"
        >
          {{ successMessage }}
        </div>

        <!-- Form -->
        <form
          @submit.prevent="handleResetPassword"
          class="space-y-5"
        >

          <!-- New Password -->
          <div>
            <label
              for="new-password"
              class="mb-2 block text-sm font-medium text-gray-300"
            >
              New Password
            </label>

            <div class="relative">
              <input
                id="new-password"
                v-model="newPassword"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="new-password"
                placeholder="Enter your new password"
                class="w-full rounded-xl border border-white/10 bg-[#0d0d10] px-4 py-3 pr-12 text-white outline-none transition placeholder:text-gray-600 focus:border-[#ff3158] focus:ring-2 focus:ring-[#ff3158]/20"
              />

              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-sm text-gray-500 transition hover:text-white"
              >
                {{ showPassword ? 'Hide' : 'Show' }}
              </button>
            </div>
          </div>

          <!-- Confirm Password -->
          <div>
            <label
              for="confirm-password"
              class="mb-2 block text-sm font-medium text-gray-300"
            >
              Confirm Password
            </label>

            <div class="relative">
              <input
                id="confirm-password"
                v-model="confirmPassword"
                :type="showConfirmPassword ? 'text' : 'password'"
                autocomplete="new-password"
                placeholder="Confirm your new password"
                class="w-full rounded-xl border border-white/10 bg-[#0d0d10] px-4 py-3 pr-12 text-white outline-none transition placeholder:text-gray-600 focus:border-[#ff3158] focus:ring-2 focus:ring-[#ff3158]/20"
              />

              <button
                type="button"
                @click="showConfirmPassword = !showConfirmPassword"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-sm text-gray-500 transition hover:text-white"
              >
                {{ showConfirmPassword ? 'Hide' : 'Show' }}
              </button>
            </div>
          </div>

          <!-- Submit -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full rounded-xl bg-[#ff3158] py-3.5 font-semibold text-white transition hover:bg-[#ff1748] disabled:cursor-not-allowed disabled:opacity-60"
          >
            <span v-if="!loading">
              Reset Password
            </span>

            <span v-else>
              Resetting...
            </span>
          </button>
        </form>

        <!-- Back to Login -->
        <div
          class="mt-7 border-t border-white/10 pt-6 text-center"
        >
          <p class="text-sm text-gray-400">
            Remember your password?

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
      <div class="mt-6 text-center">
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