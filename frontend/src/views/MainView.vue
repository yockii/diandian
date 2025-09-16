<script lang="ts" setup>
import { Message, Task } from '../../bindings/diandian/background/model/index';
import { Events } from '@wailsio/runtime';
import { onMounted, ref } from 'vue';
import WelcomeCard from '@/components/WelcomeCard.vue';
import DianDivider from '@/components/DianDivider.vue';
// import { PaperAirplaneIcon } from '@heroicons/vue/24/outline';
import { Bubble, MentionSender } from 'vue-element-plus-x';
import { SettingService, MessageService } from '../../bindings/diandian/background/service/index';
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

const sendMessage = async () => {
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
    // 调用后端确认方法
    if (confirmed) {
      // 如果确认执行，显示提示信息
      ElMessage.success('任务已确认，3秒后开始执行')
      setTimeout(async () => {
        if (currentTask.value){
          await MessageService.ConfirmAutomationTask(currentTask.value, confirmed)
        }
      }, 3000)
    } else {
      ElMessage.info('任务已取消')
    }
  } catch (error) {
    ElMessage.error('操作失败：' + error)
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
              <div class="flex">
                <template v-if="currentTask?.status === 'pending' && msg.contentObj && msg.contentObj.automation_task?.needs_confirm">
                  <el-button type="danger" @click="handleConfirmAutomation(false)">❌ 取消</el-button>
                  <el-button type="primary" @click="handleConfirmAutomation(true)">✅ 确认执行</el-button>
                </template>
              </div>
            </template>
          </Bubble>
        </div>
      </div>
    </div>
    <div class="pa-2 no-draggable">
      <mention-sender placeholder="说点什么，让点点来帮你……" v-model="input" clearable @submit="sendMessage" :loading="loading" :auto-size="{ minRows: 1, maxRows: 4 }" allow-speech
        :disabled="!canWork">
      </mention-sender>
    </div>

    <!-- 自动化任务确认对话框 -->
    <!-- <el-dialog v-model="showConfirmDialog" title="🤖 自动化任务确认" width="500px" :close-on-click-modal="false" :close-on-press-escape="false">
      <div v-if="confirmData">
        <div class="mb-4">
          <h4 class="text-lg font-semibold mb-2">{{ confirmData.analysis?.task_name }}</h4>
          <p class="text-gray-600 mb-3">{{ confirmData.analysis?.description }}</p>

          <div class="mb-4 p-3 bg-gray-50 rounded">
            <div class="flex items-center justify-between mb-2">
              <span class="text-sm font-medium">复杂度：</span>
              <el-tag :type="confirmData.analysis?.complexity === 'simple' ? 'success' : confirmData.analysis?.complexity === 'medium' ? 'warning' : 'danger'">
                {{ confirmData.analysis?.complexity }}
              </el-tag>
            </div>

            <div v-if="confirmData.analysis?.steps?.length" class="mb-3">
              <p class="text-sm font-medium mb-2">📋 执行步骤：</p>
              <ol class="list-decimal list-inside text-sm space-y-1 pl-2">
                <li v-for="step in confirmData.analysis.steps" :key="step" class="text-gray-700">{{ step }}</li>
              </ol>
            </div>

            <div v-if="confirmData.analysis?.risks?.length" class="mb-3">
              <p class="text-sm font-medium mb-2 text-orange-600">⚠️ 风险提示：</p>
              <ul class="list-disc list-inside text-sm space-y-1 text-orange-600 pl-2">
                <li v-for="risk in confirmData.analysis.risks" :key="risk">{{ risk }}</li>
              </ul>
            </div>
          </div>

          <div class="mb-4 p-3 bg-blue-50 rounded text-sm text-blue-700">
            <p class="font-medium mb-1">🔔 执行说明：</p>
            <ul class="space-y-1 text-xs">
              <li>• 确认后界面将切换到浮动模式</li>
              <li>• 任务将自动执行，无需手动干预</li>
              <li>• 执行期间请勿操作电脑</li>
              <li>• 可通过浮动窗口监控进度</li>
            </ul>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end space-x-3">
          <el-button @click="handleConfirmAutomation(false)">
            ❌ 取消
          </el-button>
          <el-button type="primary" @click="handleConfirmAutomation(true)">
            ✅ 确认执行
          </el-button>
        </div>
      </template>
    </el-dialog> -->
  </div>
</template>
