<template>
    <div>
        <div class="config-container">
            <el-card class="config-content" shadow="hover" :body-style="{ padding: '20px 50px', height: '100%', boxSizing: 'border-box' }">
                  <el-form class="w500" label-position="top" >
                    <el-form-item label="执行器超时时间：">
                      <el-input-number type="number" v-model="formConfig.executor_timeout" min="3" max="1800"></el-input-number>
                      <span style="margin-left: 5px;">秒</span>
                    </el-form-item>
                    <el-form-item label="清理间隔时间：">
                      <el-input-number type="number" v-model="formConfig.clear_time" min="2" max="72"></el-input-number>
                      <span style="margin-left: 5px;">小时</span>
                    </el-form-item>
                    <el-form-item label="系统消息：">
                      <el-input type="textarea" v-model="formConfig.notice"></el-input>
                    </el-form-item>
                    <el-form-item>
                      <el-button type="primary" @click="onSubmitConfig">保存</el-button>
                    </el-form-item>
                  </el-form>
            </el-card>
        </div>
    </div>
</template>

<script setup lang="ts" name="config">
import {simpleApi} from "@/api";
import {ElMessage} from "element-plus";
import {usePermissStore} from "@/store/permiss";
import {Config} from "@/types/config";
import {ref} from "vue";

const permiss = usePermissStore();

const formConfig = ref<Config>({
  notice: '',
  executor_timeout: 0,
  clear_time: 0,
})
const getConfigs = () => {
  simpleApi.get('/task/api/config/getConfigs', {}, permiss.token, (data) => {
    formConfig.value = data;
  });
}
getConfigs();

const onSubmitConfig = () => {
  const params = {
    notice: formConfig.value.notice,
    executor_timeout: formConfig.value.executor_timeout,
    clear_time: formConfig.value.clear_time,
  }
  simpleApi.post('/task/api/config/updateConfigs', params, permiss.token, (data) => {
    ElMessage.success('更新成功');
  });
};

</script>

<style scoped>
.config-container {
    display: flex;
}
.config-content {
    flex: 1;
}
.w500 {
    width: 500px;
}
</style>