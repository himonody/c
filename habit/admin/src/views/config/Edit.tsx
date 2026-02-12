import React, { useState, useEffect } from 'react'
import { Form, Input, Button, Card, Select, message, Space, Radio } from 'antd'
import { useNavigate, useParams } from 'react-router-dom'
import { ArrowLeftOutlined } from '@ant-design/icons'
import { createConfig, updateConfig, getConfig, type CreateConfigRequest, type UpdateConfigRequest } from '@/api/config'

const { TextArea } = Input
const { Option } = Select

const ConfigEdit: React.FC = () => {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [fetchLoading, setFetchLoading] = useState(false)
  const navigate = useNavigate()
  const { id } = useParams()
  const isEdit = !!id

  useEffect(() => {
    if (isEdit && id) {
      fetchConfigDetail()
    }
  }, [id, isEdit])

  const fetchConfigDetail = async () => {
    setFetchLoading(true)
    try {
      const response = await getConfig({ id: Number(id) })
      form.setFieldsValue(response.data)
    } catch (error) {
      message.error('获取配置详情失败')
      navigate('/config/list')
    } finally {
      setFetchLoading(false)
    }
  }

  const handleSubmit = async (values: CreateConfigRequest | UpdateConfigRequest) => {
    setLoading(true)
    try {
      if (isEdit && id) {
        await updateConfig({ ...values, id: Number(id) } as UpdateConfigRequest)
        message.success('更新成功')
      } else {
        await createConfig(values)
        message.success('创建成功')
      }
      navigate('/config/list')
    } catch (error: any) {
      if (error.message?.includes('配置键已存在')) {
        message.error('配置键已存在')
      } else {
        message.error(isEdit ? '更新失败' : '创建失败')
      }
    } finally {
      setLoading(false)
    }
  }

  const handleBack = () => {
    navigate('/config/list')
  }

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Button icon={<ArrowLeftOutlined />} onClick={handleBack}>
          返回
        </Button>
        <h2>{isEdit ? '编辑配置' : '新增配置'}</h2>
      </Space>

      <Card loading={fetchLoading}>
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          initialValues={{
            configType: 'string',
            isFrontend: 'N',
          }}
        >
          <Form.Item
            label="配置名称"
            name="configName"
            rules={[{ required: true, message: '请输入配置名称' }]}
          >
            <Input placeholder="请输入配置名称" />
          </Form.Item>

          <Form.Item
            label="配置键"
            name="configKey"
            rules={[
              { required: true, message: '请输入配置键' },
              { pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/, message: '配置键只能包含字母、数字和下划线，且以字母开头' }
            ]}
          >
            <Input placeholder="请输入配置键，如: app_name" />
          </Form.Item>

          <Form.Item
            label="配置类型"
            name="configType"
            rules={[{ required: true, message: '请选择配置类型' }]}
          >
            <Select placeholder="请选择配置类型">
              <Option value="string">字符串</Option>
              <Option value="number">数字</Option>
              <Option value="boolean">布尔值</Option>
              <Option value="json">JSON</Option>
            </Select>
          </Form.Item>

          <Form.Item
            label="配置值"
            name="configValue"
            rules={[{ required: true, message: '请输入配置值' }]}
          >
            <Form.Item noStyle shouldUpdate={(prevValues, currentValues) => prevValues.configType !== currentValues.configType}>
              {({ getFieldValue }) => {
                const configType = getFieldValue('configType')
                if (configType === 'boolean') {
                  return (
                    <Radio.Group>
                      <Radio value={true}>true</Radio>
                      <Radio value={false}>false</Radio>
                    </Radio.Group>
                  )
                }
                if (configType === 'json') {
                  return <TextArea rows={6} placeholder="请输入JSON格式的配置值" />
                }
                return <Input placeholder="请输入配置值" />
              }}
            </Form.Item>
          </Form.Item>

          <Form.Item
            label="前端配置"
            name="isFrontend"
            rules={[{ required: true, message: '请选择是否为前端配置' }]}
          >
            <Radio.Group>
              <Radio value="Y">是</Radio>
              <Radio value="N">否</Radio>
            </Radio.Group>
          </Form.Item>

          <Form.Item
            label="备注"
            name="remark"
          >
            <TextArea rows={3} placeholder="请输入备注信息" />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit" loading={loading}>
                {isEdit ? '更新' : '创建'}
              </Button>
              <Button onClick={handleBack}>
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default ConfigEdit
