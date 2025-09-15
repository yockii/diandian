<script lang="ts" setup>
import { onMounted } from 'vue';
import { WindowService } from '../../bindings/changeme/background/service';

defineProps<{
  canWork: boolean
}>();

const emit = defineEmits(['ask-selected']);

const openSettings = () => {
  WindowService.ShowSettings()
}

onMounted(() => {
  navigator.geolocation.getCurrentPosition((position) => {
    console.log(position);
  });
})
</script>

<template>
  <el-card shadow="always">
    <div class="text-base">🎉🎉🎉你好，我是你的桌面智能小助手点点（Diǎn Diǎn）</div>
    <div class="text-sm mt-5">
      点点可以与你聊天💬，帮你想办法💡，或者解答你的疑惑😕
    </div>
    <div class="text-sm mt-5">
      点点也可以帮你处理一些简单的任务🛠️
      <br />
      “一点小事，包在我身上！❤️”
    </div>
    <template v-if="canWork">
      <div class="text-sm mt-5">
        你可以尝试着发送以下内容👇👇👇
      </div>
      <el-card class="mt-3 bg-colored" @click="emit('ask-selected', '今天的天气如何？')" shadow="hover" :body-style="{ padding: '10px', cursor: 'pointer' }">
        <div class="text-sm">❓咨询问题</div>
        <div class="text-xs mt-2">今天的天气如何？</div>
      </el-card>
      <el-card class="mt-3 bg-colored" @click="emit('ask-selected', '把桌面上的图片都整理到一个文件夹中')" shadow="hover" :body-style="{ padding: '10px', cursor: 'pointer' }">
        <div class="text-sm">🗄️整理文件</div>
        <div class="text-xs mt-2">把桌面上的图片都整理到一个文件夹中</div>
      </el-card>
      <el-card class="mt-3 bg-colored" @click="emit('ask-selected', '帮我打开浏览器')" shadow="hover" :body-style="{ padding: '10px', cursor: 'pointer' }">
        <div class="text-sm">🧑‍💻打开应用</div>
        <div class="text-xs mt-2">帮我打开浏览器</div>
      </el-card>
    </template>
    <template v-else>
      <div class="text-base mt-5">
        你还没有配置好点点的工作环境😢
      </div>
      <div class="text-sm mt-3">
        请点击前往👉<span class="text-bold cursor-pointer" @click="openSettings">设置</span>👈界面，配置好大模型访问Token、基础URL以及对应模型后，点点才能智商上线噢❣️
      </div>
    </template>
  </el-card>
</template>

<style scoped>
.bg-colored {
  background: linear-gradient(135deg, #f3a9fa 0%, #92a8e4 100%);
  color: white;
}
</style>
