import { computed, onMounted, reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { batchCMDBHostsApi, createCMDBHostApi, deleteCMDBHostApi, getCMDBHostApi, listCMDBGroupsApi, listCMDBHostsApi, testCMDBHostSSHApi, updateCMDBHostApi, } from '@/api/cmdb';
const router = useRouter();
const loading = ref(false);
const list = ref([]);
const total = ref(0);
const selectedIDs = ref([]);
const groupOptions = ref([]);
const dialogVisible = ref(false);
const detailVisible = ref(false);
const detail = ref(null);
const query = reactive({ page: 1, page_size: 10, keyword: '', group_id: undefined, status: undefined });
const form = reactive({
    id: 0, hostname: '', inner_ip: '', outer_ip: '', os_type: '', os_version: '', arch: '',
    cpu_cores: 0, memory_gb: 0, disk_gb: 0, group_id: 0, credential_id: 0,
    cloud_provider: '', instance_id: '', region: '', agent_status: 'unknown', labels: '{}', status: 1,
});
const statusSwitch = computed({
    get: () => form.status === 1,
    set: (v) => { form.status = v ? 1 : 2; },
});
const flattenGroups = (nodes) => {
    const result = [];
    const walk = (items) => {
        items.forEach((i) => {
            result.push(i);
            if (i.children?.length)
                walk(i.children);
        });
    };
    walk(nodes);
    return result;
};
const loadGroups = async () => {
    const res = await listCMDBGroupsApi();
    groupOptions.value = flattenGroups(res.list || []);
};
const loadData = async () => {
    loading.value = true;
    try {
        const res = await listCMDBHostsApi(query);
        list.value = res.list || [];
        total.value = res.total || 0;
    }
    finally {
        loading.value = false;
    }
};
const handleSelectionChange = (rows) => {
    selectedIDs.value = rows.map((r) => r.id);
};
const handlePageChange = (page) => {
    query.page = page;
    loadData();
};
const openCreate = () => {
    Object.assign(form, {
        id: 0, hostname: '', inner_ip: '', outer_ip: '', os_type: '', os_version: '', arch: '',
        cpu_cores: 0, memory_gb: 0, disk_gb: 0, group_id: 0, credential_id: 0,
        cloud_provider: '', instance_id: '', region: '', agent_status: 'unknown', labels: '{}', status: 1,
    });
    dialogVisible.value = true;
};
const openEdit = (row) => {
    Object.assign(form, row);
    dialogVisible.value = true;
};
const openDetail = async (id) => {
    detail.value = await getCMDBHostApi(id);
    detailVisible.value = true;
};
const submit = async () => {
    if (form.id) {
        await updateCMDBHostApi(form.id, form);
    }
    else {
        await createCMDBHostApi(form);
    }
    ElMessage.success('保存成功');
    dialogVisible.value = false;
    await loadData();
};
const handleDelete = async (id) => {
    await ElMessageBox.confirm('确认删除该主机吗？', '提示', { type: 'warning' });
    await deleteCMDBHostApi(id);
    ElMessage.success('删除成功');
    await loadData();
};
const handleBatch = async (action) => {
    if (!selectedIDs.value.length) {
        ElMessage.warning('请先选择主机');
        return;
    }
    await batchCMDBHostsApi({ action, ids: selectedIDs.value });
    ElMessage.success('批量操作成功');
    await loadData();
};
const handleTest = async (id) => {
    await testCMDBHostSSHApi(id);
    ElMessage.success('SSH连接成功');
};
const openTerminal = (id) => {
    router.push(`/cmdb/terminal?host_id=${id}`);
};
onMounted(async () => {
    await loadGroups();
    await loadData();
});
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
// CSS variable injection 
// CSS variable injection end 
const __VLS_0 = {}.ElCard;
/** @type {[typeof __VLS_components.ElCard, typeof __VLS_components.elCard, typeof __VLS_components.ElCard, typeof __VLS_components.elCard, ]} */ ;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({}));
const __VLS_2 = __VLS_1({}, ...__VLS_functionalComponentArgsRest(__VLS_1));
var __VLS_4 = {};
__VLS_3.slots.default;
{
    const { header: __VLS_thisSlot } = __VLS_3.slots;
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ class: "toolbar" },
    });
    __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ class: "actions" },
    });
    const __VLS_5 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_6 = __VLS_asFunctionalComponent(__VLS_5, new __VLS_5({
        ...{ 'onClick': {} },
        type: "primary",
    }));
    const __VLS_7 = __VLS_6({
        ...{ 'onClick': {} },
        type: "primary",
    }, ...__VLS_functionalComponentArgsRest(__VLS_6));
    let __VLS_9;
    let __VLS_10;
    let __VLS_11;
    const __VLS_12 = {
        onClick: (__VLS_ctx.openCreate)
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('cmdb:host:create') }, null, null);
    __VLS_8.slots.default;
    var __VLS_8;
    const __VLS_13 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_14 = __VLS_asFunctionalComponent(__VLS_13, new __VLS_13({
        ...{ 'onClick': {} },
        type: "warning",
        plain: true,
    }));
    const __VLS_15 = __VLS_14({
        ...{ 'onClick': {} },
        type: "warning",
        plain: true,
    }, ...__VLS_functionalComponentArgsRest(__VLS_14));
    let __VLS_17;
    let __VLS_18;
    let __VLS_19;
    const __VLS_20 = {
        onClick: (...[$event]) => {
            __VLS_ctx.handleBatch('enable');
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('cmdb:host:batch') }, null, null);
    __VLS_16.slots.default;
    var __VLS_16;
    const __VLS_21 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_22 = __VLS_asFunctionalComponent(__VLS_21, new __VLS_21({
        ...{ 'onClick': {} },
        type: "info",
        plain: true,
    }));
    const __VLS_23 = __VLS_22({
        ...{ 'onClick': {} },
        type: "info",
        plain: true,
    }, ...__VLS_functionalComponentArgsRest(__VLS_22));
    let __VLS_25;
    let __VLS_26;
    let __VLS_27;
    const __VLS_28 = {
        onClick: (...[$event]) => {
            __VLS_ctx.handleBatch('disable');
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('cmdb:host:batch') }, null, null);
    __VLS_24.slots.default;
    var __VLS_24;
    const __VLS_29 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_30 = __VLS_asFunctionalComponent(__VLS_29, new __VLS_29({
        ...{ 'onClick': {} },
        type: "danger",
        plain: true,
    }));
    const __VLS_31 = __VLS_30({
        ...{ 'onClick': {} },
        type: "danger",
        plain: true,
    }, ...__VLS_functionalComponentArgsRest(__VLS_30));
    let __VLS_33;
    let __VLS_34;
    let __VLS_35;
    const __VLS_36 = {
        onClick: (...[$event]) => {
            __VLS_ctx.handleBatch('delete');
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('cmdb:host:batch') }, null, null);
    __VLS_32.slots.default;
    var __VLS_32;
}
const __VLS_37 = {}.ElForm;
/** @type {[typeof __VLS_components.ElForm, typeof __VLS_components.elForm, typeof __VLS_components.ElForm, typeof __VLS_components.elForm, ]} */ ;
// @ts-ignore
const __VLS_38 = __VLS_asFunctionalComponent(__VLS_37, new __VLS_37({
    inline: (true),
    ...{ class: "search" },
}));
const __VLS_39 = __VLS_38({
    inline: (true),
    ...{ class: "search" },
}, ...__VLS_functionalComponentArgsRest(__VLS_38));
__VLS_40.slots.default;
const __VLS_41 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_42 = __VLS_asFunctionalComponent(__VLS_41, new __VLS_41({}));
const __VLS_43 = __VLS_42({}, ...__VLS_functionalComponentArgsRest(__VLS_42));
__VLS_44.slots.default;
const __VLS_45 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_46 = __VLS_asFunctionalComponent(__VLS_45, new __VLS_45({
    modelValue: (__VLS_ctx.query.keyword),
    placeholder: "主机名/IP",
    clearable: true,
}));
const __VLS_47 = __VLS_46({
    modelValue: (__VLS_ctx.query.keyword),
    placeholder: "主机名/IP",
    clearable: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_46));
