<script setup lang="ts">
import { onMounted, onUnmounted, watch } from 'vue'
import { X } from '@lucide/vue'

const props = withDefaults(
  defineProps<{
    modelValue: boolean
    title?: string
    /** Clicking the backdrop closes the modal — disable for
     * destructive confirmations where an accidental click shouldn't dismiss. */
    closeOnBackdrop?: boolean
  }>(),
  {
    closeOnBackdrop: true,
  },
)

const emit = defineEmits<{ 'update:modelValue': [boolean] }>()

function close() {
  emit('update:modelValue', false)
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.modelValue) close()
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onUnmounted(() => window.removeEventListener('keydown', onKeydown))

// Lock page scroll while the modal is open
watch(
  () => props.modelValue,
  (open) => {
    document.body.style.overflow = open ? 'hidden' : ''
  },
)
</script>

<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        role="dialog"
        aria-modal="true"
        :aria-label="title"
      >
        <div
          class="absolute inset-0 bg-black/70"
          @click="closeOnBackdrop && close()"
        />

        <div
          class="relative z-10 w-full max-w-md rounded-lg border border-border/10 bg-card p-6 shadow-2xl"
        >
          <div v-if="title || $slots.header" class="mb-4 flex items-start justify-between gap-4">
            <slot name="header">
              <h2 class="text-lg font-semibold text-text">{{ title }}</h2>
            </slot>
            <button
              type="button"
              class="text-muted transition-colors hover:text-text"
              aria-label="Close"
              @click="close"
            >
              <X class="size-5" />
            </button>
          </div>

          <div class="text-sm text-muted">
            <slot />
          </div>

          <div v-if="$slots.footer" class="mt-6 flex justify-end gap-3">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.15s ease;
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
</style>
