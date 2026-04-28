'use client';

import { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';

export default function LoginForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isSignUp, setIsSignUp] = useState(false);
  const { login, signUp } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (isSignUp) {
        await signUp(email, password);
        alert('サインアップ成功！ログインしてください。');
        setIsSignUp(false);
      } else {
        await login(email, password);
      }
    } catch (err: any) {
      alert(err.message);
    }
  };

  return (
    <div className="max-w-md mx-auto mt-10 p-6 bg-white rounded shadow">
      <h2 className="text-2xl font-bold mb-4">{isSignUp ? 'サインアップ' : 'ログイン'}</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium">メールアドレス</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full p-2 border rounded"
            required
          />
        </div>
        <div>
          <label className="block text-sm font-medium">パスワード</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full p-2 border rounded"
            required
          />
        </div>
        <button type="submit" className="w-full py-2 bg-blue-600 text-white rounded">
          {isSignUp ? '登録' : 'ログイン'}
        </button>
      </form>
      <button
        onClick={() => setIsSignUp(!isSignUp)}
        className="mt-4 text-sm text-blue-600 hover:underline"
      >
        {isSignUp ? 'ログインはこちら' : 'アカウント作成はこちら'}
      </button>
    </div>
  );
}
