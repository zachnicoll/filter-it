import axiosInstance from "common/axiosInstance";
import { UploadResponse } from "./types";

const ENDPOINT = "/upload"

export const get = async (): Promise<UploadResponse> => {
  const res = await axiosInstance.get<UploadResponse>(ENDPOINT);
  return res.data;
};
