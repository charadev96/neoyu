<script setup lang="ts">
import { ref, onMounted } from "vue"

import { providerClient } from "@/client.ts"
import { type Provider } from "@/gen/neoyu/connection/v1/provider_pb"
import KitSelector from "@/components/KitSelector.vue"

const providers = ref<Provider[]>([])
const selected = ref<string | undefined>("")

const fetchProviders = async () => {
  const response = await providerClient.listProviders({})
  providers.value = response.providers
}

const handleCreate = async (name: string) => {
  await providerClient.setProvider({
    provider: {
      id: crypto.randomUUID(),
      name: name,
      baseUrl: "http://127.0.0.1:8080/",
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

onMounted(fetchProviders)
</script>

<template>
  <div class="view">
    <div class="sidepanel">
      <h1 class="title">Connections</h1>
      <kit-selector
        :items="providers"
        @create="handleCreate"
        @select="handleSelect"
        @rename="handleRename"
        @delete="handleDelete"
      />
    </div>
    <div class="content" v-if="selected && providers">
      <h1 class="title">{{ providers.find((p) => p.id === selected).name }}</h1>
    </div>
  </div>
</template>

<style scoped>
.title {
  padding: var(--padding) 0;
}

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
</style>
