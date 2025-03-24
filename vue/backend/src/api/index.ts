import axios from 'axios';
import {ElMessage} from "element-plus";
import request from '../utils/request';
import router from "@/router";
import {usePermissStore} from "@/store/permiss";

export const fetchData = () => {
    return request({
        url: './mock/table.json',
        method: 'get'
    });
};

export const fetchUserData = () => {
    return request({
        url: './mock/user.json',
        method: 'get'
    });
};


export function loginOut() {
    localStorage.removeItem('admin:username');
    usePermissStore().clearAuthData();
    router.push('/login').then(() => {
        window.location.reload(); // 强制刷新
    });
}

export const simpleApi = {
    // 通用的 POST 请求方法
    async post(url: string, params = {}, token = '', callback = (data: any) => {}): Promise<void> {
        try {
            const response = await axios.post(url, params, {
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer '+token,
                },
            });

            const { success, message = '', data } = response.data;

            if (success) {
                callback(data);
            } else {
                ElMessage.error(message);
            }
        } catch (error) {
            if (axios.isAxiosError(error) && error.response?.status === 401) {
                ElMessage.error('登录已过期，请重新登录');
                loginOut();
            } else {
                console.error('POST 请求失败:', error);
                ElMessage.error('请求失败，请稍后再试');
            }
        }
    },

    // 通用的 GET 请求方法
    async get(url: string, params = {}, token = '', callback = (data: any) => {}): Promise<void> {
        try {
            const response = await axios.get(url, {
                params,
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer '+token,
                },
            });

            const { success, message = '', data } = response.data;

            if (success) {
                callback(data);
            } else {
                ElMessage.error(message);
            }
        } catch (error) {
            if (axios.isAxiosError(error) && error.response?.status === 401) {
                ElMessage.error('登录已过期，请重新登录');
                loginOut();
            } else {
                console.error('POST 请求失败:', error);
                ElMessage.error('请求失败，请稍后再试');
            }
        }
    },

    // 通用的 POST Form 请求方法
    async postForm(url: string, params = {}, token = '', callback = (data: any) => {}): Promise<void> {
        const formData = new URLSearchParams();
        for (const [key, value] of Object.entries(params)) {
            formData.append(key, String(value));
        }
        try {
            const response = await axios.post(url, formData, {
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                    'Authorization': 'Bearer '+token,
                },
            });

            const { success, message = '', data } = response.data;

            if (success) {
                callback(data);
            } else {
                ElMessage.error(message);
            }
        } catch (error) {
            if (axios.isAxiosError(error) && error.response?.status === 401) {
                ElMessage.error('登录已过期，请重新登录');
                loginOut();
            } else {
                console.error('POST Form 请求失败:', error);
                ElMessage.error('请求失败，请稍后再试');
            }
        }
    }
};
