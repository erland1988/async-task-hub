import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';
import { usePermissStore } from '../store/permiss';
import Home from '../views/home.vue';
import NProgress from 'nprogress';
import 'nprogress/nprogress.css';

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        redirect: '/dashboard',
    },
    {
        path: '/',
        name: 'Home',
        component: Home,
        children: [
            {
                path: '/dashboard',
                name: 'dashboard',
                meta: {
                    title: '系统首页',
                    permiss: '0',
                },
                component: () => import('../views/dashboard.vue'),
            },
            {
                path: '/theme',
                name: 'theme',
                meta: {
                    title: '设置主题',
                    permiss: '0',
                },
                component: () => import('../views/pages/theme.vue'),
            },
            {
                path: '/ucenter',
                name: 'ucenter',
                meta: {
                    title: '个人中心',
                    permiss: '0',
                },
                component: () => import('../views/pages/ucenter.vue'),
            },
            {
                path: '/user-index',
                name: 'user-index',
                meta: {
                    title: '用户列表',
                    permiss: '11',
                },
                component: () => import('../views/user/index.vue'),
            },
            {
                path: '/app-index',
                name: 'app-index',
                meta: {
                    title: '应用列表',
                    permiss: '21',
                },
                component: () => import('../views/app/index.vue'),
            },
            {
                path: '/task-index',
                name: 'task-index',
                meta: {
                    title: '任务列表',
                    permiss: '31',
                },
                component: () => import('../views/task/index.vue'),
            },
            {
                path: '/taskqueue-index',
                name: 'taskqueue-index',
                meta: {
                    title: '队列列表',
                    permiss: '41',
                },
                component: () => import('../views/taskqueue/index.vue'),
            },
            {
                path: '/tasklog-index',
                name: 'tasklog-index',
                meta: {
                    title: '任务日志列表',
                    permiss: '51',
                },
                component: () => import('../views/tasklog/index.vue'),
            },
            {
                path: '/system-log',
                name: 'system-log',
                meta: {
                    title: '操作日志',
                    permiss: '61',
                },
                component: () => import('../views/system/log.vue'),
            },
        ],
    },
    {
        path: '/login',
        meta: {
            title: '登录',
            noAuth: true,
        },
        component: () => import('../views/pages/login.vue'),
    },
    {
        path: '/403',
        meta: {
            title: '没有权限',
            noAuth: true,
        },
        component: () => import('../views/pages/403.vue'),
    },
    {
        path: '/404',
        meta: {
            title: '找不到页面',
            noAuth: true,
        },
        component: () => import('../views/pages/404.vue'),
    },
    { path: '/:path(.*)', redirect: '/404' },
];

const router = createRouter({
    history: createWebHashHistory(),
    routes,
});

router.beforeEach((to, from, next) => {
    NProgress.start();
    const username = localStorage.getItem('admin:username');
    const permiss = usePermissStore();

    if (!username && to.meta.noAuth !== true) {
        next('/login');
    } else if (typeof to.meta.permiss == 'string' && !permiss.key.includes(to.meta.permiss)) {
        // 如果没有权限，则进入403
        next('/403');
    } else {
        next();
    }
});

router.afterEach(() => {
    NProgress.done();
});

export default router;
