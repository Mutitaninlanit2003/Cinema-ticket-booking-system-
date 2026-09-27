<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { X, Pencil, Trash2, Plus, Search, CalendarClock } from 'lucide-vue-next'
import api from '@/services/api'

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

interface Movie {
  id: number
  movie_id: string
  title: string
  name?: string
}

interface Theater {
  id: number
  theater_id: string
  name: string
  address?: string
}

interface Screen {
  id: number
  screen_id: string
  name?: string
  screen_name?: string
  theater_id?: string
  theater_ref_id?: number
}

interface Showtime {
  id: number
  showtime_id: string
  movie_id: string
  theater_id: string
  screen_id?: string | null
  show_date: string
  start_time: string
  end_time: string
  // Per-seat-type pricing (replaces the old single flat `price` field) —
  // matches models.Showtime.PriceStandard/PricePremium/PriceVIP.
  price_standard: number
  price_premium: number
  price_vip: number
  status: 'scheduled' | 'showing' | 'finished' | 'cancelled'
}

type ShowtimeForm = {
  movie_id: string
  theater_id: string
  screen_id: string
  show_date: string
  start_time: string
  end_time: string
  price_standard: number
  price_premium: number
  price_vip: number
  status: Showtime['status']
}

type FilterOption = 'all' | Showtime['status']

// ---------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------

const movies = ref<Movie[]>([])
const theaters = ref<Theater[]>([])
const screens = ref<Screen[]>([])
const showtimes = ref<Showtime[]>([])

const searchQuery = ref('')
const activeFilter = ref<FilterOption>('all')

const showForm = ref(false)
const editingShowtime = ref<Showtime | null>(null)
const movieSelectRef = ref<HTMLSelectElement | null>(null)

const loading = ref(false)
const saving = ref(false)
const deletingId = ref<number | null>(null)
const error = ref('')

const form = ref<ShowtimeForm>({
  movie_id: '',
  theater_id: '',
  screen_id: '',
  show_date: '',
  start_time: '',
  end_time: '',
  price_standard: 140,
  price_premium: 200,
  price_vip: 260,
  status: 'scheduled',
})

// ---------------------------------------------------------------------------
// Options
// ---------------------------------------------------------------------------

const STATUS_OPTIONS: Showtime['status'][] = [
  'scheduled',
  'showing',
  'finished',
  'cancelled',
]

const FILTER_OPTIONS: { value: FilterOption; label: string }[] = [
  { value: 'all', label: 'All' },
  { value: 'scheduled', label: 'Scheduled' },
  { value: 'showing', label: 'Showing' },
  { value: 'finished', label: 'Finished' },
  { value: 'cancelled', label: 'Cancelled' },
]

// ---------------------------------------------------------------------------
// Computed
// ---------------------------------------------------------------------------

const filteredShowtimes = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()

  return showtimes.value.filter((showtime) => {
    const movieName = getMovieName(showtime.movie_id).toLowerCase()
    const theaterName = getTheaterName(showtime.theater_id).toLowerCase()
    const screenName = getScreenName(showtime.screen_id).toLowerCase()

    const matchesSearch =
      !query ||
      movieName.includes(query) ||
      theaterName.includes(query) ||
      screenName.includes(query) ||
      showtime.show_date.includes(query)

    const matchesFilter =
      activeFilter.value === 'all' ||
      showtime.status === activeFilter.value

    return matchesSearch && matchesFilter
  })
})

// ---------------------------------------------------------------------------
// Theater -> Screen
// ---------------------------------------------------------------------------

const availableScreens = computed(() => {
  if (!form.value.theater_id) {
    return []
  }

  return screens.value.filter((screen) => {
    return screen.theater_id === form.value.theater_id
  })
})

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function getMovieName(movieId: string) {
  const movie = movies.value.find(
    (item) => item.movie_id === movieId,
  )

  return movie?.title || movie?.name || 'Unknown Movie'
}

function getTheaterName(theaterId: string) {
  const theater = theaters.value.find(
    (item) => item.theater_id === theaterId,
  )

  return theater?.name || 'Unknown Theater'
}

function getScreenName(screenId?: string | null) {
  if (!screenId) {
    return 'No Screen'
  }

  const screen = screens.value.find(
    (item) => item.screen_id === screenId,
  )

  return (
    screen?.name ||
    screen?.screen_name ||
    `Screen ${screenId}`
  )
}

