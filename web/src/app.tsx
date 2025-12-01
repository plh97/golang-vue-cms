import {
  access,
  defineRuntimeConfig,
} from '@fesjs/fes'
import { FMessage } from '@fesjs/fes-design'
import PageLoading from '@/components/pageLoading.vue'
import UserCenter from '@/components/userCenter.vue'
import { request } from './api'
import { getToken } from './common/utils'
import { baseURL } from './config'

export default defineRuntimeConfig({
  login: {
    loginPath: '/login', // 登陆页面路径，默认 /login，也可以用路由的 name
    hasLogin() {
      // 进入页面前判断是否登陆的逻辑，每次跳转非登陆页面都会检测，直到为 true，支持异步
      const token = getToken()
      return !!token
    },
  },
  request: {
    baseURL,
    timeout: 10000, // 默认 10s
    credentials: 'omit', // 默认 include, 'include' | 'same-origin' | 'omit'
    method: 'get', // 默认 post
    // mergeRequest: false, // 是否合并请求
    // responseType: null, // 可选 'json' | 'text' | 'blob' | 'arrayBuffer' | 'formData'，默认根据 content-type 处理
    // credentials: "include", // 默认 include, 'include' | 'same-origin' | 'omit'
    // headers: {
    //   token: await getToken() ?? '',
    // }, // 传给服务器的 header
    // cacheData: false, // 是否缓存
    // requestInterceptor: (config: Config) => Config,
    // // responseInterceptor: (response: RequestResponse) => RequestResponse,
    transformData(data) {
      // 处理响应内容异常
      if (data?.code === 2000) {
        FMessage.error({
          content: data?.message,
        })
        return Promise.reject(new Error('未登录或登录已过期，请重新登录'))
      }
      if (data?.code !== 0) {
        // Reject the promise with an error object containing code and message
        FMessage.error({
          content: data?.message,
        })
        return Promise.reject(data?.message)
      }
      return data.data
    },
    // // http 异常，和插件异常
    // errorHandler(error) {
    //     // 处理业务异常，例如上述 transformData 抛出的异常
    //     // if (error.code) {
    //     //     console.log(error.message)
    //     // } else if (error.response) {
    //     //     // 请求成功发出且服务器也响应了状态码，但状态代码超出了 2xx 的范围
    //     //     console.log(`服务异常：${error.response.status}`)
    //     // } else {
    //     //     // 请求异常
    //     //     console.log(error.message || error.message || `请求失败`)
    //     // }
    // },
    // 支持其他 fetch 配置
    // ...otherConfigs,
  },

  beforeRender: {
    loading: <PageLoading />,
    action: async () => {
      access.setRole('admin')
      if (location.pathname === '/login') {
        return
      }
      const res = await request('/profile')
      return new Promise((resolve) => {
        // 初始化应用的全局状态，可以通过 useModel('@@initialState') 获取，具体用法看@/components/UserCenter 文件
        resolve(res)
      })
    },
  },
  layout: {
    renderCustom: () => <UserCenter />,
  },
})
