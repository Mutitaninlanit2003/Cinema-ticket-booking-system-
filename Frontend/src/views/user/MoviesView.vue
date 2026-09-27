<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'

import { RouterLink, useRoute } from 'vue-router'

import { getMovies } from '@/services/movieService'

import type { Movie } from '@/types/movie'

// import MovieDetailView from '@/views/user/MovieDetailView.vue'

const route = useRoute()

const movies = ref<Movie[]>([])

const loading = ref(true)

const error = ref('')

// ค่าค้นหาปัจจุบัน มาจาก ?search= ใน URL (ตั้งโดย Navbar)
const searchQuery = ref(
  typeof route.query.search === 'string'
    ? route.query.search
    : '',
)

watch(
  () => route.query.search,
  (value) => {
    searchQuery.value =
      typeof value === 'string' ? value : ''
  },
)

const filteredMovies = computed(() => {
  if (!Array.isArray(movies.value)) {
    return []
  }

  const q = searchQuery.value.trim().toLowerCase()

  if (!q) {
    return movies.value
  }

  return movies.value.filter((movie) => {
    const title =
      movie.title?.toLowerCase() || ''

    const genres =
      (movie.genre || [])
        .join(' ')
        .toLowerCase()

    return (
      title.includes(q) ||
      genres.includes(q)
    )
  })
})

const nowShowing = computed(() => {
  return filteredMovies.value.filter(
    (movie) =>
      movie.status === 'NOW_SHOWING',
  )
})

const comingSoon = computed(() => {
  return filteredMovies.value.filter(
    (movie) =>
      movie.status === 'COMING_SOON',
  )
})

async function fetchMovies() {
  try {
    loading.value = true
    error.value = ''

    const data = await getMovies()

    console.log('Movies loaded:', data)

    movies.value = data
  } catch (err) {
    console.error(
      'Fetch movies error:',
      err,
    )

    error.value =
      err instanceof Error
        ? err.message
        : 'Unable to load movies'
  } finally {
    loading.value = false
  }
}

function formatReleaseDate(date?: string) {
  if (!date) return '-'

  return new Date(date).toLocaleDateString(
    'en-US',
    {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    },
  )
}

function getGradient(index: number) {
  const gradients = [
    'bg-gradient-to-b from-rose-900/70 to-black',
    'bg-gradient-to-b from-emerald-950/80 to-black',
    'bg-gradient-to-b from-purple-950/80 to-black',
    'bg-gradient-to-b from-orange-950/80 to-black',
    'bg-gradient-to-b from-slate-900 to-black',
    'bg-gradient-to-b from-purple-950 to-black',
    'bg-gradient-to-b from-orange-950 to-black',
    'bg-gradient-to-b from-cyan-950 to-black',
  ]

  return gradients[index % gradients.length]
}

onMounted(fetchMovies)
</script>

