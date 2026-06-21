<script setup lang="ts">
import { ref, computed, watch, onMounted, toRaw } from "vue"

import KitInput from "@/components/KitInput.vue"
import KitButton from "@/components/KitButton.vue"
import KitSelector from "@/components/KitSelector.vue"
import { type Item, type FormValue, type FormShape } from "@/types.ts"

const { form = {}, clean = {} } = defineProps<{
  form: FormShape
  clean: FormValue
}>()

const emit = defineEmits<{
  save: [value: FormValue]
}>()

const dirty = ref<FormValue>({})

watch(
  () => clean,
  async () => {
    handleRevert()
  },
)

const handleSave = () => {
  emit("save", dirty.value)
}

const handleRevert = () => {
  dirty.value = structuredClone(toRaw(clean))
}

const isDirty = computed(() => {
  return JSON.stringify(dirty.value) !== JSON.stringify(clean)
})

onMounted(handleRevert)
</script>

<template>
  <div class="form">
    <div class="actions">
      <template v-if="isDirty">
        <kit-button icon="save" @click="handleSave()">Save</kit-button>
        <kit-button icon="undo" @click="handleRevert()">Revert</kit-button>
      </template>
      <kit-button icon="check_indeterminate_small" v-else
        >No changes</kit-button
      >
    </div>
    <div class="fields">
      <template v-for="(field, key, index) in form">
        <div class="field">
          <div class="info">
            <h2>
              <span v-if="dirty[key] !== clean[key]">* </span>
              {{ field.name }}
            </h2>
            <p style="color: var(--color-fg-dim)">
              {{ field.info }}
            </p>
          </div>
          <div class="value">
            <kit-input
              v-if="field.type === 'string'"
              v-model="dirty[key]"
              type="text"
              placeholder="Undefined"
            />
            <kit-input
              v-if="field.type === 'number'"
              v-model="dirty[key]"
              type="number"
              :min="field?.min"
              :max="field?.max"
              placeholder="Undefined"
            />
            <kit-selector
              v-if="field.type == 'items'"
              v-model="dirty[key]"
              compact
              :items="field.items"
            />
          </div>
        </div>
        <hr class="separator" v-if="index < Object.keys(form).length - 1" />
      </template>
    </div>
  </div>
</template>

<style scoped>
.form {
  display: flex;
  flex-direction: column;
  gap: calc(var(--padding) / 2);
}

.actions {
  display: flex;
  gap: calc(var(--padding) / 4);
}

.fields {
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  border: var(--border) solid var(--color-border);
  border-radius: var(--radius);
}

.field {
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: var(--padding) calc(var(--padding) * 1.5);
  justify-content: space-between;
  gap: calc(var(--padding) * 3);
}

.info {
  display: flex;
  flex-direction: column;
  flex: 0 1 40%;
}

.value {
  display: flex;
  flex-direction: column;
  flex: 1 1 60%;
}

.separator {
  margin: 0;
}
</style>
