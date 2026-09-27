<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import axios from 'axios'
import {
  Search,
  Plus,
  Pencil,
  Trash2,
  MonitorPlay,
  X,
  Monitor,
} from 'lucide-vue-next'

// -----------------------------------------------------------------------
// API client
// NOTE: replace this with your project's central API client if you have
// one already (e.g. `import api from '@/api/client'`). This is a plain
// axios instance using VITE_API_BASE_URL as a fallback.
// -----------------------------------------------------------------------
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
})

// -----------------------------------------------------------------------
// Types — match backend models exactly (models/screen.go, models/theater.go)
// -----------------------------------------------------------------------
interface Screen {
  id: number
  screen_id: string
  theater_id: string
  name: string
  total_seats: number
  created_at: string
  updated_at: string
}

interface Theater {
  id: number
  theater_id: string
  name: string
  address: string
  status: string
  location: string
  capacity: number
}

const search = ref('')
const selectedTheaterId = ref('All')
const showForm = ref(false)
const editingScreen = ref<Screen | null>(null)

const screens = ref<Screen[]>([])
const theaters = ref<Theater[]>([])

const loading = ref(false)
const errorMessage = ref('')

// Errors from saving are shown inside the modal itself, because the
// page-level banner sits behind the overlay while the modal is open.
const formError = ref('')
const saving = ref(false)
const nameInput = ref<HTMLInputElement | null>(null)

const form = ref({
  name: '',
  theater_id: '',
  total_seats: 100,
})

// screen_id -> Theater lookup, since GET /screens does not return
// theater details (Theater field is json:"-" on the backend)
const theaterById = computed(() => {
  const map = new Map<string, Theater>()
  theaters.value.forEach((t) => map.set(t.theater_id, t))
  return map
})

function theaterName(theaterId: string) {
  return theaterById.value.get(theaterId)?.name ?? 'Unknown theater'
}

const theaterFilterOptions = computed(() => {
  return ['All', ...theaters.value.map((t) => t.theater_id)]
})

const filteredScreens = computed(() => {
  const keyword = search.value.toLowerCase()

  return screens.value.filter((screen) => {
    const matchesSearch =
      screen.name.toLowerCase().includes(keyword) ||
      theaterName(screen.theater_id).toLowerCase().includes(keyword)

    const matchesTheater =
      selectedTheaterId.value === 'All' ||
      screen.theater_id === selectedTheaterId.value

    return matchesSearch && matchesTheater
  })
})

// -----------------------------------------------------------------------
// Data loading
// -----------------------------------------------------------------------
async function loadTheaters() {
  const { data } = await api.get<Theater[]>('/theaters')
  theaters.value = data
}

async function loadScreens() {
  const { data } = await api.get<Screen[]>('/screens')
  screens.value = data
}

async function loadAll() {
  loading.value = true
  errorMessage.value = ''

  try {
    await Promise.all([loadTheaters(), loadScreens()])
  } catch (err) {
    errorMessage.value = 'Failed to load screens. Please try again.'
    console.error(err)
  } finally {
    loading.value = false
  }
}

// -----------------------------------------------------------------------
// Modal behaviour: Esc closes it, page behind stops scrolling, and the
// name field is focused as soon as it opens.
// -----------------------------------------------------------------------
function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape' && showForm.value) {
    closeForm()
  }
}

watch(showForm, async (open) => {
  document.body.style.overflow = open ? 'hidden' : ''

  if (open) {
    await nextTick()
    nameInput.value?.focus()
  }
})

onMounted(() => {
  loadAll()
  window.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
  document.body.style.overflow = ''
})

// -----------------------------------------------------------------------
// Form handling
// -----------------------------------------------------------------------
function openAddForm() {
  editingScreen.value = null
  formError.value = ''

  form.value = {
    name: '',
    theater_id: theaters.value[0]?.theater_id ?? '',
    total_seats: 100,
  }

  showForm.value = true
}

function openEditForm(screen: Screen) {
  editingScreen.value = screen
  formError.value = ''

  form.value = {
    name: screen.name,
    theater_id: screen.theater_id,
    total_seats: screen.total_seats,
  }

  showForm.value = true
}

function closeForm() {
  showForm.value = false
  editingScreen.value = null
  formError.value = ''
}

