import { FetchError } from "ofetch";

type UseAuthOpts = {
  autoLogin?: boolean;
  logoutUrl?: string;
  onLogin?: () => Promise<void>;
  onError?: (error: FetchError) => void;
};

export const useAuth = (
  opts: UseAuthOpts = {
    autoLogin: true,
    logoutUrl: "/login",
  },
) => {
  const error = ref<FetchError | null>(null);
  const isLoggedIn = useState("isLoggedIn", () => false);
  const router = useRouter();

  if (opts.autoLogin) {
    console.log("auto login");
    $fetch("/api/auth/session")
      .then(() => {
        isLoggedIn.value = true;
        opts.onLogin?.();
      })
      .catch(() => {
        isLoggedIn.value = false;
      });
  }

  const login = async (username: string, password: string) => {
    // login logic
    try {
      await $fetch("/api/auth/login", {
        method: "POST",
        body: JSON.stringify({ username, password }),
        headers: { "Content-Type": "application/json" },
      });
      isLoggedIn.value = true;
      await opts.onLogin?.();
    } catch (e) {
      error.value = e as FetchError;
      opts.onError?.(e as FetchError);
    }
  };

  const logout = async () => {
    //   logout logic
    try {
      await $fetch("/api/auth/logout");
      isLoggedIn.value = false;
    } catch (e) {
      error.value = e as FetchError;
    }
    await router.push({ path: opts.logoutUrl });
  };

  const checkSession = async () => {
    try {
      await $fetch("/api/auth/session");
      isLoggedIn.value = true;
      return true;
    } catch (e) {
      isLoggedIn.value = false;
      return false;
    }
  };

  return {
    login,
    logout,
    error,
    isLoggedIn,
    checkSession,
  };
};
