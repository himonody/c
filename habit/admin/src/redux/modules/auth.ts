import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface AuthState {
  token: string | null
  userInfo: any | null
  isLoggedIn: boolean
}

const initialState: AuthState = {
  token: null,
  userInfo: null,
  isLoggedIn: false,
}

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setToken: (state, action: PayloadAction<string>) => {
      state.token = action.payload
      state.isLoggedIn = !!action.payload
    },
    setUserInfo: (state, action: PayloadAction<any>) => {
      state.userInfo = action.payload
    },
    logout: (state) => {
      state.token = null
      state.userInfo = null
      state.isLoggedIn = false
    },
  },
})

export const { setToken, setUserInfo, logout } = authSlice.actions
export default authSlice.reducer
export { authSlice }
