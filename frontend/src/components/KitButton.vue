<script setup lang="ts">
const { tag = "button" } = defineProps<{
  tag?: string
  icon?: string
  selected?: boolean
}>()
</script>

<template>
  <component
    :is="tag"
    class="button"
    v-bind="$attrs"
    :class="{ selected: selected }"
  >
    <span v-if="icon" class="material-symbols-outlined">{{ icon }}</span>
    <slot />
  </component>
</template>

<style scoped>
.button {
  display: flex;
  gap: calc(var(--padding) / 1.5);
  align-items: center;

  outline: 0 solid var(--color-border);
  outline-offset: calc(var(--border) * -1);

  border-radius: var(--radius);

  padding: var(--padding);

  font-weight: 450;
  color: var(--color-fg-dim);

  background: transparent;

  transition: var(--transition-hover);
}

.button.selected {
  font-weight: 550;
  color: var(--color-fg-focus);
  outline-width: var(--border);
  background: var(--color-bg-element);
}

.button > .material-symbols-outlined {
  font-weight: 350;
}

.button.selected > .material-symbols-outlined {
  font-weight: 550;
}

.button:hover {
  color: var(--color-fg-focus);
  background: var(--color-bg-hover);
}
</style>
