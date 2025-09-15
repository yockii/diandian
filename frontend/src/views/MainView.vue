<script lang="ts" setup>
import type { Message } from '../../bindings/changeme/background/model';
import { Events } from '@wailsio/runtime';
import { NewMessage } from '../../bindings/changeme/background/service/messageservice';
import { onMounted, ref } from 'vue';
import WelcomeCard from '@/components/WelcomeCard.vue';
import DianDivider from '@/components/DianDivider.vue';
// import { PaperAirplaneIcon } from '@heroicons/vue/24/outline';
import { Bubble, MentionSender } from 'vue-element-plus-x';
import { SettingService } from '../../bindings/changeme/background/service';

const input = ref('')
const loading = ref(false)
const canWork = ref(false)

const sendTask = async () => {
  loading.value = true
  try {
    const userMsg: Message = {
      id: '0',
      content: input.value,
      role: 'user'
    }
    await NewMessage(userMsg)
    messages.value.push(userMsg)
    input.value = ''
  } finally {
    loading.value = false
  }
}

const messages = ref<Message[]>([])

const judgeCanWork = async () => {
  const result = await SettingService.CanWork()
  if (result) {
    canWork.value = true
  } else {
    canWork.value = false
  }
}

onMounted(() => {
  judgeCanWork()


  Events.On('can_work_changed', ({ data }) => {
    canWork.value = data
  })


  Events.On("new-msg", ({ data }) => {
      messages.value.push(data)
  })
})
</script>

<template>
  <div class="flex flex-col h-full -mt-6">
    <div class="font-bold text-lg text-center">点点小助理</div>

    <dian-divider line-color="#aaa" :line-height="1" position="center" v-if="messages.length === 0">
      <div class="text-xs text-center bg-transparent">历史任务已收起</div>
    </dian-divider>

    <div class="flex-1 overflow-y-auto my-4 scrollbar-thin">
      <div v-if="messages.length === 0" class="p-2 items-center justify-center flex h-full">
        <welcome-card @ask-selected="input = $event" :can-work="canWork"/>
      </div>
      <div v-else>
        <div v-for="(msg, index) in messages" :key="index" class="p-2">
          <Bubble
            :content="msg.content"
            :placement="msg.role === 'user' ? 'end' : 'start'"
            avatar-size="0px"
            avatar-gap="0px" />
        </div>
      </div>
    </div>
    <div class="pa-2 no-draggable">
      <mention-sender
        placeholder="说点什么，让点点来帮你……"
        v-model="input"
        clearable
        @submit="sendTask"
        :loading="loading"
        :auto-size="{ minRows: 1, maxRows: 4 }"
        allow-speech
        :disabled="!canWork"
      >
      </mention-sender>
    </div>
  </div>
</template>
