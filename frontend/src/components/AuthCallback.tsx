import React, { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Box, CircularProgress, Typography, Alert } from '@mui/material';
import { useAuth } from '../contexts/AuthContext';

const AuthCallback: React.FC = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { checkAuth } = useAuth();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const handleCallback = async () => {
      try {
        // The actual OAuth handling is done by the backend
        // We just need to check if the user is now authenticated
        await checkAuth();
        
        // If authentication was successful, redirect to dashboard
        navigate('/dashboard');
      } catch (err) {
        console.error('Authentication error:', err);
        setError('Authentication failed. Please try again.');
        setTimeout(() => navigate('/'), 3000);
      }
    };

    handleCallback();
  }, [searchParams, navigate, checkAuth]);

  if (error) {
    return (
      <Box
        display="flex"
        flexDirection="column"
        alignItems="center"
        justifyContent="center"
        minHeight="60vh"
        gap={2}
        px={2}
      >
        <Alert severity="error" sx={{ maxWidth: 400 }}>
          {error}
        </Alert>
        <Typography variant="body2" color="text.secondary">
          Redirecting to home page...
        </Typography>
      </Box>
    );
  }

  return (
    <Box
      display="flex"
      flexDirection="column"
      alignItems="center"
      justifyContent="center"
      minHeight="60vh"
      gap={2}
    >
      <CircularProgress size={60} />
      <Typography variant="h6" color="text.secondary">
        Completing authentication...
      </Typography>
      <Typography variant="body2" color="text.secondary">
        Please wait while we set up your account
      </Typography>
    </Box>
  );
};

export default AuthCallback; 