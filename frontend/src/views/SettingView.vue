<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { SettingService } from '../../bindings/changeme/background/service';
import { Setting } from '../../bindings/changeme/background/model/models';
import SettingGroup from '../components/SettingGroup.vue';

const settingMap = ref<Map<string, Setting[]>>(new Map());

onMounted(async () => {
  // 在组件挂载后执行的逻辑
  const settings = await SettingService.AllSettings()
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
});
</script>

<template>
  <div class="d-flex flex-column" style="height: 100%">
    <div class="text-h5 mx-2 mt-n7">设置</div>
    <div class="flex-fill pa-2 d-flex flex-column">
        <template v-for="([groupName, settingList]) in settingMap" :key="groupName">
          <v-divider :thickness="2">{{ groupName }}</v-divider>
          <setting-group :setting-list="settingList" />
        </template>
    </div>
  </div>
</template>