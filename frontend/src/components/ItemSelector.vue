<script setup lang="ts" generic="T">
import { ref, computed } from "vue"

const props = defineProps<{
  items: T[]
  getKey: (item: T) => String
}>()

const emit = defineEmits<{
  select: [item?: T]
}>()

const selected = ref<T>()
const filterQuery = ref("")

const filteredItems = computed(() => {
  if (!filterQuery.value) {
    return props.items
  }
  return props.items.filter((item: T) =>
    props.getKey(item).toLowerCase().includes(filterQuery.value.toLowerCase()),
  )
})

const handleSelect = (item?: T) => {
  if (item === selected.value) {
    return
  }
  selected.value = item
  emit("select", item)
}
</script>

<template>
  <div class="filter">
    <span class="material-symbols-outlined">filter_alt</span>
    <input v-model="filterQuery" type="text" placeholder="Filter" />
  </div>
  <div class="items">
    <button>+ New</button>
    <button
      class="item"
      v-for="item in filteredItems"
      :class="{ active: selected === item }"
      @click="handleSelect(item)"
    >
      <span>{{ props.getKey(item) }}</span>
      <div class="actions">
        <button class="action">
          <span class="material-symbols-outlined">edit</span>
        </button>
        <button class="action">
          <span class="material-symbols-outlined">delete</span>
        </button>
      </div>
    </button>
  </div>
</template>

<style scoped>
.items {
  overflow-y: auto;
  width: 100%;

  display: flex;
  flex-direction: column;
  gap: 4px;
}

.item {
  padding: 4px;

  display: flex;
  align-items: center;
  justify-content: space-between;

  & span {
    padding: 6px;
  }

  &:hover .actions {
    display: flex;
  }
}

.actions {
  display: none;

  background-color: var(--color-bg-body);
  border: var(--border-width) solid var(--color-border);
  border-radius: calc(var(--border-radius) - 4px);
  gap: 2px;
  padding: 2px;
  margin: 0;

  flex-direction: row;
}

.action {
  padding: 2px;
  margin: 0;
  font-size: 10pt;

  border-radius: calc(var(--border-radius) - 8px);

  & span {
    font-size: 11pt;
    padding: 2px;
  }
}

.filter {
  display: flex;
  align-items: center;
  padding: 0 6px;

  border-radius: calc(var(--border-radius) - 2px);
  border: var(--border-width) solid var(--color-border);
  background-color: var(--color-bg-element);

  & span {
    color: var(--color-fg-dim);
    font-size: 14pt;
  }

  & input {
    width: 0;
    flex-grow: 1;
  }

  &:focus-within {
    outline: var(--border-width) solid var(--color-fg-text);
    outline-offset: calc(var(--border-width) * -1);
    & span {
      color: var(--color-fg-text);
    }
  }
}
</style>
