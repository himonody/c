import React, { useState, useEffect } from 'react'
import { Table, Button, Space, message, Modal, Tag, Input, Select, Avatar } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, SearchOutlined, ReloadOutlined, KeyOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import type { ColumnsType } from 'antd/es/table'
import { getUserList, deleteUser, resetUserPassword, type UserInfo } from '@/api/user'
import { PAGINATION_CONFIG } from '@/config'
import { UserStatus } from '@/enums'


const { Option } = Select

const User: React.FC = () => {
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState<UserInfo[]>([])
  const [total, setTotal] = useState(0)
  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [searchParams, setSearchParams] = useState({
    username: '',
    status: undefined as number | undefined,
  })
  const navigate = useNavigate()

  const fetchData = async () => {
    setLoading(true)
    try {
      const response = await getUserList({
        page: currentPage,
        pageSize,
        ...searchParams,
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

  const handleEdit = (record: UserInfo) => {
    navigate(`/user/edit/${record.id}`)
  }

  const handleAdd = () => {
    navigate('/user/edit')
  }

  const handleDelete = (record: UserInfo) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除用户 "${record.username}" 吗？`,
      onOk: async () => {
        try {
          await deleteUser(record.id)
          message.success('删除成功')
          fetchData()
        } catch (error) {
          message.error('删除失败')
        }
      },
    })
  }

  const handleResetPassword = (record: UserInfo) => {
    Modal.confirm({
      title: '重置密码',
      content: `确定要重置用户 "${record.username}" 的密码吗？新密码将设置为 "123456"`,
      onOk: async () => {
        try {
          await resetUserPassword(record.id, '123456')
          message.success('密码重置成功')
        } catch (error) {
          message.error('密码重置失败')
        }
      },
    })
  }

  const handleSearch = () => {
    setCurrentPage(1)
    fetchData()
  }

  const handleReset = () => {
    setSearchParams({ username: '', status: undefined })
    setCurrentPage(1)
    setTimeout(fetchData, 0)
  }

  const getStatusTag = (status: number) => {
    const statusMap: Record<number, { color: string; text: string }> = {
      [UserStatus.INACTIVE]: { color: 'default', text: '未激活' },
      [UserStatus.ACTIVE]: { color: 'success', text: '正常' },
      [UserStatus.BANNED]: { color: 'error', text: '已禁用' },
    }
    const config = statusMap[status] || { color: 'default', text: '未知' }
    return <Tag color={config.color}>{config.text}</Tag>
  }

  const getRoleTag = (role: number) => {
    const roleMap: Record<number, { color: string; text: string }> = {
      1: { color: 'blue', text: '超级管理员' },
      2: { color: 'green', text: '管理员' },
      3: { color: 'orange', text: '普通用户' },
    }
    const config = roleMap[role] || { color: 'default', text: '未知' }
    return <Tag color={config.color}>{config.text}</Tag>
  }

  const columns: ColumnsType<UserInfo> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '头像',
      dataIndex: 'avatar',
      key: 'avatar',
      width: 80,
      render: (avatar: string, record: UserInfo) => (
        <Avatar src={avatar} icon={!avatar && record.username?.charAt(0)?.toUpperCase()}>
          {record.username?.charAt(0)?.toUpperCase()}
        </Avatar>
      ),
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
      width: 120,
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
      width: 180,
      ellipsis: true,
    },
    {
      title: '手机号',
      dataIndex: 'phone',
      key: 'phone',
      width: 130,
    },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      width: 100,
      render: (role: number) => getRoleTag(role),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: number) => getStatusTag(status),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 180,
    },
    {
      title: '最后登录',
      dataIndex: 'lastLoginAt',
      key: 'lastLoginAt',
      width: 180,
      render: (time: string) => time || '-',
    },
    {
      title: '操作',
      key: 'action',
      fixed: 'right',
      width: 200,
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Button
            type="link"
            icon={<KeyOutlined />}
            onClick={() => handleResetPassword(record)}
          >
            重置密码
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
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2>用户管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          新增用户
        </Button>
      </div>

      <div style={{ marginBottom: 16, display: 'flex', gap: 16 }}>
        <Input
          placeholder="用户名"
          value={searchParams.username}
          onChange={(e) => setSearchParams({ ...searchParams, username: e.target.value })}
          style={{ width: 200 }}
        />
        <Select
          placeholder="用户状态"
          value={searchParams.status}
          onChange={(value) => setSearchParams({ ...searchParams, status: value })}
          style={{ width: 150 }}
          allowClear
        >
          <Option value={UserStatus.ACTIVE}>正常</Option>
          <Option value={UserStatus.INACTIVE}>未激活</Option>
          <Option value={UserStatus.BANNED}>已禁用</Option>
        </Select>
        <Space>
          <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
            搜索
          </Button>
          <Button icon={<ReloadOutlined />} onClick={handleReset}>
            重置
          </Button>
        </Space>
      </div>
      
      <Table
        columns={columns}
        dataSource={data}
        rowKey="id"
        loading={loading}
        pagination={{
          ...PAGINATION_CONFIG,
          current: currentPage,
          pageSize,
          total,
          onChange: (page, size) => {
            setCurrentPage(page)
            setPageSize(size)
          },
        }}
        scroll={{ x: 1400 }}
      />
    </div>
  )
}

export default User
