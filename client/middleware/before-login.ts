export default defineNuxtRouteMiddleware(async (to, from) => {
  if (import.meta.server) return;
  if (useNuxtApp().isHydrating) return;
  const router = useRouter();
  const { isLoggedIn, checkSession } = useAuth({ checkInitial: false });

  if (isLoggedIn.value) {
    return router.push({ path: "/home" });
  } else {
    const isSessionValid = await checkSession();
    if (isSessionValid) {
      return router.push({ path: "/home" });
    }
    return;
  }
});
