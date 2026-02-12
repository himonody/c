import React from 'react'
import { Navigate } from 'react-router-dom'
import type { RouteObject } from 'react-router-dom'

// 页面组件懒加载
const Dashboard = React.lazy(() => import('@/views/dashboard'))
const Login = React.lazy(() => import('@/views/login'))
const ChallengeList = React.lazy(() => import('@/views/challenge/List'))
const ChallengeEdit = React.lazy(() => import('@/views/challenge/Edit'))
const ConfigList = React.lazy(() => import('@/views/config/List'))
const ConfigEdit = React.lazy(() => import('@/views/config/Edit'))
const User = React.lazy(() => import('@/views/user'))
const Settings = React.lazy(() => import('@/views/settings'))
const NotFound = React.lazy(() => import('@/views/404'))

export const routes: RouteObject[] = [
  {
    path: '/',
    element: <Navigate to="/dashboard" replace />,
  },
  {
    path: '/dashboard',
    element: <Dashboard />,
  },
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '/challenge',
    children: [
      {
        index: true,
        element: <ChallengeList />,
      },
      {
        path: 'list',
        element: <ChallengeList />,
      },
      {
        path: 'edit/:id?',
        element: <ChallengeEdit />,
      },
    ],
  },
  {
    path: '/config',
    children: [
      {
        index: true,
        element: <ConfigList />,
      },
      {
        path: 'list',
        element: <ConfigList />,
      },
      {
        path: 'edit/:id?',
        element: <ConfigEdit />,
      },
    ],
  },
  {
    path: '/user',
    element: <User />,
  },
  {
    path: '/settings',
    element: <Settings />,
  },
  {
    path: '/404',
    element: <NotFound />,
  },
  {
    path: '*',
    element: <Navigate to="/404" replace />,
  },
]