function statusLabel(status: Showtime['status']) {
  return status
}

function formatDate(date: string) {
  if (!date) return '-'

  // Backend อาจส่ง 2026-08-30T00:00:00Z
  // หรือ 2026-08-30
  const dateOnly = date.substring(0, 10)

  const parsed = new Date(`${dateOnly}T00:00:00`)

  if (Number.isNaN(parsed.getTime())) {
    return date
  }

  return parsed.toLocaleDateString('en-GB', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  })
}

function formatTime(time: string) {
  if (!time) return '-'

  return time.length >= 5
    ? time.substring(0, 5)
    : time
}

function resetForm() {
  form.value = {
    movie_id: '',
    theater_id: '',
    screen_id: '',
    show_date: '',
    start_time: '',
    end_time: '',
    price_standard: 140,
    price_premium: 200,
    price_vip: 260,
    status: 'scheduled',
  }
}

// ---------------------------------------------------------------------------
// API - Load Movies
// ---------------------------------------------------------------------------

async function loadMovies() {
  try {
    const response = await api.get('/movies')

    movies.value = Array.isArray(response.data)
      ? response.data
      : response.data?.data || []
  } catch (err) {
    console.error('Failed to load movies:', err)
    throw new Error('Failed to load movies')
  }
}

// ---------------------------------------------------------------------------
// API - Load Theaters
// ---------------------------------------------------------------------------

async function loadTheaters() {
  try {
    const response = await api.get('/theaters')

    theaters.value = Array.isArray(response.data)
      ? response.data
      : response.data?.data || []
  } catch (err) {
    console.error('Failed to load theaters:', err)
    throw new Error('Failed to load theaters')
  }
}

// ---------------------------------------------------------------------------
// API - Load Screens
// ---------------------------------------------------------------------------

async function loadScreens() {
  try {
    const response = await api.get('/screens')

    screens.value = Array.isArray(response.data)
      ? response.data
      : response.data?.data || []
  } catch (err) {
    console.error('Failed to load screens:', err)
    throw new Error('Failed to load screens')
  }
}

// ---------------------------------------------------------------------------
// API - Load Showtimes
// ---------------------------------------------------------------------------

async function loadShowtimes() {
  try {
    const response = await api.get('/showtimes')

    showtimes.value = Array.isArray(response.data)
      ? response.data
      : response.data?.data || []
  } catch (err) {
    console.error('Failed to load showtimes:', err)
    throw new Error('Failed to load showtimes')
  }
}

// ---------------------------------------------------------------------------
// Load everything
// ---------------------------------------------------------------------------

async function loadData() {
  loading.value = true
  error.value = ''

  try {
    await Promise.all([
      loadMovies(),
      loadTheaters(),
      loadScreens(),
      loadShowtimes(),
    ])
  } catch (err) {
    console.error(err)
    error.value = 'Failed to load data from server.'
  } finally {
    loading.value = false
  }
}

// ---------------------------------------------------------------------------
// Modal behaviour: Esc closes it, page behind stops scrolling, and the
// first field is focused as soon as it opens — same pattern as the
// Screen Management modal.
// ---------------------------------------------------------------------------

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape' && showForm.value) {
    closeForm()
  }
}

watch(showForm, async (open) => {
  document.body.style.overflow = open ? 'hidden' : ''

  if (open) {
    await nextTick()
    movieSelectRef.value?.focus()
  }
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
  document.body.style.overflow = ''
})

// ---------------------------------------------------------------------------
// Form actions
// ---------------------------------------------------------------------------

function openAddForm() {
  editingShowtime.value = null
  error.value = ''
  resetForm()

  showForm.value = true
}

function openEditForm(showtime: Showtime) {
  editingShowtime.value = showtime
  error.value = ''

  form.value = {
    movie_id: showtime.movie_id,
    theater_id: showtime.theater_id,
    screen_id: showtime.screen_id || '',
    show_date: showtime.show_date.substring(0, 10),
    start_time: formatTime(showtime.start_time),
    end_time: formatTime(showtime.end_time),
    price_standard: Number(showtime.price_standard),
    price_premium: Number(showtime.price_premium),
    price_vip: Number(showtime.price_vip),
    status: showtime.status,
  }

  showForm.value = true
}

