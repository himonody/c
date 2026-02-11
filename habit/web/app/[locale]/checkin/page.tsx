'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { motion } from 'framer-motion';
import Link from 'next/link';
import { useCheckinStore } from '@/store/checkinStore';
import { api } from '@/lib/api';
import { cn } from '@/lib/utils';

export default function CheckInPage() {
  const t = useTranslations();
  const { todayStats, todayList, hasCheckedInToday, setTodayStats, setTodayList, setHasCheckedInToday } = useCheckinStore();
  const [isChecking, setIsChecking] = useState(false);
  const [challengeId] = useState(1);

  useEffect(() => {
    loadTodayData();
  }, []);

  const loadTodayData = async () => {
    try {
      const [statsRes, listRes] = await Promise.all([
        api.checkin.todayStats(challengeId),
        api.checkin.todayList(challengeId, 1, 10),
      ]);

      if (statsRes.code === 0 && statsRes.data) {
        setTodayStats(statsRes.data);
      }

      if (listRes.code === 0 && listRes.data) {
        setTodayList(listRes.data);
      }
    } catch (error) {
      console.error('Failed to load today data:', error);
    }
  };

  const handleCheckIn = async () => {
    if (hasCheckedInToday || isChecking) return;

    setIsChecking(true);
    try {
      const response = await api.checkin.checkIn(challengeId);
      if (response.code === 0) {
        setHasCheckedInToday(true);
        loadTodayData();
      }
    } catch (error) {
      console.error('Check-in failed:', error);
    } finally {
      setIsChecking(false);
    }
  };

  return (
    <div className="min-h-screen">
      <div className="max-w-container mx-auto p-4 space-y-6 pb-20 md:pb-6">
        {/* 标题 */}
        <div className="text-center pt-6 pb-2">
          <h1 className="text-lg font-medium text-foreground">{t('checkin.title')}</h1>
        </div>

        {/* 打卡按钮 */}
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.2 }}
          className="card-large p-8 text-center"
        >
          {hasCheckedInToday ? (
            <div className="space-y-4">
              <motion.div
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ type: 'spring', duration: 0.3 }}
                className="w-20 h-20 mx-auto bg-success/10 rounded-full flex items-center justify-center"
              >
                <div className="text-4xl text-success">✓</div>
              </motion.div>
              <p className="text-lg font-medium">{t('checkin.todayChecked')}</p>
              <p className="text-xs text-muted-foreground">See you tomorrow!</p>
            </div>
          ) : (
            <div className="space-y-4">
              <div className="w-20 h-20 mx-auto bg-primary/10 rounded-full flex items-center justify-center mb-4">
                <div className="text-4xl">📝</div>
              </div>
              <motion.button
                whileTap={{ scale: 0.97 }}
                onClick={handleCheckIn}
                disabled={isChecking}
                className="btn-primary w-full max-w-xs mx-auto"
              >
                {isChecking ? t('common.loading') : t('checkin.checkInButton')}
              </motion.button>
            </div>
          )}
        </motion.div>

        {/* 今日统计 */}
        {todayStats && (
          <div className="grid grid-cols-2 gap-4">
            <div className="card-base p-4 text-center">
              <p className="text-xs text-muted-foreground mb-2">{t('checkin.checkedIn')}</p>
              <p className="text-3xl font-semibold text-success">{todayStats.checkedInCount}</p>
            </div>
            <div className="card-base p-4 text-center">
              <p className="text-xs text-muted-foreground mb-2">{t('checkin.notCheckedIn')}</p>
              <p className="text-3xl font-semibold text-warning">{todayStats.notCheckedInCount}</p>
            </div>
          </div>
        )}

        {/* 今日之星 */}
        <div className="card-base p-4">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-sm font-medium">{t('checkin.todayRanking')}</h3>
            <Link href="/checkin/history" className="text-xs text-primary hover:underline">
              {t('checkin.history')} →
            </Link>
          </div>

          <div className="space-y-2">
            {todayList.length > 0 ? (
              todayList.slice(0, 5).map((record, index) => (
                <motion.div
                  key={record.userId}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.03, duration: 0.2 }}
                  className="flex items-center gap-3 p-2 bg-accent/30 rounded-xl"
                >
                  <div className="w-7 h-7 rounded-full bg-primary/10 flex items-center justify-center text-xs font-semibold text-primary">
                    {index + 1}
                  </div>
                  <div className="w-8 h-8 rounded-full bg-card flex items-center justify-center text-sm">
                    👤
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium truncate">{record.nickname}</p>
                    <p className="text-[10px] text-muted-foreground">{record.checkInTime}</p>
                  </div>
                  <div className={cn(
                    'px-2 py-0.5 rounded-full text-[10px] font-medium',
                    record.status === 'success' && 'bg-success/10 text-success'
                  )}>
                    {t(`checkin.status.${record.status}`)}
                  </div>
                </motion.div>
              ))
            ) : (
              <p className="text-center text-xs text-muted-foreground py-6">No check-ins yet today</p>
            )}
          </div>
        </div>

        {/* 快捷入口 */}
        <div className="grid grid-cols-2 gap-3">
          <Link href="/checkin/calendar">
            <motion.div
              whileTap={{ scale: 0.97 }}
              className="card-base p-4 text-center hover:shadow-hover transition-all"
            >
              <div className="text-2xl mb-1">📅</div>
              <p className="text-xs font-medium">{t('checkin.calendar')}</p>
            </motion.div>
          </Link>
          <Link href="/checkin/history">
            <motion.div
              whileTap={{ scale: 0.97 }}
              className="card-base p-4 text-center hover:shadow-hover transition-all"
            >
              <div className="text-2xl mb-1">📊</div>
              <p className="text-xs font-medium">{t('checkin.history')}</p>
            </motion.div>
          </Link>
        </div>
      </div>
    </div>
  );
}
