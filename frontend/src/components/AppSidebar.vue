<script setup lang="ts">
import { ref } from "vue"

const links = ref([
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
    <RouterLink :to="'/'">
      <button :class="{ active: '/' === $route.path }">
        <span class="material-symbols-outlined">home</span>
      </button>
    </RouterLink>
    <hr />
    <RouterLink v-for="link in links" :to="link.to">
      <button :class="{ active: link.to === $route.path }">
        <span class="material-symbols-outlined">{{ link.icon }}</span>
      </button>
    </RouterLink>
  </aside>
</template>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  gap: 6px;

  height: auto;
  margin: 0px;
  padding: 12px;
  border: none;
  border-right: var(--border-width) solid var(--color-border);
  border-radius: 0px;
}
</style>
