const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';

export const COMMON_ID_LOGIN_URL = `${API_BASE_URL}/auth/login?back_path=/home`;
export const COMMON_ID_SIGNUP_URL = `${API_BASE_URL}/auth/signup?back_path=/home`;
export const COMMON_ID_LOGOUT_URL = `${API_BASE_URL}/auth/logout`;
