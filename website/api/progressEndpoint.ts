import axiosInstance from "common/axiosInstance";
import { Progress } from "./types";

const ENDPOINT = "/progress";

export const get = async (id: string): Promise<Progress> => {
  const res = await axiosInstance.get<Progress>(`${ENDPOINT}?id=${id}`);
  return res.data;
};