function closeForm() {
  if (saving.value) return

  showForm.value = false
  editingShowtime.value = null
  error.value = ''
  resetForm()
}

// ---------------------------------------------------------------------------
// Theater changed
// ---------------------------------------------------------------------------

function onTheaterChange() {
  const screenExists = availableScreens.value.some(
    (screen) => screen.screen_id === form.value.screen_id,
  )

  if (!screenExists) {
    form.value.screen_id = ''
  }
}

// ---------------------------------------------------------------------------
// Create / Update Showtime
// ---------------------------------------------------------------------------

async function saveShowtime() {
  error.value = ''

  if (!form.value.movie_id) {
    error.value = 'Please select a movie.'
    return
  }

  if (!form.value.theater_id) {
    error.value = 'Please select a theater.'
    return
  }

  if (!form.value.screen_id) {
    error.value = 'Please select a screen.'
    return
  }

  if (!form.value.show_date) {
    error.value = 'Please select a show date.'
    return
  }

  if (!form.value.start_time) {
    error.value = 'Please select a start time.'
    return
  }

  if (!form.value.end_time) {
    error.value = 'Please select an end time.'
    return
  }

  if (!form.value.price_standard || form.value.price_standard <= 0) {
    error.value = 'Standard price must be greater than 0.'
    return
  }

  if (!form.value.price_premium || form.value.price_premium <= 0) {
    error.value = 'Premium price must be greater than 0.'
    return
  }

  if (!form.value.price_vip || form.value.price_vip <= 0) {
    error.value = 'VIP price must be greater than 0.'
    return
  }

  saving.value = true

  try {
    // Backend's ShowDate field is a Go time.Time, which requires a full
    // RFC3339 timestamp when unmarshaling JSON (e.g. "2026-09-09T00:00:00Z").
    // The <input type="date"> only gives us "2026-09-09", so pad it out
    // here before sending — otherwise the backend fails to parse it.
    const payload = {
      movie_id: form.value.movie_id,
      theater_id: form.value.theater_id,
      screen_id: form.value.screen_id || null,
      show_date: `${form.value.show_date}T00:00:00Z`,
      start_time: form.value.start_time,
      end_time: form.value.end_time,
      price_standard: Number(form.value.price_standard),
      price_premium: Number(form.value.price_premium),
      price_vip: Number(form.value.price_vip),
      status: form.value.status,
    }

    if (editingShowtime.value) {
      await api.put(
        `/showtimes/${editingShowtime.value.showtime_id}`,
        payload,
      )
    } else {
      await api.post('/showtimes', payload)
    }

    await loadShowtimes()

    closeForm()
  } catch (err: any) {
    console.error('Failed to save showtime:', err)

    error.value =
      err?.response?.data?.error ||
      'Failed to save showtime.'
  } finally {
    saving.value = false
  }
}

// ---------------------------------------------------------------------------
// Delete Showtime
// ---------------------------------------------------------------------------

async function deleteShowtime(showtime: Showtime) {
  const confirmed = window.confirm(
    `Delete showtime for "${getMovieName(showtime.movie_id)}" at ${getTheaterName(showtime.theater_id)}?`,
  )

  if (!confirmed) return

  deletingId.value = showtime.id
  error.value = ''

  try {
    await api.delete(
      `/showtimes/${showtime.showtime_id}`,
    )

    showtimes.value = showtimes.value.filter(
      (item) => item.id !== showtime.id,
    )
  } catch (err: any) {
    console.error('Failed to delete showtime:', err)

    error.value =
      err?.response?.data?.error ||
      'Failed to delete showtime.'
  } finally {
    deletingId.value = null
  }
}

// ---------------------------------------------------------------------------
// Initial Load
// ---------------------------------------------------------------------------

onMounted(() => {
  loadData()
  window.addEventListener('keydown', onKeydown)
})
</script>

