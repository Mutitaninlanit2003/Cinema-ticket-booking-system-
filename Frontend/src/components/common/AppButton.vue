<script setup lang="ts">
import { Loader2 } from '@lucide/vue'

type Variant = 'primary' | 'secondary' | 'ghost' | 'danger'
type Size = 'sm' | 'md' | 'lg'

withDefaults(
  defineProps<{
    variant?: Variant
    size?: Size
    loading?: boolean
    disabled?: boolean
    fullWidth?: boolean
    /** Native button "type" attribute — defaults to "button" so a
     * button inside a form never accidentally submits it. */
    type?: 'button' | 'submit' | 'reset'
  }>(),
  {
    variant: 'primary',
    size: 'md',
    loading: false,
    disabled: false,
    fullWidth: false,
    type: 'button',
  },
)

defineEmits<{ click: [MouseEvent] }>()

const variantClasses: Record<Variant, string> = {
  primary: 'bg-primary text-white hover:bg-primary-hover',
  secondary: 'bg-card text-text hover:bg-surface border border-border/10',
  ghost: 'bg-transparent text-text hover:bg-card',
  danger: 'bg-transparent text-danger border border-danger/40 hover:bg-danger/10',
}

const sizeClasses: Record<Size, string> = {
  sm: 'h-8 px-3 text-sm gap-1.5',
  md: 'h-10 px-4 text-sm gap-2',
  lg: 'h-12 px-6 text-base gap-2',
}
</script>

<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    :class="[
      'inline-flex items-center justify-center rounded-md font-medium transition-colors',
      'disabled:cursor-not-allowed disabled:opacity-50',
      variantClasses[variant],
      sizeClasses[size],
      fullWidth ? 'w-full' : '',
    ]"
    @click="(e) => !disabled && !loading && $emit('click', e)"
  >
    <Loader2 v-if="loading" class="size-4 animate-spin" />
    <slot v-else name="icon-left" />
    <slot />
    <slot v-if="!loading" name="icon-right" />
  </button>
</template>
