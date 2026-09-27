<script setup lang="ts">
import { ref } from 'vue'

import AdminSidebar from '@/components/admin/AdminSidebar.vue'
import AdminTopbar from '@/components/admin/AdminTopbar.vue'

const sidebarOpen = ref(false)
</script>

<template>
  <div class="min-h-screen bg-[#0d0d12] text-white">
    <AdminSidebar
      :open="sidebarOpen"
      @close="sidebarOpen = false"
    />

    <div class="flex min-h-screen flex-col lg:pl-[260px]">
      <AdminTopbar
        @toggle-sidebar="sidebarOpen = !sidebarOpen"
      />

      <main class="flex-1">
        <slot />
      </main>
    </div>
  </div>
</template>