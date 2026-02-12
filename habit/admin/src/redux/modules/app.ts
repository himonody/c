import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface AppState {
  collapsed: boolean
  theme: 'light' | 'dark'
  language: 'zh-CN' | 'en-US'
}

const initialState: AppState = {
  collapsed: false,
  theme: 'light',
  language: 'zh-CN',
}

const appSlice = createSlice({
  name: 'app',
  initialState,
  reducers: {
    toggleSidebar: (state) => {
      state.collapsed = !state.collapsed
    },
    setTheme: (state, action: PayloadAction<'light' | 'dark'>) => {
      state.theme = action.payload
    },
    setLanguage: (state, action: PayloadAction<'zh-CN' | 'en-US'>) => {
      state.language = action.payload
    },
  },
})

export const { toggleSidebar, setTheme, setLanguage } = appSlice.actions
export default appSlice.reducer
