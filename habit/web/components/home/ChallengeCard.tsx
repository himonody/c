'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { motion } from 'framer-motion';

interface ChallengeCardProps {
  challenge: {
    challengeId: number;
    name: string;
    prizePool: string;
    participantCount: number;
    duration: number;
    entryFee: string;
    checkInStartTime: string;
    checkInEndTime: string;
  };
  onJoin: () => void;
}

export function ChallengeCard({ challenge, onJoin }: ChallengeCardProps) {
  const t = useTranslations('challenge');
  const [isExpanded, setIsExpanded] = useState(false);

  const participantAvatars = Array(Math.min(8, challenge.participantCount)).fill('👤');

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.25 }}
      className="card-large p-6"
    >
      {/* 挑战名称 */}
      <div className="text-center mb-6">
        <h2 className="text-lg font-medium text-foreground mb-1">{challenge.name}</h2>
        <p className="text-xs text-muted-foreground">
          {challenge.duration} {t('days')} • {challenge.checkInStartTime} - {challenge.checkInEndTime}
        </p>
      </div>

      {/* 奖池金额 - 视觉焦点 */}
      <div className="bg-gradient-to-br from-warning/10 to-warning/5 rounded-2xl p-6 mb-4 text-center">
        <p className="text-xs text-muted-foreground mb-2">{t('prizePool')}</p>
        <div className="text-5xl font-semibold text-warning mb-1">
          RM {challenge.prizePool}
        </div>
        <p className="text-xs text-muted-foreground mt-2">{t('entryFee')}: RM {challenge.entryFee}</p>
      </div>

      {/* 当前参与人数 */}
      <div className="bg-accent/50 rounded-xl p-4 mb-4 text-center">
        <p className="text-xs text-muted-foreground mb-1">{t('participants')}</p>
        <p className="text-3xl font-semibold text-primary">{challenge.participantCount}</p>
      </div>

      {/* 用户头像横向列表 */}
      <div className="mb-5">
        <div className="flex justify-center -space-x-3">
          {participantAvatars.map((_, i) => (
            <div
              key={i}
              className="w-10 h-10 rounded-full bg-card border-2 border-background flex items-center justify-center text-sm shadow-sm"
            >
              {_}
            </div>
          ))}
          {challenge.participantCount > 8 && (
            <div className="w-10 h-10 rounded-full bg-primary/10 border-2 border-background flex items-center justify-center text-xs font-semibold text-primary shadow-sm">
              +{challenge.participantCount - 8}
            </div>
          )}
        </div>
      </div>

      {/* 规则入口 */}
      <button
        onClick={() => setIsExpanded(!isExpanded)}
        className="w-full text-sm text-primary font-medium mb-4 py-2 hover:underline transition-all"
      >
        {isExpanded ? '▼' : '▶'} {t('rules')}
      </button>

      <motion.div
        initial={false}
        animate={{ height: isExpanded ? 'auto' : 0 }}
        className="overflow-hidden"
      >
        <div className="bg-accent/50 rounded-xl p-4 mb-5 text-sm text-muted-foreground space-y-2">
          <p>• Check in daily between {challenge.checkInStartTime} - {challenge.checkInEndTime}</p>
          <p>• Complete all {challenge.duration} days to get rewards</p>
          <p>• Failed check-ins result in entry fee deduction</p>
          <p>• Prize pool distributed among successful participants</p>
        </div>
      </motion.div>

      {/* 主按钮 - 56px 高度 */}
      <motion.button
        whileTap={{ scale: 0.97 }}
        onClick={onJoin}
        className="btn-primary w-full"
      >
        {t('join')}
      </motion.button>
    </motion.div>
  );
}
