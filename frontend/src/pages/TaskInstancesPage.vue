<template>
  <div class="task-page">
    <h1>To-Do List</h1>
    <div class="content-container">
      <div class="task-list-container">
        <div class="task-controls">
          <button 
            @click="showCompleted = !showCompleted"
            class="toggle-button"
          >
            {{ showCompleted ? 'Hide Completed Tasks' : 'Show Completed Tasks' }}
          </button>
        </div>
        <TaskInstanceTable :tasks="displayedTasks" />
      </div>
      <div class="form-container">
        <CreateTaskForm @task-created="handleTaskCreated" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { fetchTaskInstances } from '@/api/taskInstances';
import TaskInstanceTable from '@/components/TaskInstanceTable.vue';
import CreateTaskForm from '@/components/CreateTaskForm.vue';
import type { TaskInstance } from '@/types/taskInstance';

const taskInstances = ref<TaskInstance[]>([]);
const showCompleted = ref(false);

// Computed property for displayed tasks with sorting
const displayedTasks = computed(() => {
  // Sort tasks: uncompleted first, then completed by completion date (newest first)
  const sortedTasks = [...taskInstances.value].sort((a, b) => {
    // First, separate completed from uncompleted
    if (a.completed && !b.completed) return 1;
    if (!a.completed && b.completed) return -1;
    
    // If both are completed, sort by completion date (newest first)
    if (a.completed && b.completed && a.completed_at && b.completed_at) {
      return new Date(b.completed_at).getTime() - new Date(a.completed_at).getTime();
    }
    
    // If both are uncompleted, sort by creation date (newest first)
    return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
  });

  if (showCompleted.value) {
    return sortedTasks;
  } else {
    return sortedTasks.filter(task => !task.completed);
  }
});

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
  flex-direction: column;
  min-height: 100vh;
  background-color: #1e1e1e;
  color: #ddd;
  padding: 2rem;
}

h1 {
  text-align: center;
  margin-bottom: 2rem;
  color: #fff;
}

.content-container {
  display: flex;
  flex-direction: column;
  gap: 2rem;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.task-controls {
  margin-bottom: 1rem;
  display: flex;
  justify-content: flex-end;
}

.toggle-button {
  background-color: #3a7bd5;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background-color 0.3s;
}

.toggle-button:hover {
  background-color: #2d62aa;
}

.form-container {
  width: 100%;
}

@media (min-width: 992px) {
  .content-container {
    flex-direction: row;
    align-items: flex-start;
  }
  
  .task-list-container {
    flex: 1 1 65%;
    
  }
  
  .form-container {
    flex: 0 0 350px; /* Fixed width instead of flex-grow */
    padding: 0 1.5rem;
  }
}

@media (max-width: 991px) {
  .task-list-container {
    margin-bottom: 2rem;
  }
  
  .form-container {
    max-width: 500px;
    margin: 0 auto;
  }
}
</style>