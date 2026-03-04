import request from '@/utils/request'

export interface UserItem {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  status: number
  department_id: number
}

export interface UserListResp {
  list: UserItem[]
  total: number
}

export interface RoleItem {
  id: number
  role_code: string
  role_name: string
  description: string
}

export interface MenuItem {
  id: number
  parent_id: number
  menu_name: string
  route_path: string
  component: string
  icon: string
  sort_order: number
  menu_type: 'dir' | 'menu' | 'button'
  permission: string
  children?: MenuItem[]
}

export interface DepartmentItem {
  id: number
  parent_id: number
  dept_name: string
  sort_order: number
  leader: string
  phone: string
  email: string
  status: number
}

export const listUsersApi = (params: { page: number; page_size: number; keyword?: string }) =>
  request.get('/v1/system/users', { params }) as Promise<UserListResp>
export const createUserApi = (data: Record<string, unknown>) => request.post('/v1/system/users', data)
export const updateUserApi = (id: number, data: Record<string, unknown>) => request.put(`/v1/system/users/${id}`, data)
export const deleteUserApi = (id: number) => request.delete(`/v1/system/users/${id}`)
export const changeUserPasswordApi = (id: number, data: Record<string, unknown>) => request.put(`/v1/system/users/${id}/password`, data)
export const assignUserRolesApi = (id: number, role_ids: number[]) => request.put(`/v1/system/users/${id}/roles`, { role_ids })

export const listRolesApi = () => request.get('/v1/system/roles') as Promise<{ list: RoleItem[] }>
export const createRoleApi = (data: Record<string, unknown>) => request.post('/v1/system/roles', data)
export const updateRoleApi = (id: number, data: Record<string, unknown>) => request.put(`/v1/system/roles/${id}`, data)
export const deleteRoleApi = (id: number) => request.delete(`/v1/system/roles/${id}`)
export const assignRoleMenusApi = (id: number, menu_ids: number[]) => request.put(`/v1/system/roles/${id}/menus`, { menu_ids })

export const listMenusApi = () => request.get('/v1/system/menus') as Promise<{ list: MenuItem[] }>
export const menuTreeApi = () => request.get('/v1/system/menus/tree') as Promise<{ list: MenuItem[] }>
export const createMenuApi = (data: Record<string, unknown>) => request.post('/v1/system/menus', data)
export const updateMenuApi = (id: number, data: Record<string, unknown>) => request.put(`/v1/system/menus/${id}`, data)
export const deleteMenuApi = (id: number) => request.delete(`/v1/system/menus/${id}`)

export const listDepartmentsApi = () => request.get('/v1/system/departments') as Promise<{ list: DepartmentItem[] }>
export const createDepartmentApi = (data: Record<string, unknown>) => request.post('/v1/system/departments', data)
export const updateDepartmentApi = (id: number, data: Record<string, unknown>) => request.put(`/v1/system/departments/${id}`, data)
export const deleteDepartmentApi = (id: number) => request.delete(`/v1/system/departments/${id}`)
