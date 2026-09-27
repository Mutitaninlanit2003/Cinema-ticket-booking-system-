<script setup lang="ts">
import { computed, useId } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue: string
    label?: string
    type?: string
    placeholder?: string
    error?: string
    hint?: string
    disabled?: boolean
    required?: boolean
  }>(),
  {
    type: 'text',
    disabled: false,
    required: false,
  },
)

defineEmits<{ 'update:modelValue': [string] }>()

const inputId = useId()
const describedBy = computed(() => {
  if (props.error) return `${inputId}-error`
  if (props.hint) return `${inputId}-hint`
  return undefined
})
</script>

<template>
  <div class="flex flex-col gap-1.5">
    <label
      v-if="label"
      :for="inputId"
      class="text-sm font-medium text-text"
    >
      {{ label }}
      <span v-if="required" class="text-primary">*</span>
    </label>

    <input
      :id="inputId"
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :aria-invalid="!!error"
      :aria-describedby="describedBy"
      class="h-10 rounded-md border bg-surface px-3 text-sm text-text placeholder:text-muted
             focus:outline-none focus-visible:ring-2 focus-visible:ring-primary
             disabled:cursor-not-allowed disabled:opacity-50"
      :class="error ? 'border-danger' : 'border-border/10'"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />

    <p v-if="error" :id="`${inputId}-error`" class="text-xs text-danger">
      {{ error }}
    </p>
    <p v-else-if="hint" :id="`${inputId}-hint`" class="text-xs text-muted">
      {{ hint }}
    </p>
  </div>
</template>
