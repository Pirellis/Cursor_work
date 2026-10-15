import axios from 'axios';

export interface PageResult<T> {
  list: T[];
  total: number;
}

export interface PageQuery {
  current: number;
  pageSize: number;
  [key: string]: unknown;
}

export interface Role {
  id: number;
  code: string;
  name: string;
  description: string;
  status: number;
  builtin: boolean;
  userCount: number;
  createdAt: string;
}

export interface SysUser {
  id: number;
  username: string;
  nickname: string;
  email: string;
  phone: string;
  avatar: string;
  job: string;
  introduction: string;
  status: number;
  roleId: number;
  role: Role;
  lastLoginAt: string | null;
  createdAt: string;
}

export type SysUserForm = Partial<Omit<SysUser, 'role'>> & {
  password?: string;
};

export interface OperationLog {
  id: number;
  createdAt: string;
  username: string;
  method: string;
  path: string;
  status: number;
  code: number;
  latencyMs: number;
  ip: string;
  userAgent: string;
}

export interface Overview {
  users: number;
  activeUsers: number;
  roles: number;
  todayOperations: number;
}

export const getOverview = () => axios.get<Overview>('/api/system/overview');

export const listUsers = (params: PageQuery) =>
  axios.get<PageResult<SysUser>>('/api/system/users', { params });
export const createUser = (data: SysUserForm) =>
  axios.post('/api/system/users', data);
export const updateUser = (id: number, data: SysUserForm) =>
  axios.put(`/api/system/users/${id}`, data);
export const deleteUser = (id: number) =>
  axios.delete(`/api/system/users/${id}`);

export const listRoles = (params?: { keyword?: string }) =>
  axios.get<Role[]>('/api/system/roles', { params });
export const createRole = (data: Partial<Role>) =>
  axios.post('/api/system/roles', data);
export const updateRole = (id: number, data: Partial<Role>) =>
  axios.put(`/api/system/roles/${id}`, data);
export const deleteRole = (id: number) =>
  axios.delete(`/api/system/roles/${id}`);

export const listLogs = (params: PageQuery) =>
  axios.get<PageResult<OperationLog>>('/api/system/logs', { params });
