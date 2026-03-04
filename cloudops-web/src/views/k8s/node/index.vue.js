import { onMounted, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { cordonK8sNodeApi, drainK8sNodeApi, listK8sClustersApi, listK8sNodesApi, uncordonK8sNodeApi } from '@/api/k8s';
const loading = ref(false);
const list = ref([]);
const clusterOptions = ref([]);
const clusterId = ref(undefined);
const loadClusters = async () => {
    const res = await listK8sClustersApi();
    clusterOptions.value = res.list || [];
    if (!clusterId.value && clusterOptions.value.length) {
        clusterId.value = clusterOptions.value[0].id;
    }
};
const loadData = async () => {
    if (!clusterId.value)
        return;
    loading.value = true;
    try {
        const res = await listK8sNodesApi(clusterId.value);
        list.value = res.list || [];
    }
    finally {
        loading.value = false;
    }
};
const handleCordon = async (name) => {
    if (!clusterId.value)
        return;
    await cordonK8sNodeApi(clusterId.value, name);
    ElMessage.success('封锁成功');
    await loadData();
};
const handleUncordon = async (name) => {
    if (!clusterId.value)
        return;
    await uncordonK8sNodeApi(clusterId.value, name);
    ElMessage.success('解封成功');
    await loadData();
};
const handleDrain = async (name) => {
    if (!clusterId.value)
        return;
    await ElMessageBox.confirm(`确认对节点 ${name} 执行排水吗？`, '提示', { type: 'warning' });
    await drainK8sNodeApi(clusterId.value, name);
    ElMessage.success('排水任务已执行');
    await loadData();
};
onMounted(async () => {
    await loadClusters();
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
        placeholder: "选择集群",
        ...{ style: {} },
    }));
    const __VLS_7 = __VLS_6({
        ...{ 'onChange': {} },
        modelValue: (__VLS_ctx.clusterId),
        placeholder: "选择集群",
        ...{ style: {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_6));
    let __VLS_9;
    let __VLS_10;
    let __VLS_11;
    const __VLS_12 = {
        onChange: (__VLS_ctx.loadData)
    };
    __VLS_8.slots.default;
    for (const [c] of __VLS_getVForSourceType((__VLS_ctx.clusterOptions))) {
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
    const __VLS_17 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_18 = __VLS_asFunctionalComponent(__VLS_17, new __VLS_17({
        ...{ 'onClick': {} },
    }));
    const __VLS_19 = __VLS_18({
        ...{ 'onClick': {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_18));
    let __VLS_21;
    let __VLS_22;
    let __VLS_23;
    const __VLS_24 = {
        onClick: (__VLS_ctx.loadData)
    };
    __VLS_20.slots.default;
    var __VLS_20;
}
const __VLS_25 = {}.ElTable;
/** @type {[typeof __VLS_components.ElTable, typeof __VLS_components.elTable, typeof __VLS_components.ElTable, typeof __VLS_components.elTable, ]} */ ;
// @ts-ignore
const __VLS_26 = __VLS_asFunctionalComponent(__VLS_25, new __VLS_25({
    data: (__VLS_ctx.list),
}));
const __VLS_27 = __VLS_26({
    data: (__VLS_ctx.list),
}, ...__VLS_functionalComponentArgsRest(__VLS_26));
__VLS_asFunctionalDirective(__VLS_directives.vLoading)(null, { ...__VLS_directiveBindingRestFields, value: (__VLS_ctx.loading) }, null, null);
__VLS_28.slots.default;
const __VLS_29 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_30 = __VLS_asFunctionalComponent(__VLS_29, new __VLS_29({
    prop: "name",
    label: "节点名",
    minWidth: "220",
}));
const __VLS_31 = __VLS_30({
    prop: "name",
    label: "节点名",
    minWidth: "220",
}, ...__VLS_functionalComponentArgsRest(__VLS_30));
const __VLS_33 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_34 = __VLS_asFunctionalComponent(__VLS_33, new __VLS_33({
    label: "角色",
    minWidth: "150",
}));
const __VLS_35 = __VLS_34({
    label: "角色",
    minWidth: "150",
}, ...__VLS_functionalComponentArgsRest(__VLS_34));
__VLS_36.slots.default;
{
    const { default: __VLS_thisSlot } = __VLS_36.slots;
    const [scope] = __VLS_getSlotParams(__VLS_thisSlot);
    ((scope.row.roles || []).join(', '));
}
var __VLS_36;
const __VLS_37 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_38 = __VLS_asFunctionalComponent(__VLS_37, new __VLS_37({
    prop: "status",
    label: "状态",
    width: "120",
}));
const __VLS_39 = __VLS_38({
    prop: "status",
    label: "状态",
    width: "120",
}, ...__VLS_functionalComponentArgsRest(__VLS_38));
const __VLS_41 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_42 = __VLS_asFunctionalComponent(__VLS_41, new __VLS_41({
    prop: "pod_count",
    label: "Pod数",
    width: "90",
}));
const __VLS_43 = __VLS_42({
    prop: "pod_count",
    label: "Pod数",
    width: "90",
}, ...__VLS_functionalComponentArgsRest(__VLS_42));
const __VLS_45 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_46 = __VLS_asFunctionalComponent(__VLS_45, new __VLS_45({
    label: "可调度",
    width: "100",
}));
const __VLS_47 = __VLS_46({
    label: "可调度",
    width: "100",
}, ...__VLS_functionalComponentArgsRest(__VLS_46));
__VLS_48.slots.default;
{
    const { default: __VLS_thisSlot } = __VLS_48.slots;
    const [scope] = __VLS_getSlotParams(__VLS_thisSlot);
    (scope.row.unschedulable ? '否' : '是');
}
var __VLS_48;
const __VLS_49 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_50 = __VLS_asFunctionalComponent(__VLS_49, new __VLS_49({
    label: "操作",
    width: "260",
    fixed: "right",
}));
const __VLS_51 = __VLS_50({
    label: "操作",
    width: "260",
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
        type: "warning",
    }));
    const __VLS_55 = __VLS_54({
        ...{ 'onClick': {} },
        link: true,
        type: "warning",
    }, ...__VLS_functionalComponentArgsRest(__VLS_54));
    let __VLS_57;
    let __VLS_58;
    let __VLS_59;
    const __VLS_60 = {
        onClick: (...[$event]) => {
            __VLS_ctx.handleCordon(scope.row.name);
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('k8s:node:cordon') }, null, null);
    __VLS_56.slots.default;
    var __VLS_56;
    const __VLS_61 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_62 = __VLS_asFunctionalComponent(__VLS_61, new __VLS_61({
        ...{ 'onClick': {} },
        link: true,
        type: "success",
    }));
    const __VLS_63 = __VLS_62({
        ...{ 'onClick': {} },
        link: true,
        type: "success",
    }, ...__VLS_functionalComponentArgsRest(__VLS_62));
    let __VLS_65;
    let __VLS_66;
    let __VLS_67;
    const __VLS_68 = {
        onClick: (...[$event]) => {
            __VLS_ctx.handleUncordon(scope.row.name);
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('k8s:node:uncordon') }, null, null);
    __VLS_64.slots.default;
    var __VLS_64;
    const __VLS_69 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_70 = __VLS_asFunctionalComponent(__VLS_69, new __VLS_69({
        ...{ 'onClick': {} },
        link: true,
        type: "danger",
    }));
    const __VLS_71 = __VLS_70({
        ...{ 'onClick': {} },
        link: true,
        type: "danger",
    }, ...__VLS_functionalComponentArgsRest(__VLS_70));
    let __VLS_73;
    let __VLS_74;
    let __VLS_75;
    const __VLS_76 = {
        onClick: (...[$event]) => {
            __VLS_ctx.handleDrain(scope.row.name);
        }
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('k8s:node:drain') }, null, null);
    __VLS_72.slots.default;
    var __VLS_72;
}
var __VLS_52;
var __VLS_28;
var __VLS_3;
/** @type {__VLS_StyleScopedClasses['toolbar']} */ ;
/** @type {__VLS_StyleScopedClasses['actions']} */ ;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            loading: loading,
            list: list,
            clusterOptions: clusterOptions,
            clusterId: clusterId,
            loadData: loadData,
            handleCordon: handleCordon,
            handleUncordon: handleUncordon,
            handleDrain: handleDrain,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
