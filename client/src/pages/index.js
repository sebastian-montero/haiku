import { useState } from 'react';
import { useRouter } from 'next/router';

export default function Home() {

  const router = useRouter();
 
  return (
    <div className="min-h-screen bg-white flex items-center justify-center px-4">
      <div className="max-w-md w-full space-y-6">
        <div className="text-sm text-gray-700 mb-16 text-center mb-10">
        <h1 className="text-3xl font-bold text-black tracking-wide text-center">haiku⿻</h1>
        <p className="text-sm text-gray-700 mb-16 text-center mt-2">
        a platform to write freely and share your real-time writing process, let others tune in to experience the raw art of creation and inspiration.
          </p>
          </div>

          <div className="flex justify-center space-x-4 mt-6">
  <button
    onClick={() => router.push('/login')}
    className="w-32 py-2 px-4 text-sm text-black font-medium border-black hover:bg-black hover:text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-black"
  >
    log in
  </button>
  <button
    onClick={() => router.push('/signup')}
    className="w-32 py-2 px-4 text-sm text-black font-medium border-black hover:bg-black hover:text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-black"
  >
    sign up
  </button>
</div>
          

      </div>
    </div>
  );
}