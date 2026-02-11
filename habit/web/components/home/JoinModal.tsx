'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { motion, AnimatePresence } from 'framer-motion';
import { api } from '@/lib/api';
import { cn } from '@/lib/utils';

interface JoinModalProps {
  isOpen: boolean;
  onClose: () => void;
  challenge: {
    challengeId: number;
    name: string;
    entryFee: string;
  };
  onSuccess: () => void;
}

export function JoinModal({ isOpen, onClose, challenge, onSuccess }: JoinModalProps) {
  const t = useTranslations();
  const [payPassword, setPayPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const handleJoin = async () => {
    if (!payPassword) {
      setError('Please enter payment password');
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      const response = await api.challenge.join(challenge.challengeId, payPassword);
      if (response.code === 0) {
        onSuccess();
        onClose();
      } else {
        setError(response.msg);
      }
    } catch (err) {
      setError('Failed to join challenge');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
            className="fixed inset-0 bg-black/50 z-50"
          />

          {/* Modal */}
          <div className="fixed inset-0 flex items-center justify-center z-50 p-4">
            <motion.div
              initial={{ opacity: 0, scale: 0.9, y: 20 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.9, y: 20 }}
              transition={{ duration: 0.2 }}
              className="bg-card rounded-2xl p-6 w-full max-w-sm"
            >
              <h3 className="text-xl font-bold mb-4">{t('challenge.join')}</h3>

              <div className="space-y-4 mb-6">
                <div className="bg-background rounded-xl p-4">
                  <p className="text-sm text-muted-foreground mb-1">
                    {t('challenge.title')}
                  </p>
                  <p className="font-semibold">{challenge.name}</p>
                </div>

                <div className="bg-background rounded-xl p-4">
                  <p className="text-sm text-muted-foreground mb-1">
                    {t('challenge.entryFee')}
                  </p>
                  <p className="text-2xl font-bold text-primary">
                    RM {challenge.entryFee}
                  </p>
                </div>

                <div>
                  <label className="block text-sm font-medium mb-2">
                    {t('auth.payPassword')}
                  </label>
                  <input
                    type="password"
                    value={payPassword}
                    onChange={(e) => setPayPassword(e.target.value)}
                    placeholder="Enter payment password"
                    className={cn(
                      'w-full px-4 py-3 rounded-xl border bg-background',
                      'focus:outline-none focus:ring-2 focus:ring-primary',
                      error && 'border-destructive'
                    )}
                    maxLength={6}
                  />
                  {error && (
                    <p className="text-sm text-destructive mt-1">{error}</p>
                  )}
                </div>
              </div>

              <div className="flex gap-3">
                <button
                  onClick={onClose}
                  disabled={isLoading}
                  className="flex-1 py-3 rounded-xl border border-border hover:bg-accent transition-colors"
                >
                  {t('common.cancel')}
                </button>
                <button
                  onClick={handleJoin}
                  disabled={isLoading}
                  className="flex-1 py-3 rounded-xl bg-primary text-primary-foreground font-semibold hover:opacity-90 transition-opacity disabled:opacity-50"
                >
                  {isLoading ? t('common.loading') : t('common.confirm')}
                </button>
              </div>
            </motion.div>
          </div>
        </>
      )}
    </AnimatePresence>
  );
}
