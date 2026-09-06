const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

if (!API_BASE_URL && typeof window !== 'undefined') {
  console.warn('WARNING: NEXT_PUBLIC_API_BASE_URL is not defined. Defaulting to http://localhost:8080');
}

const BASE_URL = API_BASE_URL || 'http://localhost:8080';

export async function apiRequest(endpoint: string, options: RequestInit = {}) {
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  const response = await fetch(`${BASE_URL}${endpoint}`, {
    ...options,
    credentials: 'include',
    headers,
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'Something went wrong');
  }

  if (response.status === 204) return null;
  return response.json();
}
