import axios from "axios";

export const apiClient = axios.create({
  baseURL: window.__APP_CONFIG__?.apiBaseUrl || import.meta.env.VITE_API_BASE_URL,
});
