import { Menus } from '@/types/menu';

export const menuData: Menus[] = [
    {
        id: '0',
        title: '系统首页',
        index: '/dashboard',
        icon: 'Odometer',
    },
    {
        id: '1',
        title: '用户管理',
        index: '1',
        icon: 'HomeFilled',
        children: [
            {
                id: '11',
                pid: '1',
                index: '/user-index',
                title: '用户列表',
            },
        ],
    },
    {
        id: '2',
        title: '应用管理',
        index: '2',
        icon: 'HomeFilled',
        children: [
            {
                id: '21',
                pid: '2',
                index: '/app-index',
                title: '应用列表',
            },
        ],
    },
    {
        id: '3',
        title: '任务管理',
        index: '3',
        icon: 'HomeFilled',
        children: [
            {
                id: '31',
                pid: '3',
                index: '/task-index',
                title: '任务列表',
            },
        ],
    },
    {
        id: '4',
        title: '任务队列管理',
        index: '4',
        icon: 'HomeFilled',
        children: [
            {
                id: '41',
                pid: '4',
                index: '/taskqueue-index',
                title: '队列列表',
            },
        ],
    },
    {
        id: '5',
        title: '任务日志管理',
        index: '5',
        icon: 'HomeFilled',
        children: [
            {
                id: '51',
                pid: '5',
                index: '/tasklog-index',
                title: '任务日志列表',
            },
        ],
    },
    {
        id: '6',
        title: '系统管理',
        index: '6',
        icon: 'HomeFilled',
        children: [
            {
                id: '61',
                pid: '6',
                index: '/system-login-log',
                title: '登录日志',
            },
            {
                id: '62',
                pid: '6',
                index: '/system-log',
                title: '操作日志',
            }
        ],
    },
];
