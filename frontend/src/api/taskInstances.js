import axios from 'axios';
import { buildApiUrl } from '@/config/api';

/**
 * Fetches all task instances.
 * @returns {Promise<Array<Object>>} An array of task instances.
 */
export async function fetchTaskInstances() {
  const response = await axios.get(buildApiUrl('/tasks'));
  return response.data;
}
