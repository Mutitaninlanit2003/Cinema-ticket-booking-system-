<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import {
  Plus,
  Search,
  Pencil,
  Trash2,
  X,
  Building2,
} from 'lucide-vue-next'

interface Theater {
  theater_id?: string
  id?: number
  name: string
  address: string
  location: string
  capacity: number
  status: string
}

const theaters = ref<Theater[]>([])
const search = ref('')
const loading = ref(false)
const error = ref('')

const showForm = ref(false)
const editingTheater = ref<Theater | null>(null)
const nameInput = ref<HTMLInputElement | null>(null)

const form = ref({
  name: '',
  address: '',
  location: '',
  capacity: 0,
  status: 'active',
})

const filteredTheaters = computed(() => {
  const keyword = search.value.toLowerCase().trim()

  if (!keyword) {
    return theaters.value
  }

  return theaters.value.filter((theater) =>
    [
      theater.name,
      theater.address,
      theater.location,
      theater.status,
    ]
      .join(' ')
      .toLowerCase()
      .includes(keyword)
  )
})

async function fetchTheaters() {
  try {
    loading.value = true
    error.value = ''

    const response = await fetch('http://localhost:8080/theaters')

    if (!response.ok) {
      throw new Error('ไม่สามารถโหลดข้อมูลโรงภาพยนตร์ได้')
    }

    const data = await response.json()

    theaters.value = Array.isArray(data)
      ? data
      : data.data ?? []
  } catch (err) {
    error.value =
      err instanceof Error
        ? err.message
        : 'เกิดข้อผิดพลาดในการโหลดข้อมูล'
  } finally {
    loading.value = false
  }
}

function resetForm() {
  form.value = {
    name: '',
    address: '',
    location: '',
    capacity: 0,
    status: 'active',
  }

  editingTheater.value = null
}

function openAddForm() {
  resetForm()
  showForm.value = true
}

function openEditForm(theater: Theater) {
  editingTheater.value = theater

  form.value = {
    name: theater.name,
    address: theater.address,
    location: theater.location,
    capacity: theater.capacity,
    status: theater.status,
  }

  showForm.value = true
}

function closeForm() {
  showForm.value = false
  resetForm()
}

async function saveTheater() {
  try {
    error.value = ''

    const isEditing = !!editingTheater.value

    const theaterId =
      editingTheater.value?.theater_id ??
      editingTheater.value?.id

    const url = isEditing
      ? `http://localhost:8080/theaters/${theaterId}`
      : 'http://localhost:8080/theaters'

    const response = await fetch(url, {
      method: isEditing ? 'PUT' : 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: form.value.name,
        address: form.value.address,
        location: form.value.location,
        capacity: Number(form.value.capacity),
        status: form.value.status,
      }),
    })

    if (!response.ok) {
      throw new Error(
        isEditing
          ? 'ไม่สามารถแก้ไขโรงภาพยนตร์ได้'
          : 'ไม่สามารถเพิ่มโรงภาพยนตร์ได้'
      )
    }

    await fetchTheaters()
    closeForm()
  } catch (err) {
    error.value =
      err instanceof Error
        ? err.message
        : 'เกิดข้อผิดพลาด'
  }
}

async function deleteTheater(theater: Theater) {
  const theaterId = theater.theater_id

  if (!theaterId) {
    error.value = 'ไม่พบ Theater UUID'
    return
  }

  const confirmed = window.confirm(
    `ต้องการลบ "${theater.name}" หรือไม่?`
  )

  if (!confirmed) return

  try {
    error.value = ''

    const response = await fetch(
      `http://localhost:8080/theaters/${theaterId}`,
      {
        method: 'DELETE',
      }
    )

    if (!response.ok) {
      const data = await response.json().catch(() => null)

      throw new Error(
        data?.error ||
        data?.message ||
        'ไม่สามารถลบโรงภาพยนตร์ได้'
      )
    }

    await fetchTheaters()

  } catch (err) {
    error.value =
      err instanceof Error
        ? err.message
        : 'เกิดข้อผิดพลาดในการลบข้อมูล'
  }
}

function getStatusClass(status: string) {
  switch (status?.toLowerCase()) {
    case 'active':
      return 'bg-emerald-500/10 text-emerald-400'

    case 'inactive':
      return 'bg-gray-500/10 text-gray-400'

    default:
      return 'bg-yellow-500/10 text-yellow-400'
  }
}

// -----------------------------------------------------------------------
// Modal behaviour: Esc closes it, page behind stops scrolling, and the
// name field is focused as soon as it opens. (Matches ScreensView.vue)
// -----------------------------------------------------------------------
function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape' && showForm.value) {
    closeForm()
  }
}

watch(showForm, async (open) => {
  document.body.style.overflow = open ? 'hidden' : ''

  if (open) {
    await nextTick()
    nameInput.value?.focus()
  }
})

onMounted(() => {
  fetchTheaters()
  window.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
  document.body.style.overflow = ''
})
</script>

