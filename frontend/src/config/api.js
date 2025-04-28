/**
 * API Configuration
 * 
 * This file contains centralized configuration for API endpoints.
 * Change these values when deploying to different environments.
 */

// For local development
const LOCAL_API_URL = 'http://192.168.178.59:8080';

// For network access (using your computer's IP)
// const NETWORK_API_URL = 'http://YOUR_IP_ADDRESS:8080';

// For production
// const PRODUCTION_API_URL = 'https://api.your-domain.com';

// Set the active API base URL here
export const API_BASE_URL = LOCAL_API_URL;

/**
 * Helper function to build API URLs
 * @param {string} endpoint - The API endpoint path (e.g., '/tasks')
 * @returns {string} The complete API URL
 */
export function buildApiUrl(endpoint) {
  // Make sure endpoint starts with a slash
  const normalizedEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
  return `${API_BASE_URL}${normalizedEndpoint}`;
}
