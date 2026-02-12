import React, { useState, useEffect } from 'react'
import { Table, Button, Space, message, Modal, Input } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, SearchOutlined, ReloadOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import type { ColumnsType } from 'antd/es/table'
import { getConfigList, deleteConfig, type ConfigInfo } from '@/api/config'
import { PAGINATION_CONFIG } from '@/config'

const ConfigList: React.FC = () => {
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState<ConfigInfo[]>([])
  const [total, setTotal] = useState(0)
  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [searchParams, setSearchParams] = useState({
    configName: '',
    configKey: '',
  })
  const navigate = useNavigate()

  const fetchData = async () => {
    setLoading(true)
    try {
      const response = await getConfigList({
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

  const handleEdit = (record: ConfigInfo) => {
    navigate(`/config/edit/${record.id}`)
  }

  const handleAdd = () => {
    navigate('/config/edit')
  }

  const handleDelete = (record: ConfigInfo) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除配置 "${record.configName}" 吗？`,
      onOk: async () => {
        try {
          await deleteConfig(record.id)
          message.success('删除成功')
          fetchData()
        } catch (error) {
          message.error('删除失败')
        }
      },
    })
  }

  const handleSearch = () => {
    setCurrentPage(1)
    fetchData()
  }

  const handleReset = () => {
    setSearchParams({ configName: '', configKey: '' })
    setCurrentPage(1)
    setTimeout(fetchData, 0)
  }

  const columns: ColumnsType<ConfigInfo> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '配置名称',
      dataIndex: 'configName',
      key: 'configName',
      width: 150,
    },
    {
      title: '配置键',
      dataIndex: 'configKey',
      key: 'configKey',
      width: 200,
    },
    {
      title: '配置值',
      dataIndex: 'configValue',
      key: 'configValue',
      width: 200,
      ellipsis: true,
    },
    {
      title: '配置类型',
      dataIndex: 'configType',
      key: 'configType',
      width: 100,
      render: (value: string) => {
        const typeMap: Record<string, { text: string; color: string }> = {
          string: { text: '字符串', color: 'blue' },
          number: { text: '数字', color: 'green' },
          boolean: { text: '布尔', color: 'orange' },
          json: { text: 'JSON', color: 'purple' },
        }
        const config = typeMap[value] || { text: value, color: 'default' }
        return <span style={{ color: config.color }}>{config.text}</span>
      },
    },
    {
      title: '前端配置',
      dataIndex: 'isFrontend',
      key: 'isFrontend',
      width: 100,
      render: (value: string) => (
        <span style={{ color: value === 'Y' ? '#52c41a' : '#d9d9d9' }}>
          {value === 'Y' ? '是' : '否'}
        </span>
      ),
    },
    {
      title: '备注',
      dataIndex: 'remark',
      key: 'remark',
      width: 200,
      ellipsis: true,
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 180,
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
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2>系统配置</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          新增配置
        </Button>
      </div>

      <div style={{ marginBottom: 16, display: 'flex', gap: 16 }}>
        <Input
          placeholder="配置名称"
          value={searchParams.configName}
          onChange={(e) => setSearchParams({ ...searchParams, configName: e.target.value })}
          style={{ width: 200 }}
        />
        <Input
          placeholder="配置键"
          value={searchParams.configKey}
          onChange={(e) => setSearchParams({ ...searchParams, configKey: e.target.value })}
          style={{ width: 200 }}
        />
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

export default ConfigList
