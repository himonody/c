import React, { useState, useEffect } from 'react'
import { Card, Row, Col, Statistic, Table, Tag, Progress } from 'antd'
import { UserOutlined, TrophyOutlined, RiseOutlined, ClockCircleOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'

interface RecentActivity {
  id: number
  type: string
  description: string
  user: string
  time: string
  status: 'success' | 'pending' | 'failed'
}

interface ChallengeStats {
  id: number
  name: string
  participants: number
  successRate: number
  poolAmount: number
  status: 'active' | 'completed' | 'pending'
}

const Dashboard: React.FC = () => {
  const [recentActivities, setRecentActivities] = useState<RecentActivity[]>([])
  const [challengeStats, setChallengeStats] = useState<ChallengeStats[]>([])

  useEffect(() => {
    // 模拟数据加载
    const mockActivities: RecentActivity[] = [
      { id: 1, type: '打卡', description: '完成今日打卡挑战', user: '张三', time: '2分钟前', status: 'success' },
      { id: 2, type: '结算', description: '挑战结算完成', user: '系统', time: '5分钟前', status: 'success' },
      { id: 3, type: '注册', description: '新用户注册', user: '李四', time: '10分钟前', status: 'success' },
      { id: 4, type: '提现', description: '申请提现', user: '王五', time: '15分钟前', status: 'pending' },
      { id: 5, type: '打卡', description: '打卡失败', user: '赵六', time: '20分钟前', status: 'failed' },
    ]

    const mockChallengeStats: ChallengeStats[] = [
      { id: 1, name: '30天早起挑战', participants: 156, successRate: 85.2, poolAmount: 15600, status: 'active' },
      { id: 2, name: '21天健身挑战', participants: 89, successRate: 78.5, poolAmount: 8900, status: 'active' },
      { id: 3, name: '7天阅读挑战', participants: 234, successRate: 92.3, poolAmount: 23400, status: 'completed' },
    ]

    setRecentActivities(mockActivities)
    setChallengeStats(mockChallengeStats)
  }, [])

  const getStatusColor = (status: string) => {
    const colorMap: Record<string, string> = {
      success: 'success',
      pending: 'processing',
      failed: 'error',
      active: 'processing',
      completed: 'success',
    }
    return colorMap[status] || 'default'
  }

  const activityColumns: ColumnsType<RecentActivity> = [
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      width: 80,
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: '用户',
      dataIndex: 'user',
      key: 'user',
      width: 100,
    },
    {
      title: '时间',
      dataIndex: 'time',
      key: 'time',
      width: 100,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => (
        <Tag color={getStatusColor(status)}>
          {status === 'success' ? '成功' : status === 'pending' ? '处理中' : '失败'}
        </Tag>
      ),
    },
  ]

  const challengeColumns: ColumnsType<ChallengeStats> = [
    {
      title: '挑战名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '参与人数',
      dataIndex: 'participants',
      key: 'participants',
      width: 100,
    },
    {
      title: '成功率',
      dataIndex: 'successRate',
      key: 'successRate',
      width: 120,
      render: (rate: number) => (
        <Progress percent={rate} size="small" format={() => `${rate}%`} />
      ),
    },
    {
      title: '奖池金额',
      dataIndex: 'poolAmount',
      key: 'poolAmount',
      width: 120,
      render: (amount: number) => `¥${(amount / 100).toFixed(2)}`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => (
        <Tag color={getStatusColor(status)}>
          {status === 'active' ? '进行中' : status === 'completed' ? '已完成' : '待开始'}
        </Tag>
      ),
    },
  ]

  return (
    <div>
      <h1>仪表盘</h1>
      
      {/* 统计卡片 */}
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="总用户数"
              value={1128}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#3f8600' }}
              suffix={
                <span style={{ fontSize: 14, color: '#52c41a' }}>
                  ↑ 12.5%
                </span>
              }
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="今日挑战"
              value={93}
              prefix={<TrophyOutlined />}
              valueStyle={{ color: '#1890ff' }}
              suffix={
                <span style={{ fontSize: 14, color: '#52c41a' }}>
                  ↑ 8.2%
                </span>
              }
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="奖池金额"
              value={12890}
              prefix="¥"
              valueStyle={{ color: '#722ed1' }}
              precision={2}
              suffix={
                <span style={{ fontSize: 14, color: '#52c41a' }}>
                  ↑ 15.3%
                </span>
              }
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="平均成功率"
              value={85.6}
              prefix={<RiseOutlined />}
              valueStyle={{ color: '#eb2f96' }}
              suffix={
                <span style={{ fontSize: 14, color: '#52c41a' }}>
                  ↑ 2.1%
                </span>
              }
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={16}>
        {/* 最近活动 */}
        <Col span={12}>
          <Card title="最近活动" style={{ height: 400 }}>
            <Table
              columns={activityColumns}
              dataSource={recentActivities}
              rowKey="id"
              pagination={false}
              size="small"
              scroll={{ y: 300 }}
            />
          </Card>
        </Col>

        {/* 挑战统计 */}
        <Col span={12}>
          <Card title="挑战统计" style={{ height: 400 }}>
            <Table
              columns={challengeColumns}
              dataSource={challengeStats}
              rowKey="id"
              pagination={false}
              size="small"
              scroll={{ y: 300 }}
            />
          </Card>
        </Col>
      </Row>

      {/* 系统信息 */}
      <Row gutter={16} style={{ marginTop: 24 }}>
        <Col span={24}>
          <Card title="系统信息" extra={<ClockCircleOutlined />}>
            <Row gutter={16}>
              <Col span={6}>
                <Statistic title="系统运行时间" value="7" suffix="天" />
              </Col>
              <Col span={6}>
                <Statistic title="API 调用次数" value={128476} />
              </Col>
              <Col span={6}>
                <Statistic title="数据库连接数" value={23} />
              </Col>
              <Col span={6}>
                <Statistic title="内存使用率" value={68.5} suffix="%" />
              </Col>
            </Row>
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default Dashboard
