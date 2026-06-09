<script setup lang="ts">
import { ref, computed, nextTick } from "vue"

interface Item {
  id: string
  name: string
}

const props = defineProps<{
  items: Item[]
}>()

const emit = defineEmits<{
  create: [name: string]
  select: [id: string | undefined]
  rename: [id: string, name: string]
  delete: [id: string]
}>()

const selected = ref<Item>()
const editInput = ref()
const editModel = ref("")

const state = ref<"rename" | "create">()
const filterQuery = ref("")

const filteredItems = computed(() => {
  if (!filterQuery.value) {
    return props.items
  }
  return props.items.filter((item: Item) =>
    item.name.toLowerCase().includes(filterQuery.value.toLowerCase()),
  )
})

const handleSelect = (item?: Item) => {
  if (selected.value === item) {
    return
  }
  selected.value = item
  emit("select", item?.id)
}

const handleEdit = async (item?: Item) => {
  selected.value = item ?? { id: "", name: "" }
  state.value = item == undefined ? "create" : "rename"
  editModel.value = selected.value.name

  await nextTick()
  const input = editInput.value[0] ?? editInput.value
  input.focus()
}

const handleEditCancel = () => {
  state.value = undefined
}

const handleEditConfirm = () => {
  if (selected.value == undefined) {
    handleEditCancel()
    return
  }

  if (state.value === "create") {
    emit("create", editModel.value || "Unnamed")
  } else {
    emit("rename", selected.value.id, editModel.value || "Unnamed")
  }
  state.value = undefined
}

const handleDelete = (id: string) => {
  handleSelect(undefined)
  emit("delete", id)
}
</script>

<template>
  <div class="text-entry selected">
    <span class="material-symbols-outlined">filter_alt</span>
    <input v-model="filterQuery" type="text" placeholder="Filter" />
  </div>
  <div class="items">
    <button v-if="state !== 'create'" @click="handleEdit()" class="button">
      <span class="material-symbols-outlined">add</span>
      New
    </button>
    <div v-else class="text-entry selected">
      <span class="material-symbols-outlined">add</span>
      <input
        type="text"
        ref="editInput"
        @blur="handleEditCancel()"
        @keyup.esc="handleEditCancel()"
        @keyup.enter="handleEditConfirm()"
        v-model="editModel"
        placeholder="Unnamed"
      />
    </div>
    <template v-for="item in filteredItems">
      <div
        class="text-entry selected"
        v-if="state === 'rename' && selected === item"
      >
        <span class="material-symbols-outlined">edit</span>
        <input
          type="text"
          ref="editInput"
          @blur="handleEditCancel()"
          @keyup.esc="handleEditCancel()"
          @keyup.enter="handleEditConfirm()"
          v-model="editModel"
          placeholder="Unnamed"
        />
      </div>
      <div
        v-else
        :class="{ button: true, item: true, selected: selected === item }"
      >
        <button @click="handleSelect(item)">
          {{ item.name }}
        </button>
        <div class="actions">
          <button class="button action" @click="handleEdit(item)">
            <span class="material-symbols-outlined">edit</span>
          </button>
          <button class="button action" @click="handleDelete(item.id)">
            <span class="material-symbols-outlined">delete</span>
          </button>
        </div>
      </div>
    </template>
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

  align-items: center;
  justify-content: space-between;

  & > button {
    text-align: left;
    flex-grow: 1;
    padding: 6px;
  }

  &:hover > button {
    transition: color 0.15s ease;
    color: var(--color-fg-focus);
  }

  &:hover .actions {
    display: flex;
  }
}

.actions {
  display: none;

  background: var(--color-bg-body);
  border: var(--border-width) solid var(--color-border);
  border-radius: calc(var(--border-radius) - 2px);
  gap: 2px;
  padding: 2px;

  flex-direction: row;
}

.action {
  padding: 2px;
  margin: 0;
  font-size: 10pt;

  border-radius: calc(var(--border-radius) - 6px);

  & > span {
    font-size: 11pt;
    padding: 2px;
  }
}

.text-entry {
  display: flex;
  align-items: center;
  padding: 4px;

  border-radius: var(--border-radius);
  background: var(--color-bg-element);

  & > .material-symbols-outlined {
    color: var(--color-fg-dim);
    padding: 6px 0 6px 6px;
    font-size: 14pt;
  }

  & > input {
    padding: 6px;
    height: min-content;
    width: 0;
    flex-grow: 1;
  }

  &:focus-within {
    outline: var(--border-width) solid var(--color-fg-text);
    & > span {
      color: var(--color-fg-text);
    }
  }
}
</style>
