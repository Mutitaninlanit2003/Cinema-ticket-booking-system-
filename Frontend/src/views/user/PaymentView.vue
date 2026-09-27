<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/services/api'

// -----------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------

type BookingStatus = 'pending' | 'confirmed' | 'cancelled'

type PaymentMethod =
  | 'cash'
  | 'credit_card'
  | 'promptpay'
  | 'transfer'

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

interface Showtime {
  id: number
  showtime_id: string
  movie_id: string
  theater_id: string
  screen_id?: string | null
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

interface Seat {
  id: number
  seat_id: string
  screen_id: string
  seat_number: string
  seat_type: 'standard' | 'premium' | 'vip'
  status: string
}

interface Payment {
  id: number
  payment_id: string
  booking_id: string
  amount: number
  method: PaymentMethod
  status: string
  transaction_ref: string
}

// -----------------------------------------------------------------------
// Router
// -----------------------------------------------------------------------

const route = useRoute()
const router = useRouter()

const bookingId = computed(() => {
  return String(route.params.bookingId)
})

// -----------------------------------------------------------------------
// State
// -----------------------------------------------------------------------

const booking = ref<Booking | null>(null)
const showtime = ref<Showtime | null>(null)
const movie = ref<Movie | null>(null)
const theater = ref<Theater | null>(null)
const screen = ref<Screen | null>(null)
const seats = ref<Seat[]>([])

const paymentMethod = ref<PaymentMethod>('credit_card')

// Card fields are purely a visual "fake terminal" — the backend
// Payment model only has amount/method/status/transaction_ref, so
// none of these three are ever sent anywhere. They exist only to
// match the look of a real card form when "Credit card" is chosen.
const cardNumber = ref('')
const cardExpiry = ref('')
const cardCvv = ref('')

// Backend requires transaction_ref to be alphanumeric, 6–20 chars
// (models/payment.go: binding:"omitempty,alphanum,min=6,max=20").
// Pre-fill a plausible reference so the field isn't empty by
// default, but leave it editable.
const transactionRef = ref(generateTransactionRef())

function generateTransactionRef() {
  const random = Math.random().toString(36).slice(2, 8).toUpperCase()
  return `TXN${random}`
}

const loading = ref(false)
const payingLoading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

// -----------------------------------------------------------------------
// Load booking information
// -----------------------------------------------------------------------

async function loadPaymentData() {
  loading.value = true
  errorMessage.value = ''

  try {
    const [
      bookingsRes,
      showtimesRes,
      moviesRes,
      theatersRes,
      screensRes,
      seatsRes,
    ] = await Promise.all([
      api.get<Booking[]>('/bookings'),
      api.get<Showtime[]>('/showtimes'),
      api.get<Movie[]>('/movies'),
      api.get<Theater[]>('/theaters'),
      api.get<Screen[]>('/screens'),
      api.get<Seat[]>('/seats'),
    ])

    const foundBooking = bookingsRes.data.find(
      (item) => item.booking_id === bookingId.value,
    )

    if (!foundBooking) {
      errorMessage.value = 'Booking not found.'
      return
    }

    booking.value = foundBooking
    seats.value = seatsRes.data

    showtime.value =
      showtimesRes.data.find(
        (item) =>
          item.showtime_id === foundBooking.showtime_id,
      ) ?? null

    if (showtime.value) {
      movie.value =
        moviesRes.data.find(
          (item) =>
            item.movie_id === showtime.value?.movie_id,
        ) ?? null

      theater.value =
        theatersRes.data.find(
          (item) =>
            item.theater_id === showtime.value?.theater_id,
        ) ?? null

      screen.value =
        screensRes.data.find(
          (item) =>
            item.screen_id === showtime.value?.screen_id,
        ) ?? null
    }

    // Payment page should only be used for pending bookings.
    if (foundBooking.status !== 'pending') {
      errorMessage.value =
        'This booking is no longer available for payment.'
    }
  } catch (err) {
    console.error('LOAD PAYMENT DATA ERROR:', err)

    errorMessage.value =
      'Failed to load payment information. Please try again.'
  } finally {
    loading.value = false
  }
}

onMounted(loadPaymentData)

// -----------------------------------------------------------------------
// Derived display data
// -----------------------------------------------------------------------

// Seat type label(s) for the booked seats — shown next to "Seats" the
// same way the reference design does. If every booked seat shares the
// same type, show it once ("C7, C8 (standard)"); if they're mixed,
// show each seat with its own type instead.
const seatsLabel = computed(() => {
  if (!booking.value) return ''

  const seatNumbers = booking.value.seat_numbers
  const screenId = showtime.value?.screen_id

  const typesBySeat = seatNumbers.map((seatNumber) => {
    const match = seats.value.find(
      (seat) =>
        seat.seat_number === seatNumber &&
        (!screenId || seat.screen_id === screenId),
    )

    return { seatNumber, type: match?.seat_type ?? null }
  })

  const knownTypes = new Set(
    typesBySeat.map((s) => s.type).filter((t): t is string => !!t),
  )

  if (knownTypes.size === 0) {
    return seatNumbers.join(', ')
  }

  if (knownTypes.size === 1) {
    const [onlyType] = [...knownTypes]
    return `${seatNumbers.join(', ')} (${onlyType})`
  }

  return typesBySeat
    .map((s) => (s.type ? `${s.seatNumber} (${s.type})` : s.seatNumber))
    .join(', ')
})

// Price per seat, derived from the booking's own total so the
// "N × price" breakdown always adds up to exactly the Total price
// row below it — even if the showtime's price changed after this
// booking was made.
const pricePerSeat = computed(() => {
  if (!booking.value || booking.value.seat_numbers.length === 0) return 0
  return booking.value.total_price / booking.value.seat_numbers.length
})

const statusBadgeClass = computed(() => {
  switch (booking.value?.status) {
    case 'confirmed':
      return 'bg-emerald-500/15 text-emerald-400'
    case 'cancelled':
      return 'bg-red-500/15 text-red-400'
    default:
      return 'bg-amber-500/15 text-amber-400'
  }
})

const paymentMethods: { value: PaymentMethod; label: string }[] = [
  { value: 'credit_card', label: 'Credit card' },
  { value: 'promptpay', label: 'PromptPay' },
  { value: 'cash', label: 'Cash' },
  { value: 'transfer', label: 'Transfer' },
]

// -----------------------------------------------------------------------
// Formatting
// -----------------------------------------------------------------------

function formatPrice(price: number) {
  return `฿${price.toFixed(2)}`
}

function formatShowDate(dateStr?: string) {
  if (!dateStr) return '—'

  return new Date(dateStr).toLocaleDateString('en-US', {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

// -----------------------------------------------------------------------
// Payment
// -----------------------------------------------------------------------

async function submitPayment() {
  if (!booking.value) {
    errorMessage.value = 'Booking information is missing.'
    return
  }

  if (booking.value.status !== 'pending') {
    errorMessage.value =
      'This booking is no longer available for payment.'
    return
  }

  if (transactionRef.value.length < 6 || transactionRef.value.length > 20) {
    errorMessage.value =
      'Transaction reference must be 6–20 characters.'
    return
  }

  payingLoading.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const paymentData = {
      booking_id: booking.value.booking_id,
      amount: booking.value.total_price,
      method: paymentMethod.value,
      status: 'paid',
      transaction_ref: transactionRef.value,
    }

    const response = await api.post<Payment>(
      '/payments',
      paymentData,
    )

    successMessage.value =
      'Payment completed successfully.'

    // Give the user a short moment to see success.
    setTimeout(() => {
      router.push('/bookings')
    }, 800)
  } catch (err: any) {
    console.error('PAYMENT ERROR:', err)

    errorMessage.value =
      err.response?.data?.error ||
      err.response?.data?.message ||
      'Failed to complete payment. Please try again.'
  } finally {
    payingLoading.value = false
  }
}

// -----------------------------------------------------------------------
// Back
// -----------------------------------------------------------------------

function goBack() {
  router.push('/bookings')
}
</script>

<template>
  <div class="min-h-screen bg-[#0a0a0f] px-6 py-10 text-white">
    <div class="mx-auto max-w-5xl">

      <!-- Header -->
      <div class="mb-8">
        <button
          type="button"
          class="mb-5 text-sm text-gray-400 transition hover:text-white"
          @click="goBack"
        >
          ← Back to My Bookings
        </button>

        <h1 class="text-3xl font-bold">
          Payment
        </h1>

        <p class="mt-2 text-sm text-gray-500">
          Complete your payment to confirm your booking.
        </p>
      </div>

      <!-- Loading -->
      <div
        v-if="loading"
        class="rounded-2xl border border-white/10 bg-[#15151d] p-10 text-center"
      >
        <p class="text-sm text-gray-400">
          Loading payment information…
        </p>
      </div>

      <!-- Error -->
      <div
        v-else-if="errorMessage && !booking"
        class="rounded-2xl border border-red-500/30 bg-red-500/10 p-6"
      >
        <p class="text-sm text-red-400">
          {{ errorMessage }}
        </p>

        <button
          type="button"
          class="mt-5 rounded-xl border border-white/15 px-5 py-2.5 text-sm font-semibold transition hover:bg-white hover:text-black"
          @click="goBack"
        >
          Back to My Bookings
        </button>
      </div>

      <!-- Payment page -->
      <div
        v-else-if="booking"
        class="grid overflow-hidden rounded-2xl border border-white/10 bg-[#15151d] lg:grid-cols-2"
      >

        <!-- Booking Summary -->
        <div class="border-b border-white/10 p-6 lg:border-b-0 lg:border-r">
          <p class="text-xs font-semibold uppercase tracking-wide text-gray-500">
            Booking summary
          </p>

          <!-- Booking ID + status -->
          <div class="mt-4 flex items-center gap-2 text-sm text-gray-400">
            <span>Booking ID · {{ booking.booking_id.slice(0, 8) }}...{{ booking.booking_id.slice(-4) }}</span>

            <span
              class="rounded-full px-2.5 py-0.5 text-xs font-semibold"
              :class="statusBadgeClass"
            >
              {{ booking.status }}
            </span>
          </div>

          <!-- Movie -->
          <div class="mt-5 flex gap-5">
            <div
              v-if="movie?.poster_url"
              class="h-28 w-20 shrink-0 overflow-hidden rounded-xl bg-cover bg-center"
              :style="{ backgroundImage: `url(${movie.poster_url})` }"
            ></div>

            <div
              v-else
              class="h-28 w-20 shrink-0 rounded-xl bg-gradient-to-b from-purple-950 to-black"
            ></div>

            <div class="min-w-0">
              <h3 class="text-lg font-bold">
                {{ movie?.title ?? 'Unknown movie' }}
              </h3>

              <p class="mt-2 text-sm text-gray-400">
                {{ theater?.name ?? 'Unknown theater' }}
                <template v-if="screen"> · {{ screen.name }}</template>
              </p>

              <p class="mt-1 text-sm text-gray-400">
                {{ formatShowDate(showtime?.show_date) }} · {{ showtime?.start_time ?? '—' }}
              </p>
            </div>
          </div>

          <!-- Seats / Tickets / Total -->
          <div class="mt-6 space-y-3 border-t border-white/10 pt-5">
            <div class="flex items-start justify-between gap-4 text-sm">
              <span class="text-gray-500">Seats</span>
              <span class="text-right font-medium">{{ seatsLabel }}</span>
            </div>

            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-500">Tickets</span>
              <span class="font-medium">
                {{ booking.seat_numbers.length }} × {{ formatPrice(pricePerSeat) }}
              </span>
            </div>
          </div>

          <div class="mt-5 flex items-center justify-between border-t border-white/10 pt-5">
            <span class="text-base font-bold">Total price</span>
            <span class="text-2xl font-black">{{ formatPrice(booking.total_price) }}</span>
          </div>
        </div>

        <!-- Payment Form -->
        <div class="p-6">

          <!-- Error / Success -->
          <div
            v-if="errorMessage"
            class="mb-5 rounded-lg border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400"
          >
            {{ errorMessage }}
          </div>

          <div
            v-if="successMessage"
            class="mb-5 rounded-lg border border-emerald-500/30 bg-emerald-500/10 px-4 py-3 text-sm text-emerald-400"
          >
            {{ successMessage }}
          </div>

          <!-- Payment method -->
          <div>
            <p class="text-xs font-semibold uppercase tracking-wide text-gray-500">
              Payment method
            </p>

            <div class="mt-3 grid grid-cols-2 gap-2 sm:grid-cols-4">
              <button
                v-for="option in paymentMethods"
                :key="option.value"
                type="button"
                :disabled="payingLoading || booking.status !== 'pending'"
                class="rounded-xl border px-3 py-2.5 text-sm font-semibold transition disabled:cursor-not-allowed disabled:opacity-50"
                :class="
                  paymentMethod === option.value
                    ? 'border-pink-500 bg-pink-500/15 text-pink-400'
                    : 'border-white/15 text-gray-300 hover:border-white/30'
                "
                @click="paymentMethod = option.value"
              >
                {{ option.label }}
              </button>
            </div>
          </div>

          <!-- Fake card fields — visual only, never sent to the backend
               (Payment model has no card fields at all) -->
          <div v-if="paymentMethod === 'credit_card'" class="mt-5 space-y-4">
            <div>
              <label class="mb-1.5 block text-xs font-medium text-gray-400">Card number</label>
              <input
                v-model="cardNumber"
                type="text"
                placeholder="4242 4242 4242 4242"
                :disabled="payingLoading || booking.status !== 'pending'"
                class="h-12 w-full rounded-xl border border-white/15 bg-[#0a0a0f] px-4 text-sm text-white outline-none transition focus:border-pink-500 disabled:opacity-50"
              />
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="mb-1.5 block text-xs font-medium text-gray-400">Expiry</label>
                <input
                  v-model="cardExpiry"
                  type="text"
                  placeholder="08 / 29"
                  :disabled="payingLoading || booking.status !== 'pending'"
                  class="h-12 w-full rounded-xl border border-white/15 bg-[#0a0a0f] px-4 text-sm text-white outline-none transition focus:border-pink-500 disabled:opacity-50"
                />
              </div>

              <div>
                <label class="mb-1.5 block text-xs font-medium text-gray-400">CVV</label>
                <input
                  v-model="cardCvv"
                  type="password"
                  placeholder="•••"
                  :disabled="payingLoading || booking.status !== 'pending'"
                  class="h-12 w-full rounded-xl border border-white/15 bg-[#0a0a0f] px-4 text-sm text-white outline-none transition focus:border-pink-500 disabled:opacity-50"
                />
              </div>
            </div>
          </div>

          <!-- Transaction reference -->
          <div class="mt-5">
            <label class="mb-1.5 block text-xs font-medium text-gray-400">
              Transaction reference
            </label>
            <input
              v-model="transactionRef"
              type="text"
              maxlength="20"
              :disabled="payingLoading || booking.status !== 'pending'"
              class="h-12 w-full rounded-xl border border-white/15 bg-[#0a0a0f] px-4 text-sm text-white outline-none transition focus:border-pink-500 disabled:opacity-50"
            />
            <p class="mt-1.5 text-xs text-gray-600">6–20 letters/numbers.</p>
          </div>

          <!-- Amount -->
          <div class="mt-5 rounded-xl bg-[#0a0a0f] p-4">
            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-400">Amount to pay</span>
              <span class="font-bold">{{ formatPrice(booking.total_price) }}</span>
            </div>
          </div>

          <!-- Confirm -->
          <button
            type="button"
            :disabled="
              payingLoading ||
              booking.status !== 'pending' ||
              !!successMessage
            "
            class="mt-6 w-full rounded-xl bg-pink-500 px-6 py-3.5 font-bold text-white transition hover:bg-pink-600 disabled:cursor-not-allowed disabled:opacity-50"
            @click="submitPayment"
          >
            {{
              payingLoading
                ? 'Processing...'
                : successMessage
                  ? 'Payment completed'
                  : `Confirm payment · ${formatPrice(booking.total_price)}`
            }}
          </button>

          <p class="mt-4 text-center text-xs text-gray-600">
            Your payment will confirm this booking.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>