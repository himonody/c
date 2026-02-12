import React from 'react'
import { useRoutes } from 'react-router-dom'

import { routes } from './routers'
import AppLayout from './layouts/AppLayout'
import AuthLayout from './layouts/AuthLayout'
import ErrorBoundary from './components/ErrorBoundary'

const App: React.FC = () => {
  const element = useRoutes(routes)

  // 根据路由判断使用哪个布局
  const isAuthRoute = (pathname: string) => {
    return pathname === '/login' || pathname === '/404'
  }

  const LayoutComponent = isAuthRoute(window.location.pathname) ? AuthLayout : AppLayout

  return (
    <ErrorBoundary>
      <LayoutComponent>
        {element}
      </LayoutComponent>
    </ErrorBoundary>
  )
}

export default App
