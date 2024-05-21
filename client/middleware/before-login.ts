export default defineNuxtRouteMiddleware(async (to, from) => {
  if (import.meta.server) return;
  if (useNuxtApp().isHydrating) return;
  const router = useRouter();
  try {
    await $fetch.raw("/api/auth/session");
    return router.push({ path: "/home" });
  } catch (e) {
    return;
  }
});
