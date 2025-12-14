import { emit, on } from './event'
import axios, { type AxiosError, type AxiosRequestConfig, type AxiosResponse, type AxiosRequestHeaders, type InternalAxiosRequestConfig, type Method } from 'axios'
import { FMessage } from '@fesjs/fes-design'
import { baseURL } from './config'
import { setToken } from './common/utils'

interface ApiErrorPayload {
  message?: string
  code?: number
  [key: string]: any
}

on('unlogin', (payload: ApiErrorPayload) => {
  FMessage.error({ content: `[${payload.code ?? 2000}] ${payload.message || '未登录或登录已过期，请重新登录'}` })
  setToken('')
  location.replace('/login')
})
on('error', (payload: ApiErrorPayload) => {
  if (payload.message) FMessage.error({ content: `[${payload.code ?? 500}] ${payload.message}` })
})

// Simple cookie reader to fetch accessToken for Authorization header
function getAccessToken(): string {
  const match = document.cookie.split('; ').find((row) => row.startsWith('accessToken='))
  return match ? decodeURIComponent(match.split('=')[1]) : ''
}

const api = axios.create({
  baseURL: `${baseURL}/v1`,
  withCredentials: true,
})

api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = getAccessToken()
  if (token) {
    // backend StrictAuth expects the raw token in Authorization header
    config.headers = {
      ...(config.headers || {}),
      Authorization: token,
    } as AxiosRequestHeaders
  }
  return config
})

api.interceptors.response.use(
  (response: AxiosResponse) => response,
  (error: AxiosError) => {
    const status = error.response?.status
    const respData = (error.response?.data ?? {}) as { message?: string; code?: number }
    const msg = respData.message || error.message || 'Request failed'
    const code = respData.code || status || 500
    const payload: ApiErrorPayload = { message: msg, code, ...respData }
    if (status === 401) {
      emit('unlogin', payload)
    } else {
      emit('error', payload)
    }
    return Promise.reject(error)
  },
)

type RequestOptions = AxiosRequestConfig & {
  transformData?: (data: any) => any
}

export async function request<T = any>(url: string, data: Record<string, any> = {}, options: RequestOptions = {}) {
  const method = (options.method || 'GET').toString().toUpperCase() as Method
  const config: AxiosRequestConfig = {
    url,
    method,
    headers: options.headers,
    params: method === 'GET' ? data : options.params,
    data: method === 'GET' ? undefined : data,
    ...options,
  }

  const res = await api.request(config)
  const resp = res.data

  if (resp?.code === 2000) {
    emit('unlogin', resp)
    return Promise.reject(new Error('未登录或登录已过期，请重新登录'))
  }

  if (resp?.code !== 0) {
    emit('error', resp)
    return Promise.reject(resp?.message)
  }

  if (options.transformData) {
    return options.transformData(resp)
  }

  return resp.data as T
}
