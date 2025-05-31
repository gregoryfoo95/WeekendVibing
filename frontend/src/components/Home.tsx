import React, { useState } from 'react';
import {
  Container,
  Typography,
  Button,
  Card,
  CardContent,
  TextField,
  Box,
  Grid,
  Avatar,
} from '@mui/material';
import { motion } from 'framer-motion';
import { FitnessCenter, TrendingUp, Group, EmojiEvents } from '@mui/icons-material';
import { userAPI } from '../api/client';
import { useNavigate } from 'react-router-dom';

const Home: React.FC = () => {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [isRegistering, setIsRegistering] = useState(false);
  const navigate = useNavigate();

  const handleRegister = async () => {
    if (!username || !email) return;
    
    try {
      const user = await userAPI.create({ username, email });
      localStorage.setItem('userId', user.id.toString());
      navigate('/dashboard');
    } catch (error) {
      console.error('Registration failed:', error);
    }
  };

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
            Welcome to FitHero! 🦸‍♂️
          </Typography>
          <Typography variant="h5" color="text.secondary" paragraph sx={{ maxWidth: '600px', mx: 'auto' }}>
            Transform your fitness journey into an epic adventure! Complete daily challenges, 
            level up your character, and become the superhero of your own story.
          </Typography>
        </Box>

        {/* Registration Card */}
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
                <Typography variant="body2" color="text.secondary">
                  Create your account to begin your fitness adventure
                </Typography>
              </Box>
              
              <Box component="form" sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                <TextField
                  fullWidth
                  label="Hero Name (Username)"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  placeholder="Enter your superhero name"
                />
                <TextField
                  fullWidth
                  label="Email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="your.email@example.com"
                />
                <Button
                  variant="contained"
                  size="large"
                  onClick={handleRegister}
                  disabled={!username || !email || isRegistering}
                  sx={{ mt: 2, py: 1.5 }}
                >
                  {isRegistering ? 'Creating Hero...' : 'Begin Adventure! 🚀'}
                </Button>
              </Box>
            </CardContent>
          </Card>
        </motion.div>

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
        <Box textAlign="center" mt={6}>
          <Typography variant="h6" gutterBottom>
            Ready to transform your life? 💪
          </Typography>
          <Typography variant="body1" color="text.secondary" paragraph>
            Join thousands of heroes on their fitness journey. It's time to level up!
          </Typography>
        </Box>
      </motion.div>
    </Container>
  );
};

export default Home; 