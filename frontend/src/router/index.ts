import { createMemoryHistory, createRouter } from 'vue-router'
import MainView from '../views/MainView.vue'

const routes = [
  {
    path: '/',
    name: 'Main',
    component: MainView
  }
]

const router = createRouter({
  history: createMemoryHistory(),
  routes,
})

export default router