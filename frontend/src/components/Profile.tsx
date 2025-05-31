import React, { useState } from 'react';
import {
  Container,
  Typography,
  Card,
  CardContent,
  Box,
  Avatar,
  Grid,
  Chip,
  LinearProgress,
} from '@mui/material';
import { motion } from 'framer-motion';
import { Person, Star, TrendingUp, EmojiEvents } from '@mui/icons-material';
import { useQuery } from 'react-query';
import { userAPI, achievementAPI } from '../api/client';
import { User, UserAchievement } from '../types';

const Profile: React.FC = () => {
  const [userId] = useState(parseInt(localStorage.getItem('userId') || '1'));

  const { data: user, isLoading } = useQuery<User>(
    ['user', userId],
    () => userAPI.getById(userId),
    { enabled: !!userId }
  );

  const { data: userAchievements } = useQuery<UserAchievement[]>(
    ['userAchievements', userId],
    () => achievementAPI.getUserAchievements(userId),
    { enabled: !!userId }
  );

  const getCharacterLevel = (points: number) => {
    if (points >= 2000) return 5;
    if (points >= 1000) return 4;
    if (points >= 500) return 3;
    if (points >= 250) return 2;
    return 1;
  };

  if (isLoading || !user) {
    return (
      <Container maxWidth="md" sx={{ py: 4 }}>
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
          <Typography variant="h6">Loading profile... 👤</Typography>
        </Box>
      </Container>
    );
  }

  const currentLevel = getCharacterLevel(user.points);

  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
      >
        {/* Header */}
        <Box mb={4} textAlign="center">
          <Typography variant="h3" gutterBottom sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 2 }}>
            <Person sx={{ fontSize: '3rem', color: '#FF6B6B' }} />
            Hero Profile
          </Typography>
        </Box>

        {/* Profile Card */}
        <Card sx={{ mb: 4, background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', color: 'white' }}>
          <CardContent sx={{ textAlign: 'center', p: 4 }}>
            <Avatar sx={{ width: 120, height: 120, mx: 'auto', mb: 3, fontSize: '3rem' }}>
              🦸‍♂️
            </Avatar>
            <Typography variant="h4" gutterBottom>
              {user.username}
            </Typography>
            <Typography variant="h6" sx={{ mb: 2, opacity: 0.9 }}>
              {user.character}
            </Typography>
            <Typography variant="body1" sx={{ mb: 3, opacity: 0.8 }}>
              {user.job_title}
            </Typography>
            <Chip
              label={`Level ${currentLevel}`}
              sx={{ 
                bgcolor: '#FFD93D', 
                color: '#000', 
                fontWeight: 'bold',
                fontSize: '1rem',
                px: 2,
                py: 1
              }}
            />
          </CardContent>
        </Card>

        {/* Stats Grid */}
        <Grid container spacing={3} mb={4}>
          <Grid item xs={12} md={4}>
            <Card>
              <CardContent sx={{ textAlign: 'center' }}>
                <Star sx={{ fontSize: 48, color: '#FFD93D', mb: 2 }} />
                <Typography variant="h4" color="primary" gutterBottom>
                  {user.points}
                </Typography>
                <Typography variant="h6">Total Points</Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} md={4}>
            <Card>
              <CardContent sx={{ textAlign: 'center' }}>
                <TrendingUp sx={{ fontSize: 48, color: '#4ECDC4', mb: 2 }} />
                <Typography variant="h4" color="secondary" gutterBottom>
                  {currentLevel}
                </Typography>
                <Typography variant="h6">Current Level</Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} md={4}>
            <Card>
              <CardContent sx={{ textAlign: 'center' }}>
                <EmojiEvents sx={{ fontSize: 48, color: '#FFD93D', mb: 2 }} />
                <Typography variant="h4" color="primary" gutterBottom>
                  {userAchievements?.length || 0}
                </Typography>
                <Typography variant="h6">Achievements</Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>

        {/* Recent Achievements */}
        <Card>
          <CardContent>
            <Typography variant="h5" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <EmojiEvents color="primary" />
              Recent Achievements
            </Typography>
            
            {userAchievements && userAchievements.length > 0 ? (
              <Grid container spacing={2}>
                {userAchievements.slice(0, 6).map((userAchievement) => (
                  <Grid item xs={12} sm={6} md={4} key={userAchievement.id}>
                    <Card variant="outlined">
                      <CardContent sx={{ textAlign: 'center', p: 2 }}>
                        <Typography variant="h4" sx={{ mb: 1 }}>
                          {userAchievement.achievement.icon}
                        </Typography>
                        <Typography variant="h6" gutterBottom>
                          {userAchievement.achievement.title}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {userAchievement.achievement.description}
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>
                ))}
              </Grid>
            ) : (
              <Box textAlign="center" py={4}>
                <Typography variant="h6" color="text.secondary">
                  No achievements unlocked yet. Start completing tasks to earn your first achievement! 🎯
                </Typography>
              </Box>
            )}
          </CardContent>
        </Card>
      </motion.div>
    </Container>
  );
};

export default Profile; 