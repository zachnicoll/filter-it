import axiosInstance from "common/axiosInstance";
import { QueueRequestBody, QueueResponse } from "./types";

export const post = async (params: QueueRequestBody): Promise<QueueResponse> => {
  const { title, author, tag, image } = params;

  const res = await axiosInstance.post<QueueResponse>(
    "/queue",
    JSON.stringify({
      title,
      author,
      tag,
      image,
    })
  );

  return res.data;
};
