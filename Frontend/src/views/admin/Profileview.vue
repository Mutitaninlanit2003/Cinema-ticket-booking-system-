<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import api from '@/services/api'
import { useAuthStore } from '@/stores/auth'
// import('@/views/admin/ProfileView.vue')
import {
  User,
  Mail,
  Lock,
  Camera,
  ShieldCheck,
  Loader2,
} from 'lucide-vue-next'

const authStore = useAuthStore()

// -----------------------------------------------------------------------
// Types — match backend models.User (json:"-" hides password_hash, so
// it never appears here)
// -----------------------------------------------------------------------
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
const loading = ref(false)
const loadError = ref('')

// Backend stores profile_image as a root-relative path
// ("/uploads/avatars/xxx.jpg") — the API and the frontend dev server
// run on different ports, so a plain relative path won't resolve
// against the frontend's own origin. Prefix it with the API base URL.
const avatarUrl = computed(() => {
  if (!profile.value?.profile_image) return null

  const base = import.meta.env.VITE_API_BASE_URL ?? ''
  return `${base}${profile.value.profile_image}`
})

const initials = computed(() => {
  const name = profile.value?.name?.trim()
  if (!name) return '?'

  return name
    .split(/\s+/)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join('')
})

async function loadProfile() {
  loading.value = true
  loadError.value = ''

  try {
    const { data } = await api.get<Profile>('/users/me')
    profile.value = data
    startEditingProfile()
  } catch (err) {
    loadError.value = 'Failed to load your profile. Please refresh and try again.'
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(loadProfile)

// -----------------------------------------------------------------------
// Edit name / email
// -----------------------------------------------------------------------
const editForm = ref({ name: '', email: '' })
const savingProfile = ref(false)
const profileError = ref('')
const profileSuccess = ref('')

function startEditingProfile() {
  if (!profile.value) return

  editForm.value = {
    name: profile.value.name,
    email: profile.value.email,
  }
  profileError.value = ''
  profileSuccess.value = ''
}

async function saveProfile() {
  if (!editForm.value.name || !editForm.value.email) return

  savingProfile.value = true
  profileError.value = ''
  profileSuccess.value = ''

  try {
    const { data } = await api.put<Profile>('/users/me', {
      name: editForm.value.name,
      email: editForm.value.email,
    })

    profile.value = data
    profileSuccess.value = 'Profile updated.'

    // Keep the auth store (header, sidebar, etc.) in sync immediately
    // instead of waiting for the next login.
    if (authStore.user) {
      authStore.user.name = data.name
      authStore.user.email = data.email
    }
  } catch (err: any) {
    profileError.value =
      err?.response?.data?.error ?? 'Failed to update profile. Please try again.'
    console.error(err)
  } finally {
    savingProfile.value = false
  }
}

// -----------------------------------------------------------------------
// Change password
// -----------------------------------------------------------------------
const passwordForm = ref({
  current_password: '',
  new_password: '',
  confirm_password: '',
})
const savingPassword = ref(false)
const passwordError = ref('')
const passwordSuccess = ref('')

async function changePassword() {
  passwordError.value = ''
  passwordSuccess.value = ''

  if (passwordForm.value.new_password.length < 8) {
    passwordError.value = 'New password must be at least 8 characters.'
    return
  }

  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
    passwordError.value = 'New password and confirmation do not match.'
    return
  }

  savingPassword.value = true

  try {
    await api.put('/users/me/password', {
      current_password: passwordForm.value.current_password,
      new_password: passwordForm.value.new_password,
    })

    passwordSuccess.value = 'Password changed successfully.'
    passwordForm.value = { current_password: '', new_password: '', confirm_password: '' }
  } catch (err: any) {
    passwordError.value =
      err?.response?.data?.error ?? 'Failed to change password. Please try again.'
    console.error(err)
  } finally {
    savingPassword.value = false
  }
}

// -----------------------------------------------------------------------
// Avatar upload
// -----------------------------------------------------------------------
const avatarInput = ref<HTMLInputElement | null>(null)
const uploadingAvatar = ref(false)
const avatarError = ref('')

function pickAvatarFile() {
  avatarInput.value?.click()
}

async function onAvatarSelected(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return

  avatarError.value = ''

  const allowed = ['image/jpeg', 'image/png', 'image/webp']
  if (!allowed.includes(file.type)) {
    avatarError.value = 'Please choose a .jpg, .png, or .webp image.'
    return
  }

  if (file.size > 5 * 1024 * 1024) {
    avatarError.value = 'Image must be 5MB or smaller.'
    return
  }

  const formData = new FormData()
  formData.append('file', file)

  uploadingAvatar.value = true

  try {
    const { data } = await api.post<{ profile_image: string }>(
      '/users/me/profile-image',
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      },
    )

    if (profile.value) {
      profile.value.profile_image = data.profile_image
    }

    if (authStore.user) {
      authStore.user.profile_image = data.profile_image
    }
  } catch (err: any) {
    avatarError.value =
      err?.response?.data?.error ??
      'Failed to upload image. Please try again.'

    console.error(err)
  } finally {
    uploadingAvatar.value = false

    if (avatarInput.value) {
      avatarInput.value.value = ''
    }
  }

  uploadingAvatar.value = true

  try {
    const { data } = await api.post<{ profile_image: string }>(
      '/users/me/avatar',
      formData,
      { headers: { 'Content-Type': 'multipart/form-data' } },
    )

    if (profile.value) {
      profile.value.profile_image = data.profile_image
    }
  } catch (err: any) {
    avatarError.value =
      err?.response?.data?.error ?? 'Failed to upload image. Please try again.'
    console.error(err)
  } finally {
    uploadingAvatar.value = false

    // Reset so selecting the same file again still fires @change.
    if (avatarInput.value) avatarInput.value.value = ''
  }
}
</script>

