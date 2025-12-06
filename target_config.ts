// target_config.ts
export const target_tauri = true; // true - Tauri, false - dev сервер

export const LOCAL_IP_API_ADDR = "http://192.168.56.1:8080"; 
export const DEV_API_ADDR = "http://localhost:8080"; 

export const API_BASE = target_tauri ? LOCAL_IP_API_ADDR : DEV_API_ADDR;

export const IMG_BASE = target_tauri
  ? "http://192.168.56.1:9000/loading-time-img/img"
  : "/loading-time-img/img";