var __VLS_44;
const __VLS_49 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_50 = __VLS_asFunctionalComponent(__VLS_49, new __VLS_49({}));
const __VLS_51 = __VLS_50({}, ...__VLS_functionalComponentArgsRest(__VLS_50));
__VLS_52.slots.default;
const __VLS_53 = {}.ElSelect;
/** @type {[typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, ]} */ ;
// @ts-ignore
const __VLS_54 = __VLS_asFunctionalComponent(__VLS_53, new __VLS_53({
    modelValue: (__VLS_ctx.query.group_id),
    placeholder: "分组",
    clearable: true,
    ...{ style: {} },
}));
const __VLS_55 = __VLS_54({
    modelValue: (__VLS_ctx.query.group_id),
    placeholder: "分组",
    clearable: true,
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_54));
__VLS_56.slots.default;
for (const [g] of __VLS_getVForSourceType((__VLS_ctx.groupOptions))) {
    const __VLS_57 = {}.ElOption;
    /** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
    // @ts-ignore
    const __VLS_58 = __VLS_asFunctionalComponent(__VLS_57, new __VLS_57({
        key: (g.id),
        label: (g.group_name),
        value: (g.id),
    }));
    const __VLS_59 = __VLS_58({
        key: (g.id),
        label: (g.group_name),
        value: (g.id),
    }, ...__VLS_functionalComponentArgsRest(__VLS_58));
}
var __VLS_56;
var __VLS_52;
const __VLS_61 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_62 = __VLS_asFunctionalComponent(__VLS_61, new __VLS_61({}));
const __VLS_63 = __VLS_62({}, ...__VLS_functionalComponentArgsRest(__VLS_62));
__VLS_64.slots.default;
const __VLS_65 = {}.ElSelect;
/** @type {[typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, ]} */ ;
// @ts-ignore
const __VLS_66 = __VLS_asFunctionalComponent(__VLS_65, new __VLS_65({
    modelValue: (__VLS_ctx.query.status),
    placeholder: "状态",
    clearable: true,
    ...{ style: {} },
}));
const __VLS_67 = __VLS_66({
    modelValue: (__VLS_ctx.query.status),
    placeholder: "状态",
    clearable: true,
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_66));
__VLS_68.slots.default;
const __VLS_69 = {}.ElOption;
/** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
// @ts-ignore
const __VLS_70 = __VLS_asFunctionalComponent(__VLS_69, new __VLS_69({
    label: "启用",
    value: (1),
}));
const __VLS_71 = __VLS_70({
    label: "启用",
    value: (1),
}, ...__VLS_functionalComponentArgsRest(__VLS_70));
const __VLS_73 = {}.ElOption;
/** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
// @ts-ignore
const __VLS_74 = __VLS_asFunctionalComponent(__VLS_73, new __VLS_73({
    label: "禁用",
    value: (2),
}));
const __VLS_75 = __VLS_74({
    label: "禁用",
    value: (2),
}, ...__VLS_functionalComponentArgsRest(__VLS_74));
var __VLS_68;
var __VLS_64;
const __VLS_77 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_78 = __VLS_asFunctionalComponent(__VLS_77, new __VLS_77({}));
const __VLS_79 = __VLS_78({}, ...__VLS_functionalComponentArgsRest(__VLS_78));
__VLS_80.slots.default;
const __VLS_81 = {}.ElButton;
/** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
// @ts-ignore
const __VLS_82 = __VLS_asFunctionalComponent(__VLS_81, new __VLS_81({
    ...{ 'onClick': {} },
    type: "primary",
}));
const __VLS_83 = __VLS_82({
    ...{ 'onClick': {} },
    type: "primary",
}, ...__VLS_functionalComponentArgsRest(__VLS_82));
let __VLS_85;
let __VLS_86;
let __VLS_87;
const __VLS_88 = {
    onClick: (__VLS_ctx.loadData)
};
__VLS_84.slots.default;
var __VLS_84;
var __VLS_80;
var __VLS_40;
const __VLS_89 = {}.ElTable;
/** @type {[typeof __VLS_components.ElTable, typeof __VLS_components.elTable, typeof __VLS_components.ElTable, typeof __VLS_components.elTable, ]} */ ;
// @ts-ignore
const __VLS_90 = __VLS_asFunctionalComponent(__VLS_89, new __VLS_89({
    ...{ 'onSelectionChange': {} },
    data: (__VLS_ctx.list),
    rowKey: "id",
}));
const __VLS_91 = __VLS_90({
    ...{ 'onSelectionChange': {} },
    data: (__VLS_ctx.list),
    rowKey: "id",
}, ...__VLS_functionalComponentArgsRest(__VLS_90));
let __VLS_93;
let __VLS_94;
let __VLS_95;
const __VLS_96 = {
    onSelectionChange: (__VLS_ctx.handleSelectionChange)
};
__VLS_asFunctionalDirective(__VLS_directives.vLoading)(null, { ...__VLS_directiveBindingRestFields, value: (__VLS_ctx.loading) }, null, null);
__VLS_92.slots.default;
const __VLS_97 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_98 = __VLS_asFunctionalComponent(__VLS_97, new __VLS_97({
    type: "selection",
    width: "50",
}));
const __VLS_99 = __VLS_98({
    type: "selection",
    width: "50",
}, ...__VLS_functionalComponentArgsRest(__VLS_98));
const __VLS_101 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_102 = __VLS_asFunctionalComponent(__VLS_101, new __VLS_101({
    prop: "hostname",
    label: "主机名",
    minWidth: "130",
}));
const __VLS_103 = __VLS_102({
    prop: "hostname",
    label: "主机名",
    minWidth: "130",
}, ...__VLS_functionalComponentArgsRest(__VLS_102));
const __VLS_105 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_106 = __VLS_asFunctionalComponent(__VLS_105, new __VLS_105({
    prop: "inner_ip",
    label: "内网IP",
    minWidth: "130",
}));
const __VLS_107 = __VLS_106({
    prop: "inner_ip",
    label: "内网IP",
    minWidth: "130",
}, ...__VLS_functionalComponentArgsRest(__VLS_106));
const __VLS_109 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_110 = __VLS_asFunctionalComponent(__VLS_109, new __VLS_109({
    prop: "outer_ip",
    label: "公网IP",
    minWidth: "130",
}));
const __VLS_111 = __VLS_110({
    prop: "outer_ip",
    label: "公网IP",
    minWidth: "130",
}, ...__VLS_functionalComponentArgsRest(__VLS_110));
const __VLS_113 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_114 = __VLS_asFunctionalComponent(__VLS_113, new __VLS_113({
    prop: "os_type",
    label: "系统",
    minWidth: "100",
}));
const __VLS_115 = __VLS_114({
    prop: "os_type",
    label: "系统",
    minWidth: "100",
}, ...__VLS_functionalComponentArgsRest(__VLS_114));
const __VLS_117 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_118 = __VLS_asFunctionalComponent(__VLS_117, new __VLS_117({
    prop: "cpu_cores",
    label: "CPU",
    width: "80",
}));
const __VLS_119 = __VLS_118({
    prop: "cpu_cores",
    label: "CPU",
    width: "80",
}, ...__VLS_functionalComponentArgsRest(__VLS_118));
const __VLS_121 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_122 = __VLS_asFunctionalComponent(__VLS_121, new __VLS_121({
    prop: "memory_gb",
    label: "内存(GB)",
    width: "100",
}));
const __VLS_123 = __VLS_122({
    prop: "memory_gb",
    label: "内存(GB)",
    width: "100",
}, ...__VLS_functionalComponentArgsRest(__VLS_122));
const __VLS_125 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_126 = __VLS_asFunctionalComponent(__VLS_125, new __VLS_125({
    prop: "status",
    label: "状态",
    width: "80",
}));
const __VLS_127 = __VLS_126({
    prop: "status",
    label: "状态",
    width: "80",
}, ...__VLS_functionalComponentArgsRest(__VLS_126));
__VLS_128.slots.default;
{
    const { default: __VLS_thisSlot } = __VLS_128.slots;
    const [scope] = __VLS_getSlotParams(__VLS_thisSlot);
    (scope.row.status === 1 ? '启用' : '禁用');
}
var __VLS_128;
const __VLS_129 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_130 = __VLS_asFunctionalComponent(__VLS_129, new __VLS_129({
    label: "操作",
    minWidth: "330",
    fixed: "right",
}));
const __VLS_131 = __VLS_130({
    label: "操作",
    minWidth: "330",
    fixed: "right",
}, ...__VLS_functionalComponentArgsRest(__VLS_130));
__VLS_132.slots.default;
{
    const { default: __VLS_thisSlot } = __VLS_132.slots;
    const [scope] = __VLS_getSlotParams(__VLS_thisSlot);
    const __VLS_133 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_134 = __VLS_asFunctionalComponent(__VLS_133, new __VLS_133({
        ...{ 'onClick': {} },
        link: true,
        type: "primary",
    }));
    const __VLS_135 = __VLS_134({
        ...{ 'onClick': {} },
        link: true,
        type: "primary",
    }, ...__VLS_functionalComponentArgsRest(__VLS_134));
    let __VLS_137;
    let __VLS_138;
    let __VLS_139;
    const __VLS_140 = {
        onClick: (...[$event]) => {
            __VLS_ctx.openDetail(scope.row.id);
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('cmdb:host:detail') }, null, null);
    __VLS_136.slots.default;
    var __VLS_136;
    const __VLS_141 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_142 = __VLS_asFunctionalComponent(__VLS_141, new __VLS_141({
        ...{ 'onClick': {} },
        link: true,
        type: "primary",
    }));
    const __VLS_143 = __VLS_142({
        ...{ 'onClick': {} },
        link: true,
        type: "primary",
    }, ...__VLS_functionalComponentArgsRest(__VLS_142));
    let __VLS_145;
    let __VLS_146;
    let __VLS_147;
    const __VLS_148 = {
        onClick: (...[$event]) => {
            __VLS_ctx.openEdit(scope.row);
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('cmdb:host:update') }, null, null);
    __VLS_144.slots.default;
    var __VLS_144;
    const __VLS_149 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_150 = __VLS_asFunctionalComponent(__VLS_149, new __VLS_149({
        ...{ 'onClick': {} },
        link: true,
        type: "success",
    }));
    const __VLS_151 = __VLS_150({
        ...{ 'onClick': {} },
        link: true,
        type: "success",
    }, ...__VLS_functionalComponentArgsRest(__VLS_150));
    let __VLS_153;
    let __VLS_154;
    let __VLS_155;
    const __VLS_156 = {
        onClick: (...[$event]) => {
            __VLS_ctx.handleTest(scope.row.id);
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('cmdb:host:test') }, null, null);
    __VLS_152.slots.default;
    var __VLS_152;
    const __VLS_157 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_158 = __VLS_asFunctionalComponent(__VLS_157, new __VLS_157({
        ...{ 'onClick': {} },
        link: true,
        type: "warning",
    }));
    const __VLS_159 = __VLS_158({
        ...{ 'onClick': {} },
        link: true,
        type: "warning",
    }, ...__VLS_functionalComponentArgsRest(__VLS_158));
    let __VLS_161;
    let __VLS_162;
    let __VLS_163;
    const __VLS_164 = {
        onClick: (...[$event]) => {
            __VLS_ctx.openTerminal(scope.row.id);
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('cmdb:host:terminal') }, null, null);
    __VLS_160.slots.default;
    var __VLS_160;
    const __VLS_165 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_166 = __VLS_asFunctionalComponent(__VLS_165, new __VLS_165({
        ...{ 'onClick': {} },
        link: true,
        type: "danger",
    }));
    const __VLS_167 = __VLS_166({
        ...{ 'onClick': {} },
        link: true,
        type: "danger",
    }, ...__VLS_functionalComponentArgsRest(__VLS_166));
    let __VLS_169;
    let __VLS_170;
    let __VLS_171;
    const __VLS_172 = {
        onClick: (...[$event]) => {
            __VLS_ctx.handleDelete(scope.row.id);
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('cmdb:host:delete') }, null, null);
    __VLS_168.slots.default;
    var __VLS_168;
}
var __VLS_132;
var __VLS_92;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "pager" },
});
const __VLS_173 = {}.ElPagination;
/** @type {[typeof __VLS_components.ElPagination, typeof __VLS_components.elPagination, ]} */ ;
// @ts-ignore
const __VLS_174 = __VLS_asFunctionalComponent(__VLS_173, new __VLS_173({
    ...{ 'onCurrentChange': {} },
    background: true,
    layout: "total, prev, pager, next",
    total: (__VLS_ctx.total),
    currentPage: (__VLS_ctx.query.page),
    pageSize: (__VLS_ctx.query.page_size),
}));
const __VLS_175 = __VLS_174({
    ...{ 'onCurrentChange': {} },
    background: true,
    layout: "total, prev, pager, next",
    total: (__VLS_ctx.total),
    currentPage: (__VLS_ctx.query.page),
    pageSize: (__VLS_ctx.query.page_size),
}, ...__VLS_functionalComponentArgsRest(__VLS_174));
let __VLS_177;
let __VLS_178;
let __VLS_179;
const __VLS_180 = {
    onCurrentChange: (__VLS_ctx.handlePageChange)
};
var __VLS_176;
const __VLS_181 = {}.ElDialog;
/** @type {[typeof __VLS_components.ElDialog, typeof __VLS_components.elDialog, typeof __VLS_components.ElDialog, typeof __VLS_components.elDialog, ]} */ ;
// @ts-ignore
const __VLS_182 = __VLS_asFunctionalComponent(__VLS_181, new __VLS_181({
    modelValue: (__VLS_ctx.dialogVisible),
    title: (__VLS_ctx.form.id ? '编辑主机' : '新增主机'),
    width: "760px",
}));
const __VLS_183 = __VLS_182({
    modelValue: (__VLS_ctx.dialogVisible),
    title: (__VLS_ctx.form.id ? '编辑主机' : '新增主机'),
    width: "760px",
}, ...__VLS_functionalComponentArgsRest(__VLS_182));
__VLS_184.slots.default;
const __VLS_185 = {}.ElForm;
/** @type {[typeof __VLS_components.ElForm, typeof __VLS_components.elForm, typeof __VLS_components.ElForm, typeof __VLS_components.elForm, ]} */ ;
// @ts-ignore
const __VLS_186 = __VLS_asFunctionalComponent(__VLS_185, new __VLS_185({
    model: (__VLS_ctx.form),
    labelWidth: "100px",
}));
const __VLS_187 = __VLS_186({
    model: (__VLS_ctx.form),
    labelWidth: "100px",
}, ...__VLS_functionalComponentArgsRest(__VLS_186));
__VLS_188.slots.default;
const __VLS_189 = {}.ElRow;
/** @type {[typeof __VLS_components.ElRow, typeof __VLS_components.elRow, typeof __VLS_components.ElRow, typeof __VLS_components.elRow, ]} */ ;
// @ts-ignore
const __VLS_190 = __VLS_asFunctionalComponent(__VLS_189, new __VLS_189({
    gutter: (12),
}));
const __VLS_191 = __VLS_190({
    gutter: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_190));
