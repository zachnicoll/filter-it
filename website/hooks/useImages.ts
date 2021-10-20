import { AxiosError } from "axios";
import { toastError } from "common/toast";
import { ImageDocument } from "api/types";
import { useEffect, useState } from "react";
import API from "api";

interface HookReturn {
  images: ImageDocument[];
  loading: boolean;
}

const mockData = [
  {
    filename: "https://via.placeholder.com/1080x1080/eee?text=1:1",
    progress: 2,
    filters: [0, 1, 2],
    author: "Author 0",
    title: "Image 0",
    id: "0",
  },
  {
    filename: "https://via.placeholder.com/1080x1080/eee?text=1:1",
    progress: 2,
    filters: [2],
    author: "Author 1",
    title: "Image 1",
    id: "1",
  },
  {
    filename: "https://via.placeholder.com/1080x1080/eee?text=1:1",
    progress: 2,
    filters: [1],
    author: "Author 2",
    title: "Image 2",
    id: "2",
  },
  {
    filename: "https://via.placeholder.com/1080x1080/eee?text=1:1",
    progress: 2,
    filters: [0],
    author: "Author 3",
    title: "Image 3",
    id: "3",
  },
  {
    filename: "https://via.placeholder.com/1080x1080/eee?text=1:1",
    progress: 2,
    filters: [2],
    author: "Author 4",
    title: "Image 4",
    id: "4",
  },
];

export const useImages = (): HookReturn => {
  const [images, setImages] = useState<ImageDocument[]>([]);
  const [loading, setLoading] = useState(true);

  const fetch = async (): Promise<void> => {
    setLoading(true);

    try {
      // Replace with real API request when lambda working
      // const data = await API.feed.get();

      const data = mockData;
      setImages(data);
    } catch (e) {
      const err = e as AxiosError;
      toastError((e as AxiosError).message, err);
    }

    setLoading(false);
  };

  useEffect(() => {
    fetch();
  }, []);

  return { images, loading };
};
