import { defineStore } from 'pinia';

export const usePermissStore = defineStore('permiss', {
    state: () => ({
        key: JSON.parse(localStorage.getItem('admin:permiss') || '[]') as string[], // 初始化权限
        token: localStorage.getItem('admin:token') || '', // 初始化 token
        config: JSON.parse(localStorage.getItem('admin:config') || '{}') as Record<string, string>, // 配置映射
    }),
    actions: {
        setKey(val: string[]) {
            this.key = val;
            localStorage.setItem('admin:permiss', JSON.stringify(val));
        },
        setToken(token: string) {
            this.token = token;
            localStorage.setItem('admin:token', token);
        },
        setConfig(config: Record<string, string>) {
            this.config = { ...this.config, ...config };
            localStorage.setItem('admin:config', JSON.stringify(this.config));
        },
        clearAuthData() {
            this.key = [];
            this.token = '';
            this.config = {};
            localStorage.removeItem('admin:permiss');
            localStorage.removeItem('admin:token');
            localStorage.removeItem('admin:config');
        },
    },
});

