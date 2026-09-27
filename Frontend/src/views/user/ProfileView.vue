<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

import {
  User,
  Mail,
  Lock,
  Save,
  Eye,
  EyeOff,
  Camera,
} from 'lucide-vue-next'

import api from '@/services/api'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const loading = ref(true)
const savingProfile = ref(false)
const changingPassword = ref(false)
const uploadingAvatar = ref(false)

const errorMessage = ref('')
const successMessage = ref('')

const name = ref('')
const email = ref('')

const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')

const showCurrentPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)

const userInitial = ref('U')

interface Profile {
  id: number
  user_id: string
  name: string
  email: string
  role: 'user' | 'admin'
  profile_image?: string
  created_at: string
  updated_at: string
}

const profile = ref<Profile | null>(null)

/*
 * Profile page is for USER only.
 * Do not use admin_auth here.
 */
const isUserLoggedIn = computed(() => {
  return authStore.isUserAuthenticated
})

/*
 * Build the profile image URL.
 *
 * Two cases:
 * 1. Google login → profile_image is already a full URL
 *    (https://lh3.googleusercontent.com/...) → use as-is.
 * 2. Uploaded avatar → profile_image is a short backend path
 *    (/uploads/profiles/xxx.jpg) → prefix with the backend's
 *    base URL: http://localhost:8080/uploads/...
 */
const avatarUrl = computed(() => {
  const imagePath = profile.value?.profile_image

  if (!imagePath) {
    return null
  }

  if (
    imagePath.startsWith('http://') ||
    imagePath.startsWith('https://')
  ) {
    return imagePath
  }

  const apiBase =
    import.meta.env.VITE_API_BASE_URL ||
    'http://localhost:8080'

  const base = apiBase.endsWith('/api')
    ? apiBase.slice(0, -4)
    : apiBase

  return `${base}${imagePath}`
})

async function loadProfile() {
  try {
    // โหลดข้อมูล User จาก sessionStorage ก่อน
    authStore.initFromStorage()

    // เรียกข้อมูล User ที่ Login อยู่
    const { data } = await api.get<Profile>('/users/me')

    // Backend ส่ง Profile โดยตรง
    profile.value = data

    // เติมข้อมูลลงในช่องแก้ไขโปรไฟล์
    name.value = data.name || authStore.user?.name || ''
    email.value = data.email || authStore.user?.email || ''

    // อัปเดตตัวอักษรใน Avatar
    userInitial.value =
      name.value.charAt(0).toUpperCase() || 'U'

    // อัปเดตข้อมูลใน Pinia User session
    if (authStore.user) {
      authStore.user = {
        ...authStore.user,
        name: data.name || authStore.user.name,
        email: data.email || authStore.user.email,
        profile_image:
          data.profile_image || authStore.user.profile_image,
      }
    }
  } catch (error: unknown) {
    console.error('Load profile error:', error)

    // ถ้า API โหลดไม่ได้ ให้ใช้ข้อมูลจาก User session
    if (authStore.user) {
      name.value = authStore.user.name || ''
      email.value = authStore.user.email || ''

      userInitial.value =
        name.value.charAt(0).toUpperCase() || 'U'

      profile.value = {
        id: 0,
        user_id: authStore.user.user_id,
        name: authStore.user.name,
        email: authStore.user.email,
        role: authStore.user.role,
        profile_image: authStore.user.profile_image,
        created_at: '',
        updated_at: '',
      }
    }
  } finally {
    // ปิดสถานะ Loading เสมอ
    loading.value = false
  }
}

