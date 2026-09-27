<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import api from '@/services/api'
import {
  Search,
  Ticket,
  CalendarDays,
  Clock3,
  Users,
} from 'lucide-vue-next'

// -----------------------------------------------------------------------
// Types — match backend models exactly
// (models/booking.go, models/showtime.go, models/movie.go, models/theater.go, models/user.go)
// -----------------------------------------------------------------------
type BookingStatus = 'pending' | 'confirmed' | 'cancelled'

interface Booking {
  id: number
  booking_id: string
  user_id: string
  showtime_id: string
  seat_numbers: string[]
  total_price: number
  status: BookingStatus
  created_at: string
}

interface Showtime {
  id: number
  showtime_id: string
  movie_id: string
  theater_id: string
  screen_id?: string
  show_date: string
  start_time: string
  end_time: string
  price: number
  status: string
}

interface Movie {
  id: number
  movie_id: string
  title: string
}

interface Theater {
  id: number
  theater_id: string
  name: string
}

interface User {
  id: number
  user_id: string
  name: string
  email: string
}

const search = ref('')
const activeFilter = ref<'all' | BookingStatus>('all')

const bookings = ref<Booking[]>([])
const showtimes = ref<Showtime[]>([])
const movies = ref<Movie[]>([])
const theaters = ref<Theater[]>([])
const users = ref<User[]>([])

const loading = ref(false)
const errorMessage = ref('')

// -----------------------------------------------------------------------
// Lookups
// -----------------------------------------------------------------------
const showtimeById = computed(() => {
  const map = new Map<string, Showtime>()
  showtimes.value.forEach((s) => map.set(s.showtime_id, s))
  return map
})

const movieById = computed(() => {
  const map = new Map<string, Movie>()
  movies.value.forEach((m) => map.set(m.movie_id, m))
  return map
})

const theaterById = computed(() => {
  const map = new Map<string, Theater>()
  theaters.value.forEach((t) => map.set(t.theater_id, t))
  return map
})

const userById = computed(() => {
  const map = new Map<string, User>()
  users.value.forEach((u) => map.set(u.user_id, u))
  return map
})

function userLabel(userId: string) {
  const user = userById.value.get(userId)
  return user ? (user.email || user.name) : userId
}

function movieTitle(showtimeId: string) {
  const showtime = showtimeById.value.get(showtimeId)
  if (!showtime) return 'Unknown showtime'

  return movieById.value.get(showtime.movie_id)?.title ?? 'Unknown movie'
}

function theaterName(showtimeId: string) {
  const showtime = showtimeById.value.get(showtimeId)
  if (!showtime) return '—'

  return theaterById.value.get(showtime.theater_id)?.name ?? 'Unknown theater'
}

function showtimeLabel(showtimeId: string) {
  const showtime = showtimeById.value.get(showtimeId)
  if (!showtime) return '—'

  const date = new Date(showtime.show_date)
  const dateLabel = Number.isNaN(date.getTime())
    ? showtime.show_date
    : date.toLocaleDateString('en-US', { month: 'short', day: '2-digit' })

  return `${dateLabel} · ${showtime.start_time}`
}

