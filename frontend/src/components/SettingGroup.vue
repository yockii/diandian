<script lang="ts" setup>
import { Setting } from '../../bindings/diandian/background/model/models';
import { SettingService } from '../../bindings/diandian/background/service';

defineProps({
  settingList: {
    type: Array<Setting>,
    required: true
  }
})

const settingChanged = async (v: string, setting: Setting) => {
  try {
    await SettingService.SaveSetting({
      id: setting.id,
      value: v
    })
  } catch (e) {
    console.error('保存设置失败', e)
  }
}
</script>

<template>
  <el-row :gutter="20">
    <template v-for="setting in settingList" :key="setting?.id">
      <el-col v-if="setting?.showable" :span="setting?.cols" class="mt-2">
        <div class="text-sm mb-1 font-medium">{{ setting?.name }}</div>
        <div class="no-draggable">
          <el-select v-if="setting?.setting_type === 'select'" :label="setting?.name" :options="JSON.parse(setting?.options || '[]')" v-model="setting.value"
            @change="(v:any) => settingChanged(v, setting)" />
          <el-input v-else-if="setting?.setting_type === 'input'" :label="setting?.name" v-model="setting.value" @change="(v:any) => settingChanged(v, setting)" />
          <el-input v-else-if="setting?.setting_type === 'password'" :label="setting?.name" v-model="setting.value" show-password @change="(v:any) => settingChanged(v, setting)" />
          <el-switch v-else-if="setting?.setting_type === 'switch'" :label="setting?.name" v-model="setting.value" @change="(v:any) => settingChanged(v ? 'true' : 'false', setting)" />
        </div>
      </el-col>
    </template>
  </el-row>
</template>
