import axiosInstance from "common/axiosInstance";
import { Progress, ProgressResponse } from "./types";

const ENDPOINT = "/progress";

export const post = async (id: string): Promise<ProgressResponse> => {
  const res = await axiosInstance.post<ProgressResponse>(
    `${ENDPOINT}`,
    JSON.stringify({ id })
  );
  return res.data;
};