<template>
  <div class="min-h-screen px-6 py-8">

    <!-- Header -->
    <div class="mb-6 flex items-center justify-between">
      <div class="flex items-center gap-4">

        <div
          class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-gradient-to-br from-purple-500 to-pink-500 shadow-lg shadow-purple-500/20"
        >
          <CalendarClock class="size-6 text-white" />
        </div>

        <div>
          <h1 class="text-2xl font-bold text-white">
            Showtime Management
          </h1>

          <p class="mt-1 text-sm text-gray-400">
            Manage movie screening schedules
          </p>
        </div>

      </div>

      <button
        type="button"
        class="inline-flex items-center justify-center gap-2 rounded-lg bg-[#ff3158] px-4 py-2.5 text-sm font-semibold text-white shadow-lg shadow-[#ff3158]/20 transition hover:bg-[#ff4669]"
        @click="openAddForm"
      >
        <Plus class="size-4" />
        Add showtime
      </button>
    </div>

    <!-- Error (page-level — e.g. delete/load failures. While the modal is
         open, the same error shows inside the modal instead, since the
         page banner sits behind the overlay.) -->
    <div
      v-if="error && !showForm"
      class="mb-5 rounded-xl border border-red-500/20 bg-red-500/10 px-4 py-3 text-sm text-red-400"
    >
      {{ error }}
    </div>

    <!-- Search + Filter -->
    <div
      class="mb-5 flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between"
    >
      <div class="relative w-full lg:max-w-sm">

        <Search
          class="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-gray-500"
        />

        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search showtimes..."
          class="h-10 w-full rounded-lg border border-white/10 bg-[#15151d] pl-9 pr-4 text-sm text-white outline-none placeholder:text-gray-500 transition focus:border-white/20 focus:ring-1 focus:ring-white/10"
        />

      </div>

      <div class="flex flex-wrap gap-2">

        <button
          v-for="option in FILTER_OPTIONS"
          :key="option.value"
          type="button"
          class="rounded-lg border border-white/10 bg-[#15151d] px-4 py-2 text-xs font-semibold text-gray-400 transition hover:text-white"
          :class="{
            'border-[#ff2f5b] bg-[#ff2f5b] text-white':
              activeFilter === option.value,
          }"
          @click="activeFilter = option.value"
        >
          {{ option.label }}
        </button>

      </div>
    </div>

    <!-- Loading -->
    <div
      v-if="loading"
      class="flex min-h-[300px] items-center justify-center rounded-2xl border border-white/10 bg-[#15151d]"
    >
      <div class="text-center">
        <div
          class="mx-auto mb-3 h-8 w-8 animate-spin rounded-full border-2 border-white/10 border-t-[#ff3158]"
        ></div>

        <p class="text-sm text-gray-400">
          Loading showtimes...
        </p>
      </div>
    </div>

    <!-- Showtime Table -->
    <div
      v-else
      class="overflow-hidden rounded-2xl border border-white/10 bg-[#15151d]"
    >

      <div class="overflow-x-auto">

        <table class="w-full text-left">

          <thead
            class="border-b border-white/10 bg-[#1a1a23]"
          >
            <tr
              class="text-xs uppercase tracking-wide text-gray-400"
            >
              <th class="px-5 py-4">
                Movie
              </th>

              <th class="px-5 py-4">
                Theater / Screen
              </th>

              <th class="px-5 py-4">
                Show Date
              </th>

              <th class="px-5 py-4">
                Time
              </th>

              <th class="px-5 py-4">
                Price (Std / Prem / VIP)
              </th>

              <th class="px-5 py-4">
                Status
              </th>

              <th class="px-5 py-4 text-right">
                Actions
              </th>
            </tr>
          </thead>

          <tbody
            v-if="filteredShowtimes.length > 0"
            class="divide-y divide-white/10"
          >

            <tr
              v-for="showtime in filteredShowtimes"
              :key="showtime.showtime_id"
              class="transition hover:bg-white/[0.02]"
            >

              <!-- Movie -->
              <td class="px-5 py-5">
                <span class="text-sm font-semibold text-white">
                  {{ getMovieName(showtime.movie_id) }}
                </span>
              </td>

              <!-- Theater / Screen -->
              <td class="px-5 py-5">

                <div class="text-sm text-white">
                  {{ getTheaterName(showtime.theater_id) }}
                </div>

                <div class="mt-1 text-xs text-gray-400">
                  {{ getScreenName(showtime.screen_id) }}
                </div>

              </td>

              <!-- Date -->
              <td class="px-5 py-5 text-sm text-gray-400">
                {{ formatDate(showtime.show_date) }}
              </td>

              <!-- Time -->
              <td class="px-5 py-5 text-sm text-white">
                {{ formatTime(showtime.start_time) }}
                -
                {{ formatTime(showtime.end_time) }}
              </td>

              <!-- Price (per zone) -->
              <td class="px-5 py-5">
                <div class="flex flex-wrap gap-1.5 text-xs">
                  <span class="rounded-md border border-emerald-500/30 bg-emerald-500/10 px-2 py-1 font-medium text-emerald-300">
                    ฿{{ Number(showtime.price_standard).toFixed(0) }}
                  </span>
                  <span class="rounded-md border border-purple-500/30 bg-purple-500/10 px-2 py-1 font-medium text-purple-300">
                    ฿{{ Number(showtime.price_premium).toFixed(0) }}
                  </span>
                  <span class="rounded-md border border-blue-500/30 bg-blue-500/10 px-2 py-1 font-medium text-blue-300">
                    ฿{{ Number(showtime.price_vip).toFixed(0) }}
                  </span>
                </div>
              </td>

              <!-- Status -->
              <td class="px-5 py-5">

                <span
                  class="rounded-full px-3 py-1 text-xs font-medium"
                  :class="{
                    'bg-emerald-500/10 text-emerald-400':
                      showtime.status === 'scheduled',

                    'bg-amber-500/10 text-amber-400':
                      showtime.status === 'showing',

                    'bg-white/5 text-gray-400':
                      showtime.status === 'finished',

                    'bg-red-500/10 text-red-400':
                      showtime.status === 'cancelled',
                  }"
                >
                  {{ statusLabel(showtime.status) }}
                </span>

              </td>

              <!-- Actions -->
              <td class="px-5 py-5">

                <div class="flex justify-end gap-2">

                  <button
                    type="button"
                    class="rounded-lg border border-white/10 p-2 text-gray-400 transition hover:bg-white/5 hover:text-white disabled:opacity-50"
                    :disabled="deletingId === showtime.id"
                    @click="openEditForm(showtime)"
                  >
                    <Pencil class="size-4" />
                  </button>

                  <button
                    type="button"
                    class="rounded-lg border border-red-500/20 p-2 text-red-400 transition hover:bg-red-500/10 disabled:opacity-50"
                    :disabled="deletingId === showtime.id"
                    @click="deleteShowtime(showtime)"
                  >
                    <Trash2 class="size-4" />
                  </button>

                </div>

              </td>

            </tr>

          </tbody>

          <!-- Empty -->
          <tbody v-else>
            <tr>
              <td
                colspan="7"
                class="px-5 py-16 text-center"
              >
                <CalendarClock
                  class="mx-auto mb-3 size-10 text-gray-600"
                />

                <p class="text-sm font-medium text-gray-400">
                  No showtimes found
                </p>

                <p class="mt-1 text-xs text-gray-600">
                  Add a showtime or change your search/filter.
                </p>
              </td>
            </tr>
          </tbody>

        </table>

      </div>

    </div>

    <!-- ========================================================= -->
    <!-- ADD / EDIT MODAL (same modal pattern as Screen Management) -->
    <!-- ========================================================= -->

    <Transition name="modal">
      <div
        v-if="showForm"
        class="fixed inset-0 z-[100] flex items-start justify-center overflow-y-auto bg-black/70 p-4 backdrop-blur-sm sm:items-center"
        @click.self="closeForm"
      >
        <div
          class="modal-card my-auto w-full max-w-2xl rounded-2xl border border-white/10 bg-[#15151d] p-6 shadow-2xl shadow-black/50 sm:p-7"
          role="dialog"
          aria-modal="true"
          aria-labelledby="showtime-modal-title"
        >
          <!-- Modal Header -->
          <div class="mb-8 flex items-start justify-between gap-4">
            <div>
              <h2
                id="showtime-modal-title"
                class="text-2xl font-bold tracking-tight text-white"
              >
                {{ editingShowtime ? 'EDIT SHOWTIME' : 'ADD SHOWTIME' }}
              </h2>

              <p class="mt-1 text-sm text-gray-400">
                {{
                  editingShowtime
                    ? 'แก้ไขข้อมูลรอบฉาย'
                    : 'เพิ่มรอบฉายใหม่เข้าสู่ระบบ'
                }}
              </p>
            </div>

            <button
              type="button"
              class="rounded-lg p-2 text-gray-400 transition hover:bg-white/5 hover:text-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#ff3158]"
              aria-label="Close"
              :disabled="saving"
              @click="closeForm"
            >
              <X class="size-5" />
            </button>
          </div>

          <!-- Save error (shown here so it isn't hidden behind the overlay) -->
          <div
            v-if="error"
            class="mb-6 rounded-lg border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400"
            role="alert"
          >
            {{ error }}
          </div>

          <!-- Form -->
          <div class="grid gap-x-5 gap-y-6 sm:grid-cols-2">

            <!-- Movie -->
            <div>

              <label class="mb-2 block text-sm font-medium text-gray-400">
                Movie
              </label>

              <select
                ref="movieSelectRef"
                v-model="form.movie_id"
                class="h-12 w-full rounded-xl border border-white/10 bg-[#101017] px-4 text-sm text-white outline-none focus:border-[#ff3158]/60 focus:ring-2 focus:ring-[#ff3158]/10"
              >

                <option value="">
                  Select movie
                </option>

                <option
                  v-for="movie in movies"
                  :key="movie.movie_id"
                  :value="movie.movie_id"
                >
                  {{ movie.title || movie.name }}
                </option>

              </select>

            </div>

            <!-- Theater -->
            <div>

              <label class="mb-2 block text-sm font-medium text-gray-400">
                Theater
              </label>

              <select
                v-model="form.theater_id"
                class="h-12 w-full rounded-xl border border-white/10 bg-[#101017] px-4 text-sm text-white outline-none focus:border-[#ff3158]/60 focus:ring-2 focus:ring-[#ff3158]/10"
                @change="onTheaterChange"
              >

                <option value="">
                  Select theater
                </option>

                <option
                  v-for="theater in theaters"
                  :key="theater.theater_id"
                  :value="theater.theater_id"
                >
                  {{ theater.name }}
                </option>

              </select>

            </div>

            <!-- Screen -->
            <div>

              <label class="mb-2 block text-sm font-medium text-gray-400">
                Screen
              </label>

              <select
                v-model="form.screen_id"
                :disabled="!form.theater_id"
                class="h-12 w-full rounded-xl border border-white/10 bg-[#101017] px-4 text-sm text-white outline-none focus:border-[#ff3158]/60 focus:ring-2 focus:ring-[#ff3158]/10 disabled:cursor-not-allowed disabled:opacity-50"
              >

                <option value="">
                  {{
                    form.theater_id
                      ? 'Select screen'
                      : 'Select theater first'
                  }}
                </option>

                <option
                  v-for="screen in availableScreens"
                  :key="screen.screen_id"
                  :value="screen.screen_id"
                >
                  {{ screen.name || screen.screen_name || `Screen ${screen.screen_id}` }}
                </option>

              </select>

              <p
                v-if="form.theater_id && availableScreens.length === 0"
                class="mt-2 text-xs text-amber-400"
              >
                No screens found for this theater.
              </p>

            </div>

            <!-- Date -->
            <div>

              <label class="mb-2 block text-sm font-medium text-gray-400">
                Show date
              </label>

              <input
                v-model="form.show_date"
                type="date"
                class="h-12 w-full rounded-xl border border-white/10 bg-[#101017] px-4 text-sm text-white outline-none focus:border-[#ff3158]/60 focus:ring-2 focus:ring-[#ff3158]/10"
              />

            </div>

            <!-- Start -->
            <div>

              <label class="mb-2 block text-sm font-medium text-gray-400">
                Start time
              </label>

              <input
                v-model="form.start_time"
                type="time"
                lang="en-GB"
                class="h-12 w-full rounded-xl border border-white/10 bg-[#101017] px-4 text-sm text-white outline-none focus:border-[#ff3158]/60 focus:ring-2 focus:ring-[#ff3158]/10"
              />

            </div>

            <!-- End -->
            <div>

              <label class="mb-2 block text-sm font-medium text-gray-400">
                End time
              </label>

              <input
                v-model="form.end_time"
                type="time"
                lang="en-GB"
                class="h-12 w-full rounded-xl border border-white/10 bg-[#101017] px-4 text-sm text-white outline-none focus:border-[#ff3158]/60 focus:ring-2 focus:ring-[#ff3158]/10"
              />

            </div>

            <!-- Price per zone -->
            <div class="sm:col-span-2">

              <label class="mb-2 block text-sm font-medium text-gray-400">
                Price per seat zone
              </label>

              <div class="grid grid-cols-3 gap-3">

                <div>
                  <label class="mb-1.5 flex items-center gap-1.5 text-xs font-medium text-emerald-400">
                    <span class="h-2 w-2 rounded-full bg-emerald-400"></span>
                    Standard
                  </label>

                  <input
                    v-model.number="form.price_standard"
                    type="number"
                    min="1"
                    placeholder="e.g. 140"
                    class="h-12 w-full rounded-xl border border-white/10 bg-[#101017] px-4 text-sm text-white outline-none placeholder:text-gray-600 focus:border-emerald-500/60 focus:ring-2 focus:ring-emerald-500/10"
                    @keyup.enter="saveShowtime"
                  />
                </div>

                <div>
                  <label class="mb-1.5 flex items-center gap-1.5 text-xs font-medium text-purple-400">
                    <span class="h-2 w-2 rounded-full bg-purple-400"></span>
                    Premium
                  </label>

                  <input
                    v-model.number="form.price_premium"
                    type="number"
                    min="1"
                    placeholder="e.g. 200"
                    class="h-12 w-full rounded-xl border border-white/10 bg-[#101017] px-4 text-sm text-white outline-none placeholder:text-gray-600 focus:border-purple-500/60 focus:ring-2 focus:ring-purple-500/10"
                    @keyup.enter="saveShowtime"
                  />
                </div>

                <div>
                  <label class="mb-1.5 flex items-center gap-1.5 text-xs font-medium text-blue-400">
                    <span class="h-2 w-2 rounded-full bg-blue-400"></span>
                    VIP
                  </label>

                  <input
                    v-model.number="form.price_vip"
                    type="number"
                    min="1"
                    placeholder="e.g. 260"
                    class="h-12 w-full rounded-xl border border-white/10 bg-[#101017] px-4 text-sm text-white outline-none placeholder:text-gray-600 focus:border-blue-500/60 focus:ring-2 focus:ring-blue-500/10"
                    @keyup.enter="saveShowtime"
                  />
                </div>

              </div>

              <p class="mt-2 text-xs text-gray-500">
                ราคานี้ใช้กับที่นั่งตาม seat_type ของทุกที่นั่งในรอบฉายนี้ ไม่ว่าจะเลือก Screen ประเภทไหนก็ตาม
              </p>

            </div>

            <!-- Status -->
            <div class="sm:col-span-2">

              <label class="mb-2 block text-sm font-medium text-gray-400">
                Status
              </label>

              <select
                v-model="form.status"
                class="h-12 w-full rounded-xl border border-white/10 bg-[#101017] px-4 text-sm text-white outline-none focus:border-[#ff3158]/60 focus:ring-2 focus:ring-[#ff3158]/10"
              >

                <option
                  v-for="status in STATUS_OPTIONS"
                  :key="status"
                  :value="status"
                >
                  {{ statusLabel(status) }}
                </option>

              </select>

            </div>

          </div>

          <!-- Buttons -->
          <div class="mt-8 flex justify-end gap-3">

            <button
              type="button"
              class="rounded-xl border border-white/10 px-6 py-3 text-sm font-semibold text-gray-400 transition hover:bg-white/5 hover:text-white disabled:opacity-50"
              :disabled="saving"
              @click="closeForm"
            >
              Cancel
            </button>

            <button
              type="button"
              class="inline-flex items-center gap-2 rounded-xl bg-[#ff3158] px-6 py-3 text-sm font-semibold text-white transition hover:bg-[#ff4669] disabled:cursor-not-allowed disabled:opacity-50"
              :disabled="saving"
              @click="saveShowtime"
            >

              <span
                v-if="saving"
                class="h-4 w-4 animate-spin rounded-full border-2 border-white/30 border-t-white"
              ></span>

              {{
                saving
                  ? 'Saving...'
                  : editingShowtime
                    ? 'Save changes'
                    : 'Save showtime'
              }}

            </button>

          </div>

        </div>
      </div>
    </Transition>

  </div>
</template>

<style scoped>
/* Tailwind classes can't reach the popup list a native <select> renders,
   since that's drawn by the browser itself — only real CSS on <option>
   (and <select>) can change its color. Supported in Chrome, Edge and
   Firefox; some browsers (e.g. Safari on macOS) ignore this and keep
   the OS-native white list regardless. */
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