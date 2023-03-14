// @ts-ignore
/* eslint-disable */
import type { CodeListItemType } from '@/pages/ModuleList';
import { request } from '@umijs/max';

/** 获取当前的用户 GET /api/currentUser */
export async function currentUser(options?: { [key: string]: any }) {
  const data = {
    success: true,
    data: {
      name: 'admin',
      avatar: 'https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png',
      userid: '00000001',
      email: 'jyan@tricorder.dev',
      signature: '零成本云原生观测、助力工程师掌控云原生',
      title: '开发',
      group: 'tricorder',
      tags: [{ key: '0', label: '很有想法的' }],
      notifyCount: 12,
      unreadCount: 11,
      country: 'China',
      access: 'admin',
      geographic: {
        province: { label: '北京市', key: '100000' },
        city: { label: '朝阳区', key: '100020' },
      },
      address: '北京市海淀区静淑苑路2号全球创新社区401室',
      phone: '13120386629',
    },
  };
  return Promise.resolve(data);
  // return request<{
  //   data: API.CurrentUser;
  // }>('/api/currentUser', {
  //   method: 'GET',
  //   ...(options || {}),
  // });
}

/** 退出登录接口 POST /api/login/outLogin */
export async function outLogin(options?: { [key: string]: any }) {
  return request<Record<string, any>>('/api/login/outLogin', {
    method: 'POST',
    ...(options || {}),
  });
}

/** 登录接口 POST /api/login/account */
export async function login(body: API.LoginParams, options?: { [key: string]: any }) {
  return request<API.LoginResult>('/api/login/account', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 此处后端没有提供注释 GET /api/notices */
export async function getNotices(options?: { [key: string]: any }) {
  return request<API.NoticeIconList>('/api/notices', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 获取规则列表 GET /api/rule */
export async function rule(
  params: {
    // query
    /** 当前的页码 */
    current?: number;
    /** 页面的容量 */
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  return request<API.RuleList>('/api/rule', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 新建规则 PUT /api/rule */
export async function updateRule(options?: { [key: string]: any }) {
  return request<API.RuleListItem>('/api/rule', {
    method: 'PUT',
    ...(options || {}),
  });
}

/** 新建规则 POST /api/rule */
export async function addRule(options?: { [key: string]: any }) {
  return request<API.RuleListItem>('/api/rule', {
    method: 'POST',
    ...(options || {}),
  });
}

/** 删除规则 DELETE /api/rule */
export async function removeRule(options?: { [key: string]: any }) {
  return request<Record<string, any>>('/api/rule', {
    method: 'DELETE',
    ...(options || {}),
  });
}

// code form提交
export async function createModule(body?: { [key: string]: any }) {
  return request<API.ResponseType<any>>('/api/createModule', {
    method: 'POST',
    data: body,
  });
}

// code list
export async function listModule() {
  return request<API.ResponseType<CodeListItemType[]>>('/api/listModule', {
    method: 'GET',
    params: {},
  });
}

// code deploy
export async function deployModule(body: any) {
  return request<API.ResponseType<any>>('/api/deployModule', {
    method: 'POST',
    params: {
      id: body.Id,
    },
    data: body,
  });
}

// code undeploy
export async function undeployModule(body: any) {
  return request<API.ResponseType<any>>('/api/undeployModule', {
    method: 'POST',
    params: {
      id: body.Id,
    },
    data: body,
  });
}

// code delete
export async function deleteModule(body: any) {
  return request<API.ResponseType<any>>('/api/deleteModule', {
    method: 'GET',
    params: {
      id: body.Id,
    },
    data: body,
  });
}
