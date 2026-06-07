<script setup lang="ts">
import { ref, onMounted } from "vue"

import { providerClient } from "@/client.ts"
import { type Provider } from "@/gen/neoyu/connection/v1/provider_pb"
import ItemSelector from "@/components/ItemSelector.vue"

const providers = ref<Provider[]>([])
const selected = ref<Provider>()

onMounted(async () => {
  const response = await providerClient.listProviders({})
  providers.value = response.providers
})
</script>

<template>
  <div class="view">
    <div class="sidepanel">
      <h2>Connections</h2>
      <ItemSelector
        @select="(item) => (selected = item)"
        :items="providers"
        :getKey="(item) => item.name"
      />
    </div>
    <div class="content" v-if="selected">
      <h2>{{ selected.name }}</h2>
    </div>
  </div>
</template>

<style scoped>
.view {
  display: flex;
  width: 100%;
}

.sidepanel {
  min-width: fit-content;
  flex: 0 1 230px;

  display: flex;
  flex-direction: column;
  gap: 4px;

  padding: 4px 8px;

  border: none;
  border-right: var(--border-width) solid var(--color-border);
}

.content {
  overflow-y: auto;
  flex: 1 1 0%;

  display: flex;
  flex-direction: column;
  gap: 4px;

  padding: 4px 8px;
}
</style>
