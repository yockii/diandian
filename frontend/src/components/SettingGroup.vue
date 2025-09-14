<script lang="ts" setup>
import { Setting } from '../../bindings/changeme/background/model/models';
import { SettingService } from '../../bindings/changeme/background/service';

const props = defineProps({
    settingList: {
        type: Array<Setting>,
        required: true
    }
})

const settingChanged = async (v:string, setting: Setting) => {
    try {
        await SettingService.SaveSetting({
            id: setting.id || 0,
            value: v
        })
    } catch (e) {
        console.error('保存设置失败', e)
    }
}
</script>

<template>
    <v-row>
        <template v-for="setting in settingList" :key="setting?.id">
            <v-col v-if="setting?.showable" :cols="setting?.cols">
                <v-select v-if="setting?.setting_type === 'select'" :label="setting?.name"
                    :items="JSON.parse(setting?.options || '[]')" v-model="setting.value" variant="solo"
                    density="compact" @update:model-value="(v) => settingChanged(v, setting)"/>
            </v-col>
        </template>
    </v-row>
</template>