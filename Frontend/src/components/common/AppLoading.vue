<script setup lang="ts">
import { Loader2 } from '@lucide/vue'

withDefaults(
  defineProps<{
    /** Renders as a full-viewport overlay instead of an inline spinner —
     * for page-level loads (e.g. waiting on GET /movies). */
    fullscreen?: boolean
    label?: string
    size?: 'sm' | 'md' | 'lg'
  }>(),
  {
    fullscreen: false,
    size: 'md',
  },
)

const sizeClasses = {
  sm: 'size-4',
  md: 'size-6',
  lg: 'size-10',
}
</script>

<template>
  <div
    v-if="fullscreen"
    class="fixed inset-0 z-40 flex flex-col items-center justify-center gap-3 bg-bg/80"
  >
    <Loader2 class="animate-spin text-primary" :class="sizeClasses[size]" />
    <p v-if="label" class="text-sm text-muted">{{ label }}</p>
  </div>

  <div v-else class="inline-flex items-center gap-2">
    <Loader2 class="animate-spin text-primary" :class="sizeClasses[size]" />
    <span v-if="label" class="text-sm text-muted">{{ label }}</span>
  </div>
</template>
