<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Play, Star } from 'lucide-vue-next'
import { RouterLink, useRouter } from 'vue-router'

import { getMovies } from '@/services/movieService'
import type { Movie } from '@/types/movie'

const router = useRouter()

const bookNow = () => {
  router.push('/movies')
}

function bookMovie(movie: Movie) {
  router.push({
    name: 'showtime-selection',
    params: { movie_id: movie.movie_id },
  })
}

/*
 * ใส่ URL YouTube ของ Trailer ตรงนี้
 * ตัวอย่าง:
 * https://www.youtube.com/watch?v=XXXXXXXXXXX
 */
const trailerUrl =
  'https://www.youtube.com/embed/m-H1cyGpCN8?autoplay=1&mute=1&loop=1&playlist=m-H1cyGpCN8&controls=0&rel=0&modestbranding=1'

// ---------------------------------------------------------------------------
// Now Showing / Coming Soon — same data source and status values as
// MoviesView.vue (GET via movieService.getMovies(), filtered by
// Movie.Status: "NOW_SHOWING" / "COMING_SOON"), just capped to 4 per
// section here to match the homepage mockup.
// ---------------------------------------------------------------------------

const movies = ref<Movie[]>([])
const loading = ref(true)
const error = ref('')

const nowShowing = computed(() => {
  if (!Array.isArray(movies.value)) return []

  return movies.value
    .filter((movie) => movie.status === 'NOW_SHOWING')
    .slice(0, 4)
})

const comingSoon = computed(() => {
  if (!Array.isArray(movies.value)) return []

  return movies.value
    .filter((movie) => movie.status === 'COMING_SOON')
    .slice(0, 4)
})

async function fetchMovies() {
  try {
    loading.value = true
    error.value = ''

    movies.value = await getMovies()
  } catch (err) {
    console.error('Fetch movies error:', err)

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

  return new Date(date).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
  })
}

function getGradient(index: number) {
  const gradients = [
    'bg-gradient-to-b from-rose-900/70 to-black',
    'bg-gradient-to-b from-purple-950/80 to-black',
    'bg-gradient-to-b from-emerald-950/80 to-black',
    'bg-gradient-to-b from-orange-950/80 to-black',
  ]

  return gradients[index % gradients.length]
}

onMounted(fetchMovies)
</script>

