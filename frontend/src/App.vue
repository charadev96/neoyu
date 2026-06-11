<script setup lang="ts">
import "@/assets/main.css"
import "@fontsource-variable/space-grotesk/wght.css"
import "@fontsource-variable/material-symbols-outlined"

import AppSidebar from "@/components/AppSidebar.vue"
</script>

<template>
  <main>
    <AppSidebar />
    <RouterView />
  </main>
</template>

<style scoped>
main {
  display: flex;
  flex-direction: row;
  height: 100vh;
}
</style>
