<template>
    <div>
        <div class="user-container">
            <el-card class="user-profile" shadow="hover" :body-style="{ padding: '0px' }">
                <div class="user-profile-bg"></div>
                <div class="user-info">
                    <div class="info-name">{{ username }}</div>
                    <div class="info-desc">
                      <span>{{ role }}</span>
                      <el-divider direction="vertical" />
                      <span>{{ expiresAt }}</span>
                    </div>
                    <div class="info-desc"></div>
                </div>
            </el-card>
            <el-card class="user-content" shadow="hover" :body-style="{ padding: '20px 50px', height: '100%', boxSizing: 'border-box' }">
                <el-tabs tab-position="left" v-model="activeName">
                    <el-tab-pane name="label1" label="系统通知" class="user-tabpane">
                      <div class="plugins-tips">
                        {{ notice }}
                      </div>
                    </el-tab-pane>
                    <el-tab-pane name="label2" label="基本信息" class="user-tabpane">
                      <el-form class="w500" label-position="top">
                        <el-form-item label="真实姓名：">
                          <el-input type="text" v-model="formProfile.truename"></el-input>
                        </el-form-item>
                        <el-form-item label="手机号：">
                          <el-input type="text" v-model="formProfile.phone"></el-input>
                        </el-form-item>
                        <el-form-item label="邮箱：">
                          <el-input type="text" v-model="formProfile.email"></el-input>
                        </el-form-item>
                        <el-form-item>
                          <el-button type="primary" @click="onSubmitProfile">保存</el-button>
                        </el-form-item>
                      </el-form>
                    </el-tab-pane>
                    <el-tab-pane name="label3" label="修改密码" class="user-tabpane">
                        <el-form class="w500" label-position="top">
                            <el-form-item label="旧密码：">
                                <el-input type="password" v-model="formPassword.oldPassword"></el-input>
                            </el-form-item>
                            <el-form-item label="新密码：">
                                <el-input type="password" v-model="formPassword.newPassword"></el-input>
                            </el-form-item>
                            <el-form-item label="确认新密码：">
                                <el-input type="password" v-model="formPassword.confirmPassword"></el-input>
                            </el-form-item>
                            <el-form-item>
                                <el-button type="primary" @click="onSubmitPassword">保存</el-button>
                            </el-form-item>
                        </el-form>
                    </el-tab-pane>
                    <el-tab-pane name="label4" label="关于作者" class="user-tabpane">
                        <div class="plugins-tips">
                            如果该系统
                            <el-link href="https://github.com/erland1988/async-task-hub/" target="_blank">async-task-hub</el-link>
                            对你有帮助，或者你喜欢这个系统，可以请作者喝杯咖啡，鼓励一下作者继续维护！
                        </div>
                        <div class="plugins-tips">
                          加QQ373944668探讨问题。
                        </div>
                    </el-tab-pane>
                </el-tabs>
            </el-card>
        </div>
    </div>
</template>

<script setup lang="ts" name="ucenter">
import { reactive, ref } from 'vue';
import {loginOut, simpleApi} from "@/api";
import {ElMessage} from "element-plus";
import {usePermissStore} from "@/store/permiss";

const username = localStorage.getItem('admin:username');
const role = localStorage.getItem('admin:role');
const truename = localStorage.getItem('admin:truename') || '';
const phone = localStorage.getItem('admin:phone') || '';
const email = localStorage.getItem('admin:email') || '';
const expiresAt = localStorage.getItem('admin:expires_at') || '';

const permiss = usePermissStore();

const notice = permiss.config['notice'];

const formProfile = reactive({
  truename: truename,
  phone: phone,
  email: email,
});
const onSubmitProfile = () => {
  const params = {
    truename: formProfile.truename,
    phone: formProfile.phone,
    email: formProfile.email,
  };

  simpleApi.post('/api/admin/updateProfile', params, permiss.token, (data) => {
    ElMessage.success('基本信息更新成功');
    localStorage.setItem('admin:truename', formProfile.truename);
    localStorage.setItem('admin:phone', formProfile.phone);
    localStorage.setItem('admin:email', formProfile.email);
  });
};

const formPassword = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
});
const onSubmitPassword = () => {
  const params = {
    old_password: formPassword.oldPassword,
    new_password: formPassword.newPassword,
    confirm_password: formPassword.confirmPassword,
  };

  simpleApi.post('/api/admin/resetPassword', params, permiss.token, (data) => {
    ElMessage.success('密码修改成功，请重新登录');
    // 修改密码成功后，自动退出登录
    simpleApi.post('/api/admin/loginout', {}, permiss.token, function(data){
      loginOut();
    });
  });
};

const activeName = ref('label1');

</script>

<style scoped>
.user-container {
    display: flex;
}

.user-profile {
    position: relative;
}

.user-profile-bg {
    width: 100%;
    height: 200px;
    background-image: url('../../assets/img/ucenter-bg.jpg');
    background-size: cover;
    background-position: center;
    background-repeat: no-repeat;
}

.user-profile {
    width: 500px;
    margin-right: 20px;
    flex: 0 0 auto;
    align-self: flex-start;
}

.user-info {
    text-align: center;
    padding: 80px 0 30px;
}

.info-name {
    margin: 0 0 20px;
    font-size: 22px;
    font-weight: 500;
    color: #373a3c;
}

.info-desc {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 5px;
}

.info-desc,
.info-desc a {
    font-size: 18px;
    color: #55595c;
}
.info-icon i {
    font-size: 30px;
    margin: 0 10px;
    color: #343434;
}

.user-content {
    flex: 1;
}

.user-tabpane {
    padding: 10px 20px;
}

.w500 {
    width: 500px;
}

.user-footer > div + div {
    border-left: 1px solid rgba(83, 70, 134, 0.1);
}
</style>

<style>
.el-tabs.el-tabs--left {
    height: 100%;
}
</style>
