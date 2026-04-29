'use client';

import React, { createContext, useContext, useState, useEffect } from 'react';
import { apiRequest } from '@/lib/api';

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
    const savedUser = localStorage.getItem('auth_user');
    const token = localStorage.getItem('auth_token');
    
    if (savedUser && token) {
      setUser(JSON.parse(savedUser));
    }
    // ユーザーセットの直後にloadingを解除
    setLoading(false);
  }, []);

  const login = async (email: string, password: string) => {
    const data = await apiRequest('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    localStorage.setItem('auth_token', data.token);
    localStorage.setItem('auth_user', JSON.stringify(data.user));
    setUser(data.user);
  };

  const signUp = async (email: string, password: string) => {
    await apiRequest('/auth/signup', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    // アカウント作成後、そのままログインを実行
    await login(email, password);
  };

  const logout = () => {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('auth_user');
    setUser(null);
    window.location.href = '/'; // 確実に状態をリセットしてトップへ
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
      localStorage.setItem('auth_user', JSON.stringify(updatedUser));
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
