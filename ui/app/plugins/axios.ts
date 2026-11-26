import axios from 'axios'
import { ADMIN_TOKEN_KEY, ADMIN_USER_KEY, API_BASE_URL, normalizeApiUrl } from '~/constants/auth'

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  const rawBase = ((config.public as Record<string, any>)?.apiBase as string) || API_BASE_URL
  const baseURL = normalizeApiUrl(rawBase)

  const api = axios.create({
    baseURL,
    headers: {
      'Content-Type': 'application/json',
    },
  })

  api.interceptors.request.use((cfg) => {
    if (import.meta.client) {
      const token = localStorage.getItem(ADMIN_TOKEN_KEY)
      if (token) {
        cfg.headers.Authorization = `Bearer ${token}`
      }
    }
    return cfg
  })

  api.interceptors.response.use(
    (response) => response,
    (error) => {
      if (import.meta.client && error.response?.status === 401) {
        localStorage.removeItem(ADMIN_TOKEN_KEY)
        localStorage.removeItem(ADMIN_USER_KEY)
        navigateTo('/login')
      }
      return Promise.reject(error)
    }
  )

  return {
    provide: {
      axios: api,
    },
  }
})
