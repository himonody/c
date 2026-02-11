import { create } from 'zustand';

interface TodayStats {
  checkedInCount: number;
  notCheckedInCount: number;
  totalCount: number;
}

interface CheckInRecord {
  userId: number;
  username: string;
  nickname: string;
  avatar: string;
  checkInTime: string;
  status: string; // 'success' | 'late' | 'missed'
}

interface CalendarDay {
  date: string;
  status: string; // 'checked' | 'unchecked' | 'supplemented'
  amount?: string;
  loss?: string;
}

interface CheckInState {
  todayStats: TodayStats | null;
  todayList: CheckInRecord[];
  calendar: CalendarDay[];
  history: CheckInRecord[];
  hasCheckedInToday: boolean;
  isLoading: boolean;
  error: string | null;

  setTodayStats: (stats: TodayStats | null) => void;
  setTodayList: (list: CheckInRecord[]) => void;
  setCalendar: (calendar: CalendarDay[]) => void;
  setHistory: (history: CheckInRecord[]) => void;
  setHasCheckedInToday: (hasChecked: boolean) => void;
  setLoading: (isLoading: boolean) => void;
  setError: (error: string | null) => void;
  reset: () => void;
}

const initialState = {
  todayStats: null,
  todayList: [],
  calendar: [],
  history: [],
  hasCheckedInToday: false,
  isLoading: false,
  error: null,
};

export const useCheckinStore = create<CheckInState>((set) => ({
  ...initialState,

  setTodayStats: (todayStats) => set({ todayStats }),

  setTodayList: (todayList) => set({ todayList }),

  setCalendar: (calendar) => set({ calendar }),

  setHistory: (history) => set({ history }),

  setHasCheckedInToday: (hasCheckedInToday) => set({ hasCheckedInToday }),

  setLoading: (isLoading) => set({ isLoading }),

  setError: (error) => set({ error }),

  reset: () => set(initialState),
}));