__VLS_192.slots.default;
const __VLS_193 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_194 = __VLS_asFunctionalComponent(__VLS_193, new __VLS_193({
    span: (12),
}));
const __VLS_195 = __VLS_194({
    span: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_194));
__VLS_196.slots.default;
const __VLS_197 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_198 = __VLS_asFunctionalComponent(__VLS_197, new __VLS_197({
    label: "主机名",
}));
const __VLS_199 = __VLS_198({
    label: "主机名",
}, ...__VLS_functionalComponentArgsRest(__VLS_198));
__VLS_200.slots.default;
const __VLS_201 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_202 = __VLS_asFunctionalComponent(__VLS_201, new __VLS_201({
    modelValue: (__VLS_ctx.form.hostname),
}));
const __VLS_203 = __VLS_202({
    modelValue: (__VLS_ctx.form.hostname),
}, ...__VLS_functionalComponentArgsRest(__VLS_202));
var __VLS_200;
var __VLS_196;
const __VLS_205 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_206 = __VLS_asFunctionalComponent(__VLS_205, new __VLS_205({
    span: (12),
}));
const __VLS_207 = __VLS_206({
    span: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_206));
__VLS_208.slots.default;
const __VLS_209 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_210 = __VLS_asFunctionalComponent(__VLS_209, new __VLS_209({
    label: "内网IP",
}));
const __VLS_211 = __VLS_210({
    label: "内网IP",
}, ...__VLS_functionalComponentArgsRest(__VLS_210));
__VLS_212.slots.default;
const __VLS_213 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_214 = __VLS_asFunctionalComponent(__VLS_213, new __VLS_213({
    modelValue: (__VLS_ctx.form.inner_ip),
}));
const __VLS_215 = __VLS_214({
    modelValue: (__VLS_ctx.form.inner_ip),
}, ...__VLS_functionalComponentArgsRest(__VLS_214));
var __VLS_212;
var __VLS_208;
const __VLS_217 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_218 = __VLS_asFunctionalComponent(__VLS_217, new __VLS_217({
    span: (12),
}));
const __VLS_219 = __VLS_218({
    span: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_218));
