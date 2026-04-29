'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useRouter, usePathname } from 'next/navigation';
import { useEffect } from 'react';

export default function AuthGuard({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth();
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    // ローディング中は何もしない
    if (loading) return;

    // パス判定
    const isGuestOnlyPath = ['/login', '/signup'].includes(pathname);
    const isProtectedPath = pathname.startsWith('/dashboard');

    if (user) {
      // ログイン済みで /login や /signup にアクセスした場合は即座にダッシュボードへ
      if (isGuestOnlyPath) {
        router.replace('/dashboard');
      }
    } else {
      // 未ログインで保護されたページにいる場合はトップページへ
      if (isProtectedPath) {
        router.replace('/');
      }
    }
  }, [user, loading, pathname, router]);

  // ローディング中は真っ白な画面またはローディング表示を出して、一瞬のチラつきを防ぐ
  if (loading) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="text-sm font-medium text-gray-400 animate-pulse tracking-widest uppercase italic font-serif">Loading...</div>
      </div>
    );
  }

  // リダイレクトが必要な状態（ログイン済みでログイン画面にいる等）では中身を表示しない
  const isGuestOnlyPath = ['/', '/login', '/signup'].includes(pathname);
  const isProtectedPath = pathname.startsWith('/dashboard');
  if (user && isGuestOnlyPath) return null;
  if (!user && isProtectedPath) return null;

  return <>{children}</>;
}
