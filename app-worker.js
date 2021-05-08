const cacheName = "app-" + "415cf82a0c0e8ceeb2178adcaa2a8fb9238feb74";

self.addEventListener("install", event => {
  console.log("installing app worker 415cf82a0c0e8ceeb2178adcaa2a8fb9238feb74");
  self.skipWaiting();

  event.waitUntil(
    caches.open(cacheName).then(cache => {
      return cache.addAll([
        "/bofied",
        "/bofied/app.css",
        "/bofied/app.js",
        "/bofied/manifest.webmanifest",
        "/bofied/wasm_exec.js",
        "/bofied/web/app.wasm",
        "/bofied/web/icon.png",
        "/bofied/web/index.css",
        "https://unpkg.com/@patternfly/patternfly@4.102.2/patternfly-addons.css",
        "https://unpkg.com/@patternfly/patternfly@4.102.2/patternfly.css",
        
      ]);
    })
  );
});

self.addEventListener("activate", event => {
  event.waitUntil(
    caches.keys().then(keyList => {
      return Promise.all(
        keyList.map(key => {
          if (key !== cacheName) {
            return caches.delete(key);
          }
        })
      );
    })
  );
  console.log("app worker 415cf82a0c0e8ceeb2178adcaa2a8fb9238feb74 is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
