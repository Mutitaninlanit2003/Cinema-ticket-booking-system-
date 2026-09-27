<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import { Search, X, User, LogOut, Film } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'

const router = useRouter()
const authStore = useAuthStore()

const showProfileMenu = ref(false)

// ==========================================
// User Authentication
// ==========================================

const isLoggedIn = computed(() => {
  return authStore.isUserAuthenticated
})

const user = computed(() => {
  return authStore.user
})

const userInitial = computed(() => {
  if (!user.value?.name) {
    return 'U'
  }

  return user.value.name.charAt(0).toUpperCase()
})

const avatarUrl = computed(() => {
  const imagePath = user.value?.profile_image

  if (!imagePath) {
    return null
  }

  if (
    imagePath.startsWith('http://') ||
    imagePath.startsWith('https://')
  ) {
    return imagePath
  }

  const apiBase =
    import.meta.env.VITE_API_BASE_URL ||
    'http://localhost:8080'

  const base = apiBase.endsWith('/api')
    ? apiBase.slice(0, -4)
    : apiBase

  return `${base}${imagePath}`
})

function goToProfile() {
  showProfileMenu.value = false
  router.push('/profile')
}

function logout() {
  showProfileMenu.value = false
  authStore.logoutUser()
  router.push('/login')
}

// ==========================================
// Movie Search — rebuilt from scratch, self
// contained, no external service dependency.
// ==========================================

interface SearchMovie {
  movie_id: string
  title: string
  poster_url?: string
  genre?: string[]
}

const searchOpen = ref(false)
const searchQuery = ref('')
const searchLoading = ref(false)
const searchError = ref('')
const allMovies = ref<SearchMovie[]>([])
const moviesFetched = ref(false)

const searchPanelEl = ref<HTMLElement | null>(null)
const searchInputEl = ref<HTMLInputElement | null>(null)

const searchResults = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()

  if (!query) {
    return []
  }

  return allMovies.value
    .filter((movie) =>
      (movie.title ?? '').toLowerCase().includes(query),
    )
    .slice(0, 8)
})

async function fetchMoviesOnce() {
  if (moviesFetched.value) return

  searchLoading.value = true
  searchError.value = ''

  try {
    const { data } = await api.get('/movies')
    allMovies.value = Array.isArray(data) ? data : data?.data ?? []
    moviesFetched.value = true
  } catch (err) {
    console.error('Failed to load movies for search:', err)
    searchError.value = 'Could not load movies. Please try again.'
  } finally {
    searchLoading.value = false
  }
}

function openSearch() {
  searchOpen.value = true

  nextTick(() => {
    searchInputEl.value?.focus()
  })

  fetchMoviesOnce()
}

function closeSearch() {
  searchOpen.value = false
  searchQuery.value = ''
}

function handleSearchIconClick() {
  if (searchOpen.value) {
    closeSearch()
  } else {
    openSearch()
  }
}

function selectMovie(movie: SearchMovie) {
  closeSearch()

  router.push({
    name: 'movie-detail',
    params: { movie_id: movie.movie_id },
  })
}

function handleOutsideClick(event: MouseEvent) {
  if (!searchOpen.value) return
  if (!searchPanelEl.value) return

  if (!searchPanelEl.value.contains(event.target as Node)) {
    closeSearch()
  }
}

function handleEscapeKey(event: KeyboardEvent) {
  if (event.key === 'Escape' && searchOpen.value) {
    closeSearch()
  }
}

onMounted(() => {
  authStore.initFromStorage()
  document.addEventListener('mousedown', handleOutsideClick)
  document.addEventListener('keydown', handleEscapeKey)
})

onUnmounted(() => {
  document.removeEventListener('mousedown', handleOutsideClick)
  document.removeEventListener('keydown', handleEscapeKey)
})
</script>

