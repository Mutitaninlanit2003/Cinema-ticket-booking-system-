import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  // =========================
  // User
  // =========================
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/user/HomeView.vue'),
  },
  {
    path: '/movies',
    name: 'movies',
    component: () => import('@/views/user/MoviesView.vue'),
  },
  {
    path: '/help',
    name: 'help',
    component: () => import('@/views/user/HelpCenterView.vue'),
  },
  {
    path: '/refund',
    name: 'refund',
    component: () => import('@/views/user/RefundPolicyView.vue'),
  },
  {
    path: '/contact',
    name: 'contact',
    component: () => import('@/views/user/ContactUsView.vue'),
  },
  {
    path: '/movies/:movie_id',
    name: 'movie-detail',
    component: () => import('@/views/user/MovieDetailView.vue'),
  },
  {
    path: '/movies/:movie_id/theaters',
    name: 'ShowtimeSelection',
    component: () =>
      import('@/views/user/ShowtimeSelectionView.vue'),
  },
  {
    path: '/movies/:movie_id/showtimes',
    name: 'showtime-selection',
    component: () =>
      import('@/views/user/ShowtimeSelectionView.vue'),
  },
  {
    path: '/seats/:showtimeId',
    name: 'seat-selection',
    component: () =>
      import('@/views/user/SeatSelectionView.vue'),
  },
  {
    path: '/booking/summary',
    name: 'booking-summary',
    component: () =>
      import('@/views/user/BookingSummaryView.vue'),
  },
  {
    path: '/payment/:bookingId',
    name: 'payment',
    component: () =>
      import('@/views/user/PaymentView.vue'),
  },
  {
    path: '/booking/success/:bookingId',
    name: 'booking-success',
    component: () =>
      import('@/views/user/BookingSuccessView.vue'),
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/user/LoginView.vue'),
  },
  {
    path: '/oauth/callback',
    name: 'oauth-callback',
    component: () =>
      import('@/views/user/OAuthCallbackView.vue'),
  },
  {
    path: '/register',
    name: 'register',
    component: () =>
      import('@/views/user/RegisterView.vue'),
  },
  {
    path: '/forgot-password',
    name: 'forgot-password',
    component: () =>
      import('@/views/user/ForgotPasswordView.vue'),
  },
  {
    path: '/reset-password',
    component: () =>
      import('@/views/user/ResetPasswordView.vue'),
  },
  {
    path: '/profile',
    name: 'profile',
    component: () =>
      import('@/views/user/ProfileView.vue'),
    meta: {
      requiresAuth: true,
      role: 'user',
    },
  },
  {
    path: '/bookings',
    name: 'my-bookings',
    component: () =>
      import('@/views/user/MyBookingsView.vue'),
    meta: {
      requiresAuth: true,
      role: 'user',
    },
  },
  {
    path: '/bookings/:id',
    name: 'booking-detail',
    component: () =>
      import('@/views/user/BookingDetailView.vue'),
    meta: {
      requiresAuth: true,
      role: 'user',
    },
    
  },

  // =========================
  // Admin
  // =========================
  {
    path: '/admin/login',
    name: 'admin-login',
    component: () =>
      import('@/views/admin/AdminLoginView.vue'),
  },
  {
    path: '/admin',
    name: 'admin-dashboard',
    component: () =>
      import('@/views/admin/DashboardView.vue'),
    meta: {
      requiresAuth: true,
      role: 'admin',
    },
  },
  {
    path: '/admin/profile',
    name: 'admin-profile',
    component: () =>
      import('@/views/admin/Profileview.vue'),
    meta: {
      requiresAuth: true,
      role: 'admin',
    },
  },
  {
    path: '/admin/movies',
    name: 'admin-movies',
    component: () =>
      import('@/views/admin/MoviesView.vue'),
    meta: {
      requiresAuth: true,
      role: 'admin',
    },
  },
  {
    path: '/admin/theaters',
    name: 'admin-theaters',
    component: () =>
      import('@/views/admin/TheatersView.vue'),
    meta: {
      requiresAuth: true,
      role: 'admin',
    },
  },
  {
    path: '/admin/screens',
    name: 'admin-screens',
    component: () =>
      import('@/views/admin/ScreensView.vue'),
    meta: {
      requiresAuth: true,
      role: 'admin',
    },
  },
  {
    path: '/admin/seats',
    name: 'admin-seats',
    component: () =>
      import('@/views/admin/SeatsView.vue'),
    meta: {
      requiresAuth: true,
      role: 'admin',
    },
  },
  {
    path: '/admin/showtimes',
    name: 'admin-showtimes',
    component: () =>
      import('@/views/admin/ShowtimesView.vue'),
    meta: {
      requiresAuth: true,
      role: 'admin',
    },
  },
  {
    path: '/admin/bookings',
    name: 'admin-bookings',
    component: () =>
      import('@/views/admin/BookingsView.vue'),
    meta: {
      requiresAuth: true,
      role: 'admin',
    },
  },
  {
    path: '/admin/payments',
    name: 'admin-payments',
    component: () =>
      import('@/views/admin/PaymentsView.vue'),
    meta: {
      requiresAuth: true,
      role: 'admin',
    },
  },
  {
    path: '/admin/audit-logs',
    name: 'admin-audit-logs',
    component: () =>
      import('@/views/admin/AuditLogsView.vue'),
    meta: {
      requiresAuth: true,
      role: 'admin',
    },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,

  scrollBehavior() {
    return { top: 0 }
  },
})

router.beforeEach((to) => {
  const authStore = useAuthStore()

  /*
   * Load both sessions:
   *
   * user_auth  -> authStore.user
   * admin_auth -> authStore.admin
   */
  authStore.initFromStorage()

  const requiresAuth =
    to.meta.requiresAuth === true

  const requiredRole =
    to.meta.role as string | undefined

  /*
   * Public pages
   */
  if (!requiresAuth) {
    return true
  }

  /*
   * =========================
   * ADMIN ROUTES
   * =========================
   */
  if (requiredRole === 'admin') {
    if (!authStore.isAdminAuthenticated) {
      return {
        path: '/admin/login',
        query: {
          redirect: to.fullPath,
        },
      }
    }

    return true
  }

  /*
   * =========================
   * USER ROUTES
   * =========================
   */
  if (requiredRole === 'user') {
    if (!authStore.isUserAuthenticated) {
      return {
        path: '/login',
        query: {
          redirect: to.fullPath,
        },
      }
    }

    return true
  }

  /*
   * Fallback for protected routes
   * that do not specify a role.
   */
  if (!authStore.isAuthenticated) {
    return {
      path: '/login',
      query: {
        redirect: to.fullPath,
      },
    }
  }

  return true
})

export default router