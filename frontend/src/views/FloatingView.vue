<script lang="ts" setup>
import { Events } from '@wailsio/runtime'
import { WindowService } from '../../bindings/diandian/background/service'
import Mascot from '@/assets/mascot.png'
import { ref, onMounted, onUnmounted } from 'vue'
import { EVENT_NAMES } from '@/constants/events';

const imageRef = ref<any>(null)
let animationInterval: number | null = null
const isMouseInWindow = ref(false)
let leaveTimeout: number | null = null

const stickySide = ref<'' | 'left' | 'right' | 'top' | 'bottom'>('')

const mouseenter = () => {
    isMouseInWindow.value = true
     // 清除之前的延迟隐藏
    if (leaveTimeout) {
        clearTimeout(leaveTimeout)
        leaveTimeout = null
    }
    Events.Emit(EVENT_NAMES.MOUSE_ENTER_FLOATING)
    // 鼠标进入时，如果有贴边状态，需要移除旋转
    if (stickySide.value) {
        updateRotation('')
    }
}
const mouseleave = async () => {
    isMouseInWindow.value = false
    // 延迟200ms后再隐藏，避免闪烁
    leaveTimeout = setTimeout(async () => {
        await Events.Emit(EVENT_NAMES.MOUSE_LEAVE_FLOATING)
        const r = await WindowService.FloatingStickySide()
        updateStickySide(r)
        // 鼠标离开时，如果有贴边状态，需要设置旋转
        if (stickySide.value) {
            updateRotation(stickySide.value)
        }
    }, 200)
}

const updateStickySide = (sideValue: number) => {
    switch(sideValue) {
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

// 更新旋转角度
const updateRotation = (side: string) => {
    if (!imageRef.value || !imageRef.value.$el) return

    const el = imageRef.value.$el
    // 移除之前的旋转类
    el.classList.remove('rotate-left', 'rotate-right', 'rotate-top', 'rotate-bottom', 'rotate-none')

    // 根据贴边方向添加相应的旋转类
    switch(side) {
        case 'left':
            el.classList.add('rotate-left')
            break
        case 'right':
            el.classList.add('rotate-right')
            break
        case 'top':
            el.classList.add('rotate-top')
            break
        case 'bottom':
            el.classList.add('rotate-bottom')
            break
        default:
            el.classList.add('rotate-none')
    }
}

const playRandomAnimation = () => {
    // 如果有贴边状态且鼠标不在窗口中，不播放动画
    if (stickySide.value && !isMouseInWindow.value) {
        return
    }

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

onMounted(async () => {
    // 初始化时检查贴边状态
    const r = await WindowService.FloatingStickySide()
    updateStickySide(r)

    // 如果有贴边状态，设置旋转
    if (stickySide.value) {
        updateRotation(stickySide.value)
    }

    // 监听贴边状态变化事件
    Events.On(EVENT_NAMES.STICKY_SIDE_CHANGED, ({data}) => {
        updateStickySide(data)
        // 如果鼠标不在窗口中，根据新状态更新旋转
        if (!isMouseInWindow.value) {
            updateRotation(stickySide.value)
        }
    })

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
    if (leaveTimeout) {
        clearTimeout(leaveTimeout)
    }
})
</script>

<template>
    <div class="w-full h-full p-1" @mouseenter="mouseenter" @mouseleave="mouseleave" @dblclick="showMainWindow">
        <div class="draggable container overflow-hidden">
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

/* 旋转样式 */
.rotate-left {
    transform: rotate(90deg);
}

.rotate-right {
    transform: rotate(-90deg);
}

.rotate-top {
    transform: rotate(180deg);
}

.rotate-bottom {
    transform: rotate(0deg);
}

.rotate-none {
    transform: rotate(0deg);
}
</style>