async function saveScreen() {
  if (saving.value) return

  if (!form.value.name.trim()) {
    formError.value = 'Enter a screen name.'
    return
  }

  if (!form.value.theater_id) {
    formError.value = 'Select a theater.'
    return
  }

  if (!form.value.total_seats || form.value.total_seats < 1) {
    formError.value = 'Total seats must be at least 1.'
    return
  }

  formError.value = ''
  errorMessage.value = ''
  saving.value = true

  try {
    if (editingScreen.value) {
      // Backend's UpdateScreen repository only PERSISTS name + total_seats
      // (theater_id changes are ignored there) — but the controller still
      // runs full struct validation via ShouldBindJSON first, and Screen's
      // TheaterID field is tagged `binding:"required,uuid"`. Omitting it
      // makes validation fail with 400 before the request even reaches
      // the repository. So we still send the screen's current theater_id
      // here just to satisfy that check, even though it has no effect.
      await api.put(`/screens/${editingScreen.value.screen_id}`, {
        theater_id: form.value.theater_id,
        name: form.value.name,
        total_seats: form.value.total_seats,
      })
    } else {
      await api.post('/screens', {
        name: form.value.name,
        theater_id: form.value.theater_id,
        total_seats: form.value.total_seats,
      })
    }

    await loadScreens()
    closeForm()
  } catch (err) {
    if (axios.isAxiosError(err) && err.response?.data?.error) {
      formError.value = err.response.data.error
    } else {
      formError.value = 'Failed to save screen. Please check the form and try again.'
    }
    console.error(err)
  } finally {
    saving.value = false
  }
}

async function deleteScreen(screen: Screen) {
  const confirmed = window.confirm(
    `Are you sure you want to delete "${screen.name}"?`,
  )

  if (!confirmed) return

  errorMessage.value = ''

  try {
    await api.delete(`/screens/${screen.screen_id}`)
    // Only remove from local state after the backend confirms deletion.
    screens.value = screens.value.filter((s) => s.screen_id !== screen.screen_id)
  } catch (err) {
    if (axios.isAxiosError(err) && err.response?.status === 409) {
      errorMessage.value =
        err.response.data?.error ?? 'Cannot delete screen because it is used by showtimes.'
    } else if (axios.isAxiosError(err) && err.response?.status === 404) {
      errorMessage.value = 'Screen not found. It may have already been deleted.'
      await loadScreens()
    } else {
      errorMessage.value = 'Failed to delete screen. Please try again.'
    }
    console.error(err)
  }
}
</script>

