'use client';

import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { motion } from 'framer-motion';

interface TopBarProps {
  title?: string;
  showBack?: boolean;
  actions?: React.ReactNode;
}

export function TopBar({ title, showBack = false, actions }: TopBarProps) {
  const t = useTranslations();
  const router = useRouter();

  return (
    <div className="sticky top-0 left-0 right-0 bg-card/80 backdrop-blur-lg border-b border-border z-40">
      <div className="flex items-center justify-between h-14 px-4 max-w-screen-lg mx-auto">
        {/* Left */}
        <div className="flex items-center gap-2 min-w-0">
          {showBack && (
            <motion.button
              whileTap={{ scale: 0.9 }}
              onClick={() => router.back()}
              className="p-2 -ml-2 hover:bg-accent rounded-lg transition-colors"
            >
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M15 19l-7-7 7-7"
                />
              </svg>
            </motion.button>
          )}
          {title && (
            <h1 className="text-lg font-semibold truncate">{title}</h1>
          )}
        </div>

        {/* Right */}
        {actions && <div className="flex items-center gap-2">{actions}</div>}
      </div>
    </div>
  );
}
