'use client';

import { usePathname } from 'next/navigation';
import Navigation from "@/components/Navigation";
import AuthGuard from "@/components/auth/AuthGuard";
import { AuthProvider } from "@/contexts/AuthContext";
import "./globals.css";
import { Geist, Geist_Mono } from "next/font/google";

const geistSans = Geist({ variable: "--font-geist-sans", subsets: ["latin"] });
const geistMono = Geist_Mono({ variable: "--font-geist-mono", subsets: ["latin"] });

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const currentPath = pathname || '';
  // サイドバーを表示しないページの判定
  const noSidebarPaths = ['/', '/login', '/signup'];
  const shouldShowSidebar = !noSidebarPaths.includes(currentPath);

  return (
    <html lang="ja">
      <body className={`${geistSans.variable} ${geistMono.variable} antialiased bg-white`}>
        <AuthProvider>
          <AuthGuard>
            <div className={`flex flex-col ${shouldShowSidebar ? 'lg:flex-row' : ''} min-h-screen`}>
              {shouldShowSidebar && <Navigation />}
              <main className={`flex-1 ${shouldShowSidebar ? 'lg:ml-64' : ''}`}>
                {children}
              </main>
            </div>
          </AuthGuard>
        </AuthProvider>
      </body>
    </html>
  );
}
