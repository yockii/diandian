export const EVENT_NAMES = {
  THEME_CHANGED: 'theme-changed',
  CAN_WORK_CHANGED: 'can-work-changed',
  STICKY_SIDE_CHANGED: 'sticky-side-changed',
  MOUSE_ENTER_FLOATING: 'mouse-enter-floating',
  MOUSE_LEAVE_FLOATING: 'mouse-leave-floating',

  MESSAGE_RESPONSED: "message-responsed",
  TASK_STATUS_CHANGED: "task-status-changed",
  OPERATE_FAILED: "operate-failed",

  NOTIFY: "notify",

  // 任务执行相关事件
  TASK_EXECUTION_STARTED: "task-execution-started",
  TASK_EXECUTION_COMPLETED: "task-execution-completed",
} as const
