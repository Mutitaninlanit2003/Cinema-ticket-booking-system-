import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

// =====================================================
// Get JWT token from sessionStorage
// =====================================================

function getToken(storageKey: 'user_auth' | 'admin_auth'): string | null {
  const raw = sessionStorage.getItem(storageKey)

  if (!raw) {
    console.warn(`⚠️ ${storageKey} not found in sessionStorage`)
    return null
  }

  try {
    const parsed = JSON.parse(raw)

    // ---------------------------------------------
    // Normal format
    // {
    //   token: "eyJ..."
    // }
    // ---------------------------------------------

    if (
      parsed &&
      typeof parsed.token === 'string' &&
      parsed.token.trim() !== ''
    ) {
      return parsed.token.trim()
    }

    // ---------------------------------------------
    // In case token is nested
    // {
    //   data: {
    //     token: "eyJ..."
    //   }
    // }
    // ---------------------------------------------

    if (
      parsed?.data &&
      typeof parsed.data.token === 'string' &&
      parsed.data.token.trim() !== ''
    ) {
      return parsed.data.token.trim()
    }

    console.warn(`⚠️ JWT token not found inside ${storageKey}`)

    return null
  } catch (error) {
    console.error(`❌ Failed to parse ${storageKey}:`, error)

    return null
  }
}

// =====================================================
// Request Interceptor
// =====================================================

api.interceptors.request.use(
  (config) => {
    const currentPath = window.location.pathname

    // ---------------------------------------------
    // Determine whether current page is Admin
    // ---------------------------------------------

    const isAdminPage = currentPath.startsWith('/admin')

    const storageKey: 'user_auth' | 'admin_auth' = isAdminPage
      ? 'admin_auth'
      : 'user_auth'

    const token = getToken(storageKey)

    // ---------------------------------------------
    // Debug
    // ---------------------------------------------

    console.log('========== API REQUEST ==========')
    console.log('URL:', config.url)
    console.log('PATH:', currentPath)
    console.log('ADMIN PAGE:', isAdminPage)
    console.log('AUTH STORAGE:', storageKey)
    console.log('TOKEN EXISTS:', !!token)

    // ---------------------------------------------
    // Attach JWT
    // ---------------------------------------------

    if (token) {
      config.headers = config.headers ?? {}

      config.headers.Authorization = `Bearer ${token}`

      console.log('Authorization header attached')
    } else {
      if (config.headers) {
        delete config.headers.Authorization
      }

      console.warn(
        `⚠️ No JWT token available from sessionStorage key: ${storageKey}`,
      )
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)

// =====================================================
// Response Interceptor
// =====================================================

api.interceptors.response.use(
  (response) => {
    return response
  },

  (error) => {
    const status = error?.response?.status
    const currentPath = window.location.pathname

    const isAdminPage = currentPath.startsWith('/admin')

    const storageKey: 'user_auth' | 'admin_auth' = isAdminPage
      ? 'admin_auth'
      : 'user_auth'

    console.error('========== API ERROR ==========')
    console.error('URL:', error?.config?.url)
    console.error('STATUS:', status)
    console.error('DATA:', error?.response?.data)
    console.error('AUTH STORAGE:', storageKey)

    // ---------------------------------------------
    // Unauthorized
    // ---------------------------------------------

    if (status === 401) {
      console.error('❌ Unauthorized')

      /*
       * Do NOT immediately remove the token here.
       *
       * ตอนนี้เรากำลัง debug ระบบ Admin JWT
       * ถ้าลบทันทีจะทำให้ตรวจสอบ token ได้ยาก
       *
       * หลังจากระบบทำงานปกติแล้วค่อยเพิ่ม
       * sessionStorage.removeItem(storageKey)
       * กลับมาได้
       */

      if (isAdminPage) {
        console.warn('⚠️ Admin API returned 401')
      } else {
        console.warn('⚠️ User API returned 401')
      }
    }

    const message =
      error?.response?.data?.error ??
      error?.response?.data?.message ??
      error?.message ??
      'Unexpected error'

    return Promise.reject(new Error(message))
  },
)

export default api