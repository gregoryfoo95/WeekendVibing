import React from 'react';
import {
  Container,
  Typography,
  Card,
  CardContent,
  Box,
  Avatar,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  Chip,
  LinearProgress,
} from '@mui/material';
import { motion } from 'framer-motion';
import { EmojiEvents, Star, TrendingUp } from '@mui/icons-material';
import { useQuery } from 'react-query';
import { leaderboardAPI } from '../api/client';
import { User } from '../types';

const Leaderboard: React.FC = () => {
  const { data: topUsers, isLoading } = useQuery<User[]>(
    'leaderboard',
    () => leaderboardAPI.getTop(20)
  );

  const getRankEmoji = (rank: number) => {
    switch (rank) {
      case 1: return '🥇';
      case 2: return '🥈';
      case 3: return '🥉';
      default: return '🏅';
    }
  };

  const getRankColor = (rank: number) => {
    switch (rank) {
      case 1: return '#FFD700';
      case 2: return '#C0C0C0';
      case 3: return '#CD7F32';
      default: return '#9E9E9E';
    }
  };

  if (isLoading) {
    return (
      <Container maxWidth="md" sx={{ py: 4 }}>
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
          <Typography variant="h6">Loading leaderboard... 🏆</Typography>
        </Box>
      </Container>
    );
  }

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
            <EmojiEvents sx={{ fontSize: '3rem', color: '#FFD93D' }} />
            Leaderboard
          </Typography>
          <Typography variant="h6" color="text.secondary">
            See how you stack up against other fitness heroes!
          </Typography>
        </Box>

        {/* Top 3 Podium */}
        {topUsers && topUsers.length >= 3 && (
          <Card sx={{ mb: 4, background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', color: 'white' }}>
            <CardContent>
              <Box display="flex" justifyContent="center" alignItems="end" gap={2} py={2}>
                {/* 2nd Place */}
                <motion.div
                  initial={{ opacity: 0, y: 30 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.5, delay: 0.1 }}
                >
                  <Box textAlign="center">
                    <Avatar sx={{ width: 60, height: 60, mx: 'auto', mb: 1, bgcolor: getRankColor(2) }}>
                      🥈
                    </Avatar>
                    <Typography variant="h6">{topUsers[1].username}</Typography>
                    <Typography variant="body2">{topUsers[1].points} XP</Typography>
                    <Box sx={{ width: 80, height: 60, bgcolor: getRankColor(2), mt: 1, borderRadius: 1 }} />
                  </Box>
                </motion.div>

                {/* 1st Place */}
                <motion.div
                  initial={{ opacity: 0, y: 50 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.5, delay: 0.2 }}
                >
                  <Box textAlign="center">
                    <Avatar sx={{ width: 80, height: 80, mx: 'auto', mb: 1, bgcolor: getRankColor(1) }}>
                      🥇
                    </Avatar>
                    <Typography variant="h5">{topUsers[0].username}</Typography>
                    <Typography variant="body1">{topUsers[0].points} XP</Typography>
                    <Box sx={{ width: 80, height: 80, bgcolor: getRankColor(1), mt: 1, borderRadius: 1 }} />
                  </Box>
                </motion.div>

                {/* 3rd Place */}
                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.5, delay: 0.3 }}
                >
                  <Box textAlign="center">
                    <Avatar sx={{ width: 50, height: 50, mx: 'auto', mb: 1, bgcolor: getRankColor(3) }}>
                      🥉
                    </Avatar>
                    <Typography variant="h6">{topUsers[2].username}</Typography>
                    <Typography variant="body2">{topUsers[2].points} XP</Typography>
                    <Box sx={{ width: 80, height: 40, bgcolor: getRankColor(3), mt: 1, borderRadius: 1 }} />
                  </Box>
                </motion.div>
              </Box>
            </CardContent>
          </Card>
        )}

        {/* Full Leaderboard */}
        <Card>
          <CardContent>
            <Typography variant="h5" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <TrendingUp color="primary" />
              Full Rankings
            </Typography>
            
            <List>
              {topUsers?.map((user, index) => (
                <motion.div
                  key={user.id}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ duration: 0.3, delay: index * 0.05 }}
                >
                  <ListItem
                    sx={{
                      borderRadius: 2,
                      mb: 1,
                      bgcolor: index < 3 ? 'action.hover' : 'transparent',
                      border: index < 3 ? `2px solid ${getRankColor(index + 1)}` : 'none',
                    }}
                  >
                    <ListItemAvatar>
                      <Avatar sx={{ bgcolor: getRankColor(index + 1) }}>
                        {getRankEmoji(index + 1)}
                      </Avatar>
                    </ListItemAvatar>
                    
                    <ListItemText
                      primary={
                        <Box display="flex" alignItems="center" gap={2}>
                          <Typography variant="h6" sx={{ minWidth: 40 }}>
                            #{index + 1}
                          </Typography>
                          <Typography variant="h6" sx={{ flexGrow: 1 }}>
                            {user.username}
                          </Typography>
                          <Chip
                            label={user.character}
                            size="small"
                            color="primary"
                            variant="outlined"
                          />
                        </Box>
                      }
                      secondary={
                        <Box mt={1}>
                          <Box display="flex" justifyContent="space-between" alignItems="center" mb={1}>
                            <Typography variant="body2" color="text.secondary">
                              {user.job_title}
                            </Typography>
                            <Box display="flex" alignItems="center" gap={1}>
                              <Star sx={{ color: '#FFD93D', fontSize: 16 }} />
                              <Typography variant="body2" sx={{ fontWeight: 'bold' }}>
                                {user.points} XP
                              </Typography>
                            </Box>
                          </Box>
                          <LinearProgress
                            variant="determinate"
                            value={Math.min((user.points / (topUsers[0]?.points || 1)) * 100, 100)}
                            sx={{ 
                              height: 6, 
                              borderRadius: 3,
                              bgcolor: 'action.hover',
                            }}
                          />
                        </Box>
                      }
                    />
                  </ListItem>
                </motion.div>
              ))}
            </List>

            {(!topUsers || topUsers.length === 0) && (
              <Box textAlign="center" py={8}>
                <Typography variant="h6" color="text.secondary">
                  No heroes on the leaderboard yet. Be the first! 🚀
                </Typography>
              </Box>
            )}
          </CardContent>
        </Card>
      </motion.div>
    </Container>
  );
};

export default Leaderboard; 