<script setup lang="ts">
const close = async () => {
};
</script>

<template>
  <v-app theme="light" class="draggable colored">
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

<style scoped>
.colored {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}
.v-footer {
  flex: 0;
}
</style>
