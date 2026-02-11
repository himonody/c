import { create } from 'zustand';

interface InviteStats {
  todayInviteCount: number;
  totalInviteCount: number;
  inviteIncome: string;
}

interface InviteUser {
  userId: number;
  username: string;
  nickname: string;
  avatar: string;
  inviteTime: string;
}

interface InviteState {
  stats: InviteStats | null;
  inviteList: InviteUser[];
  inviteCode: string | null;
  isLoading: boolean;
  error: string | null;

  setStats: (stats: InviteStats | null) => void;
  setInviteList: (inviteList: InviteUser[]) => void;
  setInviteCode: (inviteCode: string | null) => void;
  setLoading: (isLoading: boolean) => void;
  setError: (error: string | null) => void;
  reset: () => void;
}

const initialState = {
  stats: null,
  inviteList: [],
  inviteCode: null,
  isLoading: false,
  error: null,
};

export const useInviteStore = create<InviteState>((set) => ({
  ...initialState,

  setStats: (stats) => set({ stats }),

  setInviteList: (inviteList) => set({ inviteList }),

  setInviteCode: (inviteCode) => set({ inviteCode }),

  setLoading: (isLoading) => set({ isLoading }),

  setError: (error) => set({ error }),

  reset: () => set(initialState),
}));
