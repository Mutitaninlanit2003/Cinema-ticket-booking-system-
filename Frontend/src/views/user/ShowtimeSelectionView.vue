<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/services/api'

const route = useRoute()
const router = useRouter()

// =========================
// Movie ID
// =========================

const movieId = computed(() => String(route.params.movie_id || ''))

// =========================
// Interfaces
// =========================

interface Theater {
  theater_id: string
  name: string
  address?: string
}

interface Showtime {
  id: number
  showtime_id: string
  movie_id: string
  theater_id: string
  screen_id?: string
  show_date: string
  start_time: string
  end_time?: string
  price: number
  status: 'scheduled' | 'showing' | 'finished' | 'cancelled'
}

// =========================
// Data
// =========================

const theaters = ref<Theater[]>([])
const showtimes = ref<Showtime[]>([])

const selectedTheater = ref<string | null>(null)
const selectedDate = ref<string | null>(null)
const selectedShowtime = ref<string | null>(null)

const loading = ref(true)
const error = ref('')

// =========================
// Movie Showtimes
// =========================

const movieShowtimes = computed(() => {
  return showtimes.value.filter(
    (showtime) => showtime.movie_id === movieId.value,
  )
})

// =========================
// Available Theaters
// =========================

const availableTheaters = computed(() => {
  return theaters.value
})

// =========================
// Available Dates
// =========================

// Zero-padded YYYY-MM-DD in LOCAL time (not toISOString, which shifts
// to UTC and can land on the wrong day depending on timezone).
function toISODate(date: Date) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')

  return `${year}-${month}-${day}`
}

// Every day from today through the end of the current calendar month.
// (We deliberately do NOT base this on the movie's release_date — for
// catalog/re-released titles release_date can be years in the past,
// which produced zero valid days and "No dates available".  "This
// month, from today onward" is what's actually bookable regardless.)
const availableDates = computed(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  const year = today.getFullYear()
  const month = today.getMonth()
  const daysInMonth = new Date(year, month + 1, 0).getDate()

  const dates: string[] = []

  for (let day = today.getDate(); day <= daysInMonth; day++) {
    dates.push(toISODate(new Date(year, month, day)))
  }

  return dates
})

// A date has a real, bookable showtime for the selected theater.
function dateHasShowtime(date: string) {
  return movieShowtimes.value.some(
    (showtime) =>
      showtime.theater_id === selectedTheater.value &&
      showtime.show_date.slice(0, 10) === date,
  )
}

// =========================
// Filtered Showtimes
// =========================

const filteredShowtimes = computed(() => {
  return movieShowtimes.value.filter((showtime) => {
    if (
      selectedTheater.value &&
      showtime.theater_id !== selectedTheater.value
    ) {
      return false
    }

    if (
      selectedDate.value &&
      showtime.show_date.slice(0, 10) !== selectedDate.value
    ) {
      return false
    }

    return true
  })
})

// =========================
// Selected Theater
// =========================

const selectedTheaterData = computed(() => {
  return theaters.value.find(
    (theater) =>
      theater.theater_id === selectedTheater.value,
  )
})

// =========================
// Select Theater
// =========================

function selectTheater(theaterId: string) {
  selectedTheater.value = theaterId

  selectedDate.value = null
  selectedShowtime.value = null

  // Prefer a date that already has a real showtime; otherwise just
  // land on the first day of the month so there's something selected.
  const dateWithShowtime = availableDates.value.find(dateHasShowtime)
  const firstDate = dateWithShowtime ?? availableDates.value[0]

  if (firstDate) {
    selectedDate.value = firstDate
  }
}

// =========================
// Select Date
// =========================

function selectDate(date: string) {
  selectedDate.value = date
  selectedShowtime.value = null
}

// =========================
// Select Showtime
// =========================

function selectShowtime(showtimeId: string) {
  const showtime = showtimes.value.find(
    (item) => item.showtime_id === showtimeId,
  )

  if (!showtime) return

  if (
    showtime.status === 'finished' ||
    showtime.status === 'cancelled'
  ) {
    return
  }

  selectedShowtime.value = showtimeId
}

// =========================
// Continue To Seats
// =========================

function continueToSeats() {
  if (!selectedShowtime.value) return

  router.push(`/seats/${selectedShowtime.value}`)
}

// =========================
// Get Theater Name
// =========================

function getTheaterName(theaterId: string) {
  const theater = theaters.value.find(
    (item) => item.theater_id === theaterId,
  )

  return theater?.name || 'Unknown Theater'
}

