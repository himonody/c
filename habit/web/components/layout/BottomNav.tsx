'use client';

import { useTranslations } from 'next-intl';
import { usePathname } from 'next/navigation';
import Link from 'next/link';
import { motion } from 'framer-motion';
import { cn } from '@/lib/utils';

interface NavItem {
  key: string;
  icon: string;
  href: string;
}

export function BottomNav({ locale }: { locale: string }) {
  const t = useTranslations();
  const pathname = usePathname();

  const navItems: NavItem[] = [
    { key: 'home', icon: '🏠', href: `/${locale}` },
    { key: 'checkin', icon: '✓', href: `/${locale}/checkin` },
    { key: 'ranking', icon: '🏆', href: `/${locale}/ranking` },
    { key: 'profile', icon: '👤', href: `/${locale}/profile` },
  ];

  const isActive = (href: string) => {
    if (href === `/${locale}`) {
      return pathname === href;
    }
    return pathname.startsWith(href);
  };

  return (
    <nav className="fixed bottom-0 left-0 right-0 bg-card border-t border-border z-50 md:hidden">
      <div className="flex items-center justify-around h-14 max-w-container mx-auto">
        {navItems.map((item) => {
          const active = isActive(item.href);
          return (
            <Link
              key={item.key}
              href={item.href}
              className="flex-1 flex flex-col items-center justify-center relative h-full"
            >
              <motion.div
                whileTap={{ scale: 0.9 }}
                transition={{ duration: 0.15 }}
                className="flex flex-col items-center gap-0.5"
              >
                <div className={cn(
                  'text-xl transition-all duration-200',
                  active ? 'scale-110' : 'scale-100 opacity-60'
                )}>
                  {item.icon}
                </div>
                <span className={cn(
                  'text-[10px] font-medium transition-all duration-200',
                  active ? 'text-primary' : 'text-muted-foreground'
                )}>
                  {t(`nav.${item.key}`)}
                </span>
              </motion.div>
              {active && (
                <motion.div
                  layoutId="bottomNavIndicator"
                  className="absolute top-0 left-1/2 -translate-x-1/2 w-8 h-0.5 bg-primary rounded-b-full"
                  transition={{ type: 'spring', bounce: 0.2, duration: 0.25 }}
                />
              )}
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
