import { createRouter, createWebHistory } from 'vue-router'
import LoadingView from '@/views/LoadingView.vue'
import MainLayout from '@/layout/MainLayout.vue'
import MainView from '@/views/MainView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: MainLayout,
      children: [
        {
          path: '/',
          name: 'loading',
          component: LoadingView,
          meta: {
            bgClass: 'app-background',
            showSettings: false,
          },
        },
        {
          path: '/main',
          name: 'Main',
          component: MainView,
          meta: {
            bgClass: 'app-background',
            showSettings: true,
          },
        },
        {
          path: '/settings',
          name: 'Settings',
          component: () => import('@/views/SettingView.vue'),
          meta: {
            bgClass: 'app-background',
            showSettings: false,
          },
        }
      ],
    },
  ],
})

export default router
