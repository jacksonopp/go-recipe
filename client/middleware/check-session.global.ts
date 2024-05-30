import { useAuth } from "~/composables/useAuth";

export default defineNuxtRouteMiddleware(async (to, from) => {
  if (import.meta.server) return;

  const { isLoggedIn, checkSession } = useAuth({ autoLogin: false });

  if (!isLoggedIn.value) {
    await checkSession();
  }
});
