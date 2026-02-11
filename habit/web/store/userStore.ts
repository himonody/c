import { create } from 'zustand';

interface UserInfo {
  userId: number;
  username: string;
  nickname: string;
  avatar: string;
}

interface WalletInfo {
  balance: string;
  frozen: string;
  experienceBalance: string;
  challengeBalance: string;
  totalIncome: string;
}

interface UserState {
  userInfo: UserInfo | null;
  walletInfo: WalletInfo | null;
  isLoading: boolean;
  error: string | null;

  setUserInfo: (userInfo: UserInfo | null) => void;
  setWalletInfo: (walletInfo: WalletInfo | null) => void;
  setLoading: (isLoading: boolean) => void;
  setError: (error: string | null) => void;
  reset: () => void;
}

const initialState = {
  userInfo: null,
  walletInfo: null,
  isLoading: false,
  error: null,
};

export const useUserStore = create<UserState>((set) => ({
  ...initialState,

  setUserInfo: (userInfo) => set({ userInfo }),

  setWalletInfo: (walletInfo) => set({ walletInfo }),

  setLoading: (isLoading) => set({ isLoading }),

  setError: (error) => set({ error }),

  reset: () => set(initialState),
}));
