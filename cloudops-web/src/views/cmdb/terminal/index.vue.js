import { nextTick, onBeforeUnmount, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { ElMessage } from 'element-plus';
import { listCMDBSSHRecordsApi } from '@/api/cmdb';
const route = useRoute();
const terminalBox = ref();
const hostID = ref(Number(route.query.host_id || 1));
const command = ref('');
const output = ref('');
const ws = ref(null);
const records = ref([]);
const total = ref(0);
const query = ref({ page: 1, page_size: 10 });
const appendOutput = async (text) => {
    output.value += text;
    await nextTick();
    if (terminalBox.value) {
        terminalBox.value.scrollTop = terminalBox.value.scrollHeight;
    }
};
const connect = () => {
    disconnect();
    const token = localStorage.getItem('opsnexus_access_token') || '';
    if (!token) {
        ElMessage.warning('请先登录');
        return;
    }
    const protocol = location.protocol === 'https:' ? 'wss' : 'ws';
    const wsURL = `${protocol}://${location.host}/ws/cmdb/terminal/${hostID.value}?token=${encodeURIComponent(token)}`;
    ws.value = new WebSocket(wsURL);
    ws.value.onopen = () => appendOutput(`\r\n[系统] 已连接主机 ${hostID.value}\r\n`);
    ws.value.onmessage = (ev) => appendOutput(String(ev.data || ''));
    ws.value.onerror = () => appendOutput('\r\n[系统] 连接异常\r\n');
    ws.value.onclose = () => appendOutput('\r\n[系统] 连接关闭\r\n');
};
const disconnect = () => {
    if (ws.value) {
        ws.value.close();
        ws.value = null;
    }
};
const sendCommand = () => {
    if (!ws.value || ws.value.readyState !== WebSocket.OPEN) {
        ElMessage.warning('终端未连接');
        return;
    }
    const text = `${command.value}\n`;
    ws.value.send(text);
    command.value = '';
};
const loadRecords = async () => {
    const res = await listCMDBSSHRecordsApi(query.value);
    records.value = res.list || [];
    total.value = res.total || 0;
};
const handlePageChange = (page) => {
    query.value.page = page;
    loadRecords();
};
onMounted(loadRecords);
onBeforeUnmount(disconnect);
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
// CSS variable injection 
// CSS variable injection end 
const __VLS_0 = {}.ElRow;
/** @type {[typeof __VLS_components.ElRow, typeof __VLS_components.elRow, typeof __VLS_components.ElRow, typeof __VLS_components.elRow, ]} */ ;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
    gutter: (12),
}));
const __VLS_2 = __VLS_1({
    gutter: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
var __VLS_4 = {};
__VLS_3.slots.default;
const __VLS_5 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_6 = __VLS_asFunctionalComponent(__VLS_5, new __VLS_5({
    span: (16),
}));
const __VLS_7 = __VLS_6({
    span: (16),
}, ...__VLS_functionalComponentArgsRest(__VLS_6));
__VLS_8.slots.default;
const __VLS_9 = {}.ElCard;
/** @type {[typeof __VLS_components.ElCard, typeof __VLS_components.elCard, typeof __VLS_components.ElCard, typeof __VLS_components.elCard, ]} */ ;
// @ts-ignore
const __VLS_10 = __VLS_asFunctionalComponent(__VLS_9, new __VLS_9({}));
const __VLS_11 = __VLS_10({}, ...__VLS_functionalComponentArgsRest(__VLS_10));
__VLS_12.slots.default;
{
    const { header: __VLS_thisSlot } = __VLS_12.slots;
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ class: "toolbar" },
    });
    __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({});
    const __VLS_13 = {}.ElInputNumber;
    /** @type {[typeof __VLS_components.ElInputNumber, typeof __VLS_components.elInputNumber, ]} */ ;
    // @ts-ignore
    const __VLS_14 = __VLS_asFunctionalComponent(__VLS_13, new __VLS_13({
        modelValue: (__VLS_ctx.hostID),
        min: (1),
        controlsPosition: "right",
        ...{ style: {} },
    }));
    const __VLS_15 = __VLS_14({
        modelValue: (__VLS_ctx.hostID),
        min: (1),
        controlsPosition: "right",
        ...{ style: {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_14));
    const __VLS_17 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_18 = __VLS_asFunctionalComponent(__VLS_17, new __VLS_17({
        ...{ 'onClick': {} },
        type: "primary",
    }));
    const __VLS_19 = __VLS_18({
        ...{ 'onClick': {} },
        type: "primary",
    }, ...__VLS_functionalComponentArgsRest(__VLS_18));
    let __VLS_21;
    let __VLS_22;
    let __VLS_23;
    const __VLS_24 = {
        onClick: (__VLS_ctx.connect)
    };
    __VLS_asFunctionalDirective(__VLS_directives.vPermission)(null, { ...__VLS_directiveBindingRestFields, value: ('cmdb:host:terminal') }, null, null);
    __VLS_20.slots.default;
    var __VLS_20;
    const __VLS_25 = {}.ElButton;
    /** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
    // @ts-ignore
    const __VLS_26 = __VLS_asFunctionalComponent(__VLS_25, new __VLS_25({
        ...{ 'onClick': {} },
    }));
    const __VLS_27 = __VLS_26({
        ...{ 'onClick': {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_26));
    let __VLS_29;
    let __VLS_30;
    let __VLS_31;
    const __VLS_32 = {
        onClick: (__VLS_ctx.disconnect)
    };
    __VLS_28.slots.default;
    var __VLS_28;
}
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ref: "terminalBox",
    ...{ class: "terminal-output" },
});
/** @type {typeof __VLS_ctx.terminalBox} */ ;
(__VLS_ctx.output);
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "input-row" },
});
const __VLS_33 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
const __VLS_34 = __VLS_asFunctionalComponent(__VLS_33, new __VLS_33({
    ...{ 'onKeyup': {} },
    modelValue: (__VLS_ctx.command),
    placeholder: "输入命令并回车",
}));
const __VLS_35 = __VLS_34({
    ...{ 'onKeyup': {} },
    modelValue: (__VLS_ctx.command),
    placeholder: "输入命令并回车",
}, ...__VLS_functionalComponentArgsRest(__VLS_34));
let __VLS_37;
let __VLS_38;
let __VLS_39;
const __VLS_40 = {
    onKeyup: (__VLS_ctx.sendCommand)
};
var __VLS_36;
const __VLS_41 = {}.ElButton;
/** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
// @ts-ignore
const __VLS_42 = __VLS_asFunctionalComponent(__VLS_41, new __VLS_41({
    ...{ 'onClick': {} },
    type: "primary",
}));
const __VLS_43 = __VLS_42({
    ...{ 'onClick': {} },
    type: "primary",
}, ...__VLS_functionalComponentArgsRest(__VLS_42));
let __VLS_45;
let __VLS_46;
let __VLS_47;
const __VLS_48 = {
    onClick: (__VLS_ctx.sendCommand)
};
__VLS_44.slots.default;
var __VLS_44;
var __VLS_12;
var __VLS_8;
const __VLS_49 = {}.ElCol;
/** @type {[typeof __VLS_components.ElCol, typeof __VLS_components.elCol, typeof __VLS_components.ElCol, typeof __VLS_components.elCol, ]} */ ;
// @ts-ignore
const __VLS_50 = __VLS_asFunctionalComponent(__VLS_49, new __VLS_49({
    span: (8),
}));
const __VLS_51 = __VLS_50({
    span: (8),
}, ...__VLS_functionalComponentArgsRest(__VLS_50));
__VLS_52.slots.default;
const __VLS_53 = {}.ElCard;
/** @type {[typeof __VLS_components.ElCard, typeof __VLS_components.elCard, typeof __VLS_components.ElCard, typeof __VLS_components.elCard, ]} */ ;
// @ts-ignore
const __VLS_54 = __VLS_asFunctionalComponent(__VLS_53, new __VLS_53({}));
const __VLS_55 = __VLS_54({}, ...__VLS_functionalComponentArgsRest(__VLS_54));
__VLS_56.slots.default;
{
    const { header: __VLS_thisSlot } = __VLS_56.slots;
    __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
}
const __VLS_57 = {}.ElTable;
/** @type {[typeof __VLS_components.ElTable, typeof __VLS_components.elTable, typeof __VLS_components.ElTable, typeof __VLS_components.elTable, ]} */ ;
// @ts-ignore
const __VLS_58 = __VLS_asFunctionalComponent(__VLS_57, new __VLS_57({
    data: (__VLS_ctx.records),
    size: "small",
    height: "520",
}));
const __VLS_59 = __VLS_58({
    data: (__VLS_ctx.records),
    size: "small",
    height: "520",
}, ...__VLS_functionalComponentArgsRest(__VLS_58));
__VLS_60.slots.default;
const __VLS_61 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_62 = __VLS_asFunctionalComponent(__VLS_61, new __VLS_61({
    prop: "host_id",
    label: "主机",
    width: "72",
}));
const __VLS_63 = __VLS_62({
    prop: "host_id",
    label: "主机",
    width: "72",
}, ...__VLS_functionalComponentArgsRest(__VLS_62));
const __VLS_65 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_66 = __VLS_asFunctionalComponent(__VLS_65, new __VLS_65({
    prop: "session_id",
    label: "会话ID",
    minWidth: "140",
    showOverflowTooltip: true,
}));
const __VLS_67 = __VLS_66({
    prop: "session_id",
    label: "会话ID",
    minWidth: "140",
    showOverflowTooltip: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_66));
