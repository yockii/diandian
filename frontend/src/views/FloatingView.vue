<script lang="ts" setup>
import { Events } from '@wailsio/runtime'
import { WindowService } from '../../bindings/changeme/background/service'
import Mascot from '@/assets/mascot.png'

const mouseenter = () => {
    Events.Emit('mouse-enter-floating')
}
const mouseleave = () => {
    Events.Emit('mouse-leave-floating')
}
const showMainWindow = () => {
    WindowService.ShowMainWindow()
}
</script>

<template>
    <div class="draggable container overflow-hidden" @mouseenter="mouseenter" @mouseleave="mouseleave" @dblclick="showMainWindow">
        <el-image :src="Mascot" fit="contain" style="height: 100%; width: 100%;"></el-image>
    </div>
</template>

<style scoped>
.container {
    --custom-contextmenu: floating-context-menu;
    height: 100%;
    width: 100%;
}
</style>
