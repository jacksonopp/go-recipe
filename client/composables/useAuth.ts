import { FetchError } from "ofetch";

type UseAuthOpts = {
  checkInitial?: boolean;
};

export const useAuth = (opts: UseAuthOpts = { checkInitial: true }) => {
  const error = ref<FetchError | null>(null);
  const isLoggedIn = useState("isLoggedIn", () => false);

  if (opts.checkInitial) {
    console.log("checking initial");
    $fetch("/api/auth/session")
      .then(() => {
        isLoggedIn.value = true;
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
    } catch (e) {
      error.value = e as FetchError;
    }
  };

  const logout = () => {
    // logout logic
    console.log("logging out");
  };

  const onLogin = (callback: () => void) => {
    watch(isLoggedIn, (newLoggedIn) => {
      if (newLoggedIn) {
        callback();
      }
    });
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

  const onError = (callback: (error: FetchError) => void) => {
    watch(error, (newError) => {
      if (newError !== null) {
        callback(newError);
      }
    });
  };

  return {
    login,
    logout,
    error,
    onError,
    onLogin,
    isLoggedIn,
    checkSession,
  };
};