<template>
  <div
    class="min-h-screen overflow-x-hidden bg-[#0d0d12] px-4 py-6 text-white sm:px-6 sm:py-8"
  >
    <div class="mx-auto w-full max-w-7xl">

      <!-- ========================= -->
      <!-- SEARCH STATUS -->
      <!-- ========================= -->

      <div
        v-if="searchQuery.trim()"
        class="mb-6 flex flex-col gap-3 rounded-xl border border-white/10 bg-white/5 px-4 py-3 sm:mb-8 sm:flex-row sm:items-center sm:justify-between sm:px-5"
      >
        <p
          class="min-w-0 break-words text-xs text-gray-300 sm:text-sm"
        >
          ผลการค้นหาสำหรับ

          <span class="font-semibold text-white">
            "{{ searchQuery }}"
          </span>

          — พบ {{ filteredMovies.length }} เรื่อง
        </p>

        <RouterLink
          to="/movies"
          class="shrink-0 text-xs font-semibold text-[#ff3158] hover:underline sm:text-sm"
        >
          ล้างการค้นหา
        </RouterLink>
      </div>

      <!-- ========================= -->
      <!-- NOW SHOWING -->
      <!-- ========================= -->

      <section>
        <div
          class="mb-5 flex flex-wrap items-center justify-between gap-3 sm:mb-6"
        >
          <h2
            class="text-lg font-black uppercase tracking-tight sm:text-xl"
          >
            Now Showing
          </h2>

          <span
            class="rounded-full border border-white/20 px-2.5 py-1 text-[10px] text-gray-400 sm:px-3 sm:text-xs"
          >
            status = NOW_SHOWING
          </span>
        </div>

        <!-- Loading -->
        <div
          v-if="loading"
          class="flex min-h-[250px] items-center justify-center sm:min-h-[300px]"
        >
          <p class="text-sm text-gray-500">
            Loading movies...
          </p>
        </div>

        <!-- Error -->
        <div
          v-else-if="error"
          class="rounded-xl border border-red-500/20 bg-red-500/5 p-5 text-center sm:p-8"
        >
          <p class="text-sm text-red-400">
            {{ error }}
          </p>

          <button
            type="button"
            class="mt-4 rounded-lg border border-white/10 px-4 py-2 text-sm text-gray-300 transition hover:bg-white hover:text-black sm:px-5"
            @click="fetchMovies"
          >
            Try again
          </button>
        </div>

        <!-- Movies -->
        <div
          v-else-if="nowShowing.length"
          class="grid grid-cols-1 gap-4 min-[420px]:grid-cols-2 lg:grid-cols-4"
        >
          <article
            v-for="(movie, index) in nowShowing"
            :key="movie.movie_id"
            class="min-w-0 overflow-hidden rounded-xl border border-white/10 bg-[#15151d] transition duration-300 hover:-translate-y-1 hover:border-white/20"
          >
            <!-- Poster -->
            <div
              :class="[
                'relative aspect-[2/3] overflow-hidden',
                getGradient(index),
              ]"
            >
              <img
                v-if="movie.poster_url"
                :src="movie.poster_url"
                :alt="movie.title"
                class="h-full w-full object-cover transition duration-300 sm:hover:scale-105"
              />

              <!-- fallback -->
              <div
                v-else
                class="absolute inset-0 flex items-end p-3"
              >
                <h3
                  class="max-w-[150px] text-sm font-black uppercase leading-tight"
                >
                  {{ movie.title }}
                </h3>
              </div>
            </div>

            <!-- Information -->
            <div class="p-3">
              <h3
                class="truncate text-sm font-bold"
              >
                {{ movie.title }}
              </h3>

              <p
                class="mt-2 truncate text-xs text-gray-500"
              >
                {{ movie.genre.join(' · ') }}
                · {{ movie.duration }} min
              </p>

              <p
                class="mt-1 truncate text-xs text-gray-500"
              >
                {{ movie.language }}
                · {{ movie.subtitle }}
                · {{ movie.rating }}
              </p>

              <RouterLink
                :to="{
                  name: 'movie-detail',
                  params: {
                    movie_id:
                      movie.movie_id,
                  },
                }"
                class="mt-3 block rounded-lg border border-white/20 py-2.5 text-center text-xs font-semibold transition hover:bg-white hover:text-black"
              >
                View detail
              </RouterLink>
            </div>
          </article>
        </div>

        <p
          v-else
          class="py-12 text-center text-sm text-gray-500 sm:py-16"
        >
          No movies currently showing.
        </p>
      </section>

      <!-- ========================= -->
      <!-- COMING SOON -->
      <!-- ========================= -->

      <section class="mt-10 sm:mt-12">
        <div
          class="mb-5 flex flex-wrap items-center justify-between gap-3 sm:mb-6"
        >
          <h2
            class="text-lg font-black uppercase tracking-tight sm:text-xl"
          >
            Coming Soon
          </h2>

          <span
            class="rounded-full border border-white/20 px-2.5 py-1 text-[10px] text-gray-400 sm:px-3 sm:text-xs"
          >
            status = COMING_SOON
          </span>
        </div>

        <div
          v-if="comingSoon.length"
          class="grid grid-cols-1 gap-4 min-[420px]:grid-cols-2 lg:grid-cols-4"
        >
          <article
            v-for="(movie, index) in comingSoon"
            :key="movie.movie_id"
            class="min-w-0 overflow-hidden rounded-xl border border-white/10 bg-[#15151d] transition duration-300 hover:-translate-y-1 hover:border-white/20"
          >
            <!-- Poster -->
            <div
              :class="[
                'relative aspect-[2/3] overflow-hidden',
                getGradient(index + 4),
              ]"
            >
              <img
                v-if="movie.poster_url"
                :src="movie.poster_url"
                :alt="movie.title"
                class="h-full w-full object-cover transition duration-300 sm:hover:scale-105"
              />

              <!-- Coming soon badge -->
              <span
                class="absolute left-0 top-3 rounded-r-md bg-violet-500 px-2 py-1 text-[9px] font-bold sm:px-2.5 sm:py-1 sm:text-[10px]"
              >
                COMING SOON
              </span>

              <!-- fallback -->
              <div
                v-if="!movie.poster_url"
                class="absolute inset-0 flex items-end p-3"
              >
                <h3
                  class="max-w-[150px] text-sm font-black uppercase leading-tight"
                >
                  {{ movie.title }}
                </h3>
              </div>
            </div>

            <!-- Information -->
            <div class="p-3">
              <h3
                class="truncate text-sm font-bold"
              >
                {{ movie.title }}
              </h3>

              <p
                class="mt-2 truncate text-xs text-gray-500"
              >
                {{ movie.genre.join(' · ') }}
                · {{ movie.duration }} min
              </p>

              <p
                class="mt-1 truncate text-xs text-gray-500"
              >
                Release:
                {{ formatReleaseDate(movie.release_date) }}
              </p>

              <RouterLink
                :to="`/movies/${movie.movie_id}`"
                class="mt-3 block rounded-lg border border-white/20 py-2.5 text-center text-xs font-semibold transition hover:bg-white hover:text-black"
              >
                View detail
              </RouterLink>
            </div>
          </article>
        </div>

        <p
          v-else
          class="py-12 text-center text-sm text-gray-500 sm:py-16"
        >
          No upcoming movies.
        </p>
      </section>

    </div>
  </div>
</template>