import { onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { deleteConfigMapApi, listConfigMapsApi, listK8sClustersApi, listK8sNamespacesApi, saveConfigMapApi } from '@/api/k8s';
const loading = ref(false);
const clusters = ref([]);
const namespaces = ref([]);
const clusterId = ref();
const namespace = ref('default');
const list = ref([]);
const dialogVisible = ref(false);
const form = reactive({ name: '', dataText: '{\n  "key": "value"\n}' });
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
    if (!namespaces.value.includes(namespace.value))
        namespace.value = namespaces.value[0] || 'default';
};
const loadData = async () => {
    if (!clusterId.value)
        return;
    loading.value = true;
    try {
        const res = await listConfigMapsApi(clusterId.value, namespace.value);
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
    form.name = '';
    form.dataText = '{\n  "key": "value"\n}';
    dialogVisible.value = true;
};
const openEdit = (row) => {
    form.name = row.metadata?.name || '';
    form.dataText = JSON.stringify(row.data || {}, null, 2);
    dialogVisible.value = true;
};
const submit = async () => {
    if (!clusterId.value)
        return;
    const data = JSON.parse(form.dataText || '{}');
    await saveConfigMapApi(clusterId.value, namespace.value, { metadata: { name: form.name, namespace: namespace.value }, data });
    ElMessage.success('保存成功');
    dialogVisible.value = false;
    await loadData();
};
const remove = async (name) => {
    if (!clusterId.value)
        return;
    await ElMessageBox.confirm(`确认删除 ${name} 吗？`, '提示', { type: 'warning' });
    await deleteConfigMapApi(clusterId.value, namespace.value, name);
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
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('k8s:config:configmap:save') }, null, null);
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
    minWidth: "220",
}));
const __VLS_43 = __VLS_42({
    prop: "metadata.name",
    label: "名称",
    minWidth: "220",
}, ...__VLS_functionalComponentArgsRest(__VLS_42));
const __VLS_45 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_46 = __VLS_asFunctionalComponent(__VLS_45, new __VLS_45({
    label: "键数量",
    width: "120",
}));
const __VLS_47 = __VLS_46({
    label: "键数量",
    width: "120",
}, ...__VLS_functionalComponentArgsRest(__VLS_46));
__VLS_48.slots.default;
{
    const { default: __VLS_thisSlot } = __VLS_48.slots;
    const [scope] = __VLS_getSlotParams(__VLS_thisSlot);
    (Object.keys(scope.row.data || {}).length);
}
var __VLS_48;
const __VLS_49 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_50 = __VLS_asFunctionalComponent(__VLS_49, new __VLS_49({
    label: "操作",
    width: "220",
    fixed: "right",
}));
const __VLS_51 = __VLS_50({
    label: "操作",
    width: "220",
    fixed: "right",
}, ...__VLS_functionalComponentArgsRest(__VLS_50));
__VLS_52.slots.default;
{
    const { default: __VLS_thisSlot } = __VLS_52.slots;
    const [scope] = __VLS_getSlotParams(__VLS_thisSlot);
    const __VLS_53 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_54 = __VLS_asFunctionalComponent(__VLS_53, new __VLS_53({
        ...{ 'onClick': {} },
        link: true,
        type: "primary",
    }));
    const __VLS_55 = __VLS_54({
        ...{ 'onClick': {} },
        link: true,
        type: "primary",
    }, ...__VLS_functionalComponentArgsRest(__VLS_54));
    let __VLS_57;
    let __VLS_58;
    let __VLS_59;
    const __VLS_60 = {
        onClick: (...[$event]) => {
            __VLS_ctx.openEdit(scope.row);
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('k8s:config:configmap:save') }, null, null);
    __VLS_56.slots.default;
    var __VLS_56;
    const __VLS_61 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_62 = __VLS_asFunctionalComponent(__VLS_61, new __VLS_61({
        ...{ 'onClick': {} },
        link: true,
        type: "danger",
    }));
    const __VLS_63 = __VLS_62({
        ...{ 'onClick': {} },
        link: true,
        type: "danger",
    }, ...__VLS_functionalComponentArgsRest(__VLS_62));
    let __VLS_65;
    let __VLS_66;
    let __VLS_67;
    const __VLS_68 = {
        onClick: (...[$event]) => {
            __VLS_ctx.remove(scope.row.metadata?.name);
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('k8s:config:configmap:delete') }, null, null);
    __VLS_64.slots.default;
    var __VLS_64;
}
var __VLS_52;
var __VLS_40;
const __VLS_69 = {}.ElDialog;
/** @type {[typeof __VLS_components.ElDialog, typeof __VLS_components.elDialog, typeof __VLS_components.ElDialog, typeof __VLS_components.elDialog, ]} */ ;
// @ts-ignore
const __VLS_70 = __VLS_asFunctionalComponent(__VLS_69, new __VLS_69({
    modelValue: (__VLS_ctx.dialogVisible),
    title: "保存 ConfigMap",
    width: "760px",
}));
const __VLS_71 = __VLS_70({
    modelValue: (__VLS_ctx.dialogVisible),
    title: "保存 ConfigMap",
    width: "760px",
}, ...__VLS_functionalComponentArgsRest(__VLS_70));
__VLS_72.slots.default;
const __VLS_73 = {}.ElForm;
/** @type {[typeof __VLS_components.ElForm, typeof __VLS_components.elForm, typeof __VLS_components.ElForm, typeof __VLS_components.elForm, ]} */ ;
// @ts-ignore
const __VLS_74 = __VLS_asFunctionalComponent(__VLS_73, new __VLS_73({
    model: (__VLS_ctx.form),
    labelWidth: "90px",
}));
const __VLS_75 = __VLS_74({
    model: (__VLS_ctx.form),
    labelWidth: "90px",
}, ...__VLS_functionalComponentArgsRest(__VLS_74));
__VLS_76.slots.default;
const __VLS_77 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_78 = __VLS_asFunctionalComponent(__VLS_77, new __VLS_77({
    label: "名称",
}));
const __VLS_79 = __VLS_78({
    label: "名称",
}, ...__VLS_functionalComponentArgsRest(__VLS_78));
__VLS_80.slots.default;
const __VLS_81 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_82 = __VLS_asFunctionalComponent(__VLS_81, new __VLS_81({
    modelValue: (__VLS_ctx.form.name),
}));
const __VLS_83 = __VLS_82({
    modelValue: (__VLS_ctx.form.name),
}, ...__VLS_functionalComponentArgsRest(__VLS_82));
var __VLS_80;
const __VLS_85 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
const __VLS_86 = __VLS_asFunctionalComponent(__VLS_85, new __VLS_85({
    label: "内容JSON",
}));
const __VLS_87 = __VLS_86({
    label: "内容JSON",
}, ...__VLS_functionalComponentArgsRest(__VLS_86));
__VLS_88.slots.default;
const __VLS_89 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_90 = __VLS_asFunctionalComponent(__VLS_89, new __VLS_89({
    modelValue: (__VLS_ctx.form.dataText),
    type: "textarea",
    rows: (14),
}));
const __VLS_91 = __VLS_90({
    modelValue: (__VLS_ctx.form.dataText),
    type: "textarea",
    rows: (14),
}, ...__VLS_functionalComponentArgsRest(__VLS_90));
var __VLS_88;
var __VLS_76;
{
    const { footer: __VLS_thisSlot } = __VLS_72.slots;
    const __VLS_93 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_94 = __VLS_asFunctionalComponent(__VLS_93, new __VLS_93({
        ...{ 'onClick': {} },
    }));
    const __VLS_95 = __VLS_94({
        ...{ 'onClick': {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_94));
    let __VLS_97;
    let __VLS_98;
    let __VLS_99;
    const __VLS_100 = {
        onClick: (...[$event]) => {
            __VLS_ctx.dialogVisible = false;
        }
    };
    __VLS_96.slots.default;
    var __VLS_96;
    const __VLS_101 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_102 = __VLS_asFunctionalComponent(__VLS_101, new __VLS_101({
        ...{ 'onClick': {} },
        type: "primary",
    }));
    const __VLS_103 = __VLS_102({
        ...{ 'onClick': {} },
        type: "primary",
    }, ...__VLS_functionalComponentArgsRest(__VLS_102));
    let __VLS_105;
    let __VLS_106;
    let __VLS_107;
    const __VLS_108 = {
        onClick: (__VLS_ctx.submit)
    };
    __VLS_104.slots.default;
    var __VLS_104;
}
var __VLS_72;
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
            loadData: loadData,
            handleClusterChange: handleClusterChange,
            openCreate: openCreate,
            openEdit: openEdit,
            submit: submit,
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
