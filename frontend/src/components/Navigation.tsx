import React from 'react';
import { 
  AppBar, 
  Toolbar, 
  Typography, 
  Button, 
  Box, 
  Avatar, 
  IconButton, 
  Menu, 
  MenuItem, 
  Chip,
  Divider
} from '@mui/material';
import { Link, useLocation } from 'react-router-dom';
import { 
  Home as HomeIcon, 
  Dashboard as DashboardIcon, 
  EmojiEvents as AchievementsIcon,
  Leaderboard as LeaderboardIcon,
  Person as ProfileIcon,
  Login as LoginIcon,
  Logout as LogoutIcon,
  KeyboardArrowDown as ArrowDownIcon
} from '@mui/icons-material';
import { useAuth } from '../contexts/AuthContext';

const Navigation: React.FC = () => {
  const location = useLocation();
  const { user, isAuthenticated, isLoading, login, logout } = useAuth();
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

  const handleUserMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleUserMenuClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    handleUserMenuClose();
    logout();
  };

  const publicNavigationItems = [
    { path: '/', label: 'Home', icon: <HomeIcon /> },
  ];

  const authenticatedNavigationItems = [
    { path: '/', label: 'Home', icon: <HomeIcon /> },
    { path: '/dashboard', label: 'Dashboard', icon: <DashboardIcon /> },
    { path: '/achievements', label: 'Achievements', icon: <AchievementsIcon /> },
    { path: '/leaderboard', label: 'Leaderboard', icon: <LeaderboardIcon /> },
    { path: '/profile', label: 'Profile', icon: <ProfileIcon /> },
  ];

  const navigationItems = isAuthenticated ? authenticatedNavigationItems : publicNavigationItems;

  return (
    <AppBar position="static" elevation={0} sx={{ background: 'linear-gradient(45deg, #FF6B6B 30%, #4ECDC4 90%)' }}>
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1, fontWeight: 'bold' }}>
          🦸‍♂️ FitHero
        </Typography>
        
        <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
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

          {isLoading ? (
            <Box sx={{ width: 40, height: 40 }} />
          ) : isAuthenticated ? (
            <>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, ml: 2 }}>
                <Chip
                  label={`Level ${user?.level || 1}`}
                  size="small"
                  sx={{
                    backgroundColor: 'rgba(255,255,255,0.2)',
                    color: 'white',
                    fontWeight: 'bold',
                  }}
                />
                <Chip
                  label={`${user?.points || 0} pts`}
                  size="small"
                  sx={{
                    backgroundColor: 'rgba(255,255,255,0.15)',
                    color: 'white',
                  }}
                />
              </Box>
              
              <IconButton
                onClick={handleUserMenuOpen}
                sx={{ color: 'white', ml: 1 }}
              >
                <Avatar
                  src={user?.picture}
                  alt={user?.first_name || user?.username}
                  sx={{ width: 32, height: 32 }}
                >
                  {(user?.first_name?.[0] || user?.username?.[0] || 'U').toUpperCase()}
                </Avatar>
                <ArrowDownIcon sx={{ fontSize: 16, ml: 0.5 }} />
              </IconButton>

              <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleUserMenuClose}
                anchorOrigin={{
                  vertical: 'bottom',
                  horizontal: 'right',
                }}
                transformOrigin={{
                  vertical: 'top',
                  horizontal: 'right',
                }}
              >
                <Box sx={{ px: 2, py: 1, minWidth: 200 }}>
                  <Typography variant="subtitle2" sx={{ fontWeight: 'bold' }}>
                    {user?.first_name ? `${user.first_name} ${user.last_name || ''}`.trim() : user?.username}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {user?.email}
                  </Typography>
                  <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
                    {user?.character} • {user?.job_title}
                  </Typography>
                </Box>
                <Divider />
                <MenuItem component={Link} to="/profile" onClick={handleUserMenuClose}>
                  <ProfileIcon sx={{ mr: 1 }} />
                  Profile
                </MenuItem>
                <MenuItem onClick={handleLogout}>
                  <LogoutIcon sx={{ mr: 1 }} />
                  Logout
                </MenuItem>
              </Menu>
            </>
          ) : (
            <Button
              onClick={login}
              startIcon={<LoginIcon />}
              color="inherit"
              variant="outlined"
              sx={{
                borderColor: 'rgba(255,255,255,0.5)',
                '&:hover': {
                  borderColor: 'white',
                  backgroundColor: 'rgba(255,255,255,0.1)',
                },
                borderRadius: 2,
                px: 3,
              }}
            >
              Sign In with Google
            </Button>
          )}
        </Box>
      </Toolbar>
    </AppBar>
  );
};

export default Navigation; 