import axios from 'axios';
import { buildApiUrl } from '@/config/api';

/**
 * Marks a task instance as completed by its ID.
 * @param {string} id - The ID of the task instance.
 */
export async function completeTaskInstance(id) {
  const response = await axios.put(buildApiUrl(`/tasks/${id}/complete`));
  return response.data;
}