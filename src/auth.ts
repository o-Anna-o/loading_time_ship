export function saveToken(token:string) {
  localStorage.setItem('lt_token', token)
}
export function getToken(): string | null {
  return localStorage.getItem('lt_token')
}
export function clearToken() {
  localStorage.removeItem('lt_token')
}
export function isLoggedIn() {
  return !!getToken()
}
