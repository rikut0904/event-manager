'use client';

import Link from 'next/link';

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-white">
      <header className="max-w-6xl mx-auto px-6 h-16 flex items-center justify-between">
        <span className="font-bold text-gray-900 tracking-tight">Event Manager</span>
        <div className="flex gap-4 items-center">
          <Link href="/login" className="text-sm font-medium text-gray-500 hover:text-gray-900 transition-colors">ログイン</Link>
          <Link href="/signup" className="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors">無料で始める</Link>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-6 py-24 lg:py-32 text-center lg:text-left">
        <h1 className="text-5xl lg:text-6xl font-bold text-gray-900 mb-8 tracking-tight leading-tight">
          Manage your events <br />
        </h1>
        <p className="text-xl text-gray-500 mb-12 leading-relaxed max-w-2xl">
          イベント運営をするためのプラットフォームです。あなたはイベントの告知・集客から、参加者の管理、当日の運営まで、イベントに関するあらゆる業務をスムーズに行うことができます。
        </p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center lg:justify-start">
          <Link
            href="/signup"
            className="px-8 py-4 bg-blue-600 text-white font-bold rounded-xl hover:bg-blue-700 transition-all shadow-xl shadow-blue-100"
          >
            無料でアカウント作成
          </Link>
          <Link
            href="/login"
            className="px-8 py-4 bg-white text-gray-600 border border-gray-100 font-bold rounded-xl hover:bg-gray-50 transition-all"
          >
            ログイン
          </Link>
        </div>
      </main>
    </div>
  );
}
