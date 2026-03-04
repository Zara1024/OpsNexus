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

    (20, 0, 'CMDB资产', '/cmdb', '', 'Cpu', 20, 'dir', ''),
    (21, 20, '主机管理', '/cmdb/host', 'views/cmdb/host/index.vue', 'Monitor', 21, 'menu', 'cmdb:host:list'),
    (22, 20, '分组管理', '/cmdb/group', 'views/cmdb/group/index.vue', 'Collection', 22, 'menu', 'cmdb:group:list'),
    (23, 20, '凭据管理', '/cmdb/credential', 'views/cmdb/credential/index.vue', 'Key', 23, 'menu', 'cmdb:credential:list'),
    (24, 20, 'Web终端', '/cmdb/terminal', 'views/cmdb/terminal/index.vue', 'Connection', 24, 'menu', 'cmdb:host:terminal'),
    (25, 20, 'SSH审计', '/cmdb/terminal', 'views/cmdb/terminal/index.vue', 'Document', 25, 'menu', 'cmdb:ssh_record:list'),

    (211, 21, '新增主机', '', '', '', 211, 'button', 'cmdb:host:create'),
    (212, 21, '编辑主机', '', '', '', 212, 'button', 'cmdb:host:update'),
    (213, 21, '删除主机', '', '', '', 213, 'button', 'cmdb:host:delete'),
    (214, 21, '主机详情', '', '', '', 214, 'button', 'cmdb:host:detail'),
    (215, 21, 'SSH测试', '', '', '', 215, 'button', 'cmdb:host:test'),
    (216, 21, '批量操作', '', '', '', 216, 'button', 'cmdb:host:batch'),
    (217, 21, '打开终端', '', '', '', 217, 'button', 'cmdb:host:terminal'),

    (221, 22, '新增分组', '', '', '', 221, 'button', 'cmdb:group:create'),
    (222, 22, '编辑分组', '', '', '', 222, 'button', 'cmdb:group:update'),
    (223, 22, '删除分组', '', '', '', 223, 'button', 'cmdb:group:delete'),

    (231, 23, '新增凭据', '', '', '', 231, 'button', 'cmdb:credential:create'),
    (232, 23, '编辑凭据', '', '', '', 232, 'button', 'cmdb:credential:update'),
    (233, 23, '删除凭据', '', '', '', 233, 'button', 'cmdb:credential:delete'),

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
    (143, 14, '删除部门', '', '', '', 143, 'button', 'system:department:delete'),

    (30, 0, 'K8s容器管理', '/k8s', '', 'Grid', 30, 'dir', ''),
    (31, 30, '集群管理', '/k8s/cluster', 'views/k8s/cluster/index.vue', 'Cpu', 31, 'menu', 'k8s:cluster:list'),
    (32, 30, '节点管理', '/k8s/node', 'views/k8s/node/index.vue', 'Monitor', 32, 'menu', 'k8s:node:list'),
    (33, 30, '命名空间', '/k8s/namespace', 'views/k8s/cluster/index.vue', 'Collection', 33, 'menu', 'k8s:namespace:list'),
    (34, 30, '工作负载-Deployment', '/k8s/workload/deployment', 'views/k8s/workload/deployment.vue', 'SetUp', 34, 'menu', 'k8s:workload:deployment:list'),
    (35, 30, '工作负载-Pod', '/k8s/workload/pod', 'views/k8s/workload/pod.vue', 'Memo', 35, 'menu', 'k8s:workload:pod:list'),
    (36, 30, '网络-Service', '/k8s/network/service', 'views/k8s/network/service.vue', 'Share', 36, 'menu', 'k8s:network:service:list'),
    (37, 30, '网络-Ingress', '/k8s/network/ingress', 'views/k8s/network/ingress.vue', 'Connection', 37, 'menu', 'k8s:network:ingress:list'),
    (38, 30, '配置-ConfigMap', '/k8s/config/configmap', 'views/k8s/config/configmap.vue', 'Document', 38, 'menu', 'k8s:config:configmap:list'),
    (39, 30, '配置-Secret', '/k8s/config/secret', 'views/k8s/config/secret.vue', 'Key', 39, 'menu', 'k8s:config:secret:list'),

    (311, 31, '新增集群', '', '', '', 311, 'button', 'k8s:cluster:create'),
    (312, 31, '集群详情', '', '', '', 312, 'button', 'k8s:cluster:detail'),
    (313, 31, '编辑集群', '', '', '', 313, 'button', 'k8s:cluster:update'),
    (314, 31, '删除集群', '', '', '', 314, 'button', 'k8s:cluster:delete'),
    (315, 31, '连接测试', '', '', '', 315, 'button', 'k8s:cluster:test'),

    (321, 32, '封锁节点', '', '', '', 321, 'button', 'k8s:node:cordon'),
    (322, 32, '解封节点', '', '', '', 322, 'button', 'k8s:node:uncordon'),
    (323, 32, '排水节点', '', '', '', 323, 'button', 'k8s:node:drain'),

    (341, 34, '创建Deployment', '', '', '', 341, 'button', 'k8s:workload:deployment:create'),
    (342, 34, '查看Deployment详情', '', '', '', 342, 'button', 'k8s:workload:deployment:detail'),
    (343, 34, '更新Deployment', '', '', '', 343, 'button', 'k8s:workload:deployment:update'),
    (344, 34, '删除Deployment', '', '', '', 344, 'button', 'k8s:workload:deployment:delete'),
    (345, 34, '伸缩Deployment', '', '', '', 345, 'button', 'k8s:workload:deployment:scale'),
    (346, 34, '重启Deployment', '', '', '', 346, 'button', 'k8s:workload:deployment:restart'),

    (351, 35, 'Pod详情', '', '', '', 351, 'button', 'k8s:workload:pod:detail'),
    (352, 35, 'Pod日志', '', '', '', 352, 'button', 'k8s:workload:pod:log'),
    (353, 35, '删除Pod', '', '', '', 353, 'button', 'k8s:workload:pod:delete'),

    (381, 38, '保存ConfigMap', '', '', '', 381, 'button', 'k8s:config:configmap:save'),
    (382, 38, '删除ConfigMap', '', '', '', 382, 'button', 'k8s:config:configmap:delete'),

    (391, 39, '保存Secret', '', '', '', 391, 'button', 'k8s:config:secret:save'),
    (392, 39, '删除Secret', '', '', '', 392, 'button', 'k8s:config:secret:delete'),

    (40, 0, '监控告警', '/monitor', '', 'Bell', 40, 'dir', ''),
    (41, 40, '告警规则', '/monitor/alert-rule', 'views/monitor/alert-rule/index.vue', 'BellFilled', 41, 'menu', 'monitor:rule:list'),
    (42, 40, '活跃告警', '/monitor/alert', 'views/monitor/alert/index.vue', 'Warning', 42, 'menu', 'monitor:alert:list'),
    (43, 40, '历史告警', '/monitor/alert/history', 'views/monitor/alert/history.vue', 'List', 43, 'menu', 'monitor:alert:history'),
    (44, 40, '通知渠道', '/monitor/channel', 'views/monitor/channel/index.vue', 'Message', 44, 'menu', 'monitor:channel:list'),

    (411, 41, '新增规则', '', '', '', 411, 'button', 'monitor:rule:create'),
    (412, 41, '编辑规则', '', '', '', 412, 'button', 'monitor:rule:update'),
    (413, 41, '删除规则', '', '', '', 413, 'button', 'monitor:rule:delete'),
    (414, 41, '切换规则状态', '', '', '', 414, 'button', 'monitor:rule:toggle'),

    (421, 42, '确认告警', '', '', '', 421, 'button', 'monitor:alert:ack'),
    (422, 42, '静默告警', '', '', '', 422, 'button', 'monitor:alert:silence'),

    (441, 44, '新增渠道', '', '', '', 441, 'button', 'monitor:channel:create'),
    (442, 44, '测试渠道', '', '', '', 442, 'button', 'monitor:channel:test'),

    (50, 0, '审计日志', '/audit', '', 'Document', 50, 'dir', ''),
    (51, 50, '操作日志', '/audit/operation', 'views/audit/operation/index.vue', 'Tickets', 51, 'menu', 'audit:log:list'),
    (52, 50, '登录日志', '/audit/login', 'views/audit/login/index.vue', 'User', 52, 'menu', 'audit:login:list'),

    (511, 51, '导出日志', '', '', '', 511, 'button', 'audit:log:export')
ON CONFLICT (id) DO NOTHING;
