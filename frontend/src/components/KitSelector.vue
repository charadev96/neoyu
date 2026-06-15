<script setup lang="ts">
import { ref, computed, nextTick } from "vue"

import KitInputText from "@/components/KitInputText.vue"
import KitButton from "@/components/KitButton.vue"
import KitActions from "@/components/KitActions.vue"

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

const handleAction = async (item: Item, action: string) => {
  switch (action) {
    case "edit":
      await handleEdit(item)
      break
    case "delete":
      handleDelete(item.id)
      break
  }
}
</script>

<template>
  <div class="selector">
    <kit-input-text
      v-model="filterQuery"
      icon="filter_alt"
      placeholder="Filter"
    />
    <div class="items">
      <kit-button icon="add" v-if="state !== 'create'" @click="handleEdit()"
        >Add</kit-button
      >
      <kit-input-text
        v-else
        v-model="editModel"
        icon="add"
        ref="editInput"
        @blur="handleEditCancel()"
        @keyup.esc="handleEditCancel()"
        @keyup.enter="handleEditConfirm()"
        placeholder="Unnamed"
      />
      <template v-for="item in filteredItems">
        <kit-input-text
          v-if="state === 'rename' && selected === item"
          ref="editInput"
          icon="edit"
          @blur="handleEditCancel()"
          @keyup.esc="handleEditCancel()"
          @keyup.enter="handleEditConfirm()"
          v-model="editModel"
          placeholder="Unnamed"
        />
        <kit-button v-else tag="div" class="item" :selected="selected === item">
          <button class="item-button" @click="handleSelect(item)">
            {{ item.name }}
          </button>
          <kit-actions
            :actions="['edit', 'delete']"
            @action="(action) => handleAction(item, action)"
          />
        </kit-button>
      </template>
    </div>
  </div>
</template>

<style scoped>
.selector {
  flex: 1 1 0;
  display: flex;
  flex-direction: column;
  min-height: 0;
  gap: calc(var(--padding) / 2);
}

.items {
  overflow-y: auto;
  width: 100%;
  flex: 1 1 0;
  min-height: 0;

  display: flex;
  flex-direction: column;
  gap: calc(var(--padding) / 2);
}

.item {
  padding: calc(var(--padding) / 4);
  align-items: center;
  justify-content: space-between;

  &:hover > .actions {
    display: flex;
  }
}

.item > button {
  padding: calc(var(--padding) * 0.75);
  text-align: left;
  flex-grow: 1;
  font-weight: inherit;
  color: inherit;
}
</style>
