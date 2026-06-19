<script setup lang="ts">
import { ref } from "vue"

const {} = defineProps<{
  type?: "text" | "number"
  icon?: string
}>()

const model = defineModel()
const input = ref()

const focus = () => {
  input.value?.focus()
}
const blur = () => {
  input.value?.blur()
}

defineExpose({ focus, blur })

const handleValidate = () => {}
</script>

<template>
  <div class="input">
    <span v-if="icon" class="material-symbols-outlined">{{ icon }}</span>
    <input
      ref="input"
      v-model="model"
      v-bind="$attrs"
      :type="type"
      @keydown.enter="handleValidate()"
      @keydown.esc="handleValidate()"
      @blur.esc="handleValidate()"
    />
  </div>
</template>

<style scoped>
.input {
  display: flex;
  gap: calc(var(--padding) / 1.5);
  align-items: center;
  padding: var(--padding);

  outline: var(--border) solid var(--color-border);
  outline-offset: calc(var(--border) * -1);

  border-radius: var(--radius);

  transition: var(--transition-hover);

  background: var(--color-bg-element);

  &:hover {
    background: var(--color-bg-hover);
  }

  &:focus-within {
    outline: var(--border) solid var(--color-fg-text);
    & > .material-symbols-outlined {
      color: var(--color-fg-focus);
    }
  }
}

.input > input {
  height: min-content;
  width: 0;
  flex-grow: 1;
}

.input > .material-symbols-outlined {
  color: var(--color-fg-dim);
}
</style>
