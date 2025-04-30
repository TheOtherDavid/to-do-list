<template>
  <table>
    <thead>
      <tr>
        <th>Title</th>
        <th>Description</th>
        <th>Completed</th>
        <th>Created At</th>
        <th>Action</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="task in tasks" :key="task.id">
        <td>{{ task.title }}</td>
        <td>{{ task.description }}</td>
        <td>{{ task.completed ? "✅" : "❌" }}</td>
        <td>{{ new Date(task.created_at).toLocaleString() }}</td>
        <td>
          <button
            v-if="!task.completed"
            @click="markCompleted(task.id)"
            :disabled="loadingId === task.id"
          >
            {{ loadingId === task.id ? 'Completing...' : 'Mark as Completed' }}
          </button>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { completeTaskInstance } from '@/api/completeTask';
import type { TaskInstance } from '@/types/taskInstance';

const loadingId = ref<string | null>(null);

const props = defineProps<{
  tasks: TaskInstance[];
}>();

async function markCompleted(id: string) {
  loadingId.value = id;
  try {
    await completeTaskInstance(id);
    // Find the task and mark as completed locally
    const task = props.tasks.find(t => t.id === id);
    if (task) {
      task.completed = true;
      task.completed_at = new Date().toISOString();
    }
  } finally {
    loadingId.value = null;
  }
}
</script>

<style scoped>
table {
  width: 100%;
  border-collapse: collapse;
}
th, td {
  padding: 0.5rem;
  border: 1px solid #ccc;
}
</style>
