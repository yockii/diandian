<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { WindowService, SettingService } from '../../bindings/changeme/background/service'
import { useRoute, onBeforeRouteUpdate } from 'vue-router'
import { useTheme } from 'vuetify'
import { Events } from '@wailsio/runtime'

const route = useRoute()
const showSettings = ref(true)
const bgClass = ref<string>('')
const settedTheme = ref<'light' | 'dark' | 'system'>('light')
const theme = useTheme()

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

onMounted(async () => {
  const metaBgCls = route.meta.bgClass || ''
  bgClass.value = typeof metaBgCls === 'string' ? metaBgCls : ''
  
  showSettings.value = route.meta.showSettings === true

  const themeSetting = await SettingService.GetThemeSetting()
  if (themeSetting) {
    switch (themeSetting.value) {
      case 'light':
        settedTheme.value = 'light'
        break
      case 'dark':
        settedTheme.value = 'dark'
        break
      default:
        settedTheme.value = 'system'
    }
  }

  console.log('当前主题：', theme.name.value)

  Events.On('theme-change', ({data}) => {
    console.log('收到主题变更事件：', data)
    switch (data) {
      case 'light':
        settedTheme.value = 'light'
        break
      case 'dark':
        settedTheme.value = 'dark'
        break
      default:
        settedTheme.value = 'system'
    }
  })
})
</script>

<template>
  <v-app :theme="settedTheme" class="draggable" :class="[bgClass, theme.name.value]">
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
    <v-app-bar elevation="0" height="32" color="rgba(0,0,0,0)" dense flat>
      <template v-slot:prepend>
        <v-btn v-if="showSettings" variant="plain" size="x-small" @click="openSettings">
          <v-icon icon="fas fa-cog"></v-icon>
        </v-btn>
      </template>
      <template v-slot:append>
        <v-btn variant="plain" size="x-small" @click="close">
          <v-icon icon="fas fa-xmark"></v-icon>
        </v-btn>
      </template>
    </v-app-bar>
    <v-main>
      <router-view></router-view>
    </v-main>
    <v-footer color="rgba(0,0,0,0)" height="18">
      <div class="text-caption text-center" style="width: 100%;">点点虽小，能动乾坤</div>
    </v-footer>
  </v-app>
</template>

<style lang="scss" scoped>
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
