'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { motion } from 'framer-motion';
import { useInviteStore } from '@/store/inviteStore';
import { api } from '@/lib/api';

export default function InvitePage() {
  const t = useTranslations();
  const { stats, inviteList, inviteCode, setStats, setInviteList, setInviteCode } = useInviteStore();
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadInviteData();
  }, []);

  const loadInviteData = async () => {
    setIsLoading(true);
    try {
      const [statsRes, listRes, codeRes] = await Promise.all([
        api.invite.stats(),
        api.invite.myInvites(1, 20),
        api.invite.getCode(),
      ]);

      if (statsRes.code === 0 && statsRes.data) {
        setStats(statsRes.data);
      }

      if (listRes.code === 0 && listRes.data) {
        setInviteList(listRes.data);
      }

      if (codeRes.code === 0 && codeRes.data) {
        setInviteCode(codeRes.data.code);
      }
    } catch (error) {
      console.error('Failed to load invite data:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCopyCode = () => {
    if (inviteCode) {
      navigator.clipboard.writeText(inviteCode);
      alert(t('invite.inviteSuccess'));
    }
  };

  const handleShare = () => {
    const shareUrl = `${window.location.origin}?invite=${inviteCode}`;
    if (navigator.share) {
      navigator.share({
        title: 'Join Challenge',
        text: 'Join me in this exciting challenge!',
        url: shareUrl,
      });
    } else {
      navigator.clipboard.writeText(shareUrl);
      alert('Link copied to clipboard!');
    }
  };

  return (
    <div className="min-h-screen">
      <div className="max-w-container mx-auto p-4 space-y-6 pb-20 md:pb-6">
        {/* 标题 */}
        <div className="text-center pt-6 pb-2">
          <h1 className="text-lg font-medium text-foreground">{t('invite.title')}</h1>
        </div>

        {/* Invite Stats */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.2 }}
          className="card-large p-6"
        >
          <h3 className="text-center text-sm font-medium mb-4">{t('invite.inviteStats')}</h3>
          <div className="grid grid-cols-3 gap-3">
            <div className="text-center p-3 bg-accent/30 rounded-xl">
              <p className="text-3xl font-semibold text-primary mb-1">
                {stats?.todayInviteCount || 0}
              </p>
              <p className="text-xs text-muted-foreground">{t('profile.todayInvite')}</p>
            </div>
            <div className="text-center p-3 bg-accent/30 rounded-xl">
              <p className="text-3xl font-semibold text-primary mb-1">
                {stats?.totalInviteCount || 0}
              </p>
              <p className="text-xs text-muted-foreground">{t('profile.totalInvite')}</p>
            </div>
            <div className="text-center p-3 bg-accent/30 rounded-xl">
              <p className="text-3xl font-semibold text-success mb-1">
                {stats?.inviteIncome || '0.00'}
              </p>
              <p className="text-xs text-muted-foreground">RM</p>
            </div>
          </div>
        </motion.div>

        {/* Invite Code Card */}
        <div className="card-base p-6">
          <h3 className="text-sm font-medium mb-4 text-center">{t('invite.inviteCode')}</h3>

          {/* QR Code Placeholder */}
          <div className="bg-primary/5 rounded-2xl p-6 mb-4">
            <div className="aspect-square max-w-[180px] mx-auto bg-white rounded-xl flex items-center justify-center border border-primary/10">
              <div className="text-center">
                <div className="text-5xl mb-2">📱</div>
                <p className="text-xs text-muted-foreground">QR Code</p>
              </div>
            </div>
          </div>

          {/* Code Display */}
          <div className="bg-accent/30 rounded-xl p-4 mb-4 text-center">
            <p className="text-xs text-muted-foreground mb-1">Your Code</p>
            <p className="text-2xl font-semibold tracking-wider">{inviteCode || 'LOADING...'}</p>
          </div>

          {/* Action Buttons */}
          <div className="grid grid-cols-2 gap-3">
            <motion.button
              whileTap={{ scale: 0.97 }}
              transition={{ duration: 0.15 }}
              onClick={handleCopyCode}
              className="py-3 bg-primary/10 text-primary rounded-xl text-sm font-medium transition-all duration-200"
            >
              {t('invite.copy')}
            </motion.button>
            <motion.button
              whileTap={{ scale: 0.97 }}
              transition={{ duration: 0.15 }}
              onClick={handleShare}
              className="py-3 bg-primary text-primary-foreground rounded-xl text-sm font-medium transition-all duration-200"
            >
              {t('invite.share')}
            </motion.button>
          </div>
        </div>

        {/* My Invites List */}
        <div className="card-base p-4">
          <h3 className="text-sm font-medium mb-4">{t('invite.myInvites')}</h3>

          {isLoading ? (
            <div className="text-center py-8">
              <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
          ) : inviteList.length > 0 ? (
            <div className="space-y-2">
              {inviteList.map((user, index) => (
                <motion.div
                  key={user.userId}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.03, duration: 0.2 }}
                  className="flex items-center gap-3 p-3 bg-accent/30 rounded-xl"
                >
                  <div className="w-9 h-9 rounded-full bg-card flex items-center justify-center text-lg">
                    👤
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium truncate">{user.nickname}</p>
                    <p className="text-xs text-muted-foreground">{user.inviteTime}</p>
                  </div>
                </motion.div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8 text-muted-foreground">
              <p className="text-3xl mb-2">👥</p>
              <p className="text-sm">No invites yet</p>
              <p className="text-xs mt-1">Share your code to get started!</p>
            </div>
          )}
        </div>

        {/* Invite Rewards Info */}
        <div className="card-base p-4 bg-success/5">
          <h3 className="text-sm font-medium mb-3 flex items-center gap-2">
            <span>🎁</span>
            <span>Invite Rewards</span>
          </h3>
          <ul className="space-y-2 text-xs text-muted-foreground">
            <li>• Get RM 5 for each successful invite</li>
            <li>• Earn 10% commission from their check-ins</li>
            <li>• Unlock bonus rewards at 10, 50, 100 invites</li>
            <li>• Top inviters win monthly prizes</li>
          </ul>
        </div>
      </div>
    </div>
  );
}
