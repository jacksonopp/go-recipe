import { useAuth } from "~/composables/useAuth";

export default defineNuxtRouteMiddleware(async (to, from) => {
  if (import.meta.server) return;

  const { isLoggedIn, checkSession } = useAuth({ checkInitial: false });

  if (!isLoggedIn.value) {
    const isSessionValid = await checkSession();
    if (!isSessionValid) {
      return navigateTo("/login");
    }
  }
  // try {
  //   await $fetch.raw("/api/auth/session");
  //   return;
  // } catch (e) {
  //   return navigateTo("/login");
  // }
});
