<script setup lang="ts">
import { Instagram, Twitter } from 'lucide-vue-next'

const year = new Date().getFullYear()
</script>

<template>
  <footer class="border-t border-white/10 bg-[#09090d]">
    <div class="mx-auto max-w-7xl px-8 py-14">

      <!-- Main Footer -->
      <div class="grid grid-cols-1 gap-12 md:grid-cols-4">

        <!-- Brand -->
        <div>
          <router-link
            to="/"
            class="inline-flex items-center gap-3"
          >
            <div
              class="flex h-8 w-8 items-center justify-center rounded-lg
                     bg-gradient-to-br from-pink-500 to-purple-500"
            >
              <span class="text-sm font-bold text-white">▦</span>
            </div>

            <span class="text-2xl font-bold tracking-tight text-white">
              CINEMA<span class="text-pink-500">.</span>
            </span>
          </router-link>

          <p class="mt-6 max-w-xs text-sm leading-7 text-gray-400">
            Book your seats in seconds. Premium screens,
            real-time seat maps, and mobile tickets —
            all in one place.
          </p>
        </div>

        <!-- Explore -->
        <div>
          <h3
            class="text-sm font-semibold uppercase tracking-[0.18em] text-gray-400"
          >
            Explore
          </h3>

          <div class="mt-6 flex flex-col gap-4">
            <router-link
              to="/movies"
              class="text-sm text-gray-300 transition hover:text-white"
            >
              Now showing
            </router-link>

            <router-link
              to="/movies"
              class="text-sm text-gray-300 transition hover:text-white"
            >
              Coming soon
            </router-link>

            <router-link
              to="/cinemas"
              class="text-sm text-gray-300 transition hover:text-white"
            >
              Cinemas
            </router-link>

            <router-link
              to="/promotions"
              class="text-sm text-gray-300 transition hover:text-white"
            >
              Promotions
            </router-link>
          </div>
        </div>

        <!-- Support -->
        <div>
          <h3
            class="text-sm font-semibold uppercase tracking-[0.18em] text-gray-400"
          >
            Support
          </h3>

          <div class="mt-6 flex flex-col gap-4">
            <router-link
              to="/help"
              class="text-sm text-gray-300 transition hover:text-white"
            >
              Help center
            </router-link>

            <router-link
              to="/refund"
              class="text-sm text-gray-300 transition hover:text-white"
            >
              Refund policy
            </router-link>

            <router-link
              to="/contact"
              class="text-sm text-gray-300 transition hover:text-white"
            >
              Contact us
            </router-link>
          </div>
        </div>

        <!-- Follow -->
        <div>
          <h3
            class="text-sm font-semibold uppercase tracking-[0.18em] text-gray-400"
          >
            Follow
          </h3>

          <div class="mt-6 flex gap-3">

            <a
              href="#"
              aria-label="Instagram"
              class="flex h-10 w-10 items-center justify-center
                     rounded-full border border-white/15
                     text-gray-300 transition
                     hover:border-white/30 hover:text-white"
            >
              <Instagram :size="18" :stroke-width="1.6" />
            </a>

            <a
              href="#"
              aria-label="Twitter"
              class="flex h-10 w-10 items-center justify-center
                     rounded-full border border-white/15
                     text-gray-300 transition
                     hover:border-white/30 hover:text-white"
            >
              <Twitter :size="18" :stroke-width="1.6" />
            </a>

          </div>
        </div>

      </div>

      <!-- Divider -->
      <div class="my-12 h-px bg-white/10"></div>

      <!-- Bottom -->
      <div
        class="flex flex-col gap-4 text-sm text-gray-500
               md:flex-row md:items-center md:justify-between"
      >
        <p>
          © {{ year }} Cinema Booking. All rights reserved.
        </p>

        <p>
          Design system v1.0 —
          <span class="text-gray-400">
            dark / marquee red
          </span>
        </p>
      </div>

    </div>
  </footer>
</template>

<style scoped>
footer {
  width: 100%;
  max-width: 100%;
  overflow: hidden;
}
</style>