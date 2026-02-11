type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE';

interface ApiOptions {
  method?: HttpMethod;
  body?: any;
  headers?: Record<string, string>;
}

interface ApiResponse<T = any> {
  code: number;
  msg: string;
  data?: T;
}

class ApiClient {
  private baseURL: string;

  constructor(baseURL: string = '/api') {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: ApiOptions = {}
  ): Promise<ApiResponse<T>> {
    const { method = 'POST', body, headers = {} } = options;

    const config: RequestInit = {
      method,
      headers: {
        'Content-Type': 'application/json',
        'X-Lang': this.getLanguage(),
        ...headers,
      },
      credentials: 'include', // 发送 HttpOnly Cookie
    };

    if (body) {
      config.body = JSON.stringify(body);
    }

    const url = `${this.baseURL}${endpoint}`;
    const response = await fetch(url, config);

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  private getLanguage(): string {
    if (typeof window === 'undefined') return 'en';

    const pathname = window.location.pathname;
    const locale = pathname.split('/')[1];

    if (locale === 'ms') return 'ms';
    if (locale === 'zh-CN') return 'zh-CN';
    return 'en';
  }

  // Auth APIs
  auth = {
    register: (username: string, password: string) =>
      this.request('/app/auth/register', {
        body: { username, password },
      }),

    login: (username: string, password: string) =>
      this.request<{ token: string }>('/app/auth/login', {
        body: { username, password },
      }),

    logout: () =>
      this.request('/app/auth/logout'),

    changePassword: (oldPassword: string, newPassword: string) =>
      this.request('/app/auth/change-password', {
        body: { oldPassword, newPassword },
      }),

    setPayPassword: (oldPayPassword: string, newPayPassword: string) =>
      this.request('/app/auth/set-pay-password', {
        body: { oldPayPassword, newPayPassword },
      }),

    updateProfile: (nickname?: string, avatar?: string) =>
      this.request('/app/auth/update-profile', {
        body: { nickname, avatar },
      }),

    getUserInfo: () =>
      this.request('/app/auth/user-info'),
  };

  // Wallet APIs
  wallet = {
    getInfo: () =>
      this.request('/app/wallet/info'),
  };

  // Challenge APIs
  challenge = {
    list: () =>
      this.request('/app/challenge/list'),

    detail: (challengeId: number) =>
      this.request('/app/challenge/detail', {
        body: { challengeId },
      }),

    join: (challengeId: number, payPassword: string) =>
      this.request('/app/challenge/join', {
        body: { challengeId, payPassword },
      }),

    myChallenges: () =>
      this.request('/app/challenge/my-challenges'),

    quit: (challengeId: number, payPassword: string) =>
      this.request('/app/challenge/quit', {
        body: { challengeId, payPassword },
      }),

    participants: (challengeId: number, page: number = 1, pageSize: number = 20) =>
      this.request('/app/challenge/participants', {
        body: { challengeId, page, pageSize },
      }),
  };

  // Check-in APIs
  checkin = {
    checkIn: (challengeId: number) =>
      this.request('/app/checkin/check-in', {
        body: { challengeId },
      }),

    todayStats: (challengeId: number) =>
      this.request('/app/checkin/today-stats', {
        body: { challengeId },
      }),

    todayList: (challengeId: number, page: number = 1, pageSize: number = 20) =>
      this.request('/app/checkin/today-list', {
        body: { challengeId, page, pageSize },
      }),

    calendar: (challengeId: number, year: number, month: number) =>
      this.request('/app/checkin/calendar', {
        body: { challengeId, year, month },
      }),

    history: (challengeId: number, page: number = 1, pageSize: number = 20) =>
      this.request('/app/checkin/history', {
        body: { challengeId, page, pageSize },
      }),

    supplement: (challengeId: number, date: string) =>
      this.request('/app/checkin/supplement', {
        body: { challengeId, date },
      }),
  };

  // Ranking APIs
  ranking = {
    invite: (page: number = 1, pageSize: number = 20) =>
      this.request('/app/ranking/invite', {
        body: { page, pageSize },
      }),

    wealth: (page: number = 1, pageSize: number = 20) =>
      this.request('/app/ranking/wealth', {
        body: { page, pageSize },
      }),

    persistence: (page: number = 1, pageSize: number = 20) =>
      this.request('/app/ranking/persistence', {
        body: { page, pageSize },
      }),
  };

  // Invite APIs
  invite = {
    stats: () =>
      this.request('/app/invite/stats'),

    myInvites: (page: number = 1, pageSize: number = 20) =>
      this.request('/app/invite/my-invites', {
        body: { page, pageSize },
      }),

    getCode: () =>
      this.request('/app/invite/get-code'),
  };

  // Transaction APIs
  transaction = {
    list: (type?: string, page: number = 1, pageSize: number = 20) =>
      this.request('/app/transaction/list', {
        body: { type, page, pageSize },
      }),

    todayStats: () =>
      this.request('/app/transaction/today-stats'),
  };

  // Withdraw APIs
  withdraw = {
    apply: (amount: string, payPassword: string, withdrawType: string, account: string) =>
      this.request('/app/withdraw/apply', {
        body: { amount, payPassword, withdrawType, account },
      }),

    list: (page: number = 1, pageSize: number = 20) =>
      this.request('/app/withdraw/list', {
        body: { page, pageSize },
      }),
  };
}

export const api = new ApiClient();
export type { ApiResponse };
