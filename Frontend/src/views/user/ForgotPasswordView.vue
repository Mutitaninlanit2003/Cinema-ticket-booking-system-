<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../../services/api'

const router = useRouter()

const email = ref('')

const loading = ref(false)

const errorMessage = ref('')

const successMessage = ref('')

async function handleForgotPassword() {
  errorMessage.value = ''
  successMessage.value = ''

  if (!email.value.trim()) {
    errorMessage.value = 'Please enter your email address.'
    return
  }

  loading.value = true

  try {
    const response = await api.post('/auth/forgot-password', {
      email: email.value.trim(),
    })

    successMessage.value =
      response.data?.message ||
      'If the email exists, a password reset link has been created.'

    // Get reset token from backend
    const resetToken = response.data?.reset_token

    // Go to Reset Password page
    if (resetToken) {
      setTimeout(() => {
        router.push({
          path: '/reset-password',
          query: {
            token: resetToken,
          },
        })
      }, 1000)
    }
  } catch (error: any) {
    console.error('Forgot password error:', error)

    if (error.response?.data?.error) {
      errorMessage.value = error.response.data.error
    } else {
      errorMessage.value =
        'Unable to process your request. Please try again.'
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

      <!-- Forgot Password Card -->
      <div
        class="rounded-2xl border border-white/10 bg-[#15151d] p-8 shadow-2xl"
      >

        <!-- Heading -->
        <div class="mb-7">
          <h2 class="text-2xl font-bold text-white">
            Forgot password?
          </h2>

          <p class="mt-2 text-sm leading-6 text-gray-400">
            Enter the email address associated with your account and we'll
            help you reset your password.
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
          @submit.prevent="handleForgotPassword"
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
              class="w-full rounded-xl border border-white/10 bg-[#0d0d10] px-4 py-3 text-white outline-none transition placeholder:text-gray-600 focus:border-[#ff3158] focus:ring-2 focus:ring-[#ff3158]/20"
            />
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full rounded-xl bg-[#ff3158] py-3.5 font-semibold text-white transition hover:bg-[#ff1748] disabled:cursor-not-allowed disabled:opacity-60"
          >
            <span v-if="!loading">
              Send Reset Link
            </span>

            <span v-else>
              Sending...
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