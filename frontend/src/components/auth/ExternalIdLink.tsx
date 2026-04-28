'use client';

import { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';

export default function ConnpassLink() {
  const [connpassID, setConnpassID] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { user, linkConnpass } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    try {
      await linkConnpass(connpassID);
      alert('connpass IDを連携しました！');
    } catch (err: any) {
      alert(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  if (!user) return null;

  return (
    <div className="">
      {user.connpass_id ? (
        <div className="flex items-center gap-3 p-4 bg-green-50 rounded-xl border border-green-100">
          <span className="text-2xl">✅</span>
          <div>
            <p className="text-xs text-green-600 font-bold uppercase tracking-wider">連携済み</p>
            <p className="text-lg font-bold text-green-900">{user.connpass_id}</p>
          </div>
        </div>
      ) : (
        <form onSubmit={handleSubmit} className="space-y-4">
          <p className="text-sm text-gray-500 leading-relaxed">
            connpassの「表示名」または「ID」を入力して連携してください。
          </p>
          <div className="relative">
            <input
              type="text"
              value={connpassID}
              onChange={(e) => setConnpassID(e.target.value)}
              className="w-full p-3 border border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 outline-none"
              placeholder="例: yamada_tarou"
              required
            />
          </div>
          <button
            type="submit"
            disabled={isLoading}
            className="w-full py-3 bg-red-500 text-white rounded-xl font-bold hover:bg-red-600 transition-shadow hover:shadow-lg shadow-red-100 disabled:bg-red-300"
          >
            {isLoading ? '連携中...' : 'connpassと連携する'}
          </button>
        </form>
      )}
    </div>
  );
}