// =========================
// Format Date
// =========================

// Parses a plain "YYYY-MM-DD" string as local midnight (not UTC), so
// the displayed day never shifts depending on the browser's timezone.
function parseLocalDate(date: string) {
  return new Date(`${date}T00:00:00`)
}

function formatDate(date: string) {
  if (!date) return '-'

  return parseLocalDate(date.slice(0, 10)).toLocaleDateString('en-US', {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
  })
}

// =========================
// Format Time
// =========================

function formatTime(time: string) {
  if (!time) return '-'

  // ถ้า Backend ส่งมาเป็น HH:mm:ss
  if (time.length >= 5 && time.includes(':')) {
    return time.substring(0, 5)
  }

  return time
}

// =========================
// Load Theaters
// =========================

async function loadTheaters() {
  try {
    const response = await api.get('/theaters')

    theaters.value = Array.isArray(response.data)
      ? response.data
      : response.data.data || []
  } catch (err) {
    console.error('Failed to load theaters:', err)
  }
}

// =========================
// Load Showtimes
// =========================

async function loadShowtimes() {
  try {
    const response = await api.get('/showtimes')

    showtimes.value = Array.isArray(response.data)
      ? response.data
      : response.data.data || []
  } catch (err) {
    console.error('Failed to load showtimes:', err)

    error.value = 'Failed to load showtimes'
  }
}

// =========================
// Initial Load
// =========================

