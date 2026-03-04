import { onMounted, reactive, ref } from 'vue';
import { listLoginLogsApi } from '@/api/audit';
const loading = ref(false);
const total = ref(0);
const list = ref([]);
const query = reactive({ page: 1, page_size: 10, keyword: '' });
const loadData = async () => {
    loading.value = true;
    try {
        const res = await listLoginLogsApi(query);
        list.value = res.list || [];
        total.value = res.total || 0;
    }
    finally {
        loading.value = false;
    }
};
const handlePageChange = (page) => {
    query.page = page;
    loadData();
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
}
const __VLS_5 = {}.ElTable;
/** @type {[typeof __VLS_components.ElTable, typeof __VLS_components.elTable, typeof __VLS_components.ElTable, typeof __VLS_components.elTable, ]} */ ;
// @ts-ignore
const __VLS_6 = __VLS_asFunctionalComponent(__VLS_5, new __VLS_5({
    data: (__VLS_ctx.list),
    size: "small",
}));
const __VLS_7 = __VLS_6({
    data: (__VLS_ctx.list),
    size: "small",
}, ...__VLS_functionalComponentArgsRest(__VLS_6));
__VLS_asFunctionalDirective(__VLS_directives.vLoading)(null, { ...__VLS_directiveBindingRestFields, value: (__VLS_ctx.loading) }, null, null);
__VLS_8.slots.default;
const __VLS_9 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_10 = __VLS_asFunctionalComponent(__VLS_9, new __VLS_9({
    prop: "id",
    label: "ID",
    width: "80",
}));
const __VLS_11 = __VLS_10({
    prop: "id",
    label: "ID",
    width: "80",
}, ...__VLS_functionalComponentArgsRest(__VLS_10));
const __VLS_13 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_14 = __VLS_asFunctionalComponent(__VLS_13, new __VLS_13({
    prop: "username",
    label: "用户",
    width: "120",
}));
const __VLS_15 = __VLS_14({
    prop: "username",
    label: "用户",
    width: "120",
}, ...__VLS_functionalComponentArgsRest(__VLS_14));
const __VLS_17 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_18 = __VLS_asFunctionalComponent(__VLS_17, new __VLS_17({
    prop: "ip",
    label: "IP",
    width: "140",
}));
const __VLS_19 = __VLS_18({
    prop: "ip",
    label: "IP",
    width: "140",
}, ...__VLS_functionalComponentArgsRest(__VLS_18));
const __VLS_21 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_22 = __VLS_asFunctionalComponent(__VLS_21, new __VLS_21({
    prop: "method",
    label: "方法",
    width: "90",
}));
const __VLS_23 = __VLS_22({
    prop: "method",
    label: "方法",
    width: "90",
}, ...__VLS_functionalComponentArgsRest(__VLS_22));
const __VLS_25 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_26 = __VLS_asFunctionalComponent(__VLS_25, new __VLS_25({
    prop: "path",
    label: "路径",
    minWidth: "220",
}));
const __VLS_27 = __VLS_26({
    prop: "path",
    label: "路径",
    minWidth: "220",
}, ...__VLS_functionalComponentArgsRest(__VLS_26));
const __VLS_29 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_30 = __VLS_asFunctionalComponent(__VLS_29, new __VLS_29({
    prop: "response_code",
    label: "响应码",
    width: "90",
}));
const __VLS_31 = __VLS_30({
    prop: "response_code",
    label: "响应码",
    width: "90",
}, ...__VLS_functionalComponentArgsRest(__VLS_30));
const __VLS_33 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_34 = __VLS_asFunctionalComponent(__VLS_33, new __VLS_33({
    prop: "created_at",
    label: "时间",
    width: "180",
}));
const __VLS_35 = __VLS_34({
    prop: "created_at",
    label: "时间",
    width: "180",
}, ...__VLS_functionalComponentArgsRest(__VLS_34));
var __VLS_8;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "pager" },
});
const __VLS_37 = {}.ElPagination;
/** @type {[typeof __VLS_components.ElPagination, typeof __VLS_components.elPagination, ]} */ ;
// @ts-ignore
const __VLS_38 = __VLS_asFunctionalComponent(__VLS_37, new __VLS_37({
    ...{ 'onCurrentChange': {} },
    background: true,
    layout: "total, prev, pager, next",
    total: (__VLS_ctx.total),
    currentPage: (__VLS_ctx.query.page),
    pageSize: (__VLS_ctx.query.page_size),
}));
const __VLS_39 = __VLS_38({
    ...{ 'onCurrentChange': {} },
    background: true,
    layout: "total, prev, pager, next",
    total: (__VLS_ctx.total),
    currentPage: (__VLS_ctx.query.page),
    pageSize: (__VLS_ctx.query.page_size),
}, ...__VLS_functionalComponentArgsRest(__VLS_38));
let __VLS_41;
let __VLS_42;
let __VLS_43;
const __VLS_44 = {
    onCurrentChange: (__VLS_ctx.handlePageChange)
};
var __VLS_40;
var __VLS_3;
/** @type {__VLS_StyleScopedClasses['pager']} */ ;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            loading: loading,
            total: total,
            list: list,
            query: query,
            handlePageChange: handlePageChange,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
