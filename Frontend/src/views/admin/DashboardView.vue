<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import api from '@/services/api'

// -----------------------------------------------------------------------
// Types
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
  showtime_id: string
  movie_id: string
  theater_id: string
  show_date: string
  start_time: string
}

interface Movie {
  movie_id: string
  title: string
}

interface Theater {
  theater_id: string
  name: string
}

interface User {
  user_id: string
  name: string
  email: string
}

// -----------------------------------------------------------------------
// State
// -----------------------------------------------------------------------

const bookings = ref<Booking[]>([])
const showtimes = ref<Showtime[]>([])
const movies = ref<Movie[]>([])
const theaters = ref<Theater[]>([])
const users = ref<User[]>([])

const loading = ref(false)
const errorMessage = ref('')

// -----------------------------------------------------------------------
// Helper
// Safely normalize API responses into arrays.
// Supports:
//   []
//   { data: [] }
//   { items: [] }
//   { results: [] }
// -----------------------------------------------------------------------

function toArray<T>(payload: unknown): T[] {
  if (Array.isArray(payload)) {
    return payload as T[]
  }

  if (payload && typeof payload === 'object') {
    const obj = payload as Record<string, unknown>

    if (Array.isArray(obj.data)) {
      return obj.data as T[]
    }

    if (Array.isArray(obj.items)) {
      return obj.items as T[]
    }

    if (Array.isArray(obj.results)) {
      return obj.results as T[]
    }
  }

  console.warn('Expected array but got:', payload)

  return []
}

// -----------------------------------------------------------------------
// Lookups
// -----------------------------------------------------------------------

const showtimeById = computed(() => {
  const map = new Map<string, Showtime>()

  showtimes.value.forEach((showtime) => {
    map.set(showtime.showtime_id, showtime)
  })

  return map
})

const movieById = computed(() => {
  const map = new Map<string, Movie>()

  movies.value.forEach((movie) => {
    map.set(movie.movie_id, movie)
  })

  return map
})

const theaterById = computed(() => {
  const map = new Map<string, Theater>()

  theaters.value.forEach((theater) => {
    map.set(theater.theater_id, theater)
  })

  return map
})

const userById = computed(() => {
  const map = new Map<string, User>()

  users.value.forEach((user) => {
    map.set(user.user_id, user)
  })

  return map
})

// -----------------------------------------------------------------------
// Display helpers
// -----------------------------------------------------------------------

function userLabel(userId: string) {
  const user = userById.value.get(userId)

  if (!user) {
    return userId || 'Unknown user'
  }

  return user.email || user.name || user.user_id
}

function movieTitle(showtimeId: string) {
  const showtime = showtimeById.value.get(showtimeId)

  if (!showtime) {
    return 'Unknown movie'
  }

  return (
    movieById.value.get(showtime.movie_id)?.title ??
    'Unknown movie'
  )
}

function showtimeLabel(showtimeId: string) {
  const showtime = showtimeById.value.get(showtimeId)

  if (!showtime) {
    return '—'
  }

  const date = new Date(showtime.show_date)

  const dateLabel = Number.isNaN(date.getTime())
    ? showtime.show_date
    : date.toLocaleDateString('en-US', {
        month: 'short',
        day: '2-digit',
      })

  return `${dateLabel} ${showtime.start_time}`
}

function formatPrice(price: number) {
  return `฿${Number(price || 0).toFixed(2)}`
}

// -----------------------------------------------------------------------
// Data loading
// -----------------------------------------------------------------------

