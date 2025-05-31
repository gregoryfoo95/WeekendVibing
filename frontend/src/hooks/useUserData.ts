import { useQuery, useMutation, useQueryClient } from 'react-query';
import { authAPI, achievementAPI, taskAPI } from '../api/client';
import { AuthUser, UserAchievement, DailyTask } from '../types';
import { useAuth } from '../contexts/AuthContext';

export const useUserData = () => {
  const { user } = useAuth();
  const queryClient = useQueryClient();

  // User profile data
  const {
    data: userProfile,
    isLoading: profileLoading,
    refetch: refetchProfile
  } = useQuery<AuthUser>(
    ['userProfile', user?.id],
    () => authAPI.me(),
    {
      enabled: !!user?.id,
      staleTime: 1000 * 60 * 5, // 5 minutes
    }
  );

  // User achievements
  const {
    data: userAchievements,
    isLoading: achievementsLoading,
    refetch: refetchAchievements
  } = useQuery<UserAchievement[]>(
    ['userAchievements', user?.id],
    () => achievementAPI.getUserAchievements(),
    {
      enabled: !!user?.id,
      staleTime: 1000 * 60 * 2, // 2 minutes
    }
  );

  // Daily tasks
  const {
    data: dailyTasks,
    isLoading: tasksLoading,
    refetch: refetchTasks
  } = useQuery<DailyTask[]>(
    ['dailyTasks', user?.id],
    () => taskAPI.getDailyTasks(),
    {
      enabled: !!user?.id,
      staleTime: 1000 * 30, // 30 seconds
    }
  );

  // Update user profile in cache after point changes
  const updateUserPoints = (newPoints: number) => {
    if (userProfile) {
      queryClient.setQueryData(['userProfile', user?.id], {
        ...userProfile,
        points: newPoints
      });
    }
  };

  // Invalidate all user-related data
  const invalidateAllUserData = () => {
    queryClient.invalidateQueries(['userProfile', user?.id]);
    queryClient.invalidateQueries(['userAchievements', user?.id]);
    queryClient.invalidateQueries(['dailyTasks', user?.id]);
    queryClient.invalidateQueries('leaderboard');
  };

  // Invalidate specific data types
  const invalidateUserData = {
    profile: () => queryClient.invalidateQueries(['userProfile', user?.id]),
    achievements: () => queryClient.invalidateQueries(['userAchievements', user?.id]),
    tasks: () => queryClient.invalidateQueries(['dailyTasks', user?.id]),
    leaderboard: () => queryClient.invalidateQueries('leaderboard'),
    all: invalidateAllUserData
  };

  return {
    // Data
    userProfile: userProfile || user,
    userAchievements: userAchievements || [],
    dailyTasks: dailyTasks || [],
    
    // Loading states
    isLoading: profileLoading || achievementsLoading || tasksLoading,
    profileLoading,
    achievementsLoading,
    tasksLoading,
    
    // Actions
    refetchProfile,
    refetchAchievements,
    refetchTasks,
    updateUserPoints,
    invalidateUserData,
    
    // Computed values
    completedTasks: dailyTasks?.filter(task => task.is_completed).length || 0,
    totalTasks: dailyTasks?.length || 0,
    achievementCount: userAchievements?.length || 0,
    currentPoints: userProfile?.points || user?.points || 0,
  };
}; 