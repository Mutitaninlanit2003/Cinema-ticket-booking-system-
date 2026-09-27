// Module augmentation for vue-router's RouteMeta.
// Extend this as routes need more meta (e.g. `title` for AdminTopbar,
// `requiresAuth` once an auth strategy replaces the dropped Google OAuth).
import 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    layout?: 'user' | 'admin'
    title?: string
  }
}
