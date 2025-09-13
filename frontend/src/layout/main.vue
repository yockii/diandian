<script setup lang="ts">
import { WindowService } from '../../bindings/changeme/background/service'
const close = async () => {
    await WindowService.HideMainAndShowFloating()
};
</script>

<template>
  <v-app theme="light" class="draggable">
    <v-app-bar elevation="0" height="32" color="rgba(0,0,0,0)" dense flat>
      <template v-slot:append>
        <v-btn variant="plain" size="x-small" @click="close">
          <v-icon icon="fas fa-xmark" ></v-icon>
        </v-btn>
      </template>
    </v-app-bar>
    <v-main>
      <router-view></router-view>
    </v-main>
    <v-footer color="indigo-lighten-1" height="16">ccc</v-footer>
  </v-app>
</template>

<style lang="scss" scoped>
.v-footer {
  flex: 0;
}
</style>
