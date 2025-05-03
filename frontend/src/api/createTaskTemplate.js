import axios from 'axios';
import { buildApiUrl } from '@/config/api';

/**
 * Creates a new task template.
 * @param {Object} taskTemplateData - The task template data to create.
 * @returns {Promise<Object>} The created task template.
 */
export async function createTaskTemplate(taskTemplateData) {
  const response = await axios.post(buildApiUrl('/task-templates'), taskTemplateData);
  return response.data;
}