<template>
  <div class="min-h-screen overflow-x-hidden bg-[#0B0B0F] text-white">

    <!-- Hero Section -->
    <section class="border-b border-white/10">

      <div
        class="mx-auto grid min-h-[calc(100vh-80px)] w-full max-w-[1500px] items-center gap-10 px-4 py-10 sm:gap-12 sm:px-6 sm:py-14 lg:grid-cols-2 lg:px-12 lg:py-16"
      >

        <!-- LEFT -->
        <div class="min-w-0 max-w-2xl">

          <!-- Tags -->
          <div class="mb-5 flex flex-wrap gap-2 sm:mb-6 sm:gap-3">

            <span
              class="rounded-full border border-[#ff3158]/60 bg-[#ff3158]/10 px-3 py-1.5 text-xs font-semibold sm:px-4 sm:py-2 sm:text-sm"
            >
              NOW SHOWING
            </span>

            <span
              class="rounded-full border border-white/20 bg-white/5 px-3 py-1.5 text-xs font-semibold text-gray-300 sm:px-4 sm:py-2 sm:text-sm"
            >
              ROMANTIC . COMDY FAMILY
            </span>

          </div>

          <!-- Movie title -->
          <h1
            class="break-words text-4xl font-black uppercase leading-[0.92] tracking-tight sm:text-6xl md:text-7xl lg:text-8xl"
          >
            Crimson
            <br />
            HENRY's First Dates
          </h1>

          <!-- Movie information -->
          <div class="mt-6 flex flex-wrap items-center gap-2.5 sm:mt-8 sm:gap-4">

            <div class="flex items-center gap-2 text-yellow-400">
              <Star
                :size="18"
                fill="currentColor"
                class="sm:h-5 sm:w-5"
              />

              <span class="text-sm font-bold sm:text-base">
                8.7
              </span>
            </div>

            <span
              class="rounded-full border border-white/20 bg-white/5 px-3 py-1.5 text-xs text-gray-300 sm:px-4 sm:py-2 sm:text-sm"
            >
              1h 55m
            </span>

            <span
              class="rounded-full border border-white/20 bg-white/5 px-3 py-1.5 text-xs text-gray-300 sm:px-4 sm:py-2 sm:text-sm"
            >
              PG-13
            </span>

          </div>

          <!-- Description -->
          <p
            class="mt-6 max-w-2xl text-sm leading-6 text-gray-300 sm:mt-8 sm:text-base sm:leading-7 lg:text-lg lg:leading-8"
          >
            *HENRY's First Dates* is a new feel-good romance from GDH. After "Henry" (Nadech Kugimiya) suffers an accident that resets his memory every morning, "Lucy" (Minnie Nicha)—a woman who has fallen for him—must wake up each day to embark on the mission of wooing him all over again. Come cheer on this adorable romance in theaters!
          </p>

          <!-- Buttons -->
          <div class="mt-7 flex flex-col gap-3 sm:mt-10 sm:flex-row sm:flex-wrap sm:gap-4">

            <!-- Book now -->
            <button
              type="button"
              class="flex w-full items-center justify-center gap-3 rounded-xl bg-[#ff3158] px-6 py-3.5 text-sm font-bold shadow-lg shadow-[#ff3158]/20 transition hover:-translate-y-0.5 hover:bg-[#ff4568] sm:w-auto sm:px-7 sm:py-4 sm:text-base"
              @click="bookNow"
            >
              <Play
                :size="18"
                fill="currentColor"
              />

              Book now
            </button>

          </div>

        </div>

        <!-- RIGHT : TRAILER VIDEO -->
        <div
          class="relative mx-auto w-full max-w-2xl overflow-hidden rounded-2xl border border-white/20 shadow-2xl sm:rounded-3xl"
        >

          <!-- Video container -->
          <div
            class="relative aspect-[16/10] bg-black"
          >

            <iframe
              :src="trailerUrl"
              title="Crimson Horizon Trailer"
              class="absolute inset-0 h-full w-full"
              frameborder="0"
              allow="autoplay; encrypted-media; picture-in-picture"
              allowfullscreen
            ></iframe>

            <!-- Dark gradient -->
            <div
              class="pointer-events-none absolute inset-x-0 bottom-0 h-32 bg-gradient-to-t from-black/90 via-black/30 to-transparent sm:h-48"
            ></div>

            <!-- IMAX -->
            <div
              class="pointer-events-none absolute left-3 top-3 rounded-full border border-purple-400/40 bg-purple-500/10 px-3 py-1.5 text-xs font-semibold text-purple-300 backdrop-blur-sm sm:left-6 sm:top-6 sm:px-5 sm:py-2 sm:text-sm"
            >
              IMAX
            </div>

            <!-- Rating -->
            <div
              class="pointer-events-none absolute right-3 top-3 rounded-full border border-yellow-500/30 bg-yellow-500/10 px-3 py-1.5 text-xs font-bold text-yellow-400 backdrop-blur-sm sm:right-6 sm:top-6 sm:px-4 sm:py-2 sm:text-sm"
            >
              ★ 8.7
            </div>

            <!-- Movie title -->
            <div
              class="pointer-events-none absolute bottom-4 left-4 right-4 sm:bottom-7 sm:left-7 sm:right-auto"
            >
              <h2
                class="text-lg font-black uppercase leading-tight tracking-wide sm:text-2xl md:text-3xl"
              >
                Crimson HENRY's First Dates
              </h2>

              <p
                class="mt-1 text-[10px] font-semibold text-gray-300 sm:text-sm"
              >
                IN CINEMAS NOW
              </p>
            </div>

          </div>

        </div>

      </div>

    </section>

    <!-- ============================================================= -->
    <!-- NOW SHOWING + COMING SOON — same data as MoviesView.vue        -->
    <!-- ============================================================= -->

    <section
      class="mx-auto w-full max-w-[1500px] px-4 py-10 sm:px-6 sm:py-14 lg:px-12"
    >

      <!-- Error banner -->
      <div
        v-if="error"
        class="mb-6 rounded-xl border border-red-500/20 bg-red-500/5 p-4 text-center sm:mb-8 sm:p-6"
      >
        <p class="text-sm text-red-400">
          {{ error }}
        </p>

        <button
          type="button"
          class="mt-3 rounded-lg border border-white/10 px-4 py-2 text-sm text-gray-300 transition hover:bg-white hover:text-black sm:px-5"
          @click="fetchMovies"
        >
          Try again
        </button>
      </div>

      <!-- ========================= -->
      <!-- NOW SHOWING -->
      <!-- ========================= -->

      <div class="mb-5 flex items-center justify-between gap-4 sm:mb-6">
        <h2 class="text-lg font-black uppercase tracking-tight sm:text-xl">
          Now Showing
        </h2>

        <RouterLink
          to="/movies"
          class="shrink-0 text-xs font-semibold text-gray-400 transition hover:text-white sm:text-sm"
        >
          See all →
        </RouterLink>
      </div>

      <div
        v-if="loading"
        class="mb-10 grid grid-cols-1 gap-3 sm:grid-cols-2 sm:gap-4 lg:mb-14 lg:grid-cols-4"
      >
        <div
          v-for="n in 4"
          :key="n"
          class="h-64 animate-pulse rounded-xl border border-white/10 bg-[#15151d]"
        ></div>
      </div>

      <div
        v-else-if="nowShowing.length"
        class="mb-10 grid grid-cols-1 gap-3 sm:grid-cols-2 sm:gap-4 lg:mb-14 lg:grid-cols-4"
      >
        <article
          v-for="(movie, index) in nowShowing"
          :key="movie.movie_id"
          class="min-w-0 overflow-hidden rounded-xl border border-white/10 bg-[#15151d] transition duration-300 hover:-translate-y-1 hover:border-white/20"
        >
          <!-- Poster -->
          <RouterLink
            :to="{ name: 'movie-detail', params: { movie_id: movie.movie_id } }"
            :class="[
              'relative block h-48 overflow-hidden sm:h-40',
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
          </RouterLink>

          <!-- Information -->
          <div class="p-3">
            <h3 class="truncate text-sm font-bold">
              {{ movie.title }}
            </h3>

            <p
              class="mt-2 flex min-w-0 items-center gap-1 truncate text-xs text-gray-500"
            >
              <span class="min-w-0 truncate">
                {{ movie.genre.join(' · ') }}
              </span>

              <span
                v-if="movie.rating"
                class="ml-auto flex shrink-0 items-center gap-1 text-yellow-400"
              >
                <Star
                  :size="12"
                  fill="currentColor"
                />
                {{ movie.rating }}
              </span>
            </p>

            <button
              type="button"
              class="mt-3 block w-full rounded-lg border border-white/20 py-2.5 text-center text-xs font-semibold transition hover:bg-white hover:text-black"
              @click="bookMovie(movie)"
            >
              Book now
            </button>
          </div>
        </article>
      </div>

      <p
        v-else-if="!loading"
        class="mb-10 py-8 text-center text-sm text-gray-500 lg:mb-14 lg:py-10"
      >
        No movies currently showing.
      </p>

      <!-- ========================= -->
      <!-- COMING SOON -->
      <!-- ========================= -->

      <div class="mb-5 flex items-center justify-between gap-4 sm:mb-6">
        <h2 class="text-lg font-black uppercase tracking-tight sm:text-xl">
          Coming Soon
        </h2>

        <RouterLink
          to="/movies"
          class="shrink-0 text-xs font-semibold text-gray-400 transition hover:text-white sm:text-sm"
        >
          See all →
        </RouterLink>
      </div>

      <div
        v-if="!loading && comingSoon.length"
        class="grid grid-cols-1 gap-3 sm:grid-cols-2 sm:gap-4 lg:grid-cols-4"
      >
        <article
          v-for="(movie, index) in comingSoon"
          :key="movie.movie_id"
          class="min-w-0 overflow-hidden rounded-xl border border-white/10 bg-[#15151d] transition duration-300 hover:-translate-y-1 hover:border-white/20"
        >
          <!-- Poster -->
          <RouterLink
            :to="{ name: 'movie-detail', params: { movie_id: movie.movie_id } }"
            :class="[
              'relative block h-48 overflow-hidden sm:h-40',
              getGradient(index),
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
          </RouterLink>

          <!-- Information -->
          <div class="p-3">
            <h3 class="truncate text-sm font-bold">
              {{ movie.title }}
            </h3>

            <p class="mt-2 truncate text-xs text-gray-500">
              {{ movie.genre.join(' · ') }}
            </p>

            <p class="mt-1 text-xs text-gray-500">
              {{ formatReleaseDate(movie.release_date) }}
            </p>

            <!-- Not yet bookable -->
            <button
              type="button"
              disabled
              class="mt-3 block w-full cursor-not-allowed rounded-lg border border-white/10 py-2.5 text-center text-xs font-semibold text-gray-600"
            >
              Book now
            </button>
          </div>
        </article>
      </div>

      <p
        v-else-if="!loading"
        class="py-8 text-center text-sm text-gray-500 sm:py-10"
      >
        No upcoming movies.
      </p>

    </section>

  </div>
</template>