__VLS_220.slots.default;
const __VLS_221 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_222 = __VLS_asFunctionalComponent(__VLS_221, new __VLS_221({
    label: "公网IP",
}));
const __VLS_223 = __VLS_222({
    label: "公网IP",
}, ...__VLS_functionalComponentArgsRest(__VLS_222));
__VLS_224.slots.default;
const __VLS_225 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_226 = __VLS_asFunctionalComponent(__VLS_225, new __VLS_225({
    modelValue: (__VLS_ctx.form.outer_ip),
}));
const __VLS_227 = __VLS_226({
    modelValue: (__VLS_ctx.form.outer_ip),
}, ...__VLS_functionalComponentArgsRest(__VLS_226));
var __VLS_224;
var __VLS_220;
const __VLS_229 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_230 = __VLS_asFunctionalComponent(__VLS_229, new __VLS_229({
    span: (12),
}));
const __VLS_231 = __VLS_230({
    span: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_230));
__VLS_232.slots.default;
const __VLS_233 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_234 = __VLS_asFunctionalComponent(__VLS_233, new __VLS_233({
    label: "系统类型",
}));
const __VLS_235 = __VLS_234({
    label: "系统类型",
}, ...__VLS_functionalComponentArgsRest(__VLS_234));
__VLS_236.slots.default;
const __VLS_237 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_238 = __VLS_asFunctionalComponent(__VLS_237, new __VLS_237({
    modelValue: (__VLS_ctx.form.os_type),
}));
const __VLS_239 = __VLS_238({
    modelValue: (__VLS_ctx.form.os_type),
}, ...__VLS_functionalComponentArgsRest(__VLS_238));
var __VLS_236;
var __VLS_232;
const __VLS_241 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_242 = __VLS_asFunctionalComponent(__VLS_241, new __VLS_241({
    span: (12),
}));
const __VLS_243 = __VLS_242({
    span: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_242));
