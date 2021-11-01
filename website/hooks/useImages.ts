import { AxiosError } from "axios";
import { toastError } from "common/toast";
import { Filter, ImageDocument } from "api/types";
import { useCallback, useEffect, useState } from "react";
import API from "api";
import { useSearch } from "../context";

interface HookReturn {
  images: ImageDocument[];
  loading: boolean;
}

export const useImages = (): HookReturn => {
  const { searchState } = useSearch();
  const [images, setImages] = useState<ImageDocument[]>([]);
  const [loading, setLoading] = useState(true);

  const fetch = useCallback(async (): Promise<void> => {
    setLoading(true);

    try {
      const data = await API.feed.get(searchState.search);
      setImages(data);
    } catch (e) {
      const err = e as AxiosError;
      toastError((e as AxiosError).message, err);
    }

    setLoading(false);
  }, [searchState]);

  useEffect(() => {
    fetch();
  }, [fetch, searchState]);

  return { images, loading };
};
