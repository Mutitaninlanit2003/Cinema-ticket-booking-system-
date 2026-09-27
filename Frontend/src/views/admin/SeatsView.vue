<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import axios from 'axios'
import { Armchair } from 'lucide-vue-next'

// -----------------------------------------------------------------------
// API client
// -----------------------------------------------------------------------

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
})

// -----------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------

type SeatType = 'standard' | 'premium' | 'vip'

interface Seat {
  id: number
  seat_id: string
  screen_id: string
  seat_number: string
  seat_type: SeatType
  status: string
  created_at: string
  updated_at: string
}

interface Screen {
  id: number
  screen_id: string
  theater_id: string
  name: string
  total_seats: number
}

interface Theater {
  id: number
  theater_id: string
  name: string
}

// -----------------------------------------------------------------------
// State
// -----------------------------------------------------------------------

const showForm = ref(false)
const editingSeat = ref<Seat | null>(null)

const selectedTheaterId = ref('')
const selectedScreenId = ref('')

const theaters = ref<Theater[]>([])
const screens = ref<Screen[]>([])
const seats = ref<Seat[]>([])

const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

// -----------------------------------------------------------------------
// Single seat form
// -----------------------------------------------------------------------

const form = ref({
  row: '',
  number: 1,
  seat_type: 'standard' as SeatType,
  status: 'AVAILABLE',
})

// -----------------------------------------------------------------------
// Multi-seat edit
//
// เมื่อเลือก Seat Type แล้ว สามารถคลิกที่นั่งหลายตัวบน Seat Map
// แล้วกด Save changes ครั้งเดียวได้
// -----------------------------------------------------------------------

const selectedSeatIds = ref<string[]>([])
const batchSeatType = ref<SeatType | null>(null)

// -----------------------------------------------------------------------
// Bulk generate seats
// -----------------------------------------------------------------------

const showGenerateForm = ref(false)
const generating = ref(false)

const generateForm = ref({
  start_row: '',
  end_row: '',
  seats_per_row: 10,
  seat_type: 'standard' as SeatType,
  status: 'AVAILABLE',
})

function openGenerateForm() {
  if (!selectedScreenId.value) return

  generateForm.value = {
    start_row: '',
    end_row: '',
    seats_per_row: 10,
    seat_type: 'standard',
    status: 'AVAILABLE',
  }

  showGenerateForm.value = true
}

function closeGenerateForm() {
  showGenerateForm.value = false
}

// -----------------------------------------------------------------------
// Generate seats
// -----------------------------------------------------------------------

async function generateSeats() {
  if (
    !selectedScreenId.value ||
    !generateForm.value.start_row ||
    !generateForm.value.end_row ||
    !generateForm.value.seats_per_row
  ) {
    errorMessage.value = 'Please fill in all generate seat fields.'
    return
  }

  errorMessage.value = ''
  successMessage.value = ''
  generating.value = true

  try {
    let startRow = generateForm.value.start_row
      .trim()
      .toUpperCase()

    let endRow = generateForm.value.end_row
      .trim()
      .toUpperCase()

    if (startRow.length !== 1 || endRow.length !== 1) {
      errorMessage.value =
        'Row must be a single letter from A to Z.'
      return
    }

    if (
      !/^[A-Z]$/.test(startRow) ||
      !/^[A-Z]$/.test(endRow)
    ) {
      errorMessage.value =
        'Row must be between A and Z.'
      return
    }

    // Swap reversed rows automatically
    if (
      startRow.charCodeAt(0) >
      endRow.charCodeAt(0)
    ) {
      const temp = startRow
      startRow = endRow
      endRow = temp
    }

    const { data } = await api.post('/seats/generate', {
      screen_id: selectedScreenId.value,
      start_row: startRow,
      end_row: endRow,
      seats_per_row: Number(
        generateForm.value.seats_per_row
      ),
      seat_type: generateForm.value.seat_type,
      status: generateForm.value.status,
    })

    await reloadSeats()

    closeGenerateForm()

    successMessage.value =
      `${data.count} seats generated successfully.`
  } catch (err: any) {
    errorMessage.value =
      err.response?.data?.error ||
      'Failed to generate seats. Please check the form and try again.'

    console.error(err)
  } finally {
    generating.value = false
  }
}

// -----------------------------------------------------------------------
// Data loading
// -----------------------------------------------------------------------

