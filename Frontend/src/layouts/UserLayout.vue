<script setup lang="ts">
import Navbar from '@/components/user/Navbar.vue'
import Footer from '@/components/user/Footer.vue'
</script>

<template>
  <div
    class="flex min-h-screen w-full max-w-full flex-col overflow-x-hidden bg-[#0b0b0f] text-white"
  >
    <Navbar />

    <main class="min-w-0 flex-1 overflow-x-hidden">
      <slot />
    </main>

    <Footer />
  </div>
</template>