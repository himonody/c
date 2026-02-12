import { configureStore } from '@reduxjs/toolkit'
import { persistStore, persistReducer } from 'redux-persist'
import storage from 'redux-persist/lib/storage'
import { combineReducers } from '@reduxjs/toolkit'

import authReducer from './modules/auth'
import userReducer from './modules/user'
import appReducer from './modules/app'

const persistConfig = {
  key: 'root',
  storage,
  whitelist: ['auth', 'app'], // 只持久化 auth 和 app 模块
}

const rootReducer = combineReducers({
  auth: authReducer,
  user: userReducer,
  app: appReducer,
})

const persistedReducer = persistReducer(persistConfig, rootReducer)

export const store = configureStore({
  reducer: persistedReducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: ['persist/PERSIST', 'persist/REHYDRATE'],
      },
    }),
})

export const persistor = persistStore(store)

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