async function loadAll() {
  loading.value = true
  errorMessage.value = ''

  try {
    const [
      theatersRes,
      screensRes,
      seatsRes,
    ] = await Promise.all([
      api.get<Theater[]>('/theaters'),
      api.get<Screen[]>('/screens'),
      api.get<Seat[]>('/seats'),
    ])

    theaters.value = theatersRes.data
    screens.value = screensRes.data
    seats.value = seatsRes.data
  } catch (err) {
    errorMessage.value =
      'Failed to load seats. Please try again.'

    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)

async function reloadSeats() {
  const { data } =
    await api.get<Seat[]>('/seats')

  seats.value = data
}

// -----------------------------------------------------------------------
// Theater -> Screen
// -----------------------------------------------------------------------

const filteredScreens = computed(() => {
  if (!selectedTheaterId.value) return []

  return screens.value.filter(
    (screen) =>
      screen.theater_id ===
      selectedTheaterId.value
  )
})

function onTheaterChange() {
  selectedScreenId.value = ''

  // Reset edit state when changing theater
  closeForm()
}

const selectedScreen = computed(() => {
  return (
    screens.value.find(
      (screen) =>
        screen.screen_id ===
        selectedScreenId.value
    ) ?? null
  )
})

const selectedTheater = computed(() => {
  return (
    theaters.value.find(
      (theater) =>
        theater.theater_id ===
        selectedTheaterId.value
    ) ?? null
  )
})

// -----------------------------------------------------------------------
// Seat parsing
// -----------------------------------------------------------------------

function parseSeatNumber(
  seatNumber: string
) {
  const match =
    seatNumber.match(/^([A-Za-z]+)(\d+)$/)

  if (!match) {
    return {
      row: seatNumber,
      number: 0,
    }
  }

  return {
    row: match[1],
    number: Number(match[2]),
  }
}

// -----------------------------------------------------------------------
// Seats for selected screen
// -----------------------------------------------------------------------

const screenSeats = computed(() => {
  if (!selectedScreenId.value) return []

  return seats.value.filter(
    (seat) =>
      seat.screen_id ===
      selectedScreenId.value
  )
})

const rows = computed(() => {
  const rowSet = new Set(
    screenSeats.value.map(
      (seat) =>
        parseSeatNumber(
          seat.seat_number
        ).row
    )
  )

  return [...rowSet].sort()
})

function seatsInRow(row: string) {
  return screenSeats.value
    .filter(
      (seat) =>
        parseSeatNumber(
          seat.seat_number
        ).row === row
    )
    .sort(
      (a, b) =>
        parseSeatNumber(
          a.seat_number
        ).number -
        parseSeatNumber(
          b.seat_number
        ).number
    )
}

// -----------------------------------------------------------------------
// VIP row
// -----------------------------------------------------------------------

function isVipRow(row: string) {
  return seatsInRow(row).some(
    (seat) =>
      seat.seat_type === 'vip'
  )
}

const labelOrderRows = computed(() => {
  const vip = rows.value.filter(isVipRow)

  const rest = rows.value.filter(
    (row) => !isVipRow(row)
  )

  return [...vip, ...rest]
})

const renderRows = computed(() => {
  return [...rows.value].reverse()
})

const rowDisplayLabel = computed(() => {
  const map = new Map<
    string,
    string
  >()

  labelOrderRows.value.forEach(
    (row, index) => {
      map.set(
        row,
        String.fromCharCode(
          65 + index
        )
      )
    }
  )

  return map
})

function displayLabel(row: string) {
  return (
    rowDisplayLabel.value.get(row) ??
    row
  )
}

// -----------------------------------------------------------------------
// Add / Edit seat
// -----------------------------------------------------------------------

function openAddForm() {
  if (!selectedScreenId.value) return

  editingSeat.value = null

  selectedSeatIds.value = []
  batchSeatType.value = null

  form.value = {
    row: '',
    number: 1,
    seat_type: 'standard',
    status: 'AVAILABLE',
  }

  showForm.value = true
}

function editSeat(seat: Seat) {
  // ---------------------------------------------------------------
  // ถ้ายังไม่ได้อยู่ใน Edit Seat
  // ให้เปิด Edit Seat ตามปกติ
  // ---------------------------------------------------------------
  if (!showForm.value || !editingSeat.value) {
    editingSeat.value = seat

    selectedSeatIds.value = [seat.seat_id]

    batchSeatType.value =
      seat.seat_type

    const parsed =
      parseSeatNumber(
        seat.seat_number
      )

    form.value = {
      row: parsed.row,
      number: parsed.number,
      seat_type: seat.seat_type,
      status: seat.status,
    }

    showForm.value = true

    return
  }

  // ---------------------------------------------------------------
  // ถ้า Edit Seat เปิดอยู่แล้ว
  //
  // คลิกที่นั่งบน Seat Map =
  // เพิ่ม/เอาออกจากรายการที่ต้องการแก้
  // ---------------------------------------------------------------

  if (
    seat.status.toUpperCase() ===
    'BOOKED'
  ) {
    return
  }

  if (
    seat.status.toUpperCase() ===
    'DISABLED'
  ) {
    return
  }

  toggleSeatSelection(seat)
}

// -----------------------------------------------------------------------
// เลือก Seat Type
// -----------------------------------------------------------------------

function selectSeatType(
  type: SeatType
) {
  batchSeatType.value = type

  form.value.seat_type = type

  // ถ้ายังไม่มี selected seat
  // ให้เลือก seat ที่กำลัง Edit อยู่ก่อน
  if (
    editingSeat.value &&
    selectedSeatIds.value.length === 0
  ) {
    selectedSeatIds.value = [
      editingSeat.value.seat_id,
    ]
  }
}

// -----------------------------------------------------------------------
// Click seat while Edit Seat is open
// -----------------------------------------------------------------------

function toggleSeatSelection(
  seat: Seat
) {
  // ไม่ให้เลือก BOOKED / DISABLED
  const normalized =
    seat.status.toUpperCase()

  if (
    normalized === 'BOOKED' ||
    normalized === 'DISABLED'
  ) {
    return
  }

  // ต้องเลือก Seat Type ก่อน
  if (!batchSeatType.value) {
    return
  }

  const index =
    selectedSeatIds.value.indexOf(
      seat.seat_id
    )

  if (index === -1) {
    // เพิ่มที่นั่ง
    selectedSeatIds.value.push(
      seat.seat_id
    )
  } else {
    // ถ้ากดซ้ำ ให้เอาออก
    selectedSeatIds.value.splice(
      index,
      1
    )
  }
}

// -----------------------------------------------------------------------
// Check selected seat
// -----------------------------------------------------------------------

function isSeatSelected(
  seat: Seat
) {
  return selectedSeatIds.value.includes(
    seat.seat_id
  )
}

// -----------------------------------------------------------------------
// Get seat preview type
//
// ถ้าที่นั่งถูกเลือกอยู่ ให้ใช้ประเภทที่เลือกจาก Edit Seat
// ทำให้กด Premium แล้วคลิกเลข 1,2,3,4... เป็นสีม่วงทันที
// -----------------------------------------------------------------------

function getSeatPreviewType(
  seat: Seat
): SeatType {
  if (
    isSeatSelected(seat) &&
    batchSeatType.value
  ) {
    return batchSeatType.value
  }

  return seat.seat_type
}

// -----------------------------------------------------------------------
// Seat type button class
// -----------------------------------------------------------------------

function seatTypeButtonClass(
  type: SeatType
) {
  const selected =
    batchSeatType.value === type

  if (type === 'standard') {
    return selected
      ? 'border-green-500 bg-green-500/20 text-green-300 ring-2 ring-green-500/30'
      : 'border-border/20 bg-surface text-text hover:bg-green-500/10 hover:text-green-300'
  }

  if (type === 'premium') {
    return selected
      ? 'border-purple-500 bg-purple-500/20 text-purple-300 ring-2 ring-purple-500/30'
      : 'border-border/20 bg-surface text-text hover:bg-purple-500/10 hover:text-purple-300'
  }

  return selected
    ? 'border-blue-500 bg-blue-500/20 text-blue-300 ring-2 ring-blue-500/30'
    : 'border-border/20 bg-surface text-text hover:bg-blue-500/10 hover:text-blue-300'
}

// -----------------------------------------------------------------------
// Close form
// -----------------------------------------------------------------------

function closeForm() {
  showForm.value = false
  editingSeat.value = null

  selectedSeatIds.value = []
  batchSeatType.value = null
}

// -----------------------------------------------------------------------
// Save seat / multiple seats
// -----------------------------------------------------------------------

async function saveSeat() {
  if (
    !selectedScreenId.value ||
    !batchSeatType.value
  ) {
    return
  }

  if (
    selectedSeatIds.value.length === 0
  ) {
    errorMessage.value =
      'Please select at least one seat.'
    return
  }

  errorMessage.value = ''
  successMessage.value = ''

  try {
    // ---------------------------------------------------------------
    // Update ทุก seat ที่เลือก
    // ---------------------------------------------------------------
    await Promise.all(
      selectedSeatIds.value.map(
        async (seatId) => {
          // หา seat ตัวจริงจาก seats ที่โหลดมาจาก backend
          const seat = seats.value.find(
            (s) => s.seat_id === seatId
          )

          if (!seat) {
            throw new Error(
              `Seat ${seatId} not found`
            )
          }

          // Backend ต้องการ seat_number
          // เช่น A1, A2, A3
          await api.put(
            `/seats/${seatId}`,
            {
              screen_id:
                selectedScreenId.value,

              seat_number:
                seat.seat_number,

              seat_type:
                batchSeatType.value,
            }
          )
        }
      )
    )

    const count =
      selectedSeatIds.value.length

    await reloadSeats()

    closeForm()

    successMessage.value =
      `${count} seat${count > 1 ? 's' : ''} updated successfully.`

  } catch (err: any) {
    // แสดง error จาก Backend ถ้ามี
    errorMessage.value =
      err.response?.data?.error ||
      'Failed to save seat changes. Please try again.'

    console.error(
      'Failed to save seat changes:',
      err
    )
  }
}

// -----------------------------------------------------------------------
// Delete single seat
// -----------------------------------------------------------------------

async function deleteSelectedSeats() {
  if (selectedSeatIds.value.length === 0) {
    errorMessage.value = 'Please select at least one seat.'
    return
  }

  const confirmed = window.confirm(
    `Are you sure you want to delete ${selectedSeatIds.value.length} selected seat(s)?`
  )

  if (!confirmed) {
    return
  }

  try {
    errorMessage.value = ''
    successMessage.value = ''

    await Promise.all(
      selectedSeatIds.value.map((seatId) =>
        api.delete(`/seats/${seatId}`)
      )
    )

    const count = selectedSeatIds.value.length

    selectedSeatIds.value = []

    closeForm()

    await reloadSeats()

    successMessage.value =
      `${count} seat${count > 1 ? 's' : ''} deleted successfully.`

  } catch (err: any) {
    console.error(
      'Failed to delete selected seats:',
      err
    )

    errorMessage.value =
      err.response?.data?.error ||
      'Failed to delete selected seats. Please try again.'
  }
}

// -----------------------------------------------------------------------
// Seat color
// -----------------------------------------------------------------------

function seatClass(
  status: string,
  seatType: SeatType,
  selected = false
) {
  const normalized =
    status.toUpperCase()

  // ---------------------------------------------------------------
  // ถ้า seat ถูกเลือกจาก Edit Seat
  // ให้เปลี่ยนเป็นสีตาม Seat Type ที่เลือกทันที
  // ---------------------------------------------------------------

  if (selected && batchSeatType.value) {
    if (
      batchSeatType.value ===
      'vip'
    ) {
      return 'border-blue-500 bg-blue-500/30 text-blue-300 ring-2 ring-blue-500/50'
    }

    if (
      batchSeatType.value ===
      'premium'
    ) {
      return 'border-purple-500 bg-purple-500/30 text-purple-300 ring-2 ring-purple-500/50'
    }

    return 'border-green-500 bg-green-500/30 text-green-300 ring-2 ring-green-500/50'
  }

  // BOOKED
  if (normalized === 'BOOKED') {
    return 'border-primary/40 bg-primary/20 text-primary'
  }

  // DISABLED
  if (normalized === 'DISABLED') {
    return 'border-border/20 bg-white/5 text-muted'
  }

  // AVAILABLE
  if (seatType === 'vip') {
    return 'border-blue-500/30 bg-blue-500/10 text-blue-400'
  }

  if (seatType === 'premium') {
    return 'border-purple-500/30 bg-purple-500/10 text-purple-400'
  }

  return 'border-green-500/20 bg-green-500/10 text-green-400'
}
</script>

<template>
  <div
    class="mx-auto max-w-7xl px-6 py-8"
  >
    <!-- Page Header -->
    <div
      class="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between"
    >
      <div
        class="flex items-center gap-4"
      >
        <div
          class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-gradient-to-br from-purple-500 to-pink-500 shadow-lg shadow-purple-500/20"
        >
          <Armchair
            class="size-6"
          />
        </div>

        <div>
          <h1
            class="text-2xl font-bold text-white"
          >
            Manage Seats
          </h1>

          <p
            class="mt-1 text-sm text-gray-400"
          >
            Arrange seating
          </p>
        </div>
      </div>
    </div>

    <!-- Error -->
    <div
      v-if="errorMessage"
      class="mt-5 rounded-lg border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400"
    >
      {{ errorMessage }}
    </div>

    <!-- Success -->
    <div
      v-if="successMessage"
      class="mt-5 rounded-lg border border-green-500/30 bg-green-500/10 px-4 py-3 text-sm text-green-400"
    >
      {{ successMessage }}
    </div>

    <!-- Filters -->
    <div
      class="mt-7 flex flex-col gap-4 sm:flex-row sm:items-end"
    >
      <!-- Theater -->
      <div
        class="w-full sm:max-w-xs"
      >
        <label
          class="mb-2 block text-xs font-medium text-muted"
        >
          Theater
        </label>

        <select
          v-model="selectedTheaterId"
          class="h-11 w-full rounded-lg border border-border/20 bg-surface px-3 text-sm text-text outline-none transition focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
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
      <div
        class="w-full sm:max-w-xs"
      >
        <label
          class="mb-2 block text-xs font-medium text-muted"
        >
          Screen
        </label>

        <select
          v-model="selectedScreenId"
          :disabled="!selectedTheaterId"
          class="h-11 w-full rounded-lg border border-border/20 bg-surface px-3 text-sm text-text outline-none transition focus:border-primary/60 focus:ring-2 focus:ring-primary/10 disabled:cursor-not-allowed disabled:opacity-50"
        >
          <option value="">
            Select screen
          </option>

          <option
            v-for="screen in filteredScreens"
            :key="screen.screen_id"
            :value="screen.screen_id"
          >
            {{ screen.name }}
          </option>
        </select>
      </div>

      <!-- Add -->
      <button
        type="button"
        class="h-11 rounded-lg bg-[#ff3158] px-5 text-sm font-semibold text-white transition hover:bg-[#ff4669] disabled:cursor-not-allowed disabled:opacity-50"
        :disabled="!selectedScreenId"
        @click="openAddForm"
      >
        + Add seat
      </button>

      <!-- Generate -->
      <button
        type="button"
        class="h-11 rounded-lg border border-white/15 px-5 text-sm font-semibold text-white transition hover:bg-white/5 disabled:cursor-not-allowed disabled:opacity-50"
        :disabled="!selectedScreenId"
        @click="openGenerateForm"
      >
        ⚡ Generate seats
      </button>
    </div>

    <!-- Seat Map -->
    <div
      class="mt-7 rounded-2xl border border-border/20 bg-card bg-[#15151d] p-6"
    >
      <!-- Header -->
      <div
        class="mb-6 flex items-center justify-between"
      >
        <div>
          <h2
            class="text-lg font-bold text-text"
          >
            Seat Map
          </h2>

          <p
            class="mt-1 text-xs text-muted"
          >
            {{
              selectedScreen &&
              selectedTheater
                ? `${selectedScreen.name} · ${selectedTheater.name}`
                : 'Select a theater and screen to view seats'
            }}
          </p>
        </div>

        <div
          class="text-xs text-muted"
        >
          {{ screenSeats.length }}
          seats
        </div>
      </div>

      <template
        v-if="selectedScreenId"
      >
        <!-- Screen -->
        <div
          class="mb-10 flex justify-center"
        >
          <div
            class="w-full max-w-2xl rounded-full border border-primary/30 bg-primary/5 py-3 text-center text-xs font-semibold tracking-widest text-primary"
          >
            SCREEN
          </div>
        </div>

        <!-- Seats -->
        <div
          class="space-y-5"
        >
          <div
            v-for="row in renderRows"
            :key="row"
            class="flex items-center justify-center gap-3"
          >
            <!-- Row -->
            <div
              class="w-6 text-center text-sm font-semibold text-muted"
            >
              {{ displayLabel(row) }}
            </div>

            <!-- Seat -->
            <button
              v-for="seat in seatsInRow(row)"
              :key="seat.seat_id"
              type="button"
              :class="[
                'group relative flex h-11 w-11 items-center justify-center rounded-lg border text-xs font-semibold transition hover:-translate-y-0.5',
                seatClass(
                  seat.status,
                  getSeatPreviewType(seat),
                  isSeatSelected(seat)
                ),
                showForm &&
                batchSeatType &&
                !isSeatSelected(seat) &&
                seat.status.toUpperCase() !== 'BOOKED' &&
                seat.status.toUpperCase() !== 'DISABLED'
                  ? 'cursor-pointer'
                  : '',
              ]"
              @click="editSeat(seat)"
            >
              {{ parseSeatNumber(seat.seat_number).number }}

              <!-- Tooltip -->
              <span
                class="pointer-events-none absolute -top-9 left-1/2 hidden -translate-x-1/2 whitespace-nowrap rounded-md bg-black px-2 py-1 text-[10px] text-white shadow-xl group-hover:block"
              >
                {{ seat.seat_number }}
                ·
                {{
                  isSeatSelected(seat) &&
                  batchSeatType
                    ? batchSeatType
                    : seat.seat_type
                }}
              </span>
            </button>
          </div>

          <p
            v-if="rows.length === 0"
            class="py-8 text-center text-sm text-muted"
          >
            No seats yet for this screen.
            Click "Add seat" to create one.
          </p>
        </div>

        <!-- Selection Info -->
        <div
          v-if="
            showForm &&
            editingSeat &&
            batchSeatType
          "
          class="mt-8 rounded-xl border border-purple-500/30 bg-purple-500/10 px-5 py-4"
        >
          <div
            class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between"
          >
            <div>
              <p
                class="text-sm font-semibold text-purple-300"
              >
                {{ batchSeatType.toUpperCase() }}
                selected
              </p>

              <p
                class="mt-1 text-xs text-purple-200/70"
              >
                คลิกเลขที่นั่งที่ต้องการเปลี่ยนประเภทได้เลย
                แล้วกด Save changes ครั้งเดียว
              </p>
            </div>

            <div
              class="text-sm font-semibold text-white"
            >
              {{
                selectedSeatIds.length
              }}
              selected
            </div>
          </div>
        </div>

        <!-- Legend -->
        <div
          class="mt-8 flex flex-wrap justify-center gap-6 border-t border-border/20 pt-6 text-xs text-muted"
        >
          <!-- Standard -->
          <div
            class="flex items-center gap-2"
          >
            <span
              class="size-3 rounded border border-green-500/20 bg-green-500/10"
            ></span>

            Available (Standard)
          </div>

          <!-- Premium -->
          <div
            class="flex items-center gap-2"
          >
            <span
              class="size-3 rounded border border-purple-500/30 bg-purple-500/10"
            ></span>

            Available (Premium)
          </div>

          <!-- VIP -->
          <div
            class="flex items-center gap-2"
          >
            <span
              class="size-3 rounded border border-blue-500/30 bg-blue-500/10"
            ></span>

            Available (VIP)
          </div>

          <!-- Booked -->
          <div
            class="flex items-center gap-2"
          >
            <span
              class="size-3 rounded border border-primary/40 bg-primary/20"
            ></span>

            Booked
          </div>

          <!-- Disabled -->
          <div
            class="flex items-center gap-2"
          >
            <span
              class="size-3 rounded border border-border/20 bg-white/5"
            ></span>

            Disabled
          </div>
        </div>
      </template>

      <p
        v-else
        class="py-12 text-center text-sm text-muted"
      >
        Choose a theater and screen above to see its seat map.
      </p>
    </div>

    <!-- Add / Edit Form -->
    <div
      v-if="showForm"
      class="mx-auto mt-6 max-w-3xl rounded-2xl border border-border/20 bg-card p-6 sm:p-7"
    >
      <!-- Form Header -->
      <div
        class="mb-6 flex items-center justify-between"
      >
        <div>
          <h2
            class="text-xl font-bold tracking-tight text-text"
          >
            {{ editingSeat ? 'EDIT SEAT' : 'ADD SEAT' }}
          </h2>

          <p
            class="mt-1 text-xs text-muted"
          >
            {{
              editingSeat
                ? 'แก้ไขประเภทที่นั่งหลายตัวพร้อมกัน'
                : 'เพิ่มที่นั่งใหม่เข้าสู่ระบบ'
            }}
          </p>
        </div>

        <button
          type="button"
          class="rounded-lg p-2 text-muted transition hover:bg-white/5 hover:text-text"
          @click="closeForm"
        >
          ✕
        </button>
      </div>

      <!-- ========================================================= -->
      <!-- EDIT SEAT -->
      <!-- ========================================================= -->

      <template v-if="editingSeat">
        <!-- Selected seats -->
        <div
          class="mb-6 rounded-lg border border-white/10 bg-white/5 px-4 py-3"
        >
          <p
            class="text-xs text-muted"
          >
            Selected seats
          </p>

          <p
            class="mt-1 text-sm font-semibold text-white"
          >
            {{
              selectedSeatIds.length
            }}
            seat{{
              selectedSeatIds.length > 1
                ? 's'
                : ''
            }}
          </p>
        </div>

        <!-- Row / Number -->
        <div
          class="grid gap-x-5 gap-y-5 sm:grid-cols-2"
        >
          <!-- Row -->
          <div>
            <label
              class="mb-2 block text-xs font-medium text-muted"
            >
              Row
            </label>

            <input
              v-model="form.row"
              type="text"
              maxlength="1"
              placeholder="e.g. A"
              class="h-11 w-full rounded-lg border border-border/20 bg-surface px-4 text-sm uppercase text-text outline-none transition placeholder:text-muted focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
              readonly
            />
          </div>

          <!-- Seat Number -->
          <div>
            <label
              class="mb-2 block text-xs font-medium text-muted"
            >
              Seat Number
            </label>

            <input
              v-model.number="form.number"
              type="number"
              min="1"
              placeholder="e.g. 1"
              class="h-11 w-full rounded-lg border border-border/20 bg-surface px-4 text-sm text-text outline-none transition placeholder:text-muted focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
              readonly
            />
          </div>
        </div>

        <!-- Seat Type -->
        <div
          class="mt-5"
        >
          <label
            class="mb-3 block text-xs font-medium text-muted"
          >
            Seat type
          </label>

          <div
            class="grid grid-cols-3 gap-2"
          >
            <!-- Standard -->
            <button
              type="button"
              :class="[
                'h-12 rounded-lg border text-sm font-semibold transition',
                seatTypeButtonClass('standard'),
              ]"
              @click="
                selectSeatType('standard')
              "
            >
              Standard
            </button>

            <!-- Premium -->
            <button
              type="button"
              :class="[
                'h-12 rounded-lg border text-sm font-semibold transition',
                seatTypeButtonClass('premium'),
              ]"
              @click="
                selectSeatType('premium')
              "
            >
              Premium
            </button>

            <!-- VIP -->
            <button
              type="button"
              :class="[
                'h-12 rounded-lg border text-sm font-semibold transition',
                seatTypeButtonClass('vip'),
              ]"
              @click="
                selectSeatType('vip')
              "
            >
              VIP
            </button>
          </div>

          <!-- Instruction -->
          <p
            v-if="batchSeatType"
            class="mt-3 rounded-lg border border-purple-500/20 bg-purple-500/10 px-4 py-3 text-xs text-purple-300"
          >
            เลือก
            <strong>
              {{
                batchSeatType.toUpperCase()
              }}
            </strong>
            แล้วสามารถคลิกเลขที่นั่งบน Seat Map ได้หลายตัว
            เช่น 1, 2, 3, 4, 5, 6
            แล้วค่อยกด Save changes
          </p>
        </div>

        <!-- Status -->
        <div
          class="mt-5"
        >
          <p
            class="rounded-lg border border-amber-500/20 bg-amber-500/10 px-4 py-3 text-xs text-amber-300"
          >
            Status (currently
            <strong>
              {{ editingSeat.status }}
            </strong>
            ) can't be changed after a seat is created — the backend doesn't support updating it here.
          </p>
        </div>

        <!-- Buttons -->
        <div
          class="mt-7 flex items-center justify-between gap-3"
        >
          <!-- Delete -->
          <button
            type="button"
             class="inline-flex items-center gap-2 rounded-xl bg-[#ff3158] px-6 py-3 text-sm font-semibold text-white transition hover:bg-[#ff4669] disabled:cursor-not-allowed disabled:opacity-50"
            @click="deleteSelectedSeats"
          >
            Delete seat
          </button>

          <div
            class="flex gap-3"
          >
            <button
              type="button"
              class="rounded-lg border border-border/20 px-5 py-2.5 text-sm font-medium text-muted transition hover:bg-white/5 hover:text-text"
              @click="closeForm"
            >
              Cancel
            </button>

            <button
              type="button"
              :disabled="
                !batchSeatType ||
                selectedSeatIds.length === 0
              "
              class="rounded-lg bg-[#ff3158] px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-[#ff4669] disabled:cursor-not-allowed disabled:opacity-50"
              @click="saveSeat"
            >
              Save changes
            </button>
          </div>
        </div>
      </template>

      <!-- ========================================================= -->
      <!-- ADD SEAT -->
      <!-- ========================================================= -->

      <template v-else>
        <div
          class="grid gap-x-5 gap-y-5 sm:grid-cols-2"
        >
          <!-- Row -->
          <div>
            <label
              class="mb-2 block text-xs font-medium text-muted"
            >
              Row
            </label>

            <input
              v-model="form.row"
              type="text"
              maxlength="1"
              placeholder="e.g. A"
              class="h-11 w-full rounded-lg border border-border/20 bg-surface px-4 text-sm uppercase text-text outline-none transition placeholder:text-muted focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
            />
          </div>

          <!-- Seat Number -->
          <div>
            <label
              class="mb-2 block text-xs font-medium text-muted"
            >
              Seat Number
            </label>

            <input
              v-model.number="form.number"
              type="number"
              min="1"
              placeholder="e.g. 1"
              class="h-11 w-full rounded-lg border border-border/20 bg-surface px-4 text-sm text-text outline-none transition placeholder:text-muted focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
            />
          </div>

          <!-- Seat Type -->
          <div>
            <label
              class="mb-2 block text-xs font-medium text-muted"
            >
              Seat type
            </label>

            <select
              v-model="form.seat_type"
              class="h-11 w-full rounded-lg border border-border/20 bg-surface px-4 text-sm text-text outline-none transition focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
            >
              <option value="standard">
                Standard
              </option>

              <option value="premium">
                Premium
              </option>

              <option value="vip">
                VIP
              </option>
            </select>
          </div>

          <!-- Status -->
          <div>
            <label
              class="mb-2 block text-xs font-medium text-muted"
            >
              Status
            </label>

            <select
              v-model="form.status"
              class="h-11 w-full rounded-lg border border-border/20 bg-surface px-4 text-sm text-text outline-none transition focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
            >
              <option value="AVAILABLE">
                Available
              </option>

              <option value="DISABLED">
                Disabled
              </option>
            </select>

            <p
              class="mt-1.5 text-xs text-muted"
            >
              "Booked" is set automatically when a customer books this seat.
            </p>
          </div>
        </div>

        <!-- Buttons -->
        <div
          class="mt-7 flex justify-end gap-3"
        >
          <button
            type="button"
            class="rounded-lg border border-border/20 px-5 py-2.5 text-sm font-medium text-muted transition hover:bg-white/5 hover:text-text"
            @click="closeForm"
          >
            Cancel
          </button>

          <button
            type="button"
            class="rounded-lg bg-[#ff3158] px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-[#ff4669]"
            @click="saveSeat"
          >
            Save seat
          </button>
        </div>
      </template>
    </div>

    <!-- Generate Seats Form -->
    <div
      v-if="showGenerateForm"
      class="mx-auto mt-6 max-w-2xl rounded-2xl border border-border/20 bg-card p-6 sm:p-7"
    >
      <!-- Header -->
      <div
        class="mb-6 flex items-center justify-between"
      >
        <div>
          <h2
            class="text-xl font-bold tracking-tight text-text"
          >
            GENERATE SEATS
          </h2>

          <p
            class="mt-1 text-xs text-muted"
          >
            สร้างที่นั่งหลายแถว/หลายที่นั่งพร้อมกันในครั้งเดียว
          </p>
        </div>

        <button
          type="button"
          class="rounded-lg p-2 text-muted transition hover:bg-white/5 hover:text-text"
          @click="closeGenerateForm"
        >
          ✕
        </button>
      </div>

      <div
        class="grid gap-x-5 gap-y-5 sm:grid-cols-2"
      >
        <!-- Start Row -->
        <div>
          <label
            class="mb-2 block text-xs font-medium text-muted"
          >
            Start row
          </label>

          <input
            v-model="generateForm.start_row"
            type="text"
            maxlength="1"
            placeholder="e.g. A"
            class="h-11 w-full rounded-lg border border-border/20 bg-surface px-4 text-sm uppercase text-text outline-none transition placeholder:text-muted focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
          />
        </div>

        <!-- End Row -->
        <div>
          <label
            class="mb-2 block text-xs font-medium text-muted"
          >
            End row
          </label>

          <input
            v-model="generateForm.end_row"
            type="text"
            maxlength="1"
            placeholder="e.g. L"
            class="h-11 w-full rounded-lg border border-border/20 bg-surface px-4 text-sm uppercase text-text outline-none transition placeholder:text-muted focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
          />
        </div>

        <!-- Seats per row -->
        <div>
          <label
            class="mb-2 block text-xs font-medium text-muted"
          >
            Seats per row
          </label>

          <input
            v-model.number="
              generateForm.seats_per_row
            "
            type="number"
            min="1"
            max="99"
            placeholder="e.g. 14"
            class="h-11 w-full rounded-lg border border-border/20 bg-surface px-4 text-sm text-text outline-none transition placeholder:text-muted focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
          />
        </div>

        <!-- Seat Type -->
        <div>
          <label
            class="mb-2 block text-xs font-medium text-muted"
          >
            Seat type
          </label>

          <select
            v-model="
              generateForm.seat_type
            "
            class="h-11 w-full rounded-lg border border-border/20 bg-surface px-4 text-sm text-text outline-none transition focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
          >
            <option value="standard">
              Standard
            </option>

            <option value="premium">
              Premium
            </option>

            <option value="vip">
              VIP
            </option>
          </select>
        </div>

        <!-- Status -->
        <div
          class="sm:col-span-2"
        >
          <label
            class="mb-2 block text-xs font-medium text-muted"
          >
            Status
          </label>

          <select
            v-model="
              generateForm.status
            "
            class="h-11 w-full rounded-lg border border-border/20 bg-surface px-4 text-sm text-text outline-none transition focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
          >
            <option value="AVAILABLE">
              Available
            </option>

            <option value="DISABLED">
              Disabled
            </option>
          </select>
        </div>
      </div>

      <!-- Generate Info -->
      <p
        class="mt-5 rounded-lg border border-blue-500/20 bg-blue-500/10 px-4 py-3 text-xs text-blue-300"
      >
        จะสร้างที่นั่งตั้งแต่แถว
        <strong>
          {{
            generateForm.start_row ||
            '?'
          }}
        </strong>
        ถึง
        <strong>
          {{
            generateForm.end_row ||
            '?'
          }}
        </strong>
        แถวละ
        <strong>
          {{
            generateForm.seats_per_row ||
            0
          }}
        </strong>
        ที่นั่ง
      </p>

      <!-- Buttons -->
      <div
        class="mt-7 flex justify-end gap-3"
      >
        <button
          type="button"
          class="rounded-lg border border-border/20 px-5 py-2.5 text-sm font-medium text-muted transition hover:bg-white/5 hover:text-text"
          @click="closeGenerateForm"
        >
          Cancel
        </button>

        <button
          type="button"
          :disabled="generating"
          class="rounded-lg bg-[#ff3158] px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-[#ff4669] disabled:cursor-not-allowed disabled:opacity-60"
          @click="generateSeats"
        >
          {{
            generating
              ? 'Generating...'
              : 'Generate seats'
          }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
select {
  color-scheme: dark;
}

select option {
  background-color: #15151d;
  color: #e5e7eb;
}
</style>