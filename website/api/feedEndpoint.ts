import axiosInstance from "common/axiosInstance";
import { ImageDocument } from "./types";

const ENDPOINT = "/feed";

export const get = async (): Promise<ImageDocument[]> => {
  const res = await axiosInstance.get<ImageDocument[]>(ENDPOINT);
  return res.data;
};
