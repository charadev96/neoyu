<script setup lang="ts">
import { ref, computed, onMounted } from "vue"

import { providerClient } from "@/client.ts"
import { type Provider } from "@/gen/neoyu/connection/v1/provider_pb"

import KitSelector from "@/components/KitSelector.vue"
import KitForm from "@/components/KitForm.vue"
import { type FormValue, type FormShape } from "@/types.ts"

const providers = ref<Provider[]>([])
const selected = ref<string | undefined>("")

const form = ref<FormShape>({
  type: {
    type: "items",
    name: "Provider Type",
    info: "",
    items: [
      { id: 0, name: "Custom" },
      { id: 1, name: "OpenRouter" },
    ],
  },
  baseUrl: {
    type: "string",
    name: "Base URL",
    info: "Leave empty for provider default.",
  },
  apiKey: {
    type: "string",
    name: "API Key",
  },
  model: {
    type: "string",
    name: "Model ID",
    info: "Name of the model, e.g. llama3.",
  },
})

const selectedProvider = computed(() => {
  return providers.value.find((p) => p.id === selected.value)
})

const fetchProviders = async () => {
  const response = await providerClient.listProviders({})
  providers.value = response.providers
}

const handleCreate = async (name: string) => {
  await providerClient.setProvider({
    provider: {
      id: crypto.randomUUID(),
      name: name,
    },
  })
  fetchProviders()
}

const handleSelect = (id: string | undefined) => {
  selected.value = id
}

const handleRename = async (id: string, name: string) => {
  const provider = providers.value.find((p) => p.id === id)
  if (provider == undefined) {
    return
  }
  await providerClient.setProvider({ provider: { ...provider, name: name } })
  fetchProviders()
}

const handleDelete = async (id: string) => {
  await providerClient.deleteProvider({ id: id })
  fetchProviders()
}

const handleSave = async (value: FormValue) => {
  await providerClient.setProvider({ provider: { ...value } })
  fetchProviders()
}

onMounted(fetchProviders)
</script>

<template>
  <div class="view">
    <div class="sidepanel">
      <h1>Connections</h1>
      <kit-selector
        :items="providers"
        @create="handleCreate"
        @select="handleSelect"
        @rename="handleRename"
        @delete="handleDelete"
      />
    </div>
    <div class="content" v-if="selected && providers">
      <h1>
        {{ selectedProvider?.name }}
      </h1>
      <kit-form
        class="form"
        :clean="selectedProvider ?? {}"
        :form="form"
        @save="
          (value) =>
            handleSave({
              $typeName: 'neoyu.connection.v1.Provider',
              ...value,
            })
        "
      />
    </div>
  </div>
</template>

<style scoped>
.view {
  display: flex;
  width: 100%;
}

.sidepanel,
.content {
  display: flex;
  flex-direction: column;
  padding: calc(var(--padding) * 1.5);
  gap: calc(var(--padding) / 2);
}

.sidepanel {
  min-width: fit-content;
  flex: 0 1 230px;

  border-right: var(--border) solid var(--color-border);
}

.content {
  flex: 1 1 0%;
}

.form {
  align-self: center;
  min-width: min(550px, 100%);
}
</style>
