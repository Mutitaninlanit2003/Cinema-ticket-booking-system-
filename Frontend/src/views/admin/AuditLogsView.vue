<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import axios from 'axios'
import { ClipboardList } from 'lucide-vue-next'
import {
  Search,
  CalendarDays,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  Plus,
  Pencil,
  Trash2,
  LogIn,
  LogOut,
} from 'lucide-vue-next'

// -----------------------------------------------------------------------
// API client
// NOTE: replace with your project's central API client if you have one
// (e.g. `import api from '@/api/client'`).
// -----------------------------------------------------------------------
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
})

// -----------------------------------------------------------------------
// Types — match backend models exactly (models/audit_log.go, models/user.go)
// -----------------------------------------------------------------------
interface AuditLog {
  id: number
  audit_log_id: string
  user_id: string
  action: string
  resource: string
  resource_id: string
  details: string
  created_at: string
}

interface User {
  id: number
  user_id: string
  name: string
  email: string
  role: string
}

const searchQuery = ref('')
const selectedAction = ref('All actions')
const selectedUser = ref('All users')
const selectedDate = ref('')

const currentPage = ref(1)
const itemsPerPage = 8

const logs = ref<AuditLog[]>([])
const users = ref<User[]>([])

const loading = ref(false)
const errorMessage = ref('')

// user_id -> User lookup, since audit logs only carry user_id (uuid)
const userById = computed(() => {
  const map = new Map<string, User>()
  users.value.forEach((u) => map.set(u.user_id, u))
  return map
})

function userName(userId: string) {
  return userById.value.get(userId)?.name ?? userId
}

function formatTime(isoString: string) {
  const date = new Date(isoString)

  if (Number.isNaN(date.getTime())) {
    return isoString
  }

  return date.toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  })
}

const actionOptions = computed(() => {
  const distinctActions = new Set(logs.value.map((log) => log.action))
  return ['All actions', ...distinctActions]
})

const userOptions = computed(() => {
  return ['All users', ...users.value.map((u) => u.user_id)]
})

// -----------------------------------------------------------------------
// Data loading
// -----------------------------------------------------------------------
async function loadUsers() {
  const { data } = await api.get<User[]>('/users')
  users.value = data
}

async function loadLogs() {
  const { data } = await api.get<AuditLog[]>('/audit_logs')
  logs.value = data
}

async function loadAll() {
  loading.value = true
  errorMessage.value = ''

  try {
    await Promise.all([loadUsers(), loadLogs()])
  } catch (err) {
    errorMessage.value = 'Failed to load audit logs. Please try again.'
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)

const filteredLogs = computed(() => {
  const query = searchQuery.value.toLowerCase().trim()

  return logs.value
    .filter((log) => {
      const name = userName(log.user_id)

      const matchesSearch =
        !query ||
        name.toLowerCase().includes(query) ||
        log.action.toLowerCase().includes(query) ||
        log.resource.toLowerCase().includes(query) ||
        log.details.toLowerCase().includes(query)

      const matchesAction =
        selectedAction.value === 'All actions' ||
        log.action === selectedAction.value

      const matchesUser =
        selectedUser.value === 'All users' ||
        log.user_id === selectedUser.value

      const matchesDate =
        !selectedDate.value ||
        log.created_at.slice(0, 10) === selectedDate.value

      return matchesSearch && matchesAction && matchesUser && matchesDate
    })
    .sort(
      (a, b) =>
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
    )
})

const totalPages = computed(() => {
  return Math.max(
    1,
    Math.ceil(filteredLogs.value.length / itemsPerPage)
  )
})

const paginatedLogs = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage

  return filteredLogs.value.slice(
    start,
    start + itemsPerPage
  )
})

const showingFrom = computed(() => {
  if (filteredLogs.value.length === 0) {
    return 0
  }

  return (currentPage.value - 1) * itemsPerPage + 1
})

const showingTo = computed(() => {
  return Math.min(
    currentPage.value * itemsPerPage,
    filteredLogs.value.length
  )
})

function resetFilters() {
  searchQuery.value = ''
  selectedAction.value = 'All actions'
  selectedUser.value = 'All users'
  selectedDate.value = ''
  currentPage.value = 1
}

function changePage(page: number) {
  if (page < 1 || page > totalPages.value) {
    return
  }

  currentPage.value = page
}

function getActionClass(action: string) {
  switch (action) {
    case 'CREATE':
      return 'bg-emerald-500/15 text-emerald-400'

    case 'UPDATE':
      return 'bg-amber-500/15 text-amber-400'

    case 'DELETE':
      return 'bg-red-500/15 text-red-400'

    case 'LOGIN':
      return 'bg-blue-500/15 text-blue-400'

    case 'LOGOUT':
      return 'bg-gray-500/15 text-gray-400'

    default:
      return 'bg-white/10 text-gray-300'
  }
}

