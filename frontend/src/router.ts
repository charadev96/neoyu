import { createWebHistory, createRouter } from "vue-router"

import ViewConnections from "@/views/ViewConnections.vue"

const routes = [{ path: "/connections", component: ViewConnections }]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})
