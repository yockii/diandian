<script lang="ts" setup>
import { Message, Task } from '../../bindings/diandian/background/model/index';
import { Events } from '@wailsio/runtime';
import { onMounted, onUnmounted, ref } from 'vue';
import WelcomeCard from '@/components/WelcomeCard.vue';
import DianDivider from '@/components/DianDivider.vue';
// import { PaperAirplaneIcon } from '@heroicons/vue/24/outline';
import { Bubble, MentionSender } from 'vue-element-plus-x';
import { SettingService, MessageService, WindowService } from '../../bindings/diandian/background/service/index';
import { EVENT_NAMES } from '@/constants/events';
import { ElMessage } from 'element-plus';

// 定义一个自己的Message类型，包含所有Message的字段，额外扩展一个将content转为json对象的字段
class MyMessage extends Message {
  contentObj?: UnifiedMessageResponse;
}
interface UnifiedMessageResponse {
  message_type: string;
  chat_response: string;
  automation_task: AutomationTaskResponse;
  confidence: number;
  explanation: string;
}
interface AutomationTaskResponse {
  task_name: string;
  description: string;
  steps: string[];
  complexity: string;
  risks: string[];
  needs_confirm: boolean;
}

const input = ref('')
const loading = ref(false)
const canWork = ref(false)
const currentTask = ref<Task | null>(null)
const countdown = ref(0)
const countdownTimer = ref<number | null>(null)
const isCountingDown = ref(false)
const isTaskExecuting = ref(false) // 任务执行状态
const isChatLoading = ref(false)   // 聊天加载状态

const sendMessage = async () => {
  // 聊天时只设置聊天加载状态
  isChatLoading.value = true
  loading.value = true
  try {
    const userMsg: Message = {
      id: '0',
      content: input.value,
      role: 'user'
    }
    const assistantMsg : Message = {
      id: '0',
      content: '',
      role: 'assistant'
    }
    await MessageService.NewMessage(userMsg)
    messages.value.push(userMsg, assistantMsg)
    input.value = ''
  } finally {
    // 聊天完成后立即恢复输入框可用状态
    isChatLoading.value = false
    loading.value = false
  }
}

const messages = ref<MyMessage[]>([])

const judgeCanWork = async () => {
  const result = await SettingService.CanWork()
  if (result) {
    canWork.value = true
  } else {
    canWork.value = false
  }
}

const handleConfirmAutomation = (confirmed: boolean) => {
  if (!currentTask.value) return

  try {
    if (confirmed) {
      // 开始5秒倒计时
      startCountdown()
    } else {
      // 取消任务
      cancelTask()
    }
  } catch (error) {
    ElMessage.error('操作失败：' + error)
  }
}

const startCountdown = () => {
  countdown.value = 5
  isCountingDown.value = true

  countdownTimer.value = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      // 倒计时结束，执行任务
      executeTask()
    }
  }, 1000)
}

const cancelCountdown = () => {
  if (countdownTimer.value) {
    clearInterval(countdownTimer.value)
    countdownTimer.value = null
  }
  isCountingDown.value = false
  countdown.value = 0
}

const cancelTask = () => {
  cancelCountdown()
  if (currentTask.value) {
    MessageService.ConfirmAutomationTask(currentTask.value, false)
    ElMessage.info('任务已取消')
  }
}

const executeTask = async () => {
  cancelCountdown()
  if (currentTask.value) {
    try {
      await MessageService.ConfirmAutomationTask(currentTask.value, true)
      ElMessage.success('任务开始执行')
    } catch (error) {
      ElMessage.error('执行失败：' + error)
    }
  }
}

onMounted(() => {
  judgeCanWork()

  Events.On(EVENT_NAMES.CAN_WORK_CHANGED, ({ data }) => {
    canWork.value = data
  })

  Events.On(EVENT_NAMES.MESSAGE_RESPONSED, ({ data }) => {
    const msg = new MyMessage(data)

    if (msg.role === 'assistant') {
      const lastUserMsg = messages.value.findLast((msg) => msg.role === 'user')
      if (lastUserMsg){
        lastUserMsg.conversation_id = data.conversation_id
      }
      const lastMessage = messages.value[messages.value.length - 1]
      if (lastMessage.role === 'assistant') {
        try {
          msg.contentObj = JSON.parse(msg.content)
          if (msg.contentObj) {
            msg.content = msg.contentObj.chat_response
          }
        } catch (error) {
          console.log('解析响应内容失败', error)
        }
        // 替换掉
        messages.value = messages.value.slice(0, messages.value.length - 1)
      }
      messages.value.push(msg)
    }
  })

  Events.On(EVENT_NAMES.TASK_STATUS_CHANGED, ({ data }) => {
    currentTask.value = data
  })

  // 监听任务执行开始事件，切换到浮动窗口
  Events.On(EVENT_NAMES.TASK_EXECUTION_STARTED, ({ data }) => {
    console.log('任务执行开始，切换到浮动窗口')
    isTaskExecuting.value = true  // 设置任务执行状态
    WindowService.HideMainAndShowFloating()
  })

  // 监听任务执行完成事件，恢复主窗口
  Events.On(EVENT_NAMES.TASK_EXECUTION_COMPLETED, ({ data }) => {
    console.log('任务执行完成，恢复主窗口')
    isTaskExecuting.value = false  // 清除任务执行状态
    WindowService.ShowMainWindow()
  })
})

onUnmounted(() => {
  // 清理倒计时定时器
  cancelCountdown()
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
        <welcome-card @ask-selected="input = $event" :can-work="canWork" />
      </div>
      <div v-else>
        <div v-for="(msg, index) in messages" :key="index" class="p-2">
          <Bubble :content="msg.content" :placement="msg.role === 'user' ? 'end' : 'start'" avatar-size="0px" avatar-gap="0px">
            <template #footer>
              <div class="flex flex-col gap-2">
                <template v-if="currentTask?.status === 'pending' && msg.contentObj && msg.contentObj.automation_task?.needs_confirm">
                  <!-- 倒计时界面 -->
                  <div v-if="isCountingDown" class="flex flex-col items-center gap-2 p-3 bg-blue-50 rounded-lg border border-blue-200">
                    <div class="text-lg font-bold text-blue-600">
                      {{ countdown }}秒后自动执行任务
                    </div>
                    <div class="text-sm text-gray-600">
                      点击取消可以停止执行
                    </div>
                    <el-button type="danger" @click="cancelTask" size="small">
                      ❌ 取消执行
                    </el-button>
                  </div>
                  <!-- 确认按钮 -->
                  <div v-else class="flex gap-2">
                    <el-button type="danger" @click="handleConfirmAutomation(false)">❌ 取消</el-button>
                    <el-button type="primary" @click="handleConfirmAutomation(true)">✅ 确认执行</el-button>
                  </div>
                </template>
              </div>
            </template>
          </Bubble>
        </div>
      </div>
    </div>
    <div class="pa-2 no-draggable">
      <mention-sender placeholder="说点什么，让点点来帮你……" v-model="input" clearable @submit="sendMessage" :loading="isChatLoading" :auto-size="{ minRows: 1, maxRows: 4 }" allow-speech
        :disabled="!canWork || isTaskExecuting || isCountingDown">
      </mention-sender>
    </div>


  </div>
</template>
