// 性能监控工具

// 页面加载性能监控
export const reportWebVitals = (onPerfEntry?: (metric: any) => void) => {
  if (onPerfEntry && onPerfEntry instanceof Function) {
    import('web-vitals').then(({ getCLS, getFID, getFCP, getLCP, getTTFB }) => {
      getCLS(onPerfEntry)
      getFID(onPerfEntry)
      getFCP(onPerfEntry)
      getLCP(onPerfEntry)
      getTTFB(onPerfEntry)
    })
  }
}

// 错误监控
export const setupErrorMonitoring = () => {
  // 全局错误捕获
  window.addEventListener('error', (event) => {
    console.error('Global error:', event.error)
    // 这里可以发送错误到监控服务
  })

  // Promise 错误捕获
  window.addEventListener('unhandledrejection', (event) => {
    console.error('Unhandled promise rejection:', event.reason)
    // 这里可以发送错误到监控服务
  })
}

// 性能指标收集
export const collectPerformanceMetrics = () => {
  const metrics = {
    // 页面加载时间
    loadTime: performance.timing.loadEventEnd - performance.timing.navigationStart,
    // DNS 查询时间
    dnsTime: performance.timing.domainLookupEnd - performance.timing.domainLookupStart,
    // TCP 连接时间
    tcpTime: performance.timing.connectEnd - performance.timing.connectStart,
    // 请求响应时间
    requestTime: performance.timing.responseEnd - performance.timing.requestStart,
    // DOM 解析时间
    domParseTime: performance.timing.domComplete - performance.timing.domLoading,
  }

  return metrics
}

// 资源加载监控
export const monitorResourceLoading = () => {
  const observer = new PerformanceObserver((list) => {
    const entries = list.getEntries()
    entries.forEach((entry) => {
      if (entry.entryType === 'resource') {
        const resource = entry as PerformanceResourceTiming
        console.log('Resource loaded:', {
          name: resource.name,
          duration: resource.duration,
          size: resource.transferSize,
          type: resource.initiatorType,
        })
      }
    })
  })

  observer.observe({ entryTypes: ['resource'] })
}

// 长任务监控
export const monitorLongTasks = () => {
  if ('PerformanceObserver' in window) {
    const observer = new PerformanceObserver((list) => {
      const entries = list.getEntries()
      entries.forEach((entry) => {
        if (entry.duration > 50) { // 超过 50ms 的任务
          console.warn('Long task detected:', {
            duration: entry.duration,
            startTime: entry.startTime,
          })
        }
      })
    })

    try {
      observer.observe({ entryTypes: ['longtask'] })
    } catch (e) {
      console.log('Long task monitoring not supported')
    }
  }
}
