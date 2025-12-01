import { defineBuildConfig } from '@fesjs/fes'

export default defineBuildConfig({
  title: 'OHAYO 后台管理系统',
  router: {
    mode: 'history',
  },
  access: {
    roles: {
      admin: ['*'],
      manager: ['/'],
    },
  },
  layout: {
    title: 'OHAYO 后台管理系统',
    logo: '/favicon_io/favicon.ico',
    // footer: 'Created by MumbleFE',
    navigation: 'mixin',
    multiTabs: false,
    menus: [
      {
        name: 'user/list',
        match: ['user/*'],
        icon: 'UserOutlined',
        // icon: '/user.svg',
      },
      {
        name: 'account/list',
        match: ['account/*'],
        icon: 'UserManagementOutlined',
        // icon: '/server.svg',
      },
      {
        name: 'activity/list',
        match: ['activity/*'],
        icon: 'AppstoreOutlined',
        // icon: '/server.svg',
      },
      {
        name: 'role/list',
        match: ['role/*'],
        icon: 'UserOutlined',
        // icon: '/user.svg',
      },
    ],
  },
  enums: {
    status: [
      ['0', '无效的'],
      ['1', '有效的'],
    ],
  },
})
