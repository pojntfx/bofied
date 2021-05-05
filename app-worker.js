const cacheName = "app-" + "6059d284c042a48fa31be4f3625cb0776dad509d";

self.addEventListener("install", event => {
  console.log("installing app worker 6059d284c042a48fa31be4f3625cb0776dad509d");
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
  console.log("app worker 6059d284c042a48fa31be4f3625cb0776dad509d is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
