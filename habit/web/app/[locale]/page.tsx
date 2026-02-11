'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { motion } from 'framer-motion';
import { ChallengeCard } from '@/components/home/ChallengeCard';
import { JoinModal } from '@/components/home/JoinModal';
import { useChallengeStore } from '@/store/challengeStore';
import { api } from '@/lib/api';

export default function Home() {
  const t = useTranslations();
  const { challenges, setChallenges } = useChallengeStore();
  const [selectedChallenge, setSelectedChallenge] = useState<any>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadChallenges();
  }, []);

  const loadChallenges = async () => {
    setIsLoading(true);
    try {
      const response = await api.challenge.list();
      if (response.code === 0 && response.data) {
        setChallenges(response.data);
      }
    } catch (error) {
      console.error('Failed to load challenges:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleJoinClick = (challenge: any) => {
    setSelectedChallenge(challenge);
    setIsModalOpen(true);
  };

  const handleJoinSuccess = () => {
    loadChallenges();
  };

  const currentChallenge = challenges[0];

  return (
    <main className="min-h-screen">
      <div className="max-w-container mx-auto p-4 space-y-6 pb-20 md:pb-6">
        {/* 顶部标题 - 简洁 */}
        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.2 }}
          className="text-center pt-6 pb-2"
        >
          <h1 className="text-lg font-medium text-foreground">{t('home.welcome')}</h1>
        </motion.div>

        {/* 当前挑战卡片 */}
        {isLoading ? (
          <div className="card-large p-12 text-center">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-2 border-primary border-t-transparent"></div>
            <p className="mt-4 text-sm text-muted-foreground">{t('common.loading')}</p>
          </div>
        ) : currentChallenge ? (
          <ChallengeCard
            challenge={currentChallenge}
            onJoin={() => handleJoinClick(currentChallenge)}
          />
        ) : (
          <div className="card-large p-12 text-center text-muted-foreground">
            <p className="text-sm">No active challenges</p>
          </div>
        )}
      </div>

      {/* 参与弹窗 */}
      {selectedChallenge && (
        <JoinModal
          isOpen={isModalOpen}
          onClose={() => setIsModalOpen(false)}
          challenge={selectedChallenge}
          onSuccess={handleJoinSuccess}
        />
      )}
    </main>
  );
}