async function loadAll() {
  loading.value = true
  errorMessage.value = ''

  try {
    const [
      bookingsRes,
      showtimesRes,
      moviesRes,
      theatersRes,
      usersRes,
    ] = await Promise.all([
      // IMPORTANT:
      // Admin Dashboard must use the admin booking endpoint
      // so it can see all bookings.
      api.get('/admin/bookings'),

      api.get('/showtimes'),
      api.get('/movies'),
      api.get('/theaters'),
      api.get('/users'),
    ])

    bookings.value = toArray<Booking>(bookingsRes.data)

    showtimes.value = toArray<Showtime>(showtimesRes.data)

    movies.value = toArray<Movie>(moviesRes.data)

    theaters.value = toArray<Theater>(theatersRes.data)

    users.value = toArray<User>(usersRes.data)

    console.log('========== ADMIN DASHBOARD ==========')
    console.log('Bookings:', bookings.value)
    console.log('Showtimes:', showtimes.value)
    console.log('Movies:', movies.value)
    console.log('Theaters:', theaters.value)
    console.log('Users:', users.value)
  } catch (err) {
    errorMessage.value =
      'Failed to load dashboard data. Please try again.'

    console.error('========== DASHBOARD ERROR ==========')
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)

// -----------------------------------------------------------------------
// Dashboard statistics
// -----------------------------------------------------------------------

const totalBookings = computed(() => {
  return bookings.value.length
})

const confirmedBookings = computed(() => {
  return bookings.value.filter(
    (booking) => booking.status === 'confirmed',
  ).length
})

const pendingBookings = computed(() => {
  return bookings.value.filter(
    (booking) => booking.status === 'pending',
  ).length
})

const cancelledBookings = computed(() => {
  return bookings.value.filter(
    (booking) => booking.status === 'cancelled',
  ).length
})

const stats = computed(() => {
  return [
    {
      label: 'TOTAL BOOKINGS',
      value: totalBookings.value.toLocaleString('en-US'),
    },
    {
      label: 'CONFIRMED',
      value: confirmedBookings.value.toLocaleString('en-US'),
    },
    {
      label: 'PENDING PAYMENT',
      value: pendingBookings.value.toLocaleString('en-US'),
    },
    {
      label: 'CANCELLED',
      value: cancelledBookings.value.toLocaleString('en-US'),
    },
  ]
})

// -----------------------------------------------------------------------
// Recent bookings
// Latest 5 bookings by created_at
// -----------------------------------------------------------------------

const recentBookings = computed(() => {
  return [...bookings.value]
    .sort((a, b) => {
      const dateA = new Date(a.created_at).getTime()
      const dateB = new Date(b.created_at).getTime()

      return dateB - dateA
    })
    .slice(0, 5)
    .map((booking) => {
      return {
        id: booking.booking_id,
        user: userLabel(booking.user_id),
        movie: movieTitle(booking.showtime_id),
        showtime: showtimeLabel(booking.showtime_id),
        total: formatPrice(booking.total_price),
        status: booking.status,
      }
    })
})
</script>

<template>
  <div class="min-h-screen bg-[#0d0d12] px-8 py-7">
    <div class="w-full">

      <!-- ============================================================= -->
      <!-- ERROR -->
      <!-- ============================================================= -->

      <div
        v-if="errorMessage"
        class="mb-6 rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400"
      >
        {{ errorMessage }}
      </div>

      <!-- ============================================================= -->
      <!-- STAT CARDS -->
      <!-- ============================================================= -->

      <section
        class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4"
      >
        <div
          v-for="stat in stats"
          :key="stat.label"
          class="h-[162px] rounded-xl border border-white/10 bg-[#15151d] px-5 py-8"
        >
          <p
            class="text-xs font-medium tracking-wide text-[#8190ad]"
          >
            {{ stat.label }}
          </p>

          <p
            class="mt-7 text-3xl font-black tracking-tight text-white"
          >
            {{ loading ? '—' : stat.value }}
          </p>
        </div>
      </section>

      <!-- ============================================================= -->
      <!-- RECENT BOOKINGS -->
      <!-- ============================================================= -->

      <section class="mt-9">

        <h2
          class="mb-4 text-xs font-medium tracking-wide text-[#8190ad]"
        >
          RECENT BOOKINGS
        </h2>

        <div
          class="overflow-hidden rounded-xl border border-white/10 bg-[#15151d]"
        >
          <div class="overflow-x-auto">

            <table class="w-full min-w-[900px]">

              <!-- ===================================================== -->
              <!-- TABLE HEADER -->
              <!-- ===================================================== -->

              <thead>
                <tr
                  class="border-b border-white/10 bg-[#1b1b24]"
                >
                  <th
                    class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]"
                  >
                    BOOKING ID
                  </th>

                  <th
                    class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]"
                  >
                    USER
                  </th>

                  <th
                    class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]"
                  >
                    MOVIE
                  </th>

                  <th
                    class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]"
                  >
                    SHOWTIME
                  </th>

                  <th
                    class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]"
                  >
                    TOTAL
                  </th>

                  <th
                    class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]"
                  >
                    STATUS
                  </th>
                </tr>
              </thead>

              <!-- ===================================================== -->
              <!-- TABLE BODY -->
              <!-- ===================================================== -->

              <tbody>

                <!-- Loading -->

                <template v-if="loading">
                  <tr>
                    <td
                      colspan="6"
                      class="px-5 py-12 text-center text-sm text-gray-500"
                    >
                      Loading recent bookings...
                    </td>
                  </tr>
                </template>

                <!-- Data -->

                <template v-else-if="recentBookings.length > 0">

                  <tr
                    v-for="booking in recentBookings"
                    :key="booking.id"
                    class="border-b border-white/10 last:border-b-0"
                  >

                    <!-- Booking ID -->

                    <td
                      class="px-5 py-4 text-sm font-medium text-white"
                    >
                      <span
                        class="block max-w-[140px] truncate"
                        :title="booking.id"
                      >
                        {{ booking.id }}
                      </span>
                    </td>

                    <!-- User -->

                    <td
                      class="px-5 py-4 text-sm text-gray-200"
                    >
                      {{ booking.user }}
                    </td>

                    <!-- Movie -->

                    <td
                      class="px-5 py-4 text-sm text-gray-200"
                    >
                      {{ booking.movie }}
                    </td>

                    <!-- Showtime -->

                    <td
                      class="px-5 py-4 text-sm text-gray-200"
                    >
                      {{ booking.showtime }}
                    </td>

                    <!-- Total -->

                    <td
                      class="px-5 py-4 text-sm font-medium text-white"
                    >
                      {{ booking.total }}
                    </td>

                    <!-- Status -->

                    <td class="px-5 py-4">

                      <!-- Confirmed -->

                      <span
                        v-if="booking.status === 'confirmed'"
                        class="inline-flex rounded-full bg-emerald-500/15 px-3 py-1 text-xs font-semibold text-emerald-400"
                      >
                        confirmed
                      </span>

                      <!-- Pending -->

                      <span
                        v-else-if="booking.status === 'pending'"
                        class="inline-flex rounded-full bg-orange-500/15 px-3 py-1 text-xs font-semibold text-orange-400"
                      >
                        pending
                      </span>

                      <!-- Cancelled -->

                      <span
                        v-else
                        class="inline-flex rounded-full bg-rose-500/15 px-3 py-1 text-xs font-semibold text-rose-400"
                      >
                        cancelled
                      </span>

                    </td>

                  </tr>

                </template>

                <!-- Empty -->

                <template v-else>
                  <tr>
                    <td
                      colspan="6"
                      class="px-5 py-12 text-center text-sm text-gray-500"
                    >
                      No bookings yet.
                    </td>
                  </tr>
                </template>

              </tbody>

            </table>

          </div>
        </div>

      </section>

    </div>
  </div>
</template>