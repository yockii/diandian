<script lang="ts" setup>
import { Events } from '@wailsio/runtime'
import { WindowService } from '../../bindings/changeme/background/service'
import Mascot from '@/assets/mascot.png'
import { ref, onMounted, onUnmounted } from 'vue'

const imageRef = ref<any>(null)
let animationInterval: number | null = null

const stickySide = ref<'' | 'left' | 'right' | 'top' | 'bottom'>('')

const mouseenter = () => {
    Events.Emit('mouse-enter-floating')
}
const mouseleave = async () => {
    await Events.Emit('mouse-leave-floating')
    const r = await WindowService.FloatingStickySide()
    switch(r) {
        case 1:
            stickySide.value = 'left'
            break
        case 2:
            stickySide.value = 'right'
            break
        case 3:
            stickySide.value = 'top'
            break
        case 4:
            stickySide.value = 'bottom'
            break
        default:
            stickySide.value = '' 
    }
}
const showMainWindow = () => {
    WindowService.ShowMainWindow()
}

// 史莱姆动画类型
const animations = ['bounce', 'wobble', 'jiggle', 'squish']

const playRandomAnimation = () => {
    if (!imageRef.value) return
    const el = imageRef.value.$el
    if (!el) return
    // 移除之前的动画类
    animations.forEach(anim => el.classList.remove(anim))
    // 随机选择动画
    const randomAnim = animations[Math.floor(Math.random() * animations.length)]
    el.classList.add(randomAnim)
    // 动画持续时间随机 2-5 秒
    const duration = Math.random() * 3000 + 2000
    setTimeout(() => {
        if (imageRef.value && imageRef.value.$el) {
            imageRef.value.$el.classList.remove(randomAnim)
        }
    }, duration)
}

onMounted(() => {
    // 每隔 5-15 秒随机播放一次动画
    const startAnimation = () => {
        const delay = Math.random() * 10000 + 5000
        animationInterval = setTimeout(() => {
            playRandomAnimation()
            startAnimation() // 递归调用以继续循环
        }, delay)
    }
    startAnimation()
})

onUnmounted(() => {
    if (animationInterval) {
        clearTimeout(animationInterval)
    }
})
</script>

<template>
    <div class="w-full h-full p-1">
        <div class="draggable container overflow-hidden" @mouseenter="mouseenter" @mouseleave="mouseleave" @dblclick="showMainWindow">
            <el-image ref="imageRef" :src="Mascot" fit="contain" style="height: 100%; width: 100%;"></el-image>
        </div>
    </div>
</template>

<style scoped>
.container {
    --custom-contextmenu: floating-context-menu;
    height: 100%;
    width: 100%;
}

/* 史莱姆动画定义 */
@keyframes bounce {
    0%, 20%, 50%, 80%, 100% { transform: translateY(0) scaleY(1); }
    40% { transform: translateY(-8px) scaleY(0.95); }
    60% { transform: translateY(-4px) scaleY(1.02); }
}

@keyframes wobble {
    0% { transform: scaleX(1) scaleY(1) rotate(0deg); }
    15% { transform: scaleX(0.9) scaleY(1.1) rotate(-2deg); }
    30% { transform: scaleX(1.1) scaleY(0.9) rotate(2deg); }
    45% { transform: scaleX(0.95) scaleY(1.05) rotate(-1deg); }
    60% { transform: scaleX(1.05) scaleY(0.95) rotate(1deg); }
    75% { transform: scaleX(0.98) scaleY(1.02) rotate(-0.5deg); }
    100% { transform: scaleX(1) scaleY(1) rotate(0deg); }
}

@keyframes jiggle {
    0%, 100% { transform: scale(1) scaleY(1); }
    25% { transform: scale(1.03) rotate(1deg) scaleY(0.97); }
    50% { transform: scale(0.97) rotate(-1deg) scaleY(1.03); }
    75% { transform: scale(1.01) rotate(0.5deg) scaleY(0.99); }
}

@keyframes squish {
    0%, 100% { transform: scaleX(1) scaleY(1); }
    25% { transform: scaleX(1.05) scaleY(0.95); }
    50% { transform: scaleX(0.95) scaleY(1.05); }
    75% { transform: scaleX(1.02) scaleY(0.98); }
}

.bounce {
    animation: bounce 1s ease-in-out;
}

.wobble {
    animation: wobble 1s ease-in-out;
}

.jiggle {
    animation: jiggle 1s ease-in-out;
}

.squish {
    animation: squish 1s ease-in-out;
}
</style>
