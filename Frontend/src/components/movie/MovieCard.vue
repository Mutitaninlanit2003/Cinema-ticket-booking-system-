<script setup lang="ts">
interface Movie {
  id: number
  title: string
  genre: string
  rating?: number
  releaseDate?: string
  status: 'NOW_SHOWING' | 'COMING_SOON'
  gradient: string
}

defineProps<{
  movie: Movie
}>()
</script>

<template>
  <div
    class="overflow-hidden rounded-xl border border-white/10 bg-[#15151d] transition duration-300 hover:-translate-y-1 hover:border-white/20"
  >
    <!-- Poster -->
    <div
      :class="movie.gradient"
      class="relative h-52"
    >
      <!-- Coming Soon Badge -->
      <span
        v-if="movie.status === 'COMING_SOON'"
        class="absolute left-0 top-3 rounded-r-md bg-violet-500 px-2.5 py-1 text-xs font-bold"
      >
        COMING SOON
      </span>

      <!-- Movie Title -->
      <div class="absolute bottom-4 left-4 right-4">
        <h3
          class="max-w-[180px] text-lg font-black uppercase leading-tight"
        >
          {{ movie.title }}
        </h3>
      </div>
    </div>

    <!-- Movie Info -->
    <div class="p-3.5">

      <!-- NOW SHOWING -->
      <template v-if="movie.status === 'NOW_SHOWING'">
        <h3 class="truncate font-semibold">
          {{ movie.title }}
        </h3>

        <div class="mt-1 flex items-center justify-between">
          <span class="text-sm text-gray-500">
            {{ movie.genre }}
          </span>

          <span
            v-if="movie.rating"
            class="flex items-center gap-1 text-sm font-semibold text-yellow-400"
          >
            ★ {{ movie.rating }}
          </span>
        </div>

        <button
          class="mt-3 w-full rounded-lg border border-white/20 py-2 text-sm font-semibold transition hover:bg-white hover:text-black"
        >
          Book now
        </button>
      </template>

      <!-- COMING SOON -->
      <template v-else>
        <div class="flex items-center justify-between">
          <h3 class="truncate font-semibold">
            {{ movie.title }}
          </h3>

          <span class="ml-2 whitespace-nowrap text-xs text-gray-500">
            {{ movie.releaseDate }}
          </span>
        </div>

        <p class="mt-1 text-sm text-gray-500">
          {{ movie.genre }}
        </p>
      </template>

    </div>
  </div>
</template>