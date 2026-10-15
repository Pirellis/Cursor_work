import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const SYSTEM: AppRouteRecordRaw = {
  path: '/system',
  name: 'system',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.system',
    requiresAuth: true,
    icon: 'icon-safe',
    order: 1,
  },
  children: [
    {
      path: 'user',
      name: 'SystemUser',
      component: () => import('@/views/system/user/index.vue'),
      meta: {
        locale: 'menu.system.user',
        requiresAuth: true,
        roles: ['admin'],
      },
    },
    {
      path: 'role',
      name: 'SystemRole',
      component: () => import('@/views/system/role/index.vue'),
      meta: {
        locale: 'menu.system.role',
        requiresAuth: true,
        roles: ['admin'],
      },
    },
    {
      path: 'log',
      name: 'SystemLog',
      component: () => import('@/views/system/log/index.vue'),
      meta: {
        locale: 'menu.system.log',
        requiresAuth: true,
        roles: ['admin'],
      },
    },
  ],
};

export default SYSTEM;
