<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
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

interface Payment {
  id: number
  payment_id: string
  booking_id: string
  amount: number
  method: PaymentMethod
  status: string
  transaction_ref: string
}

interface EnrichedBooking {
  booking: Booking
  showtime: Showtime | null
  movie: Movie | null
  theater: Theater | null
  screen: Screen | null
}

// -----------------------------------------------------------------------
// State
// -----------------------------------------------------------------------

const activeFilter = ref<'all' | BookingStatus>('all')

const bookings = ref<Booking[]>([])
const showtimes = ref<Showtime[]>([])
const movies = ref<Movie[]>([])
const theaters = ref<Theater[]>([])
const screens = ref<Screen[]>([])

const loading = ref(false)
const errorMessage = ref('')

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
      screensRes,
    ] = await Promise.all([
      api.get<Booking[]>('/bookings'),
      api.get<Showtime[]>('/showtimes'),
      api.get<Movie[]>('/movies'),
      api.get<Theater[]>('/theaters'),
      api.get<Screen[]>('/screens'),
    ])

    bookings.value = bookingsRes.data
    showtimes.value = showtimesRes.data
    movies.value = moviesRes.data
    theaters.value = theatersRes.data
    screens.value = screensRes.data
  } catch (err) {
    console.error('========== MY BOOKINGS ERROR ==========')
    console.error(err)

    errorMessage.value =
      'Failed to load bookings. Please try again.'
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)

// -----------------------------------------------------------------------
// Reload bookings
// -----------------------------------------------------------------------

async function reloadBookings() {
  try {
    const { data } = await api.get<Booking[]>('/bookings')
    bookings.value = data
  } catch (err) {
    console.error('Failed to reload bookings:', err)
  }
}

// -----------------------------------------------------------------------
// Find related data
// -----------------------------------------------------------------------

function findShowtime(showtimeId: string) {
  return (
    showtimes.value.find(
      (showtime) =>
        showtime.showtime_id === showtimeId,
    ) ?? null
  )
}

function findMovie(movieId: string | undefined) {
  if (!movieId) return null

  return (
    movies.value.find(
      (movie) =>
        movie.movie_id === movieId,
    ) ?? null
  )
}

function findTheater(theaterId: string | undefined) {
  if (!theaterId) return null

  return (
    theaters.value.find(
      (theater) =>
        theater.theater_id === theaterId,
    ) ?? null
  )
}

function findScreen(
  screenId: string | null | undefined,
) {
  if (!screenId) return null

  return (
    screens.value.find(
      (screen) =>
        screen.screen_id === screenId,
    ) ?? null
  )
}

// -----------------------------------------------------------------------
// Enriched bookings
// -----------------------------------------------------------------------

const enrichedBookings = computed<EnrichedBooking[]>(() => {
  return bookings.value.map((booking) => {
    const showtime = findShowtime(
      booking.showtime_id,
    )

    return {
      booking,
      showtime,
      movie: findMovie(showtime?.movie_id),
      theater: findTheater(showtime?.theater_id),
      screen: findScreen(showtime?.screen_id),
    }
  })
})

// -----------------------------------------------------------------------
// Filters
// -----------------------------------------------------------------------

const filteredBookings = computed(() => {
  if (activeFilter.value === 'all') {
    return enrichedBookings.value
  }

  return enrichedBookings.value.filter(
    (item) =>
      item.booking.status ===
      activeFilter.value,
  )
})

function setFilter(
  filter: 'all' | BookingStatus,
) {
  activeFilter.value = filter
}

// -----------------------------------------------------------------------
// Formatting
// -----------------------------------------------------------------------

function formatPrice(price: number) {
  return `฿${price.toFixed(2)}`
}

function formatShowDate(
  dateStr: string | undefined,
) {
  if (!dateStr) return '—'

  return new Date(dateStr).toLocaleDateString(
    'en-US',
    {
      weekday: 'short',
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    },
  )
}

