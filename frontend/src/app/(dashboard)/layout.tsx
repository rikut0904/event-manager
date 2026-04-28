'use client';

import Navigation from "@/components/Navigation";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex flex-col lg:flex-row min-h-screen">
      <Navigation />
      <div className="flex-1 lg:ml-64 bg-gray-50/30">
        {children}
      </div>
    </div>
  );
}
