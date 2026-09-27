<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()
</script>

<template>
  <div class="mx-auto max-w-7xl px-6 py-10 text-center">
    <h1 class="text-3xl font-bold">
      Booking Successful
    </h1>

    <p class="mt-4 text-gray-600">
      การจองของคุณสำเร็จแล้ว
    </p>

    <p class="mt-2">
      Booking ID: {{ route.params.bookingId }}
    </p>

    <router-link
      to="/bookings"
      class="mt-6 inline-block rounded-lg bg-black px-6 py-3 text-white"
    >
      ดูรายการจอง
    </router-link>
  </div>
</template>