__VLS_244.slots.default;
const __VLS_245 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_246 = __VLS_asFunctionalComponent(__VLS_245, new __VLS_245({
    label: "系统版本",
}));
const __VLS_247 = __VLS_246({
    label: "系统版本",
}, ...__VLS_functionalComponentArgsRest(__VLS_246));
__VLS_248.slots.default;
const __VLS_249 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_250 = __VLS_asFunctionalComponent(__VLS_249, new __VLS_249({
    modelValue: (__VLS_ctx.form.os_version),
}));
const __VLS_251 = __VLS_250({
    modelValue: (__VLS_ctx.form.os_version),
}, ...__VLS_functionalComponentArgsRest(__VLS_250));
var __VLS_248;
var __VLS_244;
const __VLS_253 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_254 = __VLS_asFunctionalComponent(__VLS_253, new __VLS_253({
    span: (12),
}));
const __VLS_255 = __VLS_254({
    span: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_254));
__VLS_256.slots.default;
const __VLS_257 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_258 = __VLS_asFunctionalComponent(__VLS_257, new __VLS_257({
    label: "架构",
}));
const __VLS_259 = __VLS_258({
    label: "架构",
}, ...__VLS_functionalComponentArgsRest(__VLS_258));
__VLS_260.slots.default;
const __VLS_261 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_262 = __VLS_asFunctionalComponent(__VLS_261, new __VLS_261({
    modelValue: (__VLS_ctx.form.arch),
}));
const __VLS_263 = __VLS_262({
    modelValue: (__VLS_ctx.form.arch),
}, ...__VLS_functionalComponentArgsRest(__VLS_262));
var __VLS_260;
var __VLS_256;
const __VLS_265 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_266 = __VLS_asFunctionalComponent(__VLS_265, new __VLS_265({
    span: (12),
}));
const __VLS_267 = __VLS_266({
    span: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_266));
__VLS_268.slots.default;
const __VLS_269 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_270 = __VLS_asFunctionalComponent(__VLS_269, new __VLS_269({
    label: "CPU核数",
}));
const __VLS_271 = __VLS_270({
    label: "CPU核数",
}, ...__VLS_functionalComponentArgsRest(__VLS_270));
__VLS_272.slots.default;
const __VLS_273 = {}.ElInputNumber;
/** @type {[typeof __VLS_components.ElInputNumber, typeof __VLS_components.elInputNumber, ]} */ ;
// @ts-ignore
const __VLS_274 = __VLS_asFunctionalComponent(__VLS_273, new __VLS_273({
    modelValue: (__VLS_ctx.form.cpu_cores),
    min: (0),
}));
const __VLS_275 = __VLS_274({
    modelValue: (__VLS_ctx.form.cpu_cores),
    min: (0),
}, ...__VLS_functionalComponentArgsRest(__VLS_274));
var __VLS_272;
var __VLS_268;
const __VLS_277 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_278 = __VLS_asFunctionalComponent(__VLS_277, new __VLS_277({
    span: (12),
}));
const __VLS_279 = __VLS_278({
    span: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_278));
