import React, { useState, useEffect } from 'react';
import {
  Container,
  Typography,
  Button,
  Card,
  CardContent,
  Box,
  Grid,
  LinearProgress,
  Chip,
  Avatar,
  Alert,
  Snackbar,
} from '@mui/material';
import { motion } from 'framer-motion';
import { 
  FitnessCenter,
  CheckCircle,
  TrendingUp,
  Stars,
  Timer,
  EmojiEvents
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { userAPI, taskAPI } from '../api/client';
import { User, DailyTask } from '../types';

const Dashboard: React.FC = () => {
  const [userId] = useState(parseInt(localStorage.getItem('userId') || '1'));
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' as 'success' | 'error' });
  const queryClient = useQueryClient();

  // Fetch user data
  const { data: user, isLoading: userLoading } = useQuery<User>(
    ['user', userId],
    () => userAPI.getById(userId),
    { enabled: !!userId }
  );

  // Fetch daily tasks
  const { data: dailyTasks, isLoading: tasksLoading } = useQuery<DailyTask[]>(
    ['dailyTasks', userId],
    () => taskAPI.getDailyTasks(userId),
    { enabled: !!userId }
  );

  // Generate daily tasks mutation
  const generateTasksMutation = useMutation(
    () => taskAPI.generateDailyTasks({ user_id: userId }),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['dailyTasks', userId]);
        setSnackbar({ open: true, message: 'New daily tasks generated! 🎯', severity: 'success' });
      },
      onError: () => {
        setSnackbar({ open: true, message: 'Failed to generate tasks', severity: 'error' });
      }
    }
  );

  // Complete task mutation
  const completeTaskMutation = useMutation(
    (taskId: number) => taskAPI.completeTask(taskId),
    {
      onSuccess: (data) => {
        queryClient.invalidateQueries(['dailyTasks', userId]);
        queryClient.invalidateQueries(['user', userId]);
        setSnackbar({ 
          open: true, 
          message: `Great job! You earned ${data.points_earned} points! 🌟`, 
          severity: 'success' 
        });
      },
      onError: () => {
        setSnackbar({ open: true, message: 'Failed to complete task', severity: 'error' });
      }
    }
  );

  // Auto-generate tasks if none exist
  useEffect(() => {
    if (dailyTasks && dailyTasks.length === 0) {
      generateTasksMutation.mutate();
    }
  }, [dailyTasks]);

  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case 'easy': return '#4CAF50';
      case 'medium': return '#FF9800';
      case 'hard': return '#F44336';
      default: return '#9E9E9E';
    }
  };

  const getCategoryIcon = (category: string) => {
    switch (category) {
      case 'cardio': return '❤️';
      case 'strength': return '💪';
      case 'flexibility': return '🤸';
      case 'wellness': return '🧘';
      default: return '🏃';
    }
  };

  const getCharacterLevel = (points: number) => {
    if (points >= 2000) return 5;
    if (points >= 1000) return 4;
    if (points >= 500) return 3;
    if (points >= 250) return 2;
    return 1;
  };

  const getNextLevelPoints = (level: number) => {
    const levels = [0, 250, 500, 1000, 2000];
    return levels[level] || 2000;
  };

  if (userLoading || tasksLoading) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
          <Typography variant="h6">Loading your hero dashboard... 🦸‍♂️</Typography>
        </Box>
      </Container>
    );
  }

  if (!user) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Alert severity="error">User not found. Please register first.</Alert>
      </Container>
    );
  }

  const completedTasks = dailyTasks?.filter(task => task.completed).length || 0;
  const totalTasks = dailyTasks?.length || 0;
  const progressPercentage = totalTasks > 0 ? (completedTasks / totalTasks) * 100 : 0;
  const currentLevel = getCharacterLevel(user.points);
  const nextLevelPoints = getNextLevelPoints(currentLevel);
  const levelProgress = currentLevel < 5 ? ((user.points % nextLevelPoints) / nextLevelPoints) * 100 : 100;

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
      >
        {/* Header */}
        <Box mb={4}>
          <Typography variant="h3" gutterBottom>
            Welcome back, {user.username}! 🦸‍♂️
          </Typography>
          <Typography variant="h6" color="text.secondary">
            Ready for today's fitness adventure?
          </Typography>
        </Box>

        {/* Stats Grid */}
        <Grid container spacing={3} mb={4}>
          {/* Character Card */}
          <Grid item xs={12} md={4}>
            <Card sx={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', color: 'white' }}>
              <CardContent sx={{ textAlign: 'center', p: 3 }}>
                <Avatar sx={{ width: 80, height: 80, mx: 'auto', mb: 2, fontSize: '2rem' }}>
                  🦸‍♂️
                </Avatar>
                <Typography variant="h6" gutterBottom>
                  {user.character}
                </Typography>
                <Typography variant="body2" sx={{ mb: 2, opacity: 0.9 }}>
                  {user.job_title}
                </Typography>
                <Box sx={{ mb: 2 }}>
                  <Typography variant="h4" sx={{ fontWeight: 'bold' }}>
                    Level {currentLevel}
                  </Typography>
                  <LinearProgress 
                    variant="determinate" 
                    value={levelProgress} 
                    sx={{ 
                      mt: 1, 
                      height: 8, 
                      borderRadius: 4,
                      backgroundColor: 'rgba(255,255,255,0.3)',
                      '& .MuiLinearProgress-bar': {
                        backgroundColor: '#FFD93D'
                      }
                    }} 
                  />
                  <Typography variant="caption" sx={{ opacity: 0.9 }}>
                    {user.points} / {nextLevelPoints} XP
                  </Typography>
                </Box>
              </CardContent>
            </Card>
          </Grid>

          {/* Points Card */}
          <Grid item xs={12} md={4}>
            <Card sx={{ height: '100%' }}>
              <CardContent sx={{ textAlign: 'center', p: 3 }}>
                <Stars sx={{ fontSize: 48, color: '#FFD93D', mb: 2 }} />
                <Typography variant="h4" color="primary" gutterBottom>
                  {user.points}
                </Typography>
                <Typography variant="h6" gutterBottom>
                  Total Points
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Use points to unlock achievements and upgrades!
                </Typography>
              </CardContent>
            </Card>
          </Grid>

          {/* Progress Card */}
          <Grid item xs={12} md={4}>
            <Card sx={{ height: '100%' }}>
              <CardContent sx={{ textAlign: 'center', p: 3 }}>
                <TrendingUp sx={{ fontSize: 48, color: '#4ECDC4', mb: 2 }} />
                <Typography variant="h4" color="secondary" gutterBottom>
                  {completedTasks}/{totalTasks}
                </Typography>
                <Typography variant="h6" gutterBottom>
                  Today's Progress
                </Typography>
                <LinearProgress 
                  variant="determinate" 
                  value={progressPercentage} 
                  sx={{ mt: 2, height: 8, borderRadius: 4 }} 
                />
                <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                  {Math.round(progressPercentage)}% Complete
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>

        {/* Daily Tasks */}
        <Card sx={{ mb: 4 }}>
          <CardContent>
            <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
              <Typography variant="h5" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <FitnessCenter color="primary" />
                Today's Challenges
              </Typography>
              <Button
                variant="outlined"
                onClick={() => generateTasksMutation.mutate()}
                disabled={generateTasksMutation.isLoading}
                startIcon={<Timer />}
              >
                {generateTasksMutation.isLoading ? 'Generating...' : 'New Tasks'}
              </Button>
            </Box>

            <Grid container spacing={2}>
              {dailyTasks?.map((dailyTask) => (
                <Grid item xs={12} md={6} lg={4} key={dailyTask.id}>
                  <motion.div
                    whileHover={{ scale: 1.02 }}
                    whileTap={{ scale: 0.98 }}
                  >
                    <Card 
                      sx={{ 
                        height: '100%',
                        opacity: dailyTask.completed ? 0.7 : 1,
                        border: dailyTask.completed ? '2px solid #4CAF50' : 'none',
                        position: 'relative'
                      }}
                    >
                      {dailyTask.completed && (
                        <CheckCircle 
                          sx={{ 
                            position: 'absolute', 
                            top: 8, 
                            right: 8, 
                            color: '#4CAF50',
                            zIndex: 1
                          }} 
                        />
                      )}
                      <CardContent>
                        <Box display="flex" alignItems="center" gap={1} mb={2}>
                          <Typography variant="h6" sx={{ fontSize: '1.5rem' }}>
                            {getCategoryIcon(dailyTask.task.category)}
                          </Typography>
                          <Typography variant="h6" sx={{ flexGrow: 1 }}>
                            {dailyTask.task.title}
                          </Typography>
                        </Box>
                        
                        <Typography variant="body2" color="text.secondary" paragraph>
                          {dailyTask.task.description}
                        </Typography>

                        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                          <Chip
                            label={dailyTask.task.difficulty}
                            size="small"
                            sx={{
                              backgroundColor: getDifficultyColor(dailyTask.task.difficulty),
                              color: 'white',
                              fontWeight: 'bold'
                            }}
                          />
                          <Typography variant="body2" sx={{ fontWeight: 'bold', color: '#FFD93D' }}>
                            +{dailyTask.task.points} XP
                          </Typography>
                        </Box>

                        <Button
                          fullWidth
                          variant="contained"
                          disabled={dailyTask.completed || completeTaskMutation.isLoading}
                          onClick={() => completeTaskMutation.mutate(dailyTask.id)}
                          startIcon={dailyTask.completed ? <CheckCircle /> : <EmojiEvents />}
                          sx={{
                            backgroundColor: dailyTask.completed ? '#4CAF50' : undefined,
                            '&:hover': {
                              backgroundColor: dailyTask.completed ? '#4CAF50' : undefined,
                            }
                          }}
                        >
                          {dailyTask.completed ? 'Completed! 🎉' : 'Complete Task'}
                        </Button>
                      </CardContent>
                    </Card>
                  </motion.div>
                </Grid>
              ))}
            </Grid>

            {(!dailyTasks || dailyTasks.length === 0) && (
              <Box textAlign="center" py={4}>
                <Typography variant="h6" color="text.secondary">
                  No tasks for today. Generate some new challenges! 🎯
                </Typography>
              </Box>
            )}
          </CardContent>
        </Card>
      </motion.div>

      {/* Snackbar for notifications */}
      <Snackbar
        open={snackbar.open}
        autoHideDuration={4000}
        onClose={() => setSnackbar({ ...snackbar, open: false })}
      >
        <Alert 
          severity={snackbar.severity} 
          onClose={() => setSnackbar({ ...snackbar, open: false })}
          sx={{ width: '100%' }}
        >
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Container>
  );
};

export default Dashboard; 