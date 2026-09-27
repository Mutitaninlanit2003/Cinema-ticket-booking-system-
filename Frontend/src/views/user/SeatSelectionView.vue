<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()

const showtimeId = computed(() => String(route.params.showtimeId))

// -----------------------------------------------------------------------
// API client
// NOTE: replace with your project's central API client if you have one.
// -----------------------------------------------------------------------
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
})

// -----------------------------------------------------------------------
// Types — match backend models exactly
// (models/seat.go, models/showtime.go, models/movie.go, models/theater.go)
// -----------------------------------------------------------------------
type SeatType = 'standard' | 'premium' | 'vip'

// Seat.Status is a free-form string on the backend (default "AVAILABLE"),
// but ManageSeats.vue only ever writes "AVAILABLE" / "DISABLED", and a
// seat flips to "BOOKED" once a Booking references it. There is no
// "locked/held" state in the models provided — no reservation table was
// given — so a seat someone else has mid-checkout will still read as
// AVAILABLE here. See the note on `toggleSeat` below.
type BackendSeatStatus = 'AVAILABLE' | 'BOOKED' | 'DISABLED'

interface Seat {
  id: number
  seat_id: string
  screen_id: string
  seat_number: string
  seat_type: SeatType
  status: BackendSeatStatus
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

const showtimes = ref<Showtime[]>([])
const seats = ref<Seat[]>([])
const movies = ref<Movie[]>([])
const theaters = ref<Theater[]>([])

const loading = ref(false)
const errorMessage = ref('')

// -----------------------------------------------------------------------
// Data loading — same "fetch the full lists, join client-side" pattern
// used in ManageSeats.vue, since there's no single "showtime details"
// endpoint among the models provided.
// -----------------------------------------------------------------------
async function loadAll() {
  loading.value = true
  errorMessage.value = ''

  try {
    const [showtimesRes, seatsRes, moviesRes, theatersRes] = await Promise.all([
      api.get<Showtime[]>('/showtimes'),
      api.get<Seat[]>('/seats'),
      api.get<Movie[]>('/movies'),
      api.get<Theater[]>('/theaters'),
    ])

    showtimes.value = showtimesRes.data
    seats.value = seatsRes.data
    movies.value = moviesRes.data
    theaters.value = theatersRes.data
  } catch (err) {
    errorMessage.value = 'Failed to load seats for this showtime. Please try again.'
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)

// -----------------------------------------------------------------------
// Showtime + its screen's seats
// -----------------------------------------------------------------------
const showtime = computed(
  () => showtimes.value.find((s) => s.showtime_id === showtimeId.value) ?? null,
)

const movie = computed(
  () => movies.value.find((m) => m.movie_id === showtime.value?.movie_id) ?? null,
)

const theater = computed(
  () => theaters.value.find((t) => t.theater_id === showtime.value?.theater_id) ?? null,
)

const screenSeats = computed(() => {
  if (!showtime.value?.screen_id) return []
  return seats.value.filter((seat) => seat.screen_id === showtime.value!.screen_id)
})

function parseSeatNumber(seatNumber: string) {
  const match = seatNumber.match(/^([A-Za-z]+)(\d+)$/)

  if (!match) {
    return { row: seatNumber, number: 0 }
  }

  return { row: match[1], number: Number(match[2]) }
}

const rows = computed(() => {
  const rowSet = new Set(
    screenSeats.value.map((seat) => parseSeatNumber(seat.seat_number).row),
  )

  // Highest letter nearest the screen, descending toward the back —
  // matches the reference layout (arrangement only, no other change).
  return [...rowSet].sort().reverse()
})

function seatsInRow(row: string) {
  return screenSeats.value
    .filter((seat) => parseSeatNumber(seat.seat_number).row === row)
    .sort(
      (a, b) =>
        parseSeatNumber(a.seat_number).number - parseSeatNumber(b.seat_number).number,
    )
}

// Splits a row's seats into a left and right half so the template can
// draw a center aisle gap between them — purely a layout/arrangement
// change, doesn't touch seat colors, status, or any other behavior.
function seatGroups(row: string) {
  const rowSeats = seatsInRow(row)
  const mid = Math.ceil(rowSeats.length / 2)

  return [rowSeats.slice(0, mid), rowSeats.slice(mid)]
}

// Grid-template-columns for one row: a fixed-width label column, one
// flexible (minmax(0, 1fr)) column per seat in the left half, a fixed
// aisle column, one flexible column per seat in the right half, and a
// fixed label column. Because every seat column is `minmax(0, 1fr)`
// rather than a fixed pixel width, the whole row — and every seat
// button in it — shrinks to fit the container's actual width instead
// of overflowing and requiring a horizontal scrollbar.
function rowGridColumns(row: string) {
  const [left, right] = seatGroups(row)

  return `auto repeat(${left.length}, minmax(0, 1fr)) auto repeat(${right.length}, minmax(0, 1fr)) auto`
}

// -----------------------------------------------------------------------
// Per-zone pricing — the price for a given seat depends on its
// seat_type AND the current showtime's price_standard/price_premium/
// price_vip (not a single flat price anymore).
// -----------------------------------------------------------------------
function seatPrice(seat: Seat) {
  if (!showtime.value) return 0

  if (seat.seat_type === 'vip') return showtime.value.price_vip
  if (seat.seat_type === 'premium') return showtime.value.price_premium

  return showtime.value.price_standard
}

// -----------------------------------------------------------------------
// Selection — kept as local UI state (a set of seat_id), separate from
// each seat's backend `status`. Only AVAILABLE seats are selectable.
// -----------------------------------------------------------------------
const selectedSeatIds = ref<Set<string>>(new Set())

const selectedSeats = computed(() =>
  screenSeats.value.filter((seat) => selectedSeatIds.value.has(seat.seat_id)),
)

const totalPrice = computed(() =>
  selectedSeats.value.reduce((sum, seat) => sum + seatPrice(seat), 0),
)

// Groups the current selection by seat_type, for the "N × price" style
// breakdown in the Booking Summary — since seats can now be a mix of
// different zones/prices, a single "Price per seat" line no longer
// makes sense.
const SEAT_TYPE_LABELS: Record<SeatType, string> = {
  standard: 'Standard',
  premium: 'Premium',
  vip: 'VIP',
}

const selectionBreakdown = computed(() => {
  const groups = new Map<SeatType, { count: number; unitPrice: number }>()

  for (const seat of selectedSeats.value) {
    const existing = groups.get(seat.seat_type)

    if (existing) {
      existing.count += 1
    } else {
      groups.set(seat.seat_type, { count: 1, unitPrice: seatPrice(seat) })
    }
  }

  return (['vip', 'premium', 'standard'] as SeatType[])
    .filter((type) => groups.has(type))
    .map((type) => {
      const group = groups.get(type)!

      return {
        type,
        label: SEAT_TYPE_LABELS[type],
        count: group.count,
        unitPrice: group.unitPrice,
        subtotal: group.count * group.unitPrice,
      }
    })
})

function toggleSeat(seat: Seat) {
  if (seat.status !== 'AVAILABLE') return

  const next = new Set(selectedSeatIds.value)

  if (next.has(seat.seat_id)) {
    next.delete(seat.seat_id)
  } else {
    next.add(seat.seat_id)
  }

  selectedSeatIds.value = next
}

function seatClass(seat: Seat) {
  if (selectedSeatIds.value.has(seat.seat_id)) {
    return 'border-[#ff3158] bg-[#ff3158] text-white shadow-lg shadow-[#ff3158]/30'
  }

  if (seat.status === 'BOOKED') {
    return 'cursor-not-allowed border-white/5 bg-white/5 text-gray-700'
  }

  if (seat.status === 'DISABLED') {
    return 'cursor-not-allowed border-violet-500/20 bg-violet-500/10 text-violet-500'
  }

  // Available — colored by seat_type (zone) rather than one generic color.
  if (seat.seat_type === 'vip') {
    return 'border-blue-500/30 bg-blue-500/10 text-blue-300 hover:border-blue-400/60 hover:bg-blue-500/20'
  }

  if (seat.seat_type === 'premium') {
    return 'border-purple-500/30 bg-purple-500/10 text-purple-300 hover:border-purple-400/60 hover:bg-purple-500/20'
  }

  // standard
  return 'border-emerald-500/30 bg-emerald-500/10 text-emerald-300 hover:border-emerald-400/60 hover:bg-emerald-500/20'
}

function formatPrice(price: number) {
  return `฿${price.toFixed(0)}`
}

// -----------------------------------------------------------------------
// Continue — hands off seat_number strings and showtime_id, matching
// Booking.SeatNumbers / Booking.ShowtimeID so the summary page can
// POST /bookings directly with this data.
// -----------------------------------------------------------------------
function continueToSummary() {
  if (selectedSeats.value.length === 0) return

  router.push({
    path: '/booking/summary',
    query: {
      showtimeId: showtimeId.value,
      seats: selectedSeats.value.map((seat) => seat.seat_number).join(','),
      totalPrice: String(totalPrice.value),
    },
  })
}
</script>

<template>
  <div class="min-h-screen bg-[#0d0d12] px-6 py-10 text-white">
    <div class="mx-auto max-w-7xl">

      <!-- Header -->
      <div class="mb-8">
        <p class="text-xs font-medium uppercase tracking-[0.25em] text-gray-500">
          {{ movie?.title ?? 'Showtime' }}
          <template v-if="theater"> · {{ theater.name }}</template>
        </p>

        <h1 class="mt-2 text-3xl font-black">
          Select Your Seats
        </h1>

        <p class="mt-2 text-sm text-gray-500">
          <template v-if="showtime">
            {{ showtime.show_date }} · {{ showtime.start_time }} –
          </template>
          Choose your preferred seats before continuing.
        </p>
      </div>

      <!-- Error banner -->
      <div
        v-if="errorMessage"
        class="mb-8 rounded-lg border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400"
      >
        {{ errorMessage }}
      </div>

      <!-- Step -->
      <div class="mb-10 flex items-center gap-4 text-sm">
        <div class="flex items-center gap-2 text-gray-500">
          <span
            class="flex h-6 w-6 items-center justify-center rounded-full bg-white/10 text-xs"
          >
            1
          </span>
          Showtime
        </div>

        <div class="h-px w-10 bg-white/20"></div>

        <div class="flex items-center gap-2 font-semibold text-white">
          <span
            class="flex h-6 w-6 items-center justify-center rounded-full bg-[#ff3158] text-xs font-bold"
          >
            2
          </span>
          Seats
        </div>

        <div class="h-px w-10 bg-white/20"></div>

        <div class="flex items-center gap-2 text-gray-500">
          <span
            class="flex h-6 w-6 items-center justify-center rounded-full border border-white/20 text-xs"
          >
            3
          </span>
          Summary
        </div>
      </div>

      <p v-if="loading" class="py-12 text-center text-sm text-gray-500">
        Loading seats…
      </p>

      <p
        v-else-if="!showtime"
        class="py-12 text-center text-sm text-gray-500"
      >
        This showtime couldn't be found.
      </p>

      <p
        v-else-if="!showtime.screen_id"
        class="py-12 text-center text-sm text-gray-500"
      >
        This showtime has no screen assigned yet, so no seat map is available.
      </p>

      <div v-else class="grid gap-8 lg:grid-cols-[1fr_320px]">

        <!-- Seat Area -->
        <!-- min-w-0 is required here: grid items default to min-width:auto,
             which means this section refuses to shrink below the intrinsic
             width of its content. Without it, this section (and the whole
             1fr track) would be pushed past the grid's available space,
             shoving the 320px Booking Summary column off screen instead of
             letting it sit fixed on the right. -->
        <section
          class="min-w-0 rounded-2xl border border-white/10 bg-[#121218] p-6 sm:p-10"
        >

          <!-- Screen -->
          <div class="mb-12">
            <div
              class="mx-auto h-2 max-w-xl rounded-full bg-gradient-to-r from-transparent via-white/70 to-transparent"
            ></div>

            <p
              class="mt-3 text-center text-xs font-medium uppercase tracking-[0.4em] text-gray-600"
            >
              Screen
            </p>
          </div>

          <!-- Seats -->
          <!-- Each row is a CSS grid whose seat columns are minmax(0, 1fr)
               (see rowGridColumns), so seat buttons shrink to fit the
               container's width instead of overflowing horizontally —
               no scrollbar needed regardless of how many seats a row has. -->
          <div class="w-full space-y-1.5 sm:space-y-2">

            <div
              v-for="row in rows"
              :key="row"
              class="grid items-center gap-1 sm:gap-1.5"
              :style="{ gridTemplateColumns: rowGridColumns(row) }"
            >
              <span class="w-4 text-center text-[10px] font-bold text-gray-600 sm:w-5 sm:text-xs">
                {{ row }}
              </span>

              <button
                v-for="seat in seatGroups(row)[0]"
                :key="seat.seat_id"
                type="button"
                class="aspect-square w-full min-w-0 rounded-md border text-[9px] font-bold transition sm:rounded-lg sm:text-xs"
                :class="seatClass(seat)"
                :disabled="seat.status !== 'AVAILABLE'"
                @click="toggleSeat(seat)"
              >
                {{ seat.seat_number }}
              </button>

              <!-- Aisle gap between the two halves of the row -->
              <div></div>

              <button
                v-for="seat in seatGroups(row)[1]"
                :key="seat.seat_id"
                type="button"
                class="aspect-square w-full min-w-0 rounded-md border text-[9px] font-bold transition sm:rounded-lg sm:text-xs"
                :class="seatClass(seat)"
                :disabled="seat.status !== 'AVAILABLE'"
                @click="toggleSeat(seat)"
              >
                {{ seat.seat_number }}
              </button>

              <span class="w-4 text-center text-[10px] font-bold text-gray-600 sm:w-5 sm:text-xs">
                {{ row }}
              </span>
            </div>

            <p v-if="rows.length === 0" class="py-8 text-center text-sm text-gray-600">
              No seats have been configured for this screen yet.
            </p>

          </div>

          <!-- Legend -->
          <!-- Zone chips now show that zone's actual price for this
               showtime, since VIP/Premium/Standard no longer share one
               flat price. -->
          <div
            class="mt-12 flex flex-wrap justify-center gap-5 text-xs text-gray-500"
          >
            <div class="flex items-center gap-2">
              <span class="h-3 w-3 rounded border border-emerald-500/30 bg-emerald-500/10"></span>
              Standard
              <span v-if="showtime" class="text-gray-600">
                ({{ formatPrice(showtime.price_standard) }})
              </span>
            </div>

            <div class="flex items-center gap-2">
              <span class="h-3 w-3 rounded border border-purple-500/30 bg-purple-500/10"></span>
              Premium
              <span v-if="showtime" class="text-gray-600">
                ({{ formatPrice(showtime.price_premium) }})
              </span>
            </div>

            <div class="flex items-center gap-2">
              <span class="h-3 w-3 rounded border border-blue-500/30 bg-blue-500/10"></span>
              VIP
              <span v-if="showtime" class="text-gray-600">
                ({{ formatPrice(showtime.price_vip) }})
              </span>
            </div>

            <div class="flex items-center gap-2">
              <span class="h-3 w-3 rounded bg-[#ff3158]"></span>
              Selected
            </div>

            <div class="flex items-center gap-2">
              <span class="h-3 w-3 rounded bg-white/5"></span>
              Booked
            </div>

            <div class="flex items-center gap-2">
              <span class="h-3 w-3 rounded bg-violet-500/20"></span>
              Disabled
            </div>
          </div>
        </section>

        <!-- Booking Summary -->
        <aside
          class="h-fit rounded-2xl border border-white/10 bg-[#121218] p-6"
        >
          <p
            class="text-xs font-medium uppercase tracking-[0.2em] text-gray-500"
          >
            Your Selection
          </p>

          <h2 class="mt-2 text-xl font-bold">
            Booking Summary
          </h2>

          <div class="my-6 h-px bg-white/10"></div>

          <div>
            <p class="text-sm text-gray-500">
              Selected seats
            </p>

            <p
              v-if="selectedSeats.length"
              class="mt-2 font-semibold text-white"
            >
              {{ selectedSeats.map((seat) => seat.seat_number).join(', ') }}
            </p>

            <p
              v-else
              class="mt-2 text-sm text-gray-600"
            >
              No seats selected
            </p>
          </div>

          <div class="my-6 h-px bg-white/10"></div>

          <!-- Per-zone breakdown — replaces the old single "Seats" +
               "Price per seat" rows, since a selection can now mix
               different-priced zones (e.g. 2 VIP + 1 Standard). -->
          <div v-if="selectionBreakdown.length" class="space-y-2">
            <div
              v-for="group in selectionBreakdown"
              :key="group.type"
              class="flex items-center justify-between text-sm"
            >
              <span class="text-gray-400">
                {{ group.label }} × {{ group.count }}
                <span class="text-gray-600">
                  ({{ formatPrice(group.unitPrice) }} each)
                </span>
              </span>

              <span class="font-semibold text-white">
                {{ formatPrice(group.subtotal) }}
              </span>
            </div>
          </div>

          <p v-else class="text-sm text-gray-600">
            Select seats to see the price breakdown.
          </p>

          <div class="my-6 h-px bg-white/10"></div>

          <div class="flex items-center justify-between">
            <span class="text-gray-400">
              Total
            </span>

            <span class="text-2xl font-black">
              {{ formatPrice(totalPrice) }}
            </span>
          </div>

          <button
            type="button"
            class="mt-6 w-full rounded-lg bg-[#ff3158] py-3 text-sm font-bold text-white transition hover:bg-[#ff4568] disabled:cursor-not-allowed disabled:opacity-40"
            :disabled="selectedSeats.length === 0"
            @click="continueToSummary"
          >
            Continue to summary
          </button>
        </aside>

      </div>
    </div>
  </div>
</template>