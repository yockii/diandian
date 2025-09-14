import { createWebHistory, createRouter } from 'vue-router'
import LoadingView from '../views/LoadingView.vue'
import MainLayout from '../layout/main.vue'
import MainView from '../views/MainView.vue'
import FloatingView from '../views/FloatingView.vue'
import SettingView from '../views/SettingView.vue'

const routes = [
  {
    path: '/',
    component: MainLayout,
    children: [
      {
        path: '/',
        name: 'Loading',
        component: LoadingView,
        meta: {
          bgClass: 'app-background',
          showSettings: false,
        }
      },
      {
        path: '/main',
        name: 'Main',
        component: MainView,
        meta: {
          bgClass: 'app-background',
          showSettings: true,
        }
      },
      {
        path: '/settings',
        name: 'Settings',
        component: SettingView,
        meta: {
          showSettings: false,
        }
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