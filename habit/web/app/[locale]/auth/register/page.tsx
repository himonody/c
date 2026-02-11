'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { motion } from 'framer-motion';
import Link from 'next/link';
import { api } from '@/lib/api';
import { cn } from '@/lib/utils';

export default function RegisterPage() {
  const t = useTranslations();
  const router = useRouter();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const validateForm = () => {
    if (!username || !password || !confirmPassword) {
      setError('Please fill in all fields');
      return false;
    }

    if (username.length < 6 || username.length > 12) {
      setError('Username must be 6-12 characters');
      return false;
    }

    if (password.length < 6 || password.length > 12) {
      setError('Password must be 6-12 characters');
      return false;
    }

    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return false;
    }

    return true;
  };

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      const response = await api.auth.register(username, password);
      if (response.code === 0) {
        router.push('/auth/login');
      } else {
        setError(response.msg);
      }
    } catch (err) {
      setError('Registration failed. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-background to-accent/20 flex items-center justify-center p-4">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
        className="w-full max-w-md"
      >
        {/* Logo/Header */}
        <div className="text-center mb-8">
          <div className="text-6xl mb-4">🎯</div>
          <h1 className="text-3xl font-bold mb-2">Create Account</h1>
          <p className="text-muted-foreground">Join the challenge today</p>
        </div>

        {/* Register Form */}
        <div className="bg-card rounded-3xl p-6 shadow-lg">
          <form onSubmit={handleRegister} className="space-y-4">
            {/* Username */}
            <div>
              <label className="block text-sm font-medium mb-2">
                {t('auth.username')}
              </label>
              <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="6-12 characters"
                className={cn(
                  'w-full px-4 py-3 rounded-xl border bg-background',
                  'focus:outline-none focus:ring-2 focus:ring-primary',
                  error && 'border-destructive'
                )}
                minLength={6}
                maxLength={12}
              />
            </div>

            {/* Password */}
            <div>
              <label className="block text-sm font-medium mb-2">
                {t('auth.password')}
              </label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="6-12 characters"
                className={cn(
                  'w-full px-4 py-3 rounded-xl border bg-background',
                  'focus:outline-none focus:ring-2 focus:ring-primary',
                  error && 'border-destructive'
                )}
                minLength={6}
                maxLength={12}
              />
            </div>

            {/* Confirm Password */}
            <div>
              <label className="block text-sm font-medium mb-2">
                {t('auth.confirmPassword')}
              </label>
              <input
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder="Re-enter password"
                className={cn(
                  'w-full px-4 py-3 rounded-xl border bg-background',
                  'focus:outline-none focus:ring-2 focus:ring-primary',
                  error && 'border-destructive'
                )}
                minLength={6}
                maxLength={12}
              />
            </div>

            {error && (
              <p className="text-sm text-destructive bg-destructive/10 p-3 rounded-lg">
                {error}
              </p>
            )}

            {/* Register Button */}
            <motion.button
              whileTap={{ scale: 0.95 }}
              type="submit"
              disabled={isLoading}
              className="w-full py-4 bg-primary text-primary-foreground rounded-xl font-semibold hover:opacity-90 transition-opacity disabled:opacity-50"
            >
              {isLoading ? t('common.loading') : t('auth.register')}
            </motion.button>
          </form>

          {/* Login Link */}
          <div className="mt-6 text-center">
            <p className="text-sm text-muted-foreground">
              Already have an account?{' '}
              <Link href="/auth/login" className="text-primary font-medium hover:underline">
                {t('auth.login')}
              </Link>
            </p>
          </div>
        </div>

        {/* Back to Home */}
        <div className="mt-6 text-center">
          <Link href="/" className="text-sm text-muted-foreground hover:text-foreground">
            ← {t('common.back')}
          </Link>
        </div>
      </motion.div>
    </div>
  );
}
