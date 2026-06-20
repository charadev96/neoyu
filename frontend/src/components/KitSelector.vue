<script setup lang="ts" generic="T extends ID">
import { ref, computed, nextTick } from "vue"

import KitInput from "@/components/KitInput.vue"
import KitButton from "@/components/KitButton.vue"
import KitActions from "@/components/KitActions.vue"
import { type ID, type Item } from "@/types.ts"

const { items = [], compact = false } = defineProps<{
  compact?: boolean
  items: Item<T>[]
}>()

const emit = defineEmits<{
  create: [name: string]
  select: [id: T | undefined]
  rename: [id: T, name: string]
  delete: [id: T]
}>()

const selected = defineModel<T>()

const open = ref(false)
const state = ref<"rename" | "create">()

const editRef = ref()
const editName = ref("")

const filterRef = ref<HTMLElement>()
const filterQuery = ref("")

const itemsRef = ref<HTMLElement>()

const filteredItems = computed(() => {
  if (!filterQuery.value) return items
  return items.filter((item: Item<T>) =>
    item.name.toLowerCase().includes(filterQuery.value.toLowerCase()),
  )
})

const selectedName = computed(() => {
  return items.find((item: Item<T>) => selected.value === item.id)?.name
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

const handleSelect = (id: T | undefined) => {
  if (compact) filterRef.value?.blur()
  if (selected.value === id) return
  selected.value = id
  emit("select", id)
}

const handleEdit = async (id: T | undefined) => {
  handleSelect(id)
  state.value = id == undefined ? "create" : "rename"
  editName.value = selectedName.value ?? ""

  nextTick(() => {
    const el = editRef.value[0] ?? editRef.value
    el.focus()
  })
}

const handleEditCancel = () => {
  state.value = undefined
}

const handleEditConfirm = () => {
  switch (state.value) {
    case "create":
      emit("create", editName.value || "Unnamed")
      break
    case "rename":
      if (selected.value == undefined) return
      emit("rename", selected.value, editName.value || "Unnamed")
      break
  }
  state.value = undefined
}

const handleDelete = (id: T) => {
  handleSelect(undefined)
  emit("delete", id)
}

const handleAction = async (item: Item<T>, action: string) => {
  switch (action) {
    case "edit":
      await handleEdit(item.id)
      break
    case "delete":
      handleDelete(item.id)
      break
  }
}
</script>

<template>
  <div class="selector" :class="{ compact: compact }">
    <kit-input
      class="filter"
      v-model="filterQuery"
      ref="filterRef"
      :icon="compact && !open ? 'keyboard_arrow_down' : 'filter_alt'"
      @focus="handleOpen()"
      @blur="handleClose()"
      :class="{ compact: open && compact }"
      :placeholder="!open && compact ? (selectedName ?? 'None') : 'Filter'"
    />
    <div
      v-if="open || !compact"
      class="items"
      ref="itemsRef"
      :class="{ compact: open && compact }"
    >
      <template v-if="!compact">
        <kit-button
          icon="add"
          v-if="state !== 'create'"
          @click="handleEdit(undefined)"
          >Add</kit-button
        >
        <kit-input
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
        <kit-input
          v-if="state === 'rename' && selected === item.id"
          ref="editRef"
          icon="edit"
          @blur="handleEditCancel()"
          @keyup.esc="handleEditCancel()"
          @keyup.enter="handleEditConfirm()"
          v-model="editName"
          placeholder="Unnamed"
        />
        <kit-button
          v-else
          tag="div"
          class="item"
          :selected="selected === item.id"
        >
          <button
            class="item-button"
            @mousedown.prevent="handleSelect(item.id)"
          >
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