<template>
  <div class="min-h-screen bg-[#0d0d12] px-6 py-8 lg:px-8">

    <!-- Page Header -->
    <div class="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
      <div class="flex items-center gap-4">

        <!-- Screen Icon -->
        <div
          class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-gradient-to-br from-purple-500 to-pink-500 shadow-lg shadow-purple-500/20"
        >
          <Monitor class="size-6" />
        </div>

        <div>
          <h1 class="text-2xl font-bold text-white">
            Screen Management
          </h1>

          <p class="mt-1 text-sm text-gray-400">
            Manage cinema screens
          </p>
        </div>

      </div>

      <button
        type="button"
        class="inline-flex items-center justify-center gap-2 rounded-lg bg-[#ff3158] px-4 py-2.5 text-sm font-semibold text-white shadow-lg shadow-[#ff3158]/20 transition hover:bg-[#ff4669]"
        @click="openAddForm"
      >
        <Plus class="size-4" />
        Add screen
      </button>
    </div>

    <!-- Error banner -->
    <div
      v-if="errorMessage"
      class="mt-5 rounded-lg border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400"
    >
      {{ errorMessage }}
    </div>

    <!-- Search + Filter -->
    <div class="mt-7 flex flex-col gap-3 xl:flex-row xl:items-center xl:justify-between">

      <div class="relative w-full xl:max-w-xs">
        <Search
          class="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-[#68758f]"
        />
        <input
          v-model="search"
          type="text"
          placeholder="Search by name or theater..."
          class="w-full rounded-lg border border-white/10 bg-[#15151d] py-2.5 pl-9 pr-4 text-sm text-white outline-none transition placeholder:text-[#68758f] focus:border-[#ff3158]/60"
        />
      </div>

      <select
        v-model="selectedTheaterId"
        class="h-10 w-full rounded-lg border border-border/10 bg-card px-4 text-sm text-text outline-none focus:border-primary/50 xl:w-56"
      >
        <option value="All">All theaters</option>
        <option
          v-for="theaterId in theaterFilterOptions.filter((t) => t !== 'All')"
          :key="theaterId"
          :value="theaterId"
        >
          {{ theaterName(theaterId) }}
        </option>
      </select>

    </div>

    <!-- Screens Table -->
    <div class="mt-5 overflow-hidden rounded-xl border border-white/10 bg-[#15151d]">

      <div class="overflow-x-auto">

         <table class="w-full min-w-[900px]">

          <thead>
            <tr class="border-b border-white/10 bg-[#1b1b24]">

              <th class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]">
                Screen
              </th>

              <th class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]">
                Theater
              </th>

              <th class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]">
                Total seats
              </th>

              <th class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]">
                Actions
              </th>
            </tr>
          </thead>

          <tbody class="divide-y divide-border/10">

            <tr v-if="loading">
              <td colspan="4" class="px-5 py-12 text-center text-sm text-muted">
                Loading screens...
              </td>
            </tr>

            <tr
              v-for="screen in filteredScreens"
              v-else
              :key="screen.screen_id"
              class="transition hover:bg-white/[0.02]"
            >

              <!-- Screen -->
              <td class="px-5 py-4">
                <div class="flex items-center gap-3">

                  <div
                    class="flex size-9 items-center justify-center rounded-lg bg-primary/10 text-primary"
                  >
                    <MonitorPlay class="size-4" />
                  </div>

                  <div>
                    <p class="text-sm font-semibold text-text">
                      {{ screen.name }}
                    </p>

                    <p class="mt-0.5 text-xs text-muted">
                      ID #{{ screen.id }}
                    </p>
                  </div>

                </div>
              </td>

              <!-- Theater -->
              <td class="px-5 py-4 text-sm text-text">
                {{ theaterName(screen.theater_id) }}
              </td>

              <!-- Total seats -->
              <td class="px-5 py-4 text-sm font-medium text-text">
                {{ screen.total_seats }}
                <span class="text-xs text-muted">
                  seats
                </span>
              </td>

              <!-- Actions -->
              <td class="px-5 py-4">
                <div class="flex justify-end gap-2">

                  <button
                    type="button"
                    class="rounded-md border border-border/10 p-2 text-muted transition hover:border-primary/30 hover:bg-primary/10 hover:text-primary"
                    title="Edit"
                    @click="openEditForm(screen)"
                  >
                    <Pencil class="size-4" />
                  </button>

                  <button
                    type="button"
                    class="rounded-md border border-red-500/20 p-2 text-red-400 transition hover:bg-red-500/10"
                    title="Delete"
                    @click="deleteScreen(screen)"
                  >
                    <Trash2 class="size-4" />
                  </button>

                </div>
              </td>

            </tr>

            <!-- Empty -->
            <tr v-if="!loading && filteredScreens.length === 0">
              <td
                colspan="4"
                class="px-5 py-12 text-center text-sm text-muted"
              >
                No screens found.
              </td>
            </tr>

          </tbody>

        </table>
      </div>
    </div>

    <!-- Add / Edit modal -->
    <Transition name="modal">
      <div
        v-if="showForm"
        class="fixed inset-0 z-[100] flex items-start justify-center overflow-y-auto bg-black/70 p-4 backdrop-blur-sm sm:items-center"
        @click.self="closeForm"
      >
        <div
          class="modal-card my-auto w-full max-w-2xl rounded-2xl border border-border/20 bg-card p-6 shadow-2xl shadow-black/50 sm:p-7"
          role="dialog"
          aria-modal="true"
          aria-labelledby="screen-modal-title"
        >
          <!-- Modal Header -->
          <div class="mb-8 flex items-start justify-between gap-4">
            <div>
              <h2
                id="screen-modal-title"
                class="text-2xl font-bold tracking-tight text-text"
              >
                {{ editingScreen ? 'EDIT SCREEN' : 'ADD SCREEN' }}
              </h2>

              <p class="mt-1 text-sm text-muted">
                {{
                  editingScreen
                    ? 'แก้ไขข้อมูลห้องฉาย'
                    : 'เพิ่มห้องฉายใหม่เข้าสู่ระบบ'
                }}
              </p>
            </div>

            <button
              type="button"
              class="rounded-lg p-2 text-muted transition hover:bg-white/5 hover:text-text focus-visible:outline focus-visible:outline-2 focus-visible:outline-primary"
              aria-label="Close"
              @click="closeForm"
            >
              <X class="size-5" />
            </button>
          </div>

          <!-- Save error (shown here so it isn't hidden behind the overlay) -->
          <div
            v-if="formError"
            class="mb-6 rounded-lg border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400"
            role="alert"
          >
            {{ formError }}
          </div>

          <!-- Form -->
          <div class="grid gap-x-5 gap-y-6 sm:grid-cols-2">

            <!-- Screen Name -->
            <div>
              <label class="mb-2 block text-sm font-medium text-muted">
                Screen name
              </label>

              <input
                ref="nameInput"
                v-model="form.name"
                type="text"
                placeholder="e.g. Screen 1"
                class="h-14 w-full rounded-xl border border-border/20 bg-surface px-4 text-sm text-text outline-none transition placeholder:text-muted focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
                @keyup.enter="saveScreen"
              />
            </div>

            <!-- Theater -->
            <div>
              <label class="mb-2 block text-sm font-medium text-muted">
                Theater
              </label>

              <select
                v-model="form.theater_id"
                :disabled="!!editingScreen"
                class="h-14 w-full rounded-xl border border-border/20 bg-surface px-4 text-sm text-text outline-none transition focus:border-primary/60 focus:ring-2 focus:ring-primary/10 disabled:opacity-50"
              >
                <option
                  v-for="theater in theaters"
                  :key="theater.theater_id"
                  :value="theater.theater_id"
                >
                  {{ theater.name }}
                </option>
              </select>

              <p v-if="editingScreen" class="mt-1.5 text-xs text-muted">
                Theater cannot be changed after a screen is created.
              </p>
            </div>

            <!-- Total seats -->
            <div class="sm:col-span-2">
              <label class="mb-2 block text-sm font-medium text-muted">
                Total seats
              </label>

              <input
                v-model.number="form.total_seats"
                type="number"
                min="1"
                placeholder="e.g. 120"
                class="h-14 w-full rounded-xl border border-border/20 bg-surface px-4 text-sm text-text outline-none transition placeholder:text-muted focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
                @keyup.enter="saveScreen"
              />
            </div>
          </div>

          <!-- Buttons -->
          <div class="mt-8 flex justify-end gap-3">
            <button
              type="button"
              class="rounded-xl border border-border/20 px-6 py-3 text-sm font-semibold text-muted transition hover:bg-white/5 hover:text-text"
              @click="closeForm"
            >
              Cancel
            </button>

            <button
              type="button"
              :disabled="saving"
              class="rounded-lg bg-[#ff3158] px-5 py-2 text-sm font-semibold text-white transition hover:bg-[#ff4669] disabled:cursor-not-allowed disabled:opacity-60"
              @click="saveScreen"
            >
              {{
                saving
                  ? 'Saving...'
                  : editingScreen
                    ? 'Save changes'
                    : 'Save screen'
              }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
   </div>

</template>

<style scoped>
/* Tailwind classes can't reach the popup list a native <select>
   renders, since that's drawn by the browser itself — only real CSS
   on <option> (and <select>) can change its color. Supported in
   Chrome, Edge and Firefox; some browsers (e.g. Safari on macOS)
   ignore this and keep the OS-native white list regardless. */
select {
  color-scheme: dark;
}

select option {
  background-color: #15151d;
  color: #e5e7eb;
}

/* Modal open/close: a short fade for the backdrop and a small rise for
   the card. This only plays in response to the person's click. */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.18s ease;
}

.modal-enter-active .modal-card,
.modal-leave-active .modal-card {
  transition: transform 0.18s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-card,
.modal-leave-to .modal-card {
  transform: translateY(12px) scale(0.98);
}

@media (prefers-reduced-motion: reduce) {
  .modal-enter-active,
  .modal-leave-active,
  .modal-enter-active .modal-card,
  .modal-leave-active .modal-card {
    transition: none;
  }
}
</style>