import axios from 'axios';
import { buildApiUrl } from '@/config/api';

/**
 * Fetches all task templates.
 * @returns {Promise<Array<Object>>} An array of task templates.
 */
export async function fetchTaskTemplates() {
  const response = await axios.get(buildApiUrl('/task-templates'));
  return response.data;
}