const __VLS_69 = {}.ElTableColumn;
/** @type {[typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ]} */ ;
// @ts-ignore
const __VLS_70 = __VLS_asFunctionalComponent(__VLS_69, new __VLS_69({
    prop: "start_time",
    label: "开始时间",
    minWidth: "160",
}));
const __VLS_71 = __VLS_70({
    prop: "start_time",
    label: "开始时间",
    minWidth: "160",
}, ...__VLS_functionalComponentArgsRest(__VLS_70));
var __VLS_60;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "pager" },
});
const __VLS_73 = {}.ElPagination;
/** @type {[typeof __VLS_components.ElPagination, typeof __VLS_components.elPagination, ]} */ ;
// @ts-ignore
const __VLS_74 = __VLS_asFunctionalComponent(__VLS_73, new __VLS_73({
    ...{ 'onCurrentChange': {} },
    background: true,
    small: true,
    layout: "prev, pager, next",
    total: (__VLS_ctx.total),
    currentPage: (__VLS_ctx.query.page),
    pageSize: (__VLS_ctx.query.page_size),
}));
const __VLS_75 = __VLS_74({
    ...{ 'onCurrentChange': {} },
    background: true,
    small: true,
    layout: "prev, pager, next",
    total: (__VLS_ctx.total),
    currentPage: (__VLS_ctx.query.page),
    pageSize: (__VLS_ctx.query.page_size),
}, ...__VLS_functionalComponentArgsRest(__VLS_74));
let __VLS_77;
let __VLS_78;
let __VLS_79;
const __VLS_80 = {
    onCurrentChange: (__VLS_ctx.handlePageChange)
};
var __VLS_76;
var __VLS_56;
var __VLS_52;
var __VLS_3;
/** @type {__VLS_StyleScopedClasses['toolbar']} */ ;
/** @type {__VLS_StyleScopedClasses['terminal-output']} */ ;
/** @type {__VLS_StyleScopedClasses['input-row']} */ ;
/** @type {__VLS_StyleScopedClasses['pager']} */ ;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            terminalBox: terminalBox,
            hostID: hostID,
            command: command,
            output: output,
            records: records,
            total: total,
            query: query,
            connect: connect,
            disconnect: disconnect,
            sendCommand: sendCommand,
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
