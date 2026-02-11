'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { motion } from 'framer-motion';
import Link from 'next/link';
import { useUserStore } from '@/store/userStore';
import { useInviteStore } from '@/store/inviteStore';
import { api } from '@/lib/api';

export default function ProfilePage() {
  const t = useTranslations();
  const { userInfo, walletInfo } = useUserStore();
  const { stats } = useInviteStore();
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadProfileData();
  }, []);

  const loadProfileData = async () => {
    setIsLoading(true);
    try {
      await Promise.all([
        api.auth.getUserInfo(),
        api.wallet.getInfo(),
        api.invite.stats(),
      ]);
    } catch (error) {
      console.error('Failed to load profile data:', error);
    } finally {
      setIsLoading(false);
    }
  };


  const dataItems = [
    { key: 'todayInvite', label: t('profile.todayInvite'), value: stats?.todayInviteCount || 0, icon: '👥' },
    { key: 'totalInvite', label: t('profile.totalInvite'), value: stats?.totalInviteCount || 0, icon: '🎯' },
    { key: 'inviteIncome', label: t('profile.inviteIncome'), value: `RM ${stats?.inviteIncome || '0.00'}`, icon: '💰' },
    { key: 'streakDays', label: t('profile.streakDays'), value: '0 days', icon: '🔥' },
  ];

  const actionItems = [
    { key: 'accountDetail', label: t('profile.accountDetail'), icon: '📊', href: '/profile/transactions' },
    { key: 'myInvite', label: t('profile.myInvite'), icon: '👥', href: '/invite' },
    { key: 'supplementCard', label: t('profile.supplementCard'), icon: '🎫', href: '/profile/supplement' },
    { key: 'checkInQuery', label: t('profile.checkInQuery'), icon: '🔍', href: '/checkin/history' },
  ];

  return (
    <div className="min-h-screen">
      <div className="max-w-container mx-auto p-4 space-y-6 pb-20 md:pb-6">
        {/* User Card */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.2 }}
          className="card-large p-6"
        >
          <div className="flex items-center gap-4">
            <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center text-3xl">
              👤
            </div>
            <div className="flex-1 min-w-0">
              <h2 className="text-lg font-medium truncate">{userInfo?.nickname || 'Guest'}</h2>
              <p className="text-xs text-muted-foreground">@{userInfo?.username || 'username'}</p>
            </div>
            <Link
              href="/profile/edit"
              className="px-4 py-2 bg-primary text-primary-foreground rounded-xl text-xs font-medium"
            >
              {t('profile.updateProfile')}
            </Link>
          </div>
        </motion.div>

        {/* Assets Section */}
        <div className="card-base p-4">
          <h3 className="text-sm font-medium mb-4">{t('profile.assets')}</h3>
          <div className="grid grid-cols-2 gap-3">
            <div className="text-center p-3 bg-accent/30 rounded-xl">
              <p className="text-xs text-muted-foreground mb-1">{t('profile.balance')}</p>
              <p className="text-2xl font-semibold text-success">RM {walletInfo?.balance || '0.00'}</p>
            </div>
            <div className="text-center p-3 bg-accent/30 rounded-xl">
              <p className="text-xs text-muted-foreground mb-1">{t('profile.totalIncome')}</p>
              <p className="text-2xl font-semibold text-primary">RM {walletInfo?.totalIncome || '0.00'}</p>
            </div>
          </div>
          <Link href="/wallet">
            <motion.button
              whileTap={{ scale: 0.97 }}
              className="btn-primary w-full mt-4"
            >
              {t('wallet.withdraw')} →
            </motion.button>
          </Link>
        </div>

        {/* Data Section */}
        <div className="card-base p-4">
          <h3 className="text-sm font-medium mb-4">{t('profile.data')}</h3>
          <div className="grid grid-cols-2 gap-3">
            {dataItems.map((item) => (
              <div key={item.key} className="flex items-center gap-2 p-3 bg-accent/30 rounded-xl">
                <div className="text-xl">{item.icon}</div>
                <div className="flex-1 min-w-0">
                  <p className="text-[10px] text-muted-foreground truncate">{item.label}</p>
                  <p className="text-sm font-semibold truncate">{item.value}</p>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Actions Section */}
        <div className="card-base p-4">
          <h3 className="text-sm font-medium mb-4">{t('profile.actions')}</h3>
          <div className="space-y-2">
            {actionItems.map((item) => (
              <Link key={item.key} href={item.href}>
                <motion.div
                  whileTap={{ scale: 0.98 }}
                  transition={{ duration: 0.15 }}
                  className="flex items-center gap-3 p-3 bg-accent/30 rounded-xl"
                >
                  <div className="text-xl">{item.icon}</div>
                  <div className="flex-1 text-sm font-medium">{item.label}</div>
                  <div className="text-muted-foreground text-xs">→</div>
                </motion.div>
              </Link>
            ))}
          </div>
        </div>

        {/* Settings */}
        <div className="card-base p-4">
          <div className="space-y-2">
            <Link href="/profile/change-password">
              <motion.div
                whileTap={{ scale: 0.98 }}
                transition={{ duration: 0.15 }}
                className="flex items-center gap-3 p-3 bg-accent/30 rounded-xl"
              >
                <div className="text-xl">🔒</div>
                <div className="flex-1 text-sm font-medium">{t('profile.changePassword')}</div>
                <div className="text-muted-foreground text-xs">→</div>
              </motion.div>
            </Link>
            <Link href="/profile/set-pay-password">
              <motion.div
                whileTap={{ scale: 0.98 }}
                transition={{ duration: 0.15 }}
                className="flex items-center gap-3 p-3 bg-accent/30 rounded-xl"
              >
                <div className="text-xl">🔐</div>
                <div className="flex-1 text-sm font-medium">{t('profile.setPayPassword')}</div>
                <div className="text-muted-foreground text-xs">→</div>
              </motion.div>
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
