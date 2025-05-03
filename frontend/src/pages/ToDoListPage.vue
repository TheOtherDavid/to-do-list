<template>
  <div class="todo-page">
    <h1>To-Do List</h1>
    <div class="tab-container">
      <button 
        :class="['tab-button', { active: activeTab === 'instances' }]"
        @click="activeTab = 'instances'"
      >
        Task Instances
      </button>
      <button 
        :class="['tab-button', { active: activeTab === 'templates' }]"
        @click="activeTab = 'templates'"
      >
        Task Templates
      </button>
    </div>
    <div class="content-container">
      <div class="todo-list-container">
        <div class="todo-controls">
          <button 
            v-if="activeTab === 'instances'"
            @click="showCompleted = !showCompleted"
            class="toggle-button"
          >
            {{ showCompleted ? 'Hide Completed Tasks' : 'Show Completed Tasks' }}
          </button>
        </div>
        <TaskInstanceTable 
          v-if="activeTab === 'instances'" 
          :tasks="displayedTasks" 
        />
        <TaskTemplateTable 
          v-else 
          :templates="taskTemplates" 
        />
      </div>
      <div class="form-container">
        <CreateTaskForm 
          v-if="activeTab === 'instances'" 
          @task-created="handleTaskCreated" 
        />
        <CreateTaskTemplateForm 
          v-else 
          @task-template-created="handleTaskTemplateCreated" 
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { fetchTaskInstances } from '@/api/taskInstances';
import { fetchTaskTemplates } from '@/api/taskTemplates';
import TaskInstanceTable from '@/components/TaskInstanceTable.vue';
import CreateTaskForm from '@/components/CreateTaskForm.vue';
import CreateTaskTemplateForm from '@/components/CreateTaskTemplateForm.vue';
import TaskTemplateTable from '@/components/TaskTemplateTable.vue';
import type { TaskInstance } from '@/types/taskInstance';
import type { TaskTemplate } from '@/types/taskTemplate';


const taskInstances = ref<TaskInstance[]>([]);
const taskTemplates = ref<TaskTemplate[]>([]);
const showCompleted = ref(false);
const activeTab = ref('instances');

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
  taskTemplates.value = await fetchTaskTemplates();
}

function handleTaskCreated(newTask: TaskInstance) {
  // Add the new task to the beginning of the list
  taskInstances.value = [newTask, ...taskInstances.value];
}

function handleTaskTemplateCreated(newTaskTemplate: TaskTemplate) {
  // Add the new task template to the beginning of the list
  taskTemplates.value = [newTaskTemplate, ...taskTemplates.value];
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

.tab-container {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  padding: 0.5rem;
  background-color: #2d2d2d;
  border-radius: 8px;
}

.tab-button {
  padding: 0.5rem 1rem;
  border: none;
  background: none;
  color: #888;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s;
}

.tab-button.active {
  background-color: #4CAF50;
  color: white;
}

.tab-button:hover:not(.active) {
  color: #ddd;
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