<script setup lang="ts">
import { ref } from 'vue'

const name = ref('')
const email = ref('')
const message = ref('')
const submitted = ref(false)

function handleSubmit() {
  // NOTE: There is currently no backend endpoint for contact messages.
  // This just shows a confirmation for now — wire this up to a real
  // API call (e.g. POST /contact) once one exists.
  submitted.value = true
}
</script>

<template>
  <div class="min-h-screen bg-[#0d0d12] px-6 py-14 text-white">
    <div class="mx-auto max-w-2xl">

      <h1 class="text-3xl font-black uppercase tracking-tight">
        Contact Us
      </h1>

      <p class="mt-3 text-sm text-gray-400">
        Have a question or an issue with a booking? Send us a message
        and we'll get back to you.
      </p>

      <!-- Success -->
      <div
        v-if="submitted"
        class="mt-10 rounded-xl border border-green-500/20 bg-green-500/10 p-6 text-center"
      >
        <p class="text-sm font-semibold text-green-400">
          Thanks — your message has been sent.
        </p>

        <p class="mt-2 text-sm text-gray-400">
          We'll get back to you as soon as we can.
        </p>
      </div>

      <!-- Form -->
      <form
        v-else
        class="mt-10 space-y-5 rounded-xl border border-white/10 bg-[#15151d] p-6"
        @submit.prevent="handleSubmit"
      >

        <div>
          <label class="mb-2 block text-sm font-medium text-gray-300">
            Name
          </label>

          <input
            v-model="name"
            type="text"
            required
            placeholder="Your name"
            class="w-full rounded-lg border border-white/10 bg-[#0d0d12] px-4 py-3 text-sm text-white outline-none placeholder:text-gray-600 focus:border-pink-500/60"
          />
        </div>

        <div>
          <label class="mb-2 block text-sm font-medium text-gray-300">
            Email
          </label>

          <input
            v-model="email"
            type="email"
            required
            placeholder="you@example.com"
            class="w-full rounded-lg border border-white/10 bg-[#0d0d12] px-4 py-3 text-sm text-white outline-none placeholder:text-gray-600 focus:border-pink-500/60"
          />
        </div>

        <div>
          <label class="mb-2 block text-sm font-medium text-gray-300">
            Message
          </label>

          <textarea
            v-model="message"
            required
            rows="5"
            placeholder="How can we help?"
            class="w-full resize-none rounded-lg border border-white/10 bg-[#0d0d12] px-4 py-3 text-sm text-white outline-none placeholder:text-gray-600 focus:border-pink-500/60"
          ></textarea>
        </div>

        <button
          type="submit"
          class="w-full rounded-lg bg-pink-500 py-3 text-sm font-semibold text-white transition hover:bg-pink-600"
        >
          Send message
        </button>

      </form>

      <div class="mt-8 grid gap-4 sm:grid-cols-2">
        <div class="rounded-xl border border-white/10 bg-[#15151d] p-5">
          <p class="text-xs font-semibold uppercase tracking-wide text-gray-500">
            Email
          </p>

          <p class="mt-1 text-sm text-white">
            support@cinema.example
          </p>
        </div>

        <div class="rounded-xl border border-white/10 bg-[#15151d] p-5">
          <p class="text-xs font-semibold uppercase tracking-wide text-gray-500">
            Hours
          </p>

          <p class="mt-1 text-sm text-white">
            Daily, 9:00 – 21:00
          </p>
        </div>
      </div>

    </div>
  </div>
</template>