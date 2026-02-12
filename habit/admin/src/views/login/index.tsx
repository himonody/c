import React, { useState } from 'react'
import { Form, Input, Button, Card, message } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { useDispatch } from 'react-redux'
import { setToken, setUserInfo } from '@/redux/modules/auth'
import { login, type AdminLoginRequest } from '@/api/auth'

const Login: React.FC = () => {
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const dispatch = useDispatch()

  const handleLogin = async (values: AdminLoginRequest) => {
    setLoading(true)
    try {
      const response = await login(values)
      const { token, adminInfo } = response.data
      
      dispatch(setToken(token))
      dispatch(setUserInfo(adminInfo))
      localStorage.setItem('token', token)
      
      message.success('登录成功')
      navigate('/dashboard')
    } catch (error: any) {
      if (error.code === 'invalid_credentials') {
        message.error('用户名或密码错误')
      } else if (error.code === 'account_disabled') {
        message.error('账户已被禁用')
      } else {
        message.error('登录失败，请稍后重试')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card title="Habit Admin 登录" style={{ width: '100%' }}>
      <Form
        name="login"
        onFinish={handleLogin}
        autoComplete="off"
        size="large"
      >
        <Form.Item
          name="username"
          rules={[{ required: true, message: '请输入用户名!' }]}
        >
          <Input
            prefix={<UserOutlined />}
            placeholder="用户名"
          />
        </Form.Item>

        <Form.Item
          name="password"
          rules={[{ required: true, message: '请输入密码!' }]}
        >
          <Input.Password
            prefix={<LockOutlined />}
            placeholder="密码"
          />
        </Form.Item>

        <Form.Item>
          <Button
            type="primary"
            htmlType="submit"
            loading={loading}
            block
          >
            登录
          </Button>
        </Form.Item>
      </Form>
    </Card>
  )
}

export default Login
