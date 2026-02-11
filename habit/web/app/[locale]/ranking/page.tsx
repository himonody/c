'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { motion } from 'framer-motion';
import { api } from '@/lib/api';
import { cn } from '@/lib/utils';

type RankingType = 'invite' | 'wealth' | 'persistence';

interface RankingUser {
  userId: number;
  nickname: string;
  avatar?: string;
  value: string | number;
  rank: number;
}

export default function RankingPage() {
  const t = useTranslations();
  const [activeTab, setActiveTab] = useState<RankingType>('invite');
  const [rankingList, setRankingList] = useState<RankingUser[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    loadRanking();
  }, [activeTab]);

  const loadRanking = async () => {
    setIsLoading(true);
    try {
      let response;
      if (activeTab === 'invite') {
        response = await api.ranking.invite(1, 20);
      } else if (activeTab === 'wealth') {
        response = await api.ranking.wealth(1, 20);
      } else {
        response = await api.ranking.persistence(1, 20);
      }

      if (response.code === 0 && response.data) {
        setRankingList(response.data);
      }
    } catch (error) {
      console.error('Failed to load ranking:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const tabs: { key: RankingType; label: string; icon: string }[] = [
    { key: 'invite', label: t('ranking.invite'), icon: '👥' },
    { key: 'wealth', label: t('ranking.wealth'), icon: '💰' },
    { key: 'persistence', label: t('ranking.persistence'), icon: '🔥' },
  ];

  const getValueLabel = (type: RankingType) => {
    if (type === 'invite') return t('ranking.inviteCount');
    if (type === 'wealth') return t('ranking.totalIncome');
    return t('ranking.streakDays');
  };


  return (
    <div className="min-h-screen">
      <div className="max-w-container mx-auto p-4 space-y-6 pb-20 md:pb-6">
        {/* 标题 */}
        <div className="text-center pt-6 pb-2">
          <h1 className="text-lg font-medium text-foreground">{t('ranking.title')}</h1>
        </div>

        {/* Tabs */}
        <div className="card-base p-2 flex gap-2">
          {tabs.map((tab) => (
            <button
              key={tab.key}
              onClick={() => setActiveTab(tab.key)}
              className={cn(
                'flex-1 py-3 rounded-xl font-medium transition-all duration-200 relative',
                activeTab === tab.key
                  ? 'text-primary-foreground'
                  : 'text-muted-foreground'
              )}
            >
              {activeTab === tab.key && (
                <motion.div
                  layoutId="activeTab"
                  className="absolute inset-0 bg-primary rounded-xl"
                  transition={{ type: 'spring', bounce: 0.2, duration: 0.25 }}
                />
              )}
              <span className="relative z-10 flex items-center justify-center gap-2">
                <span>{tab.icon}</span>
                <span className="text-sm">{tab.label}</span>
              </span>
            </button>
          ))}
        </div>

        {/* Ranking List */}
        {isLoading ? (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
          </div>
        ) : (
          <div className="space-y-3">
            {rankingList.map((user, index) => {
              const rank = index + 1;

              return (
                <motion.div
                  key={user.userId}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.03, duration: 0.2 }}
                  className={cn(
                    'card-base p-4 flex items-center gap-3',
                    rank <= 3 && 'bg-gradient-to-r',
                    rank === 1 && 'from-warning/10 to-warning/5',
                    rank === 2 && 'from-muted/50 to-muted/20',
                    rank === 3 && 'from-warning/5 to-transparent'
                  )}
                >
                  {/* Rank */}
                  <div className={cn(
                    'w-10 h-10 rounded-full flex items-center justify-center font-bold text-sm',
                    rank === 1 && 'bg-warning text-white',
                    rank === 2 && 'bg-muted text-foreground',
                    rank === 3 && 'bg-warning/30 text-warning',
                    rank > 3 && 'bg-accent text-muted-foreground'
                  )}>
                    {rank}
                  </div>

                  {/* Avatar */}
                  <div className="w-10 h-10 rounded-full bg-card flex items-center justify-center text-lg">
                    👤
                  </div>

                  {/* Info */}
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium truncate">{user.nickname}</p>
                    <p className="text-xs text-muted-foreground">{getValueLabel(activeTab)}</p>
                  </div>

                  {/* Value */}
                  <div className="text-right">
                    <p className={cn(
                      'text-lg font-semibold',
                      rank <= 3 ? 'text-success' : 'text-primary'
                    )}>
                      {activeTab === 'wealth' && 'RM '}
                      {user.value}
                      {activeTab === 'persistence' && ' days'}
                    </p>
                  </div>
                </motion.div>
              );
            })}

            {rankingList.length === 0 && (
              <div className="text-center py-12 text-xs text-muted-foreground">
                No ranking data yet
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
