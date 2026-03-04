CREATE EXTENSION IF NOT EXISTS pgcrypto;

INSERT INTO sys_departments (id, parent_id, dept_name, sort_order, leader, phone, email, status)
VALUES
    (1, 0, '总部', 1, '系统管理员', '13800000000', 'admin@opsnexus.local', 1),
    (2, 1, '平台研发部', 10, '研发负责人', '13800000001', 'rd@opsnexus.local', 1),
    (3, 1, '运维保障部', 20, '运维负责人', '13800000002', 'ops@opsnexus.local', 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_roles (id, role_code, role_name, description)
VALUES
    (1, 'super_admin', '超级管理员', '系统最高权限角色'),
    (2, 'admin', '管理员', '平台管理员角色'),
    (3, 'viewer', '只读用户', '仅可查看数据')
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_users (id, username, password_hash, nickname, email, phone, status, department_id)
VALUES
    (1, 'admin', crypt('OpsNexus@2026', gen_salt('bf')), '超级管理员', 'admin@opsnexus.local', '13800000000', 1, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_menus (id, parent_id, menu_name, route_path, component, icon, sort_order, menu_type, permission)
VALUES
    (1, 0, '仪表盘', '/dashboard', 'views/dashboard/index.vue', 'DataBoard', 1, 'menu', 'dashboard:view'),
    (10, 0, '系统管理', '/system', '', 'Setting', 10, 'dir', ''),
    (11, 10, '用户管理', '/system/user', 'views/system/user/index.vue', 'User', 11, 'menu', 'system:user:list'),
    (12, 10, '角色管理', '/system/role', 'views/system/role/index.vue', 'UserFilled', 12, 'menu', 'system:role:list'),
    (13, 10, '菜单管理', '/system/menu', 'views/system/menu/index.vue', 'Menu', 13, 'menu', 'system:menu:list'),
    (14, 10, '部门管理', '/system/department', 'views/system/department/index.vue', 'OfficeBuilding', 14, 'menu', 'system:department:list'),

    (111, 11, '新增用户', '', '', '', 111, 'button', 'system:user:create'),
    (112, 11, '编辑用户', '', '', '', 112, 'button', 'system:user:update'),
    (113, 11, '删除用户', '', '', '', 113, 'button', 'system:user:delete'),
    (114, 11, '分配角色', '', '', '', 114, 'button', 'system:user:assign_roles'),
    (115, 11, '修改密码', '', '', '', 115, 'button', 'system:user:password'),

    (121, 12, '新增角色', '', '', '', 121, 'button', 'system:role:create'),
    (122, 12, '编辑角色', '', '', '', 122, 'button', 'system:role:update'),
    (123, 12, '删除角色', '', '', '', 123, 'button', 'system:role:delete'),
    (124, 12, '分配菜单', '', '', '', 124, 'button', 'system:role:assign_menus'),

    (131, 13, '新增菜单', '', '', '', 131, 'button', 'system:menu:create'),
    (132, 13, '编辑菜单', '', '', '', 132, 'button', 'system:menu:update'),
    (133, 13, '删除菜单', '', '', '', 133, 'button', 'system:menu:delete'),

    (141, 14, '新增部门', '', '', '', 141, 'button', 'system:department:create'),
    (142, 14, '编辑部门', '', '', '', 142, 'button', 'system:department:update'),
    (143, 14, '删除部门', '', '', '', 143, 'button', 'system:department:delete')
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_user_roles (id, user_id, role_id)
VALUES
    (1, 1, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_role_menus (id, role_id, menu_id)
SELECT ROW_NUMBER() OVER () AS id, 1 AS role_id, m.id AS menu_id
FROM sys_menus m
WHERE m.deleted_at IS NULL
ON CONFLICT (id) DO NOTHING;
