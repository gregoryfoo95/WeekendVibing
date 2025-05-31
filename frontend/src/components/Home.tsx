import React from 'react';
import {
  Container,
  Typography,
  Button,
  Card,
  CardContent,
  Box,
  Grid,
  Avatar,
  Chip,
} from '@mui/material';
import { motion } from 'framer-motion';
import { 
  FitnessCenter, 
  TrendingUp, 
  Group, 
  EmojiEvents, 
  Dashboard as DashboardIcon,
  Login as LoginIcon 
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const Home: React.FC = () => {
  const navigate = useNavigate();
  const { user, isAuthenticated, login, isLoading } = useAuth();

  const features = [
    {
      icon: <FitnessCenter sx={{ fontSize: 40, color: '#FF6B6B' }} />,
      title: 'Daily Challenges',
      description: 'Get personalized fitness tasks every day to keep you moving and motivated!',
    },
    {
      icon: <TrendingUp sx={{ fontSize: 40, color: '#4ECDC4' }} />,
      title: 'Level Up System',
      description: 'Earn points and level up your superhero character as you complete tasks.',
    },
    {
      icon: <EmojiEvents sx={{ fontSize: 40, color: '#FFD93D' }} />,
      title: 'Achievements',
      description: 'Unlock badges, character upgrades, and job advancements with your points.',
    },
    {
      icon: <Group sx={{ fontSize: 40, color: '#6BCF7F' }} />,
      title: 'Community',
      description: 'Compete with others on the leaderboard and share your fitness journey.',
    },
  ];

  const AuthenticatedWelcome = () => (
    <motion.div
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.5, delay: 0.2 }}
    >
      <Card sx={{ maxWidth: 500, mx: 'auto', mb: 6 }}>
        <CardContent sx={{ p: 4 }}>
          <Box textAlign="center" mb={3}>
            <Avatar 
              src={user?.picture} 
              sx={{ width: 80, height: 80, mx: 'auto', mb: 2 }}
            >
              {(user?.first_name?.[0] || user?.username?.[0] || 'H').toUpperCase()}
            </Avatar>
            <Typography variant="h5" gutterBottom>
              Welcome back, {user?.first_name || user?.username}! 🎉
            </Typography>
            <Typography variant="body1" color="text.secondary" paragraph>
              Ready to continue your hero journey?
            </Typography>
            
            <Box sx={{ display: 'flex', justifyContent: 'center', gap: 1, mb: 3 }}>
              <Chip
                label={`Level ${user?.level || 1}`}
                color="primary"
                sx={{ fontWeight: 'bold' }}
              />
              <Chip
                label={user?.character || 'Rookie Hero'}
                color="secondary"
              />
              <Chip
                label={`${user?.points || 0} Points`}
                variant="outlined"
              />
            </Box>
          </Box>
          
          <Button
            variant="contained"
            size="large"
            startIcon={<DashboardIcon />}
            onClick={() => navigate('/dashboard')}
            fullWidth
            sx={{ py: 1.5 }}
          >
            Go to Dashboard
          </Button>
        </CardContent>
      </Card>
    </motion.div>
  );

  const UnauthenticatedWelcome = () => (
    <motion.div
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.5, delay: 0.2 }}
    >
      <Card sx={{ maxWidth: 400, mx: 'auto', mb: 6 }}>
        <CardContent sx={{ p: 4 }}>
          <Box textAlign="center" mb={3}>
            <Avatar sx={{ width: 80, height: 80, mx: 'auto', mb: 2, bgcolor: 'primary.main' }}>
              🦸‍♂️
            </Avatar>
            <Typography variant="h6" gutterBottom>
              Start Your Hero Journey
            </Typography>
            <Typography variant="body2" color="text.secondary" paragraph>
              Sign in with Google to begin your fitness adventure and track your progress!
            </Typography>
          </Box>
          
          <Button
            variant="contained"
            size="large"
            startIcon={<LoginIcon />}
            onClick={login}
            disabled={isLoading}
            fullWidth
            sx={{ py: 1.5 }}
          >
            {isLoading ? 'Loading...' : 'Sign In with Google 🚀'}
          </Button>
          
          <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mt: 2, textAlign: 'center' }}>
            Secure authentication powered by Google OAuth
          </Typography>
        </CardContent>
      </Card>
    </motion.div>
  );

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
      >
        {/* Hero Section */}
        <Box textAlign="center" mb={6}>
          <Typography variant="h1" gutterBottom sx={{ 
            background: 'linear-gradient(45deg, #FF6B6B 30%, #4ECDC4 90%)',
            backgroundClip: 'text',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent',
            mb: 2
          }}>
            {isAuthenticated ? `Welcome back to FitHero! 🦸‍♂️` : 'Welcome to FitHero! 🦸‍♂️'}
          </Typography>
          <Typography variant="h5" color="text.secondary" paragraph sx={{ maxWidth: '600px', mx: 'auto' }}>
            {isAuthenticated 
              ? 'Your fitness adventure continues! Check your dashboard for today\'s challenges.'
              : 'Transform your fitness journey into an epic adventure! Complete daily challenges, level up your character, and become the superhero of your own story.'
            }
          </Typography>
        </Box>

        {/* Authentication-based Welcome Section */}
        {isAuthenticated ? <AuthenticatedWelcome /> : <UnauthenticatedWelcome />}

        {/* Features Grid */}
        <Grid container spacing={4} sx={{ mt: 4 }}>
          {features.map((feature, index) => (
            <Grid item xs={12} sm={6} md={3} key={index}>
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: 0.3 + index * 0.1 }}
              >
                <Card sx={{ height: '100%', textAlign: 'center' }}>
                  <CardContent sx={{ p: 3 }}>
                    <Box mb={2}>
                      {feature.icon}
                    </Box>
                    <Typography variant="h6" gutterBottom>
                      {feature.title}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {feature.description}
                    </Typography>
                  </CardContent>
                </Card>
              </motion.div>
            </Grid>
          ))}
        </Grid>

        {/* Call to Action */}
        {!isAuthenticated && (
          <Box textAlign="center" mt={6}>
            <Typography variant="h6" gutterBottom>
              Ready to transform your life? 💪
            </Typography>
            <Typography variant="body1" color="text.secondary" paragraph>
              Join thousands of heroes on their fitness journey. It's time to level up!
            </Typography>
            <Button
              variant="outlined"
              size="large"
              startIcon={<LoginIcon />}
              onClick={login}
              sx={{ mt: 2, px: 4, py: 1.5 }}
            >
              Get Started Today
            </Button>
          </Box>
        )}
      </motion.div>
    </Container>
  );
};

export default Home; 