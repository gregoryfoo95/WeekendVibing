import axios from 'axios';
import {
  User,
  Task,
  DailyTask,
  Achievement,
  UserAchievement,
  CreateUserRequest,
  GenerateDailyTasksRequest,
  UnlockAchievementRequest
} from '../types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const apiClient = axios.create({
  baseURL: `${API_BASE_URL}/api`,
  headers: {
    'Content-Type': 'application/json',
  },
});

// User API
export const userAPI = {
  create: (userData: CreateUserRequest): Promise<User> =>
    apiClient.post('/users', userData).then(res => res.data),
  
  getById: (id: number): Promise<User> =>
    apiClient.get(`/users/${id}`).then(res => res.data),
  
  update: (id: number, userData: Partial<User>): Promise<void> =>
    apiClient.put(`/users/${id}`, userData).then(res => res.data),
};

// Task API
export const taskAPI = {
  getAll: (): Promise<Task[]> =>
    apiClient.get('/tasks').then(res => res.data),
  
  getDailyTasks: (userId: number): Promise<DailyTask[]> =>
    apiClient.get(`/tasks/daily/${userId}`).then(res => res.data),
  
  generateDailyTasks: (request: GenerateDailyTasksRequest): Promise<void> =>
    apiClient.post('/tasks/daily', request).then(res => res.data),
  
  completeTask: (taskId: number): Promise<{ message: string; points_earned: number }> =>
    apiClient.put(`/tasks/daily/${taskId}/complete`).then(res => res.data),
};

// Achievement API
export const achievementAPI = {
  getAll: (): Promise<Achievement[]> =>
    apiClient.get('/achievements').then(res => res.data),
  
  getUserAchievements: (userId: number): Promise<UserAchievement[]> =>
    apiClient.get(`/achievements/user/${userId}`).then(res => res.data),
  
  unlock: (request: UnlockAchievementRequest): Promise<void> =>
    apiClient.post('/achievements/unlock', request).then(res => res.data),
};

// Leaderboard API
export const leaderboardAPI = {
  getTop: (limit: number = 10): Promise<User[]> =>
    apiClient.get(`/leaderboard?limit=${limit}`).then(res => res.data),
};

export default apiClient; 