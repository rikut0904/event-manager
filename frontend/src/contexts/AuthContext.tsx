'use client';

import React, { createContext, useContext, useState, useEffect } from 'react';
import { apiRequest } from '@/lib/api';
import { COMMON_ID_LOGIN_URL, COMMON_ID_LOGOUT_URL, COMMON_ID_SIGNUP_URL } from '@/lib/auth';

interface User {
  id: string;
  email: string;
  name?: string;
  username?: string;
  avatar_url?: string;
  bio?: string;
  connpass_id?: string;
}

interface AuthContextType {
  user: User | null;
  login: (email: string, password: string) => Promise<void>;
  signUp: (email: string, password: string) => Promise<void>;
  logout: () => void;
  linkConnpass: (connpassID: string) => Promise<void>;
  loading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    apiRequest('/api/v1/users/me')
      .then((currentUser) => setUser(currentUser))
      .catch(() => setUser(null))
      .finally(() => setLoading(false));
  }, []);

  const login = async (_email: string, _password: string) => {
    window.location.href = COMMON_ID_LOGIN_URL;
  };

  const signUp = async (_email: string, _password: string) => {
    window.location.href = COMMON_ID_SIGNUP_URL;
  };

  const logout = () => {
    setUser(null);
    window.location.href = COMMON_ID_LOGOUT_URL;
  };

  const linkConnpass = async (connpassID: string) => {
    await apiRequest('/auth/link-connpass', {
      method: 'POST',
      body: JSON.stringify({ connpass_id: connpassID }),
    });
    // ユーザー情報を再取得または更新
    if (user) {
      const updatedUser = { ...user, connpass_id: connpassID };
      setUser(updatedUser);
    }
  };

  return (
    <AuthContext.Provider value={{ user, login, signUp, logout, linkConnpass, loading }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
