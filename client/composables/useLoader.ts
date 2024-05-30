import { FetchError } from "ofetch";

export const useLoader = () => {
  const loading = ref<boolean>(false);
  const error = ref<FetchError | null>(null);

  async function withLoader<T = any>(fn: () => Promise<T>) {
    showLoader();
    try {
      return await fn();
    } catch (e) {
      error.value = e as FetchError;
      hideLoader();
    } finally {
      hideLoader();
    }
  }

  const showLoader = () => {
    loading.value = true;
  };

  const hideLoader = () => {
    loading.value = false;
  };

  return { loading, error, withLoader, showLoader, hideLoader };
};
