'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { motion, AnimatePresence } from 'framer-motion';
import { useCheckinStore } from '@/store/checkinStore';
import { api } from '@/lib/api';
import { cn } from '@/lib/utils';

export function CheckInCard() {
  const t = useTranslations('checkin');
  const [isChecking, setIsChecking] = useState(false);

  const {
    streak,
    totalDays,
    hasCheckedInToday,
    setHasCheckedInToday,
    setStreak,
    setTotalDays,
  } = useCheckinStore();

  const handleCheckIn = async () => {
    if (hasCheckedInToday || isChecking) return;

    setIsChecking(true);
    try {
      const response = await api.checkin.checkIn();

      if (response.code === 0) {
        setHasCheckedInToday(true);
        // 更新连续天数和总天数（从后端返回的数据）
        if (response.data) {
          setStreak(response.data.streak || streak + 1);
          setTotalDays(response.data.totalDays || totalDays + 1);
        }
      }
    } catch (error) {
      console.error('Check-in failed:', error);
    } finally {
      setIsChecking(false);
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
      className="bg-card rounded-2xl p-6 shadow-lg"
    >
      {/* 连续天数显示 */}
      <div className="text-center mb-6">
        <motion.div
          key={streak}
          initial={{ scale: 0.8, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ duration: 0.2 }}
          className="inline-block"
        >
          <div className="text-5xl font-bold text-primary mb-2">
            {streak}
          </div>
          <div className="text-sm text-muted-foreground">
            {t('streak')}
          </div>
        </motion.div>
      </div>

      {/* 签到按钮 */}
      <AnimatePresence mode="wait">
        {hasCheckedInToday ? (
          <motion.div
            key="checked"
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.9 }}
            transition={{ duration: 0.2 }}
            className="w-full py-4 px-6 bg-secondary text-secondary-foreground rounded-xl text-center font-medium"
          >
            ✓ {t('todayChecked')}
          </motion.div>
        ) : (
          <motion.button
            key="unchecked"
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.9 }}
            transition={{ duration: 0.2 }}
            whileTap={{ scale: 0.95 }}
            onClick={handleCheckIn}
            disabled={isChecking}
            className={cn(
              'w-full py-4 px-6 bg-primary text-primary-foreground rounded-xl font-medium',
              'transition-opacity duration-200',
              'hover:opacity-90',
              'disabled:opacity-50 disabled:cursor-not-allowed'
            )}
          >
            {isChecking ? t('common.loading') : t('checkInButton')}
          </motion.button>
        )}
      </AnimatePresence>

      {/* 总天数统计 */}
      <div className="mt-6 pt-6 border-t border-border">
        <div className="flex justify-between items-center text-sm">
          <span className="text-muted-foreground">{t('totalDays')}</span>
          <span className="font-semibold">{totalDays}</span>
        </div>
      </div>
    </motion.div>
  );
}
