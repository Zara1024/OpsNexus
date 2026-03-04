import { onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { createDeploymentApi, deleteDeploymentApi, listDeploymentsApi, listK8sClustersApi, listK8sNamespacesApi, restartDeploymentApi, scaleDeploymentApi, } from '@/api/k8s';
const loading = ref(false);
const clusters = ref([]);
const namespaces = ref([]);
const clusterId = ref(undefined);
const namespace = ref('default');
const list = ref([]);
const dialogVisible = ref(false);
const form = reactive({ name: '', image: 'nginx:latest', replicas: 1 });
const scaleVisible = ref(false);
const scaleName = ref('');
const scaleReplicas = ref(1);
const loadClusters = async () => {
    const res = await listK8sClustersApi();
    clusters.value = res.list || [];
    if (!clusterId.value && clusters.value.length)
        clusterId.value = clusters.value[0].id;
};
const loadNamespaces = async () => {
    if (!clusterId.value)
        return;
    const res = await listK8sNamespacesApi(clusterId.value);
    namespaces.value = res.list || [];
    if (!namespaces.value.includes(namespace.value)) {
        namespace.value = namespaces.value[0] || 'default';
    }
};
const loadData = async () => {
    if (!clusterId.value || !namespace.value)
        return;
    loading.value = true;
    try {
        const res = await listDeploymentsApi(clusterId.value, namespace.value);
        list.value = res.list || [];
    }
    finally {
        loading.value = false;
    }
};
const handleClusterChange = async () => {
    await loadNamespaces();
    await loadData();
};
const openCreate = () => {
    Object.assign(form, { name: '', image: 'nginx:latest', replicas: 1 });
    dialogVisible.value = true;
};
const submitCreate = async () => {
    if (!clusterId.value)
        return;
    const payload = {
        apiVersion: 'apps/v1',
        kind: 'Deployment',
        metadata: { name: form.name, namespace: namespace.value },
        spec: {
            replicas: form.replicas,
            selector: { matchLabels: { app: form.name } },
            template: {
                metadata: { labels: { app: form.name } },
                spec: { containers: [{ name: form.name, image: form.image, ports: [{ containerPort: 80 }] }] },
            },
        },
    };
    await createDeploymentApi(clusterId.value, namespace.value, payload);
    ElMessage.success('创建成功');
    dialogVisible.value = false;
    await loadData();
};
const openScale = (row) => {
    scaleName.value = row.metadata?.name || '';
    scaleReplicas.value = row.spec?.replicas || 1;
    scaleVisible.value = true;
};
const submitScale = async () => {
    if (!clusterId.value || !scaleName.value)
        return;
    await scaleDeploymentApi(clusterId.value, namespace.value, scaleName.value, scaleReplicas.value);
    ElMessage.success('伸缩成功');
    scaleVisible.value = false;
    await loadData();
};
const restart = async (name) => {
    if (!clusterId.value)
        return;
    await restartDeploymentApi(clusterId.value, namespace.value, name);
    ElMessage.success('重启成功');
    await loadData();
};
const remove = async (name) => {
    if (!clusterId.value)
        return;
    await ElMessageBox.confirm(`确认删除 ${name} 吗？`, '提示', { type: 'warning' });
    await deleteDeploymentApi(clusterId.value, namespace.value, name);
    ElMessage.success('删除成功');
    await loadData();
};
onMounted(async () => {
    await loadClusters();
    await loadNamespaces();
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
    const __VLS_5 = {}.ElSelect;
    /** @type {[typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, ]} */ ;
    // @ts-ignore
    const __VLS_6 = __VLS_asFunctionalComponent(__VLS_5, new __VLS_5({
        ...{ 'onChange': {} },
        modelValue: (__VLS_ctx.clusterId),
        placeholder: "集群",
        ...{ style: {} },
    }));
    const __VLS_7 = __VLS_6({
        ...{ 'onChange': {} },
        modelValue: (__VLS_ctx.clusterId),
        placeholder: "集群",
        ...{ style: {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_6));
    let __VLS_9;
    let __VLS_10;
    let __VLS_11;
    const __VLS_12 = {
        onChange: (__VLS_ctx.handleClusterChange)
    };
    __VLS_8.slots.default;
    for (const [c] of __VLS_getVForSourceType((__VLS_ctx.clusters))) {
        const __VLS_13 = {}.ElOption;
        /** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
        // @ts-ignore
        const __VLS_14 = __VLS_asFunctionalComponent(__VLS_13, new __VLS_13({
            key: (c.id),
            label: (c.name),
            value: (c.id),
        }));
        const __VLS_15 = __VLS_14({
            key: (c.id),
            label: (c.name),
            value: (c.id),
        }, ...__VLS_functionalComponentArgsRest(__VLS_14));
    }
    var __VLS_8;
    const __VLS_17 = {}.ElSelect;
    /** @type {[typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, ]} */ ;
    // @ts-ignore
    const __VLS_18 = __VLS_asFunctionalComponent(__VLS_17, new __VLS_17({
        ...{ 'onChange': {} },
        modelValue: (__VLS_ctx.namespace),
        placeholder: "命名空间",
        ...{ style: {} },
    }));
    const __VLS_19 = __VLS_18({
        ...{ 'onChange': {} },
        modelValue: (__VLS_ctx.namespace),
        placeholder: "命名空间",
        ...{ style: {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_18));
    let __VLS_21;
    let __VLS_22;
    let __VLS_23;
    const __VLS_24 = {
        onChange: (__VLS_ctx.loadData)
    };
    __VLS_20.slots.default;
    for (const [ns] of __VLS_getVForSourceType((__VLS_ctx.namespaces))) {
        const __VLS_25 = {}.ElOption;
        /** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
        // @ts-ignore
        const __VLS_26 = __VLS_asFunctionalComponent(__VLS_25, new __VLS_25({
            key: (ns),
            label: (ns),
            value: (ns),
        }));
        const __VLS_27 = __VLS_26({
            key: (ns),
            label: (ns),
            value: (ns),
        }, ...__VLS_functionalComponentArgsRest(__VLS_26));
    }
    var __VLS_20;
    const __VLS_29 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_30 = __VLS_asFunctionalComponent(__VLS_29, new __VLS_29({
        ...{ 'onClick': {} },
        type: "primary",
    }));
    const __VLS_31 = __VLS_30({
        ...{ 'onClick': {} },
        type: "primary",
    }, ...__VLS_functionalComponentArgsRest(__VLS_30));
    let __VLS_33;
    let __VLS_34;
    let __VLS_35;
    const __VLS_36 = {
        onClick: (__VLS_ctx.openCreate)
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('k8s:workload:deployment:create') }, null, null);
    __VLS_32.slots.default;
    var __VLS_32;
}
const __VLS_37 = {}.ElTable;
/** @type {[typeof __VLS_components.ElTable, typeof __VLS_components.elTable, typeof __VLS_components.ElTable, typeof __VLS_components.elTable, ]} */ ;
// @ts-ignore
const __VLS_38 = __VLS_asFunctionalComponent(__VLS_37, new __VLS_37({
    data: (__VLS_ctx.list),
}));
const __VLS_39 = __VLS_38({
    data: (__VLS_ctx.list),
}, ...__VLS_functionalComponentArgsRest(__VLS_38));
__VLS_asFunctionalDirective(__VLS_directives.vLoading)(null, { ...__VLS_directiveBindingRestFields, value: (__VLS_ctx.loading) }, null, null);
__VLS_40.slots.default;
const __VLS_41 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_42 = __VLS_asFunctionalComponent(__VLS_41, new __VLS_41({
    prop: "metadata.name",
    label: "名称",
    minWidth: "180",
}));
const __VLS_43 = __VLS_42({
    prop: "metadata.name",
    label: "名称",
    minWidth: "180",
}, ...__VLS_functionalComponentArgsRest(__VLS_42));
const __VLS_45 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_46 = __VLS_asFunctionalComponent(__VLS_45, new __VLS_45({
    label: "副本",
    width: "120",
}));
const __VLS_47 = __VLS_46({
    label: "副本",
    width: "120",
}, ...__VLS_functionalComponentArgsRest(__VLS_46));
__VLS_48.slots.default;
{
    const { default: __VLS_thisSlot } = __VLS_48.slots;
    const [scope] = __VLS_getSlotParams(__VLS_thisSlot);
    (scope.row.status?.readyReplicas || 0);
    (scope.row.spec?.replicas || 0);
}
var __VLS_48;
const __VLS_49 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_50 = __VLS_asFunctionalComponent(__VLS_49, new __VLS_49({
    prop: "status.availableReplicas",
    label: "可用副本",
    width: "120",
}));
const __VLS_51 = __VLS_50({
    prop: "status.availableReplicas",
    label: "可用副本",
    width: "120",
}, ...__VLS_functionalComponentArgsRest(__VLS_50));
const __VLS_53 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_54 = __VLS_asFunctionalComponent(__VLS_53, new __VLS_53({
    label: "操作",
    minWidth: "300",
    fixed: "right",
}));
const __VLS_55 = __VLS_54({
    label: "操作",
    minWidth: "300",
    fixed: "right",
}, ...__VLS_functionalComponentArgsRest(__VLS_54));
__VLS_56.slots.default;
{
    const { default: __VLS_thisSlot } = __VLS_56.slots;
    const [scope] = __VLS_getSlotParams(__VLS_thisSlot);
    const __VLS_57 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_58 = __VLS_asFunctionalComponent(__VLS_57, new __VLS_57({
        ...{ 'onClick': {} },
        link: true,
        type: "warning",
    }));
    const __VLS_59 = __VLS_58({
        ...{ 'onClick': {} },
        link: true,
        type: "warning",
    }, ...__VLS_functionalComponentArgsRest(__VLS_58));
    let __VLS_61;
    let __VLS_62;
    let __VLS_63;
    const __VLS_64 = {
        onClick: (...[$event]) => {
            __VLS_ctx.restart(scope.row.metadata?.name);
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('k8s:workload:deployment:restart') }, null, null);
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
            __VLS_ctx.openScale(scope.row);
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('k8s:workload:deployment:scale') }, null, null);
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
            __VLS_ctx.remove(scope.row.metadata?.name);
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('k8s:workload:deployment:delete') }, null, null);
    __VLS_76.slots.default;
    var __VLS_76;
}
var __VLS_56;
var __VLS_40;
const __VLS_81 = {}.ElDialog;
/** @type {[typeof __VLS_components.ElDialog, typeof __VLS_components.elDialog, typeof __VLS_components.ElDialog, typeof __VLS_components.elDialog, ]} */ ;
// @ts-ignore
const __VLS_82 = __VLS_asFunctionalComponent(__VLS_81, new __VLS_81({
    modelValue: (__VLS_ctx.dialogVisible),
    title: "创建 Deployment",
    width: "760px",
}));
const __VLS_83 = __VLS_82({
    modelValue: (__VLS_ctx.dialogVisible),
    title: "创建 Deployment",
    width: "760px",
}, ...__VLS_functionalComponentArgsRest(__VLS_82));
__VLS_84.slots.default;
const __VLS_85 = {}.ElForm;
/** @type {[typeof __VLS_components.ElForm, typeof __VLS_components.elForm, typeof __VLS_components.ElForm, typeof __VLS_components.elForm, ]} */ ;
// @ts-ignore
const __VLS_86 = __VLS_asFunctionalComponent(__VLS_85, new __VLS_85({
    model: (__VLS_ctx.form),
    labelWidth: "100px",
}));
const __VLS_87 = __VLS_86({
    model: (__VLS_ctx.form),
    labelWidth: "100px",
}, ...__VLS_functionalComponentArgsRest(__VLS_86));
__VLS_88.slots.default;
const __VLS_89 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_90 = __VLS_asFunctionalComponent(__VLS_89, new __VLS_89({
    label: "名称",
}));
const __VLS_91 = __VLS_90({
    label: "名称",
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
    label: "镜像",
}));
const __VLS_99 = __VLS_98({
    label: "镜像",
}, ...__VLS_functionalComponentArgsRest(__VLS_98));
__VLS_100.slots.default;
const __VLS_101 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_102 = __VLS_asFunctionalComponent(__VLS_101, new __VLS_101({
    modelValue: (__VLS_ctx.form.image),
}));
const __VLS_103 = __VLS_102({
    modelValue: (__VLS_ctx.form.image),
}, ...__VLS_functionalComponentArgsRest(__VLS_102));
var __VLS_100;
const __VLS_105 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_106 = __VLS_asFunctionalComponent(__VLS_105, new __VLS_105({
    label: "副本数",
}));
const __VLS_107 = __VLS_106({
    label: "副本数",
}, ...__VLS_functionalComponentArgsRest(__VLS_106));
__VLS_108.slots.default;
const __VLS_109 = {}.ElInputNumber;
/** @type {[typeof __VLS_components.ElInputNumber, typeof __VLS_components.elInputNumber, ]} */ ;
// @ts-ignore
const __VLS_110 = __VLS_asFunctionalComponent(__VLS_109, new __VLS_109({
    modelValue: (__VLS_ctx.form.replicas),
    min: (0),
}));
const __VLS_111 = __VLS_110({
    modelValue: (__VLS_ctx.form.replicas),
    min: (0),
}, ...__VLS_functionalComponentArgsRest(__VLS_110));
var __VLS_108;
var __VLS_88;
{
    const { footer: __VLS_thisSlot } = __VLS_84.slots;
    const __VLS_113 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_114 = __VLS_asFunctionalComponent(__VLS_113, new __VLS_113({
        ...{ 'onClick': {} },
    }));
    const __VLS_115 = __VLS_114({
        ...{ 'onClick': {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_114));
    let __VLS_117;
    let __VLS_118;
    let __VLS_119;
    const __VLS_120 = {
        onClick: (...[$event]) => {
            __VLS_ctx.dialogVisible = false;
        }
    };
    __VLS_116.slots.default;
    var __VLS_116;
    const __VLS_121 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_122 = __VLS_asFunctionalComponent(__VLS_121, new __VLS_121({
        ...{ 'onClick': {} },
        type: "primary",
    }));
    const __VLS_123 = __VLS_122({
        ...{ 'onClick': {} },
        type: "primary",
    }, ...__VLS_functionalComponentArgsRest(__VLS_122));
    let __VLS_125;
    let __VLS_126;
    let __VLS_127;
    const __VLS_128 = {
        onClick: (__VLS_ctx.submitCreate)
    };
    __VLS_124.slots.default;
    var __VLS_124;
}
var __VLS_84;
const __VLS_129 = {}.ElDialog;
/** @type {[typeof __VLS_components.ElDialog, typeof __VLS_components.elDialog, typeof __VLS_components.ElDialog, typeof __VLS_components.elDialog, ]} */ ;
// @ts-ignore
const __VLS_130 = __VLS_asFunctionalComponent(__VLS_129, new __VLS_129({
    modelValue: (__VLS_ctx.scaleVisible),
    title: "伸缩 Deployment",
    width: "420px",
}));
const __VLS_131 = __VLS_130({
    modelValue: (__VLS_ctx.scaleVisible),
    title: "伸缩 Deployment",
    width: "420px",
}, ...__VLS_functionalComponentArgsRest(__VLS_130));
__VLS_132.slots.default;
const __VLS_133 = {}.ElForm;
/** @type {[typeof __VLS_components.ElForm, typeof __VLS_components.elForm, typeof __VLS_components.ElForm, typeof __VLS_components.elForm, ]} */ ;
// @ts-ignore
const __VLS_134 = __VLS_asFunctionalComponent(__VLS_133, new __VLS_133({}));
const __VLS_135 = __VLS_134({}, ...__VLS_functionalComponentArgsRest(__VLS_134));
__VLS_136.slots.default;
const __VLS_137 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_138 = __VLS_asFunctionalComponent(__VLS_137, new __VLS_137({
    label: "副本数",
    labelWidth: "80px",
}));
const __VLS_139 = __VLS_138({
    label: "副本数",
    labelWidth: "80px",
}, ...__VLS_functionalComponentArgsRest(__VLS_138));
__VLS_140.slots.default;
const __VLS_141 = {}.ElInputNumber;
/** @type {[typeof __VLS_components.ElInputNumber, typeof __VLS_components.elInputNumber, ]} */ ;
// @ts-ignore
const __VLS_142 = __VLS_asFunctionalComponent(__VLS_141, new __VLS_141({
    modelValue: (__VLS_ctx.scaleReplicas),
    min: (0),
}));
const __VLS_143 = __VLS_142({
    modelValue: (__VLS_ctx.scaleReplicas),
    min: (0),
}, ...__VLS_functionalComponentArgsRest(__VLS_142));
var __VLS_140;
var __VLS_136;
{
    const { footer: __VLS_thisSlot } = __VLS_132.slots;
    const __VLS_145 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_146 = __VLS_asFunctionalComponent(__VLS_145, new __VLS_145({
        ...{ 'onClick': {} },
    }));
    const __VLS_147 = __VLS_146({
        ...{ 'onClick': {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_146));
    let __VLS_149;
    let __VLS_150;
    let __VLS_151;
    const __VLS_152 = {
        onClick: (...[$event]) => {
            __VLS_ctx.scaleVisible = false;
        }
    };
    __VLS_148.slots.default;
    var __VLS_148;
    const __VLS_153 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_154 = __VLS_asFunctionalComponent(__VLS_153, new __VLS_153({
        ...{ 'onClick': {} },
        type: "primary",
    }));
    const __VLS_155 = __VLS_154({
        ...{ 'onClick': {} },
        type: "primary",
    }, ...__VLS_functionalComponentArgsRest(__VLS_154));
    let __VLS_157;
    let __VLS_158;
    let __VLS_159;
    const __VLS_160 = {
        onClick: (__VLS_ctx.submitScale)
    };
    __VLS_156.slots.default;
    var __VLS_156;
}
var __VLS_132;
var __VLS_3;
/** @type {__VLS_StyleScopedClasses['toolbar']} */ ;
/** @type {__VLS_StyleScopedClasses['actions']} */ ;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            loading: loading,
            clusters: clusters,
            namespaces: namespaces,
            clusterId: clusterId,
            namespace: namespace,
            list: list,
            dialogVisible: dialogVisible,
            form: form,
            scaleVisible: scaleVisible,
            scaleReplicas: scaleReplicas,
            loadData: loadData,
            handleClusterChange: handleClusterChange,
            openCreate: openCreate,
            submitCreate: submitCreate,
            openScale: openScale,
            submitScale: submitScale,
            restart: restart,
            remove: remove,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
