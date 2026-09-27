<script setup lang="ts">
import { ref } from 'vue'

const openIndex = ref<number | null>(0)

const faqs = [
  {
    question: 'How do I book a ticket?',
    answer:
      'Browse "Now Showing" or "Coming Soon", pick a movie, choose a showtime and theater, select your seats, then complete payment. Your ticket will appear under "My Bookings".',
  },
  {
    question: 'Can I change or cancel my booking?',
    answer:
      'Pending bookings can be cancelled from the "My Bookings" page before payment is completed. Once a booking is confirmed and paid, please see our Refund Policy for cancellation options.',
  },
  {
    question: 'How do I get my ticket?',
    answer:
      'Once payment is confirmed, your e-ticket is available under "My Bookings" → "View ticket". You can show this on your phone at the theater entrance.',
  },
  {
    question: 'I didn\'t receive a confirmation. What should I do?',
    answer:
      'Check "My Bookings" first — your booking may already be confirmed there even if an email didn\'t arrive. If it still shows as pending, complete payment from that page, or contact us below for help.',
  },
  {
    question: 'Can I book seats for a group?',
    answer:
      'Yes — on the seat selection screen you can select multiple seats in a single booking, as long as they\'re available.',
  },
]

function toggle(index: number) {
  openIndex.value = openIndex.value === index ? null : index
}
</script>

<template>
  <div class="min-h-screen bg-[#0d0d12] px-6 py-14 text-white">
    <div class="mx-auto max-w-3xl">

      <h1 class="text-3xl font-black uppercase tracking-tight">
        Help Center
      </h1>

      <p class="mt-3 text-sm text-gray-400">
        Answers to common questions about booking, tickets, and payments.
      </p>

      <div class="mt-10 divide-y divide-white/10 rounded-xl border border-white/10 bg-[#15151d]">

        <div
          v-for="(faq, index) in faqs"
          :key="faq.question"
        >
          <button
            type="button"
            class="flex w-full items-center justify-between gap-4 px-5 py-4 text-left"
            @click="toggle(index)"
          >
            <span class="text-sm font-semibold text-white">
              {{ faq.question }}
            </span>

            <span
              class="shrink-0 text-lg text-gray-500 transition"
              :class="{ 'rotate-45': openIndex === index }"
            >
              +
            </span>
          </button>

          <p
            v-if="openIndex === index"
            class="px-5 pb-5 text-sm leading-6 text-gray-400"
          >
            {{ faq.answer }}
          </p>
        </div>

      </div>

      <p class="mt-8 text-sm text-gray-500">
        Still need help?
        <router-link to="/contact" class="text-pink-500 hover:text-pink-400">
          Contact us
        </router-link>
        directly.
      </p>

    </div>
  </div>
</template>