import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface UserState {
  list: any[]
  total: number
  loading: boolean
}

const initialState: UserState = {
  list: [],
  total: 0,
  loading: false,
}

const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    setUsers: (state, action: PayloadAction<{ list: any[]; total: number }>) => {
      state.list = action.payload.list
      state.total = action.payload.total
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload
    },
  },
})

export const { setUsers, setLoading } = userSlice.actions
export default userSlice.reducer
