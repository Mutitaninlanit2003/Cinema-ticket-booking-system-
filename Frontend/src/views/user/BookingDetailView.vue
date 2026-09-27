<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()
</script>

<template>
  <div class="mx-auto max-w-7xl px-6 py-10">
    <h1 class="text-3xl font-bold">
      Booking Detail
    </h1>

    <p class="mt-4">
      Booking ID: {{ route.params.id }}
    </p>
  </div>
</template>