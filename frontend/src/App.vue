<script setup lang="ts">
const close = async () => {
};
</script>

<template>
  <v-app theme="light" class="draggable">
    <v-app-bar elevation="0" height="32">
      <template v-slot:append>
        <v-btn variant="plain" size="x-small" @click="close">
          <v-icon icon="fas fa-xmark" ></v-icon>
        </v-btn>
      </template>
    </v-app-bar>
  </v-app>
</template>

<style scoped>
</style>
