<script setup lang="ts">
import { ref } from "vue"
import KitButton from "@/components/KitButton.vue"

const links = ref([
  { to: "/", icon: "home" },
  //{ to: "/sessions", icon: "chat_bubble" },
  { to: "/connections", icon: "conversion_path" },
  //{ to: "/knowledge", icon: "book_2" },
  //{ to: "/prompts", icon: "convert_to_text" },
  //{ to: "/settings", icon: "settings" },
  //{ to: "/about", icon: "info" },
])
</script>

<template>
  <aside class="sidebar">
    <kit-button
      class="button"
      tag="router-link"
      v-for="link in links"
      :icon="link.icon"
      :to="link.to"
      :selected="link.to === $route.path"
    />
  </aside>
</template>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  gap: var(--padding);

  padding: calc(var(--padding) * 1.5);
  border-right: var(--border) solid var(--color-border);
}

.button * {
  font-size: 1.6rem;
}
</style>
