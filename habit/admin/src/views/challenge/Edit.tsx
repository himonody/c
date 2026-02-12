import React, { useState, useEffect } from 'react'
import { Form, Input, Button, Card, Switch, InputNumber, message, Space } from 'antd'
import { useNavigate, useParams } from 'react-router-dom'
import { ArrowLeftOutlined } from '@ant-design/icons'
import { createChallenge, updateChallenge, type ChallengeUpsertRequest } from '@/api/challenge'

const ChallengeEdit: React.FC = () => {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const { id } = useParams()
  const isEdit = !!id

  useEffect(() => {
    if (isEdit && id) {
      // 这里应该获取挑战详情并填充表单
      // fetchChallengeDetail(id)
    }
  }, [id, isEdit])

  const handleSubmit = async (values: ChallengeUpsertRequest) => {
    setLoading(true)
    try {
      if (isEdit && id) {
        await updateChallenge({ ...values, id: Number(id) })
        message.success('更新成功')
      } else {
        await createChallenge(values)
        message.success('创建成功')
      }
      navigate('/challenge/list')
    } catch (error) {
      message.error(isEdit ? '更新失败' : '创建失败')
    } finally {
      setLoading(false)
    }
  }

  const handleBack = () => {
    navigate('/challenge/list')
  }

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Button icon={<ArrowLeftOutlined />} onClick={handleBack}>
          返回
        </Button>
        <h2>{isEdit ? '编辑挑战' : '新增挑战'}</h2>
      </Space>

      <Card>
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          initialValues={{
            isAutoSettle: true,
            settleTime: '06:10:00',
            cycleDays: 1,
            startTime: '06:00:00',
            endTime: '06:10:00',
            maxDepositAmount: 0,
            minWithdrawAmount: 0,
            maxDailyProfit: 0,
            excessTaxRate: 98,
            minDailyProfit: 0,
            dailyPlatformSubsidy: 0,
            uncheckDeductRate: 100,
            minUncheckUsers: 2,
            commissionFollow: 0,
            commissionJoin: 0,
            commissionL1: 0,
            commissionL2: 0,
            commissionL3: 0,
          }}
        >
          <Form.Item
            label="自动结算"
            name="isAutoSettle"
            valuePropName="checked"
          >
            <Switch />
          </Form.Item>

          <Form.Item
            label="结算时间"
            name="settleTime"
            rules={[{ required: true, message: '请输入结算时间' }]}
          >
            <Input placeholder="06:10:00" />
          </Form.Item>

          <Form.Item
            label="挑战天数"
            name="cycleDays"
            rules={[{ required: true, message: '请输入挑战天数' }]}
          >
            <InputNumber min={1} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            label="打卡开始时间"
            name="startTime"
            rules={[{ required: true, message: '请输入开始时间' }]}
          >
            <Input placeholder="06:00:00" />
          </Form.Item>

          <Form.Item
            label="打卡结束时间"
            name="endTime"
            rules={[{ required: true, message: '请输入结束时间' }]}
          >
            <Input placeholder="06:10:00" />
          </Form.Item>

          <Form.Item
            label="挑战金上限"
            name="maxDepositAmount"
          >
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            label="最低提现"
            name="minWithdrawAmount"
          >
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            label="个人每日最高收益上限"
            name="maxDailyProfit"
            rules={[{ required: true, message: '请输入最高收益上限' }]}
          >
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            label="超过部分扣除比例(%)"
            name="excessTaxRate"
          >
            <InputNumber min={0} max={100} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            label="个人每日最低收益"
            name="minDailyProfit"
          >
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            label="每日平台补贴"
            name="dailyPlatformSubsidy"
            rules={[{ required: true, message: '请输入平台补贴' }]}
          >
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            label="未打卡扣除金比例(%)"
            name="uncheckDeductRate"
          >
            <InputNumber min={0} max={100} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            label="人数不足不扣除的阈值"
            name="minUncheckUsers"
          >
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item label="推荐关注佣金" name="commissionFollow">
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item label="推荐参加挑战佣金" name="commissionJoin">
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item label="一级推荐佣金" name="commissionL1">
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item label="二级推荐佣金" name="commissionL2">
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item label="三级推荐佣金" name="commissionL3">
            <InputNumber min={0} style={{ width: '100%' }} />
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

export default ChallengeEdit
