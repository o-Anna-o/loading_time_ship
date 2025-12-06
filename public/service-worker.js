const CACHE_NAME = "loading-time-cache-v2";

// Файлы, которые точно существуют в dist/ после сборки Vite + public/
const FILES_TO_CACHE = [
  "/loading-time-frontend/",
  "/loading-time-frontend/index.html",
  "/loading-time-frontend/manifest.json",
  "/loading-time-frontend/logo192.png",
  "/loading-time-frontend/logo512.png",
  "/loading-time-frontend/default.png"
];

// Установка — кладём статические файлы в кэш
self.addEventListener("install", event => {
  event.waitUntil(
    caches.open(CACHE_NAME).then(cache => cache.addAll(FILES_TO_CACHE))
  );
  self.skipWaiting();
});

// Активация — очищаем старые версии кэша
self.addEventListener("activate", event => {
  event.waitUntil(
    caches.keys().then(keys =>
      Promise.all(keys.map(key => key !== CACHE_NAME && caches.delete(key)))
    )
  );
  self.clients.claim();
});

// Fetch — cache first 
self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(cached => {
      return cached || fetch(event.request);
    })
  );
});
