<template>
  <div class="task-table-container">
    <h2>Task Instances</h2>
    <table>
      <thead>
        <tr>
          <th class="title-col">Title</th>
          <th class="status-col">Completed</th>
          <th class="date-col">Created At</th>
          <th class="date-col">Completed At</th>
          <th class="action-col">Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="task in tasks" :key="task.id" :class="{ 'completed-task': task.completed }">
          <td class="title-col">{{ task.title }}</td>
          <td class="status-col">{{ task.completed ? "✅" : "❌" }}</td>
          <td class="date-col">{{ formatDate(task.created_at) }}</td>
          <td class="date-col">{{ task.completed ? formatDate(task.completed_at) : '-' }}</td>
          <td class="action-col">
            <button
              v-if="!task.completed"
              @click="markCompleted(task.id)"
              :disabled="loadingId === task.id"
              class="action-button"
            >
              {{ loadingId === task.id ? 'Completing...' : 'Mark as Completed' }}
            </button>
            <span v-else class="completed-text">Completed</span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { completeTaskInstance } from '@/api/completeTask';
import type { TaskInstance } from '@/types/taskInstance';

const loadingId = ref<string | null>(null);

const props = defineProps<{
  tasks: TaskInstance[];
}>();

function formatDate(dateString: string) {
  if (!dateString) return '-';
  return new Date(dateString).toLocaleString();
}

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
.task-table-container {
  background-color: #2e2e2e;
  color: #ddd;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #444;
  width: 100%;
  margin: 0 auto;
  box-sizing: border-box;
}

h2 {
  margin-top: 0;
  margin-bottom: 1rem;
  color: #ddd;
}

.table-responsive {
  width: 100%;
  overflow-x: auto;
}

table {
  width: 100%;
  min-width: 750px; /* Ensure a minimum width for the table */
  border-collapse: collapse;
  border-radius: 4px;
  overflow: hidden;
}

th {
  background-color: #444;
  color: #fff;
  font-weight: 500;
  text-align: left;
}

th, td {
  padding: 0.75rem;
  border: 1px solid #555;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Set specific column width classes */
.title-col {
  width: 25%;
  min-width: 150px;
}

.status-col {
  width: 10%;
  min-width: 90px;
  text-align: center;
}

.date-col {
  width: 20%;
  min-width: 160px;
}

.action-col {
  width: 25%;
  min-width: 160px;
}

tr:nth-child(even) {
  background-color: #333;
}

tr:hover {
  background-color: #3a3a3a;
}

.completed-task {
  opacity: 0.7;
}

.completed-text {
  color: #8bc34a;
  font-style: italic;
}

.action-button {
  background-color: #3a7bd5;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background-color 0.3s;
  width: 100%;
  max-width: 160px;
  box-sizing: border-box;
}

.action-button:hover {
  background-color: #2d62aa;
}

.action-button:disabled {
  background-color: #666;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .action-button {
    max-width: 100%;
  }
}
</style>
