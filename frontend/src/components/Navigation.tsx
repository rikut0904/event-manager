'use client';

import { useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';

export default function Navigation() {
  const [isOpen, setIsOpen] = useState(false);
  const { user, logout } = useAuth();
  const pathname = usePathname();

  const navItems = user ? [
    { name: 'ダッシュボード', href: '/dashboard' },
  ] : [
    { name: 'ホーム', href: '/' },
    { name: 'ログイン', href: '/login' },
    { name: 'アカウント作成', href: '/signup' },
  ];

  const NavLink = ({ item, onClick }: { item: any, onClick?: () => void }) => {
    const isActive = pathname === item.href;
    return (
      <Link
        href={item.href}
        onClick={onClick}
        className={`block px-4 py-2.5 text-sm font-medium transition-colors rounded-lg ${
          isActive 
            ? 'bg-blue-50 text-blue-600' 
            : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
        }`}
      >
        {item.name}
      </Link>
    );
  };

  return (
    <>
      {/* Mobile Header */}
      <div className="lg:hidden flex items-center justify-between px-6 h-16 bg-white border-b border-gray-100 sticky top-0 z-50">
        <Link href="/" className="text-xl font-bold text-gray-900 tracking-tight">
          Event Manager
        </Link>
        <button
          onClick={() => setIsOpen(!isOpen)}
          className="p-2 text-gray-600 hover:text-gray-900 focus:outline-none"
        >
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            {isOpen ? (
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            ) : (
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16m-7 6h7" />
            )}
          </svg>
        </button>
      </div>

      {/* Mobile Overlay Menu */}
      {isOpen && (
        <div className="lg:hidden fixed inset-0 z-40 bg-white pt-20 px-6 space-y-2">
          {navItems.map((item) => (
            <NavLink key={item.href} item={item} onClick={() => setIsOpen(false)} />
          ))}
          {user && (
            <button
              onClick={() => { logout(); setIsOpen(false); }}
              className="w-full text-left px-4 py-2.5 text-sm font-medium text-gray-600 hover:text-gray-900"
            >
              ログアウト
            </button>
          )}
        </div>
      )}

      {/* Desktop Sidebar */}
      <aside className="hidden lg:flex flex-col w-64 h-screen bg-white border-r border-gray-100 fixed left-0 top-0 z-30">
        <div className="p-8 pb-12">
          <Link href="/" className="text-2xl font-bold text-gray-900 tracking-tight">
            Event Manager
          </Link>
        </div>
        
        <nav className="flex-1 px-4 space-y-1.5">
          <div className="text-xs font-bold text-gray-400 mb-4 px-4 uppercase tracking-[0.2em]">Menu</div>
          {navItems.map((item) => (
            <NavLink key={item.href} item={item} />
          ))}
        </nav>

        {user && (
          <div className="p-6 border-t border-gray-50">
            <div className="px-4 mb-4">
              <p className="text-xs text-gray-400 font-medium truncate">{user.email}</p>
            </div>
            <button
              onClick={logout}
              className="w-full text-left px-4 py-2.5 text-sm font-medium text-gray-600 hover:text-gray-900 hover:bg-gray-50 rounded-lg transition-colors"
            >
              ログアウト
            </button>
          </div>
        )}
      </aside>
    </>
  );
}