__VLS_280.slots.default;
const __VLS_281 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_282 = __VLS_asFunctionalComponent(__VLS_281, new __VLS_281({
    label: "内存GB",
}));
const __VLS_283 = __VLS_282({
    label: "内存GB",
}, ...__VLS_functionalComponentArgsRest(__VLS_282));
__VLS_284.slots.default;
const __VLS_285 = {}.ElInputNumber;
/** @type {[typeof __VLS_components.ElInputNumber, typeof __VLS_components.elInputNumber, ]} */ ;
// @ts-ignore
const __VLS_286 = __VLS_asFunctionalComponent(__VLS_285, new __VLS_285({
    modelValue: (__VLS_ctx.form.memory_gb),
    min: (0),
    precision: (2),
}));
const __VLS_287 = __VLS_286({
    modelValue: (__VLS_ctx.form.memory_gb),
    min: (0),
    precision: (2),
}, ...__VLS_functionalComponentArgsRest(__VLS_286));
var __VLS_284;
var __VLS_280;
const __VLS_289 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_290 = __VLS_asFunctionalComponent(__VLS_289, new __VLS_289({
    span: (12),
}));
const __VLS_291 = __VLS_290({
    span: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_290));
__VLS_292.slots.default;
const __VLS_293 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_294 = __VLS_asFunctionalComponent(__VLS_293, new __VLS_293({
    label: "磁盘GB",
}));
const __VLS_295 = __VLS_294({
    label: "磁盘GB",
}, ...__VLS_functionalComponentArgsRest(__VLS_294));
__VLS_296.slots.default;
const __VLS_297 = {}.ElInputNumber;
/** @type {[typeof __VLS_components.ElInputNumber, typeof __VLS_components.elInputNumber, ]} */ ;
// @ts-ignore
const __VLS_298 = __VLS_asFunctionalComponent(__VLS_297, new __VLS_297({
    modelValue: (__VLS_ctx.form.disk_gb),
    min: (0),
    precision: (2),
}));
const __VLS_299 = __VLS_298({
    modelValue: (__VLS_ctx.form.disk_gb),
    min: (0),
    precision: (2),
}, ...__VLS_functionalComponentArgsRest(__VLS_298));
var __VLS_296;
var __VLS_292;
const __VLS_301 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_302 = __VLS_asFunctionalComponent(__VLS_301, new __VLS_301({
    span: (12),
}));
const __VLS_303 = __VLS_302({
    span: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_302));
__VLS_304.slots.default;
const __VLS_305 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_306 = __VLS_asFunctionalComponent(__VLS_305, new __VLS_305({
    label: "分组ID",
}));
const __VLS_307 = __VLS_306({
    label: "分组ID",
}, ...__VLS_functionalComponentArgsRest(__VLS_306));
__VLS_308.slots.default;
const __VLS_309 = {}.ElInputNumber;
/** @type {[typeof __VLS_components.ElInputNumber, typeof __VLS_components.elInputNumber, ]} */ ;
// @ts-ignore
const __VLS_310 = __VLS_asFunctionalComponent(__VLS_309, new __VLS_309({
    modelValue: (__VLS_ctx.form.group_id),
    min: (0),
}));
const __VLS_311 = __VLS_310({
    modelValue: (__VLS_ctx.form.group_id),
    min: (0),
}, ...__VLS_functionalComponentArgsRest(__VLS_310));
var __VLS_308;
var __VLS_304;
const __VLS_313 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_314 = __VLS_asFunctionalComponent(__VLS_313, new __VLS_313({
    span: (12),
}));
const __VLS_315 = __VLS_314({
    span: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_314));
__VLS_316.slots.default;
const __VLS_317 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_318 = __VLS_asFunctionalComponent(__VLS_317, new __VLS_317({
    label: "凭据ID",
}));
const __VLS_319 = __VLS_318({
    label: "凭据ID",
}, ...__VLS_functionalComponentArgsRest(__VLS_318));
__VLS_320.slots.default;
const __VLS_321 = {}.ElInputNumber;
/** @type {[typeof __VLS_components.ElInputNumber, typeof __VLS_components.elInputNumber, ]} */ ;
// @ts-ignore
const __VLS_322 = __VLS_asFunctionalComponent(__VLS_321, new __VLS_321({
    modelValue: (__VLS_ctx.form.credential_id),
    min: (0),
}));
const __VLS_323 = __VLS_322({
    modelValue: (__VLS_ctx.form.credential_id),
    min: (0),
}, ...__VLS_functionalComponentArgsRest(__VLS_322));
var __VLS_320;
var __VLS_316;
const __VLS_325 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_326 = __VLS_asFunctionalComponent(__VLS_325, new __VLS_325({
    span: (12),
}));
const __VLS_327 = __VLS_326({
    span: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_326));
__VLS_328.slots.default;
const __VLS_329 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_330 = __VLS_asFunctionalComponent(__VLS_329, new __VLS_329({
    label: "云厂商",
}));
const __VLS_331 = __VLS_330({
    label: "云厂商",
}, ...__VLS_functionalComponentArgsRest(__VLS_330));
__VLS_332.slots.default;
const __VLS_333 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_334 = __VLS_asFunctionalComponent(__VLS_333, new __VLS_333({
    modelValue: (__VLS_ctx.form.cloud_provider),
}));
const __VLS_335 = __VLS_334({
    modelValue: (__VLS_ctx.form.cloud_provider),
}, ...__VLS_functionalComponentArgsRest(__VLS_334));
var __VLS_332;
var __VLS_328;
const __VLS_337 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_338 = __VLS_asFunctionalComponent(__VLS_337, new __VLS_337({
    span: (12),
}));
const __VLS_339 = __VLS_338({
    span: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_338));