function formatTicketDate(
  dateStr: string | undefined,
) {
  if (!dateStr) return '—'

  return new Date(dateStr).toLocaleDateString(
    'en-US',
    {
      weekday: 'short',
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    },
  )
}

// -----------------------------------------------------------------------
// Placeholder poster gradients
// -----------------------------------------------------------------------

const gradients = [
  'bg-gradient-to-b from-purple-950 via-purple-900 to-black',
  'bg-gradient-to-b from-emerald-950 via-emerald-900 to-black',
  'bg-gradient-to-b from-pink-950 via-pink-900 to-black',
  'bg-gradient-to-b from-blue-950 via-blue-900 to-black',
]

function gradientFor(bookingId: number) {
  return gradients[
    bookingId % gradients.length
  ]
}

// -----------------------------------------------------------------------
// E-Ticket
// -----------------------------------------------------------------------

const viewingTicket =
  ref<EnrichedBooking | null>(null)

function viewTicket(item: EnrichedBooking) {
  viewingTicket.value = item
}

function closeTicket() {
  viewingTicket.value = null
}

// -----------------------------------------------------------------------
// QR Code
// -----------------------------------------------------------------------

function getQrCodeUrl(
  bookingId: string,
) {
  const qrData = encodeURIComponent(
    `CINEMA-TICKET:${bookingId}`,
  )

  return `https://quickchart.io/qr?text=${qrData}&size=220&margin=2`
}

// -----------------------------------------------------------------------
// Payment
// -----------------------------------------------------------------------

const payingBooking =
  ref<EnrichedBooking | null>(null)

const paymentMethod =
  ref<PaymentMethod>('credit_card')

const payingLoading = ref(false)

function openPaymentModal(
  item: EnrichedBooking,
) {
  payingBooking.value = item
  paymentMethod.value = 'credit_card'
}

function closePaymentModal() {
  payingBooking.value = null
}

async function submitPayment() {
  if (!payingBooking.value) return

  payingLoading.value = true
  errorMessage.value = ''

  try {
    const paymentData = {
      booking_id:
        payingBooking.value.booking.booking_id,

      amount:
        payingBooking.value.booking.total_price,

      method: paymentMethod.value,

      status: 'paid',
    }

    console.log('========== PAY NOW ==========')
    console.log('PAYMENT DATA:', paymentData)

    const response = await api.post<Payment>(
      '/payments',
      paymentData,
    )

    console.log(
      'PAYMENT SUCCESS:',
      response.data,
    )

    await reloadBookings()

    closePaymentModal()
  } catch (err: any) {
    console.error(
      '========== PAYMENT ERROR ==========',
    )
    console.error(err)

    errorMessage.value =
      err?.response?.data?.error ??
      err?.message ??
      'Failed to complete payment. Please try again.'
  } finally {
    payingLoading.value = false
  }
}

// -----------------------------------------------------------------------
// Cancel booking
// -----------------------------------------------------------------------

async function cancelBooking(
  item: EnrichedBooking,
) {
  const confirmed = window.confirm(
    `Cancel this booking for "${
      item.movie?.title ??
      'this movie'
    }"?`,
  )

  if (!confirmed) return

  errorMessage.value = ''

  try {
    await api.put(
      `/bookings/${item.booking.booking_id}`,
      {
        status: 'cancelled',
      },
    )

    await reloadBookings()
  } catch (err: any) {
    console.error(
      '========== CANCEL BOOKING ERROR ==========',
    )
    console.error(err)

    errorMessage.value =
      err?.response?.data?.error ??
      err?.message ??
      'Failed to cancel booking. Please try again.'
  }
}

// -----------------------------------------------------------------------
// Prevent background scrolling while modal is open
// -----------------------------------------------------------------------

const isAnyModalOpen = computed(() => {
  return (
    payingBooking.value !== null ||
    viewingTicket.value !== null
  )
})

