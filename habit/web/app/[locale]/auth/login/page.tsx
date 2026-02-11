'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { motion } from 'framer-motion';
import Link from 'next/link';
import { api } from '@/lib/api';
import { cn } from '@/lib/utils';

export default function LoginPage() {
  const t = useTranslations();
  const router = useRouter();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!username || !password) {
      setError('Please fill in all fields');
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      const response = await api.auth.login(username, password);
      if (response.code === 0) {
        router.push('/');
      } else {
        setError(response.msg);
      }
    } catch (err) {
      setError('Login failed. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.2 }}
        className="w-full max-w-md"
      >
        {/* Logo/Header */}
        <div className="text-center mb-8">
          <div className="text-5xl mb-4">🎯</div>
          <h1 className="text-2xl font-medium mb-1">{t('home.welcome')}</h1>
          <p className="text-xs text-muted-foreground">Sign in to continue</p>
        </div>

        {/* Login Form */}
        <div className="card-large p-6">
          <form onSubmit={handleLogin} className="space-y-4">
            {/* Username */}
            <div>
              <label className="block text-xs text-muted-foreground mb-2">
                {t('auth.username')}
              </label>
              <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="Enter username"
                className={cn(
                  'w-full px-4 py-3 rounded-xl border bg-background text-sm',
                  'focus:outline-none focus:ring-2 focus:ring-primary',
                  error && 'border-destructive'
                )}
                minLength={6}
                maxLength={12}
              />
            </div>

            {/* Password */}
            <div>
              <label className="block text-xs text-muted-foreground mb-2">
                {t('auth.password')}
              </label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Enter password"
                className={cn(
                  'w-full px-4 py-3 rounded-xl border bg-background text-sm',
                  'focus:outline-none focus:ring-2 focus:ring-primary',
                  error && 'border-destructive'
                )}
                minLength={6}
                maxLength={12}
              />
            </div>

            {error && (
              <p className="text-xs text-destructive bg-destructive/10 p-3 rounded-xl">
                {error}
              </p>
            )}

            {/* Login Button */}
            <motion.button
              whileTap={{ scale: 0.97 }}
              type="submit"
              disabled={isLoading}
              className="btn-primary w-full disabled:opacity-50"
            >
              {isLoading ? t('common.loading') : t('auth.login')}
            </motion.button>
          </form>

          {/* Register Link */}
          <div className="mt-6 text-center">
            <p className="text-xs text-muted-foreground">
              Don't have an account?{' '}
              <Link href="/auth/register" className="text-primary font-medium hover:underline">
                {t('auth.register')}
              </Link>
            </p>
          </div>
        </div>

        {/* Back to Home */}
        <div className="mt-6 text-center">
          <Link href="/" className="text-xs text-muted-foreground hover:text-foreground">
            ← {t('common.back')}
          </Link>
        </div>
      </motion.div>
    </div>
  );
}
