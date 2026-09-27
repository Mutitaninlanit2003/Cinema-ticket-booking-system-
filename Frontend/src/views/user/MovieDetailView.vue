<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getMovieById } from '../../services/movieService'
import type { Movie } from '../../types/movie'

const route = useRoute()
const router = useRouter()

const movie = ref<Movie | null>(null)
const loading = ref(true)
const error = ref('')

const movieId = computed(() => {
  return String(route.params.movie_id)
})

async function fetchMovie() {
  try {
    loading.value = true
    error.value = ''

    movie.value = await getMovieById(movieId.value)

    if (!movie.value) {
      throw new Error('Movie not found')
    }
  } catch (err) {
    error.value =
      err instanceof Error
        ? err.message
        : 'Unable to load movie'
  } finally {
    loading.value = false
  }
}

function goBack() {
  router.push('/movies')
}

function bookNow() {
  if (!movie.value) return

  router.push(
    `/movies/${movie.value.movie_id}/showtimes`,
  )
}

function watchTrailer() {
  if (!movie.value?.trailer_url) return

  window.open(
    movie.value.trailer_url,
    '_blank',
  )
}

function formatReleaseDate(date: string) {
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

onMounted(fetchMovie)
</script>

<template>
  <div
    class="min-h-screen overflow-x-hidden bg-[#0d0d12] px-4 py-5 text-white sm:px-6 sm:py-6"
  >
    <div class="mx-auto w-full max-w-7xl">

      <!-- Breadcrumb -->
      <button
        type="button"
        class="mb-5 flex min-w-0 max-w-full items-center gap-2 text-xs text-gray-400 transition hover:text-white sm:mb-6 sm:text-sm"
        @click="goBack"
      >
        <span class="shrink-0 text-lg">
          ‹
        </span>

        <span class="shrink-0">
          Movies /
        </span>

        <span
          v-if="movie"
          class="min-w-0 truncate text-gray-300"
        >
          {{ movie.title }}
        </span>
      </button>

      <!-- Loading -->
      <div
        v-if="loading"
        class="flex min-h-[300px] items-center justify-center sm:min-h-[500px]"
      >
        <p class="text-sm text-gray-500 sm:text-base">
          Loading movie...
        </p>
      </div>

      <!-- Error -->
      <div
        v-else-if="error"
        class="rounded-2xl border border-red-500/20 bg-red-500/5 p-5 text-center sm:p-10"
      >
        <p class="break-words text-sm text-red-400 sm:text-base">
          {{ error }}
        </p>

        <button
          type="button"
          class="mt-5 rounded-lg border border-white/10 px-4 py-2 text-sm transition hover:bg-white hover:text-black sm:px-5"
          @click="fetchMovie"
        >
          Try again
        </button>
      </div>

      <!-- Movie Detail -->
      <div
        v-else-if="movie"
        class="overflow-hidden rounded-2xl border border-white/10 bg-[#101016] sm:rounded-3xl"
      >
        <div
          class="grid gap-7 p-4 sm:gap-8 sm:p-6 lg:grid-cols-[0.85fr_1.15fr] lg:p-10"
        >

          <!-- LEFT : POSTER -->
          <div class="min-w-0">
            <div
              class="relative aspect-[2/3] w-full overflow-hidden rounded-2xl border border-white/10 bg-gradient-to-b from-purple-950/70 to-black sm:aspect-auto sm:min-h-[500px]"
            >
              <img
                v-if="movie.poster_url"
                :src="movie.poster_url"
                :alt="movie.title"
                class="absolute inset-0 h-full w-full object-cover"
              />

              <div
                v-else
                class="absolute inset-0 flex items-center justify-center p-6 sm:p-8"
              >
                <h2
                  class="text-center text-2xl font-black uppercase sm:text-3xl"
                >
                  {{ movie.title }}
                </h2>
              </div>

              <!-- Dark overlay -->
              <div
                class="absolute inset-0 bg-gradient-to-t from-black/70 via-transparent to-transparent"
              ></div>

              <!-- Trailer button -->
              <button
                v-if="movie.trailer_url"
                type="button"
                class="absolute left-1/2 top-1/2 flex h-14 w-14 -translate-x-1/2 -translate-y-1/2 items-center justify-center rounded-full border border-white/40 bg-black/20 text-lg text-white backdrop-blur-sm transition hover:scale-110 hover:bg-white hover:text-black sm:h-16 sm:w-16 sm:text-xl"
                @click="watchTrailer"
              >
                ▶
              </button>

              <!-- Rating -->
              <div
                class="absolute right-3 top-3 rounded-full border border-amber-400/40 bg-amber-500/10 px-3 py-1.5 text-xs font-bold text-amber-300 sm:right-5 sm:top-5 sm:px-4 sm:py-2 sm:text-sm"
              >
                {{ movie.rating }}
              </div>
            </div>
          </div>

          <!-- RIGHT : INFORMATION -->
          <div
            class="flex min-w-0 flex-col justify-center"
          >

            <!-- Genres + Status -->
            <div
              class="mb-4 flex flex-wrap gap-2 sm:mb-5"
            >
              <span
                v-for="genre in movie.genre"
                :key="genre"
                class="rounded-full border border-pink-500/50 bg-pink-500/10 px-3 py-1.5 text-[10px] font-bold uppercase text-pink-300 sm:px-4 sm:text-xs"
              >
                {{ genre }}
              </span>

              <span
                class="rounded-full bg-emerald-500/10 px-3 py-1.5 text-[10px] font-bold uppercase text-emerald-400 sm:px-4 sm:text-xs"
              >
                {{ movie.status }}
              </span>
            </div>

            <!-- Title -->
            <h1
              class="break-words text-2xl font-black uppercase leading-tight tracking-tight sm:text-4xl md:text-5xl"
            >
              {{ movie.title }}
            </h1>

            <!-- Description -->
            <p
              class="mt-5 max-w-3xl text-sm leading-6 text-gray-300 sm:mt-6 sm:text-base sm:leading-7"
            >
              {{ movie.description }}
            </p>

            <div
              class="my-6 h-px bg-white/10 sm:my-7"
            ></div>

            <!-- Movie information -->
            <div
              class="grid grid-cols-2 gap-x-5 gap-y-6 sm:gap-x-8 sm:gap-y-7 md:grid-cols-3"
            >
              <div class="min-w-0">
                <p
                  class="text-[10px] uppercase tracking-wider text-gray-500 sm:text-xs"
                >
                  Duration
                </p>

                <p
                  class="mt-2 break-words text-sm font-bold text-white sm:text-base"
                >
                  {{ movie.duration }} min
                </p>
              </div>

              <div class="min-w-0">
                <p
                  class="text-[10px] uppercase tracking-wider text-gray-500 sm:text-xs"
                >
                  Language
                </p>

                <p
                  class="mt-2 break-words text-sm font-bold text-white sm:text-base"
                >
                  {{ movie.language }}
                </p>
              </div>

              <div class="min-w-0">
                <p
                  class="text-[10px] uppercase tracking-wider text-gray-500 sm:text-xs"
                >
                  Subtitle
                </p>

                <p
                  class="mt-2 break-words text-sm font-bold text-white sm:text-base"
                >
                  {{ movie.subtitle }}
                </p>
              </div>

              <div class="min-w-0">
                <p
                  class="text-[10px] uppercase tracking-wider text-gray-500 sm:text-xs"
                >
                  Rating
                </p>

                <p
                  class="mt-2 break-words text-sm font-bold text-white sm:text-base"
                >
                  {{ movie.rating }}
                </p>
              </div>

              <div class="min-w-0">
                <p
                  class="text-[10px] uppercase tracking-wider text-gray-500 sm:text-xs"
                >
                  Release Date
                </p>

                <p
                  class="mt-2 break-words text-sm font-bold text-white sm:text-base"
                >
                  {{ formatReleaseDate(movie.release_date) }}
                </p>
              </div>

              <div class="min-w-0">
                <p
                  class="text-[10px] uppercase tracking-wider text-gray-500 sm:text-xs"
                >
                  Genre
                </p>

                <p
                  class="mt-2 break-words text-sm font-bold text-white sm:text-base"
                >
                  {{ movie.genre.join(', ') }}
                </p>
              </div>
            </div>

            <div
              class="my-6 h-px bg-white/10 sm:my-7"
            ></div>

            <!-- Actions -->
            <div
              class="flex w-full flex-col gap-3 sm:flex-row"
            >
              <!-- Watch trailer -->
              <button
                v-if="movie.trailer_url"
                type="button"
                class="flex w-full items-center justify-center gap-2 rounded-xl border border-white/10 px-5 py-3.5 text-sm font-bold text-white transition hover:bg-white/5 sm:w-auto sm:px-6 sm:py-4"
                @click="watchTrailer"
              >
                <span>▷</span>
                Watch trailer
              </button>

              <!-- Book now -->
              <button
                type="button"
                class="w-full flex-1 rounded-xl bg-[#ff3158] px-5 py-3.5 text-sm font-bold text-white shadow-lg shadow-pink-500/20 transition hover:scale-[1.01] hover:bg-[#ff1745] sm:px-6 sm:py-4"
                @click="bookNow"
              >
                Book now
              </button>
            </div>

          </div>
        </div>
      </div>

    </div>
  </div>
</template>