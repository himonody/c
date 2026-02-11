import { create } from 'zustand';

interface Challenge {
  challengeId: number;
  name: string;
  duration: number;
  startTime: string;
  endTime: string;
  checkInStartTime: string;
  checkInEndTime: string;
  entryFee: string;
  prizePool: string;
  participantCount: number;
  status: string;
  rules: string;
}

interface ChallengeState {
  challenges: Challenge[];
  currentChallenge: Challenge | null;
  myChallenges: Challenge[];
  isLoading: boolean;
  error: string | null;

  setChallenges: (challenges: Challenge[]) => void;
  setCurrentChallenge: (challenge: Challenge | null) => void;
  setMyChallenges: (myChallenges: Challenge[]) => void;
  setLoading: (isLoading: boolean) => void;
  setError: (error: string | null) => void;
  reset: () => void;
}

const initialState = {
  challenges: [],
  currentChallenge: null,
  myChallenges: [],
  isLoading: false,
  error: null,
};

export const useChallengeStore = create<ChallengeState>((set) => ({
  ...initialState,

  setChallenges: (challenges) => set({ challenges }),

  setCurrentChallenge: (currentChallenge) => set({ currentChallenge }),

  setMyChallenges: (myChallenges) => set({ myChallenges }),

  setLoading: (isLoading) => set({ isLoading }),

  setError: (error) => set({ error }),

  reset: () => set(initialState),
}));
