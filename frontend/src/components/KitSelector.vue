<script setup lang="ts">
import { ref, computed, nextTick } from "vue"

import KitInputText from "@/components/KitInputText.vue"
import KitButton from "@/components/KitButton.vue"
import KitActions from "@/components/KitActions.vue"

interface Item {
  id: string
  name: string
}

const { items = [], compact = false } = defineProps<{
  compact?: boolean
  items: Item[]
}>()

const emit = defineEmits<{
  create: [name: string]
  select: [id: string | undefined]
  rename: [id: string, name: string]
  delete: [id: string]
}>()

const selected = defineModel<Item>()

const open = ref(false)
const state = ref<"rename" | "create">()

const editRef = ref<HTMLElement>()
const editName = ref("")

const filterRef = ref<HTMLElement>()
const filterQuery = ref("")

const itemsRef = ref<HTMLElement>()

const filteredItems = computed(() => {
  if (!filterQuery.value) return items
  return items.filter((item: Item) =>
    item.name.toLowerCase().includes(filterQuery.value.toLowerCase()),
  )
})

const handleOpen = () => {
  if (!compact) return
  open.value = true
  nextTick(() => clampPopup())
}

const handleClose = () => {
  if (!compact) return
  open.value = false
}

const clampPopup = () => {
  const el = itemsRef.value
  if (!el) return
  const rect = el.getBoundingClientRect()
  const available = window.innerHeight - rect.top - 8
  el.style.maxHeight = `${Math.max(available, 64)}px`
}

const handleSelect = (item?: Item) => {
  if (compact) filterRef.value?.blur()
  if (selected.value === item) return

  selected.value = item
  emit("select", item?.id)
}

const handleEdit = async (item?: Item) => {
  selected.value = item ?? { id: "", name: "" }
  state.value = item == undefined ? "create" : "rename"
  editName.value = selected.value.name

  nextTick(() => {
    const el = editRef.value[0] ?? editRef.value
    el.focus()
  })
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
    emit("create", editName.value || "Unnamed")
  } else {
    emit("rename", selected.value.id, editName.value || "Unnamed")
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
  <div class="selector" :class="{ compact: compact }">
    <kit-input-text
      class="filter"
      v-model="filterQuery"
      ref="filterRef"
      :icon="compact && !open ? 'keyboard_arrow_down' : 'filter_alt'"
      @focus="handleOpen()"
      @blur="handleClose()"
      :class="{ compact: open && compact }"
      :placeholder="!open && compact ? (selected?.name ?? 'None') : 'Filter'"
    />
    <div
      v-if="open || !compact"
      class="items"
      ref="itemsRef"
      :class="{ compact: open & compact }"
    >
      <template v-if="!compact">
        <kit-button icon="add" v-if="state !== 'create'" @click="handleEdit()"
          >Add</kit-button
        >
        <kit-input-text
          v-else
          v-model="editName"
          icon="add"
          ref="editRef"
          @blur="handleEditCancel()"
          @keyup.esc="handleEditCancel()"
          @keyup.enter="handleEditConfirm()"
          placeholder="Unnamed"
        />
      </template>
      <template v-for="item in filteredItems" :key="item.id">
        <kit-input-text
          v-if="state === 'rename' && selected === item"
          ref="editRef"
          icon="edit"
          @blur="handleEditCancel()"
          @keyup.esc="handleEditCancel()"
          @keyup.enter="handleEditConfirm()"
          v-model="editName"
          placeholder="Unnamed"
        />
        <kit-button v-else tag="div" class="item" :selected="selected === item">
          <button class="item-button" @mousedown.prevent="handleSelect(item)">
            {{ item.name }}
          </button>
          <kit-actions
            v-if="!compact"
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
  display: flex;
  flex-direction: column;
  gap: calc(var(--padding) / 2);
  min-height: 0;
}

.selector.compact {
  gap: 0;
  position: relative;
  min-height: unset;
}

.filter.compact {
  border-bottom-left-radius: 0;
  border-bottom-right-radius: 0;
}

.items {
  display: flex;
  flex-direction: column;
  gap: calc(var(--padding) / 2);

  overflow-y: auto;
  min-height: 0;
}

.items.compact {
  position: absolute;
  top: 100%;
  left: 0;
  width: 100%;
  padding: calc(var(--padding) / 2);

  background: var(--color-bg-body);
  border: var(--border) solid var(--color-border);
  border-top: none;
  border-bottom-left-radius: var(--radius);
  border-bottom-right-radius: var(--radius);
  z-index: 10;
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
