import React from 'react';
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
import { useUserData } from '../hooks';

const Profile: React.FC = () => {
  const {
    userProfile,
    userAchievements,
    currentPoints,
    isLoading
  } = useUserData();

  const getCharacterLevel = (points: number) => {
    if (points >= 2000) return 5;
    if (points >= 1000) return 4;
    if (points >= 500) return 3;
    if (points >= 250) return 2;
    return 1;
  };

  if (isLoading) {
    return (
      <Container maxWidth="md" sx={{ py: 4 }}>
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
          <Typography variant="h6">Loading your profile... 👤</Typography>
        </Box>
      </Container>
    );
  }

  if (!userProfile) {
    return (
      <Container maxWidth="md" sx={{ py: 4 }}>
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
          <Typography variant="h6">Authentication required... 👤</Typography>
        </Box>
      </Container>
    );
  }

  const currentLevel = getCharacterLevel(currentPoints);

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
            <Person sx={{ fontSize: '3rem', color: '#4ECDC4' }} />
            Hero Profile
          </Typography>
          <Typography variant="h6" color="text.secondary">
            Your fitness journey at a glance
          </Typography>
        </Box>

        {/* Profile Card */}
        <Card sx={{ mb: 4, background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', color: 'white' }}>
          <CardContent sx={{ textAlign: 'center', p: 4 }}>
            <Avatar 
              src={userProfile.picture}
              sx={{ width: 120, height: 120, mx: 'auto', mb: 3, fontSize: '3rem' }}
            >
              {(userProfile.first_name?.[0] || userProfile.username?.[0] || '🦸').toUpperCase()}
            </Avatar>
            
            <Typography variant="h4" gutterBottom>
              {userProfile.first_name || userProfile.username}
            </Typography>
            
            <Box display="flex" justifyContent="center" gap={2} mb={3}>
              <Chip
                label={userProfile.character}
                sx={{ 
                  backgroundColor: 'rgba(255,255,255,0.2)', 
                  color: 'white',
                  fontWeight: 'bold'
                }}
              />
              <Chip
                label={`Level ${currentLevel}`}
                sx={{ 
                  backgroundColor: '#FFD93D', 
                  color: '#333',
                  fontWeight: 'bold'
                }}
              />
            </Box>
            
            <Typography variant="h6" sx={{ opacity: 0.9 }}>
              {userProfile.job_title}
            </Typography>
            
            {userProfile.last_login_at && (
              <Typography variant="body2" sx={{ mt: 2, opacity: 0.8 }}>
                Last login: {new Date(userProfile.last_login_at).toLocaleDateString()}
              </Typography>
            )}
          </CardContent>
        </Card>

        {/* Stats Grid */}
        <Grid container spacing={3} mb={4}>
          <Grid item xs={12} md={4}>
            <Card>
              <CardContent sx={{ textAlign: 'center' }}>
                <Star sx={{ fontSize: 48, color: '#FFD93D', mb: 2 }} />
                <Typography variant="h4" color="primary" gutterBottom>
                  {currentPoints}
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
                        <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mt: 1 }}>
                          Unlocked: {new Date(userAchievement.unlocked_at).toLocaleDateString()}
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