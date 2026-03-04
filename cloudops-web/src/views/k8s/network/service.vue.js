import { onMounted, ref } from 'vue';
import { listK8sClustersApi, listK8sNamespacesApi, listServicesApi } from '@/api/k8s';
const loading = ref(false);
const clusters = ref([]);
const namespaces = ref([]);
const clusterId = ref();
const namespace = ref('default');
const list = ref([]);
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
        const res = await listServicesApi(clusterId.value, namespace.value);
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
}
const __VLS_29 = {}.ElTable;
/** @type {[typeof __VLS_components.ElTable, typeof __VLS_components.elTable, typeof __VLS_components.ElTable, typeof __VLS_components.elTable, ]} */ ;
// @ts-ignore
const __VLS_30 = __VLS_asFunctionalComponent(__VLS_29, new __VLS_29({
    data: (__VLS_ctx.list),
}));
const __VLS_31 = __VLS_30({
    data: (__VLS_ctx.list),
}, ...__VLS_functionalComponentArgsRest(__VLS_30));
__VLS_asFunctionalDirective(__VLS_directives.vLoading)(null, { ...__VLS_directiveBindingRestFields, value: (__VLS_ctx.loading) }, null, null);
__VLS_32.slots.default;
const __VLS_33 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_34 = __VLS_asFunctionalComponent(__VLS_33, new __VLS_33({
    prop: "metadata.name",
    label: "名称",
    minWidth: "180",
}));
const __VLS_35 = __VLS_34({
    prop: "metadata.name",
    label: "名称",
    minWidth: "180",
}, ...__VLS_functionalComponentArgsRest(__VLS_34));
const __VLS_37 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_38 = __VLS_asFunctionalComponent(__VLS_37, new __VLS_37({
    prop: "spec.type",
    label: "类型",
    width: "120",
}));
const __VLS_39 = __VLS_38({
    prop: "spec.type",
    label: "类型",
    width: "120",
}, ...__VLS_functionalComponentArgsRest(__VLS_38));
const __VLS_41 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_42 = __VLS_asFunctionalComponent(__VLS_41, new __VLS_41({
    label: "ClusterIP",
    minWidth: "150",
}));
const __VLS_43 = __VLS_42({
    label: "ClusterIP",
    minWidth: "150",
}, ...__VLS_functionalComponentArgsRest(__VLS_42));
__VLS_44.slots.default;
{
    const { default: __VLS_thisSlot } = __VLS_44.slots;
    const [scope] = __VLS_getSlotParams(__VLS_thisSlot);
    (scope.row.spec?.clusterIP || '-');
}
var __VLS_44;
const __VLS_45 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_46 = __VLS_asFunctionalComponent(__VLS_45, new __VLS_45({
    label: "端口",
    minWidth: "200",
}));
const __VLS_47 = __VLS_46({
    label: "端口",
    minWidth: "200",
}, ...__VLS_functionalComponentArgsRest(__VLS_46));
__VLS_48.slots.default;
{
    const { default: __VLS_thisSlot } = __VLS_48.slots;
    const [scope] = __VLS_getSlotParams(__VLS_thisSlot);
    ((scope.row.spec?.ports || []).map((p) => `${p.port}/${p.protocol}`).join(', '));
}
var __VLS_48;
var __VLS_32;
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
            loadData: loadData,
            handleClusterChange: handleClusterChange,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
