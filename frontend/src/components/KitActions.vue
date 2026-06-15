<script setup lang="ts">
import KitButton from "@/components/KitButton.vue"

const {} = defineProps<{
  actions: string[]
}>()

const emit = defineEmits<{
  action: [id: string]
}>()

const handleAction = (action: string) => {
  emit("action", action)
}
</script>

<template>
  <div class="actions">
    <kit-button
      v-for="action in actions"
      :icon="action"
      @click="handleAction(action)"
      class="action"
    />
  </div>
</template>

<style scoped>
.actions {
  display: none;

  background: var(--color-bg-body);
  border: var(--border) solid var(--color-border);
  border-radius: calc(var(--radius) * 0.75);
  gap: calc(var(--padding) / 4);
  padding: calc(var(--padding) / 4);

  flex-direction: row;
}

.action {
  padding: calc(var(--padding) / 4);
  margin: 0;

  border-radius: calc(var(--radius) - (var(--padding) * 0.75));
}
</style>
