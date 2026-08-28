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

    const currentPath = pathname || '';
    // パス判定
    const isGuestOnlyPath = ['/', '/login', '/signup'].includes(currentPath);
    // /events/[id] は公開。それ以外（/events, /events/new, /events/[id]/edit）は保護。
    const isPublicEventView = currentPath.startsWith('/events/') && currentPath.split('/').length === 3 && !currentPath.endsWith('/new');
    const isProtectedPath = currentPath.startsWith('/home') || (currentPath.startsWith('/events') && !isPublicEventView);

    if (user) {
      if (isGuestOnlyPath) {
        router.replace('/home');
      }
    } else {
      // 未ログインで保護されたページにいる場合はトップページへ
      if (isProtectedPath) {
        router.replace('/');
      }
    }
  }, [user, loading, pathname, router]);

  // ローディング中はローディング表示を出す
  if (loading) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="text-sm font-medium text-gray-400 animate-pulse tracking-widest italic font-serif">読み込み中...</div>
      </div>
    );
  }

  const currentPath = pathname || '';
  const isGuestOnlyPath = ['/', '/login', '/signup'].includes(currentPath);
  // 保護対象かつ未ログインの場合は何も表示せずuseEffectのリダイレクトを待つ
  const isPublicEventView = currentPath.startsWith('/events/') && currentPath.split('/').length === 3 && !currentPath.endsWith('/new');
  const isProtectedPath = currentPath.startsWith('/home') || (currentPath.startsWith('/events') && !isPublicEventView);

  if (user && isGuestOnlyPath) return null;
  if (!user && isProtectedPath) return null;

  return <>{children}</>;
}
