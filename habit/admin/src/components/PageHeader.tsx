import React from 'react'
import { Card, Space, Button } from 'antd'
import { ReloadOutlined } from '@ant-design/icons'

interface PageHeaderProps {
  title: string
  subtitle?: string
  extra?: React.ReactNode
  onRefresh?: () => void
  showBack?: boolean
  onBack?: () => void
}

const PageHeader: React.FC<PageHeaderProps> = ({
  title,
  subtitle,
  extra,
  onRefresh,
  showBack = false,
  onBack
}) => {
  return (
    <Card
      style={{ marginBottom: 16 }}
      bodyStyle={{ padding: '16px 24px' }}
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <Space>
            {showBack && (
              <Button onClick={onBack}>
                ← 返回
              </Button>
            )}
            <div>
              <h2 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>
                {title}
              </h2>
              {subtitle && (
                <div style={{ color: '#666', fontSize: 14, marginTop: 4 }}>
                  {subtitle}
                </div>
              )}
            </div>
          </Space>
        </div>
        <Space>
          {onRefresh && (
            <Button icon={<ReloadOutlined />} onClick={onRefresh}>
              刷新
            </Button>
          )}
          {extra}
        </Space>
      </div>
    </Card>
  )
}

export default PageHeader