onMounted(async () => {
  loading.value = true
  error.value = ''

  try {
    await Promise.all([
      loadTheaters(),
      loadShowtimes(),
    ])

    // เลือกโรงแรกให้อัตโนมัติ
    if (availableTheaters.value.length > 0) {
      selectTheater(
        availableTheaters.value[0].theater_id,
      )
    }
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="min-h-screen bg-[#0d0d12] px-4 py-8 text-white md:px-8">
    <div
      class="mx-auto max-w-[1280px] rounded-3xl border border-white/10 bg-[#0d0d12] p-6 shadow-2xl md:p-10"
    >

      <!-- ========================= -->
      <!-- HEADER -->
      <!-- ========================= -->

      <div class="mb-10">
        <button
          type="button"
          class="mb-6 flex items-center gap-2 text-sm text-gray-400 transition hover:text-white"
          @click="router.back()"
        >
          <span class="text-lg">←</span>
          <span>Back</span>
        </button>

        <h1
          class="text-3xl font-black tracking-tight md:text-4xl"
        >
          Select Showtime
        </h1>

        <p class="mt-3 text-sm text-gray-500 md:text-base">
          Choose a theater, date and showtime to continue to seat selection
        </p>
      </div>

      <!-- ========================= -->
      <!-- STEP INDICATOR -->
      <!-- ========================= -->

      <div class="mb-10 flex items-center gap-3 text-sm">
        <!-- Step 1 -->
        <div class="flex items-center gap-2">
          <span
            class="flex h-7 w-7 items-center justify-center rounded-full bg-[#ff3158] text-xs font-bold text-white"
          >
            1
          </span>

          <span class="font-semibold text-white">
            Theater
          </span>
        </div>

        <div class="h-px w-10 bg-white/20"></div>

        <!-- Step 2 -->
        <div class="flex items-center gap-2">
          <span
            class="flex h-7 w-7 items-center justify-center rounded-full bg-[#ff3158] text-xs font-bold text-white"
          >
            2
          </span>

          <span class="font-semibold text-white">
            Date
          </span>
        </div>

        <div class="h-px w-10 bg-white/20"></div>

        <!-- Step 3 -->
        <div class="flex items-center gap-2">
          <span
            class="flex h-7 w-7 items-center justify-center rounded-full border border-white/20 text-xs text-gray-400"
          >
            3
          </span>

          <span
            class="font-medium"
            :class="
              selectedShowtime
                ? 'text-white'
                : 'text-gray-500'
            "
          >
            Showtime
          </span>
        </div>
      </div>

      <!-- ========================= -->
      <!-- LOADING -->
      <!-- ========================= -->

      <div
        v-if="loading"
        class="flex min-h-[400px] items-center justify-center"
      >
        <p class="text-sm text-gray-500">
          Loading showtimes...
        </p>
      </div>

      <!-- ========================= -->
      <!-- ERROR -->
      <!-- ========================= -->

      <div
        v-else-if="error"
        class="flex min-h-[400px] items-center justify-center"
      >
        <div
          class="rounded-2xl border border-red-500/20 bg-red-500/5 px-10 py-8 text-center"
        >
          <p class="text-red-400">
            {{ error }}
          </p>

          <button
            type="button"
            class="mt-5 rounded-lg border border-white/10 px-6 py-2 text-sm text-gray-300 transition hover:bg-white hover:text-black"
            @click="loadShowtimes"
          >
            Try again
          </button>
        </div>
      </div>

      <!-- ========================= -->
      <!-- CONTENT -->
      <!-- ========================= -->

      <div v-else>

        <!-- ========================= -->
        <!-- SELECT THEATER -->
        <!-- ========================= -->

        <section>
          <p
            class="mb-5 text-xs font-semibold uppercase tracking-[0.18em] text-gray-400"
          >
            Select Theater
          </p>

          <div
            v-if="availableTheaters.length === 0"
            class="rounded-2xl border border-white/10 bg-[#15151d] px-6 py-12 text-center"
          >
            <p class="text-gray-400">
              No theaters available for this movie.
            </p>
          </div>

          <div
            v-else
            class="flex gap-4 overflow-x-auto pb-2"
            style="scrollbar-width: thin;"
          >
            <button
              v-for="theater in availableTheaters"
              :key="theater.theater_id"
              type="button"
              class="w-[260px] flex-shrink-0 rounded-2xl border p-5 text-left transition duration-200"
              :class="
                selectedTheater === theater.theater_id
                  ? 'border-[#ff3158] bg-[#32111c] shadow-lg shadow-[#ff3158]/10'
                  : 'border-white/10 bg-[#15151d] hover:border-white/25 hover:bg-[#191922]'
              "
              @click="selectTheater(theater.theater_id)"
            >
              <div class="flex items-start justify-between gap-3">
                <div>
                  <h3 class="text-base font-bold text-white">
                    {{ theater.name }}
                  </h3>

                  <p
                    v-if="theater.address"
                    class="mt-2 text-xs leading-5 text-gray-500"
                  >
                    {{ theater.address }}
                  </p>
                </div>

                <span
                  v-if="selectedTheater === theater.theater_id"
                  class="mt-1 h-2.5 w-2.5 rounded-full bg-[#ff3158]"
                ></span>
              </div>
            </button>
          </div>
        </section>

        <!-- ========================= -->
        <!-- SELECT DATE -->
        <!-- ========================= -->

        <section class="mt-10">
          <p
            class="mb-5 text-xs font-semibold uppercase tracking-[0.18em] text-gray-400"
          >
            Select Date
          </p>

          <div
            v-if="availableDates.length === 0"
            class="rounded-2xl border border-white/10 bg-[#15151d] px-6 py-12 text-center"
          >
            <p class="text-gray-400">
              No dates available.
            </p>
          </div>

          <div
            v-else
            class="flex flex-wrap gap-3"
          >
            <button
              v-for="date in availableDates"
              :key="date"
              type="button"
              class="flex min-w-[72px] flex-col items-center justify-center rounded-xl border px-4 py-3 transition duration-200"
              :class="
                selectedDate === date
                  ? 'border-[#ff3158] bg-[#ff3158] text-white shadow-lg shadow-[#ff3158]/20'
                  : dateHasShowtime(date)
                    ? 'border-white/10 bg-[#15151d] text-gray-400 hover:border-white/25 hover:text-white'
                    : 'border-white/5 bg-[#15151d] text-gray-600 hover:border-white/15 hover:text-gray-400'
              "
              @click="selectDate(date)"
            >
              <span class="text-xs font-medium">
                {{ formatDate(date).split(' ')[0] }}
              </span>

              <span class="mt-1 text-lg font-bold">
                {{ parseLocalDate(date).getDate() }}
              </span>

              <span
                v-if="dateHasShowtime(date)"
                class="mt-1 h-1 w-1 rounded-full bg-emerald-400"
              ></span>
            </button>
          </div>
        </section>

        <!-- ========================= -->
        <!-- SELECT SHOWTIME -->
        <!-- ========================= -->

        <section class="mt-10">
          <div class="mb-5 flex flex-wrap items-center gap-2">
            <p
              class="text-xs font-semibold uppercase tracking-[0.18em] text-gray-400"
            >
              Select Showtime
            </p>

            <span
              v-if="selectedTheaterData"
              class="text-xs uppercase tracking-wider text-gray-600"
            >
              —
              {{ selectedTheaterData.name }}
            </span>

            <span
              v-if="selectedDate"
              class="text-xs uppercase tracking-wider text-gray-600"
            >
              ·
              {{ formatDate(selectedDate) }}
            </span>
          </div>

          <!-- No Showtime -->

          <div
            v-if="
              !selectedTheater ||
              !selectedDate ||
              filteredShowtimes.length === 0
            "
            class="rounded-2xl border border-white/10 bg-[#15151d] px-6 py-14 text-center"
          >
            <p class="text-gray-400">
              No showtimes available for this theater and date.
            </p>
          </div>

          <!-- Showtime Cards -->

          <div
            v-else
            class="grid gap-4 md:grid-cols-2 lg:grid-cols-3"
          >
            <button
              v-for="showtime in filteredShowtimes"
              :key="showtime.showtime_id"
              type="button"
              :disabled="
                showtime.status === 'finished' ||
                showtime.status === 'cancelled'
              "
              class="relative rounded-2xl border p-5 text-left transition duration-200"
              :class="
                showtime.status === 'finished' ||
                showtime.status === 'cancelled'
                  ? 'cursor-not-allowed border-white/10 bg-[#121218] opacity-40'
                  : selectedShowtime === showtime.showtime_id
                    ? 'border-violet-500 bg-[#211a35] shadow-lg shadow-violet-500/10'
                    : 'border-white/10 bg-[#15151d] hover:border-white/25 hover:bg-[#191922]'
              "
              @click="
                selectShowtime(showtime.showtime_id)
              "
            >

              <!-- Time + Status -->

              <div class="flex items-center justify-between gap-3">
                <span class="text-2xl font-black tracking-tight">
                  {{ formatTime(showtime.start_time) }}
                </span>

                <span
                  class="rounded-full px-3 py-1 text-[10px] font-bold uppercase tracking-wide"
                  :class="{
                    'bg-emerald-500/15 text-emerald-400':
                      showtime.status === 'scheduled',

                    'bg-orange-500/15 text-orange-400':
                      showtime.status === 'showing',

                    'bg-white/5 text-gray-500':
                      showtime.status === 'finished',

                    'bg-red-500/15 text-red-400':
                      showtime.status === 'cancelled',
                  }"
                >
                  {{ showtime.status }}
                </span>
              </div>

              <!-- Theater -->

              <div class="mt-5">
                <p class="text-sm font-semibold text-white">
                  {{ getTheaterName(showtime.theater_id) }}
                </p>

                <p class="mt-1 text-xs text-gray-500">
                  {{ formatDate(showtime.show_date) }}
                </p>
              </div>

              <!-- Screen / End -->

              <div
                class="mt-4 rounded-xl border border-white/5 bg-white/[0.03] px-4 py-3 text-center"
              >
                <p
                  class="text-[10px] tracking-[0.18em] text-gray-500"
                >
                  <span v-if="showtime.screen_id">
                    Screen {{ showtime.screen_id }}
                  </span>

                  <span v-if="showtime.end_time">
                    · Ends {{ formatTime(showtime.end_time) }}
                  </span>
                </p>
              </div>

              <!-- Price -->

              <div class="mt-5 flex items-end justify-between">
                <div>
                  <p class="text-[10px] uppercase tracking-wider text-gray-500">
                    Ticket Price
                  </p>

                  <p class="mt-1 text-base font-bold text-white">
                    ฿{{ Number(showtime.price).toFixed(2) }}
                  </p>
                </div>

                <!-- Selected -->

                <span
                  v-if="selectedShowtime === showtime.showtime_id"
                  class="flex h-7 w-7 items-center justify-center rounded-full bg-violet-500 text-xs font-bold text-white"
                >
                  ✓
                </span>
              </div>
            </button>
          </div>
        </section>

        <!-- ========================= -->
        <!-- CONTINUE -->
        <!-- ========================= -->

        <div class="mt-10">
          <button
            type="button"
            class="w-full rounded-xl bg-[#ff3158] py-4 text-sm font-bold text-white shadow-lg shadow-[#ff3158]/20 transition duration-200 hover:bg-[#ff4568] hover:shadow-xl hover:shadow-[#ff3158]/25 disabled:cursor-not-allowed disabled:opacity-30"
            :disabled="!selectedShowtime"
            @click="continueToSeats"
          >
            Continue to seat selection
          </button>
        </div>

      </div>
    </div>
  </div>
</template>