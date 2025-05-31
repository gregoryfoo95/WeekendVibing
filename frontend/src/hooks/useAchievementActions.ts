import { useMutation, useQueryClient } from 'react-query';
import { achievementAPI } from '../api/client';
import { useAuth } from '../contexts/AuthContext';
import { useState } from 'react';
import { AchievementUnlockResponse } from '../types';
import { useRealTimeUpdates } from './useRealTimeUpdates';

export const useAchievementActions = () => {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const { updatePointsEverywhere, updateAchievementCount } = useRealTimeUpdates();
  const [notification, setNotification] = useState<{
    open: boolean;
    message: string;
    severity: 'success' | 'error' | 'info';
  }>({ open: false, message: '', severity: 'success' });

  const showNotification = (message: string, severity: 'success' | 'error' | 'info' = 'success') => {
    setNotification({ open: true, message, severity });
  };

  const hideNotification = () => {
    setNotification(prev => ({ ...prev, open: false }));
  };

  // Unlock achievement mutation with immediate real-time updates
  const unlockAchievementMutation = useMutation<AchievementUnlockResponse, Error, number>(
    (achievementId: number) => achievementAPI.unlock(achievementId),
    {
      onSuccess: (data, achievementId) => {
        const achievement = data.user_achievement.achievement;
        const pointsSpent = achievement.points_cost;

        // Apply points deduction immediately across all components
        updatePointsEverywhere(-pointsSpent);
        
        // Update achievement count
        updateAchievementCount(1);

        // Update character/job_title in profile immediately
        queryClient.setQueryData(['userProfile', user?.id], (oldProfile: any) => {
          if (!oldProfile) return oldProfile;
          return {
            ...oldProfile,
            ...(achievement.type === 'character' && { character: achievement.title }),
            ...(achievement.type === 'upgrade' && { job_title: achievement.title }),
          };
        });

        // Add the achievement to user achievements cache immediately
        queryClient.setQueryData(['userAchievements', user?.id], (oldAchievements: any) => {
          if (!oldAchievements) return [data.user_achievement];
          return [...oldAchievements, data.user_achievement];
        });

        // Show success notification
        let message = `Achievement unlocked! 🎉`;
        if (achievement.type === 'character') {
          message += ` You are now ${achievement.title}! 🦸‍♂️`;
        }
        if (achievement.type === 'upgrade') {
          message += ` Your new title: ${achievement.title}! 💼`;
        }
        message += ` (${pointsSpent} points spent)`;

        showNotification(message, 'success');
      },
      onError: (error: any) => {
        const message = error.response?.data?.error || 'Failed to unlock achievement';
        showNotification(message, 'error');
      }
    }
  );

  return {
    // Mutations
    unlockAchievement: unlockAchievementMutation.mutate,
    
    // Loading states
    isUnlockingAchievement: unlockAchievementMutation.isLoading,
    
    // Notification state
    notification,
    showNotification,
    hideNotification,
    
    // Achievement result
    lastUnlockResult: unlockAchievementMutation.data,
  };
}; 