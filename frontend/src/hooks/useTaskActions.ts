import { useMutation, useQueryClient } from 'react-query';
import { taskAPI } from '../api/client';
import { useAuth } from '../contexts/AuthContext';
import { useState } from 'react';
import { TaskCompletionResponse } from '../types';
import { useRealTimeUpdates } from './useRealTimeUpdates';

interface TaskActionResult {
  success: boolean;
  pointsEarned?: number;
  newLevel?: number;
  achievementUnlocked?: boolean;
  message: string;
}

export const useTaskActions = () => {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const { updatePointsEverywhere, updateTaskCompletion } = useRealTimeUpdates();
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

  // Complete task mutation with immediate real-time updates
  const completeTaskMutation = useMutation<TaskCompletionResponse, Error, number>(
    (taskId: number) => taskAPI.completeTask(taskId),
    {
      onMutate: async (taskId: number) => {
        // IMMEDIATE optimistic update - no waiting for server
        updateTaskCompletion(taskId, true);
        return { taskId };
      },
      onSuccess: (data, taskId) => {
        // Apply points immediately across all components
        if (data.points_earned) {
          updatePointsEverywhere(data.points_earned);
        }

        // Show success notification
        let message = `Great job! Task completed! 🎉`;
        if (data.points_earned) {
          message += ` You earned ${data.points_earned} points! ⭐`;
        }
        if (data.level_up) {
          message += ` Level up! You're now level ${data.new_level}! 🆙`;
        }
        if (data.achievement_unlocked) {
          message += ` New achievement unlocked! 🏆`;
        }

        showNotification(message, 'success');
      },
      onError: (error: any, taskId) => {
        // Rollback optimistic update on error
        updateTaskCompletion(taskId, false);
        const message = error.response?.data?.error || 'Failed to complete task';
        showNotification(message, 'error');
      }
    }
  );

  // Generate daily tasks mutation
  const generateTasksMutation = useMutation(
    () => taskAPI.generateDailyTasks(),
    {
      onSuccess: () => {
        // Invalidate tasks to show new ones
        queryClient.invalidateQueries(['dailyTasks', user?.id]);
        showNotification('New daily tasks generated! Get ready for your fitness adventure! 🎯', 'success');
      },
      onError: (error: any) => {
        const message = error.response?.data?.error || 'Failed to generate tasks';
        showNotification(message, 'error');
      }
    }
  );

  return {
    // Mutations
    completeTask: completeTaskMutation.mutate,
    generateTasks: generateTasksMutation.mutate,
    
    // Loading states
    isCompletingTask: completeTaskMutation.isLoading,
    isGeneratingTasks: generateTasksMutation.isLoading,
    
    // Notification state
    notification,
    showNotification,
    hideNotification,
    
    // Action results
    lastTaskResult: completeTaskMutation.data,
  };
}; 