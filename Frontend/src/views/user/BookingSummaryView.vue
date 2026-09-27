<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/services/api'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// -----------------------------------------------------------------------
// Incoming state from SeatSelectionView
// -----------------------------------------------------------------------

const showtimeId = computed(() => {
  return String(route.query.showtimeId ?? '')
})

const seatNumbers = computed(() => {
  const raw = route.query.seats

  if (typeof raw !== 'string' || raw === '') {
    return []
  }

  return raw.split(',')
})

// -----------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------

type BookingStatus =
  | 'pending'
  | 'confirmed'
  | 'cancelled'

interface Booking {
  id: number
  booking_id: string
  user_id: string
  showtime_id: string
  seat_numbers: string[]
  total_price: number
  status: BookingStatus
  created_at: string
  updated_at: string
}

type SeatType =
  | 'standard'
  | 'premium'
  | 'vip'

interface Seat {
  id: number
  seat_id: string
  screen_id: string
  seat_number: string
  seat_type: SeatType
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
  price_standard: number
  price_premium: number
  price_vip: number
  status: string
}

interface Movie {
  id: number
  movie_id: string
  title: string
  poster_url: string
}

interface Theater {
  id: number
  theater_id: string
  name: string
}

interface Screen {
  id: number
  screen_id: string
  theater_id: string
  name: string
}

// -----------------------------------------------------------------------
// State
// -----------------------------------------------------------------------

const showtimes = ref<Showtime[]>([])
const seats = ref<Seat[]>([])
const movies = ref<Movie[]>([])
const theaters = ref<Theater[]>([])
const screens = ref<Screen[]>([])

const loading = ref(false)
const submitting = ref(false)
const errorMessage = ref('')

// -----------------------------------------------------------------------
// Load booking details
// -----------------------------------------------------------------------

async function loadAll() {
  loading.value = true
  errorMessage.value = ''

  try {
    const [
      showtimesRes,
      seatsRes,
      moviesRes,
      theatersRes,
      screensRes,
    ] = await Promise.all([
      api.get<Showtime[]>('/showtimes'),
      api.get<Seat[]>('/seats'),
      api.get<Movie[]>('/movies'),
      api.get<Theater[]>('/theaters'),
      api.get<Screen[]>('/screens'),
    ])

    showtimes.value = showtimesRes.data
    seats.value = seatsRes.data
    movies.value = moviesRes.data
    theaters.value = theatersRes.data
    screens.value = screensRes.data
  } catch (err) {
    errorMessage.value =
      err instanceof Error
        ? err.message
        : 'Failed to load booking details. Please try again.'

    console.error(
      'Load booking details error:',
      err,
    )
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  authStore.initFromStorage()
  loadAll()
})

// -----------------------------------------------------------------------
// Find related data
// -----------------------------------------------------------------------

const showtime = computed(() => {
  return (
    showtimes.value.find(
      (item) => item.showtime_id === showtimeId.value,
    ) ?? null
  )
})

const movie = computed(() => {
  return (
    movies.value.find(
      (item) =>
        item.movie_id === showtime.value?.movie_id,
    ) ?? null
  )
})

const theater = computed(() => {
  return (
    theaters.value.find(
      (item) =>
        item.theater_id === showtime.value?.theater_id,
    ) ?? null
  )
})

const screen = computed(() => {
  return (
    screens.value.find(
      (item) =>
        item.screen_id === showtime.value?.screen_id,
    ) ?? null
  )
})

// -----------------------------------------------------------------------
// Selected seats resolved against the real Seat records
// -----------------------------------------------------------------------

const selectedSeatRecords = computed(() => {
  if (!showtime.value?.screen_id) {
    return []
  }

  return seats.value.filter(
    (seat) =>
      seat.screen_id === showtime.value!.screen_id &&
      seatNumbers.value.includes(seat.seat_number),
  )
})

function seatPrice(seat: Seat) {
  if (!showtime.value) {
    return 0
  }

  if (seat.seat_type === 'vip') {
    return showtime.value.price_vip
  }

  if (seat.seat_type === 'premium') {
    return showtime.value.price_premium
  }

  return showtime.value.price_standard
}

const SEAT_TYPE_LABELS: Record<
  SeatType,
  string
> = {
  standard: 'Standard',
  premium: 'Premium',
  vip: 'VIP',
}

// -----------------------------------------------------------------------
// Price breakdown
// -----------------------------------------------------------------------

const priceBreakdown = computed(() => {
  const groups = new Map<
    SeatType,
    {
      count: number
      unitPrice: number
    }
  >()

  for (const seat of selectedSeatRecords.value) {
    const existing = groups.get(seat.seat_type)

    if (existing) {
      existing.count += 1
    } else {
      groups.set(seat.seat_type, {
        count: 1,
        unitPrice: seatPrice(seat),
      })
    }
  }

  return (
    ['vip', 'premium', 'standard'] as SeatType[]
  )
    .filter((type) => groups.has(type))
    .map((type) => {
      const group = groups.get(type)!

      return {
        type,
        label: SEAT_TYPE_LABELS[type],
        count: group.count,
        unitPrice: group.unitPrice,
        subtotal:
          group.count * group.unitPrice,
      }
    })
})

