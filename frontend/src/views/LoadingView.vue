<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import mascot from '@/assets/mascot.png'
import { WindowService } from '../../bindings/changeme/background/service'
import { useRouter } from 'vue-router'

const title = ref('点点正在启动中...')
const router = useRouter()

onMounted(() => {
  const checkInterval = setInterval(() => {
      WindowService.IsInitializeSuccess().then((success: boolean) => {
          if (success) {
              clearInterval(checkInterval)
              title.value = "点点启动成功，正在进入新世界..."
              setTimeout(() => {
                  router.push('/main')
              }, 500)
          }
      })
  }, 3000)
})
</script>

<template>
  <div class="flex flex-col fill-height items-center justify-center h-full">
    <div class="loading-spinner">
      <el-image :src="mascot" style="width: 64px; height: 64px;" fit="contain"></el-image>
    </div>
    <div class="loading-text">{{ title }}</div>
  </div>
</template>

<style scoped>
.loading-spinner {
  animation: bounce 1.5s ease-in-out infinite;
}

@keyframes bounce {
  0% {
    transform: translateY(0) scaleY(1);
  }

  25% {
    transform: translateY(-20px) scaleY(0.9);
  }

  50% {
    transform: translateY(0) scaleY(1.1);
  }

  75% {
    transform: translateY(-10px) scaleY(0.95);
  }

  100% {
    transform: translateY(0) scaleY(1);
  }
}
</style>
