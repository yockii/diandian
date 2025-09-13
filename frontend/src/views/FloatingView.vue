<script lang="ts" setup>
import { Events } from '@wailsio/runtime';
const mouseenter = () => {
    Events.Emit({name: 'mouse-enter-floating'} as Events.WailsEvent)
}
const mouseleave = () => {
    Events.Emit({name: 'mouse-leave-floating'} as Events.WailsEvent)
}
</script>

<template>
    <div class="draggable container" @mouseenter="mouseenter" @mouseleave="mouseleave">
        cccccccccc
    </div>
</template>

<style scoped>
.container {
    --custom-contextmenu:floating-context-menu;
    height:100%;
    width: 100%;
    background-color: rgba(255, 255, 255, 0.2);
}
</style>
