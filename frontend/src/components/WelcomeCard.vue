<script lang="ts" setup>
import { onMounted } from 'vue';

const emit = defineEmits(['ask-selected']);

onMounted(() => {
    navigator.geolocation.getCurrentPosition((position) => {
        console.log(position);
    });
})
</script>

<template>
    <v-card elevation="4" >
        <v-card-text>
            <div class="text-subtitle-1">🎉🎉🎉你好，我是你的桌面智能小助手点点（Diǎn Diǎn）</div>
            <div class="text-body mt-5">
                点点可以与你聊天💬，帮你想办法💡，或者解答你的疑惑😕
            </div>
            <div class="text-body mt-5">
                点点也可以帮你处理一些简单的任务🛠️
                <br />
                “一点小事，包在我身上！❤️”
            </div>
            <div class="text-body mt-5">
                你可以尝试着发送以下内容👇👇👇
            </div>
            <v-card class="mt-3 bg-colored" @click="emit('ask-selected', '今天的天气如何？')">
                <v-card-text>
                    <div class="text-subtitle-2">❓咨询问题</div>
                    <div class="text-body-2 mt-2">今天的天气如何？</div>
                </v-card-text>
            </v-card>
            <v-card class="mt-3 bg-colored" @click="emit('ask-selected', '把桌面上的图片都整理到一个文件夹中')">
                <v-card-text>
                    <div class="text-subtitle-2">🗄️整理文件</div>
                    <div class="text-body-2 mt-2">把桌面上的图片都整理到一个文件夹中</div>
                </v-card-text>
            </v-card>
            <v-card class="mt-3 bg-colored" @click="emit('ask-selected', '帮我打开浏览器')">
                <v-card-text>
                    <div class="text-subtitle-2">🧑‍💻打开应用</div>
                    <div class="text-body-2 mt-2">帮我打开浏览器</div>
                </v-card-text>
            </v-card>
        </v-card-text>
    </v-card>
</template>

<style scoped>
.bg-colored {
  background: linear-gradient(135deg, #f3a9fa 0%, #92a8e4 100%);
  color: white;
}
</style>