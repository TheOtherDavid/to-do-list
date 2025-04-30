import axios from 'axios';
import { buildApiUrl } from '@/config/api';

/**
 * Creates a new task instance.
 * @param {Object} taskData - The task data to create.
 * @returns {Promise<Object>} The created task instance.
 */
export async function createTaskInstance(taskData) {
  const response = await axios.post(buildApiUrl('/tasks'), taskData);
  return response.data;
}