// -----------------------------------------------------------------------
// Price calculation
// -----------------------------------------------------------------------

const totalPrice = computed(() => {
  if (selectedSeatRecords.value.length > 0) {
    return selectedSeatRecords.value.reduce(
      (sum, seat) => sum + seatPrice(seat),
      0,
    )
  }

  const fromQuery = Number(
    route.query.totalPrice,
  )

  return Number.isFinite(fromQuery)
    ? fromQuery
    : 0
})

function formatPrice(price: number) {
  return `฿${price.toFixed(2)}`
}

function formatShowDate(
  dateString: string | undefined,
) {
  if (!dateString) {
    return '—'
  }

  return new Date(
    dateString,
  ).toLocaleDateString('en-US', {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
  })
}

// -----------------------------------------------------------------------
// Get current logged-in User ID
// -----------------------------------------------------------------------

function getCurrentUserId(): string | null {
  authStore.initFromStorage()

  if (!authStore.isUserAuthenticated) {
    return null
  }

  return authStore.user?.user_id ?? null
}

// -----------------------------------------------------------------------
// Confirm booking
// -----------------------------------------------------------------------

async function confirmBooking() {
  if (
    !showtime.value ||
    seatNumbers.value.length === 0
  ) {
    return
  }

  authStore.initFromStorage()

  if (
    !authStore.isUserAuthenticated ||
    !authStore.user
  ) {
    errorMessage.value =
      'You need to be signed in to complete a booking.'

    return
  }

  const userId = authStore.user.user_id

  if (!userId) {
    errorMessage.value =
      'User information is missing. Please log in again.'

    return
  }

  submitting.value = true
  errorMessage.value = ''

  try {
    await api.post<Booking>('/bookings', {
      user_id: userId,
      showtime_id: showtimeId.value,
      seat_numbers: seatNumbers.value,
      total_price: totalPrice.value,
    })

    await router.push('/bookings')
  } catch (err) {
    errorMessage.value =
      err instanceof Error
        ? err.message
        : 'Failed to create the booking. Please try again.'

    console.error(
      'Create booking error:',
      err,
    )
  } finally {
    submitting.value = false
  }
}

// -----------------------------------------------------------------------
// Navigation
// -----------------------------------------------------------------------

function backToSeats() {
  router.push({
    path: `/seats/${showtimeId.value}`,
  })
}
</script>

