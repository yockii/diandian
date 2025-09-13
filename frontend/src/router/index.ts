import { createWebHistory, createRouter } from 'vue-router'
import MainLayout from '../layout/main.vue'
import MainView from '../views/MainView.vue'
import FloatingView from '../views/FloatingView.vue'

const routes = [
  {
    path: '/',
    component: MainLayout,
    children: [
      {
        path: '',
        name: 'Main',
        component: MainView
      }
    ]
  },
  {
    path: '/floating',
    name: 'Floating',
    component: FloatingView
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router