async function saveProfile() {
  errorMessage.value = ''
  successMessage.value = ''

  if (!name.value.trim()) {
    errorMessage.value = 'Please enter your name.'
    return
  }

  if (!email.value.trim()) {
    errorMessage.value = 'Please enter your email.'
    return
  }

  if (!authStore.isUserAuthenticated) {
    errorMessage.value =
      'Please log in to update your profile.'
    return
  }

  savingProfile.value = true

  try {
    const response = await api.put('/users/me', {
      name: name.value.trim(),
      email: email.value.trim(),
    })

    const updatedUser = response.data?.user || response.data

    if (updatedUser) {
      name.value = updatedUser.name || name.value
      email.value = updatedUser.email || email.value

      userInitial.value = name.value
        .charAt(0)
        .toUpperCase()

      if (profile.value) {
        profile.value.name = name.value
        profile.value.email = email.value
      }

      /*
       * Update USER Pinia state.
       */
      if (authStore.user) {
        authStore.user = {
          ...authStore.user,
          name: name.value,
          email: email.value,
        }
      }

      /*
       * Persist the updated USER session
       * back to sessionStorage.
       */
      const userAuth = sessionStorage.getItem('user_auth')

      if (userAuth) {
        try {
          const parsed = JSON.parse(userAuth)

          parsed.user = {
            ...parsed.user,
            name: name.value,
            email: email.value,
          }

          sessionStorage.setItem(
            'user_auth',
            JSON.stringify(parsed),
          )
        } catch {
          console.warn(
            'Unable to update user_auth in sessionStorage.',
          )
        }
      }
    }

    successMessage.value =
      'Profile updated successfully.'
  } catch (error: unknown) {
    console.error('Update profile error:', error)

    if (error instanceof Error) {
      errorMessage.value = error.message
    } else {
      errorMessage.value =
        'Unable to update your profile. Please try again.'
    }
  } finally {
    savingProfile.value = false
  }
}

async function changePassword() {
  errorMessage.value = ''
  successMessage.value = ''

  if (!currentPassword.value) {
    errorMessage.value =
      'Please enter your current password.'
    return
  }

  if (!newPassword.value) {
    errorMessage.value =
      'Please enter your new password.'
    return
  }

  if (newPassword.value.length < 8) {
    errorMessage.value =
      'New password must be at least 8 characters.'
    return
  }

  if (newPassword.value !== confirmPassword.value) {
    errorMessage.value =
      'New password and confirmation password do not match.'
    return
  }

  if (!authStore.isUserAuthenticated) {
    errorMessage.value =
      'Please log in to change your password.'
    return
  }

  changingPassword.value = true

  try {
    await api.put('/users/me/password', {
      current_password: currentPassword.value,
      new_password: newPassword.value,
    })

    currentPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''

    successMessage.value =
      'Password changed successfully.'
  } catch (error: unknown) {
    console.error('Change password error:', error)

    if (error instanceof Error) {
      errorMessage.value = error.message
    } else {
      errorMessage.value =
        'Unable to change your password. Please try again.'
    }
  } finally {
    changingPassword.value = false
  }
}

async function uploadAvatar(event: Event) {
  const target = event.target as HTMLInputElement

  if (!target.files || target.files.length === 0) {
    return
  }

  const file = target.files[0]

  errorMessage.value = ''
  successMessage.value = ''

  const allowedTypes = [
    'image/jpeg',
    'image/png',
    'image/webp',
  ]

  if (!allowedTypes.includes(file.type)) {
    errorMessage.value =
      'Please select a JPG, PNG, or WebP image.'

    target.value = ''
    return
  }

  if (file.size > 5 * 1024 * 1024) {
    errorMessage.value =
      'Profile image must be smaller than 5 MB.'

    target.value = ''
    return
  }

  if (!authStore.isUserAuthenticated) {
    errorMessage.value =
      'Please log in to upload a profile image.'

    target.value = ''
    return
  }

  const formData = new FormData()

  /*
   * Backend expects:
   * multipart field name = "file"
   */
  formData.append('file', file)

  uploadingAvatar.value = true

  try {
    const response = await api.post<{
      profile_image: string
    }>(
      '/users/me/profile-image',
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      },
    )

    const imagePath = response.data.profile_image

    if (profile.value) {
      profile.value.profile_image = imagePath
    }

    /*
     * Update USER Pinia state.
     */
    if (authStore.user) {
      authStore.user = {
        ...authStore.user,
        profile_image: imagePath,
      }
    }

    /*
     * Persist profile image into USER session.
     */
    const userAuth = sessionStorage.getItem('user_auth')

    if (userAuth) {
      try {
        const parsed = JSON.parse(userAuth)

        parsed.user = {
          ...parsed.user,
          profile_image: imagePath,
        }

        sessionStorage.setItem(
          'user_auth',
          JSON.stringify(parsed),
        )
      } catch {
        console.warn(
          'Unable to update user_auth in sessionStorage.',
        )
      }
    }

    successMessage.value =
      'Profile image uploaded successfully.'
  } catch (error: unknown) {
    console.error(
      'Upload profile image error:',
      error,
    )

    if (error instanceof Error) {
      errorMessage.value = error.message
    } else {
      errorMessage.value =
        'Unable to upload your profile image. Please try again.'
    }
  } finally {
    uploadingAvatar.value = false

    /*
     * Allow selecting the same file again.
     */
    target.value = ''
  }
}

