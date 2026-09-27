<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import api from '@/services/api'
import {
  Search,
  CreditCard,
  Trash2,
} from 'lucide-vue-next'

// -----------------------------------------------------------------------
// Types — match backend model exactly (models/payment.go)
// -----------------------------------------------------------------------
type PaymentStatus = 'pending' | 'paid' | 'failed' | 'refunded'
type PaymentMethod = 'cash' | 'credit_card' | 'promptpay' | 'transfer'

interface Payment {
  id: number
  payment_id: string
  booking_id: string
  amount: number
  method: PaymentMethod
  status: PaymentStatus
  transaction_ref: string
  created_at: string
}

const searchQuery = ref('')
const activeFilter = ref<'all' | PaymentStatus>('all')

const payments = ref<Payment[]>([])

const loading = ref(false)
const deletingId = ref<string | null>(null)
const errorMessage = ref('')

async function loadPayments() {
  loading.value = true
  errorMessage.value = ''

  try {
    const { data } = await api.get<Payment[]>('/payments')
    payments.value = data
  } catch (err) {
    errorMessage.value = 'Failed to load payments. Please try again.'
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(loadPayments)

async function deletePayment(payment: Payment) {
  const confirmed = window.confirm(
    `Delete payment ${payment.payment_id} (${formatAmount(payment.amount)})? This cannot be undone.`,
  )

  if (!confirmed) return

  deletingId.value = payment.payment_id
  errorMessage.value = ''

  try {
    await api.delete(`/payments/${payment.payment_id}`)
    // Only remove from local state after the backend confirms deletion.
    payments.value = payments.value.filter((p) => p.payment_id !== payment.payment_id)
  } catch (err) {
    errorMessage.value = 'Failed to delete payment. Please try again.'
    console.error(err)
  } finally {
    deletingId.value = null
  }
}

const filteredPayments = computed(() => {
  return payments.value.filter((payment) => {
    const matchesFilter =
      activeFilter.value === 'all' ||
      payment.status === activeFilter.value

    const query = searchQuery.value.toLowerCase().trim()

    const matchesSearch =
      !query ||
      payment.payment_id.toLowerCase().includes(query) ||
      payment.booking_id.toLowerCase().includes(query) ||
      payment.transaction_ref.toLowerCase().includes(query) ||
      payment.method.toLowerCase().includes(query)

    return matchesFilter && matchesSearch
  })
})

function getStatusClass(status: PaymentStatus) {
  switch (status) {
    case 'paid':
      return 'bg-emerald-500/15 text-emerald-400'

    case 'pending':
      return 'bg-amber-500/15 text-amber-400'

    case 'failed':
      return 'bg-red-500/15 text-red-400'

    case 'refunded':
      return 'bg-purple-500/15 text-purple-400'

    default:
      return 'bg-gray-500/15 text-gray-400'
  }
}

function getStatusLabel(status: PaymentStatus) {
  return status
}

// Backend stores method as a fixed snake_case enum
// (cash | credit_card | promptpay | transfer) — map to display labels.
function getMethodLabel(method: PaymentMethod) {
  const labels: Record<PaymentMethod, string> = {
    cash: 'Cash',
    credit_card: 'Credit Card',
    promptpay: 'PromptPay',
    transfer: 'Bank Transfer',
  }

  return labels[method] ?? method
}

function formatAmount(amount: number) {
  return `฿${amount.toLocaleString('en-US', {
    minimumFractionDigits: 2,
  })}`
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
</script>

<template>
  <div class="mx-auto max-w-7xl px-6 py-10">

    <!-- Header -->
    <div class="mb-8">
      <div class="flex items-center gap-4">

        <!-- Payment Icon -->
        <div
          class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-gradient-to-br from-purple-500 to-pink-500 shadow-lg shadow-purple-500/20"
        >
          <CreditCard class="size-6" />
        </div>

        <div>
          <h1 class="text-2xl font-bold text-white">
            Manage Payments
          </h1>

          <p class="mt-1 text-sm text-gray-400">
            Manage payment transactions
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
    <div class="mb-6 flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">

      <!-- Search -->
      <div class="relative w-full lg:w-[340px]">
        <Search
          class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-500"
        />

        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search by payment ID or booking..."
          class="w-full rounded-xl border border-white/10 bg-[#15151d] py-3 pl-10 pr-4 text-sm text-white outline-none placeholder:text-gray-500 transition focus:border-pink-500/50"
        />
      </div>

      <!-- Filters -->
      <div class="flex flex-wrap gap-2">

        <button
          type="button"
          class="rounded-xl border px-5 py-2.5 text-sm font-medium transition"
          :class="
            activeFilter === 'all'
              ? 'border-pink-500 bg-pink-500 text-white'
              : 'border-white/10 bg-[#15151d] text-gray-400 hover:bg-white/5 hover:text-white'
          "
          @click="activeFilter = 'all'"
        >
          All
        </button>

        <button
          type="button"
          class="rounded-xl border px-5 py-2.5 text-sm font-medium transition"
          :class="
            activeFilter === 'pending'
              ? 'border-amber-500 bg-amber-500/10 text-amber-400'
              : 'border-white/10 bg-[#15151d] text-gray-400 hover:bg-white/5 hover:text-white'
          "
          @click="activeFilter = 'pending'"
        >
          Pending
        </button>

        <button
          type="button"
          class="rounded-xl border px-5 py-2.5 text-sm font-medium transition"
          :class="
            activeFilter === 'paid'
              ? 'border-emerald-500 bg-emerald-500/10 text-emerald-400'
              : 'border-white/10 bg-[#15151d] text-gray-400 hover:bg-white/5 hover:text-white'
          "
          @click="activeFilter = 'paid'"
        >
          Paid
        </button>

        <button
          type="button"
          class="rounded-xl border px-5 py-2.5 text-sm font-medium transition"
          :class="
            activeFilter === 'failed'
              ? 'border-red-500 bg-red-500/10 text-red-400'
              : 'border-white/10 bg-[#15151d] text-gray-400 hover:bg-white/5 hover:text-white'
          "
          @click="activeFilter = 'failed'"
        >
          Failed
        </button>

        <button
          type="button"
          class="rounded-xl border px-5 py-2.5 text-sm font-medium transition"
          :class="
            activeFilter === 'refunded'
              ? 'border-purple-500 bg-purple-500/10 text-purple-400'
              : 'border-white/10 bg-[#15151d] text-gray-400 hover:bg-white/5 hover:text-white'
          "
          @click="activeFilter = 'refunded'"
        >
          Refunded
        </button>

      </div>
    </div>

    <!-- Payment Table -->
    <div
      class="overflow-hidden rounded-2xl border border-white/10 bg-[#15151d]"
    >

      <!-- Table Header -->
      <div
        class="grid min-w-[1080px] grid-cols-[1.1fr_1.1fr_0.8fr_1fr_1.1fr_0.9fr_1fr_0.6fr] border-b border-white/10 bg-white/[0.02] px-5 py-4 text-xs font-semibold uppercase tracking-wide text-gray-500"
      >
        <div>Payment ID</div>
        <div>Booking ID</div>
        <div>Amount</div>
        <div>Method</div>
        <div>Transaction Ref</div>
        <div>Status</div>
        <div>Created At</div>
        <div class="text-right">Actions</div>
      </div>

      <!-- Loading -->
      <div
        v-if="loading"
        class="px-5 py-16 text-center text-sm text-gray-400"
      >
        Loading payments...
      </div>

      <!-- Rows -->
      <div
        v-for="payment in filteredPayments"
        v-else
        :key="payment.payment_id"
        class="grid min-w-[1080px] grid-cols-[1.1fr_1.1fr_0.8fr_1fr_1.1fr_0.9fr_1fr_0.6fr] items-center border-b border-white/5 px-5 py-5 text-sm transition last:border-b-0 hover:bg-white/[0.025]"
      >

        <!-- Payment ID -->
        <div class="flex items-center gap-3">
          <div
            class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-pink-500/10"
          >
            <CreditCard class="h-4 w-4 text-pink-400" />
          </div>

          <span
            class="max-w-[120px] truncate font-medium text-white"
            :title="payment.payment_id"
          >
            {{ payment.payment_id }}
          </span>
        </div>

        <!-- Booking ID -->
        <div class="truncate text-gray-300" :title="payment.booking_id">
          {{ payment.booking_id }}
        </div>

        <!-- Amount -->
        <div class="font-semibold text-white">
          {{ formatAmount(payment.amount) }}
        </div>

        <!-- Method -->
        <div class="text-gray-300">
          {{ getMethodLabel(payment.method) }}
        </div>

        <!-- Transaction -->
        <div class="text-gray-400">
          {{ payment.transaction_ref || '—' }}
        </div>

        <!-- Status -->
        <div>
          <span
            class="inline-flex rounded-full px-3 py-1 text-xs font-semibold"
            :class="getStatusClass(payment.status)"
          >
            {{ getStatusLabel(payment.status) }}
          </span>
        </div>

        <!-- Created -->
        <div class="text-gray-400">
          {{ formatCreatedAt(payment.created_at) }}
        </div>

        <!-- Actions -->
        <div class="flex justify-end">
          <button
            type="button"
            class="rounded-md border border-red-500/20 p-2 text-red-400 transition hover:bg-red-500/10 disabled:cursor-not-allowed disabled:opacity-50"
            title="Delete"
            :disabled="deletingId === payment.payment_id"
            @click="deletePayment(payment)"
          >
            <Trash2 class="h-4 w-4" />
          </button>
        </div>

      </div>

      <!-- Empty -->
      <div
        v-if="!loading && filteredPayments.length === 0"
        class="flex flex-col items-center justify-center py-16 text-center"
      >
        <CreditCard class="mb-3 h-8 w-8 text-gray-600" />

        <p class="text-sm font-medium text-gray-400">
          No payments found
        </p>

        <p class="mt-1 text-xs text-gray-600">
          ลองค้นหาด้วย Payment ID หรือ Booking ID อื่น
        </p>
      </div>

    </div>

  </div>
</template>