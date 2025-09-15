<script lang="ts" setup>
import { computed } from 'vue';

const props = defineProps({
  position: {
    type: String,
    default: 'center',
    validator: (value: string) => ['start', 'center', 'end'].includes(value)
  },
  lineHeight: {
    type: Number || String,
    default: 1,
  },
  lineColor: {
    type: String,
    default: '#e0e0e0',
  },
  // start/end中，长的线条占比
  longLineFlex: {
    type: Number,
    default: 10,
  },
})

const useLineHeight = computed(() => {
  return typeof props.lineHeight === 'number' ? `${props.lineHeight}px` : props.lineHeight
})
</script>

<template>
  <div class="dian-divider" :class="{
    'dian-divider--start': position === 'start',
    'dian-divider--center': position === 'center',
    'dian-divider--end': position === 'end',
  }" :style="{
    '--line-height': useLineHeight,
    '--line-color': lineColor,
    '--long-line-flex': longLineFlex,
  }">
    <slot />
  </div>
</template>

<style scoped>
.dian-divider {
  display: flex;
  align-items: center;
  width: 100%;
  --line-height: 1px;
  --line-color: #e0e0e0;
  white-space: nowrap;
  padding: 0 8px;
  gap: 16px;
}

.dian-divider--start {
  justify-content: flex-start;
}
.dian-divider--start::before,
.dian-divider--start::after {
  content: '';
  flex-grow: var(--long-line-flex);
  height: var(--line-height);
  background-color: var(--line-color);
}
.dian-divider--start::before {
  flex-grow: 1;
}

.dian-divider--center::before,
.dian-divider--center::after {
  content: '';
  flex-grow: 1;
  height: var(--line-height);
  background-color: var(--line-color);
}

.dian-divider--end {
  justify-content: flex-end;
}

.dian-divider--end::before,
.dian-divider--end::after {
  content: '';
  flex-grow: var(--long-line-flex);
  height: var(--line-height);
  background-color: var(--line-color);
}
.dian-divider--end::after {
  flex-grow: 1;
}

:deep(.dian-divider > *) {
  white-space: nowrap;
  padding: 0 8px;
}
</style>