__VLS_340.slots.default;
const __VLS_341 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_342 = __VLS_asFunctionalComponent(__VLS_341, new __VLS_341({
    label: "实例ID",
}));
const __VLS_343 = __VLS_342({
    label: "实例ID",
}, ...__VLS_functionalComponentArgsRest(__VLS_342));
__VLS_344.slots.default;
const __VLS_345 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_346 = __VLS_asFunctionalComponent(__VLS_345, new __VLS_345({
    modelValue: (__VLS_ctx.form.instance_id),
}));
const __VLS_347 = __VLS_346({
    modelValue: (__VLS_ctx.form.instance_id),
}, ...__VLS_functionalComponentArgsRest(__VLS_346));
var __VLS_344;
var __VLS_340;
const __VLS_349 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_350 = __VLS_asFunctionalComponent(__VLS_349, new __VLS_349({
    span: (12),
}));
const __VLS_351 = __VLS_350({
    span: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_350));
__VLS_352.slots.default;
const __VLS_353 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_354 = __VLS_asFunctionalComponent(__VLS_353, new __VLS_353({
    label: "地域",
}));
const __VLS_355 = __VLS_354({
    label: "地域",
}, ...__VLS_functionalComponentArgsRest(__VLS_354));
__VLS_356.slots.default;
const __VLS_357 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_358 = __VLS_asFunctionalComponent(__VLS_357, new __VLS_357({
    modelValue: (__VLS_ctx.form.region),
}));
const __VLS_359 = __VLS_358({
    modelValue: (__VLS_ctx.form.region),
}, ...__VLS_functionalComponentArgsRest(__VLS_358));
var __VLS_356;
var __VLS_352;
const __VLS_361 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_362 = __VLS_asFunctionalComponent(__VLS_361, new __VLS_361({
    span: (12),
}));
const __VLS_363 = __VLS_362({
    span: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_362));
__VLS_364.slots.default;
const __VLS_365 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_366 = __VLS_asFunctionalComponent(__VLS_365, new __VLS_365({
    label: "Agent状态",
}));
const __VLS_367 = __VLS_366({
    label: "Agent状态",
}, ...__VLS_functionalComponentArgsRest(__VLS_366));
__VLS_368.slots.default;
const __VLS_369 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_370 = __VLS_asFunctionalComponent(__VLS_369, new __VLS_369({
    modelValue: (__VLS_ctx.form.agent_status),
}));
const __VLS_371 = __VLS_370({
    modelValue: (__VLS_ctx.form.agent_status),
}, ...__VLS_functionalComponentArgsRest(__VLS_370));
var __VLS_368;
var __VLS_364;
const __VLS_373 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_374 = __VLS_asFunctionalComponent(__VLS_373, new __VLS_373({
    span: (24),
}));
const __VLS_375 = __VLS_374({
    span: (24),
}, ...__VLS_functionalComponentArgsRest(__VLS_374));
__VLS_376.slots.default;
const __VLS_377 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_378 = __VLS_asFunctionalComponent(__VLS_377, new __VLS_377({
    label: "标签JSON",
}));
const __VLS_379 = __VLS_378({
    label: "标签JSON",
}, ...__VLS_functionalComponentArgsRest(__VLS_378));
__VLS_380.slots.default;
const __VLS_381 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_382 = __VLS_asFunctionalComponent(__VLS_381, new __VLS_381({
    modelValue: (__VLS_ctx.form.labels),
    type: "textarea",
    rows: (2),
}));
const __VLS_383 = __VLS_382({
    modelValue: (__VLS_ctx.form.labels),
    type: "textarea",
    rows: (2),
}, ...__VLS_functionalComponentArgsRest(__VLS_382));
var __VLS_380;
var __VLS_376;
const __VLS_385 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_386 = __VLS_asFunctionalComponent(__VLS_385, new __VLS_385({
    span: (12),
}));
const __VLS_387 = __VLS_386({
    span: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_386));
