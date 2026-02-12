import React, { Suspense } from 'react'
import { Spin } from 'antd'
import ErrorBoundary from './ErrorBoundary'

interface LazyLoadProps {
  children: React.ReactNode
  fallback?: React.ReactNode
}

const LazyLoad: React.FC<LazyLoadProps> = ({ 
  children, 
  fallback = <Spin size="large" style={{ display: 'block', margin: '20px auto' }} /> 
}) => {
  return (
    <ErrorBoundary>
      <Suspense fallback={fallback}>
        {children}
      </Suspense>
    </ErrorBoundary>
  )
}

export default LazyLoad
