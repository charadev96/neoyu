import { createWebHistory, createRouter } from "vue-router"

import ViewHome from "@/views/ViewHome.vue"
import ViewConnections from "@/views/ViewConnections.vue"

const routes = [
  { path: "/", component: ViewHome },
  { path: "/connections", component: ViewConnections },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})
