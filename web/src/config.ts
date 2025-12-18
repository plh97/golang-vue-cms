const validEnvs = ['development', 'prod', 'test']

function getBaseURL() {
  const env = process.env.FES_ENV || process.env.NODE_ENV || 'test'
  if (!validEnvs.includes(env)) {
    console.warn(`Invalid environment: ${env}. Falling back to test environment.`)
    // return 'https://test-manage-api.ohayo.date'
  }
  const envToUrlMap = {
    development: 'http://localhost:8291',
    test: 'http://18.143.194.24:8291',
    prod: 'http://18.143.194.24:8291',
  }
  return envToUrlMap[env as keyof typeof envToUrlMap] ?? ''
}

export const baseURL = getBaseURL()
