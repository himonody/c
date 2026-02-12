import React from 'react'
import ReactDOM from 'react-dom/client'
import { Provider } from 'react-redux'
import { BrowserRouter } from 'react-router-dom'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import { PersistGate } from 'redux-persist/integration/react'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'

import App from './App'
import { store, persistor } from './redux'
import './styles/index.less'
import { reportWebVitals, setupErrorMonitoring, monitorResourceLoading, monitorLongTasks } from './utils/performance'

// 设置 dayjs 中文
dayjs.locale('zh-cn')

// 设置错误监控
setupErrorMonitoring()

// 设置性能监控
monitorResourceLoading()
monitorLongTasks()

// 性能指标上报
reportWebVitals(console.log)

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <Provider store={store}>
      <PersistGate loading={null} persistor={persistor}>
        <BrowserRouter>
          <ConfigProvider locale={zhCN}>
            <App />
          </ConfigProvider>
        </BrowserRouter>
      </PersistGate>
    </Provider>
  </React.StrictMode>,
)
