<script setup lang="ts">
const sections = [
  {
    title: '1. Cancelling before payment',
    body: 'A booking that has not been paid for yet ("pending") can be cancelled at any time from "My Bookings" at no cost.',
  },
  {
    title: '2. Cancelling after payment',
    body: 'Confirmed bookings can be cancelled up to 2 hours before the showtime for a full refund to your original payment method. Cancellations made after that window are not eligible for a refund.',
  },
  {
    title: '3. Refund timing',
    body: 'Approved refunds are processed within 5–10 business days, depending on your bank or card issuer.',
  },
  {
    title: '4. Cancelled or rescheduled showtimes',
    body: 'If a theater cancels or reschedules a showtime, all affected bookings are automatically eligible for a full refund, regardless of timing.',
  },
  {
    title: '5. Non-refundable situations',
    body: 'Tickets already used (scanned at entry), and no-shows, are not eligible for a refund.',
  },
]
</script>

<template>
  <div class="min-h-screen bg-[#0d0d12] px-6 py-14 text-white">
    <div class="mx-auto max-w-3xl">

      <h1 class="text-3xl font-black uppercase tracking-tight">
        Refund Policy
      </h1>

      <p class="mt-3 text-sm text-gray-400">
        Please review the terms below before requesting a cancellation or refund.
      </p>

      <div class="mt-10 space-y-6">
        <div
          v-for="section in sections"
          :key="section.title"
          class="rounded-xl border border-white/10 bg-[#15151d] p-5"
        >
          <h2 class="text-sm font-bold text-white">
            {{ section.title }}
          </h2>

          <p class="mt-2 text-sm leading-6 text-gray-400">
            {{ section.body }}
          </p>
        </div>
      </div>

      <p class="mt-8 text-sm text-gray-500">
        Have a question about a specific booking?
        <router-link to="/contact" class="text-pink-500 hover:text-pink-400">
          Contact us
        </router-link>
        and we'll help you out.
      </p>

    </div>
  </div>
</template>