<template>
  <div class="min-h-screen bg-[#0d0d12] px-6 py-8 lg:px-8">

    <!-- Page Header -->
    <div class="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
      <div class="flex items-center gap-4">
        <!-- Theater Icon -->
        <div
          class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-gradient-to-br from-purple-500 to-pink-500 shadow-lg shadow-purple-500/20"
        >
          <Building2 class="size-6" />
        </div>

        <div>
          <h1 class="text-2xl font-bold text-white">
            Theater Management
          </h1>

          <p class="mt-1 text-sm text-gray-400">
            Manage a movie theater
          </p>
        </div>
      </div>
    </div>


    <!-- Main Content -->
    <div class="mx-auto w-full max-w-[1200px] px-4 pb-10 lg:px-8">

      <!-- Error -->
      <div
        v-if="error"
        class="mb-5 rounded-lg border border-red-500/20 bg-red-500/10 px-4 py-3 text-sm text-red-400"
      >
        {{ error }}
      </div>

      <!-- Search -->
      <div class="mt-7 flex flex-col gap-3 xl:flex-row xl:items-center xl:justify-between">
        <div class="relative w-full xl:max-w-xs">
          <Search
            class="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-[#68758f]"
          />

          <input
            v-model="search"
            type="text"
            placeholder="Search theaters..."
            class="w-full rounded-lg border border-white/10 bg-[#15151d] py-2.5 pl-9 pr-4 text-sm text-white outline-none transition placeholder:text-[#68758f] focus:border-[#ff3158]/60"
          />
        </div>
         <!-- Add Theater -->
        <button
          type="button"
          class="inline-flex items-center justify-center gap-2 rounded-lg bg-[#ff3158] px-4 py-2.5 text-sm font-semibold text-white shadow-lg shadow-[#ff3158]/20 transition hover:bg-[#ff4669]"
          @click="openAddForm"
        >
          <Plus class="size-4" />
          Add theater
        </button>
      </div>

      <!-- Theater Table -->
      <div class="mt-5 overflow-hidden rounded-xl border border-white/10 bg-[#15151d]">

        <div class="overflow-x-auto">

          <table class="w-full min-w-[1000px]">
            <thead>
              <tr class="border-b border-white/10 bg-[#1b1b24]">

                <th class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]">
                  Name
                </th>

                <th class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]">
                  Address
                </th>

                <th class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]">
                  Location
                </th>

                <th class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]">
                  Capacity
                </th>

                <th class="px-5 py-4 text-left text-xs font-medium text-[#8190ad]">
                  Status
                </th>

                <th class="px-5 py-4 text-right text-xs font-medium uppercase tracking-wide text-[#8190ad]">
                  Actions
                </th>
              </tr>
            </thead>

            <tbody>
              <!-- Loading -->
              <tr v-if="loading">
                <td
                  colspan="6"
                  class="px-5 py-10 text-center text-sm text-muted"
                >
                  Loading theaters...
                </td>
              </tr>

              <!-- Empty -->
              <tr v-else-if="filteredTheaters.length === 0">
                <td
                  colspan="6"
                  class="px-5 py-12 text-center"
                >
                  <p class="text-sm text-muted">
                    No theaters found.
                  </p>
                </td>
              </tr>

              <!-- Data -->
              <tr
                v-for="theater in filteredTheaters"
                :key="theater.theater_id ?? theater.id"
                class="border-b border-border/10 last:border-0 transition hover:bg-white/[0.02]"
              >
                <td class="px-5 py-4">
                  <span class="font-semibold text-text">
                    {{ theater.name }}
                  </span>
                </td>

                <td class="px-5 py-4 text-muted">
                  {{ theater.address || '-' }}
                </td>

                <td class="px-5 py-4 text-muted">
                  {{ theater.location || '-' }}
                </td>

                <td class="px-5 py-4 text-muted">
                  {{ theater.capacity ?? 0 }}
                </td>

                <td class="px-5 py-4">
                  <span
                    class="inline-flex rounded-full px-3 py-1 text-xs font-medium"
                    :class="getStatusClass(theater.status)"
                  >
                    {{ theater.status }}
                  </span>
                </td>

                <td class="px-5 py-4">
                  <div class="flex justify-end gap-2">
                    <button
                      type="button"
                      class="rounded-md border border-border/10 p-2 text-muted transition hover:border-primary/30 hover:bg-primary/10 hover:text-primary"
                      title="Edit"
                      @click="openEditForm(theater)"
                    >
                      <Pencil class="size-4" />
                    </button>

                    <button
                      type="button"
                      class="rounded-md border border-red-500/20 p-2 text-red-400 transition hover:bg-red-500/10"
                      title="Delete"
                      @click="deleteTheater(theater)"
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

      <!-- Add / Edit Theater modal (design matches ScreensView.vue) -->
      <Transition name="modal">
        <div
          v-if="showForm"
          class="fixed inset-0 z-[100] flex items-start justify-center overflow-y-auto bg-black/70 p-4 backdrop-blur-sm sm:items-center"
          @click.self="closeForm"
        >
          <div
            class="modal-card my-auto w-full max-w-2xl rounded-2xl border border-border/20 bg-card p-6 shadow-2xl shadow-black/50 sm:p-7"
            role="dialog"
            aria-modal="true"
            aria-labelledby="theater-modal-title"
          >
            <!-- Modal Header -->
            <div class="mb-8 flex items-start justify-between gap-4">
              <div>
                <h2
                  id="theater-modal-title"
                  class="text-2xl font-bold tracking-tight text-text"
                >
                  {{ editingTheater ? 'EDIT THEATER' : 'ADD THEATER' }}
                </h2>

                <p class="mt-1 text-sm text-muted">
                  {{
                    editingTheater
                      ? 'แก้ไขข้อมูลโรงภาพยนตร์'
                      : 'เพิ่มโรงภาพยนตร์ใหม่เข้าสู่ระบบ'
                  }}
                </p>
              </div>

              <button
                type="button"
                class="rounded-lg p-2 text-muted transition hover:bg-white/5 hover:text-text focus-visible:outline focus-visible:outline-2 focus-visible:outline-primary"
                aria-label="Close"
                @click="closeForm"
              >
                <X class="size-5" />
              </button>
            </div>

            <!-- Form -->
            <div class="grid gap-x-5 gap-y-6 sm:grid-cols-2">

              <!-- Name -->
              <div class="sm:col-span-2">
                <label class="mb-2 block text-sm font-medium text-muted">
                  Name
                </label>

                <input
                  ref="nameInput"
                  v-model="form.name"
                  type="text"
                  placeholder="e.g. Cinema Complex B"
                  class="h-14 w-full rounded-xl border border-border/20 bg-surface px-4 text-sm text-text outline-none transition placeholder:text-muted focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
                  @keyup.enter="saveTheater"
                />
              </div>

              <!-- Address -->
              <div class="sm:col-span-2">
                <label class="mb-2 block text-sm font-medium text-muted">
                  Address
                </label>

                <input
                  v-model="form.address"
                  type="text"
                  placeholder="Street address"
                  class="h-14 w-full rounded-xl border border-border/20 bg-surface px-4 text-sm text-text outline-none transition placeholder:text-muted focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
                  @keyup.enter="saveTheater"
                />
              </div>

              <!-- Location -->
              <div>
                <label class="mb-2 block text-sm font-medium text-muted">
                  Location
                </label>

                <input
                  v-model="form.location"
                  type="text"
                  placeholder="Area / district"
                  class="h-14 w-full rounded-xl border border-border/20 bg-surface px-4 text-sm text-text outline-none transition placeholder:text-muted focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
                  @keyup.enter="saveTheater"
                />
              </div>

              <!-- Capacity -->
              <div>
                <label class="mb-2 block text-sm font-medium text-muted">
                  Capacity
                </label>

                <input
                  v-model.number="form.capacity"
                  type="number"
                  min="0"
                  placeholder="e.g. 400"
                  class="h-14 w-full rounded-xl border border-border/20 bg-surface px-4 text-sm text-text outline-none transition placeholder:text-muted focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
                  @keyup.enter="saveTheater"
                />
              </div>

              <!-- Status -->
              <div class="sm:col-span-2">
                <label class="mb-2 block text-sm font-medium text-muted">
                  Status
                </label>

                <select
                  v-model="form.status"
                  class="h-14 w-full rounded-xl border border-border/20 bg-surface px-4 text-sm text-text outline-none transition focus:border-primary/60 focus:ring-2 focus:ring-primary/10"
                >
                  <option value="active">
                    Active
                  </option>

                  <option value="inactive">
                    Inactive
                  </option>
                </select>
              </div>

            </div>

            <!-- Buttons -->
            <div class="mt-8 flex justify-end gap-3">
              <button
                type="button"
                class="rounded-xl border border-border/20 px-6 py-3 text-sm font-semibold text-muted transition hover:bg-white/5 hover:text-text"
                @click="closeForm"
              >
                Cancel
              </button>

              <button
                type="button"
                class="rounded-lg bg-[#ff3158] px-5 py-2 text-sm font-semibold text-white transition hover:bg-[#ff4669]"
                @click="saveTheater"
              >
                {{ editingTheater ? 'Save changes' : 'Save theater' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>

    </div>
  </div>
</template>

<style scoped>
/* Tailwind classes can't reach the popup list a native <select>
   renders, since that's drawn by the browser itself — only real CSS
   on <option> (and <select>) can change its color. Supported in
   Chrome, Edge and Firefox; some browsers (e.g. Safari on macOS)
   ignore this and keep the OS-native white list regardless. */
select {
  color-scheme: dark;
}

select option {
  background-color: #15151d;
  color: #e5e7eb;
}

/* Modal open/close: a short fade for the backdrop and a small rise for
   the card. This only plays in response to the person's click.
   (Matches ScreensView.vue) */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.18s ease;
}

.modal-enter-active .modal-card,
.modal-leave-active .modal-card {
  transition: transform 0.18s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-card,
.modal-leave-to .modal-card {
  transform: translateY(12px) scale(0.98);
}

@media (prefers-reduced-motion: reduce) {
  .modal-enter-active,
  .modal-leave-active,
  .modal-enter-active .modal-card,
  .modal-leave-active .modal-card {
    transition: none;
  }
}
</style>