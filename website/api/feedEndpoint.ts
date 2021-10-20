import axiosInstance from "common/axiosInstance";
import { ImageDocument } from "./types";

export const get = async (): Promise<ImageDocument[]> => {
  const res = await axiosInstance.get<ImageDocument[]>("/feed");
  return res.data;
};
