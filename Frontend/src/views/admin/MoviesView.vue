<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Search, Film, Plus, Pencil, Trash2, X } from '@lucide/vue'
import api from '@/services/api'

interface Movie {
  movie_id: number
  title: string
  description: string
  genre: string[]
  duration: number
  language: string
  subtitle: string
  rating: string
  release_date: string
  poster_url: string
  trailer_url: string
  status: 'NOW_SHOWING' | 'COMING_SOON' | 'ENDED'
}

interface MovieForm {
  title: string
  description: string
  genre: string
  duration: number
  language: string
  subtitle: string
  rating: string
  release_date: string
  poster_url: string
  trailer_url: string
  status: string
}

const movies = ref<Movie[]>([])
const search = ref('')
const activeFilter = ref('ALL')

const showForm = ref(false)
const editingMovie = ref<Movie | null>(null)
const loading = ref(false)
const error = ref('')

const form = ref<MovieForm>({
  title: '',
  description: '',
  genre: '',
  duration: 0,
  language: '',
  subtitle: '',
  rating: '',
  release_date: '',
  poster_url: '',
  trailer_url: '',
  status: 'NOW_SHOWING',
})

const filters = [
  { label: 'All', value: 'ALL' },
  { label: 'Now showing', value: 'NOW_SHOWING' },
  { label: 'Coming soon', value: 'COMING_SOON' },
  { label: 'Ended', value: 'ENDED' },
]

const filteredMovies = computed(() => {
  return movies.value.filter((movie) => {
    const matchSearch = movie.title
      .toLowerCase()
      .includes(search.value.toLowerCase())

    const matchStatus =
      activeFilter.value === 'ALL' ||
      movie.status === activeFilter.value

    return matchSearch && matchStatus
  })
})

async function fetchMovies() {
  try {
    loading.value = true
    error.value = ''

    const response = await api.get('/movies')

    const data = response.data

    movies.value = Array.isArray(data)
      ? data
      : Array.isArray(data?.data)
        ? data.data
        : []
  } catch (err) {
    console.error(err)
    error.value = 'ไม่สามารถโหลดข้อมูลภาพยนตร์ได้'
  } finally {
    loading.value = false
  }
}

function openAddForm() {
  editingMovie.value = null

  form.value = {
    title: '',
    description: '',
    genre: '',
    duration: 0,
    language: '',
    subtitle: '',
    rating: '',
    release_date: '',
    poster_url: '',
    trailer_url: '',
    status: 'NOW_SHOWING',
  }

  showForm.value = true
}

function openEditForm(movie: Movie) {
  editingMovie.value = movie

  form.value = {
    title: movie.title,
    description: movie.description,
    genre: movie.genre.join(', '),
    duration: movie.duration,
    language: movie.language,
    subtitle: movie.subtitle,
    rating: movie.rating,
    release_date: movie.release_date
      ? movie.release_date.substring(0, 10)
      : '',
    poster_url: movie.poster_url,
    trailer_url: movie.trailer_url,
    status: movie.status,
  }

  showForm.value = true
}

function closeForm() {
  showForm.value = false
  editingMovie.value = null
}

async function saveMovie() {
  try {
    error.value = ''

    if (!form.value.title.trim()) {
      error.value = 'กรุณากรอกชื่อภาพยนตร์'
      return
    }

    if (!form.value.genre.trim()) {
      error.value = 'กรุณากรอกประเภทภาพยนตร์'
      return
    }

    if (!form.value.release_date) {
      error.value = 'กรุณาเลือกวันเข้าฉาย'
      return
    }

    // แปลงข้อมูลจาก Form ให้ตรงกับ Backend
    const payload = {
      title: form.value.title,
      description: form.value.description,

      // "Action, Adventure, Sci-Fi"
      // ↓
      // ["Action", "Adventure", "Sci-Fi"]
      genre: form.value.genre
        .split(',')
        .map((item) => item.trim())
        .filter(Boolean),

      duration: Number(form.value.duration),
      language: form.value.language,
      subtitle: form.value.subtitle,
      rating: form.value.rating,

      // input type="date"
      // ↓
      // รูปแบบที่ time.Time อ่านได้
      release_date: `${form.value.release_date}T00:00:00Z`,

      poster_url: form.value.poster_url,
      trailer_url: form.value.trailer_url,
      status: form.value.status,
    }

    console.log('Sending movie:', payload)

    if (editingMovie.value) {
      await api.put(
        `/movies/${editingMovie.value.movie_id}`,
        payload,
      )
    } else {
      await api.post(
        '/movies',
        payload,
      )
    }

    await fetchMovies()
    closeForm()

  } catch (err: any) {
    console.error('Save movie error:', err)

    console.error(
      'Backend response:',
      err.response?.data,
    )

    error.value =
      err.response?.data?.error ||
      'ไม่สามารถบันทึกข้อมูลภาพยนตร์ได้'
  }
}

