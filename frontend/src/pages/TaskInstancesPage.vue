<template>
  <div class="task-page">
    <h1>Task Instances</h1>
    <div class="content-container">
      <div class="task-list-container">
        <TaskInstanceTable :tasks="taskInstances" />
      </div>
      <div class="form-container">
        <CreateTaskForm @task-created="handleTaskCreated" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { fetchTaskInstances } from '@/api/taskInstances';
import TaskInstanceTable from '@/components/TaskInstanceTable.vue';
import CreateTaskForm from '@/components/CreateTaskForm.vue';
import type { TaskInstance } from '@/types/taskInstance';

const taskInstances = ref<TaskInstance[]>([]);

async function loadTasks() {
  taskInstances.value = await fetchTaskInstances();
}

function handleTaskCreated(newTask: TaskInstance) {
  // Add the new task to the beginning of the list
  taskInstances.value = [newTask, ...taskInstances.value];
}

onMounted(async () => {
  await loadTasks();
});
</script>

<style scoped>
.task-page {
  display: flex;
  gap: 2rem;
  align-items: flex-start;
  padding: 2rem;
}

.table-container {
  flex: 2; /* or 1 */
}

.create-task-form {
  flex: 1;
  max-width: 500px;
}

@media (min-width: 768px) {
  .content-container {
    flex-direction: row;
    align-items: flex-start;
    justify-content: space-between;
  }
  
  .task-list-container {
    flex: 1 1 65%;
  }
  
  .form-container {
    flex: 1;
    display: flex;
    justify-content: center;
    align-items: flex-start;
    padding: 1rem 0;
  }
}

@media (max-width: 767px) {
  .task-list-container {
    padding-bottom: 2rem;
  }
  
  .form-container {
    width: 100%;
  }
}
</style>