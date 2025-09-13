<script lang="ts" setup>
import { Events } from '@wailsio/runtime';
import { NewTask } from '../../bindings/changeme/background/service/taskservice';
import { onMounted, ref } from 'vue';

const input = ref('')
const loading = ref(false)

const sendTask = async () => {
    loading.value = true
    try {
        await NewTask(input.value)
        input.value = ''
    } finally {
        loading.value = false
    }
}

const taskSteps = ref<any[]>([])

onMounted(() => {
    Events.On("new-step", ({data}) => {
        taskSteps.value.push(data)
    })
})
</script>

<template>
    <div class="d-flex flex-column" style="height: 100%">
        <div class="text-h5 text-center mb-2">点点小助理</div>
        <v-divider >
            <div class="text-body-2 text-center">历史任务已收起</div>
        </v-divider>
        <div>
            <div v-for="step in taskSteps">
                {{ step }}
            </div>
        </div>
        <div class="flex-fill"></div>
        <div class="pa-2 no-dragable">
            <v-text-field
                density="compact"
                hide-details
                variant="solo"
                placeholder="有任务就点点，点点帮你做"
                append-inner-icon="fa-regular fa-paper-plane"
                @click:append-inner="sendTask"
                v-model="input"
                :loading="loading"
            ></v-text-field>
        </div>
    </div>
</template>