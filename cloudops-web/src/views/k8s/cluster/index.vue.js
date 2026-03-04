import { computed, onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { createK8sClusterApi, deleteK8sClusterApi, getK8sClusterApi, listK8sClustersApi, testK8sClusterApi, updateK8sClusterApi, } from '@/api/k8s';
const loading = ref(false);
const list = ref([]);
const dialogVisible = ref(false);
const detailVisible = ref(false);
const detail = ref(null);
const form = reactive({
    id: 0,
    name: '',
    api_server_url: '',
    description: '',
    kubeconfig: '',
    status: 1,
});
const statusSwitch = computed({
    get: () => form.status === 1,
    set: (v) => (form.status = v ? 1 : 2),
});
const loadData = async () => {
    loading.value = true;
    try {
        const res = await listK8sClustersApi();
        list.value = res.list || [];
    }
    finally {
        loading.value = false;
    }
};
const openCreate = () => {
    Object.assign(form, { id: 0, name: '', api_server_url: '', description: '', kubeconfig: '', status: 1 });
    dialogVisible.value = true;
};
const openEdit = (row) => {
    Object.assign(form, { id: row.id, name: row.name, api_server_url: row.api_server_url, description: row.description, kubeconfig: '', status: row.status });
    dialogVisible.value = true;
};
const openDetail = async (id) => {
    detail.value = await getK8sClusterApi(id);
    detailVisible.value = true;
};
const submit = async () => {
    if (form.id) {
        await updateK8sClusterApi(form.id, form);
    }
    else {
        await createK8sClusterApi(form);
    }
    ElMessage.success('保存成功');
    dialogVisible.value = false;
    await loadData();
};
const handleTest = async (id) => {
    await testK8sClusterApi(id);
    ElMessage.success('连接成功');
};
const handleDelete = async (id) => {
    await ElMessageBox.confirm('确认删除该集群吗？', '提示', { type: 'warning' });
    await deleteK8sClusterApi(id);
    ElMessage.success('删除成功');
    await loadData();
};
onMounted(loadData);
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
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('k8s:cluster:create') }, null, null);
    __VLS_8.slots.default;
    var __VLS_8;
}
const __VLS_13 = {}.ElTable;
/** @type {[typeof __VLS_components.ElTable, typeof __VLS_components.elTable, typeof __VLS_components.ElTable, typeof __VLS_components.elTable, ]} */ ;
// @ts-ignore
const __VLS_14 = __VLS_asFunctionalComponent(__VLS_13, new __VLS_13({
    data: (__VLS_ctx.list),
}));
const __VLS_15 = __VLS_14({
    data: (__VLS_ctx.list),
}, ...__VLS_functionalComponentArgsRest(__VLS_14));
__VLS_asFunctionalDirective(__VLS_directives.vLoading)(null, { ...__VLS_directiveBindingRestFields, value: (__VLS_ctx.loading) }, null, null);
__VLS_16.slots.default;
const __VLS_17 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_18 = __VLS_asFunctionalComponent(__VLS_17, new __VLS_17({
    prop: "name",
    label: "名称",
    minWidth: "140",
}));
const __VLS_19 = __VLS_18({
    prop: "name",
    label: "名称",
    minWidth: "140",
}, ...__VLS_functionalComponentArgsRest(__VLS_18));
const __VLS_21 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_22 = __VLS_asFunctionalComponent(__VLS_21, new __VLS_21({
    prop: "api_server_url",
    label: "API Server",
    minWidth: "220",
}));
const __VLS_23 = __VLS_22({
    prop: "api_server_url",
    label: "API Server",
    minWidth: "220",
}, ...__VLS_functionalComponentArgsRest(__VLS_22));
const __VLS_25 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_26 = __VLS_asFunctionalComponent(__VLS_25, new __VLS_25({
    prop: "version",
    label: "版本",
    width: "120",
}));
const __VLS_27 = __VLS_26({
    prop: "version",
    label: "版本",
    width: "120",
}, ...__VLS_functionalComponentArgsRest(__VLS_26));
const __VLS_29 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_30 = __VLS_asFunctionalComponent(__VLS_29, new __VLS_29({
    prop: "node_count",
    label: "节点",
    width: "90",
}));
const __VLS_31 = __VLS_30({
    prop: "node_count",
    label: "节点",
    width: "90",
}, ...__VLS_functionalComponentArgsRest(__VLS_30));
const __VLS_33 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_34 = __VLS_asFunctionalComponent(__VLS_33, new __VLS_33({
    prop: "pod_count",
    label: "Pod",
    width: "90",
}));
const __VLS_35 = __VLS_34({
    prop: "pod_count",
    label: "Pod",
    width: "90",
}, ...__VLS_functionalComponentArgsRest(__VLS_34));
const __VLS_37 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_38 = __VLS_asFunctionalComponent(__VLS_37, new __VLS_37({
    prop: "health",
    label: "健康状态",
    width: "120",
}));
const __VLS_39 = __VLS_38({
    prop: "health",
    label: "健康状态",
    width: "120",
}, ...__VLS_functionalComponentArgsRest(__VLS_38));
const __VLS_41 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_42 = __VLS_asFunctionalComponent(__VLS_41, new __VLS_41({
    label: "状态",
    width: "100",
}));
const __VLS_43 = __VLS_42({
    label: "状态",
    width: "100",
}, ...__VLS_functionalComponentArgsRest(__VLS_42));
__VLS_44.slots.default;
{
    const { default: __VLS_thisSlot } = __VLS_44.slots;
    const [scope] = __VLS_getSlotParams(__VLS_thisSlot);
    (scope.row.status === 1 ? '启用' : '禁用');
}
var __VLS_44;
const __VLS_45 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_46 = __VLS_asFunctionalComponent(__VLS_45, new __VLS_45({
    label: "操作",
    minWidth: "320",
    fixed: "right",
}));
const __VLS_47 = __VLS_46({
    label: "操作",
    minWidth: "320",
    fixed: "right",
}, ...__VLS_functionalComponentArgsRest(__VLS_46));
__VLS_48.slots.default;
{
    const { default: __VLS_thisSlot } = __VLS_48.slots;
    const [scope] = __VLS_getSlotParams(__VLS_thisSlot);
    const __VLS_49 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_50 = __VLS_asFunctionalComponent(__VLS_49, new __VLS_49({
        ...{ 'onClick': {} },
        link: true,
        type: "primary",
    }));
    const __VLS_51 = __VLS_50({
        ...{ 'onClick': {} },
        link: true,
        type: "primary",
    }, ...__VLS_functionalComponentArgsRest(__VLS_50));
    let __VLS_53;
    let __VLS_54;
    let __VLS_55;
    const __VLS_56 = {
        onClick: (...[$event]) => {
            __VLS_ctx.openDetail(scope.row.id);
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('k8s:cluster:detail') }, null, null);
    __VLS_52.slots.default;
    var __VLS_52;
    const __VLS_57 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_58 = __VLS_asFunctionalComponent(__VLS_57, new __VLS_57({
        ...{ 'onClick': {} },
        link: true,
        type: "primary",
    }));
    const __VLS_59 = __VLS_58({
        ...{ 'onClick': {} },
        link: true,
        type: "primary",
    }, ...__VLS_functionalComponentArgsRest(__VLS_58));
    let __VLS_61;
    let __VLS_62;
    let __VLS_63;
    const __VLS_64 = {
        onClick: (...[$event]) => {
            __VLS_ctx.openEdit(scope.row);
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('k8s:cluster:update') }, null, null);
    __VLS_60.slots.default;
    var __VLS_60;
    const __VLS_65 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_66 = __VLS_asFunctionalComponent(__VLS_65, new __VLS_65({
        ...{ 'onClick': {} },
        link: true,
        type: "success",
    }));
    const __VLS_67 = __VLS_66({
        ...{ 'onClick': {} },
        link: true,
        type: "success",
    }, ...__VLS_functionalComponentArgsRest(__VLS_66));
    let __VLS_69;
    let __VLS_70;
    let __VLS_71;
    const __VLS_72 = {
        onClick: (...[$event]) => {
            __VLS_ctx.handleTest(scope.row.id);
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('k8s:cluster:test') }, null, null);
    __VLS_68.slots.default;
    var __VLS_68;
    const __VLS_73 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_74 = __VLS_asFunctionalComponent(__VLS_73, new __VLS_73({
        ...{ 'onClick': {} },
        link: true,
        type: "danger",
    }));
    const __VLS_75 = __VLS_74({
        ...{ 'onClick': {} },
        link: true,
        type: "danger",
    }, ...__VLS_functionalComponentArgsRest(__VLS_74));
    let __VLS_77;
    let __VLS_78;
    let __VLS_79;
    const __VLS_80 = {
        onClick: (...[$event]) => {
            __VLS_ctx.handleDelete(scope.row.id);
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('k8s:cluster:delete') }, null, null);
    __VLS_76.slots.default;
    var __VLS_76;
}
var __VLS_48;
var __VLS_16;
const __VLS_81 = {}.ElDialog;
/** @type {[typeof __VLS_components.ElDialog, typeof __VLS_components.elDialog, typeof __VLS_components.ElDialog, typeof __VLS_components.elDialog, ]} */ ;
// @ts-ignore
const __VLS_82 = __VLS_asFunctionalComponent(__VLS_81, new __VLS_81({
    modelValue: (__VLS_ctx.dialogVisible),
    title: (__VLS_ctx.form.id ? '编辑集群' : '新增集群'),
    width: "720px",
}));
const __VLS_83 = __VLS_82({
    modelValue: (__VLS_ctx.dialogVisible),
    title: (__VLS_ctx.form.id ? '编辑集群' : '新增集群'),
    width: "720px",
}, ...__VLS_functionalComponentArgsRest(__VLS_82));
__VLS_84.slots.default;
const __VLS_85 = {}.ElForm;
/** @type {[typeof __VLS_components.ElForm, typeof __VLS_components.elForm, typeof __VLS_components.ElForm, typeof __VLS_components.elForm, ]} */ ;
// @ts-ignore
const __VLS_86 = __VLS_asFunctionalComponent(__VLS_85, new __VLS_85({
    model: (__VLS_ctx.form),
    labelWidth: "110px",
}));
const __VLS_87 = __VLS_86({
    model: (__VLS_ctx.form),
    labelWidth: "110px",
}, ...__VLS_functionalComponentArgsRest(__VLS_86));
__VLS_88.slots.default;
const __VLS_89 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_90 = __VLS_asFunctionalComponent(__VLS_89, new __VLS_89({
    label: "集群名称",
}));
const __VLS_91 = __VLS_90({
    label: "集群名称",
}, ...__VLS_functionalComponentArgsRest(__VLS_90));
__VLS_92.slots.default;
const __VLS_93 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_94 = __VLS_asFunctionalComponent(__VLS_93, new __VLS_93({
    modelValue: (__VLS_ctx.form.name),
}));
const __VLS_95 = __VLS_94({
    modelValue: (__VLS_ctx.form.name),
}, ...__VLS_functionalComponentArgsRest(__VLS_94));
var __VLS_92;
const __VLS_97 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_98 = __VLS_asFunctionalComponent(__VLS_97, new __VLS_97({
    label: "API Server",
}));
const __VLS_99 = __VLS_98({
    label: "API Server",
}, ...__VLS_functionalComponentArgsRest(__VLS_98));
__VLS_100.slots.default;
const __VLS_101 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_102 = __VLS_asFunctionalComponent(__VLS_101, new __VLS_101({
    modelValue: (__VLS_ctx.form.api_server_url),
}));
const __VLS_103 = __VLS_102({
    modelValue: (__VLS_ctx.form.api_server_url),
}, ...__VLS_functionalComponentArgsRest(__VLS_102));
var __VLS_100;
const __VLS_105 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_106 = __VLS_asFunctionalComponent(__VLS_105, new __VLS_105({
    label: "描述",
}));
const __VLS_107 = __VLS_106({
    label: "描述",
}, ...__VLS_functionalComponentArgsRest(__VLS_106));
__VLS_108.slots.default;
const __VLS_109 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_110 = __VLS_asFunctionalComponent(__VLS_109, new __VLS_109({
    modelValue: (__VLS_ctx.form.description),
    type: "textarea",
    rows: (2),
}));
const __VLS_111 = __VLS_110({
    modelValue: (__VLS_ctx.form.description),
    type: "textarea",
    rows: (2),
}, ...__VLS_functionalComponentArgsRest(__VLS_110));
var __VLS_108;
const __VLS_113 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_114 = __VLS_asFunctionalComponent(__VLS_113, new __VLS_113({
    label: "Kubeconfig",
}));
const __VLS_115 = __VLS_114({
    label: "Kubeconfig",
}, ...__VLS_functionalComponentArgsRest(__VLS_114));
__VLS_116.slots.default;
const __VLS_117 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_118 = __VLS_asFunctionalComponent(__VLS_117, new __VLS_117({
    modelValue: (__VLS_ctx.form.kubeconfig),
    type: "textarea",
    rows: (8),
    placeholder: "编辑时可留空",
}));
const __VLS_119 = __VLS_118({
    modelValue: (__VLS_ctx.form.kubeconfig),
    type: "textarea",
    rows: (8),
    placeholder: "编辑时可留空",
}, ...__VLS_functionalComponentArgsRest(__VLS_118));
var __VLS_116;
if (__VLS_ctx.form.id) {
    const __VLS_121 = {}.ElFormItem;
    /** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
    // @ts-ignore
    const __VLS_122 = __VLS_asFunctionalComponent(__VLS_121, new __VLS_121({
        label: "状态",
    }));
    const __VLS_123 = __VLS_122({
        label: "状态",
    }, ...__VLS_functionalComponentArgsRest(__VLS_122));
    __VLS_124.slots.default;
    const __VLS_125 = {}.ElSwitch;
    /** @type {[typeof __VLS_components.ElSwitch, typeof __VLS_components.elSwitch, ]} */ ;
    // @ts-ignore
    const __VLS_126 = __VLS_asFunctionalComponent(__VLS_125, new __VLS_125({
        modelValue: (__VLS_ctx.statusSwitch),
    }));
    const __VLS_127 = __VLS_126({
        modelValue: (__VLS_ctx.statusSwitch),
    }, ...__VLS_functionalComponentArgsRest(__VLS_126));
    var __VLS_124;
}
var __VLS_88;
{
    const { footer: __VLS_thisSlot } = __VLS_84.slots;
    const __VLS_129 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_130 = __VLS_asFunctionalComponent(__VLS_129, new __VLS_129({
        ...{ 'onClick': {} },
    }));
    const __VLS_131 = __VLS_130({
        ...{ 'onClick': {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_130));
    let __VLS_133;
    let __VLS_134;
    let __VLS_135;
    const __VLS_136 = {
        onClick: (...[$event]) => {
            __VLS_ctx.dialogVisible = false;
        }
    };
    __VLS_132.slots.default;
    var __VLS_132;
    const __VLS_137 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_138 = __VLS_asFunctionalComponent(__VLS_137, new __VLS_137({
        ...{ 'onClick': {} },
        type: "primary",
    }));
    const __VLS_139 = __VLS_138({
        ...{ 'onClick': {} },
        type: "primary",
    }, ...__VLS_functionalComponentArgsRest(__VLS_138));
    let __VLS_141;
    let __VLS_142;
    let __VLS_143;
    const __VLS_144 = {
        onClick: (__VLS_ctx.submit)
    };
    __VLS_140.slots.default;
    var __VLS_140;
}
var __VLS_84;
const __VLS_145 = {}.ElDialog;
/** @type {[typeof __VLS_components.ElDialog, typeof __VLS_components.elDialog, typeof __VLS_components.ElDialog, typeof __VLS_components.elDialog, ]} */ ;
// @ts-ignore
const __VLS_146 = __VLS_asFunctionalComponent(__VLS_145, new __VLS_145({
    modelValue: (__VLS_ctx.detailVisible),
    title: "集群详情",
    width: "640px",
}));
const __VLS_147 = __VLS_146({
    modelValue: (__VLS_ctx.detailVisible),
    title: "集群详情",
    width: "640px",
}, ...__VLS_functionalComponentArgsRest(__VLS_146));
__VLS_148.slots.default;
if (__VLS_ctx.detail) {
    const __VLS_149 = {}.ElDescriptions;
    /** @type {[typeof __VLS_components.ElDescriptions, typeof __VLS_components.elDescriptions, typeof __VLS_components.ElDescriptions, typeof __VLS_components.elDescriptions, ]} */ ;
    // @ts-ignore
    const __VLS_150 = __VLS_asFunctionalComponent(__VLS_149, new __VLS_149({
        column: (2),
        border: true,
    }));
    const __VLS_151 = __VLS_150({
        column: (2),
        border: true,
    }, ...__VLS_functionalComponentArgsRest(__VLS_150));
    __VLS_152.slots.default;
    const __VLS_153 = {}.ElDescriptionsItem;
    /** @type {[typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_154 = __VLS_asFunctionalComponent(__VLS_153, new __VLS_153({
        label: "名称",
    }));
    const __VLS_155 = __VLS_154({
        label: "名称",
    }, ...__VLS_functionalComponentArgsRest(__VLS_154));
    __VLS_156.slots.default;
    (__VLS_ctx.detail.name);
    var __VLS_156;
    const __VLS_157 = {}.ElDescriptionsItem;
    /** @type {[typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_158 = __VLS_asFunctionalComponent(__VLS_157, new __VLS_157({
        label: "健康状态",
    }));
    const __VLS_159 = __VLS_158({
        label: "健康状态",
    }, ...__VLS_functionalComponentArgsRest(__VLS_158));
    __VLS_160.slots.default;
    (__VLS_ctx.detail.health || '-');
    var __VLS_160;
    const __VLS_161 = {}.ElDescriptionsItem;
    /** @type {[typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_162 = __VLS_asFunctionalComponent(__VLS_161, new __VLS_161({
        label: "API Server",
    }));
    const __VLS_163 = __VLS_162({
        label: "API Server",
    }, ...__VLS_functionalComponentArgsRest(__VLS_162));
    __VLS_164.slots.default;
    (__VLS_ctx.detail.api_server_url);
    var __VLS_164;
    const __VLS_165 = {}.ElDescriptionsItem;
    /** @type {[typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_166 = __VLS_asFunctionalComponent(__VLS_165, new __VLS_165({
        label: "版本",
    }));
    const __VLS_167 = __VLS_166({
        label: "版本",
    }, ...__VLS_functionalComponentArgsRest(__VLS_166));
    __VLS_168.slots.default;
    (__VLS_ctx.detail.version || '-');
    var __VLS_168;
    const __VLS_169 = {}.ElDescriptionsItem;
    /** @type {[typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_170 = __VLS_asFunctionalComponent(__VLS_169, new __VLS_169({
        label: "节点数",
    }));
    const __VLS_171 = __VLS_170({
        label: "节点数",
    }, ...__VLS_functionalComponentArgsRest(__VLS_170));
    __VLS_172.slots.default;
    (__VLS_ctx.detail.node_count);
    var __VLS_172;
    const __VLS_173 = {}.ElDescriptionsItem;
    /** @type {[typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_174 = __VLS_asFunctionalComponent(__VLS_173, new __VLS_173({
        label: "Pod数",
    }));
    const __VLS_175 = __VLS_174({
        label: "Pod数",
    }, ...__VLS_functionalComponentArgsRest(__VLS_174));
    __VLS_176.slots.default;
    (__VLS_ctx.detail.pod_count);
    var __VLS_176;
    const __VLS_177 = {}.ElDescriptionsItem;
    /** @type {[typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, typeof __VLS_components.ElDescriptionsItem, typeof __VLS_components.elDescriptionsItem, ]} */ ;
    // @ts-ignore
    const __VLS_178 = __VLS_asFunctionalComponent(__VLS_177, new __VLS_177({
        label: "描述",
        span: (2),
    }));
    const __VLS_179 = __VLS_178({
        label: "描述",
        span: (2),
    }, ...__VLS_functionalComponentArgsRest(__VLS_178));
    __VLS_180.slots.default;
    (__VLS_ctx.detail.description || '-');
    var __VLS_180;
    var __VLS_152;
}
var __VLS_148;
var __VLS_3;
/** @type {__VLS_StyleScopedClasses['toolbar']} */ ;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            loading: loading,
            list: list,
            dialogVisible: dialogVisible,
            detailVisible: detailVisible,
            detail: detail,
            form: form,
            statusSwitch: statusSwitch,
            openCreate: openCreate,
            openEdit: openEdit,
            openDetail: openDetail,
            submit: submit,
            handleTest: handleTest,
            handleDelete: handleDelete,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
