<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const errorMessage = ref('')

onMounted(async () => {
  const token = route.query.token as string | undefined

  if (!token) {
    errorMessage.value =
      'Missing login token. Please try logging in again.'
    return
  }

  try {
    // ----------------------------------------------------
    // เก็บ token ชั่วคราวไว้ก่อน เพื่อให้ api.ts (interceptor)
    // แนบ Authorization header ตอนเรียก /users/me ด้านล่าง
    // ----------------------------------------------------

    sessionStorage.setItem(
      'user_auth',
      JSON.stringify({ token, user: null }),
    )

    // ----------------------------------------------------
    // ดึงข้อมูลโปรไฟล์ตัวเองด้วย token ที่เพิ่งได้มา
    // ----------------------------------------------------

    const { data } = await api.get('/users/me')

    authStore.setUserSession(token, {
      user_id: data.user_id,
      name: data.name,
      email: data.email,
      role: data.role,
      profile_image: data.profile_image,
    })

    await router.replace('/')
  } catch (err) {
    console.error('Google login callback failed:', err)

    sessionStorage.removeItem('user_auth')

    errorMessage.value =
      'Google login failed. Please try again.'
  }
})
</script>

<template>
  <div
    class="flex min-h-screen items-center justify-center bg-[#0d0d10] px-6"
  >
    <div class="text-center">

      <template v-if="!errorMessage">
        <div
          class="mx-auto mb-4 h-10 w-10 animate-spin rounded-full border-4 border-white/10 border-t-[#ff3158]"
        ></div>

        <p class="text-sm text-gray-400">
          Signing you in with Google…
        </p>
      </template>

      <template v-else>
        <p class="mb-4 text-sm text-red-400">
          {{ errorMessage }}
        </p>

        <router-link
          to="/login"
          class="text-sm font-medium text-[#ff3158] hover:text-[#ff5273]"
        >
          Back to login
        </router-link>
      </template>

    </div>
  </div>
</template>