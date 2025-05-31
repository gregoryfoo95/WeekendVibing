import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { AuthUser, AuthContextType } from '../types';
import { authAPI } from '../api/client';

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const checkAuth = async () => {
    try {
      setIsLoading(true);
      const userData = await authAPI.checkAuth();
      setUser(userData);
    } catch (error) {
      setUser(null);
      console.log('Not authenticated');
    } finally {
      setIsLoading(false);
    }
  };

  const login = async () => {
    try {
      await authAPI.googleLogin();
    } catch (error) {
      console.error('Login failed:', error);
      // Could show a toast notification here
    }
  };

  const logout = async () => {
    try {
      await authAPI.logout();
      setUser(null);
      // Redirect to home page after logout
      window.location.href = '/';
    } catch (error) {
      console.error('Logout error:', error);
      // Even if logout fails on server, clear local state
      setUser(null);
      window.location.href = '/';
    }
  };

  useEffect(() => {
    checkAuth();
  }, []);

  const value: AuthContextType = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    logout,
    checkAuth,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export default AuthContext; 