<template>
  <header
    class="relative z-50 border-b border-white/10 bg-[#0B0B0F]"
  >
    <nav
      class="mx-auto flex h-[116px] max-w-[1500px] items-center justify-between px-6 lg:px-12"
    >

      <!-- Logo -->
      <router-link
        to="/"
        class="flex items-center gap-3"
      >
        <div
          class="flex h-10 w-10 items-center justify-center rounded-lg bg-gradient-to-br from-pink-500 to-purple-500"
        >
          <span class="text-lg font-bold">▦</span>
        </div>

        <div class="text-2xl font-black tracking-tight text-white">
          CINEMA
          <span class="text-[#ff3158]">.</span>
        </div>
      </router-link>

      <!-- Menu -->
      <div class="hidden items-center gap-10 text-[16px] font-medium md:flex">
        <router-link to="/" class="text-white transition hover:text-[#ff3158]">
          Home
        </router-link>

        <router-link to="/movies" class="text-gray-300 transition hover:text-white">
          Movies
        </router-link>

        <router-link to="/bookings" class="text-gray-300 transition hover:text-white">
          My Bookings
        </router-link>
      </div>

      <!-- Right -->
      <div class="flex items-center gap-6">

        <!-- ============================= -->
        <!-- SEARCH (rebuilt) -->
        <!-- ============================= -->
        <div ref="searchPanelEl" class="relative">

          <button
            type="button"
            class="flex h-10 w-10 items-center justify-center rounded-full text-gray-300 transition hover:bg-white/5 hover:text-white"
            :class="{ 'bg-white/10 text-white': searchOpen }"
            aria-label="Search movies"
            @click="handleSearchIconClick"
          >
            <component :is="searchOpen ? X : Search" :size="21" />
          </button>

          <Transition name="search-panel">
            <div
              v-if="searchOpen"
              class="absolute right-0 top-full mt-3 w-80 rounded-2xl border border-white/10 bg-[#15151d] p-3 shadow-2xl sm:w-96"
            >

              <input
                ref="searchInputEl"
                v-model="searchQuery"
                type="text"
                placeholder="Search movies by title..."
                autocomplete="off"
                class="w-full rounded-xl border border-white/10 bg-[#0d0d10] px-4 py-3 text-sm text-white outline-none transition placeholder:text-gray-600 focus:border-[#ff3158]/60"
              />

              <div class="mt-2">

                <p
                  v-if="searchLoading"
                  class="px-2 py-4 text-center text-xs text-gray-500"
                >
                  Loading movies...
                </p>

                <p
                  v-else-if="searchError"
                  class="px-2 py-4 text-center text-xs text-red-400"
                >
                  {{ searchError }}
                </p>

                <p
                  v-else-if="!searchQuery.trim()"
                  class="px-2 py-4 text-center text-xs text-gray-500"
                >
                  Start typing a movie title...
                </p>

                <p
                  v-else-if="searchResults.length === 0"
                  class="px-2 py-4 text-center text-xs text-gray-500"
                >
                  No movies found for "{{ searchQuery }}".
                </p>

                <ul v-else class="max-h-80 space-y-1 overflow-y-auto">
                  <li v-for="movie in searchResults" :key="movie.movie_id">
                    <button
                      type="button"
                      class="flex w-full items-center gap-3 rounded-xl px-2 py-2 text-left transition hover:bg-white/5"
                      @click="selectMovie(movie)"
                    >
                      <div
                        class="flex h-14 w-10 shrink-0 items-center justify-center overflow-hidden rounded-lg bg-gradient-to-br from-purple-900 to-black"
                      >
                        <img
                          v-if="movie.poster_url"
                          :src="movie.poster_url"
                          :alt="movie.title"
                          class="h-full w-full object-cover"
                        />

                        <Film v-else :size="16" class="text-gray-500" />
                      </div>

                      <div class="min-w-0">
                        <p class="truncate text-sm font-semibold text-white">
                          {{ movie.title }}
                        </p>

                        <p
                          v-if="movie.genre?.length"
                          class="truncate text-xs text-gray-500"
                        >
                          {{ movie.genre.join(' · ') }}
                        </p>
                      </div>
                    </button>
                  </li>
                </ul>

              </div>

            </div>
          </Transition>

        </div>

        <!-- Before Login -->
        <router-link
          v-if="!isLoggedIn"
          to="/login"
          class="rounded-xl bg-gradient-to-r from-[#ff2d55] to-[#ff4d6d] px-7 py-3.5 font-bold text-white shadow-lg shadow-[#ff2d55]/25 transition hover:-translate-y-0.5"
        >
          Login
        </router-link>

        <!-- After Login -->
        <div v-else class="relative">

          <button
            type="button"
            @click="showProfileMenu = !showProfileMenu"
            class="group flex h-11 w-11 items-center justify-center overflow-hidden rounded-full border-2 border-white/10 bg-gradient-to-br from-[#ff3158] to-purple-500 shadow-lg shadow-[#ff3158]/20 transition hover:scale-105 hover:border-[#ff3158]/60"
          >
            <img
              v-if="avatarUrl"
              :src="avatarUrl"
              alt="Profile"
              class="h-full w-full object-cover"
            />

            <span v-else class="text-base font-bold text-white">
              {{ userInitial }}
            </span>
          </button>

          <div
            v-if="showProfileMenu"
            class="absolute right-0 mt-3 w-60 overflow-hidden rounded-2xl border border-white/10 bg-[#15151d] shadow-2xl"
          >

            <div class="border-b border-white/10 px-4 py-4">
              <div class="flex items-center gap-3">

                <div
                  class="flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-full bg-gradient-to-br from-[#ff3158] to-purple-500"
                >
                  <img
                    v-if="avatarUrl"
                    :src="avatarUrl"
                    alt="Profile"
                    class="h-full w-full object-cover"
                  />

                  <span v-else class="font-bold text-white">
                    {{ userInitial }}
                  </span>
                </div>

                <div class="min-w-0">
                  <p class="truncate font-semibold text-white">
                    {{ user?.name || 'User' }}
                  </p>

                  <p class="truncate text-xs text-gray-400">
                    {{ user?.email || '' }}
                  </p>
                </div>

              </div>
            </div>

            <button
              type="button"
              @click="goToProfile"
              class="flex w-full items-center gap-3 px-4 py-3 text-left text-sm text-gray-300 transition hover:bg-white/5 hover:text-white"
            >
              <User :size="18" />
              Profile
            </button>

            <button
              type="button"
              @click="logout"
              class="flex w-full items-center gap-3 border-t border-white/10 px-4 py-3 text-left text-sm text-red-400 transition hover:bg-red-500/10"
            >
              <LogOut :size="18" />
              Logout
            </button>

          </div>

        </div>

      </div>
    </nav>
  </header>
</template>

<style scoped>
header,
nav {
  max-width: 100%;
}

.search-panel-enter-active,
.search-panel-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.search-panel-enter-from,
.search-panel-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>