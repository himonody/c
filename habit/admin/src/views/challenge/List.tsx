import React, { useState, useEffect } from 'react'
import { Table, Button, Space, message, Modal } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import type { ColumnsType } from 'antd/es/table'
import { getChallengeList, deleteChallenge, type ChallengeInfo } from '@/api/challenge'

const ChallengeList: React.FC = () => {
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState<ChallengeInfo[]>([])
  const [total, setTotal] = useState(0)
  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const navigate = useNavigate()

  const fetchData = async () => {
    setLoading(true)
    try {
      const response = await getChallengeList({
        page: currentPage,
        pageSize,
      })
      setData(response.data.list)
      setTotal(response.data.total)
    } catch (error) {
      message.error('获取数据失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [currentPage, pageSize])

  const handleEdit = (record: ChallengeInfo) => {
    navigate(`/challenge/edit/${record.id}`)
  }

  const handleAdd = () => {
    navigate('/challenge/edit')
  }

  const handleDelete = (record: ChallengeInfo) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除挑战 ID: ${record.id} 吗？`,
      onOk: async () => {
        try {
          await deleteChallenge(record.id)
          message.success('删除成功')
          fetchData()
        } catch (error) {
          message.error('删除失败')
        }
      },
    })
  }

  const columns: ColumnsType<ChallengeInfo> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '自动结算',
      dataIndex: 'isAutoSettle',
      key: 'isAutoSettle',
      render: (value: boolean) => (value ? '是' : '否'),
      width: 100,
    },
    {
      title: '结算时间',
      dataIndex: 'settleTime',
      key: 'settleTime',
      width: 100,
    },
    {
      title: '挑战天数',
      dataIndex: 'cycleDays',
      key: 'cycleDays',
      width: 100,
    },
    {
      title: '开始时间',
      dataIndex: 'startTime',
      key: 'startTime',
      width: 100,
    },
    {
      title: '结束时间',
      dataIndex: 'endTime',
      key: 'endTime',
      width: 100,
    },
    {
      title: '最高收益',
      dataIndex: 'maxDailyProfit',
      key: 'maxDailyProfit',
      width: 100,
    },
    {
      title: '平台补贴',
      dataIndex: 'dailyPlatformSubsidy',
      key: 'dailyPlatformSubsidy',
      width: 100,
    },
    {
      title: '操作',
      key: 'action',
      fixed: 'right',
      width: 150,
      render: (_, record) => (
        <Space size="middle">
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Button
            type="link"
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleDelete(record)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2>挑战管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          新增挑战
        </Button>
      </div>
      
      <Table
        columns={columns}
        dataSource={data}
        rowKey="id"
        loading={loading}
        pagination={{
          current: currentPage,
          pageSize,
          total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条记录`,
          onChange: (page, size) => {
            setCurrentPage(page)
            setPageSize(size)
          },
        }}
        scroll={{ x: 1200 }}
      />
    </div>
  )
}

export default ChallengeList
