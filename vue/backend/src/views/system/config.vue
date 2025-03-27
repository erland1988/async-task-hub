<template>
  <div>
    <div class="config-container">
      <el-card class="config-content" shadow="hover" :body-style="{ padding: '20px 50px', height: '100%', boxSizing: 'border-box' }">
        <el-form class="w500" label-position="top" >
          <el-form-item label="执行器超时时间：">
            <!-- 使用计算属性绑定值 -->
            <el-input-number type="number" v-model="executorTimeoutNumber" :min=3 :max=1800></el-input-number>
            <span style="margin-left: 5px;">秒</span>
          </el-form-item>
          <el-form-item label="清理间隔时间：">
            <!-- 使用计算属性绑定值 -->
            <el-input-number type="number" v-model="clearTimeNumber" :min=2 :max=72></el-input-number>
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
import {ref, computed} from "vue";

const permiss = usePermissStore();

const formConfig = ref<Config>({
  notice: "",
  executor_timeout: "",
  clear_time: "",
})

// 定义计算属性来处理 executor_timeout 的类型转换
const executorTimeoutNumber = computed({
  get() {
    return Number(formConfig.value.executor_timeout);
  },
  set(newValue) {
    formConfig.value.executor_timeout = String(newValue);
  }
});

// 定义计算属性来处理 clear_time 的类型转换
const clearTimeNumber = computed({
  get() {
    return Number(formConfig.value.clear_time);
  },
  set(newValue) {
    formConfig.value.clear_time = String(newValue);
  }
});

const getConfigs = () => {
  simpleApi.get('/task/api/config/getConfigs', {}, permiss.token, (data) => {
    formConfig.value = {
      notice: data.notice,
      executor_timeout: data.executor_timeout,
      clear_time: data.clear_time,
    }
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