import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/services/api'

const USER_STORAGE_KEY = 'user_auth'
const ADMIN_STORAGE_KEY = 'admin_auth'

export interface AuthUser {
  user_id: string
  name: string
  email: string
  role: 'user' | 'admin'
  profile_image?: string
}

interface LoginResponse {
  token: string
  user: AuthUser
}

interface StoredAuth {
  token: string
  user: AuthUser
}

export const useAuthStore = defineStore('auth', () => {
  // =========================
  // User Authentication
  // =========================

  const userToken = ref<string | null>(null)
  const user = ref<AuthUser | null>(null)

  const isUserAuthenticated = computed(() => {
    return !!userToken.value && user.value?.role === 'user'
  })

  // =========================
  // Admin Authentication
  // =========================

  const adminToken = ref<string | null>(null)
  const admin = ref<AuthUser | null>(null)

  const isAdminAuthenticated = computed(() => {
    return !!adminToken.value && admin.value?.role === 'admin'
  })

  // =========================
  // Backward Compatibility
  // =========================

  const token = computed(() => {
    return userToken.value
  })

  const currentUser = computed(() => {
    return user.value
  })

  const isAuthenticated = computed(() => {
    return isUserAuthenticated.value
  })

  const isAdmin = computed(() => {
    return isAdminAuthenticated.value
  })

  // =========================
  // Load User Auth
  // =========================
  // User ใช้ sessionStorage
  // เมื่อปิดแท็บหรือหน้าต่าง session จะหาย

  function loadUserAuth() {
    const raw = sessionStorage.getItem(USER_STORAGE_KEY)

    if (!raw) {
      userToken.value = null
      user.value = null
      return
    }

    try {
      const parsed = JSON.parse(raw) as StoredAuth

      if (
        !parsed.token ||
        !parsed.user ||
        parsed.user.role !== 'user'
      ) {
        sessionStorage.removeItem(USER_STORAGE_KEY)

        userToken.value = null
        user.value = null
        return
      }

      userToken.value = parsed.token
      user.value = parsed.user
    } catch {
      sessionStorage.removeItem(USER_STORAGE_KEY)

      userToken.value = null
      user.value = null
    }
  }

  // =========================
  // Set User Session (ใช้ตอน Google OAuth callback)
  // =========================
  // ต่างจาก loginUser() ตรงที่ไม่เรียก API /auth/login
  // เพราะ token มาจาก backend redirect กลับมาให้แล้ว

  function setUserSession(token: string, userData: AuthUser) {
    userToken.value = token
    user.value = userData

    persistUser()
  }

  // =========================
  // Load Admin Auth
  // =========================
  // Admin ยังคงใช้ localStorage
  // เพื่อแยก session ออกจาก User

  function loadAdminAuth() {
    const raw = sessionStorage.getItem(ADMIN_STORAGE_KEY)   // เปลี่ยนจาก localStorage

    if (!raw) {
      adminToken.value = null
      admin.value = null
      return
    }

    try {
      const parsed = JSON.parse(raw) as StoredAuth

      if (
        !parsed.token ||
        !parsed.user ||
        parsed.user.role !== 'admin'
      ) {
        sessionStorage.removeItem(ADMIN_STORAGE_KEY)   // เปลี่ยนจาก localStorage

        adminToken.value = null
        admin.value = null
        return
      }

      adminToken.value = parsed.token
      admin.value = parsed.user
    } catch {
      sessionStorage.removeItem(ADMIN_STORAGE_KEY)   // เปลี่ยนจาก localStorage

      adminToken.value = null
      admin.value = null
    }
  }

  // =========================
  // Initialize
  // =========================

  function initFromStorage() {
    // โหลด User จาก sessionStorage
    loadUserAuth()

    // โหลด Admin จาก sessionStorage
    loadAdminAuth()

    // ล้าง storage เก่าที่เคยใช้ร่วมกัน
    localStorage.removeItem('auth')
    localStorage.removeItem('token')
    localStorage.removeItem('user')

    // ล้าง User session แบบเก่าที่เคยเก็บใน localStorage
    localStorage.removeItem(USER_STORAGE_KEY)

    // ล้าง Admin session แบบเก่าที่เคยเก็บใน localStorage
    // (เพราะตอนนี้ย้ายไปใช้ sessionStorage แล้ว)
    localStorage.removeItem(ADMIN_STORAGE_KEY)
  }

  // =========================
  // Save User Auth
  // =========================

  function persistUser() {
    if (userToken.value && user.value) {
      sessionStorage.setItem(
        USER_STORAGE_KEY,
        JSON.stringify({
          token: userToken.value,
          user: user.value,
        }),
      )
    } else {
      sessionStorage.removeItem(USER_STORAGE_KEY)
    }
  }

  // =========================
  // Save Admin Auth
  // =========================

  function persistAdmin() {
    if (adminToken.value && admin.value) {
      sessionStorage.setItem(   // เปลี่ยนจาก localStorage
        ADMIN_STORAGE_KEY,
        JSON.stringify({
          token: adminToken.value,
          user: admin.value,
        }),
      )
    } else {
      sessionStorage.removeItem(ADMIN_STORAGE_KEY)   // เปลี่ยนจาก localStorage
    }
  }

  // =========================
  // User Login
  // =========================

  async function loginUser(
    email: string,
    password: string,
  ) {
    const response = await api.post<LoginResponse>(
      '/auth/login',
      {
        email,
        password,
      },
    )

    const loggedInUser = response.data.user

    // ป้องกันไม่ให้ Admin token ถูกเก็บเป็น User session
    if (loggedInUser.role !== 'user') {
      throw new Error(
        'This account is not a user account.',
      )
    }

    userToken.value = response.data.token
    user.value = loggedInUser

    persistUser()

    return loggedInUser
  }

  // =========================
  // Admin Login
  // =========================

  async function loginAdmin(
    email: string,
    password: string,
  ) {
    const response = await api.post<LoginResponse>(
      '/auth/login',
      {
        email,
        password,
      },
    )

    const loggedInAdmin = response.data.user

    // ป้องกันไม่ให้ User token ถูกเก็บเป็น Admin session
    if (loggedInAdmin.role !== 'admin') {
      throw new Error(
        'This account does not have admin access.',
      )
    }

    adminToken.value = response.data.token
    admin.value = loggedInAdmin

    persistAdmin()

    return loggedInAdmin
  }

  // =========================
  // User Logout
  // =========================

  function logoutUser() {
    userToken.value = null
    user.value = null

    persistUser()
  }

  // =========================
  // Admin Logout
  // =========================

  function logoutAdmin() {
    adminToken.value = null
    admin.value = null

    persistAdmin()
  }

  // =========================
  // Backward Compatibility
  // =========================

  // login() เดิมให้ถือว่าเป็น User Login
  async function login(
    email: string,
    password: string,
  ) {
    return loginUser(email, password)
  }

  // logout() เดิมให้ถือว่าเป็น User Logout
  function logout() {
    logoutUser()
  }

  return {
    // User
    userToken,
    user,
    isUserAuthenticated,
    loginUser,
    logoutUser,
    setUserSession,

    // Admin
    adminToken,
    admin,
    isAdminAuthenticated,
    loginAdmin,
    logoutAdmin,

    // Shared / backward compatibility
    token,
    currentUser,
    isAuthenticated,
    isAdmin,
    initFromStorage,
    login,
    logout,
  }
})