onMounted(() => {
  loadProfile()
})
</script>

<template>
  <div
    class="min-h-screen bg-[#0b0b0f] px-4 py-8 sm:px-6 sm:py-12"
  >
    <div class="mx-auto w-full max-w-5xl">

      <!-- Page Header -->
      <div class="mb-6 sm:mb-10">
        <h1
          class="text-2xl font-bold text-white sm:text-3xl"
        >
          My Profile
        </h1>

        <p
          class="mt-2 max-w-xl text-sm leading-6 text-gray-400 sm:text-base"
        >
          Manage your account information and password.
        </p>
      </div>

      <!-- Loading -->
      <div
        v-if="loading"
        class="rounded-2xl border border-white/10 bg-[#15151d] p-6 text-center sm:p-10"
      >
        <p class="text-sm text-gray-400 sm:text-base">
          Loading your profile...
        </p>
      </div>

      <div v-else>

        <!-- Error -->
        <div
          v-if="errorMessage"
          class="mb-4 break-words rounded-xl border border-red-500/20 bg-red-500/10 px-4 py-3 text-sm leading-6 text-red-400 sm:mb-6 sm:px-5 sm:py-4"
        >
          {{ errorMessage }}
        </div>

        <!-- Success -->
        <div
          v-if="successMessage"
          class="mb-4 break-words rounded-xl border border-green-500/20 bg-green-500/10 px-4 py-3 text-sm leading-6 text-green-400 sm:mb-6 sm:px-5 sm:py-4"
        >
          {{ successMessage }}
        </div>

        <!-- Profile Card -->
        <section
          class="overflow-hidden rounded-2xl border border-white/10 bg-[#15151d] shadow-2xl"
        >

          <!-- Profile Header -->
          <div
            class="border-b border-white/10 px-5 py-6 sm:px-8 sm:py-8"
          >
            <div
              class="flex flex-col items-center gap-4 text-center sm:flex-row sm:items-center sm:gap-5 sm:text-left"
            >

              <!-- Avatar -->
              <div class="relative shrink-0">

                <div
                  class="flex h-20 w-20 items-center justify-center overflow-hidden rounded-full bg-gradient-to-br from-[#ff3158] to-purple-500 text-2xl font-bold text-white shadow-lg shadow-[#ff3158]/20 sm:h-24 sm:w-24 sm:text-3xl"
                >
                  <img
                    v-if="avatarUrl"
                    :src="avatarUrl"
                    alt="Profile"
                    class="h-full w-full object-cover"
                  />

                  <span v-else>
                    {{ userInitial }}
                  </span>
                </div>

                <!-- Upload Button -->
                <label
                  class="absolute bottom-0 right-0 flex h-8 w-8 cursor-pointer items-center justify-center rounded-full border-2 border-[#15151d] bg-[#ff3158] text-white transition hover:bg-[#ff4d6d] sm:h-7 sm:w-7"
                  title="Change profile image"
                >
                  <Camera :size="14" />

                  <input
                    type="file"
                    accept="image/jpeg,image/png,image/webp"
                    class="hidden"
                    :disabled="uploadingAvatar"
                    @change="uploadAvatar"
                  />
                </label>
              </div>

              <!-- Profile Info -->
              <div class="min-w-0 max-w-full">

                <h2
                  class="truncate text-lg font-bold text-white sm:text-xl"
                >
                  {{ name || 'User' }}
                </h2>

                <p
                  class="mt-1 break-all text-sm leading-5 text-gray-400"
                >
                  {{ email }}
                </p>

                <p
                  v-if="uploadingAvatar"
                  class="mt-2 text-xs text-gray-500"
                >
                  Uploading image...
                </p>

              </div>
            </div>
          </div>

          <!-- Personal Information -->
          <div class="p-5 sm:p-8">

            <div class="mb-5 sm:mb-6">
              <h3
                class="text-lg font-semibold text-white"
              >
                Personal Information
              </h3>

              <p
                class="mt-1 text-sm leading-6 text-gray-400"
              >
                Update your basic account information.
              </p>
            </div>

            <!-- Name + Email -->
            <div
              class="grid gap-5 md:grid-cols-2 md:gap-6"
            >

              <!-- Name -->
              <div class="min-w-0">
                <label
                  for="name"
                  class="mb-2 block text-sm font-medium text-gray-300"
                >
                  Name
                </label>

                <div class="relative">

                  <User
                    :size="19"
                    class="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-gray-500"
                  />

                  <input
                    id="name"
                    v-model="name"
                    type="text"
                    placeholder="Enter your name"
                    class="w-full min-w-0 rounded-xl border border-white/10 bg-[#0b0b0f] py-3.5 pl-12 pr-4 text-base text-white outline-none transition placeholder:text-gray-600 focus:border-[#ff3158] focus:ring-2 focus:ring-[#ff3158]/20"
                  />

                </div>
              </div>

              <!-- Email -->
              <div class="min-w-0">
                <label
                  for="email"
                  class="mb-2 block text-sm font-medium text-gray-300"
                >
                  Email
                </label>

                <div class="relative">

                  <Mail
                    :size="19"
                    class="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-gray-500"
                  />

                  <input
                    id="email"
                    v-model="email"
                    type="email"
                    placeholder="Enter your email"
                    class="w-full min-w-0 rounded-xl border border-white/10 bg-[#0b0b0f] py-3.5 pl-12 pr-4 text-base text-white outline-none transition placeholder:text-gray-600 focus:border-[#ff3158] focus:ring-2 focus:ring-[#ff3158]/20"
                  />

                </div>
              </div>

            </div>

            <!-- Save -->
            <div
              class="mt-6 flex justify-stretch sm:mt-7 sm:justify-end"
            >

              <button
                type="button"
                @click="saveProfile"
                :disabled="savingProfile"
                class="flex w-full items-center justify-center gap-2 rounded-xl bg-gradient-to-r from-[#ff2d55] to-[#ff4d6d] px-6 py-3.5 font-semibold text-white shadow-lg shadow-[#ff2d55]/20 transition hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-60 sm:w-auto sm:py-3"
              >

                <Save :size="18" />

                <span v-if="!savingProfile">
                  Save Changes
                </span>

                <span v-else>
                  Saving...
                </span>

              </button>

            </div>
          </div>
        </section>

        <!-- Change Password -->
        <section
          class="mt-6 rounded-2xl border border-white/10 bg-[#15151d] p-5 shadow-2xl sm:mt-8 sm:p-8"
        >

          <div class="mb-5 sm:mb-6">

            <h2
              class="text-lg font-semibold text-white"
            >
              Change Password
            </h2>

            <p
              class="mt-1 text-sm leading-6 text-gray-400"
            >
              Keep your account secure by using a strong password.
            </p>

          </div>

          <div class="space-y-5">

            <!-- Current Password -->
            <div class="min-w-0">

              <label
                for="current-password"
                class="mb-2 block text-sm font-medium text-gray-300"
              >
                Current Password
              </label>

              <div class="relative">

                <Lock
                  :size="19"
                  class="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-gray-500"
                />

                <input
                  id="current-password"
                  v-model="currentPassword"
                  :type="
                    showCurrentPassword
                      ? 'text'
                      : 'password'
                  "
                  placeholder="Enter current password"
                  class="w-full min-w-0 rounded-xl border border-white/10 bg-[#0b0b0f] py-3.5 pl-12 pr-12 text-base text-white outline-none transition placeholder:text-gray-600 focus:border-[#ff3158] focus:ring-2 focus:ring-[#ff3158]/20"
                />

                <button
                  type="button"
                  @click="
                    showCurrentPassword =
                      !showCurrentPassword
                  "
                  class="absolute right-3 top-1/2 flex h-10 w-10 -translate-y-1/2 items-center justify-center text-gray-500 transition hover:text-white sm:right-2"
                  aria-label="Toggle current password visibility"
                >

                  <EyeOff
                    v-if="showCurrentPassword"
                    :size="19"
                  />

                  <Eye
                    v-else
                    :size="19"
                  />

                </button>

              </div>
            </div>

            <!-- New Password -->
            <div class="min-w-0">

              <label
                for="new-password"
                class="mb-2 block text-sm font-medium text-gray-300"
              >
                New Password
              </label>

              <div class="relative">

                <Lock
                  :size="19"
                  class="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-gray-500"
                />

                <input
                  id="new-password"
                  v-model="newPassword"
                  :type="
                    showNewPassword
                      ? 'text'
                      : 'password'
                  "
                  placeholder="Enter new password"
                  class="w-full min-w-0 rounded-xl border border-white/10 bg-[#0b0b0f] py-3.5 pl-12 pr-12 text-base text-white outline-none transition placeholder:text-gray-600 focus:border-[#ff3158] focus:ring-2 focus:ring-[#ff3158]/20"
                />

                <button
                  type="button"
                  @click="
                    showNewPassword =
                      !showNewPassword
                  "
                  class="absolute right-3 top-1/2 flex h-10 w-10 -translate-y-1/2 items-center justify-center text-gray-500 transition hover:text-white sm:right-2"
                  aria-label="Toggle new password visibility"
                >

                  <EyeOff
                    v-if="showNewPassword"
                    :size="19"
                  />

                  <Eye
                    v-else
                    :size="19"
                  />

                </button>

              </div>
            </div>

            <!-- Confirm Password -->
            <div class="min-w-0">

              <label
                for="confirm-password"
                class="mb-2 block text-sm font-medium text-gray-300"
              >
                Confirm New Password
              </label>

              <div class="relative">

                <Lock
                  :size="19"
                  class="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-gray-500"
                />

                <input
                  id="confirm-password"
                  v-model="confirmPassword"
                  :type="
                    showConfirmPassword
                      ? 'text'
                      : 'password'
                  "
                  placeholder="Confirm new password"
                  class="w-full min-w-0 rounded-xl border border-white/10 bg-[#0b0b0f] py-3.5 pl-12 pr-12 text-base text-white outline-none transition placeholder:text-gray-600 focus:border-[#ff3158] focus:ring-2 focus:ring-[#ff3158]/20"
                />

                <button
                  type="button"
                  @click="
                    showConfirmPassword =
                      !showConfirmPassword
                  "
                  class="absolute right-3 top-1/2 flex h-10 w-10 -translate-y-1/2 items-center justify-center text-gray-500 transition hover:text-white sm:right-2"
                  aria-label="Toggle confirm password visibility"
                >

                  <EyeOff
                    v-if="showConfirmPassword"
                    :size="19"
                  />

                  <Eye
                    v-else
                    :size="19"
                  />

                </button>

              </div>
            </div>

          </div>

          <!-- Change Password Button -->
          <div
            class="mt-6 flex justify-stretch sm:mt-7 sm:justify-end"
          >

            <button
              type="button"
              @click="changePassword"
              :disabled="changingPassword"
              class="w-full rounded-xl border border-white/10 bg-[#0b0b0f] px-6 py-3.5 font-semibold text-white transition hover:border-[#ff3158]/50 hover:bg-[#1c1c26] disabled:cursor-not-allowed disabled:opacity-60 sm:w-auto sm:py-3"
            >

              <span v-if="!changingPassword">
                Change Password
              </span>

              <span v-else>
                Changing...
              </span>

            </button>

          </div>
        </section>

      </div>
    </div>
  </div>
</template>