__VLS_388.slots.default;
const __VLS_389 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_390 = __VLS_asFunctionalComponent(__VLS_389, new __VLS_389({
    label: "状态",
}));
const __VLS_391 = __VLS_390({
    label: "状态",
}, ...__VLS_functionalComponentArgsRest(__VLS_390));
__VLS_392.slots.default;
const __VLS_393 = {}.ElSwitch;
/** @type {[typeof __VLS_components.ElSwitch, typeof __VLS_components.elSwitch, ]} */ ;
// @ts-ignore
const __VLS_394 = __VLS_asFunctionalComponent(__VLS_393, new __VLS_393({
    modelValue: (__VLS_ctx.statusSwitch),
}));
const __VLS_395 = __VLS_394({
    modelValue: (__VLS_ctx.statusSwitch),
}, ...__VLS_functionalComponentArgsRest(__VLS_394));
var __VLS_392;
var __VLS_388;
var __VLS_192;
var __VLS_188;
{
    const { footer: __VLS_thisSlot } = __VLS_184.slots;
    const __VLS_397 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_398 = __VLS_asFunctionalComponent(__VLS_397, new __VLS_397({
        ...{ 'onClick': {} },
    }));
    const __VLS_399 = __VLS_398({
        ...{ 'onClick': {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_398));
    let __VLS_401;
    let __VLS_402;
    let __VLS_403;
    const __VLS_404 = {
        onClick: (...[$event]) => {
            __VLS_ctx.dialogVisible = false;
        }
    };
    __VLS_400.slots.default;
    var __VLS_400;
    const __VLS_405 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_406 = __VLS_asFunctionalComponent(__VLS_405, new __VLS_405({
        ...{ 'onClick': {} },
        type: "primary",
    }));
    const __VLS_407 = __VLS_406({
        ...{ 'onClick': {} },
        type: "primary",
    }, ...__VLS_functionalComponentArgsRest(__VLS_406));
    let __VLS_409;
    let __VLS_410;
    let __VLS_411;
    const __VLS_412 = {
        onClick: (__VLS_ctx.submit)
    };
    __VLS_408.slots.default;
    var __VLS_408;
}
var __VLS_184;
const __VLS_413 = {}.ElDialog;
/** @type {[typeof __VLS_components.ElDialog, typeof __VLS_components.elDialog, typeof __VLS_components.ElDialog, typeof __VLS_components.elDialog, ]} */ ;
// @ts-ignore
const __VLS_414 = __VLS_asFunctionalComponent(__VLS_413, new __VLS_413({
    modelValue: (__VLS_ctx.detailVisible),
    title: "主机详情",
    width: "620px",
}));
const __VLS_415 = __VLS_414({
    modelValue: (__VLS_ctx.detailVisible),
    title: "主机详情",
    width: "620px",
}, ...__VLS_functionalComponentArgsRest(__VLS_414));
__VLS_416.slots.default;
if (__VLS_ctx.detail) {
    const __VLS_417 = {}.ElDescriptions;
    /** @type {[typeof __VLS_components.ElDescriptions, typeof __VLS_components.elDescriptions, typeof __VLS_components.ElDescriptions, typeof __VLS_components.elDescriptions, ]} */ ;
    // @ts-ignore
    const __VLS_418 = __VLS_asFunctionalComponent(__VLS_417, new __VLS_417({
        column: (2),
        border: true,
    }));
    const __VLS_419 = __VLS_418({
        column: (2),
        border: true,
    }, ...__VLS_functionalComponentArgsRest(__VLS_418));
    __VLS_420.slots.default;
    const __VLS_421 = {}.ElDescriptionsItem;
    /** @type {[typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_422 = __VLS_asFunctionalComponent(__VLS_421, new __VLS_421({
        label: "主机名",
    }));
    const __VLS_423 = __VLS_422({
        label: "主机名",
    }, ...__VLS_functionalComponentArgsRest(__VLS_422));
    __VLS_424.slots.default;
    (__VLS_ctx.detail.hostname);
    var __VLS_424;
    const __VLS_425 = {}.ElDescriptionsItem;
    /** @type {[typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_426 = __VLS_asFunctionalComponent(__VLS_425, new __VLS_425({
        label: "内网IP",
    }));
    const __VLS_427 = __VLS_426({
        label: "内网IP",
    }, ...__VLS_functionalComponentArgsRest(__VLS_426));
    __VLS_428.slots.default;
    (__VLS_ctx.detail.inner_ip);
    var __VLS_428;
    const __VLS_429 = {}.ElDescriptionsItem;
    /** @type {[typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_430 = __VLS_asFunctionalComponent(__VLS_429, new __VLS_429({
        label: "公网IP",
    }));
    const __VLS_431 = __VLS_430({
        label: "公网IP",
    }, ...__VLS_functionalComponentArgsRest(__VLS_430));
    __VLS_432.slots.default;
    (__VLS_ctx.detail.outer_ip);
    var __VLS_432;
    const __VLS_433 = {}.ElDescriptionsItem;
    /** @type {[typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_434 = __VLS_asFunctionalComponent(__VLS_433, new __VLS_433({
        label: "系统",
    }));
    const __VLS_435 = __VLS_434({
        label: "系统",
    }, ...__VLS_functionalComponentArgsRest(__VLS_434));
    __VLS_436.slots.default;
    (__VLS_ctx.detail.os_type);
    (__VLS_ctx.detail.os_version);
    var __VLS_436;
    const __VLS_437 = {}.ElDescriptionsItem;
    /** @type {[typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_438 = __VLS_asFunctionalComponent(__VLS_437, new __VLS_437({
        label: "CPU/内存",
    }));
    const __VLS_439 = __VLS_438({
        label: "CPU/内存",
    }, ...__VLS_functionalComponentArgsRest(__VLS_438));
    __VLS_440.slots.default;
    (__VLS_ctx.detail.cpu_cores);
    (__VLS_ctx.detail.memory_gb);
    var __VLS_440;
    const __VLS_441 = {}.ElDescriptionsItem;
    /** @type {[typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_442 = __VLS_asFunctionalComponent(__VLS_441, new __VLS_441({
        label: "磁盘",
    }));
    const __VLS_443 = __VLS_442({
        label: "磁盘",
    }, ...__VLS_functionalComponentArgsRest(__VLS_442));
    __VLS_444.slots.default;
    (__VLS_ctx.detail.disk_gb);
    var __VLS_444;
    const __VLS_445 = {}.ElDescriptionsItem;
    /** @type {[typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_446 = __VLS_asFunctionalComponent(__VLS_445, new __VLS_445({
        label: "分组ID",
    }));
    const __VLS_447 = __VLS_446({
        label: "分组ID",
    }, ...__VLS_functionalComponentArgsRest(__VLS_446));
    __VLS_448.slots.default;
    (__VLS_ctx.detail.group_id);
    var __VLS_448;
    const __VLS_449 = {}.ElDescriptionsItem;
    /** @type {[typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_450 = __VLS_asFunctionalComponent(__VLS_449, new __VLS_449({
        label: "凭据ID",
    }));
    const __VLS_451 = __VLS_450({
        label: "凭据ID",
    }, ...__VLS_functionalComponentArgsRest(__VLS_450));
    __VLS_452.slots.default;
    (__VLS_ctx.detail.credential_id);
    var __VLS_452;
    var __VLS_420;
}
var __VLS_416;
var __VLS_3;
/** @type {__VLS_StyleScopedClasses['toolbar']} */ ;
/** @type {__VLS_StyleScopedClasses['actions']} */ ;
/** @type {__VLS_StyleScopedClasses['search']} */ ;
/** @type {__VLS_StyleScopedClasses['pager']} */ ;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            loading: loading,
            list: list,
            total: total,
            groupOptions: groupOptions,
            dialogVisible: dialogVisible,
            detailVisible: detailVisible,
            detail: detail,
            query: query,
            form: form,
            statusSwitch: statusSwitch,
            loadData: loadData,
            handleSelectionChange: handleSelectionChange,
            handlePageChange: handlePageChange,
            openCreate: openCreate,
            openEdit: openEdit,
            openDetail: openDetail,
            submit: submit,
            handleDelete: handleDelete,
            handleBatch: handleBatch,
            handleTest: handleTest,
            openTerminal: openTerminal,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
