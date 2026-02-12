import React from 'react'
import { Layout, BackTop } from 'antd'
import { Outlet } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { RootState } from '@/redux'
import Sidebar from '@/components/Sidebar'
import Header from '@/components/Header'
import CustomBreadcrumb from '@/components/Breadcrumb'

const { Content } = Layout

const AppLayout: React.FC<{ children?: React.ReactNode }> = ({ children }) => {
  const { collapsed } = useSelector((state: RootState) => state.app)

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sidebar />
      <Layout>
        <Header />
        <Content
          style={{
            margin: '16px',
            padding: '16px',
            background: '#fff',
            borderRadius: '8px',
            minHeight: 'calc(100vh - 112px)',
            marginLeft: collapsed ? 80 : 200,
            transition: 'margin-left 0.2s',
          }}
        >
          <CustomBreadcrumb />
          {children || <Outlet />}
        </Content>
      </Layout>
      <BackTop />
    </Layout>
  )
}

export default AppLayout