async function deleteMovie(movie: Movie) {
  const confirmed = window.confirm(
    `ต้องการลบ "${movie.title}" ใช่หรือไม่?`,
  )

  if (!confirmed) return

  try {
    await api.delete(`/movies/${movie.movie_id}`)
    await fetchMovies()
  } catch (err) {
    console.error(err)
    error.value = 'ไม่สามารถลบภาพยนตร์ได้'
  }
}

function statusLabel(status: string) {
  switch (status) {
    case 'NOW_SHOWING':
      return 'NOW SHOWING'
    case 'COMING_SOON':
      return 'COMING SOON'
    case 'ENDED':
      return 'ENDED'
    default:
      return status
  }
}

function statusClass(status: string) {
  switch (status) {
    case 'NOW_SHOWING':
      return 'bg-emerald-500/15 text-emerald-400'

    case 'COMING_SOON':
      return 'bg-orange-500/15 text-orange-400'

    case 'ENDED':
      return 'bg-white/10 text-gray-400'

    default:
      return 'bg-white/10 text-gray-400'
  }
}

onMounted(fetchMovies)
</script>

<template>
  <div class="min-h-screen bg-[#0d0d12] px-6 py-8 lg:px-8">

    <!-- Header -->
    <div class="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
      <div class="flex items-center gap-4">
        
        <!-- Movie Icon -->
        <div class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-gradient-to-br from-purple-500 to-pink-500 shadow-lg shadow-purple-500/20" 
          > 
          <Film class="size-6 text-white" /> </div>

        <!-- Title -->
        <div>
          <h1 class="text-2xl font-bold tracking-tight text-text">
            Movie Management
          </h1>

          <p class="mt-1 text-sm text-muted">
            Manage the film
          </p>
        </div>

      </div>
      <button
        type="button"
        class="inline-flex items-center justify-center gap-2 rounded-lg bg-[#ff3158] px-4 py-2.5 text-sm font-semibold text-white shadow-lg shadow-[#ff3158]/20 transition hover:bg-[#ff4669]"
        @click="openAddForm"
      >
        <Plus class="size-4" />
        Add movie
      </button>

    </div>


    <!-- Search + Filters -->
    <div class="mt-7 flex flex-col gap-3 xl:flex-row xl:items-center xl:justify-between">

      <div class="relative w-full xl:max-w-xs">

        <Search
          class="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-[#68758f]"
        />

        <input
          v-model="search"
          type="text"
          placeholder="Search by title..."
          class="w-full rounded-lg border border-white/10 bg-[#15151d] py-2.5 pl-9 pr-4 text-sm text-white outline-none transition placeholder:text-[#68758f] focus:border-[#ff3158]/60"
        />

      </div>


      <div class="flex flex-wrap gap-2">

        <button
          v-for="filter in filters"
          :key="filter.value"
          type="button"
          class="rounded-lg border px-4 py-2 text-xs font-semibold transition"
          :class="
            activeFilter === filter.value
              ? 'border-[#ff3158] bg-[#ff3158] text-white'
              : 'border-white/10 bg-[#15151d] text-[#a0abc0] hover:border-white/20 hover:text-white'
          "
          @click="activeFilter = filter.value"
        >
          {{ filter.label }}
        </button>

      </div>

    </div>


    <!-- Error -->
    <div
      v-if="error"
      class="mt-5 rounded-lg border border-red-500/20 bg-red-500/10 px-4 py-3 text-sm text-red-400"
    >
      {{ error }}
    </div>


    <!-- Movie Table -->
    <div class="mt-5 overflow-hidden rounded-xl border border-white/10 bg-[#15151d]">

      <div class="overflow-x-auto">

        <table class="w-full min-w-[1000px]">

          <thead>
            <tr class="border-b border-white/10 bg-[#1b1b24]">

              <th class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]">
                TITLE
              </th>

              <th class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]">
                GENRE
              </th>

              <th class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]">
                DURATION
              </th>

              <th class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]">
                LANGUAGE
              </th>

              <th class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]">
                RATING
              </th>

              <th class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]">
                RELEASE DATE
              </th>

              <th class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]">
                STATUS
              </th>

              <th class="px-5 py-4 text-right text-xs font-medium text-[#8190ad]">
                ACTION
              </th>

            </tr>
          </thead>


          <tbody>

            <tr
              v-if="loading"
              class="border-b border-white/10"
            >
              <td
                colspan="8"
                class="px-5 py-10 text-center text-sm text-[#8190ad]"
              >
                Loading movies...
              </td>
            </tr>


            <tr
              v-else-if="filteredMovies.length === 0"
              class="border-b border-white/10"
            >
              <td
                colspan="8"
                class="px-5 py-10 text-center text-sm text-[#8190ad]"
              >
                No movies found.
              </td>
            </tr>


            <tr
              v-for="movie in filteredMovies"
              :key="movie.movie_id"
              class="border-b border-white/10 transition hover:bg-white/[0.02]"
            >

              <td class="px-5 py-4">

                <div class="flex items-center gap-3">

                  <img
                    v-if="movie.poster_url"
                    :src="movie.poster_url"
                    :alt="movie.title"
                    class="h-12 w-9 rounded object-cover"
                  />

                  <div>
                    <p class="text-sm font-semibold text-white">
                      {{ movie.title }}
                    </p>
                  </div>

                </div>

              </td>


              <td class="px-5 py-4 text-sm text-gray-300">
                {{ movie.genre || '-' }}
              </td>


              <td class="px-5 py-4 text-sm text-gray-300">
                {{ movie.duration || '-' }} min
              </td>


              <td class="px-5 py-4 text-sm text-gray-300">
                {{ movie.language || '-' }}
              </td>


              <td class="px-5 py-4 text-sm text-gray-300">
                {{ movie.rating || '-' }}
              </td>


              <td class="px-5 py-4 text-sm text-gray-300">
                {{
                  movie.release_date
                    ? new Date(movie.release_date).toLocaleDateString('en-US', {
                        month: 'short',
                        day: 'numeric',
                        year: 'numeric',
                      })
                    : '-'
                }}
              </td>


              <td class="px-5 py-4">

                <span
                  class="inline-flex rounded-full px-3 py-1 text-[11px] font-semibold"
                  :class="statusClass(movie.status)"
                >
                  {{ statusLabel(movie.status) }}
                </span>

              </td>


              <td class="px-5 py-4">

                <div class="flex justify-end gap-2">

                  <button
                    type="button"
                    class="flex size-8 items-center justify-center rounded-md border border-white/10 text-[#9aa6bd] transition hover:border-white/20 hover:bg-white/5 hover:text-white"
                    title="Edit movie"
                    @click="openEditForm(movie)"
                  >
                    <Pencil class="size-4" />
                  </button>

                  <button
                    type="button"
                    class="flex size-8 items-center justify-center rounded-md border border-white/10 text-[#9aa6bd] transition hover:border-red-500/30 hover:bg-red-500/10 hover:text-red-400"
                    title="Delete movie"
                    @click="deleteMovie(movie)"
                  >
                    <Trash2 class="size-4" />
                  </button>

                </div>

              </td>

            </tr>

          </tbody>

        </table>

      </div>

    </div>


    <!-- Add / Edit Movie Form -->
    <div
      v-if="showForm"
      class="mt-6 rounded-xl border border-white/10 bg-[#15151d] p-6 shadow-xl"
    >

      <div class="mb-6 flex items-center justify-between">

        <div>
          <h2 class="text-lg font-bold text-white">
            {{ editingMovie ? 'Edit Movie' : 'Add Movie' }}
          </h2>

          <p class="mt-1 text-xs text-[#8190ad]">
            {{
              editingMovie
                ? 'แก้ไขข้อมูลภาพยนตร์'
                : 'เพิ่มภาพยนตร์เข้าสู่ระบบ'
            }}
          </p>
        </div>

        <button
          type="button"
          class="flex size-8 items-center justify-center rounded-md text-[#8190ad] hover:bg-white/5 hover:text-white"
          @click="closeForm"
        >
          <X class="size-5" />
        </button>

      </div>


      <div class="grid grid-cols-1 gap-5 md:grid-cols-2">

        <!-- Title -->
        <div class="md:col-span-2">

          <label class="mb-2 block text-xs font-medium text-[#8190ad]">
            Title
          </label>

          <input
            v-model="form.title"
            type="text"
            placeholder="Movie title"
            class="w-full rounded-lg border border-white/10 bg-[#101016] px-4 py-2.5 text-sm text-white outline-none focus:border-[#ff3158]/60"
          />

        </div>


        <!-- Description -->
        <div class="md:col-span-2">

          <label class="mb-2 block text-xs font-medium text-[#8190ad]">
            Description
          </label>

          <textarea
            v-model="form.description"
            rows="3"
            placeholder="Movie description"
            class="w-full resize-none rounded-lg border border-white/10 bg-[#101016] px-4 py-2.5 text-sm text-white outline-none focus:border-[#ff3158]/60"
          />

        </div>


        <!-- Genre -->
        <div>

          <label class="mb-2 block text-xs font-medium text-[#8190ad]">
            Genre
          </label>

          <input
            v-model="form.genre"
            type="text"
            placeholder="Sci-Fi, Action"
            class="w-full rounded-lg border border-white/10 bg-[#101016] px-4 py-2.5 text-sm text-white outline-none focus:border-[#ff3158]/60"
          />

        </div>


        <!-- Duration -->
        <div>

          <label class="mb-2 block text-xs font-medium text-[#8190ad]">
            Duration (min)
          </label>

          <input
            v-model.number="form.duration"
            type="number"
            min="0"
            placeholder="120"
            class="w-full rounded-lg border border-white/10 bg-[#101016] px-4 py-2.5 text-sm text-white outline-none focus:border-[#ff3158]/60"
          />

        </div>


        <!-- Language -->
        <div>

          <label class="mb-2 block text-xs font-medium text-[#8190ad]">
            Language
          </label>

          <input
            v-model="form.language"
            type="text"
            placeholder="English"
            class="w-full rounded-lg border border-white/10 bg-[#101016] px-4 py-2.5 text-sm text-white outline-none focus:border-[#ff3158]/60"
          />

        </div>


        <!-- Subtitle -->
        <div>

          <label class="mb-2 block text-xs font-medium text-[#8190ad]">
            Subtitle
          </label>

          <input
            v-model="form.subtitle"
            type="text"
            placeholder="Thai"
            class="w-full rounded-lg border border-white/10 bg-[#101016] px-4 py-2.5 text-sm text-white outline-none focus:border-[#ff3158]/60"
          />

        </div>


        <!-- Rating -->
        <div>

          <label class="mb-2 block text-xs font-medium text-[#8190ad]">
            Rating
          </label>

          <input
            v-model="form.rating"
            type="text"
            placeholder="PG-13"
            class="w-full rounded-lg border border-white/10 bg-[#101016] px-4 py-2.5 text-sm text-white outline-none focus:border-[#ff3158]/60"
          />

        </div>


        <!-- Release date -->
        <div>

          <label class="mb-2 block text-xs font-medium text-[#8190ad]">
            Release date
          </label>

          <input
            v-model="form.release_date"
            type="date"
            class="w-full rounded-lg border border-white/10 bg-[#101016] px-4 py-2.5 text-sm text-white outline-none focus:border-[#ff3158]/60"
          />

        </div>


        <!-- Poster -->
        <div>

          <label class="mb-2 block text-xs font-medium text-[#8190ad]">
            Poster URL
          </label>

          <input
            v-model="form.poster_url"
            type="text"
            placeholder="https://..."
            class="w-full rounded-lg border border-white/10 bg-[#101016] px-4 py-2.5 text-sm text-white outline-none focus:border-[#ff3158]/60"
          />

        </div>


        <!-- Trailer -->
        <div>

          <label class="mb-2 block text-xs font-medium text-[#8190ad]">
            Trailer URL
          </label>

          <input
            v-model="form.trailer_url"
            type="text"
            placeholder="https://..."
            class="w-full rounded-lg border border-white/10 bg-[#101016] px-4 py-2.5 text-sm text-white outline-none focus:border-[#ff3158]/60"
          />

        </div>


        <!-- Status -->
        <div class="md:col-span-2">

          <label class="mb-2 block text-xs font-medium text-[#8190ad]">
            Status
          </label>

          <select
            v-model="form.status"
            class="w-full rounded-lg border border-white/10 bg-[#101016] px-4 py-2.5 text-sm text-white outline-none focus:border-[#ff3158]/60"
          >
            <option value="NOW_SHOWING">
              NOW_SHOWING
            </option>

            <option value="COMING_SOON">
              COMING_SOON
            </option>

            <option value="ENDED">
              ENDED
            </option>
          </select>

        </div>

      </div>


      <!-- Form buttons -->
      <div class="mt-6 flex justify-end gap-3">

        <button
          type="button"
          class="rounded-lg border border-white/10 px-4 py-2 text-sm font-semibold text-gray-300 transition hover:bg-white/5 hover:text-white"
          @click="closeForm"
        >
          Cancel
        </button>

        <button
          type="button"
          class="rounded-lg bg-[#ff3158] px-5 py-2 text-sm font-semibold text-white transition hover:bg-[#ff4669]"
          @click="saveMovie"
        >
          {{ editingMovie ? 'Save changes' : 'Add movie' }}
        </button>

      </div>

    </div>

  </div>
</template>