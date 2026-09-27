<script setup lang="ts">
import {
  LayoutDashboard,
  Film,
  Building2,
  MonitorPlay,
  Armchair,
  CalendarClock,
  Ticket,
  CreditCard,
  ScrollText,
  Clapperboard,
  X,
} from '@lucide/vue'

defineProps<{
  open: boolean
}>()

defineEmits<{
  close: []
}>()

const navItems = [
  { label: 'Dashboard', to: '/admin', icon: LayoutDashboard },
  { label: 'Movies', to: '/admin/movies', icon: Film },
  { label: 'Theaters', to: '/admin/theaters', icon: Building2 },
  { label: 'Screens', to: '/admin/screens', icon: MonitorPlay },
  { label: 'Seats', to: '/admin/seats', icon: Armchair },
  { label: 'Showtimes', to: '/admin/showtimes', icon: CalendarClock },
  { label: 'Bookings', to: '/admin/bookings', icon: Ticket },
  { label: 'Payments', to: '/admin/payments', icon: CreditCard },
  { label: 'Audit Logs', to: '/admin/audit-logs', icon: ScrollText },
]
</script>

<template>
  <!-- Mobile backdrop -->
  <div
    v-if="open"
    class="fixed inset-0 z-40 bg-black/70 lg:hidden"
    @click="$emit('close')"
  ></div>

  <!-- Sidebar -->
  <aside
    class="fixed inset-y-0 left-0 z-50 flex w-[260px] flex-col border-r border-white/10 bg-[#15151d] transition-transform duration-300 lg:translate-x-0"
    :class="open ? 'translate-x-0' : '-translate-x-full'"
  >
    <!-- Logo -->
    <div
      class="flex h-[80px] items-center justify-between border-b border-white/10 px-6"
    >
      <RouterLink
        to="/admin"
        class="flex items-center gap-3 text-white"
      >
        <div
          class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-pink-500 to-purple-500"
        >
          <Clapperboard class="h-4 w-4 text-white" />
        </div>

        <span class="text-lg font-black tracking-tight">
          ADMIN
        </span>
      </RouterLink>

      <!-- Mobile close -->
      <button
        type="button"
        class="text-gray-400 hover:text-white lg:hidden"
        aria-label="ปิดเมนู"
        @click="$emit('close')"
      >
        <X class="h-5 w-5" />
      </button>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 space-y-1 px-4 py-5">
      <RouterLink
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        class="group flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-medium text-gray-300 transition-all duration-200 hover:bg-white/5 hover:text-white"
        active-class="bg-[#3b1827] text-white"
        @click="$emit('close')"
      >
        <component
          :is="item.icon"
          class="h-[18px] w-[18px] shrink-0 text-gray-400 transition-colors group-hover:text-white"
        />

        <span>
          {{ item.label }}
        </span>
      </RouterLink>
    </nav>
  </aside>
</template>