<template>
    <div>
        <TableSearch :query="query" :options="searchOpt" :search="handleSearch" />
        <div class="container">
            <TableCustom :columns="columns" :tableData="tableData" :total="page.total" :currentPage="page.index" :changePage="changePage" :viewFunc="handleView">
            </TableCustom>
        </div>
        <el-dialog title="查看详情" v-model="visible1" width="700px" destroy-on-close>
            <TableDetail :data="viewData"></TableDetail>
        </el-dialog>
    </div>
</template>

<script setup lang="ts" name="system-user">
import { ref, reactive } from 'vue';
import { TaskQueue } from '@/types/task-queue';
import {simpleApi} from '@/api';
import TableCustom from '@/components/table-custom.vue';
import TableDetail from '@/components/table-detail.vue';
import TableSearch from '@/components/table-search.vue';
import { FormOptionList } from '@/types/form-option';
import {usePermissStore} from "@/store/permiss";

const permiss = usePermissStore();
// 查询相关
const query = reactive({
  start: '',
  end: '',
});
const searchOpt = ref<FormOptionList[]>([
    { type: 'input', label: '开始时间：', prop: 'start' },
    { type: 'input', label: '结束时间：', prop: 'end' },
])
const handleSearch = () => {
    changePage(1);
};

// 表格相关
let columns = ref([
    { prop: 'id', label: 'ID', width: 55, align: 'center' },
    { prop: 'relative_delay_time', label: '延迟时间' },
    { prop: 'delay_execution_time', label: '绝对时间' },
    { prop: 'execution_status_string', label: '状态' },
    { prop: 'execution_duration', label: '执行时长(毫秒)' },
    { prop: 'execution_count', label: '执行次数' },
    { prop: 'created_at', label: '创建时间' },
    { prop: 'operator', label: '操作', width: 250, buttons: ['view'] },
])
const page = reactive({
    index: 1,
    size: 10,
    total: 0,
})
const tableData = ref<TaskQueue[]>([]);
const getData = async () => {
    simpleApi.get('/api/taskqueue/getList', { page: page.index, pageSize: page.size, start: query.start, end: query.end }, permiss.token, function(data){
      tableData.value = data.list;
      page.total = data.total;
    });
};
getData();

const changePage = (val: number) => {
    page.index = val;
    getData();
};

// 查看详情弹窗相关
const visible1 = ref(false);
const viewData = ref({
    row: {},
    list: []
});
const handleView = (row: TaskQueue) => {
    simpleApi.get('/api/taskqueue/getDetail', { id: row.id }, permiss.token, function(data){
      viewData.value.row = data;
      viewData.value.list = [
        {
            prop: 'id',
            label: 'ID',
        },
        {
            prop: 'appname',
            label: '应用',
        },
        {
          prop: 'taskname',
          label: '任务',
        },
        {
          prop: 'executor_url',
          label: '执行器URL',
        },
        {
            prop: 'parameters',
            label: '参数',
        },
        {
            prop:'relative_delay_time',
            label: '延迟时间',
        },
        {
            prop: 'delay_execution_time',
            label: '绝对时间',
        },
        {
            prop: 'execution_status_string',
            label: '状态',
        },
        {
          prop: 'execution_start',
          label: '执行开始时间',
        },
        {
          prop: 'execution_end',
          label: '执行结束时间',
        },
        {
          prop: 'execution_duration',
          label: '执行时长(毫秒)',
        },
        {
          prop: 'execution_count',
          label: '执行次数',
        },
        {
          prop: 'created_at',
          label: '创建时间',
        },
        {
          prop: 'updated_at',
          label: '更新时间',
        },
      ]
      visible1.value = true;
    });
};
</script>

<style scoped></style>