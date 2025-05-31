import { useCallback } from 'react';
import { useQueryClient } from 'react-query';
import { useAuth } from '../contexts/AuthContext';

/**
 * Hook for real-time updates across the app
 * This provides immediate state synchronization without waiting for server roundtrips
 */
export const useRealTimeUpdates = () => {
  const { user } = useAuth();
  const queryClient = useQueryClient();

  // Immediately update points across all components
  const updatePointsEverywhere = useCallback((pointsDelta: number) => {
    if (!user?.id) return;

    // Update user profile points
    queryClient.setQueryData(['userProfile', user.id], (oldData: any) => {
      if (!oldData) return oldData;
      return { ...oldData, points: (oldData.points || 0) + pointsDelta };
    });

    // Update leaderboard if user is in it
    queryClient.setQueryData('leaderboard', (oldData: any) => {
      if (!Array.isArray(oldData)) return oldData;
      return oldData.map((u: any) => 
        u.id === user.id 
          ? { ...u, points: (u.points || 0) + pointsDelta }
          : u
      ).sort((a: any, b: any) => b.points - a.points); // Re-sort by points
    });

    // Optionally invalidate for fresh data after a delay
    setTimeout(() => {
      queryClient.invalidateQueries(['userProfile', user.id]);
      queryClient.invalidateQueries('leaderboard');
    }, 5000); // Refresh from server after 5 seconds

  }, [user?.id, queryClient]);

  // Update achievement count immediately
  const updateAchievementCount = useCallback((delta: number) => {
    if (!user?.id) return;

    queryClient.setQueryData(['userProfile', user.id], (oldData: any) => {
      if (!oldData) return oldData;
      return { ...oldData, achievement_count: (oldData.achievement_count || 0) + delta };
    });
  }, [user?.id, queryClient]);

  // Update task completion status immediately
  const updateTaskCompletion = useCallback((taskId: number, completed: boolean) => {
    if (!user?.id) return;

    queryClient.setQueryData(['dailyTasks', user.id], (oldData: any) => {
      if (!Array.isArray(oldData)) return oldData;
      return oldData.map((task: any) => 
        task.id === taskId 
          ? { ...task, is_completed: completed }
          : task
      );
    });
  }, [user?.id, queryClient]);

  // Force refresh all data immediately (for critical updates)
  const forceRefreshAll = useCallback(() => {
    if (!user?.id) return;
    
    const queries = [
      ['userProfile', user.id],
      ['userAchievements', user.id], 
      ['dailyTasks', user.id],
      'leaderboard',
      'achievements'
    ];

    queries.forEach(queryKey => {
      queryClient.invalidateQueries(queryKey);
    });
  }, [user?.id, queryClient]);

  return {
    updatePointsEverywhere,
    updateAchievementCount, 
    updateTaskCompletion,
    forceRefreshAll
  };
}; 