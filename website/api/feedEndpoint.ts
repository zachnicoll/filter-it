import axiosInstance from "common/axiosInstance";
import { Filter, ImageDocument } from "./types";

const ENDPOINT = "/feed";

export const get = async (filter: Filter | null): Promise<ImageDocument[]> => {
  const res = await axiosInstance.get<ImageDocument[]>(
    filter ? `${ENDPOINT}?filter=${filter}` : ENDPOINT
  );
  return res.data;
};
