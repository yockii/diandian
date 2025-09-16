<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { SettingService } from '../../bindings/diandian/background/service';
import { Setting } from '../../bindings/diandian/background/model/models';
import SettingGroup from '../components/SettingGroup.vue';
import DianDivider from '@/components/DianDivider.vue';

const settingMap = ref<Map<string, Setting[]>>(new Map());

const fetchSettings = async () => {
  const settings = await SettingService.AllSettings();
  // 分组放入
  settings.forEach(setting => {
    if (setting === null) {
      return;
    }
    const group = setting.group_name || '未分组';
    if (!settingMap.value.has(group)) {
      settingMap.value.set(group, []);
    }
    settingMap.value.get(group)?.push(setting);
  });
};

onMounted(() => {
  fetchSettings();
});
</script>

<template>
  <div class="flex flex-col h-full">
    <div class="text-3xl font-bold -ml-4 -mt-12 absolute">设置</div>
    <div class="flex-fill pa-2 d-flex flex-column">
      <template v-for="([groupName, settingList]) in settingMap" :key="groupName">
        <dian-divider position="start">{{ groupName }}</dian-divider>
        <setting-group :setting-list="settingList" class="mb-8"/>
      </template>
    </div>
  </div>
</template>