<template>
  <div class="mx-auto max-w-4xl px-6 py-8 lg:px-8">

    <!-- Page header -->
    <div class="mb-7">
      <h1 class="text-2xl font-bold text-white">My Profile</h1>
      <p class="mt-1 text-sm text-gray-400">
        Manage your account details, password, and photo
      </p>
    </div>

    <!-- Loading -->
    <div
      v-if="loading"
      class="flex items-center justify-center gap-2 rounded-2xl border border-white/10 bg-[#15151d] py-16 text-sm text-gray-400"
    >
      <Loader2 class="size-4 animate-spin" />
      Loading your profile...
    </div>

    <!-- Load error -->
    <div
      v-else-if="loadError"
      class="rounded-2xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400"
    >
      {{ loadError }}
    </div>

    <div v-else-if="profile" class="space-y-6">

      <!-- Identity card -->
      <div class="rounded-2xl border border-white/10 bg-[#15151d] p-6">
        <div class="flex flex-col items-center gap-5 sm:flex-row">

          <!-- Avatar -->
          <div class="relative shrink-0">
            <div
              class="flex h-24 w-24 items-center justify-center overflow-hidden rounded-full border border-white/10 bg-[#0d0d12] text-2xl font-bold text-gray-400"
            >
              <img
                v-if="avatarUrl"
                :src="avatarUrl"
                alt="Profile photo"
                class="h-full w-full object-cover"
              />
              <span v-else>{{ initials }}</span>
            </div>

            <button
              type="button"
              class="absolute -bottom-1 -right-1 flex h-8 w-8 items-center justify-center rounded-full border border-white/10 bg-[#ff3158] text-white transition hover:bg-[#ff4669] disabled:opacity-50"
              title="Change photo"
              :disabled="uploadingAvatar"
              @click="pickAvatarFile"
            >
              <Loader2 v-if="uploadingAvatar" class="size-4 animate-spin" />
              <Camera v-else class="size-4" />
            </button>

            <input
              ref="avatarInput"
              type="file"
              accept="image/jpeg,image/png,image/webp"
              class="hidden"
              @change="onAvatarSelected"
            />
          </div>

          <!-- Name / email / role -->
          <div class="min-w-0 flex-1 text-center sm:text-left">
            <p class="truncate text-lg font-bold text-white">
              {{ profile.name }}
            </p>

            <p class="truncate text-sm text-gray-400">
              {{ profile.email }}
            </p>

            <span
              class="mt-2 inline-flex items-center gap-1.5 rounded-full bg-[#ff3158]/15 px-3 py-1 text-xs font-semibold text-[#ff3158]"
            >
              <ShieldCheck class="size-3.5" />
              {{ profile.role }}
            </span>
          </div>

        </div>

        <p v-if="avatarError" class="mt-4 text-xs text-red-400">
          {{ avatarError }}
        </p>
      </div>

      <!-- Edit name / email -->
      <div class="rounded-2xl border border-white/10 bg-[#15151d] p-6">
        <div class="mb-5 flex items-center gap-2">
          <User class="size-4 text-gray-400" />
          <h2 class="text-sm font-semibold text-white">Profile information</h2>
        </div>

        <div class="grid gap-4 sm:grid-cols-2">
          <div>
            <label class="mb-1.5 block text-xs font-medium text-gray-400">Name</label>
            <input
              v-model="editForm.name"
              type="text"
              class="h-11 w-full rounded-xl border border-white/10 bg-[#0d0d12] px-4 text-sm text-white outline-none transition focus:border-[#ff3158]/60"
            />
          </div>

          <div>
            <label class="mb-1.5 block text-xs font-medium text-gray-400">Email</label>
            <div class="relative">
              <Mail class="absolute left-3.5 top-1/2 size-4 -translate-y-1/2 text-gray-500" />
              <input
                v-model="editForm.email"
                type="email"
                class="h-11 w-full rounded-xl border border-white/10 bg-[#0d0d12] pl-10 pr-4 text-sm text-white outline-none transition focus:border-[#ff3158]/60"
              />
            </div>
          </div>
        </div>

        <p v-if="profileError" class="mt-4 text-xs text-red-400">{{ profileError }}</p>
        <p v-if="profileSuccess" class="mt-4 text-xs text-emerald-400">{{ profileSuccess }}</p>

        <div class="mt-5 flex justify-end">
          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-xl bg-[#ff3158] px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-[#ff4669] disabled:opacity-50"
            :disabled="savingProfile"
            @click="saveProfile"
          >
            <Loader2 v-if="savingProfile" class="size-4 animate-spin" />
            Save changes
          </button>
        </div>
      </div>

      <!-- Change password -->
      <div class="rounded-2xl border border-white/10 bg-[#15151d] p-6">
        <div class="mb-5 flex items-center gap-2">
          <Lock class="size-4 text-gray-400" />
          <h2 class="text-sm font-semibold text-white">Change password</h2>
        </div>

        <div class="grid gap-4 sm:grid-cols-3">
          <div>
            <label class="mb-1.5 block text-xs font-medium text-gray-400">Current password</label>
            <input
              v-model="passwordForm.current_password"
              type="password"
              autocomplete="current-password"
              class="h-11 w-full rounded-xl border border-white/10 bg-[#0d0d12] px-4 text-sm text-white outline-none transition focus:border-[#ff3158]/60"
            />
          </div>

          <div>
            <label class="mb-1.5 block text-xs font-medium text-gray-400">New password</label>
            <input
              v-model="passwordForm.new_password"
              type="password"
              autocomplete="new-password"
              class="h-11 w-full rounded-xl border border-white/10 bg-[#0d0d12] px-4 text-sm text-white outline-none transition focus:border-[#ff3158]/60"
            />
          </div>

          <div>
            <label class="mb-1.5 block text-xs font-medium text-gray-400">Confirm new password</label>
            <input
              v-model="passwordForm.confirm_password"
              type="password"
              autocomplete="new-password"
              class="h-11 w-full rounded-xl border border-white/10 bg-[#0d0d12] px-4 text-sm text-white outline-none transition focus:border-[#ff3158]/60"
            />
          </div>
        </div>

        <p class="mt-3 text-xs text-gray-500">At least 8 characters.</p>

        <p v-if="passwordError" class="mt-2 text-xs text-red-400">{{ passwordError }}</p>
        <p v-if="passwordSuccess" class="mt-2 text-xs text-emerald-400">{{ passwordSuccess }}</p>

        <div class="mt-5 flex justify-end">
          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-xl border border-white/10 bg-transparent px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-white/5 disabled:opacity-50"
            :disabled="
              savingPassword ||
              !passwordForm.current_password ||
              !passwordForm.new_password ||
              !passwordForm.confirm_password
            "
            @click="changePassword"
          >
            <Loader2 v-if="savingPassword" class="size-4 animate-spin" />
            Update password
          </button>
        </div>
      </div>

    </div>

  </div>
</template>