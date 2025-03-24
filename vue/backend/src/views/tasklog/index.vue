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
import { TaskLog } from '@/types/task-log';
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
    request_id: '',
});
const searchOpt = ref<FormOptionList[]>([
    { type: 'input', label: '开始时间：', prop: 'start' },
    { type: 'input', label: '结束时间：', prop: 'end' },
    { type: 'input', label: 'REQUEST_ID：', prop: 'request_id' },
])
const handleSearch = () => {
    changePage(1);
};

// 表格相关
let columns = ref([
    { prop: 'id', label: 'ID', width: 55, align: 'center' },
    { prop: 'request_id', label: '请求ID' },
    { prop: 'action_string', label: '动作' },
    { prop: 'created_at', label: '创建时间' },
    { prop: 'operator', label: '操作', width: 250, buttons: ['view'] },
])
const page = reactive({
    index: 1,
    size: 10,
    total: 0,
})
const tableData = ref<TaskLog[]>([]);
const getData = async () => {
    simpleApi.get('/api/tasklog/getList', { page: page.index, pageSize: page.size, start: query.start, end: query.end, request_id: query.request_id }, permiss.token, function(data){
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
const handleView = (row: TaskLog) => {
    simpleApi.get('/api/tasklog/getDetail', { id: row.id }, permiss.token, function(data){
        viewData.value.row = data;
        viewData.value.list = [
          {
            prop: 'id',
            label: 'ID',
          },
          {
            prop: 'app_name',
            label: '应用',
          },
          {
            prop: 'task_name',
            label: '任务',
          },
          {
            prop: 'task_queue_id',
            label: '队列ID',
          },
          {
            prop: 'request_id',
            label: '请求ID',
          },
          {
            prop: 'action_string',
            label: '动作',
          },
          {
            prop: 'message',
            label: '信息',
          },
          {
            prop: 'created_at',
            label: '创建时间',
          },
        ]
        visible1.value = true;
    });
};

</script>

<style scoped></style>