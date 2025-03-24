<template>
    <div>
        <TableSearch :query="query" :options="searchOpt" :search="handleSearch" />
        <div class="container">
            <TableCustom :columns="columns" :tableData="tableData" :total="page.total" :currentPage="page.index" :changePage="changePage" :viewFunc="handleView" :delFunc="handleDelete" :editFunc="handleEdit">
                <template #toolbarBtn>
                    <el-button type="warning" :icon="CirclePlusFilled" @click="visible = true">新增</el-button>
                </template>
            </TableCustom>
        </div>
        <el-dialog :title="isEdit ? '编辑' : '新增'" v-model="visible" width="700px" destroy-on-close
            :close-on-click-modal="false" @close="closeDialog">
            <TableEdit :form-data="rowData" :options="options" :edit="isEdit" :update="updateData" />
        </el-dialog>
        <el-dialog title="查看详情" v-model="visible1" width="700px" destroy-on-close>
            <TableDetail :data="viewData"></TableDetail>
        </el-dialog>
    </div>
</template>

<script setup lang="ts" name="system-user">
import { ref, reactive } from 'vue';
import { ElMessage } from 'element-plus';
import { CirclePlusFilled } from '@element-plus/icons-vue';
import { App } from '@/types/app';
import {simpleApi} from '@/api';
import TableCustom from '@/components/table-custom.vue';
import TableDetail from '@/components/table-detail.vue';
import TableSearch from '@/components/table-search.vue';
import { FormOption, FormOptionList } from '@/types/form-option';
import {usePermissStore} from "@/store/permiss";
import {User} from "@/types/user";

const permiss = usePermissStore();
// 查询相关
const query = reactive({
    keywords: '',
});
const searchOpt = ref<FormOptionList[]>([
    { type: 'input', label: '应用：', prop: 'keywords' }
])
const handleSearch = () => {
    changePage(1);
};

// 表格相关
let columns = ref([
    { prop: 'id', label: 'ID', width: 55, align: 'center' },
    { prop: 'name', label: '应用名称' },
    { prop: 'username', label: '用户' },
    { prop: 'app_key', label: '应用标识' },
    { prop: 'created_at', label: '创建时间' },
    { prop: 'operator', label: '操作', width: 250 },
])
const page = reactive({
    index: 1,
    size: 10,
    total: 0,
})
const tableData = ref<App[]>([]);
const getData = async () => {
    simpleApi.get('/api/app/getList', { page: page.index, pageSize: page.size, keywords: query.keywords }, permiss.token, function(data){
      tableData.value = data.list;
      page.total = data.total;
    });
};
getData();

const changePage = (val: number) => {
    page.index = val;
    getData();
};

// 新增/编辑弹窗相关
let options = ref<FormOption>({
    labelWidth: '100px',
    span: 12,
    list: [
        { type: 'input', label: '应用名称', prop: 'name', required: true },
        { type: 'input', label: '应用标识', prop: 'app_key', required: true },
        { type: 'input', label: '应用秘钥', prop: 'app_secret', required: true },
        { type: 'input', label: '备注', prop: 'remark', required: true },
    ]
})
const visible = ref(false);
const isEdit = ref(false);
const rowData = ref({});
const handleEdit = (row: App) => {
    simpleApi.get('/api/app/getDetail', { id: row.id }, permiss.token, function(data) {
      rowData.value = data;
      isEdit.value = true;
      visible.value = true;
    })
};
const updateData = (row: App) => {
    if(isEdit.value){
      const params = {
        id: row.id,
        name: row.name,
        app_key: row.app_key,
        app_secret: row.app_secret,
        remark: row.remark,
      }
      simpleApi.post('/api/app/update', params, permiss.token, function(data) {
        closeDialog();
        getData();
      })
    }else{
      const params = {
        name: row.name,
        app_key: row.app_key,
        app_secret: row.app_secret,
        remark: row.remark,
      }
      simpleApi.post('/api/app/create', params, permiss.token, function(data) {
        closeDialog();
        getData();
      })
    }
};

const closeDialog = () => {
    visible.value = false;
    isEdit.value = false;
};

// 查看详情弹窗相关
const visible1 = ref(false);
const viewData = ref({
    row: {},
    list: []
});
const handleView = (row: App) => {
    simpleApi.get('/api/app/getDetail', { id: row.id }, permiss.token, function(data){
      viewData.value.row = data;
      viewData.value.list = [
        {
          prop: 'id',
          label: 'ID',
        },
        {
          prop: 'name',
          label: '应用名称',
        },
        {
          prop: 'app_key',
          label: '应用标识',
        },
        {
          prop: 'app_secret',
          label: '应用秘钥',
        },
        {
          prop: 'created_at',
          label: '创建时间',
        },
        {
          prop: 'updated_at',
          label: '更新时间',
        },
        {
          prop: 'remark',
          label: '备注',
        },
      ]
      visible1.value = true;
    })
};

// 删除相关
const handleDelete = (row: App) => {
    simpleApi.postForm('/api/app/delete', { id: row.id }, permiss.token, function(data){
      ElMessage.success('删除成功');
      getData();
    })
}
</script>

<style scoped></style>