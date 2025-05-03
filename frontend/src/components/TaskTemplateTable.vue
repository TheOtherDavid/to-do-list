<template>
  <div class="task-table-container">
    <h2>Task Templates</h2>
    <table>
      <thead>
        <tr>
          <th class="title-col">Title</th>
          <th class="recurrence-col">Recurrence Rule</th>
          <th class="date-col">Created At</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="template in templates" :key="template.id">
          <td class="title-col">{{ template.title }}</td>
          <td class="recurrence-col">{{ template.recurrence_rule }}</td>
          <td class="date-col">{{ formatDate(template.created_at) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import type { TaskTemplate } from '@/types/taskTemplate';

interface Props {
  templates: TaskTemplate[];
}

defineProps<Props>();

function formatDate(dateString: string) {
  if (!dateString) return '-';
  return new Date(dateString).toLocaleString();
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

tr:nth-child(even) {
  background-color: #333;
}

tr:hover {
  background-color: #3a3a3a;
}

</style>