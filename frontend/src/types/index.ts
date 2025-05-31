export interface User {
  id: number;
  username: string;
  email: string;
  level: number;
  points: number;
  character: string;
  job_title: string;
  google_id?: string;
  first_name?: string;
  last_name?: string;
  picture?: string;
  is_active: boolean;
  last_login_at?: string;
  created_at: string;
  updated_at: string;
}

export interface AuthUser {
  id: number;
  username: string;
  email: string;
  first_name?: string;
  last_name?: string;
  picture?: string;
  level: number;
  points: number;
  character: string;
  job_title: string;
  is_active: boolean;
  last_login_at?: string;
}

export interface AuthResponse {
  user: AuthUser;
  token: string;
  expires_at: string;
}

export interface AuthContextType {
  user: AuthUser | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: () => Promise<void>;
  logout: () => void;
  checkAuth: () => Promise<void>;
}

export interface Task {
  id: number;
  title: string;
  description: string;
  points: number;
  category: string;
  difficulty: 'easy' | 'medium' | 'hard';
}

export interface DailyTask {
  id: number;
  user_id: number;
  task_id: number;
  task: Task;
  is_completed: boolean;
  points: number;
  created_at: string;
  updated_at: string;
}

export interface Achievement {
  id: number;
  title: string;
  description: string;
  icon: string;
  points_cost: number;
  type: 'character' | 'upgrade' | 'badge';
}

export interface UserAchievement {
  id: number;
  user_id: number;
  achievement_id: number;
  achievement: Achievement;
  unlocked_at: string;
}

export interface CreateUserRequest {
  username: string;
  email: string;
}

export interface GenerateDailyTasksRequest {
  user_id: number;
}

export interface UnlockAchievementRequest {
  user_id: number;
  achievement_id: number;
} 