<template>
  <div
    class="min-h-screen bg-[#0d0d12] px-4 py-8 text-white sm:px-6 sm:py-10"
  >
    <div class="mx-auto w-full max-w-3xl">

      <!-- Header -->
      <div class="mb-6 sm:mb-8">

        <p
          class="text-[10px] font-medium uppercase tracking-[0.2em] text-gray-500 sm:text-xs sm:tracking-[0.25em]"
        >
          Booking
        </p>

        <h1
          class="mt-2 text-2xl font-black sm:text-3xl"
        >
          Booking Summary
        </h1>

        <p
          class="mt-2 text-sm leading-6 text-gray-500"
        >
          ตรวจสอบรายละเอียดการจองก่อนยืนยัน
        </p>

      </div>

      <!-- Step -->
      <div
        class="mb-8 flex w-full items-center justify-between gap-1 text-[11px] sm:mb-10 sm:justify-start sm:gap-4 sm:text-sm"
      >

        <!-- Step 1 -->
        <div
          class="flex min-w-0 items-center gap-1.5 text-gray-500 sm:gap-2"
        >
          <span
            class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-white/10 text-[10px] sm:text-xs"
          >
            1
          </span>

          <span class="truncate">
            Showtime
          </span>
        </div>

        <div
          class="h-px min-w-4 flex-1 bg-white/20 sm:w-10 sm:flex-none"
        ></div>

        <!-- Step 2 -->
        <div
          class="flex min-w-0 items-center gap-1.5 text-gray-500 sm:gap-2"
        >
          <span
            class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-white/10 text-[10px] sm:text-xs"
          >
            2
          </span>

          <span class="truncate">
            Seats
          </span>
        </div>

        <div
          class="h-px min-w-4 flex-1 bg-white/20 sm:w-10 sm:flex-none"
        ></div>

        <!-- Step 3 -->
        <div
          class="flex min-w-0 items-center gap-1.5 font-semibold text-white sm:gap-2"
        >
          <span
            class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-[#ff3158] text-[10px] font-bold sm:text-xs"
          >
            3
          </span>

          <span class="truncate">
            Summary
          </span>
        </div>

      </div>

      <!-- Error message -->
      <div
        v-if="errorMessage"
        class="mb-5 break-words rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm leading-6 text-red-400 sm:mb-6"
      >
        {{ errorMessage }}
      </div>

      <!-- Loading -->
      <p
        v-if="loading"
        class="py-12 text-center text-sm text-gray-500"
      >
        Loading booking details…
      </p>

      <!-- Invalid booking data -->
      <p
        v-else-if="
          !showtime ||
          seatNumbers.length === 0
        "
        class="py-12 text-center text-sm leading-6 text-gray-500"
      >
        We couldn't find your seat selection. Please

        <button
          type="button"
          class="text-[#ff3158] underline underline-offset-2"
          @click="router.push('/')"
        >
          start over
        </button>

        .
      </p>

      <!-- Booking summary -->
      <template v-else>

        <div
          class="rounded-2xl border border-white/10 bg-[#121218] p-5 sm:p-8"
        >

          <!-- Movie information -->
          <div
            class="flex min-w-0 gap-4 sm:gap-5"
          >

            <!-- Poster -->
            <div
              class="h-28 w-20 shrink-0 overflow-hidden rounded-xl bg-gradient-to-b from-purple-950 to-black bg-cover bg-center sm:h-32 sm:w-[88px]"
              :style="
                movie?.poster_url
                  ? {
                      backgroundImage: `url(${movie.poster_url})`,
                    }
                  : {}
              "
            ></div>

            <!-- Movie Details -->
            <div
              class="min-w-0 flex-1"
            >

              <h2
                class="line-clamp-2 text-lg font-bold leading-6 sm:text-xl"
              >
                {{ movie?.title ?? 'Unknown movie' }}
              </h2>

              <p
                class="mt-2 break-words text-sm leading-5 text-gray-500"
              >
                {{ theater?.name ?? 'Unknown theater' }}

                <template v-if="screen">
                  · {{ screen.name }}
                </template>
              </p>

              <p
                class="mt-1 text-sm leading-5 text-gray-500"
              >
                {{ formatShowDate(showtime.show_date) }}

                <span class="mx-1">
                  ·
                </span>

                {{ showtime.start_time }}
              </p>

            </div>
          </div>

          <div
            class="my-5 h-px bg-white/10 sm:my-6"
          ></div>

          <!-- Selected seats -->
          <div>

            <p
              class="text-sm text-gray-500"
            >
              Selected seats
            </p>

            <p
              class="mt-2 break-words font-semibold leading-6 text-white"
            >
              {{ seatNumbers.join(', ') }}
            </p>

          </div>

          <div
            class="my-5 h-px bg-white/10 sm:my-6"
          ></div>

          <!-- Price Breakdown -->
          <div
            v-if="priceBreakdown.length"
            class="space-y-3"
          >

            <div
              v-for="group in priceBreakdown"
              :key="group.type"
              class="flex items-start justify-between gap-4 text-sm"
            >

              <span
                class="min-w-0 text-gray-400"
              >
                <span>
                  {{ group.label }} × {{ group.count }}
                </span>

                <span
                  class="ml-1 text-xs text-gray-600 sm:text-sm"
                >
                  ({{ formatPrice(group.unitPrice) }} each)
                </span>
              </span>

              <span
                class="shrink-0 font-semibold text-white"
              >
                {{ formatPrice(group.subtotal) }}
              </span>

            </div>

          </div>

          <p
            v-else
            class="text-sm text-gray-600"
          >
            Seats: {{ seatNumbers.length }}
          </p>

          <div
            class="my-5 h-px bg-white/10 sm:my-6"
          ></div>

          <!-- Total price -->
          <div
            class="flex items-center justify-between gap-4"
          >

            <span
              class="text-sm text-gray-400 sm:text-base"
            >
              Total
            </span>

            <span
              class="text-xl font-black sm:text-2xl"
            >
              {{ formatPrice(totalPrice) }}
            </span>

          </div>

        </div>

        <!-- Actions -->
        <div
          class="mt-5 flex flex-col gap-3 sm:mt-6 sm:flex-row-reverse"
        >

          <!-- Confirm -->
          <button
            type="button"
            :disabled="submitting"
            class="flex w-full flex-1 items-center justify-center rounded-xl bg-[#ff3158] py-3.5 text-sm font-bold text-white transition hover:bg-[#ff4568] disabled:cursor-not-allowed disabled:opacity-40 sm:w-auto sm:py-3"
            @click="confirmBooking"
          >
            {{
              submitting
                ? 'Confirming…'
                : 'Confirm booking'
            }}
          </button>

          <!-- Back -->
          <button
            type="button"
            :disabled="submitting"
            class="w-full rounded-xl border border-white/15 py-3.5 text-sm font-semibold text-gray-300 transition hover:bg-white/5 disabled:cursor-not-allowed disabled:opacity-40 sm:w-auto sm:px-6 sm:py-3"
            @click="backToSeats"
          >
            Back to seats
          </button>

        </div>

        <p
          class="mt-4 px-2 text-center text-xs leading-5 text-gray-600"
        >
          Your seats aren't held until you confirm — if
          someone else books one of your selected seats
          first, this step may fail.
        </p>

      </template>

    </div>
  </div>
</template>