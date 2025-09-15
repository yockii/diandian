<script lang="ts" setup>
import { Events } from '@wailsio/runtime'
import { WindowService } from '../../bindings/changeme/background/service/index'

const mouseenter = () => {
    Events.Emit({ name: 'mouse-enter-floating' } as Events.WailsEvent)
}
const mouseleave = () => {
    Events.Emit({ name: 'mouse-leave-floating' } as Events.WailsEvent)
}
const showMainWindow = () => {
    WindowService.ShowMainWindow()
}
</script>

<template>
    <div class="draggable container" @mouseenter="mouseenter" @mouseleave="mouseleave" @dblclick="showMainWindow">
        <el-image src="../../public/floating.png" fit="contain" style="height: 100%; width: 100%;"></el-image>
    </div>
</template>

<style scoped>
.container {
    --custom-contextmenu: floating-context-menu;
    height: 100%;
    width: 100%;
}
</style>
