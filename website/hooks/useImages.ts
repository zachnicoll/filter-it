import { AxiosError } from "axios";
import { toastError } from "common/toast";
import { Filter, ImageDocument } from "api/types";
import { useCallback, useEffect, useState } from "react";
import API from "api";

interface HookReturn {
  images: ImageDocument[];
  loading: boolean;
}

export const useImages = (filter: Filter | null): HookReturn => {
  const [images, setImages] = useState<ImageDocument[]>([]);
  const [loading, setLoading] = useState(true);

  const fetch = useCallback(async (): Promise<void> => {
    setLoading(true);

    try {
      // Replace with real API request when lambda working
      const data = await API.feed.get(filter);
      setImages(data);
    } catch (e) {
      const err = e as AxiosError;
      toastError((e as AxiosError).message, err);
    }

    setLoading(false);
  }, [filter]);

  useEffect(() => {
    fetch();
  }, [fetch]);

  return { images, loading };
};
