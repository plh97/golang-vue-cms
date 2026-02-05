const validEnvs = ['development', 'prod', 'test']

function getBaseURL() {
  const env = process.env.FES_ENV || process.env.NODE_ENV || 'test'
  if (!validEnvs.includes(env)) {
    console.warn(`Invalid environment: ${env}. Falling back to test environment.`)
    return ''
  }
  const envToUrlMap = {
    development: 'http://127.0.0.1:8291',
    // development: 'http://45.76.197.2:8291',
    // test: '',
    // prod: '',
  }
  return envToUrlMap[env as keyof typeof envToUrlMap] ?? ''
}

export const baseURL = getBaseURL()
