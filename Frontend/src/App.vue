<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

import UserLayout from '@/layouts/UserLayout.vue'
import AdminLayout from '@/layouts/AdminLayout.vue'

// Current route
const route = useRoute()

// --------------------------------------------------
// Blank Layout
// --------------------------------------------------
// ใช้สำหรับหน้า Authentication
// เช่น Login / Register / Forgot Password
//
// หน้าเหล่านี้จะไม่มี Header / Footer / Sidebar
// เพราะผู้ใช้ยังไม่ได้เข้าสู่ระบบ
const BlankLayout = {
  setup(
    _: unknown,
    {
      slots,
    }: {
      slots: Record<string, () => unknown>
    },
  ) {
    return () => (
      slots.default ? slots.default() : null
    )
  },
}

// --------------------------------------------------
// Layout selection
// --------------------------------------------------
const layout = computed(() => {
  // ------------------------------------------------
  // User Authentication Pages
  // ------------------------------------------------
  // ไม่ใช้ UserLayout
  // ไม่แสดง Header / Footer
  if (
    route.path === '/login' ||
    route.path === '/register' ||
    route.path === '/forgot-password' ||
    route.path === '/reset-password'
  ) {
    return BlankLayout
  }

  // ------------------------------------------------
  // Admin Login
  // ------------------------------------------------
  // ไม่ใช้ AdminLayout เพราะยังไม่ได้ Login
  if (route.path === '/admin/login') {
    return BlankLayout
  }

  // ------------------------------------------------
  // Admin Pages
  // ------------------------------------------------
  if (route.path.startsWith('/admin')) {
    return AdminLayout
  }

  // ------------------------------------------------
  // User Pages
  // ------------------------------------------------
  return UserLayout
})
</script>

<template>
  <component :is="layout">
    <router-view />
  </component>
</template>