<script setup lang="ts">
// Generic status pill. Deliberately generic ("success" / "warning" / ...)
// rather than domain-specific ("confirmed" / "expired") — a future
// statusToVariant() mapper (added alongside the real Booking/Payment/
// Seat types) translates domain status -> variant, keeping this
// component free of business logic.
type Variant = 'success' | 'warning' | 'danger' | 'info' | 'neutral'

withDefaults(
  defineProps<{
    variant?: Variant
  }>(),
  {
    variant: 'neutral',
  },
)

const variantClasses: Record<Variant, string> = {
  success: 'bg-success/15 text-success',
  warning: 'bg-warning/15 text-warning',
  danger: 'bg-danger/15 text-danger',
  info: 'bg-primary/15 text-primary',
  neutral: 'bg-muted/15 text-muted',
}
</script>

<template>
  <span
    :class="[
      'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium',
      variantClasses[variant],
    ]"
  >
    <slot />
  </span>
</template>
