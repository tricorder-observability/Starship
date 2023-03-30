/**
 * @name umi 的路由配置
 * @description 只支持 path,component,routes,redirect,wrappers,title 的配置
 * @param path  path 只支持两种占位符配置，第一种是动态参数 :id 的形式，第二种是 * 通配符，通配符只能出现路由字符串的最后。
 * @param component 配置 location 和 path 匹配后用于渲染的 React 组件路径。可以是绝对路径，也可以是相对路径，如果是相对路径，会从 src/pages 开始找起。
 * @param routes 配置子路由，通常在需要为多个路径增加 layout 组件时使用。
 * @param redirect 配置路由跳转
 * @param wrappers 配置路由组件的包装组件，通过包装组件可以为当前的路由组件组合进更多的功能。 比如，可以用于路由级别的权限校验
 * @doc https://umijs.org/docs/guides/routes
 */
export default [
  {
    path: '/create-module',
    name: '新建模组',
    name_en: 'Apply Observe',
    icon: 'compass',
    component: './CreateModule',
  },
  {
    path: '/module-list',
    name: '模组列表',
    name_en: 'Module List',
    icon: 'table',
    component: './ModuleList',
  },
  {
    path: '/grafana',
    name: 'Grafana',
    name_en: 'Grafana',
    icon: 'https://grafana.com/static/img/menu/grafana2.svg',
    target: '_blank'
  },
  {
    path: '/',
    redirect: '/create-module',
  },
  {
    path: '*',
    layout: false,
    component: './404',
  },
];
