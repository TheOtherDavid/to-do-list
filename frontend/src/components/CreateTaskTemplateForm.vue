<template>
    <div class="create-task-template-form">
      <h2>Create New Task Template</h2>
      <form @submit.prevent="submitForm">
        <div class="form-group">
          <label for="title">Title</label>
          <input 
            type="text" 
            id="title" 
            v-model="taskTemplateData.title" 
            required
            placeholder="Enter task title"
          />

          <!--Selector for DAILY, WEEKLY, FORTNIGHTLY, MONTHLY-->
          <label for="recurrence_rule">Recurrence Rule</label>
          <select 
            id="recurrence_rule" 
            v-model="taskTemplateData.recurrence_rule" 
            required
          >
            <option value="DAILY">Daily</option>
            <option value="WEEKLY">Weekly</option>
            <option value="FORTNIGHTLY">Fortnightly</option>
            <option value="MONTHLY">Monthly</option>
          </select>
        </div>
        
        <button 
          type="submit" 
          :disabled="isSubmitting"
        >
          {{ isSubmitting ? 'Creating...' : 'Create Task Template' }}
        </button>
      </form>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref } from 'vue';
  import { createTaskTemplate } from '@/api/createTaskTemplate';
  
  const emit = defineEmits(['task-template-created']);
  
  const taskTemplateData = ref({
    title: '',
    description: '',
    recurrence_rule: '',
  });
  
  const isSubmitting = ref(false);
  
  async function submitForm() {
    if (!taskTemplateData.value.title) return;
    
    isSubmitting.value = true;
    
    try {
      const newTaskTemplate = await createTaskTemplate(taskTemplateData.value);
      emit('task-template-created', newTaskTemplate);
      
      // Reset form
      taskTemplateData.value = {
        title: '',
        description: '',
        recurrence_rule: '',
      };
    } catch (error) {
      console.error('Error creating task template:', error);
      alert('Failed to create task template. Please try again.');
    } finally {
      isSubmitting.value = false;
    }
  }
  </script>
  
  
  <style scoped>
  .create-task-template-form {
    background-color: #2e2e2e;
    color: #ddd;
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    border: 1px solid #444;
    width: 100%;
    min-width: 300px;
    box-sizing: border-box;
  }
  
  .form-group {
    margin-bottom: 1.5rem;
  }
  
  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: #ddd;
  }
  
  input,
  textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #555;
    border-radius: 4px;
    font-size: 1rem;
    background-color: #444;
    color: #fff;
    box-sizing: border-box;
  }
  
  textarea {
    resize: vertical;
    min-height: 80px;
  }
  
  button {
    background-color: #3a7bd5;
    color: white;
    border: none;
    padding: 0.75rem 1.5rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
    transition: background-color 0.3s;
    width: 100%;
    box-sizing: border-box;
  }
  
  button:hover {
    background-color: #2d62aa;
  }
  
  button:disabled {
    background-color: #666;
    cursor: not-allowed;
  }
  </style>