export interface User {
  id: number;
  username: string;
  email: string;
  level: number;
  points: number;
  character: string;
  job_title: string;
  created_at: string;
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
  completed: boolean;
  date: string;
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