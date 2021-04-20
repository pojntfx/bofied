const cacheName = "app-" + "cc6a0c07474d19c0555566ff5b4c6d526e5ecc4f";

self.addEventListener("install", event => {
  console.log("installing app worker cc6a0c07474d19c0555566ff5b4c6d526e5ecc4f");
  self.skipWaiting();

  event.waitUntil(
    caches.open(cacheName).then(cache => {
      return cache.addAll([
        "",
        "/bofied",
        "/bofied/app.css",
        "/bofied/app.js",
        "/bofied/manifest.webmanifest",
        "/bofied/wasm_exec.js",
        "/bofied/web/app.wasm",
        "/bofied/web/icon.png",
        "/bofied/web/index.css",
        "https://unpkg.com/@patternfly/patternfly@4.96.2/patternfly-addons.css",
        "https://unpkg.com/@patternfly/patternfly@4.96.2/patternfly.css",
        
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
  console.log("app worker cc6a0c07474d19c0555566ff5b4c6d526e5ecc4f is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
