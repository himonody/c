'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { motion } from 'framer-motion';
import { useCheckinStore } from '@/store/checkinStore';
import { api } from '@/lib/api';
import { cn } from '@/lib/utils';

export default function CalendarPage() {
  const t = useTranslations();
  const { calendar, setCalendar } = useCheckinStore();
  const [currentDate, setCurrentDate] = useState(new Date());
  const [challengeId] = useState(1);

  useEffect(() => {
    loadCalendar();
  }, [currentDate]);

  const loadCalendar = async () => {
    try {
      const response = await api.checkin.calendar(
        challengeId,
        currentDate.getFullYear(),
        currentDate.getMonth() + 1
      );
      if (response.code === 0 && response.data) {
        setCalendar(response.data);
      }
    } catch (error) {
      console.error('Failed to load calendar:', error);
    }
  };

  const getDaysInMonth = () => {
    const year = currentDate.getFullYear();
    const month = currentDate.getMonth();
    const firstDay = new Date(year, month, 1).getDay();
    const daysInMonth = new Date(year, month + 1, 0).getDate();

    const days = [];
    for (let i = 0; i < firstDay; i++) {
      days.push(null);
    }
    for (let i = 1; i <= daysInMonth; i++) {
      days.push(i);
    }
    return days;
  };

  const getStatusForDay = (day: number) => {
    const dateStr = `${currentDate.getFullYear()}-${String(currentDate.getMonth() + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
    return calendar.find((d) => d.date === dateStr);
  };

  const previousMonth = () => {
    setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() - 1));
  };

  const nextMonth = () => {
    setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() + 1));
  };

  const weekDays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
  const monthYear = currentDate.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });

  return (
    <div className="min-h-screen">
      <div className="max-w-container mx-auto p-4 space-y-6 pb-20 md:pb-6">
        {/* Month Navigation */}
        <div className="flex items-center justify-between pt-6 pb-2">
          <motion.button
            whileTap={{ scale: 0.9 }}
            transition={{ duration: 0.15 }}
            onClick={previousMonth}
            className="p-2 bg-accent/30 rounded-xl text-lg"
          >
            ←
          </motion.button>
          <h2 className="text-lg font-medium">{monthYear}</h2>
          <motion.button
            whileTap={{ scale: 0.9 }}
            transition={{ duration: 0.15 }}
            onClick={nextMonth}
            className="p-2 bg-accent/30 rounded-xl text-lg"
          >
            →
          </motion.button>
        </div>

        {/* Calendar Grid */}
        <div className="card-base p-4">
          {/* Week Days */}
          <div className="grid grid-cols-7 gap-1 mb-2">
            {weekDays.map((day) => (
              <div key={day} className="text-center text-xs font-medium text-muted-foreground py-2">
                {day}
              </div>
            ))}
          </div>

          {/* Days */}
          <div className="grid grid-cols-7 gap-1">
            {getDaysInMonth().map((day, index) => {
              if (day === null) {
                return <div key={`empty-${index}`} className="aspect-square" />;
              }

              const status = getStatusForDay(day);
              const isToday =
                day === new Date().getDate() &&
                currentDate.getMonth() === new Date().getMonth() &&
                currentDate.getFullYear() === new Date().getFullYear();

              return (
                <motion.div
                  key={day}
                  whileTap={{ scale: 0.95 }}
                  transition={{ duration: 0.15 }}
                  className={cn(
                    'aspect-square flex flex-col items-center justify-center rounded-xl text-xs',
                    'transition-all duration-200',
                    isToday && 'ring-2 ring-primary',
                    !status && 'bg-accent/30',
                    status?.status === 'checked' && 'bg-success/20 text-success',
                    status?.status === 'unchecked' && 'bg-warning/20 text-warning',
                    status?.status === 'supplemented' && 'bg-primary/20 text-primary'
                  )}
                >
                  <span className="font-medium">{day}</span>
                  {status && (
                    <span className="text-[10px] mt-0.5">
                      {status.status === 'checked' && '✓'}
                      {status.status === 'unchecked' && '✗'}
                      {status.status === 'supplemented' && '↻'}
                    </span>
                  )}
                </motion.div>
              );
            })}
          </div>
        </div>

        {/* Legend */}
        <div className="grid grid-cols-3 gap-2">
          <div className="flex items-center gap-2 text-xs">
            <div className="w-4 h-4 rounded bg-success/20" />
            <span className="text-muted-foreground">{t('checkin.status.success')}</span>
          </div>
          <div className="flex items-center gap-2 text-xs">
            <div className="w-4 h-4 rounded bg-warning/20" />
            <span className="text-muted-foreground">{t('checkin.status.missed')}</span>
          </div>
          <div className="flex items-center gap-2 text-xs">
            <div className="w-4 h-4 rounded bg-primary/20" />
            <span className="text-muted-foreground">{t('checkin.status.supplemented')}</span>
          </div>
        </div>

        {/* Monthly Stats */}
        <div className="card-base p-4">
          <h3 className="text-sm font-medium mb-4">Monthly Summary</h3>
          <div className="grid grid-cols-2 gap-4">
            <div className="text-center p-3 bg-accent/30 rounded-xl">
              <p className="text-xs text-muted-foreground mb-1">{t('checkin.amount')}</p>
              <p className="text-2xl font-semibold text-success">
                RM {calendar.reduce((sum, day) => sum + parseFloat(day.amount || '0'), 0).toFixed(2)}
              </p>
            </div>
            <div className="text-center p-3 bg-accent/30 rounded-xl">
              <p className="text-xs text-muted-foreground mb-1">{t('checkin.loss')}</p>
              <p className="text-2xl font-semibold text-warning">
                RM {calendar.reduce((sum, day) => sum + parseFloat(day.loss || '0'), 0).toFixed(2)}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
