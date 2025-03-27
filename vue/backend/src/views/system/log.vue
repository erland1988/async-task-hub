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
import { Log } from '@/types/log';
import { User } from "@/types/user";
import {simpleApi} from '@/api';
import TableCustom from '@/components/table-custom.vue';
import TableDetail from '@/components/table-detail.vue';
import TableSearch from '@/components/table-search.vue';
import { FormOptionList } from '@/types/form-option';
import {usePermissStore} from "@/store/permiss";

const permiss = usePermissStore();
// 查询相关
const query = reactive({
    keywords: '',
});
const searchOpt = ref<FormOptionList[]>([
    { type: 'input', label: '用户：', prop: 'keywords' }
])
const handleSearch = () => {
    changePage(1);
};

// 表格相关
let columns = ref([
    { prop: 'id', label: 'ID', width: 55, align: 'center' },
    { prop: 'admin', label: '用户', formatter: (row: User) => row.username?row.username:'-' },
    { prop: 'operation', label: '操作' },
    { prop: 'created_at', label: '创建时间' },
    { prop: 'operator', buttons: ['view'], label: '操作', width: 250 },
])
const page = reactive({
    index: 1,
    size: 10,
    total: 0,
})
const tableData = ref<Log[]>([]);
const getData = async () => {
    simpleApi.get('/task/api/log/getList', { page: page.index, pageSize: page.size, keywords: query.keywords }, permiss.token, function(data){
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
const handleView = (row: Log) => {
    simpleApi.get('/task/api/log/getDetail', { id: row.id }, permiss.token, function(data) {
      viewData.value.row = data;
      viewData.value.list = [
        {
          prop: 'id',
          label: 'ID',
        },
        {
          prop: 'admin',
          label: '用户',
          formatter: () => row.admin.username?row.admin.username:'-',
        },
        {
            prop: 'operation',
            label: '操作',
        },
        {
            prop: 'created_at',
            label: '创建时间',
        },
        {
          prop: 'details',
          label: '内容',
        },
      ]
      visible1.value = true;
    })
};
</script>

<style scoped></style>