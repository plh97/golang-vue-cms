import { FMessage } from '@fesjs/fes-design'
import { setToken } from './common/utils'
import { baseURL } from './config'

// Helper function to convert an object into a query string
function toQueryString(data: Record<string, any>): string {
  const params = new URLSearchParams()
  for (const key in data) {
    if (Object.prototype.hasOwnProperty.call(data, key) && data[key] !== undefined && data[key] !== null) {
      params.append(key, data[key])
    }
  }
  const queryString = params.toString()
  return queryString ? `?${queryString}` : ''
}

export async function request(url: string, data: Record<string, any> = {}, options: any = {}) {
  const config: RequestInit = {
    method: options.method || 'GET',
    headers: {
      // 'Token': getToken() ?? '',
      'Content-Type': 'application/json',
      ...(options.headers ?? {}),
    },
    credentials: 'include',
  }

  let finalUrl = `${baseURL}/v1${url}`

  // 1. 🔍 HANDLE GET REQUEST PARAMETERS
  if (config.method === 'GET') {
    const queryString = toQueryString(data)
    // Append the query string to the URL
    if (queryString) {
      finalUrl += queryString
    }
    // GET requests MUST NOT have a body
    config.body = undefined
  }

  // 2. 📝 HANDLE POST/PUT/DELETE REQUEST BODY (your original logic)
  else if (options.method) { // Not GET
    if (data instanceof FormData) {
      config.body = data
      // eslint-disable-next-line ts/ban-ts-comment
      // @ts-ignore
      delete config.headers['Content-Type']
    }
    else {
      // eslint-disable-next-line ts/ban-ts-comment
      // @ts-ignore
      config.headers['Content-Type'] = 'application/json'
      config.body = JSON.stringify(data)
    }
  }

  // 3. 🚀 FETCH CALL (Using the potentially modified finalUrl)
  return fetch(finalUrl, config).then(async (response) => {
    const data = await response.json();
    if (!response.ok) {
      const errorMessage = `${response.status}! ${data.message ?? ''}`
      FMessage.error({
        content: errorMessage,
      })
      return Promise.reject(new Error(errorMessage))// filepath: /Users/plh/code/golang-tutorial/web/src/api.ts
    }
    return data
  }).then((data) => {
    // 处理响应内容异常
    if (data?.code === 2000) {
      FMessage.error({
        content: data?.message,
      })
      setToken('')
      location.replace('/login')
      return Promise.reject(new Error('未登录或登录已过期，请重新登录'))
    }
    if (data?.code !== 0) {
      // Reject the promise with an error object containing code and message
      FMessage.error({
        content: data?.message,
      })
      return Promise.reject(data?.message)
    }
    if (options.transformData) {
      return options.transformData(data)
    }
    return data.data
  })
}
