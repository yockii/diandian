<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { WindowService, SettingService } from '../../bindings/changeme/background/service'
import { useRoute, onBeforeRouteUpdate } from 'vue-router'
import { Events } from '@wailsio/runtime'
import { useDark, usePreferredDark } from '@vueuse/core'
import { Cog6ToothIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import { EVENT_NAMES } from '@/constants/events'

const route = useRoute()
const showSettings = ref(true)
const bgClass = ref<string>('')

route.meta.background = 'transparent'

const close = () => {
  if (route.meta.showSettings) {
    WindowService.HideMainAndShowFloating()
  } else {
    WindowService.HideSettings()
  }
};

const openSettings = () => {
  WindowService.ShowSettings()
}

onBeforeRouteUpdate((to, from, next) => {
  const metaBgCls = to.meta.bgClass || ''
  bgClass.value = typeof metaBgCls === 'string' ? metaBgCls : ''
  showSettings.value = to.meta.showSettings === true
  next()
})

// 主题处理
const preferredDark = usePreferredDark()
const isDark = useDark({
  selector: 'html',
  attribute: 'class',
  valueDark: 'dark',
  valueLight: '',
  storageKey: 'user-theme-mode',
})
const currentTheme = ref<string>('auto')

const applyTheme = (mode:string) => {
  currentTheme.value = mode
  if (mode === 'auto') {
    isDark.value = preferredDark.value
  } else {
    isDark.value = mode === 'dark'
  }
}

watch(preferredDark, (newVal) => {
  if (currentTheme.value === 'auto') {
    isDark.value = newVal
  }
})


onMounted(async () => {
  const metaBgCls = route.meta.bgClass || ''
  bgClass.value = typeof metaBgCls === 'string' ? metaBgCls : ''

  showSettings.value = route.meta.showSettings === true

  const themeSetting = await SettingService.GetThemeSetting()
  if (themeSetting) {
    applyTheme(themeSetting.value || '')
  }

  console.log('当前主题：', isDark.value ? 'dark' : 'light')

  Events.On(EVENT_NAMES.THEME_CHANGED, ({data}) => {
    applyTheme(data)
  })
})

onUnmounted(() => {
  Events.OffAll()
})
</script>

<template>
  <el-container class="draggable h-full" :class="[
    bgClass,
    isDark ? 'dark' : 'light',
    isDark ? 'scrollbar-thumb-scroll-thumb-dark scrollbar-track-scroll-track-dark' : 'scrollbar-thumb-scroll-thumb scrollbar-track-scroll-track'
    ]">
    <template v-if="bgClass === 'app-background'">
      <div class="particle"></div>
      <div class="particle"></div>
      <div class="particle"></div>
      <div class="particle"></div>
      <div class="particle"></div>
      <div class="particle"></div>
      <div class="particle"></div>
      <div class="particle"></div>
    </template>
    <el-header height="32px">
      <div class="flex justify-between">
        <div class="-ml-5">
          <el-button v-if="showSettings" link @click="openSettings">
            <el-icon size="24">
              <Cog6ToothIcon />
            </el-icon>
          </el-button>
        </div>
        <div class="-mr-5">
          <el-button link @click="close">
            <el-icon size="24">
              <XMarkIcon />
            </el-icon>
          </el-button>
        </div>
      </div>
    </el-header>
    <el-main class="z-10">
      <router-view></router-view>
    </el-main>
    <el-footer height="18px">
      <div class="text-xs text-center">点点虽小，能动乾坤</div>
    </el-footer>
  </el-container>
</template>

<style scoped>
.v-footer {
  flex: 0;
}

.app-background {
  background-image: linear-gradient(45deg, #ff6b6b, #4ecdc4, #45b7d1, #96ceb4, #ffeaa7, #dda0dd);
  background-size: 400% 400%;
  animation: gradientShift 6s ease infinite;
  position: relative;
  overflow: hidden;
  min-height: 100vh;
}

.dark.app-background {
  background-image: linear-gradient(45deg, #0f172a, #1e293b, #334155, #475569, #581c87, #7c3aed);
  /* 调整为深蓝、紫色等暗色调 */
}

@keyframes gradientShift {
  0% {
    background-position: 0% 50%;
  }

  50% {
    background-position: 100% 50%;
  }

  100% {
    background-position: 0% 50%;
  }
}

/* 添加粒子效果 */
.particle {
  position: absolute;
  width: 10px;
  height: 10px;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 50%;
  animation: float 4s ease-in-out infinite;
  z-index: 0;
}

.particle:nth-child(1) {
  top: 10%;
  left: 10%;
  animation-delay: 0s;
}

.particle:nth-child(2) {
  top: 20%;
  left: 80%;
  animation-delay: 1s;
}

.particle:nth-child(3) {
  top: 30%;
  left: 20%;
  animation-delay: 2s;
}

.particle:nth-child(4) {
  top: 40%;
  left: 70%;
  animation-delay: 3s;
}

.particle:nth-child(5) {
  top: 50%;
  left: 30%;
  animation-delay: 4s;
}

.particle:nth-child(6) {
  top: 60%;
  left: 90%;
  animation-delay: 5s;
}

.particle:nth-child(7) {
  top: 70%;
  left: 40%;
  animation-delay: 6s;
}

.particle:nth-child(8) {
  top: 80%;
  left: 60%;
  animation-delay: 7s;
}

@keyframes float {

  0%,
  100% {
    transform: translateY(0px) rotate(0deg);
    opacity: 0.5;
  }

  50% {
    transform: translateY(-20px) rotate(180deg);
    opacity: 1;
  }
}
</style>
