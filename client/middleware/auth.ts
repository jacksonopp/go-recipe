export default defineNuxtRouteMiddleware(async (to, from) => {
  if (import.meta.server) return;

  try {
    await $fetch.raw("/api/auth/session");
    return;
  } catch (e) {
    return navigateTo("/login");
  }
});
