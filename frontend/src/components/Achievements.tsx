import React, { useState } from 'react';
import {
  Container,
  Typography,
  Button,
  Card,
  CardContent,
  Box,
  Grid,
  Tabs,
  Tab,
  Chip,
  Avatar,
  Alert,
  Snackbar,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from '@mui/material';
import { motion } from 'framer-motion';
import { 
  EmojiEvents,
  Star,
  Lock,
  CheckCircle,
  PersonPin,
  Work,
  Shield
} from '@mui/icons-material';
import { useQuery } from 'react-query';
import { achievementAPI } from '../api/client';
import { Achievement, UserAchievement } from '../types';
import { useUserData, useAchievementActions } from '../hooks';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`achievement-tabpanel-${index}`}
      aria-labelledby={`achievement-tab-${index}`}
      {...other}
    >
      {value === index && <Box>{children}</Box>}
    </div>
  );
}

const Achievements: React.FC = () => {
  const [tabValue, setTabValue] = useState(0);
  const [selectedAchievement, setSelectedAchievement] = useState<Achievement | null>(null);
  const [confirmDialog, setConfirmDialog] = useState(false);

  const {
    userProfile,
    userAchievements,
    currentPoints
  } = useUserData();

  const {
    unlockAchievement,
    isUnlockingAchievement,
    notification,
    hideNotification
  } = useAchievementActions();

  // Fetch all achievements
  const { data: achievements } = useQuery<Achievement[]>(
    'achievements',
    achievementAPI.getAll
  );

  const handleUnlock = (achievement: Achievement) => {
    setSelectedAchievement(achievement);
    setConfirmDialog(true);
  };

  const confirmUnlock = () => {
    if (selectedAchievement) {
      unlockAchievement(selectedAchievement.id);
      setConfirmDialog(false);
      setSelectedAchievement(null);
    }
  };

  const isUnlocked = (achievementId: number) => {
    return userAchievements?.some(ua => ua.achievement_id === achievementId) || false;
  };

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'character': return <PersonPin />;
      case 'upgrade': return <Work />;
      case 'badge': return <Shield />;
      default: return <EmojiEvents />;
    }
  };

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'character': return '#FF6B6B';
      case 'upgrade': return '#4ECDC4';
      case 'badge': return '#FFD93D';
      default: return '#9E9E9E';
    }
  };

  const filteredAchievements = achievements?.filter(achievement => {
    switch (tabValue) {
      case 0: return true; // All
      case 1: return achievement.type === 'character';
      case 2: return achievement.type === 'upgrade';
      case 3: return achievement.type === 'badge';
      default: return true;
    }
  }) || [];

  const unlockedCount = userAchievements?.length || 0;
  const totalCount = achievements?.length || 0;

  if (!userProfile) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Alert severity="error">Authentication required. Please sign in.</Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
      >
        {/* Header */}
        <Box mb={4}>
          <Typography variant="h3" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            <EmojiEvents sx={{ fontSize: '3rem', color: '#FFD93D' }} />
            Achievement Store
          </Typography>
          <Typography variant="h6" color="text.secondary" paragraph>
            Spend your hard-earned points to unlock new characters, job titles, and badges!
          </Typography>
          
          {/* Progress Summary */}
          <Card sx={{ mb: 3, background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', color: 'white' }}>
            <CardContent sx={{ textAlign: 'center' }}>
              <Grid container spacing={3}>
                <Grid item xs={12} md={4}>
                  <Typography variant="h4">{currentPoints}</Typography>
                  <Typography variant="body1">Available Points</Typography>
                </Grid>
                <Grid item xs={12} md={4}>
                  <Typography variant="h4">{unlockedCount}</Typography>
                  <Typography variant="body1">Unlocked</Typography>
                </Grid>
                <Grid item xs={12} md={4}>
                  <Typography variant="h4">{totalCount}</Typography>
                  <Typography variant="body1">Total</Typography>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Box>

        {/* Tabs */}
        <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
          <Tabs value={tabValue} onChange={(_, newValue) => setTabValue(newValue)}>
            <Tab label="All Achievements" />
            <Tab label="Characters" />
            <Tab label="Job Titles" />
            <Tab label="Badges" />
          </Tabs>
        </Box>

        {/* Achievement Grid */}
        <TabPanel value={tabValue} index={tabValue}>
          <Grid container spacing={3}>
            {filteredAchievements.map((achievement) => {
              const unlocked = isUnlocked(achievement.id);
              const canAfford = currentPoints >= achievement.points_cost;
              
              return (
                <Grid item xs={12} sm={6} md={4} key={achievement.id}>
                  <motion.div
                    whileHover={{ scale: unlocked ? 1 : 1.05 }}
                    whileTap={{ scale: 0.95 }}
                  >
                    <Card 
                      sx={{ 
                        height: '100%',
                        position: 'relative',
                        border: unlocked ? '2px solid #4CAF50' : 'none',
                        opacity: unlocked ? 0.8 : 1,
                        background: unlocked ? 'linear-gradient(135deg, #e8f5e8 0%, #f1f8e9 100%)' : 'white',
                      }}
                    >
                      {unlocked && (
                        <CheckCircle 
                          sx={{ 
                            position: 'absolute', 
                            top: 12, 
                            right: 12, 
                            color: '#4CAF50',
                            zIndex: 1
                          }} 
                        />
                      )}
                      
                      <CardContent sx={{ textAlign: 'center', p: 3 }}>
                        {/* Icon */}
                        <Avatar 
                          sx={{ 
                            width: 80, 
                            height: 80, 
                            mx: 'auto', 
                            mb: 2, 
                            fontSize: '2.5rem',
                            bgcolor: getTypeColor(achievement.type),
                            filter: unlocked ? 'grayscale(50%)' : 'none'
                          }}
                        >
                          {achievement.icon}
                        </Avatar>

                        {/* Type Badge */}
                        <Box display="flex" justifyContent="center" mb={2}>
                          <Chip
                            icon={getTypeIcon(achievement.type)}
                            label={achievement.type.charAt(0).toUpperCase() + achievement.type.slice(1)}
                            size="small"
                            sx={{
                              backgroundColor: getTypeColor(achievement.type),
                              color: 'white',
                              fontWeight: 'bold'
                            }}
                          />
                        </Box>

                        {/* Title & Description */}
                        <Typography variant="h6" gutterBottom>
                          {achievement.title}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" paragraph>
                          {achievement.description}
                        </Typography>

                        {/* Points Cost */}
                        <Typography variant="h6" color="primary" gutterBottom>
                          {achievement.points_cost} Points
                        </Typography>

                        {/* Action Button */}
                        <Button
                          fullWidth
                          variant={unlocked ? "outlined" : "contained"}
                          disabled={unlocked || !canAfford || isUnlockingAchievement}
                          onClick={() => handleUnlock(achievement)}
                          sx={{
                            backgroundColor: unlocked ? undefined : 
                              canAfford ? undefined : '#ccc'
                          }}
                        >
                          {unlocked ? 'Unlocked! ✓' : 
                           canAfford ? 'Unlock' : 
                           'Not Enough Points'}
                        </Button>
                      </CardContent>
                    </Card>
                  </motion.div>
                </Grid>
              );
            })}
          </Grid>
        </TabPanel>

        {/* Confirmation Dialog */}
        <Dialog open={confirmDialog} onClose={() => setConfirmDialog(false)}>
          <DialogTitle>Unlock Achievement?</DialogTitle>
          <DialogContent>
            <Typography paragraph>
              Are you sure you want to unlock "{selectedAchievement?.title}" for {selectedAchievement?.points_cost} points?
            </Typography>
            <Typography variant="body2" color="text.secondary">
              You currently have {currentPoints} points.
            </Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setConfirmDialog(false)}>Cancel</Button>
            <Button onClick={confirmUnlock} variant="contained" disabled={isUnlockingAchievement}>
              {isUnlockingAchievement ? 'Unlocking...' : 'Unlock'}
            </Button>
          </DialogActions>
        </Dialog>

        {/* Snackbar for notifications */}
        <Snackbar
          open={notification.open}
          autoHideDuration={4000}
          onClose={hideNotification}
        >
          <Alert 
            severity={notification.severity} 
            onClose={hideNotification}
            sx={{ width: '100%' }}
          >
            {notification.message}
          </Alert>
        </Snackbar>
      </motion.div>
    </Container>
  );
};

export default Achievements; 