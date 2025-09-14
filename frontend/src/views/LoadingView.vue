<script lang="ts" setup>
import { onMounted, ref } from 'vue'
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
    <div class="d-flex flex-column fill-height align-center justify-center">
        <div class="loading-spinner">
            <v-img src="../../public/floating.png" width="64" height="64" contain></v-img>
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