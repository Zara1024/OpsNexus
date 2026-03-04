import { onMounted, ref } from 'vue';
import { listK8sClustersApi, listK8sNamespacesApi } from '@/api/k8s';
const loading = ref(false);
const clusters = ref([]);
const clusterId = ref();
const list = ref([]);
const loadClusters = async () => {
    const res = await listK8sClustersApi();
    clusters.value = res.list || [];
    if (!clusterId.value && clusters.value.length) {
        clusterId.value = clusters.value[0].id;
    }
};
const loadData = async () => {
    if (!clusterId.value)
        return;
    loading.value = true;
    try {
        const res = await listK8sNamespacesApi(clusterId.value);
        list.value = (res.list || []).map((name) => ({ name }));
    }
    finally {
        loading.value = false;
    }
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
}
const __VLS_17 = {}.ElTable;
/** @type {[typeof __VLS_components.ElTable, typeof __VLS_components.elTable, typeof __VLS_components.ElTable, typeof __VLS_components.elTable, ]} */ ;
// @ts-ignore
const __VLS_18 = __VLS_asFunctionalComponent(__VLS_17, new __VLS_17({
    data: (__VLS_ctx.list),
}));
const __VLS_19 = __VLS_18({
    data: (__VLS_ctx.list),
}, ...__VLS_functionalComponentArgsRest(__VLS_18));
__VLS_asFunctionalDirective(__VLS_directives.vLoading)(null, { ...__VLS_directiveBindingRestFields, value: (__VLS_ctx.loading) }, null, null);
__VLS_20.slots.default;
const __VLS_21 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_22 = __VLS_asFunctionalComponent(__VLS_21, new __VLS_21({
    prop: "name",
    label: "命名空间",
    minWidth: "220",
}));
const __VLS_23 = __VLS_22({
    prop: "name",
    label: "命名空间",
    minWidth: "220",
}, ...__VLS_functionalComponentArgsRest(__VLS_22));
var __VLS_20;
var __VLS_3;
/** @type {__VLS_StyleScopedClasses['toolbar']} */ ;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            loading: loading,
            clusters: clusters,
            clusterId: clusterId,
            list: list,
            loadData: loadData,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
