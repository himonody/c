'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { motion } from 'framer-motion';
import Link from 'next/link';
import { useUserStore } from '@/store/userStore';
import { api } from '@/lib/api';
import { cn } from '@/lib/utils';

export default function WalletPage() {
  const t = useTranslations();
  const { walletInfo } = useUserStore();
  const [amount, setAmount] = useState('');
  const [payPassword, setPayPassword] = useState('');
  const [withdrawType, setWithdrawType] = useState('wechat');
  const [account, setAccount] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleWithdraw = async () => {
    if (!amount || parseFloat(amount) <= 0) {
      setError('Please enter a valid amount');
      return;
    }

    if (!payPassword) {
      setError('Please enter payment password');
      return;
    }

    if (!account) {
      setError('Please enter account details');
      return;
    }

    setIsLoading(true);
    setError('');
    setSuccess('');

    try {
      const response = await api.withdraw.apply(amount, payPassword, withdrawType, account);
      if (response.code === 0) {
        setSuccess(t('wallet.withdrawSuccess'));
        setAmount('');
        setPayPassword('');
        setAccount('');
      } else {
        setError(response.msg);
      }
    } catch (err) {
      setError('Failed to submit withdrawal request');
    } finally {
      setIsLoading(false);
    }
  };

  const quickAmounts = ['10', '50', '100', '500'];

  return (
    <div className="min-h-screen">
      <div className="max-w-container mx-auto p-4 space-y-6 pb-20 md:pb-6">
        {/* 标题 */}
        <div className="text-center pt-6 pb-2">
          <h1 className="text-lg font-medium text-foreground">{t('wallet.title')}</h1>
        </div>

        {/* Balance Card */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.2 }}
          className="card-large p-6 text-center"
        >
          <p className="text-xs text-muted-foreground mb-2">{t('wallet.balance')}</p>
          <p className="text-5xl font-semibold text-success mb-4">
            RM {walletInfo?.balance || '0.00'}
          </p>
          <div className="text-xs text-muted-foreground">
            {t('wallet.frozen')}: RM {walletInfo?.frozen || '0.00'}
          </div>
        </motion.div>

        {/* Withdraw Form */}
        <div className="card-base p-6 space-y-4">
          <h3 className="text-sm font-medium">{t('wallet.withdraw')}</h3>

          {/* Amount */}
          <div>
            <label className="block text-xs text-muted-foreground mb-2">
              {t('wallet.withdrawAmount')}
            </label>
            <input
              type="number"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              placeholder="0.00"
              className="w-full px-4 py-3 rounded-xl border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-primary"
            />
            <div className="grid grid-cols-4 gap-2 mt-3">
              {quickAmounts.map((qa) => (
                <button
                  key={qa}
                  onClick={() => setAmount(qa)}
                  className="py-2 text-xs rounded-xl border border-border bg-accent/30 hover:bg-accent transition-all duration-200"
                >
                  RM {qa}
                </button>
              ))}
            </div>
          </div>

          {/* Withdraw Type */}
          <div>
            <label className="block text-xs text-muted-foreground mb-2">
              {t('wallet.withdrawType')}
            </label>
            <div className="grid grid-cols-2 gap-3">
              <button
                onClick={() => setWithdrawType('wechat')}
                className={cn(
                  'py-3 rounded-xl text-sm font-medium transition-all duration-200',
                  withdrawType === 'wechat'
                    ? 'bg-primary text-primary-foreground'
                    : 'bg-accent/30'
                )}
              >
                WeChat
              </button>
              <button
                onClick={() => setWithdrawType('alipay')}
                className={cn(
                  'py-3 rounded-xl text-sm font-medium transition-all duration-200',
                  withdrawType === 'alipay'
                    ? 'bg-primary text-primary-foreground'
                    : 'bg-accent/30'
                )}
              >
                Alipay
              </button>
            </div>
          </div>

          {/* Account */}
          <div>
            <label className="block text-xs text-muted-foreground mb-2">
              {t('wallet.account')}
            </label>
            <input
              type="text"
              value={account}
              onChange={(e) => setAccount(e.target.value)}
              placeholder={`Enter ${withdrawType} account`}
              className="w-full px-4 py-3 rounded-xl border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-primary"
            />
          </div>

          {/* Payment Password */}
          <div>
            <label className="block text-xs text-muted-foreground mb-2">
              {t('auth.payPassword')}
            </label>
            <input
              type="password"
              value={payPassword}
              onChange={(e) => setPayPassword(e.target.value)}
              placeholder="Enter payment password"
              className={cn(
                'w-full px-4 py-3 rounded-xl border bg-background text-sm',
                'focus:outline-none focus:ring-2 focus:ring-primary',
                error && 'border-destructive'
              )}
              maxLength={6}
            />
          </div>

          {error && (
            <p className="text-xs text-destructive bg-destructive/10 p-3 rounded-xl">
              {error}
            </p>
          )}

          {success && (
            <p className="text-xs text-success bg-success/10 p-3 rounded-xl">
              {success}
            </p>
          )}

          {/* Submit Button */}
          <motion.button
            whileTap={{ scale: 0.97 }}
            onClick={handleWithdraw}
            disabled={isLoading}
            className="btn-primary w-full disabled:opacity-50"
          >
            {isLoading ? t('common.loading') : t('wallet.withdraw')}
          </motion.button>
        </div>

        {/* Transaction History Link */}
        <Link href="/wallet/transactions">
          <motion.div
            whileTap={{ scale: 0.98 }}
            transition={{ duration: 0.15 }}
            className="card-base p-4 flex items-center gap-3"
          >
            <div className="text-xl">📊</div>
            <div className="flex-1 text-sm font-medium">{t('wallet.transactions')}</div>
            <div className="text-muted-foreground text-xs">→</div>
          </motion.div>
        </Link>
      </div>
    </div>
  );
}
