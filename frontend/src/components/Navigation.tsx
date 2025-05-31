import React from 'react';
import { AppBar, Toolbar, Typography, Button, Box } from '@mui/material';
import { Link, useLocation } from 'react-router-dom';
import { 
  Home as HomeIcon, 
  Dashboard as DashboardIcon, 
  EmojiEvents as AchievementsIcon,
  Leaderboard as LeaderboardIcon,
  Person as ProfileIcon
} from '@mui/icons-material';

const Navigation: React.FC = () => {
  const location = useLocation();

  const navigationItems = [
    { path: '/', label: 'Home', icon: <HomeIcon /> },
    { path: '/dashboard', label: 'Dashboard', icon: <DashboardIcon /> },
    { path: '/achievements', label: 'Achievements', icon: <AchievementsIcon /> },
    { path: '/leaderboard', label: 'Leaderboard', icon: <LeaderboardIcon /> },
    { path: '/profile', label: 'Profile', icon: <ProfileIcon /> },
  ];

  return (
    <AppBar position="static" elevation={0} sx={{ background: 'linear-gradient(45deg, #FF6B6B 30%, #4ECDC4 90%)' }}>
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1, fontWeight: 'bold' }}>
          🦸‍♂️ FitHero
        </Typography>
        <Box sx={{ display: 'flex', gap: 1 }}>
          {navigationItems.map((item) => (
            <Button
              key={item.path}
              component={Link}
              to={item.path}
              startIcon={item.icon}
              color="inherit"
              sx={{
                backgroundColor: location.pathname === item.path ? 'rgba(255,255,255,0.2)' : 'transparent',
                '&:hover': {
                  backgroundColor: 'rgba(255,255,255,0.1)',
                },
                borderRadius: 2,
                px: 2,
              }}
            >
              {item.label}
            </Button>
          ))}
        </Box>
      </Toolbar>
    </AppBar>
  );
};

export default Navigation; 