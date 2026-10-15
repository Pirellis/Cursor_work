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

export interface RoleRecord {
  id: number;
  name: string;
  code: string;
  description: string;
  status: number;
  userCount?: number;
  createdAt: string;
}

export interface UserRecord {
  id: number;
  username: string;
  nickname: string;
  email: string;
  phone: string;
  avatar: string;
  roleId: number;
  role: RoleRecord;
  status: number;
  lastLoginAt: string | null;
  createdAt: string;
}

export interface UserForm {
  username: string;
  password?: string;
  nickname: string;
  email: string;
  phone: string;
  roleId?: number;
  status: number;
}

export interface RoleForm {
  name: string;
  code: string;
  description: string;
  status: number;
}

export interface LogRecord {
  id: number;
  username: string;
  action: string;
  method: string;
  path: string;
  ip: string;
  success: boolean;
  latency: number;
  createdAt: string;
}

export interface Overview {
  userCount: number;
  activeCount: number;
  roleCount: number;
  logsToday: number;
  trend: { day: string; count: number }[];
  recentLogs: LogRecord[];
}

export function getOverview() {
  return axios.get<Overview>('/api/system/overview');
}

export function queryUsers(params: PageQuery) {
  return axios.get<PageResult<UserRecord>>('/api/system/users', { params });
}

export function createUser(data: UserForm) {
  return axios.post<UserRecord>('/api/system/users', data);
}

export function updateUser(id: number, data: UserForm) {
  return axios.put(`/api/system/users/${id}`, data);
}

export function resetUserPassword(id: number, password: string) {
  return axios.put(`/api/system/users/${id}/password`, { password });
}

export function deleteUser(id: number) {
  return axios.delete(`/api/system/users/${id}`);
}

export function queryRoles(params: PageQuery) {
  return axios.get<PageResult<RoleRecord>>('/api/system/roles', { params });
}

export function getRoleOptions() {
  return axios.get<RoleRecord[]>('/api/system/roles/options');
}

export function createRole(data: RoleForm) {
  return axios.post<RoleRecord>('/api/system/roles', data);
}

export function updateRole(id: number, data: RoleForm) {
  return axios.put(`/api/system/roles/${id}`, data);
}

export function deleteRole(id: number) {
  return axios.delete(`/api/system/roles/${id}`);
}

export function queryLogs(params: PageQuery) {
  return axios.get<PageResult<LogRecord>>('/api/system/logs', { params });
}