function handleEscape(event: KeyboardEvent) {
  if (event.key !== 'Escape') return

  if (payingBooking.value) {
    closePaymentModal()
  }

  if (viewingTicket.value) {
    closeTicket()
  }
}

onMounted(() => {
  window.addEventListener(
    'keydown',
    handleEscape,
  )
})
</script>

<template>
  <div
    class="min-h-screen bg-[#0a0a0f] px-4 py-8 text-white sm:px-6 sm:py-10"
  >
    <div class="mx-auto max-w-7xl">

      <!-- ========================================================= -->
      <!-- PAGE HEADER -->
      <!-- ========================================================= -->

      <div
        class="mb-7 flex items-center justify-between"
      >
        <div>
          <h1
            class="text-2xl font-bold sm:text-3xl"
          >
            My Bookings
          </h1>

          <p
            class="mt-1 text-sm text-gray-500"
          >
            Manage your movie bookings and tickets
          </p>
        </div>
      </div>

      <!-- ========================================================= -->
      <!-- ERROR -->
      <!-- ========================================================= -->

      <div
        v-if="errorMessage"
        class="mb-5 rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400"
      >
        {{ errorMessage }}
      </div>

      <!-- ========================================================= -->
      <!-- FILTERS -->
      <!-- ========================================================= -->

      <div class="mb-7">
        <div
          class="flex flex-wrap gap-2 sm:gap-3"
        >
          <button
            type="button"
            class="rounded-xl px-5 py-3 text-sm font-bold transition"
            :class="
              activeFilter === 'all'
                ? 'bg-pink-500 text-white shadow-lg shadow-pink-500/20'
                : 'border border-white/15 text-gray-300 hover:bg-white/5'
            "
            @click="setFilter('all')"
          >
            All
          </button>

          <button
            type="button"
            class="rounded-xl border px-5 py-3 text-sm font-bold transition"
            :class="
              activeFilter === 'confirmed'
                ? 'border-pink-500 bg-pink-500 text-white shadow-lg shadow-pink-500/20'
                : 'border-white/15 text-gray-300 hover:bg-white/5'
            "
            @click="setFilter('confirmed')"
          >
            Confirmed
          </button>

          <button
            type="button"
            class="rounded-xl border px-5 py-3 text-sm font-bold transition"
            :class="
              activeFilter === 'pending'
                ? 'border-pink-500 bg-pink-500 text-white shadow-lg shadow-pink-500/20'
                : 'border-white/15 text-gray-300 hover:bg-white/5'
            "
            @click="setFilter('pending')"
          >
            Pending
          </button>

          <button
            type="button"
            class="rounded-xl border px-5 py-3 text-sm font-bold transition"
            :class="
              activeFilter === 'cancelled'
                ? 'border-pink-500 bg-pink-500 text-white shadow-lg shadow-pink-500/20'
                : 'border-white/15 text-gray-300 hover:bg-white/5'
            "
            @click="setFilter('cancelled')"
          >
            Cancelled
          </button>
        </div>
      </div>

      <!-- ========================================================= -->
      <!-- LOADING -->
      <!-- ========================================================= -->

      <div
        v-if="loading"
        class="py-20 text-center"
      >
        <div
          class="mx-auto mb-4 h-8 w-8 animate-spin rounded-full border-2 border-white/10 border-t-pink-500"
        ></div>

        <p
          class="text-sm text-gray-500"
        >
          Loading your bookings…
        </p>
      </div>

      <!-- ========================================================= -->
      <!-- BOOKING LIST -->
      <!-- ========================================================= -->

      <div
        v-else
        class="space-y-4"
      >
        <div
          v-for="item in filteredBookings"
          :key="item.booking.booking_id"
          class="flex flex-col gap-5 rounded-2xl border border-white/10 bg-[#15151d] p-5 transition hover:border-white/20 md:flex-row md:items-center"
          :class="{
            'opacity-50':
              item.booking.status === 'cancelled',
          }"
        >

          <!-- Poster -->

          <div
            class="h-28 w-20 shrink-0 overflow-hidden rounded-xl bg-cover bg-center shadow-lg"
            :class="
              !item.movie?.poster_url
                ? gradientFor(item.booking.id)
                : ''
            "
            :style="
              item.movie?.poster_url
                ? {
                    backgroundImage: `url(${item.movie.poster_url})`,
                  }
                : {}
            "
          ></div>

          <!-- Booking info -->

          <div class="min-w-0 flex-1">

            <h2
              class="truncate text-lg font-bold"
            >
              {{
                item.movie?.title ??
                'Unknown movie'
              }}
            </h2>

            <p
              class="mt-2 text-sm text-gray-500"
            >
              {{
                item.theater?.name ??
                'Unknown theater'
              }}

              <template
                v-if="item.screen"
              >
                · {{ item.screen.name }}
              </template>
            </p>

            <p
              class="mt-1 text-sm text-gray-500"
            >
              {{
                formatShowDate(
                  item.showtime?.show_date,
                )
              }}

              ·

              {{
                item.showtime?.start_time ??
                '—'
              }}

              · Seats
              {{
                item.booking.seat_numbers.join(
                  ', ',
                )
              }}
            </p>
          </div>

          <!-- Status -->

          <div class="shrink-0">
            <span
              class="rounded-full px-4 py-2 text-xs font-bold uppercase tracking-wide"
              :class="{
                'bg-emerald-500/10 text-emerald-400':
                  item.booking.status ===
                  'confirmed',

                'bg-amber-500/10 text-amber-400':
                  item.booking.status ===
                  'pending',

                'bg-red-500/10 text-red-400':
                  item.booking.status ===
                  'cancelled',
              }"
            >
              {{ item.booking.status }}
            </span>
          </div>

          <!-- Price + actions -->

          <div
            class="flex shrink-0 flex-col items-start gap-2 sm:items-end md:w-52"
          >
            <p
              class="text-lg font-black"
            >
              {{
                formatPrice(
                  item.booking.total_price,
                )
              }}
            </p>

            <div
              class="flex flex-wrap items-center gap-2"
            >

              <!-- CONFIRMED -->

              <template
                v-if="
                  item.booking.status ===
                  'confirmed'
                "
              >
                <button
                  type="button"
                  class="rounded-xl bg-pink-500 px-5 py-2.5 text-sm font-bold text-white transition hover:bg-pink-600 hover:shadow-lg hover:shadow-pink-500/20"
                  @click="viewTicket(item)"
                >
                  View ticket
                </button>

                <button
                  type="button"
                  class="rounded-xl border border-red-500/20 px-3 py-2.5 text-xs font-semibold text-red-400 transition hover:bg-red-500/10"
                  @click="
                    cancelBooking(item)
                  "
                >
                  Cancel
                </button>
              </template>

              <!-- PENDING -->

              <template
                v-else-if="
                  item.booking.status ===
                  'pending'
                "
              >
                <button
                  type="button"
                  class="rounded-xl bg-pink-500 px-5 py-2.5 text-sm font-bold text-white transition hover:bg-pink-600 hover:shadow-lg hover:shadow-pink-500/20"
                  @click="
                    openPaymentModal(item)
                  "
                >
                  Complete payment
                </button>

                <button
                  type="button"
                  class="rounded-xl border border-red-500/20 px-3 py-2.5 text-xs font-semibold text-red-400 transition hover:bg-red-500/10"
                  @click="
                    cancelBooking(item)
                  "
                >
                  Cancel
                </button>
              </template>

              <!-- CANCELLED -->

              <button
                v-else
                type="button"
                class="rounded-xl border border-white/10 px-5 py-2.5 text-sm font-semibold text-gray-600"
                disabled
              >
                Cancelled
              </button>

            </div>
          </div>
        </div>

        <!-- Empty -->

        <div
          v-if="
            filteredBookings.length === 0
          "
          class="rounded-2xl border border-white/10 bg-[#15151d] py-20 text-center"
        >
          <div
            class="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-white/5 text-2xl"
          >
            🎟️
          </div>

          <p
            class="text-sm text-gray-500"
          >
            No bookings found.
          </p>
        </div>
      </div>
    </div>

    <!-- =========================================================== -->
    <!-- PAYMENT MODAL -->
    <!-- =========================================================== -->

    <Transition name="modal">
      <div
        v-if="payingBooking"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 px-4 backdrop-blur-sm"
        @click.self="
          closePaymentModal()
        "
      >
        <div
          class="w-full max-w-md rounded-2xl border border-white/10 bg-[#15151d] p-6 shadow-2xl"
        >

          <div
            class="mb-5 flex items-center justify-between"
          >
            <div>
              <h3
                class="text-lg font-bold"
              >
                Complete payment
              </h3>

              <p
                class="mt-1 text-xs text-gray-500"
              >
                Secure your booking
              </p>
            </div>

            <button
              type="button"
              class="flex h-9 w-9 items-center justify-center rounded-full text-gray-400 transition hover:bg-white/5 hover:text-white"
              @click="
                closePaymentModal()
              "
            >
              ✕
            </button>
          </div>

          <div
            class="mb-5 rounded-xl border border-white/10 bg-black/20 p-4"
          >
            <p
              class="font-bold"
            >
              {{
                payingBooking.movie?.title ??
                'Movie'
              }}
            </p>

            <p
              class="mt-1 text-sm text-gray-500"
            >
              {{
                payingBooking.theater?.name ??
                'Unknown theater'
              }}

              <template
                v-if="
                  payingBooking.screen
                "
              >
                ·
                {{
                  payingBooking.screen.name
                }}
              </template>
            </p>

            <div
              class="mt-3 flex items-center justify-between border-t border-white/10 pt-3"
            >
              <span
                class="text-sm text-gray-500"
              >
                Total
              </span>

              <span
                class="text-lg font-black text-pink-400"
              >
                {{
                  formatPrice(
                    payingBooking.booking
                      .total_price,
                  )
                }}
              </span>
            </div>
          </div>

          <label
            class="mb-2 block text-xs font-medium text-gray-400"
          >
            Payment method
          </label>

          <select
            v-model="paymentMethod"
            class="mb-6 h-11 w-full rounded-lg border border-white/15 bg-[#0a0a0f] px-3 text-sm text-white outline-none focus:border-pink-500"
          >
            <option value="credit_card">
              Credit card
            </option>

            <option value="promptpay">
              PromptPay
            </option>

            <option value="transfer">
              Bank transfer
            </option>

            <option value="cash">
              Cash
            </option>
          </select>

          <div
            class="flex justify-end gap-3"
          >
            <button
              type="button"
              class="rounded-xl border border-white/15 px-5 py-2.5 text-sm font-semibold text-gray-300 transition hover:bg-white/5"
              @click="
                closePaymentModal()
              "
            >
              Cancel
            </button>

            <button
              type="button"
              :disabled="payingLoading"
              class="rounded-xl bg-pink-500 px-6 py-3 font-semibold text-white transition hover:bg-pink-600 disabled:cursor-not-allowed disabled:opacity-50"
              @click="submitPayment"
            >
              {{
                payingLoading
                  ? 'Processing...'
                  : 'Pay now'
              }}
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- =========================================================== -->
    <!-- E-TICKET MODAL (landscape layout: poster+details on the left,
         perforated stub with booking ID / QR / total / close on the
         right — same content as before, just reflowed from a
         top-to-bottom stack into two side-by-side panels. Falls back
         to the original stacked layout below the `sm` breakpoint.) -->
    <!-- =========================================================== -->

    <Transition name="modal">
      <div
        v-if="viewingTicket"
        class="fixed inset-0 z-50 flex items-center justify-center overflow-y-auto bg-black/85 px-4 py-8 backdrop-blur-sm"
        @click.self="closeTicket"
      >

        <!-- Ticket -->

        <div
          class="relative flex w-full max-w-3xl flex-col overflow-hidden rounded-[28px] bg-[#f7f7f8] text-black shadow-2xl sm:flex-row"
        >

          <!-- ===================================================== -->
          <!-- Left: poster + core details -->
          <!-- ===================================================== -->

          <div class="flex flex-col sm:w-3/5">

            <!-- Poster header -->

            <div
              class="relative h-56 shrink-0 overflow-hidden bg-black sm:h-64"
            >

              <!-- Poster -->

              <div
                v-if="
                  viewingTicket.movie
                    ?.poster_url
                "
                class="absolute inset-0 bg-cover bg-center"
                :style="{
                  backgroundImage: `url(${viewingTicket.movie.poster_url})`,
                }"
              ></div>

              <div
                v-else
                class="absolute inset-0"
                :class="
                  gradientFor(
                    viewingTicket.booking.id,
                  )
                "
              ></div>

              <!-- Dark overlay -->

              <div
                class="absolute inset-0 bg-gradient-to-t from-black via-black/40 to-transparent"
              ></div>

              <!-- Close -->

              <button
                type="button"
                class="absolute right-4 top-4 z-10 flex h-9 w-9 items-center justify-center rounded-full bg-black/50 text-white backdrop-blur transition hover:bg-black/70"
                @click="closeTicket"
              >
                ✕
              </button>

              <!-- Ticket title -->

              <div
                class="absolute bottom-5 left-5 right-5"
              >
                <div
                  class="mb-2 inline-flex items-center rounded-full bg-emerald-400 px-3 py-1 text-[10px] font-black uppercase tracking-wider text-black"
                >
                  ✓ Confirmed
                </div>

                <h2
                  class="text-2xl font-black text-white drop-shadow-lg"
                >
                  {{
                    viewingTicket.movie
                      ?.title ??
                    'Movie Ticket'
                  }}
                </h2>
              </div>
            </div>

            <!-- Core details -->

            <div class="flex-1 px-6 pb-6 pt-6">

              <!-- Cinema -->

              <div class="mb-6">
                <p
                  class="text-[10px] font-bold uppercase tracking-[0.2em] text-gray-400"
                >
                  Cinema
                </p>

                <p
                  class="mt-1 text-lg font-black"
                >
                  {{
                    viewingTicket.theater
                      ?.name ??
                    'Unknown theater'
                  }}
                </p>
              </div>

              <!-- Date / Time -->

              <div
                class="mb-6 grid grid-cols-2 gap-5"
              >
                <div>
                  <p
                    class="text-[10px] font-bold uppercase tracking-[0.2em] text-gray-400"
                  >
                    Date
                  </p>

                  <p
                    class="mt-1 text-sm font-bold"
                  >
                    {{
                      formatTicketDate(
                        viewingTicket
                          .showtime
                          ?.show_date,
                      )
                    }}
                  </p>
                </div>

                <div>
                  <p
                    class="text-[10px] font-bold uppercase tracking-[0.2em] text-gray-400"
                  >
                    Time
                  </p>

                  <p
                    class="mt-1 text-sm font-bold"
                  >
                    {{
                      viewingTicket.showtime
                        ?.start_time ??
                      '—'
                    }}

                    <span
                      v-if="
                        viewingTicket
                          .showtime
                          ?.end_time
                      "
                      class="font-normal text-gray-400"
                    >
                      -
                      {{
                        viewingTicket
                          .showtime
                          .end_time
                      }}
                    </span>
                  </p>
                </div>
              </div>

              <!-- Screen / Seats -->

              <div
                class="grid grid-cols-2 gap-5"
              >
                <div>
                  <p
                    class="text-[10px] font-bold uppercase tracking-[0.2em] text-gray-400"
                  >
                    Screen
                  </p>

                  <p
                    class="mt-1 text-sm font-bold"
                  >
                    {{
                      viewingTicket.screen
                        ?.name ??
                      '—'
                    }}
                  </p>
                </div>

                <div>
                  <p
                    class="text-[10px] font-bold uppercase tracking-[0.2em] text-gray-400"
                  >
                    Seats
                  </p>

                  <p
                    class="mt-1 text-sm font-black text-pink-500"
                  >
                    {{
                      viewingTicket.booking
                        .seat_numbers.join(
                          ', ',
                        )
                    }}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- ===================================================== -->
          <!-- Perforated Divider + Right: Booking ID / QR / Price stub -->
          <!-- ===================================================== -->

          <div
            class="relative flex flex-col items-center justify-center gap-5 border-t border-dashed border-gray-300 px-6 py-8 sm:w-2/5 sm:border-l sm:border-t-0"
          >

            <!-- Notch circles: horizontal divider (mobile, stacked) -->

            <div
              class="absolute -left-3 top-0 h-6 w-6 -translate-y-1/2 rounded-full bg-[#0a0a0f] sm:hidden"
            ></div>

            <div
              class="absolute -right-3 top-0 h-6 w-6 -translate-y-1/2 rounded-full bg-[#0a0a0f] sm:hidden"
            ></div>

            <!-- Notch circles: vertical divider (sm+, side-by-side) -->

            <div
              class="absolute left-1/2 top-0 hidden h-6 w-6 -translate-x-1/2 -translate-y-1/2 rounded-full bg-[#0a0a0f] sm:block"
            ></div>

            <div
              class="absolute bottom-0 left-1/2 hidden h-6 w-6 -translate-x-1/2 translate-y-1/2 rounded-full bg-[#0a0a0f] sm:block"
            ></div>

            <!-- Booking ID -->

            <div class="w-full">
              <p
                class="text-center text-[10px] font-bold uppercase tracking-[0.2em] text-gray-400"
              >
                Booking ID
              </p>

              <p
                class="mt-2 break-all text-center font-mono text-xs font-bold text-gray-700"
              >
                {{
                  viewingTicket.booking
                    .booking_id
                }}
              </p>
            </div>

            <!-- QR Code -->

            <div
              class="flex flex-col items-center justify-center"
            >
              <div
                class="rounded-2xl border border-gray-200 bg-white p-3 shadow-sm"
              >
                <img
                  :src="
                    getQrCodeUrl(
                      viewingTicket.booking
                        .booking_id,
                    )
                  "
                  alt="Ticket QR Code"
                  class="h-36 w-36"
                />
              </div>

              <p
                class="mt-3 text-[10px] font-medium uppercase tracking-widest text-gray-400"
              >
                Scan to verify ticket
              </p>
            </div>

            <!-- Price -->

            <div
              class="flex w-full items-center justify-between border-t border-gray-200 pt-4"
            >
              <span
                class="text-xs font-bold uppercase tracking-wider text-gray-400"
              >
                Total paid
              </span>

              <span
                class="text-xl font-black"
              >
                {{
                  formatPrice(
                    viewingTicket.booking
                      .total_price,
                  )
                }}
              </span>
            </div>

            <!-- Close button -->

            <button
              type="button"
              class="w-full rounded-xl bg-black px-5 py-3 text-sm font-bold text-white transition hover:bg-gray-800"
              @click="closeTicket"
            >
              Close ticket
            </button>

            <p
              class="text-center text-[10px] text-gray-400"
            >
              Please present this ticket at the cinema.
            </p>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition:
    opacity 0.2s ease,
    transform 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from > div,
.modal-leave-to > div {
  transform: scale(0.96);
}
</style>