function formatCreatedAt(isoString: string) {
  const date = new Date(isoString)
  if (Number.isNaN(date.getTime())) return isoString

  return date.toLocaleString('en-US', {
    month: 'short',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
}

// -----------------------------------------------------------------------
// Data loading
// -----------------------------------------------------------------------
async function loadAll() {
  loading.value = true
  errorMessage.value = ''

  try {
    const [bookingsRes, showtimesRes, moviesRes, theatersRes, usersRes] =
      await Promise.all([
        api.get<Booking[]>('/bookings/'),
        api.get<Showtime[]>('/showtimes'),
        api.get<Movie[]>('/movies'),
        api.get<Theater[]>('/theaters'),
        api.get<User[]>('/users'),
      ])

    bookings.value = bookingsRes.data
    showtimes.value = showtimesRes.data
    movies.value = moviesRes.data
    theaters.value = theatersRes.data
    users.value = usersRes.data
  } catch (err) {
    errorMessage.value = 'Failed to load bookings. Please try again.'
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)

const filteredBookings = computed(() => {
  const keyword = search.value.toLowerCase().trim()

  return bookings.value.filter((booking) => {
    const matchesSearch =
      !keyword ||
      booking.booking_id.toLowerCase().includes(keyword) ||
      userLabel(booking.user_id).toLowerCase().includes(keyword)

    const matchesStatus =
      activeFilter.value === 'all' ||
      booking.status === activeFilter.value

    return matchesSearch && matchesStatus
  })
})

const statusLabel = (status: BookingStatus) => {
  const labels = {
    pending: 'pending',
    confirmed: 'confirmed',
    cancelled: 'cancelled',
  }

  return labels[status]
}
</script>

<template>
  <div class="mx-auto max-w-7xl px-6 py-8 lg:px-8">

    <!-- Header -->
    <div class="mb-8">
      <div class="flex items-center gap-3">
        <div
          class="flex h-10 w-10 items-center justify-center rounded-xl
                 bg-gradient-to-br from-pink-500 to-purple-500"
        >
          <Ticket class="h-5 w-5 text-white" />
        </div>

        <div>
          <h1 class="text-2xl font-bold text-white">
            Manage Bookings
          </h1>

          <p class="mt-1 text-sm text-gray-400">
            Manage and monitor customer bookings
          </p>
        </div>
      </div>
    </div>

    <!-- Error banner -->
    <div
      v-if="errorMessage"
      class="mb-6 rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400"
    >
      {{ errorMessage }}
    </div>

    <!-- Search + Filters -->
    <div
      class="mb-5 flex flex-col gap-4 lg:flex-row lg:items-center
             lg:justify-between"
    >

      <!-- Search -->
      <div class="relative w-full lg:w-[320px]">
        <Search
          class="absolute left-3 top-1/2 h-4 w-4
                 -translate-y-1/2 text-gray-500"
        />

        <input
          v-model="search"
          type="text"
          placeholder="Search by booking ID or user..."
          class="h-11 w-full rounded-xl border border-white/10
                 bg-[#181820] pl-10 pr-4 text-sm text-white
                 outline-none transition
                 placeholder:text-gray-500
                 focus:border-pink-500/50"
        />
      </div>

      <!-- Filters -->
      <div class="flex gap-2 overflow-x-auto">

        <button
          type="button"
          class="rounded-xl border px-5 py-2.5 text-sm font-medium
                 transition"
          :class="
            activeFilter === 'all'
              ? 'border-pink-500 bg-pink-500 text-white'
              : 'border-white/10 bg-[#15151d] text-gray-300 hover:bg-white/5'
          "
          @click="activeFilter = 'all'"
        >
          All
        </button>

        <button
          type="button"
          class="rounded-xl border px-5 py-2.5 text-sm font-medium
                 transition"
          :class="
            activeFilter === 'pending'
              ? 'border-pink-500 bg-pink-500 text-white'
              : 'border-white/10 bg-[#15151d] text-gray-300 hover:bg-white/5'
          "
          @click="activeFilter = 'pending'"
        >
          Pending
        </button>

        <button
          type="button"
          class="rounded-xl border px-5 py-2.5 text-sm font-medium
                 transition"
          :class="
            activeFilter === 'confirmed'
              ? 'border-pink-500 bg-pink-500 text-white'
              : 'border-white/10 bg-[#15151d] text-gray-300 hover:bg-white/5'
          "
          @click="activeFilter = 'confirmed'"
        >
          Confirmed
        </button>

        <button
          type="button"
          class="rounded-xl border px-5 py-2.5 text-sm font-medium
                 transition"
          :class="
            activeFilter === 'cancelled'
              ? 'border-pink-500 bg-pink-500 text-white'
              : 'border-white/10 bg-[#15151d] text-gray-300 hover:bg-white/5'
          "
          @click="activeFilter = 'cancelled'"
        >
          Cancelled
        </button>

      </div>
    </div>

    <!-- Booking Table -->
    <div class="overflow-hidden rounded-2xl border border-white/10 bg-[#15151d]">

      <!-- Loading -->
      <div
        v-if="loading"
        class="px-5 py-16 text-center text-sm text-gray-400"
      >
        Loading bookings...
      </div>

      <!-- Desktop table -->
      <div v-else class="hidden overflow-x-auto lg:block">

        <table class="w-full border-collapse">

          <thead>
            <tr class="border-b border-white/10 text-left">

              <th class="px-5 py-4 text-xs font-semibold uppercase
                         tracking-wider text-gray-400">
                Booking ID
              </th>

              <th class="px-5 py-4 text-xs font-semibold uppercase
                         tracking-wider text-gray-400">
                User
              </th>

              <th class="px-5 py-4 text-xs font-semibold uppercase
                         tracking-wider text-gray-400">
                Movie
              </th>

              <th class="px-5 py-4 text-xs font-semibold uppercase
                         tracking-wider text-gray-400">
                Theater
              </th>

              <th class="px-5 py-4 text-xs font-semibold uppercase
                         tracking-wider text-gray-400">
                Showtime
              </th>

              <th class="px-5 py-4 text-xs font-semibold uppercase
                         tracking-wider text-gray-400">
                Seats
              </th>

              <th class="px-5 py-4 text-xs font-semibold uppercase
                         tracking-wider text-gray-400">
                Total
              </th>

              <th class="px-5 py-4 text-xs font-semibold uppercase
                         tracking-wider text-gray-400">
                Status
              </th>

              <th class="px-5 py-4 text-xs font-semibold uppercase
                         tracking-wider text-gray-400">
                Created At
              </th>

            </tr>
          </thead>

          <tbody>

            <tr
              v-for="booking in filteredBookings"
              :key="booking.booking_id"
              class="border-b border-white/5 transition
                     last:border-b-0 hover:bg-white/[0.025]"
            >

              <!-- Booking ID -->
              <td class="px-5 py-5">
                <span
                  class="block max-w-[115px] truncate text-sm
                         font-medium text-white"
                  :title="booking.booking_id"
                >
                  {{ booking.booking_id }}
                </span>
              </td>

              <!-- User -->
              <td class="px-5 py-5">
                <span class="text-sm text-gray-300">
                  {{ userLabel(booking.user_id) }}
                </span>
              </td>

              <!-- Movie -->
              <td class="max-w-[180px] px-5 py-5">
                <span class="block text-sm font-medium text-white">
                  {{ movieTitle(booking.showtime_id) }}
                </span>
              </td>

              <!-- Theater -->
              <td class="max-w-[150px] px-5 py-5">
                <span class="block text-sm text-gray-300">
                  {{ theaterName(booking.showtime_id) }}
                </span>
              </td>

              <!-- Showtime -->
              <td class="px-5 py-5">
                <div class="flex items-center gap-2 whitespace-nowrap">
                  <Clock3 class="h-4 w-4 text-gray-500" />

                  <span class="text-sm text-gray-300">
                    {{ showtimeLabel(booking.showtime_id) }}
                  </span>
                </div>
              </td>

              <!-- Seats -->
              <td class="px-5 py-5">
                <div class="flex items-center gap-2">
                  <Users class="h-4 w-4 text-gray-500" />

                  <span class="text-sm text-gray-300">
                    {{ booking.seat_numbers.join(', ') }}
                  </span>
                </div>
              </td>

              <!-- Total -->
              <td class="px-5 py-5">
                <span class="whitespace-nowrap text-sm font-semibold
                             text-white">
                  ฿{{ booking.total_price.toFixed(2) }}
                </span>
              </td>

              <!-- Status -->
              <td class="px-5 py-5">

                <span
                  class="inline-flex rounded-full px-3 py-1
                         text-xs font-semibold"
                  :class="{
                    'bg-emerald-500/15 text-emerald-400':
                      booking.status === 'confirmed',

                    'bg-amber-500/15 text-amber-400':
                      booking.status === 'pending',

                    'bg-pink-500/15 text-pink-400':
                      booking.status === 'cancelled',
                  }"
                >
                  {{ statusLabel(booking.status) }}
                </span>

              </td>

              <!-- Created -->
              <td class="px-5 py-5">
                <div class="flex items-center gap-2 whitespace-nowrap">
                  <CalendarDays class="h-4 w-4 text-gray-500" />

                  <span class="text-sm text-gray-400">
                    {{ formatCreatedAt(booking.created_at) }}
                  </span>
                </div>
              </td>

            </tr>

            <!-- Empty -->
            <tr v-if="filteredBookings.length === 0">
              <td
                colspan="9"
                class="px-5 py-16 text-center"
              >
                <Ticket
                  class="mx-auto mb-3 h-8 w-8 text-gray-600"
                />

                <p class="text-sm font-medium text-gray-400">
                  No bookings found
                </p>

                <p class="mt-1 text-xs text-gray-600">
                  Try another search or filter
                </p>
              </td>
            </tr>

          </tbody>

        </table>

      </div>

      <!-- Mobile -->
      <div v-if="!loading" class="divide-y divide-white/5 lg:hidden">

        <div
          v-for="booking in filteredBookings"
          :key="booking.booking_id"
          class="p-5"
        >

          <div class="mb-4 flex items-start justify-between gap-3">

            <div>
              <p class="text-xs text-gray-500">
                BOOKING ID
              </p>

              <p class="mt-1 text-sm font-medium text-white">
                {{ booking.booking_id }}
              </p>
            </div>

            <span
              class="rounded-full px-3 py-1 text-xs font-semibold"
              :class="{
                'bg-emerald-500/15 text-emerald-400':
                  booking.status === 'confirmed',

                'bg-amber-500/15 text-amber-400':
                  booking.status === 'pending',

                'bg-pink-500/15 text-pink-400':
                  booking.status === 'cancelled',
              }"
            >
              {{ statusLabel(booking.status) }}
            </span>

          </div>

          <div class="space-y-3">

            <div>
              <p class="text-xs text-gray-500">
                USER
              </p>

              <p class="mt-1 text-sm text-gray-300">
                {{ userLabel(booking.user_id) }}
              </p>
            </div>

            <div>
              <p class="text-xs text-gray-500">
                MOVIE
              </p>

              <p class="mt-1 text-sm font-medium text-white">
                {{ movieTitle(booking.showtime_id) }}
              </p>
            </div>

            <div class="grid grid-cols-2 gap-4">

              <div>
                <p class="text-xs text-gray-500">
                  THEATER
                </p>

                <p class="mt-1 text-sm text-gray-300">
                  {{ theaterName(booking.showtime_id) }}
                </p>
              </div>

              <div>
                <p class="text-xs text-gray-500">
                  SHOWTIME
                </p>

                <p class="mt-1 text-sm text-gray-300">
                  {{ showtimeLabel(booking.showtime_id) }}
                </p>
              </div>

            </div>

            <div class="grid grid-cols-2 gap-4">

              <div>
                <p class="text-xs text-gray-500">
                  SEATS
                </p>

                <p class="mt-1 text-sm text-gray-300">
                  {{ booking.seat_numbers.join(', ') }}
                </p>
              </div>

              <div>
                <p class="text-xs text-gray-500">
                  TOTAL
                </p>

                <p class="mt-1 text-sm font-semibold text-white">
                  ฿{{ booking.total_price.toFixed(2) }}
                </p>
              </div>

            </div>

          </div>

        </div>

        <div
          v-if="filteredBookings.length === 0"
          class="px-5 py-16 text-center"
        >
          <Ticket class="mx-auto mb-3 h-8 w-8 text-gray-600" />

          <p class="text-sm font-medium text-gray-400">
            No bookings found
          </p>
        </div>

      </div>

    </div>

    <!-- Result count -->
    <div class="mt-4 flex items-center justify-between">

      <p class="text-xs text-gray-500">
        Showing {{ filteredBookings.length }}
        {{ filteredBookings.length === 1 ? 'booking' : 'bookings' }}
      </p>

    </div>

  </div>
</template>