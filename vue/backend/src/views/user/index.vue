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

<script setup lang="ts" name="user-index">
import { ref, reactive } from 'vue';
import {ElMessage, FormRules} from 'element-plus';
import { CirclePlusFilled } from '@element-plus/icons-vue';
import {User} from '@/types/user';
import {simpleApi} from '@/api';
import TableCustom from '@/components/table-custom.vue';
import TableDetail from '@/components/table-detail.vue';
import TableSearch from '@/components/table-search.vue';
import { FormOption, FormOptionList } from '@/types/form-option';
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
    { prop: 'username', label: '用户名' },
    { prop: 'truename', label: '真实姓名' },
    { prop: 'rolename', label: '角色' },
    { prop: 'expires_at', label: '到期时间' },
    { prop: 'created_at', label: '创建时间', },
    { prop: 'operator', label: '操作', width: 250 },
])
const page = reactive({
    index: 1,
    size: 10,
    total: 0,
})
const tableData = ref<User[]>([]);
const getData = async () => {
    simpleApi.get('/api/admin/getList', { page: page.index, pageSize: page.size, keywords: query.keywords }, permiss.token, function(data){
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
        { type: 'input', label: '用户名', prop: 'username', required: true },
        { type: 'input', label: '密码', prop: 'password', required: false },
        { type: 'input', label: '真实姓名', prop: 'truename', required: false },
        { type: 'input', label: '手机号', prop: 'phone', required: false },
        { type: 'input', label: '邮箱', prop: 'email', required: false },
        { type: 'select', label: '角色', prop: 'role', required: true, opts: [
            { label: '超级管理员', value: 'global_admin' },
            { label: '应用管理员', value: 'app_admin' },
          ]
        },
        { type: 'datetime', label: '到期时间', prop: 'expires_at', required: true },
    ]
})
const visible = ref(false);
const isEdit = ref(false);
const rowData = ref({});
const handleEdit = (row: User) => {
    simpleApi.get('/api/admin/getDetail', { id: row.id }, permiss.token, function(data) {
      rowData.value = data;
      isEdit.value = true;
      visible.value = true;
    })
};
const updateData = (row: User) => {
   if(isEdit.value){
     const params = {
       id: row.id,
       username: row.username,
       password: row.password,
       truename: row.truename,
       phone: row.phone,
       email: row.email,
       role: row.role,
       expires_at: row.expires_at,
     }
     simpleApi.post('/api/admin/update', params, permiss.token, function(data) {
       closeDialog();
       getData();
     })
   }else{
     const params = {
       username: row.username,
       password: row.password,
       truename: row.truename,
       phone: row.phone,
       email: row.email,
       role: row.role,
       expires_at: row.expires_at,
     }
     simpleApi.post('/api/admin/create', params, permiss.token, function(data) {
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
const handleView = (row: User) => {
    simpleApi.get('/api/admin/getDetail', { id: row.id }, permiss.token, function(data) {
      viewData.value.row = data;
      viewData.value.list = [
        {
          prop: 'id',
          label: 'ID',
        },
        {
          prop: 'username',
          label: '用户名',
        },
        {
          prop: 'password',
          label: '密码',
        },
        {
          prop: 'truename',
          label: '真实姓名',
        },
        {
          prop: 'phone',
          label: '手机号',
        },
        {
          prop: 'email',
          label: '邮箱',
        },
        {
          prop: 'rolename',
          label: '角色',
        },
        {
          prop: 'expires_at',
          label: '到期时间',
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
    })
};

// 删除相关
const handleDelete = (row: User) => {
  simpleApi.postForm('/api/admin/delete', {id: row.id}, permiss.token, function(data){
    ElMessage.success('删除成功');
    getData();
  })
}
</script>

<style scoped></style>