function getActionIcon(action: string) {
  switch (action) {
    case 'CREATE':
      return Plus

    case 'UPDATE':
      return Pencil

    case 'DELETE':
      return Trash2

    case 'LOGIN':
      return LogIn

    case 'LOGOUT':
      return LogOut

    default:
      return ClipboardList
  }
}
</script>

<template>
  <div class="mx-auto max-w-7xl px-6 py-10">

    <!-- Header -->
    <div class="mb-6">
      <div class="flex items-center gap-4">

        <!-- Audit Log Icon -->
        <div
          class="flex size-12 shrink-0 items-center justify-center rounded-xl bg-gradient-to-br from-purple-500 to-pink-500 shadow-lg shadow-purple-500/20"
        >
          <ClipboardList class="size-6" />
        </div>

        <div>
          <h1 class="text-2xl font-bold text-white">
            Audit Logs
          </h1>

          <p class="mt-1 text-sm text-gray-400">
            Check the system's transaction history
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

    <!-- Filters -->
    <div
      class="mb-6 rounded-2xl border border-white/10 bg-[#15151d] p-5"
    >
      <div
        class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-[1.3fr_1fr_1fr_1.1fr_auto]"
      >

        <!-- Search -->
        <div>
          <label
            class="mb-2 block text-xs font-medium text-gray-400"
          >
            Search
          </label>

          <div class="relative">
            <Search
              class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-500"
            />

            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search logs..."
              class="w-full rounded-xl border border-white/10 bg-[#0f0f15] py-2.5 pl-10 pr-4 text-sm text-white outline-none placeholder:text-gray-600 transition focus:border-pink-500/50"
              @input="currentPage = 1"
            />
          </div>
        </div>

        <!-- Action -->
        <div>
          <label
            class="mb-2 block text-xs font-medium text-gray-400"
          >
            Action
          </label>

          <div class="relative">
            <select
              v-model="selectedAction"
              class="w-full appearance-none rounded-xl border border-white/10 bg-[#0f0f15] px-4 py-2.5 pr-10 text-sm text-gray-300 outline-none focus:border-pink-500/50"
              @change="currentPage = 1"
            >
              <option
                v-for="action in actionOptions"
                :key="action"
                :value="action"
              >
                {{ action }}
              </option>
            </select>

            <ChevronDown
              class="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-500"
            />
          </div>
        </div>

        <!-- User -->
        <div>
          <label
            class="mb-2 block text-xs font-medium text-gray-400"
          >
            User
          </label>

          <div class="relative">
            <select
              v-model="selectedUser"
              class="w-full appearance-none rounded-xl border border-white/10 bg-[#0f0f15] px-4 py-2.5 pr-10 text-sm text-gray-300 outline-none focus:border-pink-500/50"
              @change="currentPage = 1"
            >
              <option value="All users">All users</option>
              <option
                v-for="userId in userOptions.filter((u) => u !== 'All users')"
                :key="userId"
                :value="userId"
              >
                {{ userName(userId) }}
              </option>
            </select>

            <ChevronDown
              class="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-500"
            />
          </div>
        </div>

        <!-- Date -->
        <div>
          <label
            class="mb-2 block text-xs font-medium text-gray-400"
          >
            Date
          </label>

          <div class="relative">
            <CalendarDays
              class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-500"
            />

            <input
              v-model="selectedDate"
              type="date"
              class="w-full rounded-xl border border-white/10 bg-[#0f0f15] py-2.5 pl-10 pr-4 text-sm text-gray-400 outline-none focus:border-pink-500/50"
              @change="currentPage = 1"
            />
          </div>
        </div>

        <!-- Reset -->
        <div class="flex items-end">
          <button
            type="button"
            class="w-full rounded-xl border border-white/10 bg-[#0f0f15] px-5 py-2.5 text-sm font-medium text-gray-300 transition hover:bg-white/5 hover:text-white xl:w-auto"
            @click="resetFilters"
          >
            Reset
          </button>
        </div>

      </div>
    </div>

    <!-- Logs Table -->
    <div
      class="overflow-hidden rounded-2xl border border-white/10 bg-[#15151d]"
    >

      <!-- Desktop Header -->
      <div
        class="hidden min-w-[900px] grid-cols-[1.25fr_0.9fr_0.8fr_0.8fr_2.2fr] border-b border-white/10 bg-white/[0.02] px-5 py-4 text-[11px] font-semibold uppercase tracking-wide text-gray-500 lg:grid"
      >
        <div>Time</div>
        <div>User</div>
        <div>Action</div>
        <div>Resource</div>
        <div>Details</div>
      </div>

      <!-- Loading -->
      <div
        v-if="loading"
        class="flex flex-col items-center justify-center py-16"
      >
        <p class="text-sm text-gray-400">
          Loading logs...
        </p>
      </div>

      <!-- Rows -->
      <div
        v-for="log in paginatedLogs"
        v-else
        :key="log.audit_log_id"
        class="border-b border-white/5 px-5 py-4 last:border-b-0 hover:bg-white/[0.025]"
      >

        <!-- Desktop -->
        <div
          class="hidden min-w-[900px] grid-cols-[1.25fr_0.9fr_0.8fr_0.8fr_2.2fr] items-center text-xs lg:grid"
        >

          <!-- Time -->
          <div class="text-gray-300">
            {{ formatTime(log.created_at) }}
          </div>

          <!-- User -->
          <div class="text-gray-300">
            {{ userName(log.user_id) }}
          </div>

          <!-- Action -->
          <div>
            <span
              class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-[10px] font-semibold"
              :class="getActionClass(log.action)"
            >
              <component
                :is="getActionIcon(log.action)"
                class="h-3 w-3"
              />

              {{ log.action }}
            </span>
          </div>

          <!-- Resource -->
          <div class="text-gray-300">
            {{ log.resource }}
          </div>

          <!-- Details -->
          <div class="truncate pr-4 text-gray-300">
            {{ log.details }}
          </div>

        </div>

        <!-- Mobile -->
        <div class="space-y-3 lg:hidden">

          <div class="flex items-center justify-between">
            <span class="text-xs text-gray-500">
              {{ formatTime(log.created_at) }}
            </span>

            <span
              class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-[10px] font-semibold"
              :class="getActionClass(log.action)"
            >
              <component
                :is="getActionIcon(log.action)"
                class="h-3 w-3"
              />

              {{ log.action }}
            </span>
          </div>

          <div class="text-sm font-medium text-white">
            {{ log.details }}
          </div>

          <div
            class="flex flex-wrap gap-x-5 gap-y-1 text-xs text-gray-500"
          >
            <span>
              User: {{ userName(log.user_id) }}
            </span>

            <span>
              Resource: {{ log.resource }}
            </span>
          </div>

        </div>

      </div>

      <!-- Empty -->
      <div
        v-if="!loading && paginatedLogs.length === 0"
        class="flex flex-col items-center justify-center py-16"
      >
        <Search class="mb-3 h-8 w-8 text-gray-600" />

        <p class="text-sm font-medium text-gray-400">
          No logs found
        </p>

        <p class="mt-1 text-xs text-gray-600">
          ลองเปลี่ยนเงื่อนไขการค้นหา
        </p>
      </div>

      <!-- Footer / Pagination -->
      <div
        class="flex flex-col gap-4 border-t border-white/10 px-5 py-4 sm:flex-row sm:items-center sm:justify-between"
      >

        <div class="text-xs text-gray-500">
          Showing
          <span class="text-gray-300">
            {{ showingFrom }}
          </span>
          to
          <span class="text-gray-300">
            {{ showingTo }}
          </span>
          of
          <span class="text-gray-300">
            {{ filteredLogs.length }}
          </span>
          logs
        </div>

        <div class="flex items-center gap-1">

          <!-- Previous -->
          <button
            type="button"
            class="flex h-8 w-8 items-center justify-center rounded-lg border border-white/10 text-gray-500 transition hover:bg-white/5 hover:text-white disabled:cursor-not-allowed disabled:opacity-40"
            :disabled="currentPage === 1"
            @click="changePage(currentPage - 1)"
          >
            <ChevronLeft class="h-4 w-4" />
          </button>

          <!-- Pages -->
          <button
            v-for="page in totalPages"
            :key="page"
            type="button"
            class="flex h-8 min-w-8 items-center justify-center rounded-lg border px-2 text-xs font-medium transition"
            :class="
              currentPage === page
                ? 'border-pink-500 bg-pink-500/20 text-pink-400'
                : 'border-white/10 text-gray-400 hover:bg-white/5 hover:text-white'
            "
            @click="changePage(page)"
          >
            {{ page }}
          </button>

          <!-- Next -->
          <button
            type="button"
            class="flex h-8 w-8 items-center justify-center rounded-lg border border-white/10 text-gray-500 transition hover:bg-white/5 hover:text-white disabled:cursor-not-allowed disabled:opacity-40"
            :disabled="currentPage === totalPages"
            @click="changePage(currentPage + 1)"
          >
            <ChevronRight class="h-4 w-4" />
          </button>

        </div>

      </div>

    </div>

  </div>
</template>