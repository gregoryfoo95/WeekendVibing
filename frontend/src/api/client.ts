import axios from 'axios';
import {
  User,
  Task,
  DailyTask,
  Achievement,
  UserAchievement,
  CreateUserRequest,
  GenerateDailyTasksRequest,
  UnlockAchievementRequest,
  AuthUser,
  AuthResponse
} from '../types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const apiClient = axios.create({
  baseURL: `${API_BASE_URL}/api`,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true, // Enable cookies for authentication
});

// Add request interceptor to handle authentication
apiClient.interceptors.request.use(
  (config) => {
    // The JWT token will be automatically sent via HTTP-only cookies
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add response interceptor to handle authentication errors
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Redirect to login or handle unauthorized access
      window.location.href = '/';
    }
    return Promise.reject(error);
  }
);

// Authentication API
export const authAPI = {
  // Initiate Google OAuth login
  googleLogin: async (): Promise<void> => {
    try {
      const response = await apiClient.get('/auth/google');
      const authUrl = response.data.auth_url;
      window.location.href = authUrl;
    } catch (error) {
      console.error('Failed to initiate Google login:', error);
      throw error;
    }
  },

  // Check authentication status
  checkAuth: (): Promise<AuthUser> =>
    apiClient.get('/auth/check').then(res => res.data.user),

  // Get current user profile
  me: (): Promise<AuthUser> =>
    apiClient.get('/me').then(res => res.data.user),

  // Logout
  logout: (): Promise<void> =>
    apiClient.post('/auth/logout').then(res => res.data),

  // Refresh token
  refreshToken: (): Promise<AuthResponse> =>
    apiClient.post('/auth/refresh').then(res => res.data),
};

// User API
export const userAPI = {
  create: (userData: CreateUserRequest): Promise<User> =>
    apiClient.post('/admin/users', userData).then(res => res.data),
  
  getById: (id: number): Promise<User> =>
    apiClient.get(`/users/${id}`).then(res => res.data),
  
  update: (id: number, userData: Partial<User>): Promise<void> =>
    apiClient.put(`/users/${id}`, userData).then(res => res.data),

  // Get current user profile
  getCurrentProfile: (): Promise<User> =>
    apiClient.get('/profile').then(res => res.data),

  // Update current user profile
  updateCurrentProfile: (userData: Partial<User>): Promise<void> =>
    apiClient.put('/profile', userData).then(res => res.data),

  // Get current user's tasks
  getCurrentUserTasks: (): Promise<DailyTask[]> =>
    apiClient.get('/users/tasks').then(res => res.data),

  // Get current user's achievements
  getCurrentUserAchievements: (): Promise<UserAchievement[]> =>
    apiClient.get('/users/achievements').then(res => res.data),
};

// Task API
export const taskAPI = {
  getAll: (): Promise<Task[]> =>
    apiClient.get('/public/tasks').then(res => res.data),
  
  getDailyTasks: (): Promise<DailyTask[]> =>
    apiClient.get('/tasks/daily').then(res => res.data.tasks),
  
  generateDailyTasks: (): Promise<DailyTask[]> =>
    apiClient.post('/tasks/daily/generate').then(res => res.data.tasks),
  
  completeTask: (taskId: number): Promise<{ message: string; points_earned: number }> =>
    apiClient.post(`/tasks/daily/${taskId}/complete`).then(res => res.data),
};

// Achievement API
export const achievementAPI = {
  getAll: (): Promise<Achievement[]> =>
    apiClient.get('/public/achievements').then(res => res.data),
  
  getUserAchievements: (): Promise<UserAchievement[]> =>
    apiClient.get('/achievements/user').then(res => res.data.achievements),
  
  unlock: (achievementId: number): Promise<void> =>
    apiClient.post(`/achievements/${achievementId}/unlock`).then(res => res.data),
};

// Leaderboard API (assuming this will be added to backend)
export const leaderboardAPI = {
  getTop: (limit: number = 10): Promise<User[]> =>
    apiClient.get(`/public/leaderboard?limit=${limit}`).then(res => res.data),
};

export default apiClient; 