import { AxiosError } from "axios";
import { toastError } from "common/toast";
import { Filter, ImageDocument } from "api/types";
import { useCallback, useEffect, useState } from "react";
import API from "api";

interface HookReturn {
  images: ImageDocument[];
  loading: boolean;
}

const mockData = [
  {
    image: "https://via.placeholder.com/1080x1080/eee?text=1:1",
    progress: 2,
    tag: 0,
    author: "Author 0",
    title: "Image 0",
    id: "0",
  },
  {
    image: "https://via.placeholder.com/1080x1080/eee?text=1:1",
    progress: 2,
    tag: 1,
    author: "Author 1",
    title: "Image 1",
    id: "1",
  },
  {
    image: "https://via.placeholder.com/1080x1080/eee?text=1:1",
    progress: 2,
    tag: 1,
    author: "Author 2",
    title: "Image 2",
    id: "2",
  },
  {
    image: "https://via.placeholder.com/1080x1080/eee?text=1:1",
    progress: 2,
    tag: 0,
    author: "Author 3",
    title: "Image 3",
    id: "3",
  },
  {
    image: "https://via.placeholder.com/1080x1080/eee?text=1:1",
    progress: 2,
    tag: 2,
    author: "Author 4",
    title: "Image 4",
    id: "4",
  },
];

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
